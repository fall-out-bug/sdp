package metrics

import (
	"fmt"
	"strings"
)

// generateTaxonomySection creates failure breakdown (AC2).
func (r *Reporter) generateTaxonomySection(data ReportData) string {
	t := data.Taxonomy
	if t.TotalClassifications == 0 {
		return "\n## Failure Taxonomy\n\nNo failures recorded.\n"
	}

	var sb strings.Builder
	sb.WriteString("\n## Failure Taxonomy\n\n")
	sb.WriteString("Breakdown of verification failures by type:\n\n")

	// Sort by count (most common first)
	sortedTypes := r.sortByCount(t.ByType)

	for _, failureType := range sortedTypes {
		count := t.ByType[failureType]
		percentage := float64(count) / float64(t.TotalClassifications) * 100
		sb.WriteString(fmt.Sprintf("- **%s**: %d (%.1f%%) - %s\n",
			failureType, count, percentage, r.getFailureDescription(failureType)))
	}

	sb.WriteString(fmt.Sprintf("\n### Severity Distribution\n\n"))
	sb.WriteString(fmt.Sprintf("- CRITICAL: %d (systemic issues caught by acceptance)\n",
		t.BySeverity["CRITICAL"]))
	sb.WriteString(fmt.Sprintf("- HIGH: %d (crashes, edge cases)\n", t.BySeverity["HIGH"]))
	sb.WriteString(fmt.Sprintf("- MEDIUM: %d (logic, type errors)\n", t.BySeverity["MEDIUM"]))
	sb.WriteString(fmt.Sprintf("- LOW: %d (minor issues)\n", t.BySeverity["LOW"]))

	return sb.String()
}

// getFailureDescription returns human-readable description.
func (r *Reporter) getFailureDescription(failureType string) string {
	descriptions := map[string]string{
		"wrong_logic":            "Logic errors - incorrect algorithms or flow",
		"missing_edge_case":      "Edge cases - boundary conditions, nil pointers",
		"hallucinated_api":       "Hallucinated APIs - non-existent functions or methods",
		"type_error":             "Type errors - type mismatches or incompatibilities",
		"test_passing_but_wrong": "Test-Passing-But-Wrong - unit tests pass but acceptance fails",
		"compilation_error":      "Compilation errors - syntax or build failures",
		"import_error":           "Import errors - missing or unresolved modules",
		"unknown":                "Unknown - unclassified failures",
	}
	if desc, ok := descriptions[failureType]; ok {
		return desc
	}
	return "Uncategorized failure"
}

// sortByCount sorts failure types by count (descending).
func (r *Reporter) sortByCount(byType map[string]int) []string {
	type item struct {
		key   string
		count int
	}
	var items []item
	for k, v := range byType {
		items = append(items, item{key: k, count: v})
	}
	// Simple bubble sort (descending by count)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].count > items[i].count {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.key
	}
	return result
}

// generateTrendSection creates trend analysis (AC3).
func (r *Reporter) generateTrendSection(data ReportData) string {
	if len(data.Historical) == 0 {
		return "\n## Trends Over Time\n\nHistorical data not available. Future reports will include trend analysis as more data accumulates.\n"
	}

	var sb strings.Builder
	sb.WriteString("\n## Trends Over Time\n\n")
	sb.WriteString("Comparison of catch rates across periods:\n\n")
	sb.WriteString("| Period | Catch Rate | Total Workstreams |\n")
	sb.WriteString("|--------|------------|-------------------|\n")

	// Add current period
	sb.WriteString(fmt.Sprintf("| %s (current) | %.1f%% | %d |\n",
		r.GetCurrentQuarter(),
		data.Metrics.CatchRate*100,
		len(data.Metrics.IterationCount)))

	// Add historical periods
	for _, hist := range data.Historical {
		sb.WriteString(fmt.Sprintf("| %s | %.1f%% | %d |\n",
			hist.Period, hist.CatchRate*100, hist.TotalWS))
	}

	// Calculate trend
	if len(data.Historical) >= 1 {
		latest := data.Historical[len(data.Historical)-1].CatchRate
		current := data.Metrics.CatchRate
		if current < latest {
			sb.WriteString("\n**Trend:** Improving - Catch rate decreased from previous period\n")
		} else if current > latest {
			sb.WriteString("\n**Trend:** Declining - Catch rate increased from previous period\n")
		} else {
			sb.WriteString("\n**Trend:** Stable - Catch rate unchanged from previous period\n")
		}
	}

	return sb.String()
}
