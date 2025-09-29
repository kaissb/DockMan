// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package handlers

import (
	"log"
	"net/http"
	"strconv"

	"docker-manager/api/internal/database"
	"docker-manager/api/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateEnvironmentVariable adds a new variable to an environment.
func CreateEnvironmentVariable(c *gin.Context) {
	environmentID, _ := strconv.Atoi(c.Param("id"))
	var variable models.EnvironmentVariable
	if err := c.ShouldBindJSON(&variable); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	variable.EnvironmentID = uint(environmentID)

	if err := database.DB.Create(&variable).Error; err != nil {
		// Log the detailed error to the console
		log.Printf("Error creating environment variable: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create environment variable"})
		return
	}
	c.JSON(http.StatusOK, variable)
}

// ListEnvironmentVariables lists all variables for an environment.
func ListEnvironmentVariables(c *gin.Context) {
	environmentID := c.Param("id")
	var variables []models.EnvironmentVariable
	if err := database.DB.Where("environment_id = ?", environmentID).Find(&variables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list environment variables"})
		return
	}
	c.JSON(http.StatusOK, variables)
}

// UpdateEnvironmentVariable updates an existing variable.
func UpdateEnvironmentVariable(c *gin.Context) {
	variableID := c.Param("varId")
	var variable models.EnvironmentVariable
	if err := database.DB.First(&variable, variableID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Environment variable not found"})
		return
	}

	var input models.EnvironmentVariable
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	variable.Key = input.Key
	variable.Value = input.Value

	if err := database.DB.Save(&variable).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update environment variable"})
		return
	}
	c.JSON(http.StatusOK, variable)
}

// DeleteEnvironmentVariable deletes a variable.
func DeleteEnvironmentVariable(c *gin.Context) {
	variableID := c.Param("varId")
	if err := database.DB.Delete(&models.EnvironmentVariable{}, variableID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete environment variable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Environment variable deleted"})
}
