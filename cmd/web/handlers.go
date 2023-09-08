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
  data.Form = pasteCreateForm {
    Expires: 365,
  }

  app.render(w, http.StatusOK, "create.gotmpl", data)
}

type pasteCreateForm struct {
  Title string
  Content string
  Expires int
  FieldErrors map[string]string
}

func (app *application) pasteCreatePost(w http.ResponseWriter, r *http.Request) {

  err := r.ParseForm()
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  expires, err := strconv.Atoi(r.PostForm.Get("expires"))
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }
  
  form := pasteCreateForm {
    Title: r.PostForm.Get("title"),
    Content: r.PostForm.Get("content"),
    Expires: expires,
    FieldErrors: map[string]string{},
  }

  app.infoLog.Printf("Creating paste '%s': '%s'. Expires in %d days", form.Title, form.Content, expires)


  if strings.TrimSpace(form.Title) == "" {
    form.FieldErrors["title"] = "This field cannot be blank"
  } else if utf8.RuneCountInString(form.Title) > 100 {
    form.FieldErrors["title"] = "This field cannot be over 100 characters long"
  }

  if strings.TrimSpace(form.Content) == "" {
    form.FieldErrors["content"] = "This field cannot be blank"
  }

  if expires != 1 && expires != 7 && expires != 365 {
    form.FieldErrors["expires"] = "This field must be 1, 7, or 365"
  }

  if len(form.FieldErrors) != 0 {
    data := app.newTemplateData(r)
    data.Form = form
    app.render(w, http.StatusUnprocessableEntity, "create.gotmpl", data)
    return
  }

  id, err := app.pastes.Insert(form.Title, form.Content, expires)
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
