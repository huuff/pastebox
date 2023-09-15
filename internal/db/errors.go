package db

import (
	"errors"
	"fmt"
)

var ErrRecordNotFound = errors.New("not found")

func NewRecordNotFoundError(kind, id string) error {
  return fmt.Errorf("%s %q: %w", kind, id, ErrRecordNotFound)
}
