package memory

import (
	"math"
	"strings"
)

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// containsMatch checks if text contains query (case-insensitive)
func containsMatch(text, query string) bool {
	if len(text) == 0 || len(query) == 0 {
		return false
	}
	return strings.Contains(strings.ToLower(text), strings.ToLower(query))
}
