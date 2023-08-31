package main

import (
  // TODO: I only use the model here, yet it's in the DAO module... maybe I should separate both?
  "xyz.haff/pastebox/internal/dao"
)

type templateData struct {
  Paste *dao.Paste
}
