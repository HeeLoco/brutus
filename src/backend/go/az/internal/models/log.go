package models

import "time"

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogEntry represents a log entry from the frontend
type LogEntry struct {
	Timestamp string                 `json:"timestamp" binding:"required"`
	Level     LogLevel               `json:"level" binding:"required"`
	Message   string                 `json:"message" binding:"required"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Source    string                 `json:"source" binding:"required"`
	SessionID string                 `json:"sessionId" binding:"required"`
	UserAgent string                 `json:"userAgent,omitempty"`
	URL       string                 `json:"url,omitempty"`
}

// LogBatch represents a batch of log entries
type LogBatch struct {
	Logs []LogEntry `json:"logs" binding:"required"`
}

// LogResponse represents the response after logging
type LogResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Count   int    `json:"count,omitempty"`
}