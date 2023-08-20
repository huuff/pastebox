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

  // TODO: Remove, just for testing that it works
  app.pastes.Insert("test", "test", 5)

  app.infoLog.Printf("Starting server on %s", args.Addr())
  if err := srv.ListenAndServe(); err != nil {
    app.errorLog.Fatal(err)
  }

}
