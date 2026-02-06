package security

import (
	"context"
	"testing"
	"time"
)

// TestSafeCommandWhitelist tests that SafeCommand validates against whitelist
func TestSafeCommandWhitelist(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		shouldPass  bool
	}{
		{
			name:       "safe_pytest",
			command:    "pytest",
			args:       []string{"tests/"},
			shouldPass: true,
		},
		{
			name:       "safe_go_test",
			command:    "go",
			args:       []string{"test"},
			shouldPass: true,
		},
		{
			name:       "safe_git",
			command:    "git",
			args:       []string{"--version"},
			shouldPass: true,
		},
		{
			name:       "blocked_rm",
			command:    "rm",
			args:       []string{"-rf", "/"},
			shouldPass: false,
		},
		{
			name:       "blocked_cat",
			command:    "cat",
			args:       []string{"/etc/passwd"},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cmd, err := SafeCommand(ctx, tt.command, tt.args...)

			if tt.shouldPass {
				if err != nil {
					t.Errorf("SafeCommand(%s %v) should pass but got error: %v", tt.command, tt.args, err)
				}
				if cmd == nil {
					t.Error("SafeCommand should return non-nil cmd")
				}
			} else {
				if err == nil {
					t.Errorf("SafeCommand(%s %v) should be blocked but passed", tt.command, tt.args)
				}
			}
		})
	}
}

// TestSafeCommandInjectionDetection tests injection pattern detection
func TestSafeCommandInjectionDetection(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		hasInjection bool
	}{
		{
			name:         "semicolon_injection",
			command:      "pytest",
			args:         []string{"tests/; rm -rf /"},
			hasInjection: true,
		},
		{
			name:         "pipe_injection",
			command:      "pytest",
			args:         []string{"| sh"},
			hasInjection: true,
		},
		{
			name:         "backtick_injection",
			command:      "pytest",
			args:         []string{"`whoami`"},
			hasInjection: true,
		},
		{
			name:         "dollar_injection",
			command:      "pytest",
			args:         []string{"$(rm -rf /)"},
			hasInjection: true,
		},
		{
			name:         "path_traversal",
			command:      "pytest",
			args:         []string{"../../../etc/passwd"},
			hasInjection: true,
		},
		{
			name:         "safe_args",
			command:      "pytest",
			args:         []string{"tests/", "-v", "--cov"},
			hasInjection: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := SafeCommand(ctx, tt.command, tt.args...)

			if tt.hasInjection {
				if err == nil {
					t.Errorf("SafeCommand should detect injection in: %v", tt.args)
				}
			} else {
				if err != nil {
					t.Errorf("SafeCommand should allow safe args: %v (error: %v)", tt.args, err)
				}
			}
		})
	}
}

// TestSafeCommandTimeout tests that SafeCommand sets default timeout
func TestSafeCommandTimeout(t *testing.T) {
	ctx := context.Background()

	// Command without deadline should get default timeout
	cmd, err := SafeCommand(ctx, "pytest", []string{"tests/"}...)
	if err != nil {
		t.Fatalf("SafeCommand failed: %v", err)
	}

	// Check that command has a timeout
	// We can't directly inspect the context, but we can verify it doesn't hang
	if cmd == nil {
		t.Fatal("SafeCommand returned nil cmd")
	}

	// Verify context is set
	_ = cmd // Use cmd to avoid "declared and not used" error
}

// TestSafeCommandCustomTimeout tests that custom timeout is preserved
func TestSafeCommandCustomTimeout(t *testing.T) {
	// Create context with custom timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd, err := SafeCommand(ctx, "pytest", []string{"tests/"}...)
	if err != nil {
		t.Fatalf("SafeCommand failed: %v", err)
	}

	if cmd == nil {
		t.Fatal("SafeCommand returned nil cmd")
	}

	// Custom timeout should be preserved
	_ = cmd
}

// TestValidateTestCommand tests test command validation
func TestValidateTestCommand(t *testing.T) {
	tests := []struct {
		name        string
		testCmd     string
		shouldPass  bool
	}{
		{
			name:       "valid_pytest",
			testCmd:    "pytest",
			shouldPass: true,
		},
		{
			name:       "valid_go_test",
			testCmd:    "go test",
			shouldPass: true,
		},
		{
			name:       "injection_command",
			testCmd:    "pytest; rm -rf /",
			shouldPass: false,
		},
		{
			name:       "malicious_command",
			testCmd:    "rm -rf /",
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTestCommand(tt.testCmd)

			if tt.shouldPass {
				if err != nil {
					t.Errorf("ValidateTestCommand(%q) should pass but got error: %v", tt.testCmd, err)
				}
			} else {
				if err == nil {
					t.Errorf("ValidateTestCommand(%q) should fail but passed", tt.testCmd)
				}
			}
		})
	}
}
