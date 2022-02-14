package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/home.layout.tmpl",
		"./ui/html/contact.partial.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	app.parseTemplateFiles(w, files)

}

func (app *application) showTeachers(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/prepodavateli.page.tmpl",
		"./ui/html/page.layout.tmpl",
		"./ui/html/pageheader.partial.tmpl",
		"./ui/html/contact.partial.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	app.parseTemplateFiles(w, files)
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Show a specific post with ID %d", id)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new post"))
}
