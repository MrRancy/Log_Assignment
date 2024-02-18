package store

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const (
	contentTypeJSON = "application/json"
)

func (m *MyMem) TransportLogs() bool {
	startTime := time.Now()

	// Encode logs to JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m.Log); err != nil {
		m.Logger.Error("Error encoding logs to JSON", zap.Error(err))
		return false
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", m.config.PostEndpoint, &buf)
	if err != nil {
		m.Logger.Error("Error creating HTTP request", zap.Error(err))
		return false
	}
	req.Header.Set("Content-Type", contentTypeJSON)

	// Perform HTTP request
	client := &http.Client{Timeout: time.Minute}
	resp, err := client.Do(req)

	// Calculate the duration
	duration := time.Since(startTime)

	// Log the batch size, result status code, and duration
	status := http.StatusNotFound
	if resp != nil {
		status = resp.StatusCode
	}
	m.Logger.Info("Batch Size", zap.String("size", strconv.Itoa(m.config.BatchSize)), zap.Int("Status Code", status), zap.Duration("Duration", duration))

	// Check for errors after making the HTTP request
	if err != nil {
		m.Logger.Error("Error sending logs to server", zap.Error(err))
		return false
	}
	defer resp.Body.Close()
	return status >= http.StatusOK && status <= http.StatusNoContent
}
