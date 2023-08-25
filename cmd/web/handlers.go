package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    app.notFound(w)
    return
  }

  files := []string {
    "./ui/html/base.gotmpl",
    "./ui/html/partials/nav.gotmpl",
    "./ui/html/pages/home.gotmpl",
  }
  ts, err := template.ParseFiles(files...)
  if err != nil {
    app.serverError(w, err)
    return
  }

  err = ts.ExecuteTemplate(w, "base", nil)
  if err != nil {
    app.serverError(w, err)
  }
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {
  id := r.URL.Query().Get("id")
  if id == "" {
    app.notFound(w)
    return
  }

  // TODO: This is coupled to the database's error for not found! Try to wrap it
  paste, err := app.pastes.Get(id)
  if err != nil {
    if errors.Is(err, mongo.ErrNoDocuments) {
      app.notFound(w)
    } else {
      app.serverError(w, err)
    }
    return
  }

  fmt.Fprintf(w, "%+v", paste)
}

func (app *application) pasteCreate(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", http.MethodPost)
    app.clientError(w, http.StatusMethodNotAllowed)
    return
  }

  // TODO: Actually receive params from request
  title := "O snail"
  content := `O snail
Climb Mount Fuji
But slowly, slowly!

Kobayashi Issa`
  expires := 7

  id, err := app.pastes.Insert(title, content, expires)
  if err != nil {
    app.serverError(w, err)
    return
  }

  http.Redirect(w, r, fmt.Sprintf("/paste/view?id=%s", id), http.StatusSeeOther)
}
