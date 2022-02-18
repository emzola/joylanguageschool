package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/mux"
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
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
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

	image, _ := app.FileUpload(r)
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if strings.TrimSpace(image) == "" {
		errors["image"] = "This field cannot be blank"
	}

	// If there are errors, re-populate the post form and display it again
	if len(errors) > 0 {
		app.render(w, r, "createpost.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	err = app.posts.Insert(title, content, image)
	if err != nil {
		app.serverError(w, err)
	}

	app.session.Put(r, "flash", "Post successfully created")
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

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user signup form...")
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display the user login form...")
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
