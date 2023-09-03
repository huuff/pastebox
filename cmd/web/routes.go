package main

import (
  "net/http"
  "github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
  mux := http.NewServeMux()

  
  fileServer := http.FileServer(nonIndexingFileSystem { http.Dir("./ui/static") })
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("/", app.home)
  mux.HandleFunc("/paste/view", app.pasteView)
  mux.HandleFunc("/paste/create", app.pasteCreate)

  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

  return standard.Then(mux)
}
