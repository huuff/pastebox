package main

import (
  "xyz.haff/pastebox/internal/models"
)

type templateData struct {
  Paste *models.Paste
  Pastes []*models.Paste
}
