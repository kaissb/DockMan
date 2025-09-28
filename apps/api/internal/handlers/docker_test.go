// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package handlers

import (
	"context"
	"docker-manager/api/internal/database"
	"docker-manager/api/internal/models"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockDockerClient is a mock implementation of the Docker client for testing.
// It implements the DockerClientInterface.
type MockDockerClient struct {
	mock.Mock
}

// setupTestRouter creates a new Gin router for testing and injects the mock client.
func setupTestRouter(mockClient *MockDockerClient) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Setup in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	database.DB = db

	// Migrate the schema for the test database
	db.AutoMigrate(&models.Project{}, &models.Environment{}, &models.Service{})

	router := gin.Default()

	// Inject the mock client into the handlers package
	DockerClient = mockClient

	return router
}

func TestListContainers(t *testing.T) {
	mockClient := new(MockDockerClient)
	router := setupTestRouter(mockClient)
	router.GET("/containers", ListContainers)

	expectedContainers := []types.Container{
		{ID: "test-id-1", Image: "ubuntu"},
		{ID: "test-id-2", Image: "redis"},
	}
	mockClient.On("ContainerList", mock.Anything, mock.AnythingOfType("container.ListOptions")).Return(expectedContainers, nil)

	req, _ := http.NewRequest("GET", "/containers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseContainers []types.Container
	err := json.Unmarshal(w.Body.Bytes(), &responseContainers)
	assert.NoError(t, err)
	assert.Equal(t, expectedContainers, responseContainers)
	mockClient.AssertExpectations(t)
}

func TestListImages(t *testing.T) {
	mockClient := new(MockDockerClient)
	router := setupTestRouter(mockClient)
	router.GET("/images", ListImages)

	expectedImages := []types.ImageSummary{
		{ID: "img-id-1", RepoTags: []string{"ubuntu:latest"}},
		{ID: "img-id-2", RepoTags: []string{"redis:latest"}},
	}
	mockClient.On("ImageList", mock.Anything, mock.AnythingOfType("types.ImageListOptions")).Return(expectedImages, nil)

	req, _ := http.NewRequest("GET", "/images", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseImages []types.ImageSummary
	err := json.Unmarshal(w.Body.Bytes(), &responseImages)
	assert.NoError(t, err)
	assert.Equal(t, expectedImages, responseImages)
	mockClient.AssertExpectations(t)
}

// Implement all methods of the DockerClientInterface for the mock

func (m *MockDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.Container), args.Error(1)
}

func (m *MockDockerClient) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	args := m.Called(ctx, containerID, options)
	return args.Error(0)
}

func (m *MockDockerClient) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	args := m.Called(ctx, containerID, options)
	return args.Error(0)
}

func (m *MockDockerClient) ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error {
	args := m.Called(ctx, containerID, options)
	return args.Error(0)
}

func (m *MockDockerClient) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	args := m.Called(ctx, containerID, options)
	return args.Error(0)
}

func (m *MockDockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.ImageSummary), args.Error(1)
}

func (m *MockDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, refStr, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockDockerClient) ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	args := m.Called(ctx, imageID, options)
	return args.Get(0).([]types.ImageDeleteResponseItem), args.Error(1)
}

func (m *MockDockerClient) ImageInspectWithRaw(ctx context.Context, imageID string) (types.ImageInspect, []byte, error) {
	args := m.Called(ctx, imageID)
	return args.Get(0).(types.ImageInspect), args.Get(1).([]byte), args.Error(2)
}

func (m *MockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error) {
	args := m.Called(ctx, config, hostConfig, networkingConfig, platform, containerName)
	return args.Get(0).(container.CreateResponse), args.Error(1)
}

func (m *MockDockerClient) ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error) {
	args := m.Called(ctx, container, options)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockDockerClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	args := m.Called(ctx, containerID)
	return args.Get(0).(types.ContainerJSON), args.Error(1)
}

func (m *MockDockerClient) ContainerExecCreate(ctx context.Context, container string, config types.ExecConfig) (types.IDResponse, error) {
	args := m.Called(ctx, container, config)
	return args.Get(0).(types.IDResponse), args.Error(1)
}

func (m *MockDockerClient) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	args := m.Called(ctx, execID, config)
	return args.Get(0).(types.HijackedResponse), args.Error(1)
}

func (m *MockDockerClient) ContainerExecInspect(ctx context.Context, execID string) (types.ContainerExecInspect, error) {
	args := m.Called(ctx, execID)
	return args.Get(0).(types.ContainerExecInspect), args.Error(1)
}

func (m *MockDockerClient) ContainerStats(ctx context.Context, containerID string, stream bool) (types.ContainerStats, error) {
	args := m.Called(ctx, containerID, stream)
	return args.Get(0).(types.ContainerStats), args.Error(1)
}
