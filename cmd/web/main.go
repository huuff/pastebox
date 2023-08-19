package main

import (
	"net/http"
)


func main() {
  args := ParseArgs()


  app, close := newApplication()
  defer close()

  srv := &http.Server {
    Addr: args.Addr(),
    ErrorLog: app.errorLog,
    Handler: app.routes(),
  }

  app.infoLog.Printf("Starting server on %s", args.Addr())
  if err := srv.ListenAndServe(); err != nil {
    app.errorLog.Fatal(err)
  }

}
