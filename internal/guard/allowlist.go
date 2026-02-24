package guard

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultAllowlist contains dependency files that legitimately change across workstreams.
var DefaultAllowlist = []string{
	"go.sum",
	"go.mod",
	"package-lock.json",
	"yarn.lock",
}

// AllowlistConfig is the schema for .sdp/guard-allowlist.yaml.
type AllowlistConfig struct {
	Files []string `yaml:"files"`
}

// LoadAllowlist returns allowlist from .sdp/guard-allowlist.yaml, or default if absent.
func LoadAllowlist(projectRoot string) ([]string, error) {
	path := filepath.Join(projectRoot, ".sdp", "guard-allowlist.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultAllowlist, nil
		}
		return nil, err
	}
	var cfg AllowlistConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Files) == 0 {
		return DefaultAllowlist, nil
	}
	return cfg.Files, nil
}

// IsAllowlisted returns true if the file (relative path) is in the allowlist.
// Matches exact path or basename.
func IsAllowlisted(file string, allowlist []string) bool {
	base := filepath.Base(file)
	for _, a := range allowlist {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if file == a || base == a {
			return true
		}
	}
	return false
}
