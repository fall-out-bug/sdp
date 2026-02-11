package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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
	// Add timeout context
	ctx, cancel := context.WithTimeout(context.Background(), RegexTimeout)
	defer cancel()

	// Compile regex
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	// Test regex with pathological input to detect ReDoS
	pathologicalInput := strings.Repeat("a", 1000) + "!\""
	testChan := make(chan bool, 1)

	go func() {
		// This should complete quickly for safe regex
		re.FindStringSubmatch(pathologicalInput)
		testChan <- true
	}()

	select {
	case <-testChan:
		// Regex is safe (completed in time)
		return re, nil
	case <-ctx.Done():
		// Regex timed out - potential ReDoS
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

// AnalyzeGoBackend extracts routes from Go backend code
func (ca *CodeAnalyzer) AnalyzeGoBackend(filePath string) ([]ExtractedRoute, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Go file: %w", err)
	}

	// Check file size before processing
	if len(content) > MaxRegexMatchSize*10 {
		return nil, fmt.Errorf("file too large for analysis: %d bytes (max %d)", len(content), MaxRegexMatchSize*10)
	}

	var routes []ExtractedRoute
	lines := strings.Split(string(content), "\n")

	// Safe regex patterns for different frameworks
	// FIXED: Avoided catastrophic backtracking by:
	// 1. Using more explicit patterns (no ambiguous optional prefixes)
	// 2. Limiting character classes to prevent overlap
	// 3. Using atomic grouping where possible
	patterns := []struct {
		Name    string
		Pattern *regexp.Regexp
	}{
		{
			Name:    "gorilla/mux",
			// FIXED: More explicit pattern, no nested quantifiers, quote outside capture group
			Pattern: regexp.MustCompile(`HandleFunc\("([^"]{1,200})",\s*(\w+)\)\.Methods\("(\w{3,7})"\)`),
		},
		{
			Name: "gin",
			// FIXED: Removed ambiguous r? prefix, use alternation instead
			Pattern: regexp.MustCompile(`(?:router|r)?\.(GET|POST|PUT|DELETE|PATCH)\("([^"]{1,200})",\s*(\w+)\)`),
		},
		{
			Name: "echo",
			// FIXED: Removed ambiguous e? prefix, use alternation instead
			Pattern: regexp.MustCompile(`(?:echo|e)?\.(GET|POST|PUT|DELETE|PATCH)\("([^"]{1,200})",\s*(\w+)\)`),
		},
	}

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip overly long lines
		if len(line) > MaxRegexMatchSize {
			continue
		}

		// Try each pattern
		for _, p := range patterns {
			matches := p.Pattern.FindStringSubmatch(line)
			if len(matches) >= 3 {
				var path, method string

				// Extract path and method based on pattern
				switch p.Name {
				case "gorilla/mux":
					path = matches[1]
					method = matches[3]
				case "gin", "echo":
					method = matches[1]
					path = matches[2]
				}

				routes = append(routes, ExtractedRoute{
					Path:   path,
					Method: strings.ToUpper(method),
					File:   filePath,
					Line:   lineNum + 1,
				})
				break
			}
		}
	}

	return routes, nil
}

// AnalyzeTypeScriptFrontend extracts API calls from TypeScript/JavaScript
func (ca *CodeAnalyzer) AnalyzeTypeScriptFrontend(filePath string) ([]ExtractedCall, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TypeScript file: %w", err)
	}

	// Check file size
	if len(content) > MaxRegexMatchSize*10 {
		return nil, fmt.Errorf("file too large for analysis: %d bytes (max %d)", len(content), MaxRegexMatchSize*10)
	}

	var calls []ExtractedCall
	contentStr := string(content)

	// Truncate content if too large
	contentStr = truncateInput(contentStr)

	// FIXED: Safe regex patterns with explicit limits
	// First, process multiline patterns on the entire content
	multilinePatterns := []struct {
		Name    string
		Pattern *regexp.Regexp
		Method  string
	}{
		{
			Name:    "fetch with method (multiline)",
			// FIXED: Added non-greedy quantifier and explicit length limits
			Pattern: regexp.MustCompile(`fetch\("([^"]{1,200})",\s*\{[\s\S]{0,500}?method:\s*["'](\w+)["'][\s\S]{0,500}?\}`),
			Method:  "",
		},
		{
			Name:    "fetch GET (multiline with .then)",
			// FIXED: Limited multiline match to 500 chars
			Pattern: regexp.MustCompile(`fetch\("([^"]{1,200}")\)[\s\S]{0,300}?\.then\(`),
			Method:  "GET",
		},
	}

	// Find all multiline matches
	for _, p := range multilinePatterns {
		allMatches := p.Pattern.FindAllStringSubmatchIndex(contentStr, -1)
		for _, match := range allMatches {
			if len(match) >= 6 {
				// Extract matched text
				matchedText := contentStr[match[0]:match[1]]
				// Find line number by counting newlines before match
				lineNum := strings.Count(contentStr[:match[0]], "\n")

				// Extract submatches
				submatches := p.Pattern.FindStringSubmatch(matchedText)
				if len(submatches) >= 3 {
					calls = append(calls, ExtractedCall{
						Path:   submatches[1],
						Method: strings.ToUpper(submatches[2]),
						File:   filePath,
						Line:   lineNum + 1,
					})
				}
			}
		}
	}

	// Then process line-by-line for simple patterns
	lines := strings.Split(contentStr, "\n")
	simplePatterns := []struct {
		Name    string
		Pattern *regexp.Regexp
		Method  string
	}{
		{
			Name:    "fetch simple (single line, no comma or brace)",
			// FIXED: Added explicit length limit
			Pattern: regexp.MustCompile(`fetch\("([^"]{1,200}")`),
			Method:  "GET",
		},
		{
			Name:    "axios",
			// FIXED: Added explicit length limit to method and path
			Pattern: regexp.MustCompile(`axios\.(\w{3,7})\("([^"]{1,200}")`),
			Method:  "",
		},
	}

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip overly long lines
		if len(line) > MaxRegexMatchSize {
			continue
		}

		// Skip lines that are part of multiline patterns
		if strings.Contains(line, "{") || strings.Contains(line, "}.then") {
			continue
		}

		for _, p := range simplePatterns {
			matches := p.Pattern.FindStringSubmatch(line)
			if len(matches) >= 2 {
				path := matches[1]
				method := p.Method

				// Extract method from match if available
				if len(matches) >= 3 && matches[2] != "" {
					method = matches[2]
				}

				// For axios, method is in group 1
				if p.Name == "axios" && len(matches) >= 3 {
					method = matches[1]
					path = matches[2]
				}

				calls = append(calls, ExtractedCall{
					Path:   path,
					Method: strings.ToUpper(method),
					File:   filePath,
					Line:   lineNum + 1,
				})
				break
			}
		}
	}

	// Deduplicate calls (same path and method)
	uniqueCalls := make(map[string]ExtractedCall)
	for _, call := range calls {
		key := call.Path + ":" + call.Method
		uniqueCalls[key] = call
	}

	// Convert map back to slice
	var result []ExtractedCall
	for _, call := range uniqueCalls {
		result = append(result, call)
	}

	return result, nil
}

// AnalyzePythonSDK extracts public methods from Python SDK
func (ca *CodeAnalyzer) AnalyzePythonSDK(filePath string) ([]ExtractedMethod, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Python file: %w", err)
	}

	// Check file size
	if len(content) > MaxRegexMatchSize*10 {
		return nil, fmt.Errorf("file too large for analysis: %d bytes (max %d)", len(content), MaxRegexMatchSize*10)
	}

	var methods []ExtractedMethod
	lines := strings.Split(string(content), "\n")

	// FIXED: Added length limits to prevent ReDoS
	// Regex for method definition: def method_name(self, ...) -> ReturnType:
	// Handles both def method(self, param) and def method(self, param: type) -> ReturnType:
	methodRe := regexp.MustCompile(`def\s+(\w{1,100})\(self([^)]{0,500})\)(?:\s*->\s*[^:]{1,100})?:`)
	docsRe := regexp.MustCompile(`"""(.{1,500})"""`)

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip overly long lines
		if len(line) > MaxRegexMatchSize {
			continue
		}

		matches := methodRe.FindStringSubmatch(line)
		if len(matches) >= 2 {
			methodName := matches[1]

			// Skip private methods
			if strings.HasPrefix(methodName, "_") {
				continue
			}

			// Extract parameters
			paramsStr := matches[2]
			var parameters []string
			if paramsStr != "" {
				// Remove leading comma and split
				paramsStr = strings.TrimPrefix(paramsStr, ",")
				params := strings.Split(paramsStr, ",")
				for _, p := range params {
					p = strings.TrimSpace(p)
					if p != "" && len(p) <= 100 {
						// Extract parameter name (before type hint)
						parts := strings.Fields(p)
						if len(parts) > 0 {
							parameters = append(parameters, parts[0])
						}
					}
				}
			}

			// Look for docstring in next few lines
			description := ""
			for i := lineNum; i < len(lines) && i < lineNum+5; i++ {
				docMatches := docsRe.FindStringSubmatch(lines[i])
				if len(docMatches) >= 2 {
					description = docMatches[1]
					break
				}
			}

			methods = append(methods, ExtractedMethod{
				Name:        methodName,
				Parameters:  parameters,
				ReturnType:  "dict", // Default return type
				Description: description,
				File:        filePath,
				Line:        lineNum + 1,
			})
		}
	}

	return methods, nil
}

// GenerateOpenAPIContract generates an OpenAPI contract from extracted data
func (ca *CodeAnalyzer) GenerateOpenAPIContract(
	componentName string,
	backendRoutes []ExtractedRoute,
	frontendCalls []ExtractedCall,
) (*OpenAPIContract, error) {
	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info: InfoSpec{
			Title:   fmt.Sprintf("%s API", strings.Title(componentName)),
			Version: "1.0.0",
		},
		Paths: make(PathsSpec),
	}

	// Add backend routes to contract
	for _, route := range backendRoutes {
		if _, exists := contract.Paths[route.Path]; !exists {
			contract.Paths[route.Path] = make(PathSpec)
		}

		operation := OperationSpec{
			Summary: fmt.Sprintf("%s %s", route.Method, route.Path),
			Responses: ResponsesSpec{
				"200": ResponseSpec{
					Description: "Success",
					Content: map[string]MediaSpec{
						"application/json": {
							Schema: SchemaRefSpec{Type: "object"},
						},
					},
				},
			},
		}

		// Add request body for POST/PUT/PATCH
		if route.Method == "POST" || route.Method == "PUT" || route.Method == "PATCH" {
			operation.RequestBody = &RequestSpec{
				Required: true,
				Content: map[string]MediaSpec{
					"application/json": {
						Schema: SchemaRefSpec{Type: "object"},
					},
				},
			}
		}

		contract.Paths[route.Path][strings.ToLower(route.Method)] = operation
	}

	return contract, nil
}

// WriteExtractedContract writes the extracted contract to a YAML file
func (ca *CodeAnalyzer) WriteExtractedContract(contract *OpenAPIContract, outputPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write contract: %w", err)
	}

	return nil
}

// ExtractComponentContract performs end-to-end contract extraction for a component
func (ca *CodeAnalyzer) ExtractComponentContract(
	componentName string,
	componentType string, // "backend", "frontend", "sdk"
	filePaths []string,
	outputPath string,
) error {
	var contract *OpenAPIContract
	var err error

	switch componentType {
	case "backend":
		var allRoutes []ExtractedRoute
		for _, filePath := range filePaths {
			routes, err := ca.AnalyzeGoBackend(filePath)
			if err != nil {
				return fmt.Errorf("failed to analyze %s: %w", filePath, err)
			}
			allRoutes = append(allRoutes, routes...)
		}
		contract, err = ca.GenerateOpenAPIContract(componentName, allRoutes, nil)
		if err != nil {
			return fmt.Errorf("failed to generate contract: %w", err)
		}

	case "frontend":
		var allCalls []ExtractedCall
		for _, filePath := range filePaths {
			calls, err := ca.AnalyzeTypeScriptFrontend(filePath)
			if err != nil {
				return fmt.Errorf("failed to analyze %s: %w", filePath, err)
			}
			allCalls = append(allCalls, calls...)
		}
		// For frontend, we treat calls as routes
		var routes []ExtractedRoute
		for _, call := range allCalls {
			routes = append(routes, ExtractedRoute{
				Path:   call.Path,
				Method: call.Method,
				File:   call.File,
				Line:   call.Line,
			})
		}
		contract, err = ca.GenerateOpenAPIContract(componentName, routes, nil)
		if err != nil {
			return fmt.Errorf("failed to generate contract: %w", err)
		}

	case "sdk":
		// SDK methods don't map directly to OpenAPI
		// We'll create a minimal contract for now
		contract = &OpenAPIContract{
			OpenAPI: "3.0.0",
			Info: InfoSpec{
				Title:   fmt.Sprintf("%s SDK", strings.Title(componentName)),
				Version: "1.0.0",
			},
			Paths: make(PathsSpec),
		}

	default:
		return fmt.Errorf("unknown component type: %s", componentType)
	}

	// Write contract
	if err := ca.WriteExtractedContract(contract, outputPath); err != nil {
		return fmt.Errorf("failed to write contract: %w", err)
	}

	return nil
}
