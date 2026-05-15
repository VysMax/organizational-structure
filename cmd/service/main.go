package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VysMax/organizational-structure/config"
	"github.com/VysMax/organizational-structure/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}
	fmt.Println(cfg.File)

	logger, err := logger.Init(cfg)
	if err != nil {
		log.Println("Failed to init logger:", err)
	}

	logger.Info("logger initiated")

	fmt.Println("Hello")

	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server", "port", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server start failed", "error", err)
			return
		}
	}()

	<-quit
	logger.Info("Received shutdown signal, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
	}

	logger.Info("Server shutdown complete")
}
