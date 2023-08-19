package main

import (
	"log"
  "os"
	"go.mongodb.org/mongo-driver/mongo"
  "xyz.haff/pastebox/internal/db"
)

type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  mongo *mongo.Client
}

func newApplication() (application, func()) {
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  mongo, closeMongo, err := db.ConnectToMongo(infoLog)
  if err != nil {
    errorLog.Fatal(err)
  }

  return application {
    infoLog: infoLog,
    errorLog: errorLog,
    mongo: mongo,
  }, closeMongo
}
