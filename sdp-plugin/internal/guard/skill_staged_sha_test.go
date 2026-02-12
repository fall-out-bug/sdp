package guard

import (
	"os"
	"testing"
)

// TestParseCheckOptions_ValidatesSHA ensures CI_BASE_SHA and CI_HEAD_SHA
// are validated for correct format (40 hex characters) - sdp-67l6
func TestParseCheckOptions_ValidatesSHA(t *testing.T) {
	// Save original values
	origBase := os.Getenv("CI_BASE_SHA")
	origHead := os.Getenv("CI_HEAD_SHA")
	defer func() {
		if origBase != "" {
			os.Setenv("CI_BASE_SHA", origBase)
		} else {
			os.Unsetenv("CI_BASE_SHA")
		}
		if origHead != "" {
			os.Setenv("CI_HEAD_SHA", origHead)
		} else {
			os.Unsetenv("CI_HEAD_SHA")
		}
	}()

	tests := []struct {
		name     string
		base     string
		head     string
		wantBase  string // Expected base after validation
		wantHead  string // Expected head after validation
	}{
		{
			name:    "valid SHAs with non-hex",
			base:    "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2",
			head:    "b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w",
			wantBase: "", // Rejected due to non-hex (g, h)
			wantHead: "", // Rejected due to non-hex (w)
		},
		{
			name:    "valid 40-char hex SHAs",
			base:    "0123456789abcdef0123456789abcdef01234567",
			head:    "0123456789abcdef0123456789abcdef01234567",
			wantBase: "0123456789abcdef0123456789abcdef01234567", // Accepted
			wantHead: "0123456789abcdef0123456789abcdef01234567", // Accepted
		},
		{
			name:    "short SHA - too short",
			base:    "abc123",
			head:    "def456",
			wantBase: "", // Rejected
			wantHead: "", // Rejected
		},
		{
			name:    "empty strings - valid",
			base:    "",
			head:    "",
			wantBase: "", // Accepted
			wantHead: "", // Accepted
		},
		{
			name:    "injection attempt with semicolon",
			base:    "0123456789abc0123456789abc0123456789; rm -rf",
			head:    "",
			wantBase: "", // Rejected
			wantHead: "", // Accepted (empty)
		},
		{
			name:    "injection attempt with pipe",
			base:    "0123456789abc0123456789abc0123456789 | nc",
			head:    "",
			wantBase: "", // Rejected
			wantHead: "", // Accepted (empty)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.base != "" {
				os.Setenv("CI_BASE_SHA", tt.base)
			} else {
				os.Unsetenv("CI_BASE_SHA")
			}
			if tt.head != "" {
				os.Setenv("CI_HEAD_SHA", tt.head)
			} else {
				os.Unsetenv("CI_HEAD_SHA")
			}

			// Parse options - should validate SHA format
			opts := ParseCheckOptions()

			// Check validation worked correctly
			if opts.Base != tt.wantBase {
				t.Errorf("Base SHA = %q, want %q", opts.Base, tt.wantBase)
			}
			if opts.Head != tt.wantHead {
				t.Errorf("Head SHA = %q, want %q", opts.Head, tt.wantHead)
			}
		})
	}
}

// TestIsValidSHA tests SHA validation logic
func TestIsValidSHA(t *testing.T) {
	tests := []struct {
		sha   string
		valid bool
	}{
		{"", true},                                                         // Empty is valid
		{"abc123", false},                                                    // Too short
		{"ggg" + string(make([]byte, 37)) + "h", false},                       // Non-hex chars (g, h)
		{"a1b2c3d4e5f6g7h8", false},                                     // Contains g, h (non-hex)
		{"0123456789abcdef0123456789abcdef01234567", true},                 // Valid 40 hex
		{"0123456789012345678901234567890123456", false},                // Valid numeric but 41 chars
		{"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", true},                            // Valid uppercase
		{"ffffffffffffffffffffffffffffffffffffffff", true},                            // Valid lowercase
		{"GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG", false},                        // Invalid uppercase G
	}

	for _, tt := range tests {
		t.Run(shaLabel(tt.sha), func(t *testing.T) {
			got := isValidSHA(tt.sha)
			if got != tt.valid {
				t.Errorf("isValidSHA(%q) = %v, want %v", tt.sha, got, tt.valid)
			}
		})
	}
}

// TestIsHexChar tests hex character validation
func TestIsHexChar(t *testing.T) {
	tests := []struct {
		c     rune
		valid bool
	}{
		{'0', true}, {'1', true}, {'2', true}, {'3', true}, {'4', true}, {'5', true}, {'6', true}, {'7', true}, {'8', true}, {'9', true},
		{'a', true}, {'b', true}, {'c', true}, {'d', true}, {'e', true}, {'f', true},
		{'A', true}, {'B', true}, {'C', true}, {'D', true}, {'E', true}, {'F', true},
		{'g', false}, {'G', false}, {'z', false}, {'Z', false},
		{'@', false}, {'-', false}, {' ', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.c), func(t *testing.T) {
			got := isHexChar(tt.c)
			if got != tt.valid {
				t.Errorf("isHexChar(%q) = %v, want %v", tt.c, got, tt.valid)
			}
		})
	}
}

// shaLabel creates a short label for SHA in test output
func shaLabel(sha string) string {
	if len(sha) <= 10 {
		return sha
	}
	return sha[:10] + "..."
}
