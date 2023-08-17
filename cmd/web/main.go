package main

import (
	"log"
	"net/http"
  "os"
)


type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
}

func main() {
  args := ParseArgs()

  mux := http.NewServeMux()

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  app := application {
    infoLog: infoLog,
    errorLog: errorLog,
  }
  
  fileServer := http.FileServer(nonIndexingFileSystem { http.Dir("./ui/static") })

  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("/", app.home)
  mux.HandleFunc("/paste/view", app.pasteView)
  mux.HandleFunc("/paste/create", app.pasteCreate)

  srv := &http.Server {
    Addr: args.Addr(),
    ErrorLog: errorLog,
    Handler: mux,
  }

  app.infoLog.Printf("Starting server on %s", args.Addr())
  if err := srv.ListenAndServe(); err != nil {
    app.errorLog.Fatal(err)
  }

}
