package quality

import (
	"fmt"
)

func (c *Checker) CheckComplexity() (*ComplexityResult, error) {
	result := &ComplexityResult{
		Threshold: 10,
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
