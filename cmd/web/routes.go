package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.session.Enable, app.recoverPanic, app.logRequest, secureHeaders)

	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/", app.home).Methods("GET")
	router.HandleFunc("/teachers", app.showTeachers).Methods("GET")
	router.HandleFunc("/posts", app.showPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", app.showPost).Methods("GET")
	router.HandleFunc("/admin", app.showDashboard).Methods("GET")
	router.HandleFunc("/admin/post/create", app.createPostForm).Methods("GET")
	router.HandleFunc("/admin/post/create", app.createPost).Methods("POST")
	router.HandleFunc("/admin/posts", app.showAllDashboardPosts).Methods("GET")
	router.HandleFunc("/signup", app.signupUserForm).Methods("GET")
	router.HandleFunc("/signup", app.signupUser).Methods("POST")
	router.HandleFunc("/login", app.loginUserForm).Methods("GET")
	router.HandleFunc("/login", app.loginUser).Methods("POST")
	router.HandleFunc("/logout", app.logoutUser).Methods("POST")

	// Serve static assets
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads/"))))

	return standardMiddleware.Then(router)
}
