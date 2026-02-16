package agents

import (
	"fmt"
	"os"
	"time"
)

// ValidateContractFile validates a contract file and returns issues
func (cv *ContractValidator) ValidateContractFile(contractPath string) ([]*ContractMismatch, error) {
	start := time.Now()
	success := false
	defer func() {
		duration := time.Since(start)
		cv.metrics.RecordValidation(success, duration, 0, 0, 0)
	}()

	content, err := os.ReadFile(contractPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read contract: %w", err)
	}

	contract := &OpenAPIContract{}
	parseErr := safeYAMLUnmarshal(content, contract)
	cv.metrics.RecordSchemaParse(parseErr == nil)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse contract: %w", parseErr)
	}

	var mismatches []*ContractMismatch

	if contract.OpenAPI == "" {
		mismatches = append(mismatches, &ContractMismatch{
			Severity: "ERROR",
			Type:     "invalid_contract",
			Expected: "openapi version",
			Actual:   "missing",
			Fix:      "Add openapi: 3.0.0 to contract",
		})
	}

	if len(contract.Paths) == 0 {
		mismatches = append(mismatches, &ContractMismatch{
			Severity: "WARNING",
			Type:     "invalid_contract",
			Expected: "at least one path",
			Actual:   "no paths defined",
			Fix:      "Add API paths to contract",
		})
	}

	success = true
	return mismatches, nil
}
