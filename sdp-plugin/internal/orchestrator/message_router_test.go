package orchestrator

import (
	"fmt"
	"testing"
)

func TestMessageRouter_NewMessageRouter(t *testing.T) {
	router := NewMessageRouter()

	if router == nil {
		t.Fatal("Expected non-nil router")
	}
}

func TestMessageRouter_Register(t *testing.T) {
	router := NewMessageRouter()

	ch := router.Register("agent-1")

	if ch == nil {
		t.Fatal("Expected non-nil channel")
	}

	// Verify channel receives messages
	msg := Message{ID: "test", To: "agent-1"}
	go func() {
		ch <- msg
	}()

	received, err := router.Receive("agent-1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if received.ID != "test" {
		t.Errorf("Expected message ID 'test', got '%s'", received.ID)
	}
}

func TestMessageRouter_SendProgress(t *testing.T) {
	// AC1: Agents send progress updates
	router := NewMessageRouter()
	router.Register("agent-1")
	router.Register("orchestrator") // Register orchestrator to receive messages

	err := router.SendProgress("agent-1", 50.0, "Half complete")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify message logged
	log := router.GetLog()
	if len(log) != 1 {
		t.Fatalf("Expected 1 message in log, got %d", len(log))
	}

	if log[0].Type != MessageTypeProgress {
		t.Errorf("Expected type 'progress', got '%s'", log[0].Type)
	}

	progress, ok := log[0].Content["progress"].(float64)
	if !ok {
		t.Fatal("Expected progress to be float64")
	}

	if progress != 50.0 {
		t.Errorf("Expected progress 50.0, got %v", progress)
	}
}

func TestMessageRouter_SendError(t *testing.T) {
	// AC1: Errors reported immediately
	router := NewMessageRouter()
	router.Register("agent-1")
	router.Register("orchestrator") // Register orchestrator to receive messages

	err := router.SendError("agent-1", "Something went wrong")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify message logged
	log := router.GetLog()
	if len(log) != 1 {
		t.Fatalf("Expected 1 message in log, got %d", len(log))
	}

	if log[0].Type != MessageTypeError {
		t.Errorf("Expected type 'error', got '%s'", log[0].Type)
	}

	errorMsg, ok := log[0].Content["error"].(string)
	if !ok {
		t.Fatal("Expected error to be string")
	}

	if errorMsg != "Something went wrong" {
		t.Errorf("Expected error 'Something went wrong', got '%s'", errorMsg)
	}
}

func TestMessageRouter_SendCommand(t *testing.T) {
	// AC2: Orchestrator sends pause/resume commands
	router := NewMessageRouter()
	router.Register("agent-1")

	err := router.SendCommand("agent-1", "pause", map[string]interface{}{})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify agent receives message
	received, err := router.Receive("agent-1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if received.Type != MessageTypeCommand {
		t.Errorf("Expected type 'command', got '%s'", received.Type)
	}

	command, ok := received.Content["command"].(string)
	if !ok {
		t.Fatal("Expected command to be string")
	}

	if command != "pause" {
		t.Errorf("Expected command 'pause', got '%s'", command)
	}
}

func TestMessageRouter_MessageLogging(t *testing.T) {
	// AC3: All messages logged for debugging
	router := NewMessageRouter()
	router.Register("agent-1")

	// Send multiple messages
	router.SendProgress("agent-1", 25.0, "Started")
	router.SendProgress("agent-1", 50.0, "Half complete")
	router.SendProgress("agent-1", 75.0, "Almost done")
	router.SendError("agent-1", "Failed")

	// Verify all logged
	log := router.GetLog()
	if len(log) != 4 {
		t.Fatalf("Expected 4 messages in log, got %d", len(log))
	}

	// Verify timestamps
	for i, msg := range log {
		if msg.Timestamp.IsZero() {
			t.Errorf("Message %d: Expected non-zero timestamp", i)
		}
	}
}

func TestMessageRouter_Unregister(t *testing.T) {
	router := NewMessageRouter()

	ch := router.Register("agent-1")
	router.Unregister("agent-1")

	// Verify channel closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("Expected channel to be closed")
		}
	default:
		// Channel empty, OK
	}

	// Verify cannot send to unregistered agent
	err := router.SendProgress("agent-1", 100.0, "Complete")
	if err == nil {
		t.Error("Expected error sending to unregistered agent")
	}
}

func TestMessageRouter_SendToInvalidAgent(t *testing.T) {
	router := NewMessageRouter()

	msg := Message{ID: "test", To: "nonexistent"}
	err := router.Send(msg)

	if err == nil {
		t.Fatal("Expected error sending to invalid agent")
	}
}

func TestMessageRouter_GetLogJSON(t *testing.T) {
	router := NewMessageRouter()
	router.Register("agent-1")

	router.SendProgress("agent-1", 50.0, "Half complete")

	jsonLog, err := router.GetLogJSON()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if jsonLog == "" {
		t.Fatal("Expected non-empty JSON log")
	}

	// Verify it's valid JSON (contains expected fields)
	if len(jsonLog) < 10 {
		t.Error("JSON log too short")
	}
}

func TestMessageRouter_ConcurrentAccess(t *testing.T) {
	router := NewMessageRouter()
	router.Register("agent-1")

	// Send messages concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			router.SendProgress("agent-1", float64(idx*10), fmt.Sprintf("Message %d", idx))
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all messages logged
	log := router.GetLog()
	if len(log) != 10 {
		t.Fatalf("Expected 10 messages, got %d", len(log))
	}
}
