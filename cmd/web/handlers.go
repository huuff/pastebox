package main

import (
  "net/http"
  "strconv"
  "html/template"
  "fmt"
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

  id, err := strconv.Atoi(r.URL.Query().Get("id"))
  if err != nil || id < 1 {
    app.notFound(w)
    return
  }

  fmt.Fprintf(w, "Display a specific paste with id %d...", id)
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

  app.pastes.Insert(title, content, expires)

  w.Write([]byte("Create a new paste"))
}
