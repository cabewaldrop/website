deploy: test css
  flyctl deploy

css:
  npx tailwindcss -i input.css -o static/output.css

test:
  go test ./...

up:
  docker-compose up -d

down:
  docker-compose down

local css:
  air
