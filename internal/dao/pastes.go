package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
}

func NewPasteDAO(mongo *mongo.Client) *PasteDAO {
  collection := mongo.Database(db.DatabaseName).Collection("pastes")

  return &PasteDAO { collection }
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

  // TODO: A method to do this conversion automatically?
  return result.InsertedID.(primitive.ObjectID).String(), nil
}

func (dao *PasteDAO) Get(id string) (*Paste, error) {
  return nil, nil
}

func (dao *PasteDAO) Latest() ([]*Paste, error) {
  return nil, nil
}
