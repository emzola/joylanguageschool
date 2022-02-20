package main

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"joylanguageschool.ru/pkg/models"
	"joylanguageschool.ru/pkg/validator"
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

	image, _ := app.FileUpload(r)
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// Validate form fields
	errors := validator.CreatePost(title, content, image)

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
		return
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
	app.render(w, r, "signup.page.tmpl", nil)
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// Validate form fields
	errors := validator.Signup(name, email, password)

	// If there are errors, re-populate the registration form and display it again
	if len(errors) > 0 {
		app.render(w, r, "signup.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	err = app.users.Insert(name, email, password)

	if err == models.ErrDuplicateEmail {
		errors["email"] = "Email address is already in use"
		app.render(w, r, "signup.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup successful. Please log in.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", nil)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	id, err := app.users.Authenticate(email, password)
	if err == models.ErrInvalidCredentials {
		errors := validator.Login() 
		errors["generic"] = "Email or password is incorrect"
		app.render(w, r, "login.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
