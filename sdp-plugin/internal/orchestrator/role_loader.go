package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Role represents a loaded role
type Role struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Prompt      string            `json:"prompt"`
	Permissions []string          `json:"permissions"`
	Metadata    map[string]string `json:"metadata"`
	FilePath    string            `json:"file_path"`
}

// RoleLoader handles loading and caching roles
type RoleLoader struct {
	mu        sync.RWMutex
	roles     map[string]*Role
	cache     map[string]*Role
	cacheTime map[string]time.Time
	roleDir   string
}

// NewRoleLoader creates a new role loader
func NewRoleLoader(roleDir string) *RoleLoader {
	return &RoleLoader{
		roles:     make(map[string]*Role),
		cache:     make(map[string]*Role),
		cacheTime: make(map[string]time.Time),
		roleDir:   roleDir,
	}
}

// LoadAll loads all roles from the role directory
func (rl *RoleLoader) LoadAll() (map[string]*Role, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if directory exists
	if _, err := os.Stat(rl.roleDir); os.IsNotExist(err) {
		return rl.roles, nil // No roles to load
	}

	// Read all .md files
	files, err := filepath.Glob(filepath.Join(rl.roleDir, "*.md"))
	if err != nil {
		return nil, fmt.Errorf("failed to read role files: %w", err)
	}

	// Load each role
	for _, file := range files {
		role, err := rl.loadRole(file)
		if err != nil {
			return nil, fmt.Errorf("failed to load role from %s: %w", file, err)
		}

		// Check for duplicates
		if _, exists := rl.roles[role.Name]; exists {
			return nil, fmt.Errorf("duplicate role: %s", role.Name)
		}

		rl.roles[role.Name] = role
	}

	return rl.roles, nil
}

// loadRole loads a single role from a file
func (rl *RoleLoader) loadRole(filePath string) (*Role, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Simple parsing - in production, use proper markdown frontmatter parser
	role := &Role{
		Name:        filepath.Base(filePath[:len(filePath)-3]),
		Description: string(content),
		FilePath:    filePath,
	}

	return role, nil
}

// Get retrieves a role by name
func (rl *RoleLoader) Get(name string) (*Role, error) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	// Check cache
	if role, ok := rl.cache[name]; ok {
		// Check cache age (invalidate after 5 minutes)
		if time.Since(rl.cacheTime[name]) < 5*time.Minute {
			return role, nil
		}
	}

	// Load from roles map
	role, ok := rl.roles[name]
	if !ok {
		return nil, fmt.Errorf("role not found: %s", name)
	}

	// Update cache
	rl.cache[name] = role
	rl.cacheTime[name] = time.Now()

	return role, nil
}

// InvalidateCache invalidates the cache for a role
func (rl *RoleLoader) InvalidateCache(name string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.cache, name)
	delete(rl.cacheTime, name)
}

// List returns all role names
func (rl *RoleLoader) List() []string {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	names := make([]string, 0, len(rl.roles))
	for name := range rl.roles {
		names = append(names, name)
	}
	return names
}
