package collision

import (
	"encoding/json"
	"strings"
)

// FeatureScope represents a feature's scope for boundary detection.
type FeatureScope struct {
	FeatureID  string
	ScopeFiles []string
}

// SharedBoundary represents a shared type/interface boundary between features.
type SharedBoundary struct {
	FileName       string      `json:"fileName"`
	TypeName       string      `json:"typeName"`
	Fields         []FieldInfo `json:"fields"`
	Features       []string    `json:"features"`
	Recommendation string      `json:"recommendation"`
}

// FieldInfo represents a field in a struct or interface.
type FieldInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DetectBoundaries analyzes scope files to identify shared boundaries (types/interfaces).
func DetectBoundaries(features []FeatureScope) []SharedBoundary {
	fileToFeatures := buildFileToFeatures(features)
	var boundaries []SharedBoundary

	for file, featureIDs := range fileToFeatures {
		if len(featureIDs) < 2 {
			continue // Not a shared boundary
		}

		// Resolve relative path before parsing (bug fix for sdp-zidp)
		resolvedFile := resolveFilePath(file)

		// Parse Go file to extract types
		types, err := extractGoTypes(resolvedFile)
		if err != nil {
			continue // Skip files that can't be parsed
		}

		for _, typeName := range types {
			fields := extractStructFields(resolvedFile, typeName)
			boundaries = append(boundaries, SharedBoundary{
				FileName:       file, // Store original path
				TypeName:       typeName,
				Fields:         fields,
				Features:       featureIDs,
				Recommendation: "Create shared interface contract",
			})
		}
	}

	return boundaries
}

// buildFileToFeatures maps files to the features that use them.
// Deduplicates feature IDs to avoid same-feature false positives.
func buildFileToFeatures(features []FeatureScope) map[string][]string {
	fileToFeatures := make(map[string][]string)
	for _, f := range features {
		for _, file := range f.ScopeFiles {
			file = normalizePath(file)
			if file == "" {
				continue
			}
			// Check if it's a Go file
			if !strings.HasSuffix(file, ".go") {
				continue
			}
			// Deduplicate: only add if featureID not already present
			featureIDs := fileToFeatures[file]
			if !stringSliceContains(featureIDs, f.FeatureID) {
				fileToFeatures[file] = append(featureIDs, f.FeatureID)
			}
		}
	}
	return fileToFeatures
}

// stringSliceContains checks if a string slice contains a value.
func stringSliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// BoundaryToJSON converts a SharedBoundary to JSON string.
func BoundaryToJSON(boundary SharedBoundary) (string, error) {
	data, err := json.MarshalIndent(boundary, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
