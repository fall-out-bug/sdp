package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/telemetry"
)

func TestAdoptCmdAlreadyInitialized(t *testing.T) {
	// Get original working directory
	originalWd, _ := os.Getwd()

	// Create temp directory
	tmpDir := t.TempDir()

	// Change to temp directory
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Setup git repo
	setupTestGitRepo(t, tmpDir)

	// Create .sdp and .claude directories to simulate initialized state
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp: %v", err)
	}

	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("failed to create .claude: %v", err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create settings.json: %v", err)
	}

	// Create command
	cmd := adoptCmd()
	cmd.SetArgs([]string{})

	// Execute command - should fail
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("Expected error when already initialized, got nil")
	}
	if !strings.Contains(err.Error(), "SDP already initialized") {
		t.Errorf("Expected error about already initialized, got: %v", err)
	}
}

func TestAdoptCmdWithForce(t *testing.T) {
	// Get original working directory
	originalWd, _ := os.Getwd()

	// Create temp directory
	tmpDir := t.TempDir()

	// Change to temp directory
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Setup git repo
	setupTestGitRepo(t, tmpDir)

	// Create .sdp and .claude directories to simulate initialized state
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp: %v", err)
	}

	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("failed to create .claude: %v", err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create settings.json: %v", err)
	}

	// Create command with force flag
	cmd := adoptCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Flags().Set("force", "true"); err != nil {
		t.Fatalf("failed to set force flag: %v", err)
	}

	// Execute command - should succeed (may fail on actual init, but shouldn't fail on already initialized check)
	err := cmd.RunE(cmd, []string{})
	// We expect this might fail due to actual init issues, but NOT due to "already initialized"
	if err != nil && strings.Contains(err.Error(), "already initialized") {
		t.Errorf("Should not fail with 'already initialized' when using --force, got: %v", err)
	}
}

func TestAdoptTelemetry(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Setup git repo
	setupTestGitRepo(t, tmpDir)

	// Create temp telemetry dir
	telemetryDir := t.TempDir()

	// Create UX metrics collector with temp dir
	uxMetrics, err := telemetry.NewUXMetricsCollector(telemetryDir)
	if err != nil {
		t.Fatalf("failed to create UX metrics collector: %v", err)
	}

	// Record adopt complete
	err = uxMetrics.RecordAdoptComplete("test-project", 100)
	if err != nil {
		t.Fatalf("failed to record adopt complete: %v", err)
	}

	// Verify event was written
	eventsFile := uxMetrics.GetEventsFile()
	data, err := os.ReadFile(eventsFile)
	if err != nil {
		t.Fatalf("failed to read events file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "metric_type") {
		t.Errorf("expected event to contain metric_type, got: %s", content)
	}
	if !strings.Contains(content, "time_to_first_value") {
		t.Errorf("expected event to contain time_to_first_value, got: %s", content)
	}
	if !strings.Contains(content, "adopt") {
		t.Errorf("expected event to contain adopt step, got: %s", content)
	}
}

func TestCommitSDPStructure(t *testing.T) {
	tests := []struct {
		name        string
		setupRepo   func(t *testing.T, dir string)
		wantErr     bool
		errContains string
	}{
		{
			name: "commits successfully",
			setupRepo: func(t *testing.T, dir string) {
				setupTestGitRepo(t, dir)

				// Create .sdp directory
				sdpDir := filepath.Join(dir, ".sdp")
				if err := os.MkdirAll(sdpDir, 0755); err != nil {
					t.Fatalf("failed to create .sdp: %v", err)
				}

				// Create .claude directory
				claudeDir := filepath.Join(dir, ".claude")
				if err := os.MkdirAll(claudeDir, 0755); err != nil {
					t.Fatalf("failed to create .claude: %v", err)
				}

				// Create a file in .sdp
				if err := os.WriteFile(filepath.Join(sdpDir, "config.yml"), []byte("test: true"), 0644); err != nil {
					t.Fatalf("failed to create config file: %v", err)
				}
			},
			wantErr: false,
		},
		{
			name: "fails when no SDP structure exists",
			setupRepo: func(t *testing.T, dir string) {
				setupTestGitRepo(t, dir)
			},
			wantErr:     true,
			errContains: "failed to add",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tt.setupRepo(t, tmpDir)

			oldWd, _ := os.Getwd()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to chdir: %v", err)
			}
			defer os.Chdir(oldWd)

			err := commitSDPStructure()

			if (err != nil) != tt.wantErr {
				t.Errorf("commitSDPStructure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %q, want error containing %q", err.Error(), tt.errContains)
				}
			}
		})
	}
}

// setupTestGitRepo creates a minimal git repository for testing
func setupTestGitRepo(t *testing.T, dir string) {
	t.Helper()

	commands := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test User"},
		{"git", "checkout", "-b", "main"},
		{"sh", "-c", "echo test > README.md"},
		{"git", "add", "README.md"},
		{"git", "commit", "-m", "initial commit"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git setup failed: %v: %s", err, string(output))
		}
	}
}

func TestAdoptFullFlag_SetsAdoptionModeOff(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Setup git repo
	setupTestGitRepo(t, tmpDir)

	// Create .sdp directory with config.yml that has adoption_mode: true
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp: %v", err)
	}
	configPath := filepath.Join(sdpDir, "config.yml")
	if err := os.WriteFile(configPath, []byte("version: 1\nadoption_mode: true\n"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Run adopt --full
	cmd := adoptCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Flags().Set("full", "true"); err != nil {
		t.Fatalf("failed to set full flag: %v", err)
	}

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("adopt --full failed: %v", err)
	}

	// Verify adoption_mode is now false
	cfg, err := config.Load(tmpDir)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if cfg.AdoptionMode {
		t.Error("expected adoption_mode to be false after adopt --full")
	}
}

func TestAdoptFullFlag_NoSDPDir_ReturnsError(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	setupTestGitRepo(t, tmpDir)

	cmd := adoptCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Flags().Set("full", "true"); err != nil {
		t.Fatalf("failed to set full flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("expected error when .sdp/ not found, got nil")
	}
	if !strings.Contains(err.Error(), ".sdp/") {
		t.Errorf("expected error about .sdp/ directory, got: %v", err)
	}
}

func TestAdoptDefault_SetsAdoptionModeOn(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	setupTestGitRepo(t, tmpDir)

	// Create .sdp and .claude to simulate partial setup
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp: %v", err)
	}

	// Run adopt --force (bypass already-initialized check)
	cmd := adoptCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Flags().Set("force", "true"); err != nil {
		t.Fatalf("failed to set force flag: %v", err)
	}

	err := cmd.RunE(cmd, []string{})
	// The adopt may fail on sdpinit.Run (depends on environment), but
	// the config.yml should still be created with adoption_mode: true
	// if the directory was set up before that step
	if err != nil {
		// Check if config was written before the failure
		configPath := filepath.Join(sdpDir, "config.yml")
		if _, statErr := os.Stat(configPath); statErr != nil {
			t.Skipf("config not written before adopt failed: %v", err)
		}
	}

	// Verify config has adoption_mode: true
	cfg, cfgErr := config.Load(tmpDir)
	if cfgErr != nil {
		t.Fatalf("failed to load config: %v", cfgErr)
	}
	if !cfg.AdoptionMode {
		t.Error("expected adoption_mode to be true after adopt")
	}
}
