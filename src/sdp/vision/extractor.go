package vision

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// FeatureDraft represents a feature extracted from PRD
type FeatureDraft struct {
	Title       string
	Description string
	Priority    string
}

// Slug returns URL-friendly slug from title
func (f *FeatureDraft) Slug() string {
	// Convert to lowercase and replace spaces with hyphens
	slug := strings.ToLower(f.Title)
	// Remove non-alphanumeric characters (except spaces and hyphens)
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

// ExtractFeaturesFromPRD extracts P0 and P1 features from PRD file
func ExtractFeaturesFromPRD(prdPath string) ([]FeatureDraft, error) {
	// Read PRD file
	file, err := os.Open(prdPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PRD: %w", err)
	}
	defer file.Close()

	var features []FeatureDraft
	scanner := bufio.NewScanner(file)
	currentPriority := ""
	inFeaturesSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect Features section
		if strings.HasPrefix(line, "## Features") {
			inFeaturesSection = true
			continue
		}

		// Exit Features section on next major section
		if inFeaturesSection && strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "## Features") {
			break
		}

		// Detect priority level (P0, P1, P2)
		if match := regexp.MustCompile(`### (P[012])\s+\(`).FindStringSubmatch(line); len(match) > 0 {
			currentPriority = match[1]
			// Reset if we see P2 (we only want P0 and P1)
			if match[1] == "P2" {
				currentPriority = ""
			}
			continue
		}

		// Extract feature (only P0 and P1)
		if currentPriority == "P0" || currentPriority == "P1" {
			// Try multiple patterns for feature lines
			// Pattern 1: "- Feature N: Title"
			if match := regexp.MustCompile(`-\s+Feature\s+\d+:\s+(.+)`).FindStringSubmatch(line); len(match) > 0 {
				title := strings.TrimSpace(match[1])
				features = append(features, FeatureDraft{
					Title:       title,
					Description: fmt.Sprintf("Feature: %s", title),
					Priority:    currentPriority,
				})
			} else if match := regexp.MustCompile(`-\s+\d+:\s+(.+)`).FindStringSubmatch(line); len(match) > 0 {
				// Pattern 2: "- N: Title" (without "Feature" keyword)
				title := strings.TrimSpace(match[1])
				features = append(features, FeatureDraft{
					Title:       title,
					Description: fmt.Sprintf("Feature: %s", title),
					Priority:    currentPriority,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read PRD: %w", err)
	}

	return features, nil
}
