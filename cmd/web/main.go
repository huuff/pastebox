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

  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  app := application {
    infoLog: infoLog,
    errorLog: errorLog,
  }

  srv := &http.Server {
    Addr: args.Addr(),
    ErrorLog: errorLog,
    Handler: app.routes(),
  }

  app.infoLog.Printf("Starting server on %s", args.Addr())
  if err := srv.ListenAndServe(); err != nil {
    app.errorLog.Fatal(err)
  }

}
