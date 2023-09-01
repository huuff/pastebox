package main

import (
  "html/template"
  "path/filepath"

  "xyz.haff/pastebox/internal/models"
)

type templateData struct {
  Paste *models.Paste
  Pastes []*models.Paste
}

func newTemplateCache() (map[string]*template.Template, error) {
  cache := map[string]*template.Template{}

  pages, err := filepath.Glob("./ui/html/pages/*.gotmpl")
  if err != nil {
    return nil, err
  }

  for _, page := range pages {
    name := filepath.Base(page)

    files := []string {
      "./ui/html/base.gotmpl",
      "./ui/html/partials/nav.gotmpl",
      page,
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
      return nil, err
    }

    cache[name] = ts
  }

  return cache, nil
}
