package sdpinit

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func localPromptsCandidates() []string {
	candidates := []string{"prompts", "../prompts", "sdp/prompts", "../sdp/prompts"}
	// When running from repo (go run, or binary next to prompts), try executable-relative path
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		// sdp-plugin/cmd/sdp -> sdp/prompts or ../prompts
		candidates = append(candidates,
			filepath.Join(exeDir, "..", "..", "prompts"),
			filepath.Join(exeDir, "..", "prompts"),
		)
	}
	return candidates
}

func isValidPromptsDir(dir string) bool {
	if dir == "" {
		return false
	}
	info, err := os.Stat(filepath.Join(dir, "skills"))
	return err == nil && info.IsDir()
}

func envPromptsSourceDir() (string, error) {
	envDir := strings.TrimSpace(os.Getenv("SDP_PROMPTS_SOURCE_DIR"))
	if envDir == "" {
		return "", nil
	}
	if !isValidPromptsDir(envDir) {
		return "", fmt.Errorf("SDP_PROMPTS_SOURCE_DIR is invalid (missing skills/): %s", envDir)
	}
	return envDir, nil
}

func promptsArchiveURL() string {
	if custom := strings.TrimSpace(os.Getenv("SDP_PROMPTS_ARCHIVE_URL")); custom != "" {
		return custom
	}
	repo := strings.TrimSpace(os.Getenv("SDP_PROMPTS_REPO"))
	if repo == "" {
		repo = "fall-out-bug/sdp"
	}
	ref := strings.TrimSpace(os.Getenv("SDP_PROMPTS_REF"))
	if ref == "" {
		ref = "main"
	}
	if strings.HasPrefix(ref, "v") {
		return fmt.Sprintf("https://codeload.github.com/%s/zip/refs/tags/%s", repo, ref)
	}
	return fmt.Sprintf("https://codeload.github.com/%s/zip/refs/heads/%s", repo, ref)
}

func downloadPromptsToCache() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("resolve user cache dir: %w", err)
	}
	promptsRoot := filepath.Join(cacheDir, "sdp", "prompts")
	if isValidPromptsDir(promptsRoot) {
		return promptsRoot, nil
	}
	// Clear invalid/partial cache before re-download
	_ = os.RemoveAll(promptsRoot)
	if err := os.MkdirAll(promptsRoot, 0755); err != nil {
		return "", fmt.Errorf("create prompts cache dir: %w", err)
	}

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Get(promptsArchiveURL())
	if err != nil {
		return "", fmt.Errorf("download prompts archive: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download prompts archive failed: HTTP %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "sdp-prompts-*.zip")
	if err != nil {
		return "", fmt.Errorf("create temp prompts archive: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		_ = tmpFile.Close()
		return "", fmt.Errorf("write prompts archive: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("close prompts archive: %w", err)
	}

	if err := unzipPrompts(tmpPath, promptsRoot); err != nil {
		return "", err
	}
	if !isValidPromptsDir(promptsRoot) {
		return "", fmt.Errorf("extracted prompts cache is invalid: %s", promptsRoot)
	}
	return promptsRoot, nil
}

func unzipPrompts(zipPath, destRoot string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open prompts archive: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		rel, ok := promptsArchiveRelPath(f.Name)
		if !ok {
			continue
		}

		target, err := safePromptsTarget(destRoot, rel)
		if err != nil {
			return fmt.Errorf("invalid prompts archive path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("create prompts dir %s: %w", target, err)
			}
			continue
		}

		if err := extractArchiveFile(f, target); err != nil {
			return err
		}
	}
	return nil
}

func promptsArchiveRelPath(entryName string) (string, bool) {
	idx := strings.Index(entryName, "/prompts/")
	if idx < 0 {
		return "", false
	}
	rel := strings.TrimPrefix(entryName[idx+len("/prompts/"):], "/")
	if rel == "" {
		return "", false
	}
	return rel, true
}

func safePromptsTarget(destRoot, rel string) (string, error) {
	target := filepath.Join(destRoot, rel)
	if !strings.HasPrefix(target, filepath.Clean(destRoot)+string(os.PathSeparator)) {
		return "", fmt.Errorf("unsafe archive target: %s", target)
	}
	return target, nil
}

func extractArchiveFile(f *zip.File, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return fmt.Errorf("create prompts parent dir %s: %w", target, err)
	}

	src, err := f.Open()
	if err != nil {
		return fmt.Errorf("open archive entry %s: %w", f.Name, err)
	}
	defer src.Close()

	dst, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("create target file %s: %w", target, err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		_ = dst.Close()
		return fmt.Errorf("extract %s: %w", f.Name, err)
	}
	if err := dst.Close(); err != nil {
		return fmt.Errorf("close target file %s: %w", target, err)
	}
	return nil
}
