package ciloop

import "strings"

// Classification describes how a failing CI check should be handled.
type Classification string

const (
	ClassAutoFixable Classification = "auto-fixable"
	ClassEscalate    Classification = "escalate"
)

// FixType maps check name to fix handler: "go-test", "go-build", "k8s-validate", or "".
// Shared by Classify and Fixer.applyFix (DRY: yysx).
var fixTypePatterns = map[string][]string{
	"go-test":     {"go-test", "go test"},
	"go-build":    {"go-build", "go build"},
	"k8s-validate": {"k8s-validate", "k8s validate"},
}

// FixType returns the fix handler type for a check, or "" if not auto-fixable.
func FixType(checkName string) string {
	lower := strings.ToLower(checkName)
	for ft, patterns := range fixTypePatterns {
		for _, p := range patterns {
			if strings.Contains(lower, p) {
				return ft
			}
		}
	}
	return ""
}

// Classify returns the classification for a failing CI check by name.
// Auto-fixable checks are routed to deterministic fixers first (goimports, go mod tidy),
// then to the LLM/diagnostics path if fixers don't resolve. Unknown checks default to Escalate (fail-safe).
func Classify(checkName string) Classification {
	if FixType(checkName) != "" {
		return ClassAutoFixable
	}
	return ClassEscalate
}
