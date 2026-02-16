package main

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "world",
			input:    "World",
			expected: "Hello, World!",
		},
		{
			name:     "empty",
			input:    "",
			expected: "Hello, !",
		},
		{
			name:     "SDP",
			input:    "SDP",
			expected: "Hello, SDP!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := greet(tt.input)
			if result != tt.expected {
				t.Errorf("greet(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
