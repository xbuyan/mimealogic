package main

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	wateringEvents    int
	totalWateringTime time.Duration
	mu                sync.Mutex
}

func (m *Metrics) RecordWatering(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wateringEvents++
	m.totalWateringTime += duration
}

func (m *Metrics) PrintSummary() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Using Go's built-in max() function (requires Go 1.21+)
	avg := m.totalWateringTime / time.Duration(max(1, m.wateringEvents))

	fmt.Println("\n📊 === SYSTEM METRICS ===")
	fmt.Printf("   Total watering events:    %d\n", m.wateringEvents)
	fmt.Printf("   Total watering time:      %v\n", m.totalWateringTime)
	fmt.Printf("   Average watering duration: %v\n", avg)
}
