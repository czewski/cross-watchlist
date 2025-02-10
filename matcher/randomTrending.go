package matcher

import "fmt"

func RandomFromTrending() ([]string, error) {
	// fetch from api response, no need to create another func because the pagination is different here, so it wont have multiple pages parsed
	url := "https://letterboxd.com/films/ajax/popular/this/week/?esiAllowFilters=true"
	var movies []string

	moviesResp, err := getMoviesFromList(url)
	if err != nil {
		fmt.Printf("error getting movies from list: %v\n", err)
		return movies, err
	}
	movies = append(movies, moviesResp...)

	return movies, nil
}
