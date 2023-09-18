package main

import (
  "github.com/samber/lo"

	"github.com/gookit/validate"
)

func errorsToMap(errors validate.Errors) map[string]string {
  return lo.MapValues(errors.All(), func(errs map[string]string, _ string) string {
    for _, v := range errs {
      return v
    }
    return ""
  })
}
