package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"brutus/azure-api/internal/models"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

// LogEntry handles single log entry from frontend
func (h *LogHandler) LogEntry(c *gin.Context) {
	var logEntry models.LogEntry
	if err := c.ShouldBindJSON(&logEntry); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid log entry format",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	h.writeLogToContainer(logEntry)

	c.JSON(http.StatusOK, models.LogResponse{
		Success: true,
		Message: "Log entry received",
		Count:   1,
	})
}

// LogBatch handles batch of log entries from frontend
func (h *LogHandler) LogBatch(c *gin.Context) {
	var logBatch models.LogBatch
	if err := c.ShouldBindJSON(&logBatch); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid log batch format",
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	for _, logEntry := range logBatch.Logs {
		h.writeLogToContainer(logEntry)
	}

	c.JSON(http.StatusOK, models.LogResponse{
		Success: true,
		Message: "Log batch received",
		Count:   len(logBatch.Logs),
	})
}

// writeLogToContainer writes the log entry to the container's stdout
// This ensures logs appear in Docker/Kubernetes container logs
func (h *LogHandler) writeLogToContainer(logEntry models.LogEntry) {
	// Parse timestamp to ensure proper formatting
	timestamp := logEntry.Timestamp
	if parsedTime, err := time.Parse(time.RFC3339, logEntry.Timestamp); err == nil {
		timestamp = parsedTime.Format("2006-01-02T15:04:05.000Z")
	}

	// Create structured log output for container logs
	logOutput := map[string]interface{}{
		"timestamp": timestamp,
		"level":     string(logEntry.Level),
		"source":    logEntry.Source,
		"sessionId": logEntry.SessionID,
		"message":   logEntry.Message,
		"url":       logEntry.URL,
		"userAgent": logEntry.UserAgent,
	}

	// Add context if provided
	if logEntry.Context != nil && len(logEntry.Context) > 0 {
		logOutput["context"] = logEntry.Context
	}

	// Convert to JSON for structured logging
	jsonBytes, err := json.Marshal(logOutput)
	if err != nil {
		// Fallback to simple string format if JSON marshaling fails
		log.Printf("[FRONTEND-%s] %s: %s (Session: %s)",
			string(logEntry.Level),
			timestamp,
			logEntry.Message,
			logEntry.SessionID,
		)
		return
	}

	// Write structured log to stdout (appears in container logs)
	log.Printf("[FRONTEND-LOG] %s", string(jsonBytes))
}