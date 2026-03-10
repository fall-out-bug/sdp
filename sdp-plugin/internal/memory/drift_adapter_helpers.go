package memory

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/drift"
)

// buildDriftContent creates searchable content from a DriftReport
func (a *DriftAdapter) buildDriftContent(report *drift.DriftReport) string {
	var content strings.Builder
	content.WriteString("Drift report for ")
	content.WriteString(report.WorkstreamID)
	content.WriteString("\nVerdict: ")
	content.WriteString(report.Verdict)
	content.WriteString("\nTimestamp: ")
	content.WriteString(report.Timestamp.Format(time.RFC3339))
	content.WriteString("\n")

	for _, issue := range report.Issues {
		content.WriteString("Issue: ")
		content.WriteString(issue.File)
		content.WriteString(" ")
		content.WriteString(string(issue.Status))
		content.WriteString("\nExpected: ")
		content.WriteString(issue.Expected)
		content.WriteString("\n")
		if issue.Actual != "" {
			content.WriteString("Actual: ")
			content.WriteString(issue.Actual)
			content.WriteString("\n")
		}
		if issue.Recommendation != "" {
			content.WriteString("Recommendation: ")
			content.WriteString(issue.Recommendation)
			content.WriteString("\n")
		}
	}

	return content.String()
}

// buildEnhancedDriftContent creates searchable content from an EnhancedDriftReport
func (a *DriftAdapter) buildEnhancedDriftContent(report *drift.EnhancedDriftReport, dataBytes []byte) string {
	var content strings.Builder
	content.WriteString("Enhanced drift report for ")
	content.WriteString(report.WorkstreamID)
	content.WriteString("\nVerdict: ")
	content.WriteString(report.Verdict)
	content.WriteString("\nTimestamp: ")
	content.WriteString(report.Timestamp.Format(time.RFC3339))
	content.WriteString("\n")

	for _, dt := range report.DriftTypes {
		content.WriteString("Type: ")
		content.WriteString(string(dt.Type))
		content.WriteString(" (")
		content.WriteString(string(dt.Severity))
		content.WriteString(")\n")
		for _, issue := range dt.Issues {
			content.WriteString("Issue: ")
			content.WriteString(issue.File)
			if issue.Line > 0 {
				content.WriteString(":")
				content.WriteString(strconv.Itoa(issue.Line))
			}
			content.WriteString(" ")
			content.WriteString(issue.Message)
			content.WriteString("\n")
		}
		for _, s := range dt.Suggestions {
			content.WriteString("Suggestion: ")
			content.WriteString(s)
			content.WriteString("\n")
		}
	}

	// Include full JSON for detailed search
	content.WriteString(" ")
	content.Write(dataBytes)

	return content.String()
}

// buildDriftTags builds tags from drift report
func (a *DriftAdapter) buildDriftTags(report *drift.EnhancedDriftReport) []string {
	tags := []string{"drift", report.Verdict}

	for _, dt := range report.DriftTypes {
		tags = append(tags, string(dt.Type))
		if dt.Severity == drift.SeverityError {
			tags = append(tags, "error")
		} else if dt.Severity == drift.SeverityWarning {
			tags = append(tags, "warning")
		}
	}

	return tags
}

// hashDriftReport generates a hash for a drift report
func hashDriftReport(report *drift.DriftReport) string {
	// Simple hash based on content
	data, err := json.Marshal(report)
	if err != nil {
		return "hash-error"
	}
	return simpleHash(string(data))
}

// hashEnhancedDriftReport generates a hash for an enhanced drift report
func hashEnhancedDriftReport(report *drift.EnhancedDriftReport) string {
	data, err := json.Marshal(report)
	if err != nil {
		return "hash-error"
	}
	return simpleHash(string(data))
}

// simpleHash generates a simple hash from a string
func simpleHash(s string) string {
	// Simple FNV-1a style hash
	h := uint32(2166136261)
	for _, c := range s {
		h ^= uint32(c)
		h *= 16777619
	}
	return "h" + intToStr(int(h))
}

// intToStr converts int to string without importing strconv
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
