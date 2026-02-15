package agents

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// GenerateReport generates a markdown validation report
func (cv *ContractValidator) GenerateReport(mismatches []*ContractMismatch) string {
	return cv.GenerateReportWithOptions(mismatches, false)
}

// GenerateReportWithOptions generates a markdown validation report with options
func (cv *ContractValidator) GenerateReportWithOptions(mismatches []*ContractMismatch, redact bool) string {
	var sb strings.Builder

	sb.WriteString("# Contract Validation Report\n\n")

	sort.Slice(mismatches, func(i, j int) bool {
		if mismatches[i].Severity != mismatches[j].Severity {
			return mismatches[i].Severity > mismatches[j].Severity
		}
		return mismatches[i].Type < mismatches[j].Type
	})

	errorCount := 0
	warningCount := 0
	infoCount := 0
	sensitiveCount := 0

	for _, m := range mismatches {
		switch m.Severity {
		case "ERROR":
			errorCount++
		case "WARNING":
			warningCount++
		case "INFO":
			infoCount++
		}

		if isSensitivePath(m.Path) {
			sensitiveCount++
		}
	}

	sb.WriteString("## Summary\n\n")
	sb.WriteString(fmt.Sprintf("- Total issues: %d\n", len(mismatches)))
	sb.WriteString(fmt.Sprintf("- Errors: %d\n", errorCount))
	sb.WriteString(fmt.Sprintf("- Warnings: %d\n", warningCount))
	sb.WriteString(fmt.Sprintf("- Info: %d\n", infoCount))

	if redact && sensitiveCount > 0 {
		sb.WriteString(fmt.Sprintf("- ⚠️ Sensitive endpoints redacted: %d\n\n", sensitiveCount))
	} else if sensitiveCount > 0 {
		sb.WriteString(fmt.Sprintf("- Sensitive endpoints detected: %d\n\n", sensitiveCount))
	} else {
		sb.WriteString("\n")
	}

	if errorCount > 0 {
		sb.WriteString("## Errors\n\n")
		cv.writeMismatchesTable(&sb, mismatches, "ERROR", redact)
	}

	if warningCount > 0 {
		sb.WriteString("## Warnings\n\n")
		cv.writeMismatchesTable(&sb, mismatches, "WARNING", redact)
	}

	if infoCount > 0 {
		sb.WriteString("## Info\n\n")
		cv.writeMismatchesTable(&sb, mismatches, "INFO", redact)
	}

	if len(mismatches) == 0 {
		sb.WriteString("✅ No contract mismatches found!\n")
	}

	return sb.String()
}

// writeMismatchesTable writes a markdown table for mismatches of given severity
func (cv *ContractValidator) writeMismatchesTable(sb *strings.Builder, mismatches []*ContractMismatch, severity string, redact bool) {
	sb.WriteString("| Component | Type | Expected | Actual | Fix |\n")
	sb.WriteString("|-----------|------|----------|--------|-----|\n")

	for _, m := range mismatches {
		if m.Severity != severity {
			continue
		}

		component := fmt.Sprintf("%s vs %s", m.ComponentA, m.ComponentB)
		expected := m.Expected
		actual := m.Actual
		fix := m.Fix

		if redact {
			if isSensitivePath(m.Path) {
				expected = "[REDACTED]"
				actual = "[REDACTED]"
				fix = "Review manually (sensitive endpoint)"
			}
			component = redactSensitiveInfo(component)
		}

		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			component, m.Type, expected, actual, fix))
	}

	sb.WriteString("\n")
}

// WriteReport writes the validation report to a file
func (cv *ContractValidator) WriteReport(report, outputPath string) error {
	start := time.Now()
	success := false
	defer func() {
		duration := time.Since(start)
		cv.metrics.RecordReportGeneration(success, duration)
	}()

	dir := outputPath[:strings.LastIndex(outputPath, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	success = true
	return nil
}

// isSensitivePath checks if a path contains sensitive information
func isSensitivePath(path string) bool {
	sensitivePrefixes := []string{
		"/admin", "/internal", "/private", "/config",
		"/secret", "/auth", "/login", "/logout",
		"/password", "/token", "/key", "/credentials",
	}

	lowerPath := strings.ToLower(path)
	for _, prefix := range sensitivePrefixes {
		if strings.HasPrefix(lowerPath, prefix) || strings.Contains(lowerPath, "/"+prefix+"/") {
			return true
		}
	}

	return false
}

// redactSensitiveInfo redacts sensitive information from a string
func redactSensitiveInfo(input string) string {
	if strings.Contains(input, "/") {
		parts := strings.Split(input, "/")
		if len(parts) > 1 {
			filename := parts[len(parts)-1]
			input = "***/***/" + filename
		}
	}

	return input
}
