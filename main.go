package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cabewaldrop/website/pkg/server"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var ADDRESS = fmt.Sprintf(":%s", os.Getenv("PORT"))

func main() {
	log.Info().Msgf("Starting server with port: %s", os.Getenv("PORT"))
	e := echo.New()
	server.ConfigureServer(e)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT)
	defer cancel()

	startServer(e, ctx)
	<-ctx.Done()

	shutdownCtx, shutdown := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdown()

	e.Shutdown(shutdownCtx)
}

func startServer(e *echo.Echo, ctx context.Context) {
	go e.Start(ADDRESS)
}
