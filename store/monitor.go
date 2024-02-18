package store

import (
	"time"
)

func (m *MyMem) InitializeLogMonitor() {

	defer func() {
		if r := recover(); r != nil {
			m.Logger.Error("Recovered from panic!")
		}
	}()

	// Start a single goroutine to handle both monitoring conditions and triggering transportLogs
	go func() {
		for {
			select {
			case <-time.After(time.Until(m.calculateNextTriggerTime())):
				m.Logger.Info("[Firing] Batch Interval")
				if !m.processTrigger() {
					return
				}
			case <-m.Trigger:
				m.Logger.Info("[Firing] Batch Size Overflow")
				if !m.processTrigger() {
					return
				}
			case <-m.Exit:
				close(m.Trigger)
				return
			}
		}
	}()
}

// processTrigger triggers transportLogs and handles retries
func (m *MyMem) processTrigger() bool {
	if len(m.Log) <= 0 {
		m.Logger.Warn("In Memory cache is empty, Skipping...")
		return true
	}
	success := m.TransportLogs()
	for retries := 0; !success && retries < m.config.Retry; retries++ {
		sleepDuration := time.Duration(1<<uint(retries)) * time.Second
		time.Sleep(sleepDuration)
		success = m.TransportLogs()
	}

	if !success {
		m.Logger.Error("Failed to transport logs after multiple retries.")
		close(m.Exit)
		return false
	}
	m.Purge() // Purging local memory
	return true
}

// calculateNextTriggerTime calculates the next trigger time
func (m *MyMem) calculateNextTriggerTime() time.Time {
	return time.Now().Add(time.Minute * time.Duration(m.config.BatchInterval))
}
