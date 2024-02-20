FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cabewaldrop .

FROM node:18 AS css
WORKDIR /app

COPY . .
RUN npm install && \
  npx tailwindcss -i input.css -o output.css --minify

FROM alpine:edge

WORKDIR /app

COPY --from=build /app/cabewaldrop .
COPY --from=build /app/content ./content/
COPY --from=build /app/static ./static/
COPY --from=css /app/output.css ./static/output.css

ENTRYPOINT ["/app/cabewaldrop"]
