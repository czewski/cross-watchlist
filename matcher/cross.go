package matcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type CacheItem struct {
	Movies    []string
	ExpiresAt time.Time
}

func CrossWatchlists(users []string) []string {
	profiles := getListsConcurrent(users)
	if len(profiles) == 0 {
		return nil
	}
	return intersectAll(profiles, users)
}

func intersectAll(profiles map[string][]string, users []string) []string {
	if len(users) == 0 {
		return nil
	}

	// Start with smallest slice to minimize comparisons
	smallestIndex := 0
	for i := 1; i < len(users); i++ {
		if len(profiles[users[i]]) < len(profiles[users[smallestIndex]]) {
			smallestIndex = i
		}
	}

	result := profiles[users[smallestIndex]]
	for i, user := range users {
		if i == smallestIndex {
			continue
		}
		result = hashIntersect(result, profiles[user])
		if len(result) == 0 {
			break
		}
	}
	return result
}

func getListsConcurrent(users []string) map[string][]string {
	//var wg sync.WaitGroup
	mu := sync.Mutex{}
	profiles := make(map[string][]string, len(users))

	g, ctx := errgroup.WithContext(context.Background())

	for _, user := range users {
		user := user // Capture loop variable
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				movies, err := getCachedUserMovies(user)
				if err != nil {
					return err
				}

				mu.Lock()
				profiles[user] = movies
				mu.Unlock()
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return nil
	}

	return profiles
}

func getCachedUserMovies(user string) ([]string, error) {
	if cached, ok := cache.Load(user); ok {
		item := cached.(CacheItem)
		if time.Now().Before(item.ExpiresAt) {
			return item.Movies, nil
		}
	}

	movies, err := movieSlugsFromUser(user)
	if err != nil {
		return nil, err
	}

	cache.Store(user, CacheItem{
		Movies:    movies,
		ExpiresAt: time.Now().Add(cacheTTL),
	})

	return movies, nil
}

func movieSlugsFromUser(user string) ([]string, error) {
	// first page to determine total pages
	url := fmt.Sprintf("https://letterboxd.com/%s/watchlist", user)
	body, err := fetchMoviesFromPage(url, 1)
	if err != nil {
		return nil, err
	}

	totalPages := parseTotalPages(body)
	if totalPages == 0 {
		return nil, fmt.Errorf("no pages found")
	}

	// all pages concurrently
	var mu sync.Mutex
	var movies []string
	g, ctx := errgroup.WithContext(context.Background())

	for page := 1; page <= totalPages; page++ {
		page := page
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				url := fmt.Sprintf("https://letterboxd.com/%s/watchlist", user)
				body, err := fetchMoviesFromPage(url, page)
				if err != nil {
					return err
				}

				pageMovies := parseMovies(body)
				mu.Lock()
				movies = append(movies, pageMovies...)
				mu.Unlock()
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return unique(movies), nil
}

// func fetchWatchlistPage(user string, page int) (string, error) {
// 	url := fmt.Sprintf("https://letterboxd.com/%s/watchlist/page/%d/", user, page)
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("User-Agent", userAgent)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to fetch page: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}

// 	body := &strings.Builder{}
// 	_, err = io.Copy(body, resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read body: %w", err)
// 	}

// 	return body.String(), nil
// }

func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range slice {
		if _, value := keys[item]; !value {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}
