package memory

import (
	"testing"
)

func TestQuantizeEmbedding_Basic(t *testing.T) {
	embedding := []float32{0.5, -0.3, 0.8, -0.1, 0.0}

	quantized, scale := QuantizeEmbedding(embedding)

	if len(quantized) != len(embedding) {
		t.Errorf("Length mismatch: %d vs %d", len(quantized), len(embedding))
	}

	if scale <= 0 {
		t.Error("Scale should be positive")
	}
}

func TestQuantizeEmbedding_Empty(t *testing.T) {
	embedding := []float32{}

	_, _, err := QuantizeEmbeddingSafe(embedding)
	if err == nil {
		t.Error("Expected error for empty embedding")
	}
}

func TestDequantizeEmbedding_RoundTrip(t *testing.T) {
	original := []float32{0.5, -0.3, 0.8, -0.1, 0.0, 0.99, -0.99}

	quantized, scale := QuantizeEmbedding(original)
	reconstructed := DequantizeEmbedding(quantized, scale)

	// Should be approximately equal (with quantization error)
	for i := range original {
		diff := original[i] - reconstructed[i]
		if diff < 0 {
			diff = -diff
		}
		// Quantization error should be < 1%
		if diff > 0.02 {
			t.Errorf("Round trip error at %d: %f vs %f (diff: %f)", i, original[i], reconstructed[i], diff)
		}
	}
}

func TestQuantizeEmbedding_AllPositive(t *testing.T) {
	embedding := []float32{0.1, 0.5, 0.9, 0.3}

	quantized, _ := QuantizeEmbedding(embedding)

	for i, q := range quantized {
		if q < 0 {
			t.Errorf("Quantized value at %d should be positive: %d", i, q)
		}
	}
}

func TestQuantizeEmbedding_AllNegative(t *testing.T) {
	embedding := []float32{-0.1, -0.5, -0.9, -0.3}

	quantized, _ := QuantizeEmbedding(embedding)

	for i, q := range quantized {
		if q > 0 {
			t.Errorf("Quantized value at %d should be negative: %d", i, q)
		}
	}
}

func TestQuantizationSavings(t *testing.T) {
	// 1536-dimensional embedding
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = 0.5
	}

	// Before: 1536 * 4 bytes = 6144 bytes
	beforeSize := len(embedding) * 4

	quantized, _ := QuantizeEmbedding(embedding)

	// After: 1536 * 1 byte = 1536 bytes
	afterSize := len(quantized)

	savings := float64(beforeSize-afterSize) / float64(beforeSize) * 100

	// Should be ~75% savings
	if savings < 70 || savings > 80 {
		t.Errorf("Expected ~75%% savings, got %.1f%%", savings)
	}
}

func TestCalculateQuantizationSavings(t *testing.T) {
	stats := CalculateQuantizationSavings(1000, 1536)

	if stats.OriginalSize != 6144000 {
		t.Errorf("Expected original size 6144000, got %d", stats.OriginalSize)
	}
	if stats.QuantizedSize != 1536000 {
		t.Errorf("Expected quantized size 1536000, got %d", stats.QuantizedSize)
	}
	if stats.SavingsPct < 74 || stats.SavingsPct > 76 {
		t.Errorf("Expected ~75%% savings, got %.1f%%", stats.SavingsPct)
	}
}
