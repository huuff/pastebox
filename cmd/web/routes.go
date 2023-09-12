package main

import (
  "net/http"
  "github.com/justinas/alice"
  "github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
  router := mux.NewRouter()

  router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    app.notFound(w)
  })
  
  fileServer := http.FileServer(nonIndexingFileSystem { http.Dir("./ui/static") })
  router.PathPrefix("/static").
        Handler(http.StripPrefix("/static", fileServer))

  dynamic := alice.New(app.sessionManager.LoadAndSave)

  router.Handle("/", dynamic.ThenFunc(app.home)).
          Methods(http.MethodGet)
  router.Handle("/paste/view/{id}", dynamic.ThenFunc(app.pasteView)).
          Methods(http.MethodGet)
  router.Handle("/paste/create", dynamic.ThenFunc(app.pasteCreate)).
          Methods(http.MethodGet)
  router.Handle("/paste/create", dynamic.ThenFunc(app.pasteCreatePost)).
          Methods(http.MethodPost)

  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

  return standard.Then(router)
}
