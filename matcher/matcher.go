package matcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

func Do(users []string) []string {
	profiles := getLists(users)
	matches := matchLists(profiles, users)
	return matches
}

func matchLists(movieSlugs map[string][]string, users []string) []string {
	countUsers := len(users)
	var matches []string

	for i := 0; i < countUsers; i++ {
		if i == 0 {
			matches = movieSlugs[users[i]]
		} else {
			matches = hashGeneric(matches, movieSlugs[users[i]])
		}
	}

	return matches
}

// FROM: https://github.com/juliangruber/go-intersect/blob/master/intersect.go
// Hash has complexity: O(n * x) where x is a factor of hash function efficiency (between 1 and 2)
func hashGeneric[T comparable](a []T, b []T) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := hash[v]; ok {
			set = append(set, v)
		}
	}

	return set
}

func getLists(users []string) map[string][]string {
	usersWatchlist := make(map[string][]string, len(users))
	for _, user := range users {
		usersWatchlist[user] = movieSlugsFromUser(user)
	}

	return usersWatchlist
}

func movieSlugsFromUser(user string) (movieIds []string) {
	//Request
	resp, err := http.Get("https://letterboxd.com/" + user + "/watchlist/page/1/")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Get number of pages
	pattern := `class="paginate-page"><a href="[^"]+">(\d+)</a>`
	re := regexp.MustCompile(pattern)
	match := re.FindAllStringSubmatch(string(body), -1)
	pages := match[2][1]
	numberOfPages, _ := strconv.Atoi(pages)

	// scrap page + go to next
	for i := 1; i <= numberOfPages+1; i++ {
		resp, err := http.Get("https://letterboxd.com/" + user + "/watchlist/page/" + strconv.Itoa(i) + "/")
		if err != nil {
			fmt.Println(err)
			return nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		//pattern = `data-film-id="(\d+)"`
		pattern = `data-film-slug="([a-z-]+)"`
		re = regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(string(body), -1)

		if len(matches) > 0 {
			for _, match := range matches {
				if len(match) > 1 {
					movieIds = append(movieIds, match[1])
				}
			}
		}
	}

	return movieIds
}
