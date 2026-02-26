package quality

import (
	"context"
	"fmt"

	"github.com/fall-out-bug/sdp/internal/config"
)

func (c *Checker) CheckCoverage(ctx context.Context) (*CoverageResult, error) {
	// Load threshold from project config first, then guard rules
	threshold := 80.0 // default
	projectRoot, rootErr := config.FindProjectRoot()
	if rootErr == nil {
		if cfg, err := config.Load(projectRoot); err == nil && cfg.Quality.CoverageThreshold > 0 {
			threshold = float64(cfg.Quality.CoverageThreshold)
		} else {
			guardRules, rulesErr := config.LoadGuardRules(projectRoot)
			if rulesErr == nil {
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
	}

	result := &CoverageResult{
		Threshold: threshold,
	}

	switch c.projectType {
	case Python:
		return c.checkPythonCoverage(ctx, result)
	case Go:
		return c.checkGoCoverage(ctx, result)
	case Java:
		return c.checkJavaCoverage(ctx, result)
	default:
		return result, fmt.Errorf("unsupported project type: %d", c.projectType)
	}
}
