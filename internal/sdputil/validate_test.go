package sdputil

import (
	"testing"
)

func TestValidateFeatureID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid F014", "F014", false},
		{"valid F027", "F027", false},
		{"valid F1234", "F1234", false},
		{"empty", "", true},
		{"path separator", "F014/foo", true},
		{"backslash", "F014\\x", true},
		{"dot", "F014.", true},
		{"double dot", "F014..", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFeatureID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFeatureID(%q) err = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateWSID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid 00-014-01", "00-014-01", false},
		{"valid 00-027-01", "00-027-01", false},
		{"empty", "", true},
		{"path separator", "00-014/01", true},
		{"backslash", "00-014\\01", true},
		{"dot", "00-014.01", true},
		{"double dot", "..", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWSID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWSID(%q) err = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
