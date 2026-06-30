package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

// startServer sets up the Echo router, configures server parameters, and handles graceful shutdown
func startServer(port int, db *gorm.DB) {
	e := echo.New()

	// Register Routes and Handlers
	registerHandlers(e, db)

	// Start standard net/http Server using Echo as handler
	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Log.Infof("HTTP Server starting on address %s", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Shutting down the server due to error: %v", err)
		}
	}()

	// Graceful Shutdown implementation
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Log.Info("Shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Log.Info("Server exited gracefully")
}
