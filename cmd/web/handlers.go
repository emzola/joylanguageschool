package main

import (
	"net/http"

	"joylanguageschool.ru/pkg/mail"
	"joylanguageschool.ru/pkg/models"
	"joylanguageschool.ru/pkg/validator"
)

func (app *application) NotFound(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "404.page.tmpl", nil)
}

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
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	post, err := app.posts.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Get last 5 posts for sidebar
	posts, err := app.posts.GetLastFivePosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "post.page.tmpl", &templateData{
		Post: post,
		Posts: posts,
	})
}

func (app *application) showPosts(w http.ResponseWriter, r *http.Request) {
	p, err := app.posts.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "posts.page.tmpl", &templateData{Posts: p})
}

func (app *application) showProgramme(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	programme, err := app.programmes.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "programme.page.tmpl", &templateData{
		Programme: programme,
	})
}


// Process email fron contact form
func(app *application) sendMail(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	subject := r.PostForm.Get("subject")
	message := r.PostForm.Get("message")

	errors := validator.ContactForm(name, email, subject, message)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(errors) > 0 {
		app.render(w, r, "contact.partial.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	err = mail.SendEmail(name, email, subject, message) 
	if err != nil {
		app.serverError(w, err)
	}

	app.session.Put(r, "flash", "Сообщение успешно отправлено")
	http.Redirect(w, r, "/#контакты", http.StatusSeeOther)
}

func (app *application) showDashboard(w http.ResponseWriter, r *http.Request) {
	p, err := app.posts.GetLastFivePosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "dashboard.page.tmpl", &templateData{Posts: p})
}

// Dashboard - Posts handlers

// Display posts on the dashboard
func (app *application) showAllDashboardPosts(w http.ResponseWriter, r *http.Request) {
	p, err := app.posts.GetAllPosts()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "dashboardposts.page.tmpl", &templateData{Posts: p})
}

// Display create post form
func (app *application) createPostForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createpost.page.tmpl", nil)
}

// Create post on the dashboard
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

// Display edit post form
func (app *application) editPostForm(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	post, err := app.posts.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "editpost.page.tmpl", &templateData{
		Post:   post,
	})
}

// Edit Post
func(app *application) editPost(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	image, _ := app.FileUpload(r)
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// Get post data to use to re-populate form if re-displayed due to error
	post, err := app.posts.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Validate form fields
	errors := validator.EditPost(title, content, image)

	// If there are errors, re-populate the post form and display it again
	if len(errors) > 0 {
		app.render(w, r, "editpost.page.tmpl", &templateData{
			FormErrors: errors,
			Post:   post,
		})
		return
	}

	// Check if a new image file is uploaded.
	// If not uploaded, set it's value to value already in the database
	if image == "" {
		image = post.Image
	}

	// Update post
	err = app.posts.Update(id, title, content, image)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Post successfully updated")
	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}

// Delete post
func (app *application) deletePost (w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	err = app.posts.Delete(id)
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}

	app.session.Put(r, "flash", "Post successfully deleted")
	http.Redirect(w, r, "/admin/posts", http.StatusSeeOther)
}


// Dashboard - Programmes handlers

// Display programmes on the dashboard
func (app *application) showAllDashboardProgrammes(w http.ResponseWriter, r *http.Request) {
	p, err := app.programmes.GetAllProgrammes()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "dashboardprogrammes.page.tmpl", &templateData{Programmes: p})
}

// Display create programme form
func (app *application) createProgrammeForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createprogramme.page.tmpl", nil)
}

// Create programme
func (app *application) createProgramme(w http.ResponseWriter, r *http.Request) {
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
		app.render(w, r, "createprogramme.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	err = app.programmes.Insert(title, content, image)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Programme successfully created")
	http.Redirect(w, r, "/admin/programmes", http.StatusSeeOther)
}

// Display edit post form
func (app *application) editProgrammeForm(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	programme, err := app.programmes.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "editprogramme.page.tmpl", &templateData{
		Programme:   programme,
	})
}

// Edit Post
func(app *application) editProgramme(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	image, _ := app.FileUpload(r)
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// Get post data to use to re-populate form if re-displayed due to error
	programme, err := app.programmes.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Validate form fields
	errors := validator.EditPost(title, content, image)

	// If there are errors, re-populate the post form and display it again
	if len(errors) > 0 {
		app.render(w, r, "editprogramme.page.tmpl", &templateData{
			FormErrors: errors,
			Programme:   programme,
		})
		return
	}

	// Check if a new image file is uploaded.
	// If not uploaded, set it's value to value already in the database
	if image == "" {
		image = programme.Image
	}

	// Update post
	err = app.programmes.Update(id, title, content, image)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Programme successfully updated")
	http.Redirect(w, r, "/admin/programmes", http.StatusSeeOther)
}

// Delete post
func (app *application) deleteProgramme (w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	err = app.programmes.Delete(id)
	if err == models.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.clientError(w, http.StatusInternalServerError)
		return
	}

	app.session.Put(r, "flash", "Programme successfully deleted")
	http.Redirect(w, r, "/admin/programmes", http.StatusSeeOther)
}


// User signup and login handlers
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
