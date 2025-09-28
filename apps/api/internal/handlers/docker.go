// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package handlers

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"docker-manager/api/internal/database"
	"docker-manager/api/internal/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/gin-gonic/gin"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// DockerClientInterface defines the methods we use from the Docker client.
// This makes the client mockable for testing.
type DockerClientInterface interface {
	ContainerList(context.Context, container.ListOptions) ([]types.Container, error)
	ContainerStart(context.Context, string, container.StartOptions) error
	ContainerStop(context.Context, string, container.StopOptions) error
	ContainerRestart(context.Context, string, container.StopOptions) error
	ContainerRemove(context.Context, string, container.RemoveOptions) error
	ImageList(context.Context, types.ImageListOptions) ([]types.ImageSummary, error)
	ImagePull(context.Context, string, types.ImagePullOptions) (io.ReadCloser, error)
	ImageRemove(context.Context, string, types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error)
	ImageInspectWithRaw(context.Context, string) (types.ImageInspect, []byte, error)
	ContainerCreate(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) (container.CreateResponse, error)
	ContainerLogs(context.Context, string, container.LogsOptions) (io.ReadCloser, error)
	ContainerInspect(context.Context, string) (types.ContainerJSON, error)
	ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error)
	ContainerStats(ctx context.Context, containerID string, stream bool) (types.ContainerStats, error)
}

var DockerClient DockerClientInterface

// ListContainers handles listing all containers.
func ListContainers(c *gin.Context) {
	containers, err := DockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list containers"})
		return
	}
	c.JSON(http.StatusOK, containers)
}

// StartContainer handles starting a container.
func StartContainer(c *gin.Context) {
	containerID := c.Param("id")
	if err := DockerClient.ContainerStart(context.Background(), containerID, container.StartOptions{}); err != nil {
		var errContainerAlreadyStarted interface {
			StatusCode() int
		}
		if errors.As(err, &errContainerAlreadyStarted) && errContainerAlreadyStarted.StatusCode() == http.StatusNotModified {
			c.JSON(http.StatusOK, gin.H{"status": "already running"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "started"})
}

// StopContainer handles stopping a container.
func StopContainer(c *gin.Context) {
	containerID := c.Param("id")
	timeout := 10
	if err := DockerClient.ContainerStop(context.Background(), containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}

// PullImage handles pulling a Docker image from a registry.
func PullImage(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image name"})
		return
	}
	out, err := DockerClient.ImagePull(context.Background(), body.Name, types.ImagePullOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pull image"})
		return
	}
	defer out.Close()
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(c.Writer, out)
}

// DeleteImage handles deleting a Docker image.
func DeleteImage(c *gin.Context) {
	imageID := c.Param("id")
	inspect, _, err := DockerClient.ImageInspectWithRaw(context.Background(), imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect image"})
		return
	}
	for _, tag := range inspect.RepoTags {
		if strings.Contains(tag, "dockman") {
			c.JSON(http.StatusForbidden, gin.H{"error": "This image belongs to the DockMan application and cannot be deleted."}) 
			return
		}
	}
	_, err = DockerClient.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove image"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// ListImages handles listing all Docker images.
func ListImages(c *gin.Context) {
	images, err := DockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list images"})
		return
	}
	c.JSON(http.StatusOK, images)
}

// DeleteContainer handles deleting a container.
func DeleteContainer(c *gin.Context) {
	containerID := c.Param("id")
	if err := DockerClient.ContainerRemove(context.Background(), containerID, container.RemoveOptions{Force: true}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete container"})
		return
	}
	if err := database.DB.Where("container_id = ?", containerID).Delete(&models.Service{}).Error; err != nil {
		log.Printf("Failed to delete service from database for container %s: %v", containerID, err)
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// RestartContainer handles restarting a container.
func RestartContainer(c *gin.Context) {
	containerID := c.Param("id")
	timeout := 10
	if err := DockerClient.ContainerRestart(context.Background(), containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "restarted"})
}
