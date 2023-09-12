package main

import (
	"html/template"
	"log"
	"os"

	"xyz.haff/pastebox/internal/db"
	"xyz.haff/pastebox/internal/models"
  "github.com/go-playground/form/v4"
)

type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  pastes *models.PasteDAO
  templateCache map[string]*template.Template
  formDecoder *form.Decoder
}

func newApplication() (application, func()) {
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  mongo, closeMongo, err := db.ConnectToMongo(infoLog)
  if err != nil {
    errorLog.Fatal(err)
  }

  pastes := models.NewPasteDAO(mongo, infoLog)

  templateCache, err := newTemplateCache()
  if err != nil {
    errorLog.Fatal(err)
  }

  return application {
    infoLog: infoLog,
    errorLog: errorLog,
    pastes: pastes,
    templateCache: templateCache,
    formDecoder: form.NewDecoder(),
  }, closeMongo
}
