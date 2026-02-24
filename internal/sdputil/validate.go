package sdputil

import (
	"fmt"
	"regexp"
)

var (
	// wsIDPattern: 00-XXX-YY (e.g. 00-014-01)
	wsIDPattern = regexp.MustCompile(`^[0-9]{2}-[0-9]{3}-[0-9]{2}$`)
	// featureIDPattern: F001-F9999
	featureIDPattern = regexp.MustCompile(`^F[0-9]{3,4}$`)
)

// ValidateFeatureID rejects featureID values that would allow path traversal.
// Format: F001-F9999 (allowlist).
func ValidateFeatureID(featureID string) error {
	if !featureIDPattern.MatchString(featureID) {
		return fmt.Errorf("invalid feature_id %q: must match F001-F9999", featureID)
	}
	return nil
}

// ValidateWSID rejects wsID values that would allow path traversal.
// Format: 00-XXX-YY (e.g. 00-014-01) (allowlist).
func ValidateWSID(wsID string) error {
	if !wsIDPattern.MatchString(wsID) {
		return fmt.Errorf("invalid ws_id %q: must match 00-XXX-YY", wsID)
	}
	return nil
}
