package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) newTemplateData(r *http.Request) *templateData {
  return &templateData {
    CurrentYear: time.Now().Year(),
    Flash: app.sessionManager.PopString(r.Context(), "flash"),
    IsAuthenticated: app.isAuthenticated(r),
    CSRFToken: nosurf.Token(r),
  }
}

func (app *application) serverError(w http.ResponseWriter, err error) {
  trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
  app.errorLog.Output(2, trace)

  http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
  http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
  app.clientError(w, http.StatusNotFound)
}

func (app *application) isAuthenticated(r *http.Request) bool {
  return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
