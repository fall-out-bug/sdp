package memory

import (
	"errors"
	"math"
)

// QuantizeEmbedding converts float32 embedding to int8 (AC4)
// Achieves 4x size reduction (float32 = 4 bytes, int8 = 1 byte)
func QuantizeEmbedding(embedding []float32) ([]int8, float32) {
	if len(embedding) == 0 {
		return nil, 0
	}

	// Find max absolute value for scaling
	maxAbs := float32(0)
	for _, v := range embedding {
		abs := float32(math.Abs(float64(v)))
		if abs > maxAbs {
			maxAbs = abs
		}
	}

	if maxAbs == 0 {
		// All zeros
		return make([]int8, len(embedding)), 1
	}

	// Scale factor to map to int8 range [-127, 127]
	scale := float32(127) / maxAbs

	quantized := make([]int8, len(embedding))
	for i, v := range embedding {
		// Clamp to int8 range
		scaled := v * scale
		if scaled > 127 {
			scaled = 127
		} else if scaled < -127 {
			scaled = -127
		}
		quantized[i] = int8(scaled)
	}

	return quantized, scale
}

// QuantizeEmbeddingSafe returns error for empty input
func QuantizeEmbeddingSafe(embedding []float32) ([]int8, float32, error) {
	if len(embedding) == 0 {
		return nil, 0, errors.New("empty embedding")
	}
	q, s := QuantizeEmbedding(embedding)
	return q, s, nil
}

// DequantizeEmbedding converts int8 back to float32 (AC4)
func DequantizeEmbedding(quantized []int8, scale float32) []float32 {
	if scale == 0 {
		scale = 1
	}

	result := make([]float32, len(quantized))
	for i, v := range quantized {
		result[i] = float32(v) / scale
	}

	return result
}

// QuantizationStats tracks quantization metrics
type QuantizationStats struct {
	OriginalSize  int64   `json:"original_size"`
	QuantizedSize int64   `json:"quantized_size"`
	SavingsPct    float64 `json:"savings_pct"`
}

// CalculateQuantizationSavings returns storage savings from quantization
func CalculateQuantizationSavings(embeddingCount, dimensions int) QuantizationStats {
	originalSize := int64(embeddingCount * dimensions * 4)  // float32
	quantizedSize := int64(embeddingCount * dimensions * 1) // int8

	savings := float64(originalSize-quantizedSize) / float64(originalSize) * 100

	return QuantizationStats{
		OriginalSize:  originalSize,
		QuantizedSize: quantizedSize,
		SavingsPct:    savings,
	}
}
