package agents

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/fall-out-bug/sdp/src/sdp/monitoring"
	"gopkg.in/yaml.v3"
)

const (
	// MaxYAMLFileSize is the maximum allowed YAML file size (10MB)
	MaxYAMLFileSize = 10 * 1024 * 1024
	// YAMLParseTimeout is the maximum time allowed for YAML parsing
	YAMLParseTimeout = 30 * time.Second
)

// safeYAMLUnmarshal safely unmarshals YAML with security controls
func safeYAMLUnmarshal(data []byte, v interface{}) error {
	if len(data) > MaxYAMLFileSize {
		return fmt.Errorf("YAML file size %d bytes exceeds maximum allowed size %d bytes", len(data), MaxYAMLFileSize)
	}

	ctx, cancel := context.WithTimeout(context.Background(), YAMLParseTimeout)
	defer cancel()

	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	done := make(chan error, 1)
	go func() {
		done <- decoder.Decode(v)
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("YAML parse error: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("YAML parsing timeout after %v", YAMLParseTimeout)
	}
}

// ContractMismatch represents a detected contract mismatch
type ContractMismatch struct {
	Severity   string `yaml:"severity"`
	Type       string `yaml:"type"`
	ComponentA string `yaml:"component_a"`
	ComponentB string `yaml:"component_b"`
	Path       string `yaml:"path"`
	Method     string `yaml:"method"`
	Expected   string `yaml:"expected"`
	Actual     string `yaml:"actual"`
	File       string `yaml:"file"`
	Fix        string `yaml:"fix"`
}

// ContractValidator validates contracts against each other
type ContractValidator struct {
	metrics *monitoring.MetricsCollector
}

// NewContractValidator creates a new contract validator
func NewContractValidator() *ContractValidator {
	return &ContractValidator{
		metrics: monitoring.NewMetricsCollector(),
	}
}

// NewContractValidatorWithMetrics creates a new contract validator with custom metrics collector
func NewContractValidatorWithMetrics(metrics *monitoring.MetricsCollector) *ContractValidator {
	return &ContractValidator{
		metrics: metrics,
	}
}

// GetMetrics returns the current metrics snapshot
func (cv *ContractValidator) GetMetrics() *monitoring.MetricsSnapshot {
	return cv.metrics.GetMetrics()
}
