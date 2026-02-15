package hooks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHookPack(t *testing.T) {
	pack := HookPack{
		Name:         "my-hook",
		Version:      "1.0.0",
		Events:       []string{"command:pre", "command:post"},
		Priority:     100,
		Dependencies: []string{"core>=1.0"},
		HandlerPath:  "./hooks/my-hook.js",
	}

	if pack.Name != "my-hook" {
		t.Errorf("expected name my-hook, got %s", pack.Name)
	}
	if len(pack.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(pack.Events))
	}
}

func TestHookPack_Validate(t *testing.T) {
	tests := []struct {
		name    string
		pack    HookPack
		wantErr bool
	}{
		{
			name: "valid pack",
			pack: HookPack{
				Name:    "valid-hook",
				Version: "1.0.0",
				Events:  []string{"command:pre"},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			pack: HookPack{
				Version: "1.0.0",
				Events:  []string{"command:pre"},
			},
			wantErr: true,
		},
		{
			name: "missing version",
			pack: HookPack{
				Name:   "hook",
				Events: []string{"command:pre"},
			},
			wantErr: true,
		},
		{
			name: "missing events",
			pack: HookPack{
				Name:    "hook",
				Version: "1.0.0",
			},
			wantErr: true,
		},
		{
			name: "invalid version format",
			pack: HookPack{
				Name:    "hook",
				Version: "invalid",
				Events:  []string{"command:pre"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pack.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscoverHooks(t *testing.T) {
	// Create temp directory with hook packs
	tempDir := t.TempDir()
	hooksDir := filepath.Join(tempDir, ".sdp", "hooks")
	os.MkdirAll(hooksDir, 0755)

	// Create valid hook pack
	validPack := `{
		"name": "test-hook",
		"version": "1.0.0",
		"events": ["command:pre"],
		"priority": 50,
		"handler_path": "./handler.js"
	}`
	os.WriteFile(filepath.Join(hooksDir, "test-hook.json"), []byte(validPack), 0644)

	// Create invalid hook pack (missing name)
	invalidPack := `{
		"version": "1.0.0",
		"events": ["command:pre"]
	}`
	os.WriteFile(filepath.Join(hooksDir, "invalid-hook.json"), []byte(invalidPack), 0644)

	packs, err := DiscoverHooks(hooksDir)
	if err != nil {
		t.Fatalf("DiscoverHooks failed: %v", err)
	}

	// Should only return valid packs
	if len(packs) != 1 {
		t.Errorf("expected 1 valid pack, got %d", len(packs))
	}

	if len(packs) > 0 && packs[0].Name != "test-hook" {
		t.Errorf("expected pack name test-hook, got %s", packs[0].Name)
	}
}

func TestDiscoverHooks_EmptyDir(t *testing.T) {
	tempDir := t.TempDir()
	hooksDir := filepath.Join(tempDir, ".sdp", "hooks")
	os.MkdirAll(hooksDir, 0755)

	packs, err := DiscoverHooks(hooksDir)
	if err != nil {
		t.Fatalf("DiscoverHooks failed: %v", err)
	}

	if len(packs) != 0 {
		t.Errorf("expected 0 packs from empty dir, got %d", len(packs))
	}
}

func TestDiscoverHooks_NonexistentDir(t *testing.T) {
	_, err := DiscoverHooks("/nonexistent/path")
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestHookPack_CheckEligibility(t *testing.T) {
	tests := []struct {
		name       string
		pack       HookPack
		installed  []HookPack
		eligible   bool
	}{
		{
			name: "no dependencies",
			pack: HookPack{
				Name:    "hook-a",
				Version: "1.0.0",
				Events:  []string{"command:pre"},
			},
			installed: []HookPack{},
			eligible:  true,
		},
		{
			name: "dependency satisfied",
			pack: HookPack{
				Name:         "hook-b",
				Version:      "1.0.0",
				Events:       []string{"command:pre"},
				Dependencies: []string{"hook-a>=1.0.0"},
			},
			installed: []HookPack{
				{Name: "hook-a", Version: "1.0.0"},
			},
			eligible: true,
		},
		{
			name: "dependency not satisfied",
			pack: HookPack{
				Name:         "hook-c",
				Version:      "1.0.0",
				Events:       []string{"command:pre"},
				Dependencies: []string{"hook-a>=2.0.0"},
			},
			installed: []HookPack{
				{Name: "hook-a", Version: "1.0.0"},
			},
			eligible: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pack.CheckEligibility(tt.installed)
			eligible := err == nil
			if eligible != tt.eligible {
				t.Errorf("CheckEligibility() = %v, want %v", eligible, tt.eligible)
			}
		})
	}
}

func TestListInstalledHooks(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	hooksDir := filepath.Join(tempDir, ".sdp", "hooks")
	os.MkdirAll(hooksDir, 0755)

	// Create installed hooks marker file
	installed := `[
		{"name": "hook-a", "version": "1.0.0"},
		{"name": "hook-b", "version": "2.0.0"}
	]`
	os.WriteFile(filepath.Join(hooksDir, "installed.json"), []byte(installed), 0644)

	packs, err := ListInstalledHooks(hooksDir)
	if err != nil {
		t.Fatalf("ListInstalledHooks failed: %v", err)
	}

	if len(packs) != 2 {
		t.Errorf("expected 2 installed hooks, got %d", len(packs))
	}
}

func TestInstallHook(t *testing.T) {
	tempDir := t.TempDir()
	hooksDir := filepath.Join(tempDir, ".sdp", "hooks")
	os.MkdirAll(hooksDir, 0755)

	pack := HookPack{
		Name:        "install-test",
		Version:     "1.0.0",
		Events:      []string{"command:pre"},
		Priority:    50,
		HandlerPath: "./handler.js",
	}

	err := InstallHook(hooksDir, pack)
	if err != nil {
		t.Fatalf("InstallHook failed: %v", err)
	}

	// Verify hook was added to installed.json
	installed, err := ListInstalledHooks(hooksDir)
	if err != nil {
		t.Fatalf("ListInstalledHooks failed: %v", err)
	}

	found := false
	for _, p := range installed {
		if p.Name == "install-test" {
			found = true
			break
		}
	}

	if !found {
		t.Error("hook was not found in installed list")
	}
}

func TestUninstallHook(t *testing.T) {
	tempDir := t.TempDir()
	hooksDir := filepath.Join(tempDir, ".sdp", "hooks")
	os.MkdirAll(hooksDir, 0755)

	// First install
	pack := HookPack{
		Name:    "uninstall-test",
		Version: "1.0.0",
		Events:  []string{"command:pre"},
	}

	InstallHook(hooksDir, pack)

	// Then uninstall
	err := UninstallHook(hooksDir, "uninstall-test")
	if err != nil {
		t.Fatalf("UninstallHook failed: %v", err)
	}

	// Verify hook was removed
	installed, _ := ListInstalledHooks(hooksDir)
	for _, p := range installed {
		if p.Name == "uninstall-test" {
			t.Error("hook should have been uninstalled")
		}
	}
}

func TestUninstallHook_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	hooksDir := filepath.Join(tempDir, ".sdp", "hooks")
	os.MkdirAll(hooksDir, 0755)

	err := UninstallHook(hooksDir, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent hook")
	}
}
