// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package handlers

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"

	"docker-manager/api/internal/database"
	"docker-manager/api/internal/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
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

// GetServiceDetails handles retrieving the details of a single service, including sub-services.
func GetServiceDetails(c *gin.Context) {
	serviceID := c.Param("id")
	var service models.Service
	if err := database.DB.Preload("SubServices").First(&service, serviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	if service.Type == "compose" {
		updateSubServicesFromComposeFile(&service)
		// Reload to get the newly created sub-services
		database.DB.Preload("SubServices").First(&service, serviceID)
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
		c.JSON(http.StatusOK, gin.H{"message": "Compose service started", "output": string(output)})
	}
}

// DownService handles stopping a service.
func DownService(c *gin.Context) {
	serviceID := c.Param("id") // Corrected to use ID from the URL path
	var service models.Service
	if err := database.DB.First(&service, serviceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	switch service.Type {
	case "container":
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Container stop not implemented yet"})
	case "compose":
		cmd := exec.Command("docker-compose", "-f", service.ComposePath, "down")
		cmd.Dir = filepath.Dir(service.ComposePath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run docker-compose down", "output": string(output)})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Compose service stopped", "output": string(output)})
	}
}

func updateSubServicesFromComposeFile(service *models.Service) {
	yamlFile, err := ioutil.ReadFile(service.ComposePath)
	if err != nil { return }

	var composeConfig struct {
		Services map[string]struct {
			Image string `yaml:"image"`
		} `yaml:"services"`
	}
	if yaml.Unmarshal(yamlFile, &composeConfig) != nil { return }

	database.DB.Where("parent_service_id = ?", service.ID).Delete(&models.Service{})

	for name, compService := range composeConfig.Services {
		subService := models.Service{
			Name:            name,
			Type:            "container",
			Image:           compService.Image,
			ParentServiceID: &service.ID,
			EnvironmentID:   service.EnvironmentID,
		}
		database.DB.Create(&subService)
	}
}
