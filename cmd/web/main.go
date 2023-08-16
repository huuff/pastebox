package main

import (
	"log"
	"net/http"
)


func main() {
  args := ParseArgs()

  mux := http.NewServeMux()
  
  fileServer := http.FileServer(nonIndexingFileSystem { http.Dir("./ui/static") })

  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  mux.HandleFunc("/", home)
  mux.HandleFunc("/paste/view", pasteView)
  mux.HandleFunc("/paste/create", pasteCreate)

  log.Printf("Starting server on %s", args.Addr())
  err := http.ListenAndServe(args.Addr(), mux)
  log.Fatal(err)

}
