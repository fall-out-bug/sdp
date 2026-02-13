package quality

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/config"
)

func (c *Checker) CheckCoverage() (*CoverageResult, error) {
	// Load threshold from guard rules (AC6)
	threshold := 80.0 // default
	projectRoot, rootErr := config.FindProjectRoot()
	if rootErr == nil {
		guardRules, rulesErr := config.LoadGuardRules(projectRoot)
		if rulesErr == nil {
			// Find coverage-threshold rule and get its threshold
			for _, rule := range guardRules.Rules {
				if rule.Enabled && rule.ID == "coverage-threshold" {
					if minVal, ok := rule.Config["minimum"]; ok {
						switch v := minVal.(type) {
						case int:
							threshold = float64(v)
						case float64:
							threshold = v
						}
					}
					break
				}
			}
		}
	}

	result := &CoverageResult{
		Threshold: threshold,
	}

	switch c.projectType {
	case Python:
		return c.checkPythonCoverage(result)
	case Go:
		return c.checkGoCoverage(result)
	case Java:
		return c.checkJavaCoverage(result)
	default:
		return result, fmt.Errorf("unsupported project type: %d", c.projectType)
	}
}
