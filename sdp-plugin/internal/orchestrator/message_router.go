package orchestrator

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeProgress MessageType = "progress"
	MessageTypeError    MessageType = "error"
	MessageTypeCommand  MessageType = "command"
	MessageTypeResponse MessageType = "response"
)

// Message represents a message between orchestrator and agents
type Message struct {
	ID        string                 `json:"id"`
	Type      MessageType            `json:"type"`
	From      string                 `json:"from"`
	To        string                 `json:"to"`
	Timestamp time.Time              `json:"timestamp"`
	Content   map[string]interface{} `json:"content"`
}

// MessageRouter handles message routing between orchestrator and agents
type MessageRouter struct {
	mu        sync.RWMutex
	channels  map[string]chan Message
	log       []Message
	listeners map[string][]func(Message)
}

// NewMessageRouter creates a new message router
func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		channels:  make(map[string]chan Message),
		log:       []Message{},
		listeners: make(map[string][]func(Message)),
	}
}

// Register registers a channel for an agent
func (mr *MessageRouter) Register(agentID string) chan Message {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	ch := make(chan Message, 100)
	mr.channels[agentID] = ch
	return ch
}

// Unregister removes an agent's channel
func (mr *MessageRouter) Unregister(agentID string) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if ch, ok := mr.channels[agentID]; ok {
		close(ch)
		delete(mr.channels, agentID)
	}
}

// Send sends a message to an agent
func (mr *MessageRouter) Send(msg Message) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	// Set timestamp if not set
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	// Log message
	mr.log = append(mr.log, msg)

	// Find recipient channel
	ch, ok := mr.channels[msg.To]
	if !ok {
		return fmt.Errorf("recipient not found: %s", msg.To)
	}

	// Send message
	select {
	case ch <- msg:
		return nil
	default:
		return fmt.Errorf("channel full for recipient: %s", msg.To)
	}
}

// Receive receives a message from an agent's channel
func (mr *MessageRouter) Receive(agentID string) (Message, error) {
	mr.mu.RLock()
	ch, ok := mr.channels[agentID]
	mr.mu.RUnlock()

	if !ok {
		return Message{}, fmt.Errorf("agent not registered: %s", agentID)
	}

	select {
	case msg := <-ch:
		return msg, nil
	case <-time.After(30 * time.Second):
		return Message{}, fmt.Errorf("timeout receiving message for agent: %s", agentID)
	}
}

// SendProgress sends a progress update from agent to orchestrator
func (mr *MessageRouter) SendProgress(from string, progress float64, message string) error {
	msg := Message{
		ID:   generateMessageID(),
		Type: MessageTypeProgress,
		From: from,
		To:   "orchestrator",
		Content: map[string]interface{}{
			"progress": progress,
			"message":  message,
		},
	}
	return mr.Send(msg)
}

// SendError sends an error from agent to orchestrator
func (mr *MessageRouter) SendError(from string, errMsg string) error {
	msg := Message{
		ID:   generateMessageID(),
		Type: MessageTypeError,
		From: from,
		To:   "orchestrator",
		Content: map[string]interface{}{
			"error": errMsg,
		},
	}
	return mr.Send(msg)
}

// SendCommand sends a command from orchestrator to agent
func (mr *MessageRouter) SendCommand(to string, command string, args map[string]interface{}) error {
	msg := Message{
		ID:   generateMessageID(),
		Type: MessageTypeCommand,
		From: "orchestrator",
		To:   to,
		Content: map[string]interface{}{
			"command": command,
			"args":    args,
		},
	}
	return mr.Send(msg)
}

// GetLog returns the message log
func (mr *MessageRouter) GetLog() []Message {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	// Return copy of log
	logCopy := make([]Message, len(mr.log))
	copy(logCopy, mr.log)
	return logCopy
}

// GetLogJSON returns the message log as JSON
func (mr *MessageRouter) GetLogJSON() (string, error) {
	log := mr.GetLog()
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal log: %w", err)
	}
	return string(data), nil
}

// AddListener adds a listener for messages to a recipient
func (mr *MessageRouter) AddListener(recipient string, callback func(Message)) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if mr.listeners[recipient] == nil {
		mr.listeners[recipient] = []func(Message){}
	}
	mr.listeners[recipient] = append(mr.listeners[recipient], callback)
}

// generateMessageID generates a unique message ID
func generateMessageID() string {
	return fmt.Sprintf("msg-%d", time.Now().UnixNano())
}
