package main

import (
  "bytes"
  "io"
  "log"
  "net/http"
  "net/http/httptest"
  "testing"
  "xyz.haff/pastebox/internal/assert"
)

func TestPing(t *testing.T) {
  app := &application {
    errorLog: log.New(io.Discard, "", 0),
    infoLog: log.New(io.Discard, "", 0),
  }

  ts := httptest.NewServer(app.routes())
  defer ts.Close()

  rs, err := ts.Client().Get(ts.URL + "/ping")
  if err != nil {
    t.Fatal(err)
  }

  assert.Equal(t, rs.StatusCode, http.StatusOK)

  defer rs.Body.Close()
  body, err := io.ReadAll(rs.Body)
  if err != nil {
    t.Fatal(err)
  }
  bytes.TrimSpace(body)

  assert.Equal(t, string(body), "OK")
}
