// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package handlers

import (
	"net/http"
	"strconv"

	"docker-manager/api/internal/database"
	"docker-manager/api/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateEnvironment handles the creation of a new environment for a project.
func CreateEnvironment(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var environment models.Environment
	if err := c.ShouldBindJSON(&environment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	environment.ProjectID = uint(projectID)

	if err := database.DB.Create(&environment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create environment"})
		return
	}

	c.JSON(http.StatusOK, environment)
}

// ListEnvironments handles listing all environments for a project.
func ListEnvironments(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var environments []models.Environment
	if err := database.DB.Where("project_id = ?", projectID).Find(&environments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list environments"})
		return
	}

	c.JSON(http.StatusOK, environments)
}
