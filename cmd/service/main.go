package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VysMax/organizational-structure/config"
	"github.com/VysMax/organizational-structure/controller"
	"github.com/VysMax/organizational-structure/database"
	"github.com/VysMax/organizational-structure/logger"
	"github.com/VysMax/organizational-structure/repository"
	"github.com/VysMax/organizational-structure/usecase"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	logger, err := logger.Init(cfg)
	if err != nil {
		log.Println("Failed to init logger:", err)
	}

	logger.Info("logger initiated")

	dbConn, err := database.New(cfg, logger)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
	}

	repo := repository.New(dbConn.DB, logger)
	usecase := usecase.New(repo, logger)
	handler := controller.New(usecase, logger)

	http.HandleFunc("/departments", handler.CreateDepartment)

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
