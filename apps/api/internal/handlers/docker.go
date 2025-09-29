// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// DockerClientInterface defines the methods we use from the Docker client.
// This helps in mocking the Docker client for tests.
type DockerClientInterface interface {
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error)
	ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error)
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
	ContainerStats(ctx context.Context, containerID string, stream bool) (types.ContainerStats, error)
	ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error)
	ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error)
	ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error)

	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error)
	ImageInspectWithRaw(ctx context.Context, imageID string) (types.ImageInspect, []byte, error)
}

// DockerClient is an instance of the Docker client that satisfies the DockerClientInterface.
var DockerClient DockerClientInterface

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
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// GetImageDetails returns detailed information about a single image.
func GetImageDetails(c *gin.Context) {
	imageID := c.Param("id")
	inspect, _, err := DockerClient.ImageInspectWithRaw(context.Background(), imageID)
	if err != nil {
		if client.IsErrNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to inspect image: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, inspect)
}

// ListImages handles listing all Docker images.
func ListImages(c *gin.Context) {
	images, err := DockerClient.ImageList(context.Background(), types.ImageListOptions{})
{{ ... }}
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
