package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Global tracker instance
var globalTracker *Tracker
var trackerOnce sync.Once

// Tracker manages automatic telemetry tracking for CLI commands
type Tracker struct {
	collector      *Collector
	currentCommand *CommandEvent
	mu             sync.Mutex
}

// CommandEvent represents a command execution
type CommandEvent struct {
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	StartTime time.Time         `json:"start_time"`
	EndTime   time.Time         `json:"end_time"`
	Duration  time.Duration     `json:"duration"`
	Success   bool              `json:"success"`
	Error     string            `json:"error,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// GetTracker returns the global telemetry tracker instance
func GetTracker() *Tracker {
	trackerOnce.Do(func() {
		configDir, err := os.UserConfigDir()
		if err != nil {
			// If we can't get config dir, return a disabled tracker
			globalTracker = &Tracker{collector: &Collector{}}
			return
		}

		telemetryFile := filepath.Join(configDir, "sdp", "telemetry.jsonl")

		// Check if telemetry is enabled (opt-in)
		configPath := filepath.Join(configDir, "sdp", "telemetry.json")
		enabled := false // Opt-in: disabled by default
		if data, err := os.ReadFile(configPath); err == nil {
			var config map[string]bool
			if err := json.Unmarshal(data, &config); err == nil {
				if enabledVal, ok := config["enabled"]; ok && enabledVal {
					enabled = true
				}
			}
		}

		collector, err := NewCollector(telemetryFile, enabled)
		if err != nil {
			// If collector creation fails, return a disabled tracker
			globalTracker = &Tracker{collector: &Collector{}}
			return
		}

		globalTracker = &Tracker{collector: collector}
	})

	return globalTracker
}
