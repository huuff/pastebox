package main

import (
	"context"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"xyz.haff/pastebox/internal/db"
	"xyz.haff/pastebox/internal/models"
)

type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  pastes *models.PasteDAO
  users *models.UserDAO
  templateCache map[string]*template.Template
  formDecoder *form.Decoder
  sessionManager *scs.SessionManager
}

func newApplication() (application, func()) {
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  timeoutCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  mongo, closeMongo, err := db.ConnectToMongo(infoLog, timeoutCtx)
  if err != nil {
    errorLog.Fatal(err)
  }

  sessionManager := scs.New()
  sessionManager.Store = mongodbstore.New(mongo.Database(db.DatabaseName))

  pastes := models.NewPasteDAO(mongo, infoLog)
  users := models.NewUserDAO(mongo, infoLog)

  templateCache, err := newTemplateCache()
  if err != nil {
    errorLog.Fatal(err)
  }

  return application {
    infoLog: infoLog,
    errorLog: errorLog,
    pastes: pastes,
    users: users,
    templateCache: templateCache,
    formDecoder: form.NewDecoder(),
    sessionManager: sessionManager,
  }, closeMongo
}
