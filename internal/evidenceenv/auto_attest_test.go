package evidenceenv

import "testing"

func TestMatchesAnyPrefix(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		prefixes []string
		want     bool
	}{
		{
			name:     "matches regular directory prefix",
			file:     "internal/orchestrate/hooks.go",
			prefixes: []string{"internal/orchestrate/"},
			want:     true,
		},
		{
			name:     "matches exact file path",
			file:     "internal/orchestrate/hooks.go",
			prefixes: []string{"internal/orchestrate/hooks.go"},
			want:     true,
		},
		{
			name:     "matches directory path without trailing slash",
			file:     "internal",
			prefixes: []string{"internal/"},
			want:     true,
		},
		{
			name:     "does not match unrelated path",
			file:     "cmd/sdp/main.go",
			prefixes: []string{"internal/"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesAnyPrefix(tt.file, tt.prefixes)
			if got != tt.want {
				t.Fatalf("matchesAnyPrefix(%q, %v) = %v, want %v", tt.file, tt.prefixes, got, tt.want)
			}
		})
	}
}
