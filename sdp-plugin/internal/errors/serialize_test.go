package errors

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestToSerializable(t *testing.T) {
	t.Run("sdp_error", func(t *testing.T) {
		err := New(ErrGitNotFound, nil).WithContext("file", "config.yml")
		se := ToSerializable(err)

		if se.Code != "ENV001" {
			t.Errorf("Code = %q, want %q", se.Code, "ENV001")
		}
		if se.Class != "ENV" {
			t.Errorf("Class = %q, want %q", se.Class, "ENV")
		}
		if se.Message == "" {
			t.Error("Message should not be empty")
		}
		if se.RecoveryHint == "" {
			t.Error("RecoveryHint should not be empty")
		}
		if se.Context["file"] != "config.yml" {
			t.Errorf("Context[file] = %q, want %q", se.Context["file"], "config.yml")
		}
	})

	t.Run("sdp_error_with_cause", func(t *testing.T) {
		cause := New(ErrGitNotFound, nil)
		err := New(ErrCommandFailed, cause)
		se := ToSerializable(err)

		if se.Cause == "" {
			t.Error("Cause should not be empty")
		}
		if !strings.Contains(se.Cause, "ENV001") {
			t.Errorf("Cause should contain 'ENV001', got %q", se.Cause)
		}
	})

	t.Run("non_sdp_error", func(t *testing.T) {
		err := json.Unmarshal([]byte("invalid"), &struct{}{})
		se := ToSerializable(err)

		if se.Code != "RUNTIME006" {
			t.Errorf("Code = %q, want %q", se.Code, "RUNTIME006")
		}
		if se.Class != "RUNTIME" {
			t.Errorf("Class = %q, want %q", se.Class, "RUNTIME")
		}
	})

	t.Run("nil_error", func(t *testing.T) {
		se := ToSerializable(nil)
		if se != nil {
			t.Error("ToSerializable(nil) should return nil")
		}
	})
}

func TestFromSerializable(t *testing.T) {
	t.Run("full_error", func(t *testing.T) {
		se := &SerializableError{
			Code:         "ENV001",
			Class:        "ENV",
			Message:      "Git is not installed",
			RecoveryHint: "Install Git",
			Context:      map[string]string{"file": "config.yml"},
			Cause:        "underlying error",
		}

		err := FromSerializable(se)
		if err.Code != ErrGitNotFound {
			t.Errorf("Code = %v, want %v", err.Code, ErrGitNotFound)
		}
		if err.Message != "Git is not installed" {
			t.Errorf("Message = %q, want %q", err.Message, "Git is not installed")
		}
		if err.Context["file"] != "config.yml" {
			t.Errorf("Context[file] = %q, want %q", err.Context["file"], "config.yml")
		}
	})

	t.Run("invalid_code", func(t *testing.T) {
		se := &SerializableError{
			Code:    "INVALID",
			Message: "test",
		}

		err := FromSerializable(se)
		if err.Code != ErrInternalError {
			t.Errorf("Invalid code should default to ErrInternalError, got %v", err.Code)
		}
	})

	t.Run("nil_input", func(t *testing.T) {
		err := FromSerializable(nil)
		if err != nil {
			t.Error("FromSerializable(nil) should return nil")
		}
	})
}

func TestMarshalJSON(t *testing.T) {
	err := New(ErrGitNotFound, nil).WithContext("file", "config.yml")
	data, marshalErr := json.Marshal(err)

	if marshalErr != nil {
		t.Fatalf("MarshalJSON failed: %v", marshalErr)
	}

	// Verify it's valid JSON and contains expected fields
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"code"`) {
		t.Error("JSON should contain 'code' field")
	}
	if !strings.Contains(jsonStr, `"ENV001"`) {
		t.Error("JSON should contain 'ENV001' code")
	}
	if !strings.Contains(jsonStr, `"context"`) {
		t.Error("JSON should contain 'context' field")
	}
}

func TestUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"code": "ENV001",
		"class": "ENV",
		"message": "Git is not installed",
		"recovery_hint": "Install Git",
		"context": {"file": "config.yml"}
	}`

	var err SDPError
	if unmarshalErr := json.Unmarshal([]byte(jsonStr), &err); unmarshalErr != nil {
		t.Fatalf("UnmarshalJSON failed: %v", unmarshalErr)
	}

	if err.Code != ErrGitNotFound {
		t.Errorf("Code = %v, want %v", err.Code, ErrGitNotFound)
	}
	if err.Message != "Git is not installed" {
		t.Errorf("Message = %q, want %q", err.Message, "Git is not installed")
	}
	if err.Context["file"] != "config.yml" {
		t.Errorf("Context[file] = %q, want %q", err.Context["file"], "config.yml")
	}
}

func TestToJSON(t *testing.T) {
	err := New(ErrGitNotFound, nil)
	jsonStr, jsonErr := ToJSON(err)

	if jsonErr != nil {
		t.Fatalf("ToJSON failed: %v", jsonErr)
	}

	if !strings.Contains(jsonStr, "ENV001") {
		t.Errorf("JSON should contain 'ENV001', got %s", jsonStr)
	}
}

func TestFromJSON(t *testing.T) {
	jsonStr := `{
		"code": "ENV001",
		"class": "ENV",
		"message": "Git is not installed",
		"recovery_hint": "Install Git"
	}`

	err, jsonErr := FromJSON(jsonStr)
	if jsonErr != nil {
		t.Fatalf("FromJSON failed: %v", jsonErr)
	}

	if err.Code != ErrGitNotFound {
		t.Errorf("Code = %v, want %v", err.Code, ErrGitNotFound)
	}
	if err.Message != "Git is not installed" {
		t.Errorf("Message = %q, want %q", err.Message, "Git is not installed")
	}
}

func TestFromJSON_Invalid(t *testing.T) {
	_, err := FromJSON("invalid json")
	if err == nil {
		t.Error("FromJSON with invalid JSON should return error")
	}
}

func TestValidateSerializationContract(t *testing.T) {
	t.Run("valid_serialization", func(t *testing.T) {
		err := New(ErrGitNotFound, nil).WithContext("file", "config.yml")
		if contractErr := ValidateSerializationContract(err); contractErr != nil {
			t.Errorf("ValidateSerializationContract failed: %v", contractErr)
		}
	})

	t.Run("with_cause", func(t *testing.T) {
		cause := New(ErrFileNotWritable, nil)
		err := New(ErrGitNotFound, cause)
		if contractErr := ValidateSerializationContract(err); contractErr != nil {
			t.Errorf("ValidateSerializationContract failed: %v", contractErr)
		}
	})

	t.Run("all_error_codes", func(t *testing.T) {
		codes := []ErrorCode{
			ErrGitNotFound, ErrInvalidWorkstreamID, ErrBlockedWorkstream,
			ErrCoverageLow, ErrCommandFailed,
		}

		for _, code := range codes {
			t.Run(string(code), func(t *testing.T) {
				err := New(code, nil)
				if contractErr := ValidateSerializationContract(err); contractErr != nil {
					t.Errorf("Serialization contract failed for %s: %v", code, contractErr)
				}
			})
		}
	})
}

func TestRoundTrip(t *testing.T) {
	// Test round-trip serialization for various error types
	testCases := []struct {
		name string
		err  *SDPError
	}{
		{
			name: "simple_error",
			err:  New(ErrGitNotFound, nil),
		},
		{
			name: "error_with_context",
			err:  New(ErrGitNotFound, nil).WithContext("file", "test.yml").WithContext("line", "42"),
		},
		{
			name: "error_with_cause",
			err:  New(ErrCommandFailed, New(ErrGitNotFound, nil)),
		},
		{
			name: "custom_message",
			err:  Newf(ErrGitNotFound, "Custom message: %s", "detail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Serialize
			jsonStr, err := ToJSON(tc.err)
			if err != nil {
				t.Fatalf("ToJSON failed: %v", err)
			}

			// Deserialize
			parsed, err := FromJSON(jsonStr)
			if err != nil {
				t.Fatalf("FromJSON failed: %v", err)
			}

			// Verify
			if parsed.Code != tc.err.Code {
				t.Errorf("Code mismatch: got %v, want %v", parsed.Code, tc.err.Code)
			}
			if parsed.Message != tc.err.Message {
				t.Errorf("Message mismatch: got %q, want %q", parsed.Message, tc.err.Message)
			}

			// Verify context if present
			if tc.err.Context != nil {
				for k, v := range tc.err.Context {
					if parsed.Context[k] != v {
						t.Errorf("Context[%s] = %q, want %q", k, parsed.Context[k], v)
					}
				}
			}
		})
	}
}
