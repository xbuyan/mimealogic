package main

import (
	"fmt"
	"sync"
	"time"
)

// HardwareController simulates or interfaces with actual hardware.
type HardwareController struct {
	enabled    bool
	mu         sync.Mutex
	isWatering bool
}

// NewHardwareController creates a new hardware controller.
func NewHardwareController(enabled bool) *HardwareController {
	return &HardwareController{enabled: enabled}
}

// SetWateringState turns the pump on or off.
func (h *HardwareController) SetWateringState(active bool) {
	if !h.enabled {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if active == h.isWatering {
		return
	}

	h.isWatering = active

	if active {
		fmt.Println("   💧 [PUMP]: ON - Watering plants...")
	} else {
		fmt.Println("   ⚡ [PUMP]: OFF")
	}
}

// Water runs the pump for the specified duration, printing progress.
func (h *HardwareController) Water(duration time.Duration) {
	if duration <= 0 {
		return
	}

	h.SetWateringState(true)

	if duration > 2*time.Second {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		elapsed := 0 * time.Second
		for elapsed < duration {
			<-ticker.C
			elapsed += time.Second
			progress := (elapsed * 100) / duration
			fmt.Printf("   ⏳ Watering progress: %d%%\r", progress)
		}
		fmt.Println()
	} else {
		time.Sleep(duration)
	}

	h.SetWateringState(false)
}
