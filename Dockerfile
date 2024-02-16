FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cabewaldrop .

FROM alpine:edge

WORKDIR /app

COPY --from=build /app/cabewaldrop .
COPY --from=build /app/content ./content/
COPY --from=build /app/images ./images/

EXPOSE 42069/tcp

ENTRYPOINT ["/app/cabewaldrop"]
