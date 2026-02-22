package quality

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/config"
)

func (c *Checker) CheckComplexity() (*ComplexityResult, error) {
	threshold := 10
	if projectRoot, err := config.FindProjectRoot(); err == nil {
		if cfg, err := config.Load(projectRoot); err == nil && cfg.Quality.ComplexityThreshold > 0 {
			threshold = cfg.Quality.ComplexityThreshold
		}
	}
	result := &ComplexityResult{
		Threshold: threshold,
	}

	switch c.projectType {
	case Python:
		return c.checkPythonComplexity(result)
	case Go:
		return c.checkGoComplexity(result)
	case Java:
		return c.checkJavaComplexity(result)
	default:
		return result, fmt.Errorf("unsupported project type: %d", c.projectType)
	}
}
