package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"write_base/config"
	"write_base/pkg/di"
)

func main() {
	// Loading Environment
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}
	// Create Container
	container, err := di.NewContainer(cfg)
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: container.Router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server running on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🔻 Shutting down server...")

	// Give in-flight requests 5s to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Server shutdown error: %v", err)
	}

	// Close MongoDB connection
	if err := container.MongoClient.Disconnect(ctx); err != nil {
		log.Printf("⚠️ MongoDB disconnect error: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
