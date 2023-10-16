package mocks

import (
  "xyz.haff/pastebox/internal/db"
  customErr "xyz.haff/pastebox/internal/errors"
)

type UserDAO struct {}

func (dao *UserDAO) Insert(name, email, password string) error {
  switch email {
  case "dupe@example.com":
    return customErr.DuplicateEmailError
  default:
    return nil
  }
}

func (dao *UserDAO) Authenticate(email, password string) (int, error) {
  if email == "alice@example.com" && password == "pa$$word" {
    return 1, nil
  }

  // TODO: Move this error out of db and put it in the errors package
  return 0, db.ErrInvalidCredentials
}

func (dao *UserDAO) Exists(id string) (bool, error) {
  switch id {
    case "1":
      return true, nil
    default:
      return false, nil
  }
}
