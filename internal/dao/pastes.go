package dao

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"xyz.haff/pastebox/internal/db"
)

type Paste struct {
  ID string
  Title string
  Content string
  Created time.Time
  Expires time.Time
}

type PasteDAO struct {
  collection *mongo.Collection
  infoLog *log.Logger
}

func NewPasteDAO(mongo *mongo.Client, infoLog *log.Logger) *PasteDAO {
  collection := mongo.Database(db.DatabaseName).Collection("pastes")

  return &PasteDAO { collection, infoLog }
}

func (dao *PasteDAO) Insert(title string, content string, expires int) (string, error) {
  result, err := dao.collection.InsertOne(context.TODO(), Paste {
    Title: title,
    Content: content,
    Created: time.Now(),
    Expires: time.Now().AddDate(0, 0, expires),
  } )

  if err != nil {
    return "", err
  }

  id := db.GetInsertOneStringId(result)
  dao.infoLog.Printf("Inserted %s", id)

  return id, nil
}

func (dao *PasteDAO) Get(id string) (*Paste, error) {
  return nil, nil
}

func (dao *PasteDAO) Latest() ([]*Paste, error) {
  return nil, nil
}
