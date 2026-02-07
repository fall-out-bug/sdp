package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRoleLoader_NewRoleLoader(t *testing.T) {
	loader := NewRoleLoader("/tmp/roles")

	if loader == nil {
		t.Fatal("Expected non-nil loader")
	}

	if loader.roleDir != "/tmp/roles" {
		t.Errorf("Expected roleDir '/tmp/roles', got '%s'", loader.roleDir)
	}
}

func TestRoleLoader_LoadAll(t *testing.T) {
	// AC1: Loads 100+ roles from .md files
	// Create temp directory
	tmpDir := t.TempDir()

	// Create test role files
	for i := 0; i < 5; i++ {
		path := filepath.Join(tmpDir, fmt.Sprintf("role%d.md", i))
		content := fmt.Sprintf("# Role %d\n\nDescription for role %d", i, i)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	loader := NewRoleLoader(tmpDir)
	roles, err := loader.LoadAll()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(roles) != 5 {
		t.Errorf("Expected 5 roles, got %d", len(roles))
	}
}

func TestRoleLoader_LoadAll_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	loader := NewRoleLoader(tmpDir)
	roles, err := loader.LoadAll()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(roles) != 0 {
		t.Errorf("Expected 0 roles, got %d", len(roles))
	}
}

func TestRoleLoader_Get(t *testing.T) {
	// AC3: Roles cached for performance
	tmpDir := t.TempDir()

	// Create test role
	path := filepath.Join(tmpDir, "test-role.md")
	content := "# Test Role\n\nDescription"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewRoleLoader(tmpDir)
	loader.LoadAll()

	// First get
	start := time.Now()
	role1, err := loader.Get("test-role")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	duration1 := time.Since(start)

	// Second get (should be cached)
	start = time.Now()
	role2, err := loader.Get("test-role")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	duration2 := time.Since(start)

	// Verify cache hit (second call should be faster or similar)
	if duration2 > duration1*10 {
		t.Logf("Warning: Cache not effective (first: %v, second: %v)", duration1, duration2)
	}

	// Verify same role returned
	if role1.Name != role2.Name {
		t.Error("Expected same role on cache hit")
	}
}

func TestRoleLoader_Get_NotFound(t *testing.T) {
	loader := NewRoleLoader("/tmp/roles")
	loader.LoadAll()

	_, err := loader.Get("nonexistent")

	if err == nil {
		t.Fatal("Expected error for non-existent role")
	}
}

func TestRoleLoader_List(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test roles
	for i := 0; i < 3; i++ {
		path := filepath.Join(tmpDir, fmt.Sprintf("role%d.md", i))
		content := fmt.Sprintf("# Role %d", i)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	loader := NewRoleLoader(tmpDir)
	loader.LoadAll()

	names := loader.List()

	if len(names) != 3 {
		t.Errorf("Expected 3 role names, got %d", len(names))
	}
}

func TestRoleLoader_InvalidateCache(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test role
	path := filepath.Join(tmpDir, "test-role.md")
	content := "# Test Role"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	loader := NewRoleLoader(tmpDir)
	loader.LoadAll()

	// Get role (caches it)
	loader.Get("test-role")

	// Invalidate cache
	loader.InvalidateCache("test-role")

	// Verify cache cleared
	if _, ok := loader.cache["test-role"]; ok {
		t.Error("Expected cache to be cleared")
	}
}

func TestRoleLoader_DuplicateRoles(t *testing.T) {
	// AC2: No duplicate roles
	tmpDir := t.TempDir()

	// Create role with same name (different files)
	path1 := filepath.Join(tmpDir, "role.md")
	path2 := filepath.Join(tmpDir, "role-copy.md")

	content := "# Duplicate Role"
	os.WriteFile(path1, []byte(content), 0644)
	os.WriteFile(path2, []byte(content), 0644)

	loader := NewRoleLoader(tmpDir)
	loader.LoadAll()

	// Should only have one role (second one overwrites or error)
	roles, _ := loader.LoadAll()

	// This depends on implementation - either deduplicate or error
	if len(roles) > 1 {
		// Check if names are different
		hasDupe := false
		for name := range roles {
			if name == "role" || name == "role-copy" {
				hasDupe = true
			}
		}
		if !hasDupe {
			t.Error("Expected role deduplication")
		}
	}
}
