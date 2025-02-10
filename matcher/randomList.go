package matcher

import (
	"fmt"
)

func RandomFromLists(lists []string) ([]string, error) {
	var movies []string

	// TODO: retrieve from cache?

	for _, list := range lists {
		moviesResp, err := getMoviesFromList(list)
		if err != nil {
			fmt.Printf("error getting movies from list: %v\n", err)
			return movies, err
		}
		movies = append(movies, moviesResp...)
	}

	return movies, nil
}
