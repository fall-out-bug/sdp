package agents

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// AnalyzeGoBackend extracts routes from Go backend code
func (ca *CodeAnalyzer) AnalyzeGoBackend(filePath string) ([]ExtractedRoute, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Go file: %w", err)
	}

	if len(content) > MaxRegexMatchSize*10 {
		return nil, fmt.Errorf("file too large for analysis: %d bytes (max %d)", len(content), MaxRegexMatchSize*10)
	}

	var routes []ExtractedRoute
	lines := strings.Split(string(content), "\n")

	patterns := []struct {
		Name    string
		Pattern *regexp.Regexp
	}{
		{
			Name:    "gorilla/mux",
			Pattern: regexp.MustCompile(`HandleFunc\("([^"]{1,200})",\s*(\w+)\)\.Methods\("(\w{3,7})"\)`),
		},
		{
			Name:    "gin",
			Pattern: regexp.MustCompile(`(?:router|r)?\.(GET|POST|PUT|DELETE|PATCH)\("([^"]{1,200})",\s*(\w+)\)`),
		},
		{
			Name:    "echo",
			Pattern: regexp.MustCompile(`(?:echo|e)?\.(GET|POST|PUT|DELETE|PATCH)\("([^"]{1,200})",\s*(\w+)\)`),
		},
	}

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > MaxRegexMatchSize {
			continue
		}

		for _, p := range patterns {
			matches := p.Pattern.FindStringSubmatch(line)
			if len(matches) >= 3 {
				var path, method string
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

// AnalyzePythonSDK extracts public methods from Python SDK
func (ca *CodeAnalyzer) AnalyzePythonSDK(filePath string) ([]ExtractedMethod, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Python file: %w", err)
	}

	if len(content) > MaxRegexMatchSize*10 {
		return nil, fmt.Errorf("file too large for analysis: %d bytes (max %d)", len(content), MaxRegexMatchSize*10)
	}

	var methods []ExtractedMethod
	lines := strings.Split(string(content), "\n")

	methodRe := regexp.MustCompile(`def\s+(\w{1,100})\(self([^)]{0,500})\)(?:\s*->\s*[^:]{1,100})?:`)
	docsRe := regexp.MustCompile(`"""(.{1,500})"""`)

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > MaxRegexMatchSize {
			continue
		}

		matches := methodRe.FindStringSubmatch(line)
		if len(matches) >= 2 {
			methodName := matches[1]
			if strings.HasPrefix(methodName, "_") {
				continue
			}

			paramsStr := matches[2]
			var parameters []string
			if paramsStr != "" {
				paramsStr = strings.TrimPrefix(paramsStr, ",")
				params := strings.Split(paramsStr, ",")
				for _, p := range params {
					p = strings.TrimSpace(p)
					if p != "" && len(p) <= 100 {
						parts := strings.Fields(p)
						if len(parts) > 0 {
							parameters = append(parameters, parts[0])
						}
					}
				}
			}

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
				ReturnType:  "dict",
				Description: description,
				File:        filePath,
				Line:        lineNum + 1,
			})
		}
	}

	return methods, nil
}
