package main

import (
	"log"
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
	// Start Server
	container.Router.Run(":"+cfg.ServerPort)
	log.Printf("Server running on port %s", cfg.ServerPort)
}