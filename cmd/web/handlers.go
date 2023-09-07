package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
  "unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/samber/lo"

	"xyz.haff/pastebox/internal/db"
	"xyz.haff/pastebox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
  pastes, err := app.pastes.Latest()
  if err != nil {
    app.serverError(w, err)
    return
  }
  
  data := app.newTemplateData(r)
  data.Pastes = lo.Map(pastes, func(paste models.Paste, _ int) *models.Paste {
      return &paste
  })

  app.render(w, http.StatusOK, "home.gotmpl", data)
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]
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

  data := app.newTemplateData(r)
  data.Paste = paste

  app.render(w, http.StatusOK, "view.gotmpl", data)
}

func (app *application) pasteCreate(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)

  app.render(w, http.StatusOK, "create.gotmpl", data)
}

func (app *application) pasteCreatePost(w http.ResponseWriter, r *http.Request) {

  err := r.ParseForm()
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  title := r.PostForm.Get("title")
  content := r.PostForm.Get("content")
  expires, err := strconv.Atoi(r.PostForm.Get("expires"))
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  app.infoLog.Printf("Creating paste '%s': '%s'. Expires in %d days", title, content, expires)

  fieldErrors := make(map[string]string)

  if strings.TrimSpace(title) == "" {
    fieldErrors["title"] = "This field cannot be blank"
  } else if utf8.RuneCountInString(title) > 100 {
    fieldErrors["title"] = "This field cannot be over 100 characters long"
  }

  if strings.TrimSpace(content) == "" {
    fieldErrors["content"] = "This field cannot be blank"
  }

  if expires != 1 && expires != 7 && expires != 365 {
    fieldErrors["expires"] = "This field must be 1, 7, or 365"
  }

  if len(fieldErrors) != 0 {
    fmt.Fprint(w, fieldErrors)
    return
  }

  id, err := app.pastes.Insert(title, content, expires)
  if err != nil {
    app.serverError(w, err)
    return
  }

  http.Redirect(w, r, fmt.Sprintf("/paste/view/%s", id), http.StatusSeeOther)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
  ts, ok := app.templateCache[page]
  if !ok {
    err := fmt.Errorf("the template %s does not exist", page)
    app.serverError(w, err)
    return
  }

  buf := new(bytes.Buffer)
  
  if err := ts.ExecuteTemplate(buf, "base", data); err != nil {
    app.serverError(w, err)
    return
  }


  w.WriteHeader(status)

  buf.WriteTo(w)
}
