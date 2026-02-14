package agents

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// MaxContentLength limits content size to prevent ReDoS
	MaxContentLength = 100000
	// MaxFieldCount limits number of fields to prevent resource exhaustion
	MaxFieldCount = 100
)

var (
	// allowedHTTPMethods is a whitelist of permitted HTTP methods
	allowedHTTPMethods = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
	}
)

// validateFeatureName checks if a feature name is valid
func validateFeatureName(name string) error {
	// Feature names must be lowercase alphanumeric with dashes only
	// This prevents injection attacks via malicious feature names
	matched, err := regexp.MatchString(`^[a-z0-9-]+$`, name)
	if err != nil {
		return fmt.Errorf("feature name validation failed: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid feature name %q: must contain only lowercase letters, numbers, and dashes", name)
	}
	return nil
}

// validateHTTPMethod checks if an HTTP method is allowed
func validateHTTPMethod(method string) error {
	// Method must contain only uppercase letters (no special characters, no spaces)
	if !regexp.MustCompile(`^[A-Z]+$`).MatchString(method) {
		return fmt.Errorf("invalid HTTP method %q: must contain only uppercase letters", method)
	}

	// Normalize to uppercase for whitelist check
	normalizedMethod := strings.ToUpper(method)
	if !allowedHTTPMethods[normalizedMethod] {
		return fmt.Errorf("invalid HTTP method %q: must be one of %v", method, getAllowedMethods())
	}
	return nil
}

// getAllowedMethods returns the list of allowed HTTP methods
func getAllowedMethods() []string {
	methods := make([]string, 0, len(allowedHTTPMethods))
	for method := range allowedHTTPMethods {
		methods = append(methods, method)
	}
	return methods
}

// validateEndpointPath checks if an endpoint path is valid
func validateEndpointPath(path string) error {
	// Paths must start with /
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("invalid endpoint path %q: must start with /", path)
	}

	// Paths should not contain .. to prevent directory traversal
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid endpoint path %q: must not contain ..", path)
	}

	// Limit path length to prevent DoS
	if len(path) > 500 {
		return fmt.Errorf("invalid endpoint path %q: too long (max 500 chars)", path)
	}

	// Paths should only contain safe characters
	// Allow: /, alphanumeric, -, _, :, {}, * (for path params and wildcards)
	matched, err := regexp.MatchString(`^/[a-zA-Z0-9\-_{}/*.:]+$`, path)
	if err != nil {
		return fmt.Errorf("endpoint path validation failed: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid endpoint path %q: contains invalid characters", path)
	}

	return nil
}

// sanitizeFieldName sanitizes field names to prevent injection
func sanitizeFieldName(name string) (string, error) {
	// Field names must be alphanumeric with underscores
	matched, err := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	if err != nil {
		return "", fmt.Errorf("field name validation failed: %w", err)
	}
	if !matched {
		return "", fmt.Errorf("invalid field name %q: must be alphanumeric with underscores, starting with letter or underscore", name)
	}
	return name, nil
}
