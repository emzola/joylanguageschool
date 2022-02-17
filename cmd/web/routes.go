package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/teachers", app.showTeachers)
	router.HandlerFunc(http.MethodGet, "/posts/:id", app.showPost)
	router.HandlerFunc(http.MethodGet, "/posts", app.showPosts)
	router.HandlerFunc(http.MethodGet, "/admin", app.showDashboard)
	router.HandlerFunc(http.MethodGet, "/admin/post/create", app.createPostForm)
	router.HandlerFunc(http.MethodPost, "/admin/post/create", app.createPost)
	router.HandlerFunc(http.MethodGet, "/admin/posts", app.showAllDashboardPosts)

	// fileServer := http.FileServer(http.Dir("./ui/static"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))
	router.ServeFiles("/uploads/*filepath", http.Dir("./uploads"))
	
	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}

