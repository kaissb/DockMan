// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package router

import (
	"docker-manager/api/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup initializes and configures the Gin router.
func Setup() *gin.Engine {
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	r.SetTrustedProxies(nil)

	// --- API Endpoints ---

	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Docker container endpoints
	r.GET("/containers", handlers.ListContainers)
	r.POST("/containers/:id/start", handlers.StartContainer)
	r.POST("/containers/:id/stop", handlers.StopContainer)
	r.POST("/containers/:id/restart", handlers.RestartContainer)
	r.DELETE("/containers/:id", handlers.DeleteContainer)

	// Docker image endpoints
	r.GET("/images", handlers.ListImages)
	r.POST("/images/pull", handlers.PullImage)
	r.DELETE("/images/:id", handlers.DeleteImage)

	// WebSocket endpoints
	r.GET("/ws/logs/:id", handlers.StreamLogs)
	r.GET("/ws/terminal/:id", handlers.InteractiveTerminal)
	r.GET("/ws/stats/:id", handlers.StreamStats)

	// Project endpoints
	api := r.Group("/api")
	{
		projects := api.Group("/projects")
		{
			projects.POST("", handlers.CreateProject)
			projects.GET("", handlers.ListProjects)
			projects.GET("/:id", handlers.GetProject)
			projects.POST("/:id/environments", handlers.CreateEnvironment)
			projects.GET("/:id/environments", handlers.ListEnvironments)
		}

		environments := api.Group("/environments")
		{
			environments.POST("/:id/services", handlers.CreateService)
			environments.GET("/:id/services", handlers.ListServices)

			// Environment Variables
			environments.POST("/:id/variables", handlers.CreateEnvironmentVariable)
			environments.GET("/:id/variables", handlers.ListEnvironmentVariables)
			environments.PUT("/:id/variables/:varId", handlers.UpdateEnvironmentVariable)
			environments.DELETE("/:id/variables/:varId", handlers.DeleteEnvironmentVariable)
		}

		services := api.Group("/services")
		{
			services.GET("/:id", handlers.GetServiceDetails)
			services.POST("/:id/up", handlers.UpService)
			services.POST("/:id/down", handlers.DownService)
			services.POST("/:id/scale", handlers.ScaleService)
		}
	}

	return r
}
