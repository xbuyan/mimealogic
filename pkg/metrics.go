package main

import (
	"fmt"
	"sync"
	"time"
)

// Metrics tracks watering events and total time spent watering.
type Metrics struct {
	wateringEvents    int
	totalWateringTime time.Duration
	mu                sync.Mutex
}

// RecordWatering logs a single watering event.
func (m *Metrics) RecordWatering(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wateringEvents++
	m.totalWateringTime += duration
}

// PrintSummary prints aggregate stats to stdout.
func (m *Metrics) PrintSummary() {
	m.mu.Lock()
	defer m.mu.Unlock()

	avg := m.totalWateringTime / time.Duration(max(1, m.wateringEvents))

	fmt.Println("\n📊 === SYSTEM METRICS ===")
	fmt.Printf("   Total watering events:    %d\n", m.wateringEvents)
	fmt.Printf("   Total watering time:      %v\n", m.totalWateringTime)
	fmt.Printf("   Average watering duration: %v\n", avg)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
