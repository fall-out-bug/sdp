package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstall(t *testing.T) {
	// Create temporary directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Check that hooks were created
	expectedHooks := []string{"pre-commit", "pre-push"}
	for _, hookName := range expectedHooks {
		hookPath := filepath.Join(gitDir, hookName)
		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			t.Errorf("Hook %s was not created", hookName)
			continue
		}

		// Check content
		content, err := os.ReadFile(hookPath)
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", hookPath, err)
		}

		if !strings.Contains(string(content), "# SDP Git Hook") {
			t.Errorf("Hook %s has wrong content: %s", hookName, string(content))
		}
	}
}

func TestInstall_NoGitDir(t *testing.T) {
	// Create temp directory WITHOUT .git
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install - should fail
	err := Install()
	if err == nil {
		t.Fatal("Install() should fail when .git doesn't exist")
	}

	if !strings.Contains(err.Error(), ".git directory not found") {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestInstall_SkipExisting(t *testing.T) {
	// Create temp directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create existing hook
	existingHook := filepath.Join(gitDir, "pre-commit")
	if err := os.WriteFile(existingHook, []byte("# existing hook"), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Run install - should skip existing hook
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Check that existing hook wasn't overwritten
	content, err := os.ReadFile(existingHook)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(content) != "# existing hook" {
		t.Errorf("Hook was overwritten, got: %s", string(content))
	}
}

func TestUninstall(t *testing.T) {
	// Create temp directory with hooks
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create SDP-managed hooks
	for _, hookName := range []string{"pre-commit", "pre-push"} {
		hookPath := filepath.Join(gitDir, hookName)
		// Create hook with SDP marker so uninstall will remove it
		if err := os.WriteFile(hookPath, []byte(sdpManagedMarker+"\n# test hook"), 0755); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	// Run uninstall
	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}

	// Check that hooks were removed
	for _, hookName := range []string{"pre-commit", "pre-push"} {
		hookPath := filepath.Join(gitDir, hookName)
		if _, err := os.Stat(hookPath); !os.IsNotExist(err) {
			t.Errorf("Hook %s was not removed", hookName)
		}
	}
}

func TestUninstall_NotExists(t *testing.T) {
	// Create temp directory WITHOUT hooks
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run uninstall - should not fail
	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() should not fail when hooks don't exist: %v", err)
	}
}

// TestAC1_InstallsAllCanonicalHooks tests that install creates all 4 canonical SDP hooks
func TestAC1_InstallsAllCanonicalHooks(t *testing.T) {
	// Create temporary directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Check that all canonical hooks were created
	expectedHooks := []string{"pre-commit", "pre-push", "post-merge", "post-checkout"}
	for _, hookName := range expectedHooks {
		hookPath := filepath.Join(gitDir, hookName)
		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			t.Errorf("Hook %s was not created", hookName)
			continue
		}

		// Verify hook is executable
		info, err := os.Stat(hookPath)
		if err != nil {
			t.Fatalf("Stat(%s): %v", hookPath, err)
		}
		if info.Mode().Perm()&0111 == 0 {
			t.Errorf("Hook %s is not executable", hookName)
		}
	}
}

// TestAC2_HooksContainSDPMarker tests that installed hooks have the SDP-MANAGED-HOOK marker
func TestAC2_HooksContainSDPMarker(t *testing.T) {
	// Create temporary directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Check that all hooks have the SDP marker
	expectedHooks := []string{"pre-commit", "pre-push", "post-merge", "post-checkout"}
	for _, hookName := range expectedHooks {
		hookPath := filepath.Join(gitDir, hookName)

		content, err := os.ReadFile(hookPath)
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", hookPath, err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, sdpManagedMarker) {
			t.Errorf("Hook %s does not contain SDP marker '%s'. Got: %s", hookName, sdpManagedMarker, contentStr)
		}

		// Verify marker is on the second line (after shebang)
		lines := strings.Split(contentStr, "\n")
		if len(lines) < 2 || !strings.Contains(lines[1], sdpManagedMarker) {
			t.Errorf("Hook %s does not have SDP marker on second line. First two lines: %q, %q", hookName, lines[0], lines[1])
		}
	}
}

// TestAC3_InstallIsIdempotent tests that running install multiple times doesn't duplicate hooks
func TestAC3_InstallIsIdempotent(t *testing.T) {
	// Create temporary directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install twice
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Read the hooks after first install
	expectedHooks := []string{"pre-commit", "pre-push", "post-merge", "post-checkout"}
	firstInstallContent := make(map[string]string)
	for _, hookName := range expectedHooks {
		hookPath := filepath.Join(gitDir, hookName)
		content, err := os.ReadFile(hookPath)
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", hookPath, err)
		}
		firstInstallContent[hookName] = string(content)
	}

	// Run install again
	if err := Install(); err != nil {
		t.Fatalf("Second Install() failed: %v", err)
	}

	// Verify hooks weren't duplicated or modified
	for _, hookName := range expectedHooks {
		hookPath := filepath.Join(gitDir, hookName)
		content, err := os.ReadFile(hookPath)
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", hookPath, err)
		}

		if string(content) != firstInstallContent[hookName] {
			t.Errorf("Hook %s was modified on second install. Expected:\n%s\n\nGot:\n%s",
				hookName, firstInstallContent[hookName], string(content))
		}
	}
}

// TestAC4_UninstallOnlyRemovesSDPHooks tests that uninstall only removes SDP-managed hooks
func TestAC4_UninstallOnlyRemovesSDPHooks(t *testing.T) {
	// Create temporary directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Install SDP hooks
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Create a non-SDP hook
	nonSDPHookPath := filepath.Join(gitDir, "pre-commit")
	nonSDPHookContent := `#!/bin/sh
# My custom hook
# Not managed by SDP
echo "Custom check"
`
	// First, let's rename the SDP hook and create a non-SDP one
	sdpHookPath := filepath.Join(gitDir, "pre-commit-sdp")
	if err := os.Rename(nonSDPHookPath, sdpHookPath); err != nil {
		t.Fatalf("Rename failed: %v", err)
	}
	if err := os.WriteFile(nonSDPHookPath, []byte(nonSDPHookContent), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Also create a custom post-commit hook (not managed by SDP)
	customHookPath := filepath.Join(gitDir, "post-commit")
	customHookContent := `#!/bin/sh
# Custom post-commit
echo "Post commit action"
`
	if err := os.WriteFile(customHookPath, []byte(customHookContent), 0755); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Run uninstall
	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall() failed: %v", err)
	}

	// Verify non-SDP hooks still exist
	if _, err := os.Stat(customHookPath); os.IsNotExist(err) {
		t.Errorf("Non-SDP hook post-commit was removed but should have been preserved")
	}

	// Read the pre-commit hook - should be the non-SDP one (since uninstall skips non-SDP)
	content, err := os.ReadFile(nonSDPHookPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(content) != nonSDPHookContent {
		t.Errorf("Non-SDP pre-commit hook was modified. Expected:\n%s\n\nGot:\n%s",
			nonSDPHookContent, string(content))
	}
}

// TestAC5_HooksArePOSISSafe tests that hooks use /bin/sh and avoid bashisms
func TestAC5_HooksArePOSISSafe(t *testing.T) {
	// Create temporary directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Check that all hooks use /bin/sh
	expectedHooks := []string{"pre-commit", "pre-push", "post-merge", "post-checkout"}
	for _, hookName := range expectedHooks {
		hookPath := filepath.Join(gitDir, hookName)

		content, err := os.ReadFile(hookPath)
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", hookPath, err)
		}

		contentStr := string(content)

		// Check for /bin/sh shebang
		if !strings.Contains(contentStr, "#!/bin/sh") {
			t.Errorf("Hook %s does not use #!/bin/sh shebang. Got: %s", hookName, contentStr)
		}

		// Check for bashisms that should NOT be present
		bashisms := []string{
			"[[",        // bash-only test
			"$((",       // bash arithmetic
			"==",        // bash string comparison (should be =)
			"source ",   // bash keyword (should be .)
			"function ", // bash keyword
			"echo -n",   // not portable
			"echo -e",   // not portable
		}

		for _, bashism := range bashisms {
			if strings.Contains(contentStr, bashism) {
				t.Errorf("Hook %s contains bashism '%s': %s", hookName, bashism, contentStr)
			}
		}
	}
}

// TestAC6_HooksWarnWhenSDPMissing tests that hooks check for sdp binary and warn if missing
func TestAC6_HooksWarnWhenSDPMissing(t *testing.T) {
	// Create temporary directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git", "hooks")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run install
	if err := Install(); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	// Check that hooks include a warning when sdp is not found
	expectedHooks := []string{"pre-commit", "pre-push", "post-merge", "post-checkout"}
	for _, hookName := range expectedHooks {
		hookPath := filepath.Join(gitDir, hookName)

		content, err := os.ReadFile(hookPath)
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", hookPath, err)
		}

		contentStr := string(content)

		// Check for command existence check (POSIX style: command -v or which)
		hasCommandCheck := strings.Contains(contentStr, "command -v sdp") ||
			strings.Contains(contentStr, "which sdp") ||
			strings.Contains(contentStr, "type sdp")

		if !hasCommandCheck {
			t.Errorf("Hook %s does not check for sdp binary existence. Got: %s", hookName, contentStr)
		}

		// Check for warning message when sdp is not found
		hasWarning := strings.Contains(contentStr, "not found") ||
			strings.Contains(contentStr, "not installed") ||
			strings.Contains(contentStr, "warning") ||
			strings.Contains(contentStr, "Warning")

		if !hasWarning {
			t.Errorf("Hook %s does not include warning when sdp is missing. Got: %s", hookName, contentStr)
		}
	}
}
