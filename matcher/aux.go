package matcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
	pageRe    = regexp.MustCompile(`class="paginate-page"><a href="[^"]+">(\d+)</a>`)
	movieRe   = regexp.MustCompile(`data-film-slug="([a-z0-9-]+)"`)
	cache     = sync.Map{}
	cacheTTL  = 5 * time.Minute
	userAgent = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)" //"Mozilla/5.0 (compatible; WatchlistMatcher/1.0; +https://github.com/your/repo)"

)

func getMoviesFromList(link string) ([]string, error) {
	// Get first page to determine total pages
	body, err := fetchMoviesFromPage(link, 1)
	if err != nil {
		return nil, err
	}
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	_, err = f.WriteString(body)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return nil, nil
	}

	totalPages := parseTotalPages(body)
	if totalPages == 0 {
		return nil, fmt.Errorf("no pages found")
	}

	// Collect all pages concurrently
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
				body, err := fetchMoviesFromPage(link, page)
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

func fetchMoviesFromPage(initialLink string, page int) (string, error) {
	url := fmt.Sprintf("%s/page/%d/", initialLink, page)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body := &strings.Builder{}
	_, err = io.Copy(body, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return body.String(), nil
}

func parseTotalPages(body string) int {
	matches := pageRe.FindAllStringSubmatch(body, -1)
	if len(matches) < 1 {
		return 1
	}
	lastPage, _ := strconv.Atoi(matches[len(matches)-1][1])
	return lastPage
}

func parseMovies(body string) []string {
	matches := movieRe.FindAllStringSubmatch(body, -1)
	movies := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			movies = append(movies, match[1])
		}
	}
	return movies
}

func hashIntersect(a, b []string) []string {
	set := make(map[string]struct{}, len(a))
	result := make([]string, 0, min(len(a), len(b)))

	for _, v := range a {
		set[v] = struct{}{}
	}

	for _, v := range b {
		if _, exists := set[v]; exists {
			result = append(result, v)
		}
	}
	return result
}
