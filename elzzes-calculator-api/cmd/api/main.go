package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"elzzes-calculator-api/internal/config"
	"elzzes-calculator-api/internal/logger"
	transport "elzzes-calculator-api/internal/transport/http"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	appLog := logger.New(os.Stdout)

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		appLog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	listener, err := net.Listen("tcp", cfg.Addr())
	if err != nil {
		appLog.Error("listen failed", "addr", cfg.Addr(), "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Handler:           transport.NewHandler(appLog),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ErrorLog:          log.New(logWriter{logger: appLog}, "", 0),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() { errCh <- server.Serve(listener) }()

	appLog.Info("api listening", "addr", listener.Addr().String())

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			appLog.Error("server error", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		appLog.Info("api shutting down")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			appLog.Error("graceful shutdown failed", "error", err)
			os.Exit(1)
		}
		if err := <-errCh; !errors.Is(err, http.ErrServerClosed) {
			appLog.Error("server error", "error", err)
			os.Exit(1)
		}
	}

	appLog.Info("api stopped")
}

type logWriter struct {
	logger *logger.Logger
}

func (w logWriter) Write(p []byte) (int, error) {
	w.logger.Error("http server error", "message", strings.TrimSpace(string(p)))
	return len(p), nil
}
