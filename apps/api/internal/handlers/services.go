// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package handlers
import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"

	"docker-manager/api/internal/database"
	"docker-manager/api/internal/models"
	"github.com/gin-gonic/gin"
)
// ListServices handles listing all top-level services for an environment.
func ListServices(c *gin.Context) {
	environmentID := c.Param("id")
	var services []models.Service
	if err := database.DB.Where("environment_id = ? AND parent_service_id IS NULL", environmentID).Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list services"})
		return
	}
	c.JSON(http.StatusOK, services)
}

func GetServiceDetails(c *gin.Context) {
	serviceID := c.Param("id")
	var service models.Service
	if err := database.DB.Preload("SubServices").First(&service, serviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, service)
}
// CreateService handles the creation of a new service in an environment.
func CreateService(c *gin.Context) {
	environmentID := c.Param("id")
	var service models.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	envID, _ := strconv.Atoi(environmentID)
	service.EnvironmentID = uint(envID)

	switch service.Type {
	case "container":
		if service.Image == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required for container type"})
			return
		}
	case "compose":
		if service.ComposePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ComposePath is required for compose type"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service type"})
		return
	}

	if err := database.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	c.JSON(http.StatusOK, service)
}

// UpService handles starting a service.
func UpService(c *gin.Context) {
	serviceID := c.Param("id") // Corrected to use ID from the URL path
	var service models.Service
	if err := database.DB.First(&service, serviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	switch service.Type {
	case "container":
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Container start not implemented yet"})
	case "compose":
		cmd := exec.Command("docker-compose", "-f", service.ComposePath, "up", "-d")
		cmd.Dir = filepath.Dir(service.ComposePath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run docker-compose up", "output": string(output)})
			return
		}
	}
}

// DownService handles stopping a service.
func DownService(c *gin.Context) {
	// Implementation for stopping a service
	serviceID := c.Param("id")
	var service models.Service
	if err := database.DB.First(&service, serviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	if service.Type != "compose" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Down operation is only for compose services"})
		return
	}

	cmd := exec.Command("docker-compose", "-f", service.ComposePath, "down")
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bring down service: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service brought down successfully"})
}

// ScaleService scales a specific service within a compose stack.
func ScaleService(c *gin.Context) {
	type ScaleRequest struct {
		SubServiceName string `json:"sub_service_name"`
		Replicas       int    `json:"replicas"`
	}

	serviceID := c.Param("id")
	var request ScaleRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var service models.Service
	if err := database.DB.First(&service, serviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	if service.Type != "compose" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Scale operation is only for compose services"})
		return
	}

	scaleArg := fmt.Sprintf("%s=%d", request.SubServiceName, request.Replicas)
	cmd := exec.Command("docker", "compose", "-f", service.ComposePath, "up", "-d", "--scale", scaleArg)

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scale service: %s", string(output))})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Service %s scaled to %d replicas", request.SubServiceName, request.Replicas)})
}
