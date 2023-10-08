package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"xyz.haff/pastebox/ui"
)

func (app *application) routes() http.Handler {
  router := mux.NewRouter()

  router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    app.notFound(w)
  })
  
  fileServer := http.FileServer(http.FS(ui.Files))
  router.PathPrefix("/static/").
    Handler(fileServer).
    Methods(http.MethodGet)


  dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

  protected := dynamic.Append(app.requireAuthentication)

  router.Handle("/", dynamic.ThenFunc(app.home)).
          Methods(http.MethodGet)
  router.Handle("/paste/view/{id}", dynamic.ThenFunc(app.pasteView)).
          Methods(http.MethodGet)
  router.Handle("/paste/create", protected.ThenFunc(app.pasteCreate)).
          Methods(http.MethodGet)
  router.Handle("/paste/create", protected.ThenFunc(app.pasteCreatePost)).
          Methods(http.MethodPost)
  
  router.Handle("/user/signup", dynamic.ThenFunc(app.userSignup)).
    Methods(http.MethodGet)
  router.Handle("/user/signup", dynamic.ThenFunc(app.userSignupPost)).
    Methods(http.MethodPost)
  router.Handle("/user/login", dynamic.ThenFunc(app.userLogin)).
    Methods(http.MethodGet)
  router.Handle("/user/login", dynamic.ThenFunc(app.userLoginPost)).
    Methods(http.MethodPost)
  router.Handle("/user/logout", protected.ThenFunc(app.userLogoutPost)).
    Methods(http.MethodPost)

  router.Handle("/ping", http.HandlerFunc(ping)).Methods(http.MethodGet)


  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

  return standard.Then(router)
}
