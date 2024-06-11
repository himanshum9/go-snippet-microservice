package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/himanshum9/go-snippet/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Because httprouter matches the "/" path exactly, we can now remove the manual check of r.URL.Path != "/" from this handler.
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }
	snippets, err := app.snippetModel.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year), and add the
	// snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

}

// Change the signature of the snippetView handler so it is defined as a method
// against *application.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		// http.NotFound(w, r)
		return
	}
	s, err := app.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// And do the same thing again here...
	data := app.newTemplateData(r)
	data.Snippet = s

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.tmpl", data)
}

// Rename this handler to snippetCreatePost.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// First we call r.ParseForm() which adds any data in POST request bodies
	// to the r.PostForm map. This also works in the same way for PUT and PATCH
	// requests. If there are any errors, we use our app.ClientError() helper to
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content
	// from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an integer. So we need to manually covert
	// the form data to an integer using strconv.Atoi(), and we send a 400 Bad
	// Request response if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	// http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

	// Update the redirect path to use the new clean URL format.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
