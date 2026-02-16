package coordination

import "errors"

var (
	// ErrMissingID indicates the event ID is missing
	ErrMissingID = errors.New("missing event ID")
	// ErrMissingAgentID indicates the agent ID is missing
	ErrMissingAgentID = errors.New("missing agent ID")
	// ErrInvalidEventType indicates the event type is invalid
	ErrInvalidEventType = errors.New("invalid event type")
	// ErrHashChainBroken indicates the hash chain is broken
	ErrHashChainBroken = errors.New("hash chain broken")
)
