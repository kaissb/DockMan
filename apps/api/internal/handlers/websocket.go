// Copyright (c) 2025 Bouali Consulting Inc.
// Author: Kaiss Bouali (kaissb)
// Company: Bouali Consulting Inc.
// GitHub: https://github.com/kaissb

package handlers
import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// StreamStats handles streaming live stats for a container.
func StreamStats(c *gin.Context) {
	containerID := c.Param("id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	defer conn.Close()

	stats, err := DockerClient.ContainerStats(context.Background(), containerID, true)
	if err != nil {
		log.Printf("Failed to get container stats: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error: Could not get container stats."))
		return
	}
	defer stats.Body.Close()

	decoder := json.NewDecoder(stats.Body)
	for {
		var v types.StatsJSON
		if err := decoder.Decode(&v); err != nil {
			// End of stream or error
			break
		}

		// Calculate CPU percentage
		cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
		systemDelta := float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
		cpuPercent := 0.0
		if systemDelta > 0.0 && cpuDelta > 0.0 {
			cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
		}

		// Memory usage
		memUsage := float64(v.MemoryStats.Usage) - float64(v.MemoryStats.Stats["cache"])
		memPercent := (memUsage / float64(v.MemoryStats.Limit)) * 100.0

		msg := map[string]float64{
			"cpu_percent":    cpuPercent,
			"memory_percent": memPercent,
			"memory_usage":   memUsage / (1024 * 1024), // in MiB
		}

		if err := conn.WriteJSON(msg); err != nil {
			// Client closed connection
			break
		}
	}
}

// InteractiveTerminal handles the websocket connection for an interactive terminal.
func InteractiveTerminal(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for terminal: %v", err)
		return
	}
	defer ws.Close()

	containerID := c.Param("id")
	log.Printf("Starting terminal session for container: %s", containerID)

	// 1. Check if the container is running
	inspect, err := DockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Printf("Failed to inspect container %s: %v", containerID, err)
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Could not find container."))
		return
	}
	if !inspect.State.Running {
		log.Printf("Attempted to open terminal on non-running container %s", containerID)
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Container is not running."))
		return
	}

	// 2. Create the exec instance
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/sh"}, // Try /bin/bash if /bin/sh fails
		Tty:          true,
	}

	execID, err := DockerClient.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		log.Printf("Failed to create exec instance in container %s: %v", containerID, err)
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Failed to create terminal session."))
		return
	}

	// 3. Attach to the exec instance
	hijackedResp, err := DockerClient.ContainerExecAttach(context.Background(), execID.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		log.Printf("Failed to attach to exec instance in container %s: %v", containerID, err)
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Failed to attach to terminal."))
		return
	}
	defer hijackedResp.Close()

	// Channel to signal when to close connections
	done := make(chan bool, 2)

	// Goroutine to read from WebSocket and write to container stdin
	go func() {
		defer func() {
			done <- true
		}()

		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}

			if messageType == websocket.TextMessage {
				_, err = hijackedResp.Conn.Write(message)
				if err != nil {
					log.Printf("Error writing to container stdin: %v", err)
					return
				}
			}
		}
	}()

	// Goroutine to read from container stdout/stderr and write to WebSocket
	go func() {
		defer func() {
			done <- true
		}()

		buffer := make([]byte, 4096)
		for {
			n, err := hijackedResp.Reader.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading from container: %v", err)
				}
				return
			}

			if n > 0 {
				err = ws.WriteMessage(websocket.TextMessage, buffer[:n])
				if err != nil {
					log.Printf("Error writing to WebSocket: %v", err)
					return
				}
			}
		}
	}()

	// Wait for either goroutine to finish
	<-done
	log.Printf("Terminal session ended for container: %s", containerID)
}

// StreamLogs handles the WebSocket connection for log streaming.
func StreamLogs(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer ws.Close()

	containerID := c.Param("id")

	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100",
	}

	logReader, err := DockerClient.ContainerLogs(context.Background(), containerID, options)
	if err != nil {
		log.Printf("Error getting container logs for %s: %v", containerID, err)
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Could not retrieve logs for container."))
		return
	}
	defer logReader.Close()

	header := make([]byte, 8)
	for {
		_, err := io.ReadFull(logReader, header)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Println("Log stream for container", containerID, "closed.")
			} else {
				log.Println("Error reading log stream header for container", containerID, ":", err)
			}
			break
		}

		size := binary.BigEndian.Uint32(header[4:])
		if size == 0 {
			continue
		}

		payload := make([]byte, size)
		_, err = io.ReadFull(logReader, payload)
		if err != nil {
			log.Println("Error reading log stream payload for container", containerID, ":", err)
			break
		}

		if err := ws.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Println("Error writing to websocket:", err)
			break
		}
	}
}
