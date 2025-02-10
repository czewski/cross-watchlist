package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"watchlist/matcher"
)

type FormData struct {
	Inputs []string
}

func main() {
	// x := []string{"https://letterboxd.com/dave/list/official-top-250-narrative-feature-films"}
	// matcher.RandomFromLists(x)

	_, _ = matcher.RandomFromTrending()

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/random-from-watchlists", handleCrossWatchlists) //submit-form
	http.HandleFunc("/random-from-lists", handleRandomFromLists)
	http.HandleFunc("/random-from-trending", handleRandomFromTrending)

	//TODO: Apply random to list and return a favorite one
	//TODO: add other 3 funcs, rename the submit-form (this changes the FE as well)
	// showRandomFromList() (ok), maybe add cache
	// showRandomFromTrending()
	// showSuggestionsBasedOnProfile()
	// showCrossWatchlist() //Maybe the user can add preferences, length, genre, etc
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("new.html"))
	tmpl.Execute(w, nil)
}

func handleCrossWatchlists(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var inputs []string
	for _, values := range r.Form {
		inputs = append(inputs, values...)
	}
	// todo: handle error from those calls
	matches := matcher.CrossWatchlists(inputs)
	response := strings.Join(matches, "\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func handleRandomFromLists(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var inputs []string
	for _, values := range r.Form {
		inputs = append(inputs, values...)
	}
	matches, err := matcher.RandomFromLists(inputs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := strings.Join(matches, "\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func handleRandomFromTrending(w http.ResponseWriter, r *http.Request) {
	matches, err := matcher.RandomFromTrending()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := strings.Join(matches, "\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
