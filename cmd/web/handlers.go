package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"joylanguageschool.ru/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	p, err := app.posts.GetLastThreePosts()
	if err != nil {
		app.serverError(w, err)
		return 
	}

	app.render(w, r, "home.page.tmpl", &templateData{Posts: p})
}

func (app *application) showTeachers(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "teachers.page.tmpl", nil)
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	p, err := app.posts.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "post.page.tmpl", &templateData{Post: p})
}

func (app *application) showDashboard(w http.ResponseWriter, r *http.Request) {
	p, err := app.posts.GetLastTenPosts()
	if err != nil {
		app.serverError(w, err)
		return 
	}
	app.render(w, r, "dashboard.page.tmpl", &templateData{Posts: p})
}

func (app *application) showAllDashboardPosts(w http.ResponseWriter, r *http.Request) {
	p, err := app.posts.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return 
	}
	app.render(w, r, "dashboardposts.page.tmpl", &templateData{Posts: p})
}

func (app *application) createPostForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createpost.page.tmpl", nil)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	image, err := app.FileUpload(r)
	if err != nil {
		app.badRequest(w, err)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	err = app.posts.Insert(title, content, image)
	if err != nil {
		app.serverError(w, err)
	}
	
	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

func (app *application) showPosts(w http.ResponseWriter, r *http.Request) {
	p, err := app.posts.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return 
	}

	app.render(w, r, "posts.page.tmpl", &templateData{Posts: p})
}
