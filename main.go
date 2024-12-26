package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)
type Page struct {
    Title string
    Body  []byte
}
func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, _ := template.ParseFiles(tmpl + ".html")
    t.Execute(w, p)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/games/"):]

	// construct filepath
	filepath := "assets/games/" + title + ".html"
	// Make sure we sanitize the string to prevent some bs like ../../passwords
	filepath = path.Clean(filepath)

	// Check if the file exists
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        http.NotFound(w, r)
        return
    }
	// Tell the webpage to expect html
	w.Header().Set("Content-Type", "text/html")
	// actually serve the file
	http.ServeFile(w, r, filepath)
}


func main() {
	http.HandleFunc("/games/", gameHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
