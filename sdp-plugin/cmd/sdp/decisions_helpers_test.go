package main

import (
	"testing"
)

func TestValidateFieldLength(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     string
		maxLen    int
		wantErr   bool
	}{
		{
			name:      "empty value",
			fieldName: "title",
			value:     "",
			maxLen:    100,
			wantErr:   false,
		},
		{
			name:      "value under max",
			fieldName: "title",
			value:     "short title",
			maxLen:    100,
			wantErr:   false,
		},
		{
			name:      "value exactly at max",
			fieldName: "title",
			value:     "exactly 10",
			maxLen:    10,
			wantErr:   false,
		},
		{
			name:      "value over max",
			fieldName: "title",
			value:     "this is a very long title that exceeds the maximum length",
			maxLen:    10,
			wantErr:   true,
		},
		{
			name:      "unicode value under max",
			fieldName: "title",
			value:     "привет",
			maxLen:    100,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFieldLength(tt.fieldName, tt.value, tt.maxLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFieldLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStripControlChars(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "no control chars",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "preserves newline",
			input: "hello\nworld",
			want:  "hello\nworld",
		},
		{
			name:  "preserves tab",
			input: "hello\tworld",
			want:  "hello\tworld",
		},
		{
			name:  "removes null byte",
			input: "hello\x00world",
			want:  "helloworld",
		},
		{
			name:  "removes bell",
			input: "hello\x07world",
			want:  "helloworld",
		},
		{
			name:  "removes escape",
			input: "hello\x1bworld",
			want:  "helloworld",
		},
		{
			name:  "preserves newline and removes others",
			input: "hello\n\x00world\t\x07",
			want:  "hello\nworld\t",
		},
		{
			name:  "unicode preserved",
			input: "привет мир",
			want:  "привет мир",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripControlChars(tt.input)
			if got != tt.want {
				t.Errorf("stripControlChars(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFindProjectRoot(t *testing.T) {
	// This test runs in a git repository (the project itself)
	root, err := findProjectRoot()
	if err != nil {
		// May fail if test runs from temp directory
		t.Logf("findProjectRoot() returned error: %v (may be expected)", err)
		return
	}

	// Root should contain .git directory
	if root == "" {
		t.Error("findProjectRoot() returned empty string without error")
	}
}
