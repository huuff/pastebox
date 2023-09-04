package main

import (
  "net/http"
  "github.com/justinas/alice"
  "github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
  router := mux.NewRouter()

  
  fileServer := http.FileServer(nonIndexingFileSystem { http.Dir("./ui/static") })
  router.PathPrefix("/static").
        Handler(http.StripPrefix("/static", fileServer))

  router.HandleFunc("/", app.home).
          Methods(http.MethodGet)
  router.HandleFunc("/paste/view", app.pasteView).
          Methods(http.MethodGet)
  router.HandleFunc("/paste/create", app.pasteCreate).
          Methods(http.MethodGet)
  router.HandleFunc("/paste/create", app.pasteCreatePost).
          Methods(http.MethodPost)

  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

  return standard.Then(router)
}
