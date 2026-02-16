package agents

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	// MaxRegexMatchSize limits the size of regex matches to prevent ReDoS
	MaxRegexMatchSize = 10000
	// RegexTimeout is the maximum time to spend on regex operations
	RegexTimeout = 5 * time.Second
)

// CodeAnalyzer extracts API contracts from existing code
type CodeAnalyzer struct{}

// ExtractedRoute represents a backend route
type ExtractedRoute struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	File   string `yaml:"file"`
	Line   int    `yaml:"line"`
}

// ExtractedCall represents a frontend API call
type ExtractedCall struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	File   string `yaml:"file"`
	Line   int    `yaml:"line"`
}

// ExtractedMethod represents a Python SDK method
type ExtractedMethod struct {
	Name        string   `yaml:"name"`
	Parameters  []string `yaml:"parameters"`
	ReturnType  string   `yaml:"return_type"`
	Description string   `yaml:"description"`
	File        string   `yaml:"file"`
	Line        int      `yaml:"line"`
}

// NewCodeAnalyzer creates a new code analyzer
func NewCodeAnalyzer() *CodeAnalyzer {
	return &CodeAnalyzer{}
}

// safeCompileRegex compiles a regex with timeout and size limits
func safeCompileRegex(pattern string) (*regexp.Regexp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RegexTimeout)
	defer cancel()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Test regex with pathological input to detect ReDoS
	pathologicalInput := strings.Repeat("a", 1000) + "!\""
	testChan := make(chan bool, 1)

	go func() {
		re.FindStringSubmatch(pathologicalInput)
		testChan <- true
	}()

	select {
	case <-testChan:
		return re, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("regex pattern potentially vulnerable to ReDoS: timeout after %v", RegexTimeout)
	}
}

// truncateInput limits input size to prevent ReDoS
func truncateInput(input string) string {
	if len(input) > MaxRegexMatchSize {
		return input[:MaxRegexMatchSize]
	}
	return input
}
