package doctor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// checkBeadsIntegrity validates the canonical Beads snapshot and tolerates
// a legacy/local beads.db artifact when no snapshot is present.
func checkBeadsIntegrity() DeepCheckResult {
	start := getTime()
	details := make(map[string]any)

	beadsPath, legacy := findBeadsStateArtifact()
	if beadsPath == "" {
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "warning",
			Duration: since(start),
			Message:  "Beads snapshot not found",
			Details:  details,
		}
	}

	artifact := filepath.Base(beadsPath)
	details["artifact"] = artifact
	details["legacy"] = legacy

	content, err := os.ReadFile(beadsPath)
	if err != nil {
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "error",
			Duration: since(start),
			Message:  fmt.Sprintf("Cannot read %s: %v", artifact, err),
			Details:  details,
		}
	}

	hash := sha256.Sum256(content)
	details["size"] = len(content)
	details["hash"] = hex.EncodeToString(hash[:8])

	if len(content) == 0 {
		message := "Beads snapshot is empty"
		if legacy {
			message = "Legacy beads.db is empty"
		}
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "warning",
			Duration: since(start),
			Message:  message,
			Details:  details,
		}
	}

	message := fmt.Sprintf("Snapshot OK (%d bytes)", len(content))
	if legacy {
		message = fmt.Sprintf("Legacy local DB OK (%d bytes)", len(content))
	}

	return DeepCheckResult{
		Check:    "Beads Integrity",
		Status:   "ok",
		Duration: since(start),
		Message:  message,
		Details:  details,
	}
}

// checkConfigVersion validates config version is compatible
func checkConfigVersion() DeepCheckResult {
	start := getTime()
	details := make(map[string]any)

	configPath := ".sdp/config.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DeepCheckResult{
			Check:    "Config Version",
			Status:   "ok",
			Duration: since(start),
			Message:  "No custom config (using defaults)",
			Details:  details,
		}
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return DeepCheckResult{
			Check:    "Config Version",
			Status:   "error",
			Duration: since(start),
			Message:  fmt.Sprintf("Cannot read config: %v", err),
			Details:  details,
		}
	}

	// Check for version field
	contentStr := string(content)
	if !strings.Contains(contentStr, "version:") {
		return DeepCheckResult{
			Check:    "Config Version",
			Status:   "warning",
			Duration: since(start),
			Message:  "Config missing version field",
			Details:  details,
		}
	}

	return DeepCheckResult{
		Check:    "Config Version",
		Status:   "ok",
		Duration: since(start),
		Message:  "Config version present",
		Details:  details,
	}
}
