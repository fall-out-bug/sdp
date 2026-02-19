package sdpinit

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPromptsArchiveURL(t *testing.T) {
	reset := func(key, value string) {
		_ = os.Setenv(key, value)
	}

	origCustom := os.Getenv("SDP_PROMPTS_ARCHIVE_URL")
	origRepo := os.Getenv("SDP_PROMPTS_REPO")
	origRef := os.Getenv("SDP_PROMPTS_REF")
	t.Cleanup(func() {
		reset("SDP_PROMPTS_ARCHIVE_URL", origCustom)
		reset("SDP_PROMPTS_REPO", origRepo)
		reset("SDP_PROMPTS_REF", origRef)
	})

	_ = os.Setenv("SDP_PROMPTS_ARCHIVE_URL", "https://example.com/prompts.zip")
	if got := promptsArchiveURL(); got != "https://example.com/prompts.zip" {
		t.Fatalf("custom archive URL mismatch: %s", got)
	}

	_ = os.Setenv("SDP_PROMPTS_ARCHIVE_URL", "")
	_ = os.Setenv("SDP_PROMPTS_REPO", "acme/sdp")
	_ = os.Setenv("SDP_PROMPTS_REF", "main")
	if got := promptsArchiveURL(); got != "https://codeload.github.com/acme/sdp/zip/refs/heads/main" {
		t.Fatalf("head ref URL mismatch: %s", got)
	}

	_ = os.Setenv("SDP_PROMPTS_REF", "v1.2.3")
	if got := promptsArchiveURL(); got != "https://codeload.github.com/acme/sdp/zip/refs/tags/v1.2.3" {
		t.Fatalf("tag ref URL mismatch: %s", got)
	}
}

func TestPromptsArchiveRelPathAndSafeTarget(t *testing.T) {
	rel, ok := promptsArchiveRelPath("repo-main/prompts/skills/test.md")
	if !ok || rel != "skills/test.md" {
		t.Fatalf("unexpected rel path parse: ok=%v rel=%q", ok, rel)
	}

	if _, ok := promptsArchiveRelPath("repo-main/README.md"); ok {
		t.Fatalf("expected non-prompts path to be ignored")
	}

	root := t.TempDir()
	if _, err := safePromptsTarget(root, "skills/test.md"); err != nil {
		t.Fatalf("expected safe target: %v", err)
	}
	if _, err := safePromptsTarget(root, "../escape"); err == nil {
		t.Fatalf("expected unsafe target error")
	}
}

func TestUnzipPrompts(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "prompts.zip")
	if err := writeZip(zipPath, map[string]string{
		"repo-main/prompts/skills/a.md": "# A",
		"repo-main/prompts/agents/b.md": "# B",
		"repo-main/README.md":           "ignore",
	}); err != nil {
		t.Fatalf("write zip: %v", err)
	}

	dest := filepath.Join(tmp, "out")
	if err := unzipPrompts(zipPath, dest); err != nil {
		t.Fatalf("unzipPrompts failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dest, "skills", "a.md")); err != nil {
		t.Fatalf("missing extracted skill: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "agents", "b.md")); err != nil {
		t.Fatalf("missing extracted agent: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "README.md")); err == nil {
		t.Fatalf("non-prompts file should not be extracted")
	}
}

func TestDownloadPromptsToCache_UsesExistingCache(t *testing.T) {
	cacheHome := t.TempDir()
	origXDG := os.Getenv("XDG_CACHE_HOME")
	t.Cleanup(func() { _ = os.Setenv("XDG_CACHE_HOME", origXDG) })
	if err := os.Setenv("XDG_CACHE_HOME", cacheHome); err != nil {
		t.Fatalf("set XDG_CACHE_HOME: %v", err)
	}

	promptsDir := filepath.Join(cacheHome, "sdp", "prompts", "skills")
	if err := os.MkdirAll(promptsDir, 0755); err != nil {
		t.Fatalf("mkdir prompts cache: %v", err)
	}
	if err := os.WriteFile(filepath.Join(promptsDir, "cached.md"), []byte("# cached"), 0644); err != nil {
		t.Fatalf("write cached prompt: %v", err)
	}

	resolved, err := downloadPromptsToCache()
	if err != nil {
		t.Fatalf("downloadPromptsToCache should use cache: %v", err)
	}
	if !isValidPromptsDir(resolved) {
		t.Fatalf("resolved cache path should contain skills/: %s", resolved)
	}
	if !strings.HasSuffix(resolved, filepath.Join("sdp", "prompts")) {
		t.Fatalf("resolved cache path should end with sdp/prompts: %s", resolved)
	}
}

func writeZip(path string, files map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	for name, content := range files {
		entry, err := w.Create(name)
		if err != nil {
			_ = w.Close()
			return err
		}
		if _, err := entry.Write([]byte(content)); err != nil {
			_ = w.Close()
			return err
		}
	}
	return w.Close()
}
