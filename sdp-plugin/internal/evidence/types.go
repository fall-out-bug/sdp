package evidence

// Event is the base evidence log event (AC5, AC6).
type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"` // plan, generation, verification, approval, decision, lesson
	Timestamp string      `json:"timestamp"`
	WSID      string      `json:"ws_id"`
	CommitSHA string      `json:"commit_sha,omitempty"`
	PrevHash  string      `json:"prev_hash,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// GenerationData is provenance for generation events (AC3).
type GenerationData struct {
	ModelID      string   `json:"model_id"`
	ModelVersion string   `json:"model_version"`
	PromptHash   string   `json:"prompt_hash"`
	FilesChanged []string `json:"files_changed,omitempty"`
}

// VerificationData is pass/fail for verification events (AC4).
type VerificationData struct {
	Passed   bool    `json:"passed"`
	GateName string  `json:"gate_name,omitempty"`
	Coverage float64 `json:"coverage,omitempty"`
}

// DecisionEventData is the payload for decision events (AC9).
type DecisionEventData struct {
	Question     string   `json:"question"`
	Choice       string   `json:"choice"`
	Rationale    string   `json:"rationale"`
	Alternatives []string `json:"alternatives,omitempty"`
	Confidence   float64  `json:"confidence,omitempty"`
	Reverses     *string  `json:"reverses,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

// LessonEventData is the payload for lesson events (AC10).
type LessonEventData struct {
	Category         string   `json:"category"`
	Insight          string   `json:"insight"`
	SourceWSID       string   `json:"source_ws_id"`
	Outcome          string   `json:"outcome"` // worked, failed, mixed
	RelatedDecisions []string `json:"related_decisions,omitempty"`
}

// EventTypes are the six evidence event types (AC2).
var EventTypes = []string{"plan", "generation", "verification", "approval", "decision", "lesson"}
