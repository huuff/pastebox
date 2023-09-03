package main

import "net/http"

func (app *application) routes() http.Handler {
  mux := http.NewServeMux()

  
  fileServer := http.FileServer(nonIndexingFileSystem { http.Dir("./ui/static") })
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("/", app.home)
  mux.HandleFunc("/paste/view", app.pasteView)
  mux.HandleFunc("/paste/create", app.pasteCreate)

  return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
