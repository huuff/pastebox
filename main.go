package main

import (
  "log"
  "net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  w.Write([]byte("Hello from Pastebox"))
}

func pasteView(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Display a specific paste"))
}

func pasteCreate(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.WriteHeader(http.StatusMethodNotAllowed)
    w.Write([]byte("Method Not Allowed"))
    return
  }

  w.Write([]byte("Create a new paste"))
}

func main() {
  mux := http.NewServeMux()

  mux.HandleFunc("/", home)
  mux.HandleFunc("/paste/view", pasteView)
  mux.HandleFunc("/paste/create", pasteCreate)

  log.Println("Starting server on :4000")
  err := http.ListenAndServe(":4000", mux)
  log.Fatal(err)

}
