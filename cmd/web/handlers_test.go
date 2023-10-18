package main

import (
  "net/http"
  "testing"
  "xyz.haff/pastebox/internal/assert"
)

func TestPing(t *testing.T) {
  app := newTestApplication(t)

  ts := newTestServer(t, app.routes())
  defer ts.Close()

  code, _, body := ts.get(t, "/ping")

  assert.Equal(t, code, http.StatusOK)
  assert.Equal(t, body, "OK")
}

func TestPasteView(t *testing.T) {
  app := newTestApplication(t)

  ts := newTestServer(t, app.routes())
  defer ts.Close()

  tests := []struct {
    name string
    urlPath string
    wantCode int
    wantBody string
  } {
    {
      name: "Valid ID",
      urlPath: "/paste/view/1",
      wantCode: http.StatusOK,
      wantBody: "An old silent pond...",
    },
    {
      name: "Non-existent ID",
      urlPath: "/paste/view/2",
      wantCode: http.StatusNotFound,
    },
    {
      name: "Negative ID",
      urlPath: "/paste/view/-1",
      wantCode: http.StatusNotFound,
    },
    {
      name: "Decimal ID",
      urlPath: "/paste/view/1.23",
      wantCode: http.StatusNotFound,
    },
    {
      name: "String ID",
      urlPath: "/paste/view/foo",
      wantCode: http.StatusNotFound,
    },
    {
      name: "Empty ID",
      urlPath: "/paste/view/",
      wantCode: http.StatusNotFound,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      code, _, body := ts.get(t, tt.urlPath)

      assert.Equal(t, code, tt.wantCode)

      if tt.wantBody != "" {
        assert.StringContains(t, body, tt.wantBody)
      }
    })

  }
}
