package errors

import (
	"encoding/json"
	"fmt"
)

// SerializableError is a JSON-serializable representation of an SDPError.
type SerializableError struct {
	Code         string            `json:"code"`
	Class        string            `json:"class"`
	Message      string            `json:"message"`
	RecoveryHint string            `json:"recovery_hint,omitempty"`
	Context      map[string]string `json:"context,omitempty"`
	Cause        string            `json:"cause,omitempty"`
}

// ToSerializable converts an SDPError to a serializable format.
func ToSerializable(err error) *SerializableError {
	if err == nil {
		return nil
	}

	sdpErr, ok := err.(*SDPError)
	if !ok {
		return &SerializableError{
			Code:    string(ErrInternalError),
			Class:   string(ClassRuntime),
			Message: err.Error(),
		}
	}

	result := &SerializableError{
		Code:         string(sdpErr.Code),
		Class:        string(sdpErr.Class()),
		Message:      sdpErr.Message,
		RecoveryHint: sdpErr.RecoveryHint(),
		Context:      sdpErr.Context,
	}

	if sdpErr.Cause != nil {
		result.Cause = sdpErr.Cause.Error()
	}

	return result
}

// FromSerializable creates an SDPError from a serializable format.
func FromSerializable(se *SerializableError) *SDPError {
	if se == nil {
		return nil
	}

	code := ErrorCode(se.Code)
	if !code.IsValid() {
		code = ErrInternalError
	}

	var cause error
	if se.Cause != "" {
		cause = fmt.Errorf("%s", se.Cause)
	}

	return &SDPError{
		Code:    code,
		Message: se.Message,
		Cause:   cause,
		Context: se.Context,
	}
}

// MarshalJSON implements json.Marshaler for SDPError.
func (e *SDPError) MarshalJSON() ([]byte, error) {
	return json.Marshal(ToSerializable(e))
}

// UnmarshalJSON implements json.Unmarshaler for SDPError.
func (e *SDPError) UnmarshalJSON(data []byte) error {
	var se SerializableError
	if err := json.Unmarshal(data, &se); err != nil {
		return err
	}

	parsed := FromSerializable(&se)
	e.Code = parsed.Code
	e.Message = parsed.Message
	e.Cause = parsed.Cause
	e.Context = parsed.Context
	return nil
}

// ToJSON returns a JSON representation of the error.
func ToJSON(err error) (string, error) {
	data, err := json.MarshalIndent(ToSerializable(err), "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}
	return string(data), nil
}

// FromJSON parses a JSON string into an SDPError.
func FromJSON(jsonStr string) (*SDPError, error) {
	var se SerializableError
	if err := json.Unmarshal([]byte(jsonStr), &se); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return FromSerializable(&se), nil
}

// ValidateSerializationContract ensures errors serialize/deserialize correctly.
func ValidateSerializationContract(sdpErr *SDPError) error {
	// Serialize
	jsonStr, err := ToJSON(sdpErr)
	if err != nil {
		return fmt.Errorf("serialization failed: %w", err)
	}

	// Deserialize
	parsed, err := FromJSON(jsonStr)
	if err != nil {
		return fmt.Errorf("deserialization failed: %w", err)
	}

	// Verify contract
	if parsed.Code != sdpErr.Code {
		return fmt.Errorf("code mismatch: got %s, want %s", parsed.Code, sdpErr.Code)
	}
	if parsed.Message != sdpErr.Message {
		return fmt.Errorf("message mismatch: got %s, want %s", parsed.Message, sdpErr.Message)
	}

	return nil
}
