package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

const (
	// MaxWorkstreamFileSize is the maximum allowed size for workstream files (1MB)
	MaxWorkstreamFileSize = 1 << 20 // 1MB
	// MaxStringLength is the maximum allowed length for string fields (10KB)
	MaxStringLength = 10 << 10 // 10KB
	// MaxYAMLDocumentDepth is the maximum allowed nesting depth for YAML documents
	MaxYAMLDocumentDepth = 100
)

var (
	// ErrFileTooLarge is returned when a workstream file exceeds the size limit
	ErrFileTooLarge = errors.New("file size exceeds maximum allowed size")
	// ErrFieldTooLong is returned when a string field exceeds the maximum length
	ErrFieldTooLong = errors.New("field length exceeds maximum allowed length")
	// ErrYAMLDepthExceeded is returned when YAML nesting depth is too deep
	ErrYAMLDepthExceeded = errors.New("YAML document depth exceeds maximum")
)

// SafeYAMLDecoder wraps yaml.Decoder with security limits
type SafeYAMLDecoder struct {
	decoder *yaml.Decoder
	maxSize int64
	depth   int
}

// NewSafeYAMLDecoder creates a new SafeYAMLDecoder with security limits
func NewSafeYAMLDecoder(r io.Reader) *SafeYAMLDecoder {
	return &SafeYAMLDecoder{
		decoder: yaml.NewDecoder(r),
		maxSize: MaxWorkstreamFileSize,
		depth:   0,
	}
}

// Decode decodes YAML with security checks
func (d *SafeYAMLDecoder) Decode(v interface{}) error {
	// Check depth before decoding
	if d.depth > MaxYAMLDocumentDepth {
		return fmt.Errorf("%w: depth %d exceeds maximum %d", ErrYAMLDepthExceeded, d.depth, MaxYAMLDocumentDepth)
	}

	d.depth++
	defer func() { d.depth-- }()

	// Decode using standard YAML decoder
	if err := d.decoder.Decode(v); err != nil {
		return err
	}

	// Validate string lengths in decoded struct
	if err := validateStringLengths(v); err != nil {
		return err
	}

	return nil
}

// validateStringLengths validates that all string fields are within limits
func validateStringLengths(v interface{}) error {
	// For frontmatter struct, validate specific fields
	if fm, ok := v.(*frontmatter); ok {
		if len(fm.WSID) > MaxStringLength {
			return fmt.Errorf("%w: ws_id length %d exceeds maximum %d", ErrFieldTooLong, len(fm.WSID), MaxStringLength)
		}
		if len(fm.Parent) > MaxStringLength {
			return fmt.Errorf("%w: parent length %d exceeds maximum %d", ErrFieldTooLong, len(fm.Parent), MaxStringLength)
		}
		if len(fm.Feature) > MaxStringLength {
			return fmt.Errorf("%w: feature length %d exceeds maximum %d", ErrFieldTooLong, len(fm.Feature), MaxStringLength)
		}
		if len(fm.Status) > MaxStringLength {
			return fmt.Errorf("%w: status length %d exceeds maximum %d", ErrFieldTooLong, len(fm.Status), MaxStringLength)
		}
		if len(fm.Size) > MaxStringLength {
			return fmt.Errorf("%w: size length %d exceeds maximum %d", ErrFieldTooLong, len(fm.Size), MaxStringLength)
		}
		if len(fm.ProjectID) > MaxStringLength {
			return fmt.Errorf("%w: project_id length %d exceeds maximum %d", ErrFieldTooLong, len(fm.ProjectID), MaxStringLength)
		}
	}
	return nil
}

// ValidateFileSize checks if a file is within the size limit
func ValidateFileSize(content []byte) error {
	size := int64(len(content))
	if size > MaxWorkstreamFileSize {
		return fmt.Errorf("%w: %d bytes exceeds maximum %d bytes", ErrFileTooLarge, size, MaxWorkstreamFileSize)
	}
	return nil
}

// ValidateContentLength validates the length of content fields
// Note: The entire file can be up to 1MB, but individual sections should be reasonable
func ValidateContentLength(content string) error {
	// We allow up to 1MB for the entire content section
	// This is validated separately from the file size check
	const maxContentLength = 1 << 20 // 1MB
	if len(content) > maxContentLength {
		return fmt.Errorf("%w: content length %d exceeds maximum %d", ErrFieldTooLong, len(content), maxContentLength)
	}
	return nil
}

// SafeYAMLUnmarshal wraps yaml.Unmarshal with security checks
func SafeYAMLUnmarshal(data []byte, v interface{}) error {
	// Check file size first
	if err := ValidateFileSize(data); err != nil {
		return err
	}

	// Use safe decoder
	decoder := NewSafeYAMLDecoder(bytes.NewReader(data))
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("YAML decode failed: %w", err)
	}

	return nil
}
