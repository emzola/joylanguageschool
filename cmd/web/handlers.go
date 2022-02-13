package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/home.layout.tmpl",
		"./ui/html/contact.partial.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	parseTemplateFiles(w, files)

}

func showTeachers(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/prepodavateli.page.tmpl",
		"./ui/html/page.layout.tmpl",
		"./ui/html/pageheader.partial.tmpl",
		"./ui/html/contact.partial.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	parseTemplateFiles(w, files)
}

func showPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Show a specific post with ID %d", id)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new post"))
}
