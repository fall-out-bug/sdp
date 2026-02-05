package drift

import (
	"fmt"
	"strings"
	"time"
)

// Status represents the status of a drift check
type Status string

const (
	// StatusOK means no issues found
	StatusOK Status = "OK"
	// StatusWarning means potential issues found
	StatusWarning Status = "WARNING"
	// StatusError means critical issues found
	StatusError Status = "ERROR"
)

// DriftIssue represents a single drift issue found
type DriftIssue struct {
	File          string
	Status        Status
	Expected      string
	Actual        string
	Recommendation string
}

// DriftReport represents the complete drift detection report
type DriftReport struct {
	WorkstreamID string
	Timestamp    time.Time
	Issues       []DriftIssue
	Verdict      string // PASS, FAIL, WARNING
}

// String returns a markdown-formatted report
func (r *DriftReport) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## Drift Report: %s\n\n", r.WorkstreamID))
	sb.WriteString("| File | Status | Issue |\n")
	sb.WriteString("|------|--------|-------|\n")

	for _, issue := range r.Issues {
		icon := "✅"
		if issue.Status == StatusWarning {
			icon = "⚠️"
		} else if issue.Status == StatusError {
			icon = "❌"
		}

		issueDesc := issue.Expected
		if issue.Actual != "" {
			issueDesc = fmt.Sprintf("%s (found: %s)", issue.Expected, issue.Actual)
		}

		sb.WriteString(fmt.Sprintf("| %s | %s %s | %s |\n",
			issue.File, icon, issue.Status, issueDesc))
	}

	// Count issues
	errorCount := 0
	warningCount := 0
	for _, issue := range r.Issues {
		if issue.Status == StatusError {
			errorCount++
		} else if issue.Status == StatusWarning {
			warningCount++
		}
	}

	// Verdict
	verdictIcon := "✅"
	if r.Verdict == "WARNING" {
		verdictIcon = "⚠️"
	} else if r.Verdict == "FAIL" {
		verdictIcon = "❌"
	}

	sb.WriteString(fmt.Sprintf("\n**Verdict:** %s %s", verdictIcon, r.Verdict))
	if errorCount > 0 || warningCount > 0 {
		sb.WriteString(fmt.Sprintf(" - %d error(s), %d warning(s)", errorCount, warningCount))
	}
	sb.WriteString("\n")

	// Recommendations
	if len(r.Issues) > 0 {
		sb.WriteString("\n**Recommendations:**\n")
		for i, issue := range r.Issues {
			if issue.Recommendation != "" {
				sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, issue.Recommendation))
			}
		}
	}

	return sb.String()
}

// GenerateVerdict determines the overall verdict based on issues
func (r *DriftReport) GenerateVerdict() string {
	errorCount := 0
	warningCount := 0

	for _, issue := range r.Issues {
		if issue.Status == StatusError {
			errorCount++
		} else if issue.Status == StatusWarning {
			warningCount++
		}
	}

	if errorCount > 0 {
		return "FAIL"
	} else if warningCount > 0 {
		return "WARNING"
	}
	return "PASS"
}

// AddIssue adds an issue to the report
func (r *DriftReport) AddIssue(issue DriftIssue) {
	r.Issues = append(r.Issues, issue)
	r.Verdict = r.GenerateVerdict()
}
