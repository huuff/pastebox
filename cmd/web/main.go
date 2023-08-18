package main

import (
	"context"
	"log"
	"net/http"
	"os"
  "fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
}

const mongoUri = "mongodb://root:pass@localhost:27017"
func main() {
  args := ParseArgs()

  serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

  client, err := mongo.Connect(context.TODO(), opts)
  if err != nil {
    panic(err)
  }
  defer func() {
    if err = client.Disconnect(context.TODO()); err != nil {
      panic(err)
    }
  }()

  var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
  fmt.Println("Successfully connected to mongo!")

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
