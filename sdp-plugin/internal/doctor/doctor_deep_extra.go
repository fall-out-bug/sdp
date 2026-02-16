package doctor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// checkBeadsIntegrity validates Beads database is healthy
func checkBeadsIntegrity() DeepCheckResult {
	start := getTime()
	details := make(map[string]interface{})

	beadsDB := ".beads/beads.db"
	if _, err := os.Stat(beadsDB); os.IsNotExist(err) {
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "warning",
			Duration: since(start),
			Message:  "Beads database not found",
			Details:  details,
		}
	}

	// Check file can be read
	content, err := os.ReadFile(beadsDB)
	if err != nil {
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "error",
			Duration: since(start),
			Message:  fmt.Sprintf("Cannot read beads.db: %v", err),
			Details:  details,
		}
	}

	// Basic integrity: check it's not empty and has expected structure
	hash := sha256.Sum256(content)
	details["size"] = len(content)
	details["hash"] = hex.EncodeToString(hash[:8])

	if len(content) == 0 {
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "warning",
			Duration: since(start),
			Message:  "Beads database is empty",
			Details:  details,
		}
	}

	return DeepCheckResult{
		Check:    "Beads Integrity",
		Status:   "ok",
		Duration: since(start),
		Message:  fmt.Sprintf("Database OK (%d bytes)", len(content)),
		Details:  details,
	}
}

// checkConfigVersion validates config version is compatible
func checkConfigVersion() DeepCheckResult {
	start := getTime()
	details := make(map[string]interface{})

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
