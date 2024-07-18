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
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/submit-form", handleSubmitForm)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("new.html"))
	tmpl.Execute(w, nil)
}

func handleSubmitForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var inputs []string
	for _, values := range r.Form {
		inputs = append(inputs, values...)
	}
	matches := matcher.Do(inputs)
	response := strings.Join(matches, "\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
