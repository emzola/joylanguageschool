package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/lithammer/shortuuid/v4"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.logger.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()

	// Add flash message to template data is one exists
	td.Flash = app.session.PopString(r, "flash")

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

// Get image from form and save to uploads directory
func (app *application) FileUpload(r *http.Request) (string, error) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		return "", err
	}

	defer file.Close()

	name := handler.Filename + shortuuid.New() + ".jpg"

	f, err := os.OpenFile("uploads/"+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	defer f.Close()

	io.Copy(f, file)

	return name, nil
}
