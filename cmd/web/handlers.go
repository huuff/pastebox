package main

import (
	"errors"
	"fmt"
  "html/template"
	"net/http"

	"xyz.haff/pastebox/internal/db"
)

// TODO: Just printing out every paste for testing, but bring back the template in the future
func (app *application) home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    app.notFound(w)
    return
  }

  pastes, err := app.pastes.Latest()
  if err != nil {
    app.serverError(w, err)
    return
  }

  for _, paste := range pastes {
    fmt.Fprintf(w, "%+v\n", paste)
  }

  //files := []string {
    //"./ui/html/base.gotmpl",
    //"./ui/html/partials/nav.gotmpl",
    //"./ui/html/pages/home.gotmpl",
  //}
  //ts, err := template.ParseFiles(files...)
  //if err != nil {
    //app.serverError(w, err)
    //return
  //}

  //err = ts.ExecuteTemplate(w, "base", nil)
  //if err != nil {
    //app.serverError(w, err)
  //}
}

// TODO: Give some styles to this!!
func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {
  id := r.URL.Query().Get("id")
  if id == "" {
    app.notFound(w)
    return
  }

  paste, err := app.pastes.Get(id)
  if err != nil {
    if errors.Is(err, db.ErrPasteNotFound) {
      app.notFound(w)
    } else {
      app.serverError(w, err)
    }
    return
  }

  files := []string {
    "./ui/html/base.gotmpl",
    "./ui/html/partials/nav.gotmpl",
    "./ui/html/pages/view.gotmpl",
  }

  ts, err := template.ParseFiles(files...)
  if err != nil {
    app.serverError(w, err)
    return
  }

  err = ts.ExecuteTemplate(w, "base", paste)
  if err != nil {
    app.serverError(w, err)
  }
}

func (app *application) pasteCreate(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", http.MethodPost)
    app.clientError(w, http.StatusMethodNotAllowed)
    return
  }

  // TODO: Actually receive params from request
  title := "O snail"
  content := `
O snail
Climb Mount Fuji
But slowly, slowly!

- Kobayashi Issa`
  expires := 7

  id, err := app.pastes.Insert(title, content, expires)
  if err != nil {
    app.serverError(w, err)
    return
  }

  http.Redirect(w, r, fmt.Sprintf("/paste/view?id=%s", id), http.StatusSeeOther)
}
