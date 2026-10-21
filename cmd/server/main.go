package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	catalog "github.com/bbdw-gh-action/bbw-catalog"
	"github.com/getsentry/sentry-go"
)

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	if err := catalog.InitTelemetry(); err != nil {
		log.Fatal("Configuração Sentry inválida")
	}
	defer sentry.Flush(2 * time.Second)
	handler, err := catalog.NewHandler(env("CORS_ORIGIN", "http://localhost:3000"))
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Addr: ":" + env("PORT", "3002"), Handler: handler,
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}()
	log.Printf("Catálogo disponível na porta %s", env("PORT", "3002"))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
