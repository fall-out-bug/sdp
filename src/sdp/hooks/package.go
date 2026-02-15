package hooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HookPack represents a packaged hook with metadata.
type HookPack struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Events       []string `json:"events"`
	Priority     int      `json:"priority"`
	Dependencies []string `json:"dependencies,omitempty"`
	HandlerPath  string   `json:"handler_path,omitempty"`
}

// Validate checks if the hook pack has required fields.
func (p HookPack) Validate() error {
	if p.Name == "" {
		return errors.New("hook name is required")
	}
	if p.Version == "" {
		return errors.New("hook version is required")
	}
	if len(p.Events) == 0 {
		return errors.New("at least one event is required")
	}
	// Basic semver check
	parts := strings.Split(p.Version, ".")
	if len(parts) < 2 {
		return errors.New("version must be in semver format (e.g., 1.0.0)")
	}
	return nil
}

// CheckEligibility verifies dependencies are satisfied.
func (p HookPack) CheckEligibility(installed []HookPack) error {
	for _, dep := range p.Dependencies {
		satisfied := false
		for _, inst := range installed {
			if checkDependency(inst, dep) {
				satisfied = true
				break
			}
		}
		if !satisfied {
			return fmt.Errorf("dependency not satisfied: %s", dep)
		}
	}
	return nil
}

// checkDependency checks if an installed pack satisfies a dependency string.
func checkDependency(pack HookPack, dep string) bool {
	// Simple check: "name>=version" format
	if strings.Contains(dep, ">=") {
		parts := strings.SplitN(dep, ">=", 2)
		name := strings.TrimSpace(parts[0])
		if pack.Name != name {
			return false
		}
		// For simplicity, just check exact version match
		// Full semver comparison would require additional library
		return pack.Version >= strings.TrimSpace(parts[1])
	}
	return pack.Name == dep
}

// DiscoverHooks finds all valid hook packs in a directory.
func DiscoverHooks(hooksDir string) ([]HookPack, error) {
	entries, err := os.ReadDir(hooksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read hooks directory: %w", err)
	}

	var packs []HookPack
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		path := filepath.Join(hooksDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var pack HookPack
		if err := json.Unmarshal(data, &pack); err != nil {
			continue
		}

		if err := pack.Validate(); err != nil {
			continue
		}

		packs = append(packs, pack)
	}

	return packs, nil
}

// installedFile is the name of the installed hooks manifest.
const installedFile = "installed.json"

// ListInstalledHooks returns hooks that have been installed.
func ListInstalledHooks(hooksDir string) ([]HookPack, error) {
	path := filepath.Join(hooksDir, installedFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []HookPack{}, nil
		}
		return nil, fmt.Errorf("failed to read installed hooks: %w", err)
	}

	var packs []HookPack
	if err := json.Unmarshal(data, &packs); err != nil {
		return nil, fmt.Errorf("failed to parse installed hooks: %w", err)
	}

	return packs, nil
}

// InstallHook adds a hook to the installed list.
func InstallHook(hooksDir string, pack HookPack) error {
	if err := pack.Validate(); err != nil {
		return fmt.Errorf("invalid hook pack: %w", err)
	}

	installed, err := ListInstalledHooks(hooksDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Check if already installed
	for i, p := range installed {
		if p.Name == pack.Name {
			installed[i] = pack // Update
			return saveInstalled(hooksDir, installed)
		}
	}

	// Add new
	installed = append(installed, pack)
	return saveInstalled(hooksDir, installed)
}

// UninstallHook removes a hook from the installed list.
func UninstallHook(hooksDir, name string) error {
	installed, err := ListInstalledHooks(hooksDir)
	if err != nil {
		return err
	}

	found := false
	var updated []HookPack
	for _, p := range installed {
		if p.Name == name {
			found = true
			continue
		}
		updated = append(updated, p)
	}

	if !found {
		return fmt.Errorf("hook %q not found", name)
	}

	return saveInstalled(hooksDir, updated)
}

// saveInstalled writes the installed hooks manifest.
func saveInstalled(hooksDir string, packs []HookPack) error {
	data, err := json.MarshalIndent(packs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal installed hooks: %w", err)
	}

	path := filepath.Join(hooksDir, installedFile)
	return os.WriteFile(path, data, 0644)
}
