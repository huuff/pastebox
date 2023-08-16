package main

import (
  "net/http"
  "path/filepath"
)


type nonIndexingFileSystem struct {
  fs http.FileSystem
}

func (nifs nonIndexingFileSystem) Open(path string) (http.File, error) {
  f, err := nifs.fs.Open(path)
  if err != nil {
    return nil, err
  }

  s, err := f.Stat()
  if s.IsDir() {
    index := filepath.Join(path, "index.html")
    if _, err := nifs.fs.Open(index); err != nil {
      closeErr := f.Close()
      if closeErr != nil {
        return nil, closeErr
      }

      return nil, err
    }
  }

  return f, nil
}
