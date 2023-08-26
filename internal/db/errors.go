package db

import (
	"errors"
	"fmt"
)

var ErrPasteNotFound = errors.New("not found")

func NewPasteNotFoundError(id string) error {
  return fmt.Errorf("%q: %w", id, ErrPasteNotFound)
}
