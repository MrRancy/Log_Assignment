package store

import (
	"go.uber.org/zap"
	"mrrancy/logAssignment/models"
	"time"
)

type TriggerInfo struct {
	NextTriggerTime time.Time
	Log             []models.LogPayload
}

func (m *MyMem) InitializeLogMonitor() {

	defer func() {
		if r := recover(); r != nil {
			m.Logger.Error("Recovered from panic:", zap.Any("panic", r))
		}
	}()

	// Start a single goroutine to handle both monitoring conditions and triggering transportLogs
	go m.startLogMonitor()
}

func (m *MyMem) startLogMonitor() {
	for {
		select {
		case <-time.After(time.Until(m.calculateNextTriggerTime().NextTriggerTime)):
			batch := m.calculateNextTriggerTime().Log
			m.Log = nil
			m.Logger.Info("[Firing] Batch Interval")
			if !m.processTrigger(batch) {
				m.Logger.Error("Error processing batch interval")
				return
			}
		case batch := <-m.Trigger:
			m.Logger.Info("[Firing] Batch Size Overflow")
			if !m.processTrigger(batch) {
				m.Logger.Error("Error processing batch size overflow")
				return
			}
		case <-m.Exit:
			close(m.Trigger)
			m.Logger.Info("Exiting log monitor goroutine")
			return
		}
	}
}

// processTrigger triggers transportLogs and handles retries
func (m *MyMem) processTrigger(batch []models.LogPayload) bool {
	if len(batch) <= 0 {
		m.Logger.Warn("In Memory cache is empty, Skipping...")
		return true
	}
	success := m.TransportLogs(batch)
	for retries := 0; !success && retries < m.config.Retry; retries++ {
		sleepDuration := time.Duration(1<<uint(retries)) * time.Second
		time.Sleep(sleepDuration)
		success = m.TransportLogs(batch)
	}

	if !success {
		m.Logger.Error("Failed to transport logs after multiple retries.")
		close(m.Exit)
		return false
	}
	return true
}

// calculateNextTriggerTime calculates the next trigger time and returns it along with m.Log
func (m *MyMem) calculateNextTriggerTime() *TriggerInfo {
	return &TriggerInfo{
		NextTriggerTime: time.Now().Add(time.Minute * time.Duration(m.config.BatchInterval)),
		Log:             m.Log,
	}
}
