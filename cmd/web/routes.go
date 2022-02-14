package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/prepodavateli", app.showTeachers)
	mux.HandleFunc("/post", app.showPost)
	mux.HandleFunc("/post/create", app.createPost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	return mux
}