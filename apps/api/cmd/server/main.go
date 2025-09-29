// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package main

import (
	"log"

	"docker-manager/api/internal/database"
	"docker-manager/api/internal/handlers"
	"docker-manager/api/internal/models"
	"docker-manager/api/internal/router"

	"github.com/docker/docker/client"
)

func main() {
	// Initialize Docker client
	var err error
	handlers.DockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// The Docker client from client.NewClientWithOpts satisfies the DockerClientInterface
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}

	// Initialize Database
	database.Init()
	database.Migrate(&models.Project{}, &models.Environment{}, &models.Service{}, &models.EnvironmentVariable{})

	// Setup Router
	r := router.Setup()

	// Start Server
	log.Println("Starting API server on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
