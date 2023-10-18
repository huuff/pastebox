package mocks

import (
	"time"

	"xyz.haff/pastebox/internal/db"
	"xyz.haff/pastebox/internal/models"
  "context"
)

var mockPaste = &models.Paste{
  ID: "1",
  Title: "An old silent pond",
  Content: "An old silent pond...",
  Created: time.Now(),
  Expires: time.Now(),
}

type PasteDAO struct {}

func (m *PasteDAO) Insert(title string, content string, expires int, ctx context.Context) (string, error) {
  return "2", nil
}

func (m *PasteDAO) Get(id string, ctx context.Context) (*models.Paste, error) {
  switch id {
  case "1":
    return mockPaste, nil
  default:
    return nil, db.ErrRecordNotFound
  }
}

func (m *PasteDAO) Latest(ctx context.Context) ([]models.Paste, error) {
  return []models.Paste{*mockPaste}, nil
}

