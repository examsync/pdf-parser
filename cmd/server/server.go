package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/examsync/pdf-parser/utils/config"
	"github.com/examsync/pdf-parser/utils/errors"
	"github.com/examsync/pdf-parser/utils/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// startServer sets up the Gin router engine, configures middleware and graceful shutdown
func startServer(cfg *config.Config, db *gorm.DB) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Attach Recovery, Request Delimiter Logging, and Error Handling Middleware
	r.Use(gin.Recovery())
	r.Use(logger.RequestLoggerMiddleware())
	r.Use(errors.ErrorHandlerMiddleware(cfg.Error))


	// Register Routes and Handlers
	registerHandlers(r, db)

	// Start standard net/http Server using Gin as handler
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Log.Infof("Gin HTTP Server starting on address %s", addr)

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
