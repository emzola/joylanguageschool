package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.session.Enable, app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	router := httprouter.New()
	// Routes
	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/teachers", dynamicMiddleware.ThenFunc(app.showTeachers))
	router.Handler(http.MethodGet, "/posts", dynamicMiddleware.ThenFunc(app.showPosts))
	router.Handler(http.MethodGet, "/posts/:id", dynamicMiddleware.ThenFunc(app.showPost))
	router.Handler(http.MethodGet, "/admin", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showDashboard))
	router.Handler(http.MethodGet, "/admin/post/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createPostForm))
	router.Handler(http.MethodPost, "/admin/post/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createPost))
	router.Handler(http.MethodGet, "/admin/posts", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showAllDashboardPosts))
	router.Handler(http.MethodGet, "/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	router.Handler(http.MethodPost, "/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	router.Handler(http.MethodGet, "/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	router.Handler(http.MethodPost, "/login", dynamicMiddleware.ThenFunc(app.loginUser))
	router.Handler(http.MethodPost, "/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser))

	// Serve static assets
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))))
	// router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads/"))))

	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))
	router.ServeFiles("/uploads/*filepath", http.Dir("./uploads"))

	return standardMiddleware.Then(router)
}
