package main

import (
	"os"
	"sync"
	"time"
)

// StateManager handles reading and writing the last-watered timestamp to disk.
type StateManager struct {
	filePath string
	mu       sync.RWMutex
}

// NewStateManager creates a new StateManager for the given file path.
func NewStateManager(filePath string) *StateManager {
	return &StateManager{filePath: filePath}
}

// ReadLastWateredTime reads the timestamp from disk.
func (sm *StateManager) ReadLastWateredTime() (time.Time, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	data, err := os.ReadFile(sm.filePath)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, string(data))
}

// WriteLastWateredTime persists the timestamp to disk.
func (sm *StateManager) WriteLastWateredTime(t time.Time) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	return os.WriteFile(sm.filePath, []byte(t.Format(time.RFC3339)), 0o644)
}

// GetOrDefault returns the stored timestamp, or the provided default if unavailable.
func (sm *StateManager) GetOrDefault(defaultTime time.Time) time.Time {
	lastWatered, err := sm.ReadLastWateredTime()
	if err != nil {
		return defaultTime
	}
	return lastWatered
}
