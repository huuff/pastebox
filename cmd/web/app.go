package main

import (
	"log"
	"os"

	"xyz.haff/pastebox/internal/dao"
	"xyz.haff/pastebox/internal/db"
)

type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  pastes *dao.PasteDAO
}

func newApplication() (application, func()) {
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  mongo, closeMongo, err := db.ConnectToMongo(infoLog)
  if err != nil {
    errorLog.Fatal(err)
  }

  // TODO: Or maybe just pass the mongo and let pastes create the collection?
  pastesCollection := mongo.Database("db").Collection("pastes")

  pastes := dao.NewPasteDAO(pastesCollection)

  return application {
    infoLog: infoLog,
    errorLog: errorLog,
    pastes: pastes,
  }, closeMongo
}
