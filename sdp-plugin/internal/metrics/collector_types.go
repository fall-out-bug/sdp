package metrics

// Metrics represents collected metrics from evidence events (AC2-AC5).
type Metrics struct {
	CatchRate           float64            `json:"catch_rate"`
	TotalVerifications  int                `json:"total_verifications"`
	FailedVerifications int                `json:"failed_verifications"`
	IterationCount      map[string]int     `json:"iteration_count"`
	ModelPassRate       map[string]float64 `json:"model_pass_rate"`
	TotalApprovals      int                `json:"total_approvals"`
	FailedApprovals     int                `json:"failed_approvals"`
	AcceptanceCatchRate float64            `json:"acceptance_catch_rate"`
}

// workstreamState tracks workstream iteration state.
type workstreamState struct {
	generationCount int
	lastWasGen      bool
}

// modelVerificationStats tracks verification stats per model.
type modelVerificationStats struct {
	Passed   int
	Total    int
	PassRate float64
}

// evidenceEvent represents an evidence log event.
type evidenceEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	WSID      string                 `json:"ws_id"`
	Data      map[string]interface{} `json:"data"`
}
