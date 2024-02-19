package store

import (
	"go.uber.org/zap"
	"mrrancy/logAssignment/models"
	"os"
	"sync"
)

type MyMem struct {
	Log     []models.LogPayload
	Mu      sync.RWMutex
	Trigger chan []models.LogPayload
	Exit    chan os.Signal
	Logger  *zap.Logger
	config  *MemConfig
}

type MemConfig struct {
	BatchSize     int
	BatchInterval int
	PostEndpoint  string
	Retry         int
	RetryInterval int
}

func NewCache(log *zap.Logger, postEndpoint string, batchInterval, batchSize, retry, retryInterval int) *MyMem {
	return &MyMem{
		Log:     make([]models.LogPayload, 0),
		Trigger: make(chan []models.LogPayload),
		Exit:    make(chan os.Signal, 1),
		config: &MemConfig{
			BatchInterval: batchInterval,
			BatchSize:     batchSize,
			PostEndpoint:  postEndpoint,
			Retry:         retry,
			RetryInterval: retryInterval,
		},
		Logger: log,
	}
}

func (m *MyMem) Put(payload models.LogPayload) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Log = append(m.Log, payload)
	if len(m.Log) >= m.config.BatchSize {
		m.Trigger <- m.Log
		m.Log = m.Log[m.config.BatchSize:]
	}
}
