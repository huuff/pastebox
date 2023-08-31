package main

import (
	"log"
	"os"

	"xyz.haff/pastebox/internal/models"
	"xyz.haff/pastebox/internal/db"
)

type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  pastes *models.PasteDAO
}

func newApplication() (application, func()) {
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  mongo, closeMongo, err := db.ConnectToMongo(infoLog)
  if err != nil {
    errorLog.Fatal(err)
  }

  pastes := models.NewPasteDAO(mongo, infoLog)

  return application {
    infoLog: infoLog,
    errorLog: errorLog,
    pastes: pastes,
  }, closeMongo
}
