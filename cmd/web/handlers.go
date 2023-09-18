package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/gookit/validate"
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
    if errors.Is(err, db.ErrRecordNotFound) {
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
  Title string `validate:"required|max_len:100" form:"title"`
  Content string `validate:"required" form:"content"`
  Expires int `validate:"in:1,7,365" form:"expires"`
  FieldErrors map[string]string `validate:"-" form:"-"`
}

func (app *application) pasteCreatePost(w http.ResponseWriter, r *http.Request) {
  var form pasteCreateForm

  if err := app.decodePostForm(r, &form); err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  app.infoLog.Printf("Creating paste '%s': '%s'. Expires in %d days", form.Title, form.Content, form.Expires)


  validation := validate.Struct(form)

  if !validation.Validate() {
    form.FieldErrors = errorsToMap(validation.Errors)
  }

  if len(form.FieldErrors) != 0 {
    data := app.newTemplateData(r)
    data.Form = form
    app.render(w, http.StatusUnprocessableEntity, "create.gotmpl", data)
    return
  }

  id, err := app.pastes.Insert(form.Title, form.Content, form.Expires)
  if err != nil {
    app.serverError(w, err)
    return
  }

  app.sessionManager.Put(r.Context(), "flash", "Paste successfully created!")

  http.Redirect(w, r, fmt.Sprintf("/paste/view/%s", id), http.StatusSeeOther)
}

type userSignupForm struct {
  Name string `validate:"required" form:"name"`
  Email string `validate:"required|email" form:"email"`
  Password string `validate:"required|min_len:8" form:"password"`
  FieldErrors map[string]string `validate:"-" form:"-"`
}

// TODO: Please make it pretty
func (app *application) userSignup(w http.ResponseWriter, r* http.Request) {
  data := app.newTemplateData(r)
  data.Form = userSignupForm {}
  app.render(w, http.StatusOK, "signup.gotmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
  var form userSignupForm

  err := app.decodePostForm(r, &form)
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  validation := validate.Struct(form)
  
  if !validation.Validate() {
    form.FieldErrors = errorsToMap(validation.Errors)
  }


  if len(form.FieldErrors) != 0 {
    data := app.newTemplateData(r)
    data.Form = form
    app.render(w, http.StatusUnprocessableEntity, "signup.gotmpl", data)
    return
  }

  // TODO: Finish it!
  fmt.Println(w, "Create a new user")
  
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Display an HTML form for logging-in a user")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Authenticate and login the user")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "Logout the user...")
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

func (app *application) decodePostForm(r *http.Request, dst any) error {
  if err := r.ParseForm(); err != nil {
    return err
  }

  if err := app.formDecoder.Decode(dst, r.PostForm); err != nil {
    var invalidDecoderError *form.InvalidDecoderError
    if errors.As(err, &invalidDecoderError) {
      panic(err)
    }

    return err
  }

  return nil
}
