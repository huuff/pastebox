dev *args='':
  docker compose up -d
  go run ./cmd/web {{args}}
