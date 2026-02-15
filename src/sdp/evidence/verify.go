package evidence

import (
	"bytes"
	"errors"
	"io"
)

// VerifyChain checks the hash chain integrity.
func (l *EvidenceLog) VerifyChain() error {
	events, err := l.ReadAll()
	if err != nil {
		return err
	}

	for i, event := range events {
		if i == 0 {
			if event.PrevHash != "" {
				return errors.New("genesis event should have empty prev_hash")
			}
		} else {
			if event.PrevHash != events[i-1].Hash {
				return errors.New("hash chain broken")
			}
		}
	}

	return nil
}

// CopyTo copies all events to a writer.
func (l *EvidenceLog) CopyTo(w io.Writer) error {
	return l.copyToInternal(w)
}

// VerifyHash verifies that the event's hash is correctly calculated.
func VerifyHash(event BaseEvent) bool {
	original := event.Hash
	event.Hash = ""
	event.CalculateHash()
	return event.Hash == original
}

// VerifyEvents verifies all events in a slice.
func VerifyEvents(events []BaseEvent) error {
	for i, event := range events {
		if !VerifyHash(event) {
			return errors.New("event hash mismatch")
		}
		if i == 0 {
			if event.PrevHash != "" {
				return errors.New("genesis event should have empty prev_hash")
			}
		} else {
			if event.PrevHash != events[i-1].Hash {
				return errors.New("hash chain broken")
			}
		}
	}
	return nil
}

// FindEventByID finds an event by its ID.
func FindEventByID(events []BaseEvent, id string) (BaseEvent, int) {
	for i, e := range events {
		if e.ID == id {
			return e, i
		}
	}
	return BaseEvent{}, -1
}

// FilterByType filters events by type.
func FilterByType(events []BaseEvent, eventType EventType) []BaseEvent {
	var result []BaseEvent
	for _, e := range events {
		if e.Type == eventType {
			result = append(result, e)
		}
	}
	return result
}

// EventsToJSONL converts events to JSONL format.
func EventsToJSONL(events []BaseEvent) (string, error) {
	var buf bytes.Buffer
	for _, e := range events {
		line, err := e.ToJSONL()
		if err != nil {
			return "", err
		}
		buf.WriteString(line)
	}
	return buf.String(), nil
}
