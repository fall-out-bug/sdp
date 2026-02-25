package executor

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestEvalsSchemaPathsResolve verifies schema paths in sdp/evals test_cases.jsonl exist (AC 00-053-14).
func TestEvalsSchemaPathsResolve(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	pluginDir := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	sdpRoot := filepath.Dir(pluginDir)
	evalsDir := filepath.Join(sdpRoot, "evals")
	if _, err := os.Stat(evalsDir); os.IsNotExist(err) {
		t.Skip("sdp/evals not found")
	}
	for _, sub := range []string{"idea", "build", "review"} {
		p := filepath.Join(evalsDir, sub, "test_cases.jsonl")
		if _, err := os.Stat(p); err != nil {
			continue
		}
		validateSchemaPathsInFile(t, sdpRoot, p)
	}
}

func validateSchemaPathsInFile(t *testing.T, baseDir, path string) {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := sc.Text()
		if line == "" {
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}
		schema, ok := obj["schema"].(string)
		if !ok || schema == "" {
			continue
		}
		resolved := filepath.Join(baseDir, schema)
		if _, err := os.Stat(resolved); os.IsNotExist(err) {
			t.Errorf("%s:%d: schema path %q does not exist (resolved: %s)", path, lineNum, schema, resolved)
		}
	}
}
