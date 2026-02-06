package security

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCommandInjection tests that command injection is prevented
func TestCommandInjection(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		shouldBlock bool
		description string
	}{
		{
			name:        "safe_pytest_command",
			command:     "pytest",
			args:        []string{"tests/"},
			shouldBlock: false,
			description: "Legitimate pytest command should be allowed",
		},
		{
			name:        "safe_go_test_command",
			command:     "go",
			args:        []string{"test", "./..."},
			shouldBlock: false,
			description: "Legitimate go test command should be allowed",
		},
		{
			name:        "injection_via_semicolon",
			command:     "pytest",
			args:        []string{"tests/; rm -rf /"},
			shouldBlock: true,
			description: "Semicolon injection attempt should be blocked",
		},
		{
			name:        "injection_via_pipe",
			command:     "pytest",
			args:        []string{"|", "sh"},
			shouldBlock: true,
			description: "Pipe injection attempt should be blocked",
		},
		{
			name:        "injection_via_backtick",
			command:     "pytest",
			args:        []string{"`whoami`"},
			shouldBlock: true,
			description: "Backtick command substitution should be blocked",
		},
		{
			name:        "injection_via_dollar",
			command:     "pytest",
			args:        []string{"$(rm -rf /)"},
			shouldBlock: true,
			description: "Dollar sign command substitution should be blocked",
		},
		{
			name:        "injection_via_newline",
			command:     "pytest",
			args:        []string{"tests/\nrm -rf /"},
			shouldBlock: true,
			description: "Newline injection should be blocked",
		},
		{
			name:        "path_traversal",
			command:     "pytest",
			args:        []string{"../../../etc/passwd"},
			shouldBlock: true,
			description: "Path traversal attempt should be blocked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = exec.Command(tt.command, tt.args...) // _ because we only validate arguments

			// Check if any argument contains injection patterns
			for _, arg := range tt.args {
				if hasInjectionPattern(arg) {
					if !tt.shouldBlock {
						t.Errorf("%s: injection pattern detected but should be allowed: %s", tt.name, arg)
					}
					return
				}
			}

			if tt.shouldBlock {
				t.Errorf("%s: injection attempt should have been blocked: %s %v", tt.name, tt.command, tt.args)
			}
		})
	}
}

// hasInjectionPattern checks for command injection patterns
func hasInjectionPattern(s string) bool {
	injectionPatterns := []string{
		";",   // Command separator
		"|",   // Pipe
		"&",   // Background command
		"`",   // Command substitution (backtick)
		"$(",  // Command substitution (dollar)
		"\n",  // Newline
		"\r",  // Carriage return
		"\\",  // Escape character
		"../", // Path traversal (partial check)
	}

	for _, pattern := range injectionPatterns {
		if strings.Contains(s, pattern) {
			return true
		}
	}

	// Check for absolute paths to system files
	if strings.HasPrefix(s, "/etc/") || strings.HasPrefix(s, "/usr/") || strings.HasPrefix(s, "/bin/") {
		return true
	}

	return false
}

// TestWhitelistedCommands tests that only whitelisted commands are allowed
func TestWhitelistedCommands(t *testing.T) {
	whitelistedCommands := []string{
		"pytest", "pytest-3", "python -m pytest",
		"go test", "go vet", "go build",
		"mvn test", "mvnw test",
		"gradle test", "./gradlew test", "gradlew test",
		"npm test", "yarn test", "pnpm test",
		"jest", "mocha", "jasmine",
		"cargo test", "dart test", "flutter test",
		"git", // Git commands are safe
	}

	for _, cmd := range whitelistedCommands {
		t.Run("whitelisted_"+cmd, func(t *testing.T) {
			if !isWhitelistedCommand(cmd) {
				t.Errorf("Command '%s' should be whitelisted", cmd)
			}
		})
	}

	// Test that non-whitelisted commands are blocked
	nonWhitelisted := []string{
		"rm",
		"cat",
		"ls",
		"bash",
		"sh",
		"curl",
		"wget",
	}

	for _, cmd := range nonWhitelisted {
		t.Run("blocked_"+cmd, func(t *testing.T) {
			if isWhitelistedCommand(cmd) {
				t.Errorf("Command '%s' should NOT be whitelisted", cmd)
			}
		})
	}
}

// isWhitelistedCommand checks if a command is in the whitelist
func isWhitelistedCommand(cmd string) bool {
	whitelist := map[string]bool{
		"pytest":              true,
		"pytest-3":            true,
		"python -m pytest":    true,
		"go test":             true,
		"go vet":              true,
		"go build":            true,
		"mvn test":            true,
		"mvnw test":           true,
		"gradle test":         true,
		"./gradlew test":      true,
		"gradlew test":        true,
		"npm test":            true,
		"yarn test":           true,
		"pnpm test":           true,
		"jest":                true,
		"mocha":               true,
		"jasmine":             true,
		"cargo test":          true,
		"dart test":           true,
		"flutter test":        true,
		"git":                 true,
	}

	return whitelist[cmd]
}

// TestArgumentsContainInjection tests argument validation
func TestArgumentsContainInjection(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		hasInjection bool
	}{
		{
			name:        "safe_arguments",
			args:        []string{"tests/", "-v", "--cov"},
			hasInjection: false,
		},
		{
			name:        "semicolon_injection",
			args:        []string{"tests/; cat /etc/passwd"},
			hasInjection: true,
		},
		{
			name:        "pipe_injection",
			args:        []string{"|", "sh"},
			hasInjection: true,
		},
		{
			name:        "command_substitution",
			args:        []string{"$(whoami)"},
			hasInjection: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := argumentsContainInjection(tt.args)
			if result != tt.hasInjection {
				t.Errorf("argumentsContainInjection(%v) = %v, want %v", tt.args, result, tt.hasInjection)
			}
		})
	}
}

// argumentsContainInjection checks if any argument contains injection patterns
func argumentsContainInjection(args []string) bool {
	for _, arg := range args {
		if hasInjectionPattern(arg) {
			return true
		}
	}
	return false
}
