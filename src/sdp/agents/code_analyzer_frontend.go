package agents

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// AnalyzeTypeScriptFrontend extracts API calls from TypeScript/JavaScript
func (ca *CodeAnalyzer) AnalyzeTypeScriptFrontend(filePath string) ([]ExtractedCall, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TypeScript file: %w", err)
	}

	if len(content) > MaxRegexMatchSize*10 {
		return nil, fmt.Errorf("file too large for analysis: %d bytes (max %d)", len(content), MaxRegexMatchSize*10)
	}

	var calls []ExtractedCall
	contentStr := string(content)
	contentStr = truncateInput(contentStr)

	// Multiline patterns
	multilinePatterns := []struct {
		Name    string
		Pattern *regexp.Regexp
		Method  string
	}{
		{
			Name:    "fetch with method (multiline)",
			Pattern: regexp.MustCompile(`fetch\("([^"]{1,200})",\s*\{[\s\S]{0,500}?method:\s*["'](\w+)["'][\s\S]{0,500}?\}`),
			Method:  "",
		},
		{
			Name:    "fetch GET (multiline with .then)",
			Pattern: regexp.MustCompile(`fetch\("([^"]{1,200}")\)[\s\S]{0,300}?\.then\(`),
			Method:  "GET",
		},
	}

	for _, p := range multilinePatterns {
		allMatches := p.Pattern.FindAllStringSubmatchIndex(contentStr, -1)
		for _, match := range allMatches {
			if len(match) >= 6 {
				matchedText := contentStr[match[0]:match[1]]
				lineNum := strings.Count(contentStr[:match[0]], "\n")

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

	// Simple line-by-line patterns
	lines := strings.Split(contentStr, "\n")
	simplePatterns := []struct {
		Name    string
		Pattern *regexp.Regexp
		Method  string
	}{
		{
			Name:    "fetch simple",
			Pattern: regexp.MustCompile(`fetch\("([^"]{1,200}")`),
			Method:  "GET",
		},
		{
			Name:    "axios",
			Pattern: regexp.MustCompile(`axios\.(\w{3,7})\("([^"]{1,200}")`),
			Method:  "",
		},
	}

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > MaxRegexMatchSize {
			continue
		}

		if strings.Contains(line, "{") || strings.Contains(line, "}.then") {
			continue
		}

		for _, p := range simplePatterns {
			matches := p.Pattern.FindStringSubmatch(line)
			if len(matches) >= 2 {
				path := matches[1]
				method := p.Method

				if len(matches) >= 3 && matches[2] != "" {
					method = matches[2]
				}

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

	// Deduplicate
	uniqueCalls := make(map[string]ExtractedCall)
	for _, call := range calls {
		key := call.Path + ":" + call.Method
		uniqueCalls[key] = call
	}

	var result []ExtractedCall
	for _, call := range uniqueCalls {
		result = append(result, call)
	}

	return result, nil
}
