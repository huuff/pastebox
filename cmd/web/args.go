package main

import (
  "flag"
  "strconv"
)

type Args struct {
  port int
}

func (args Args) Addr() string {
  return ":" + strconv.Itoa(args.port)
}

func ParseArgs() Args {
  port := flag.Int("port", 4000, "Port to serve the app")
  flag.Parse()

  return Args { *port }
}
