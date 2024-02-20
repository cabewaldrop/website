deploy: test css
  flyctl deploy

css:
  npx tailwindcss -i input.css -o static/output.css

test: up
  go test ./...
  docker-compose down

up:
  docker-compose build
  docker-compose up -d

down:
  docker-compose down

local css:
  air
