package dao

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Paste struct {
  ID string
  Title string
  Content string
  Created time.Time
  Expires time.Time
}

type PasteDAO struct {
  mongo *mongo.Client
}

func NewPasteDAO(mongo *mongo.Client) *PasteDAO {
  return &PasteDAO { mongo }
}

func (dao *PasteDAO) Insert(title string, content string, expires int) (string, error) {
  return "", nil
}

func (dao *PasteDAO) Get(id string) (*Paste, error) {
  return nil, nil
}

func (dao *PasteDAO) Latest() ([]*Paste, error) {
  return nil, nil
}
