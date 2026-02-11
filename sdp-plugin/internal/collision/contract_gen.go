package collision

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Contract represents a shared interface contract file.
type Contract struct {
	TypeName   string      `yaml:"typeName"`
	Fields     []FieldInfo `yaml:"fields"`
	RequiredBy []string    `yaml:"requiredBy"`
	Status     string      `yaml:"status"`
	FileName   string      `yaml:"fileName,omitempty"`
}

// GenerateContracts creates contract YAML files from shared boundaries.
func GenerateContracts(boundaries []SharedBoundary, outputDir string) ([]Contract, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	var contracts []Contract
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(boundaries))

	for _, boundary := range boundaries {
		wg.Add(1)
		go func(b SharedBoundary) {
			defer wg.Done()

			contract := Contract{
				TypeName:   b.TypeName,
				Fields:     b.Fields,
				RequiredBy: b.Features,
				Status:     "draft",
				FileName:   b.FileName,
			}

			contractPath := filepath.Join(outputDir, fmt.Sprintf("%s.yaml", b.TypeName))
			data, err := yaml.Marshal(contract)
			if err != nil {
				errChan <- fmt.Errorf("marshal contract %s: %w", b.TypeName, err)
				return
			}

			if err := os.WriteFile(contractPath, data, 0644); err != nil {
				errChan <- fmt.Errorf("write contract %s: %w", b.TypeName, err)
				return
			}

			mu.Lock()
			contracts = append(contracts, contract)
			mu.Unlock()
		}(boundary)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return contracts, err
		}
	}

	return contracts, nil
}
