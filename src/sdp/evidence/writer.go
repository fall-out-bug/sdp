package evidence

import (
	"path/filepath"
)

// DefaultLogPath is the default evidence log location.
const DefaultLogPath = ".sdp/log/events.jsonl"

// Writer provides type-safe methods for writing evidence events.
type Writer struct {
	log *EvidenceLog
}

// NewWriter creates a new evidence writer.
func NewWriter(logPath string) (*Writer, error) {
	log, err := NewEvidenceLog(logPath)
	if err != nil {
		return nil, err
	}
	return &Writer{log: log}, nil
}

// DefaultWriter creates a writer with the default log path.
func DefaultWriter() (*Writer, error) {
	return NewWriter(DefaultLogPath)
}

// WriteDecision writes a decision event with auto-chaining.
func (w *Writer) WriteDecision(event *DecisionEvent) error {
	_, err := w.log.AppendEvent(EventTypeDecision, func(b *BaseEvent) {
		b.ID = event.ID
		b.Timestamp = event.Timestamp
	})
	return err
}

// WritePlan writes a plan event with auto-chaining.
func (w *Writer) WritePlan(event *PlanEvent) error {
	_, err := w.log.AppendEvent(EventTypePlan, func(b *BaseEvent) {
		b.ID = event.ID
		b.Timestamp = event.Timestamp
	})
	return err
}

// WriteGeneration writes a generation event with auto-chaining.
func (w *Writer) WriteGeneration(event *GenerationEvent) error {
	_, err := w.log.AppendEvent(EventTypeGeneration, func(b *BaseEvent) {
		b.ID = event.ID
		b.Timestamp = event.Timestamp
	})
	return err
}

// WriteVerification writes a verification event with auto-chaining.
func (w *Writer) WriteVerification(event *VerificationEvent) error {
	_, err := w.log.AppendEvent(EventTypeVerification, func(b *BaseEvent) {
		b.ID = event.ID
		b.Timestamp = event.Timestamp
	})
	return err
}

// WriteAcceptance writes an acceptance event with auto-chaining.
func (w *Writer) WriteAcceptance(event *AcceptanceEvent) error {
	_, err := w.log.AppendEvent(EventTypeAcceptance, func(b *BaseEvent) {
		b.ID = event.ID
		b.Timestamp = event.Timestamp
	})
	return err
}

// WriteApproval writes an approval event with auto-chaining.
func (w *Writer) WriteApproval(event *ApprovalEvent) error {
	_, err := w.log.AppendEvent(EventTypeApproval, func(b *BaseEvent) {
		b.ID = event.ID
		b.Timestamp = event.Timestamp
	})
	return err
}

// Log returns the underlying evidence log.
func (w *Writer) Log() *EvidenceLog {
	return w.log
}

// Path returns the log file path.
func (w *Writer) Path() string {
	return w.log.Path()
}

// GetProjectLogPath returns the evidence log path for a project root.
func GetProjectLogPath(projectRoot string) string {
	return filepath.Join(projectRoot, DefaultLogPath)
}
