package main

import (
	"log"
	"net/http"
  "os"
)


func main() {
  args := ParseArgs()

  mux := http.NewServeMux()

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
  
  fileServer := http.FileServer(nonIndexingFileSystem { http.Dir("./ui/static") })

  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("/", home)
  mux.HandleFunc("/paste/view", pasteView)
  mux.HandleFunc("/paste/create", pasteCreate)

  srv := &http.Server {
    Addr: args.Addr(),
    ErrorLog: errorLog,
    Handler: mux,
  }

  infoLog.Printf("Starting server on %s", args.Addr())
  if err := srv.ListenAndServe(); err != nil {
    errorLog.Fatal(err)
  }

}
