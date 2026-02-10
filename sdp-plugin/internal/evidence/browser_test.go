package evidence

import (
	"testing"
	"time"
)

func TestBrowser_Pagination(t *testing.T) {
	// Setup: Create 45 events
	events := make([]Event, 45)
	now := time.Now()
	for i := 0; i < 45; i++ {
		events[i] = Event{
			ID:        string(rune('a' + i)),
			Type:      "generation",
			Timestamp: now.Add(time.Duration(i) * time.Minute).Format(time.RFC3339),
			WSID:      "00-054-03",
		}
	}

	b := NewBrowser(events)

	// AC1: Page 1 should have 20 events
	page1, total1 := b.Page(1, 20)
	if len(page1) != 20 {
		t.Errorf("Expected 20 events on page 1, got %d", len(page1))
	}
	if total1 != 45 {
		t.Errorf("Expected total 45 events, got %d", total1)
	}

	// AC1: Page 2 should have 20 events
	page2, total2 := b.Page(2, 20)
	if len(page2) != 20 {
		t.Errorf("Expected 20 events on page 2, got %d", len(page2))
	}
	if total2 != 45 {
		t.Errorf("Expected total 45 events, got %d", total2)
	}

	// AC1: Page 3 should have 5 events
	page3, total3 := b.Page(3, 20)
	if len(page3) != 5 {
		t.Errorf("Expected 5 events on page 3, got %d", len(page3))
	}
	if total3 != 45 {
		t.Errorf("Expected total 45 events, got %d", total3)
	}

	// AC1: Page 4 should be empty
	page4, total4 := b.Page(4, 20)
	if len(page4) != 0 {
		t.Errorf("Expected 0 events on page 4, got %d", len(page4))
	}
	if total4 != 45 {
		t.Errorf("Expected total 45 events, got %d", total4)
	}
}

func TestBrowser_FilterByType(t *testing.T) {
	events := []Event{
		{Type: "generation", WSID: "00-054-01"},
		{Type: "decision", WSID: "00-054-02"},
		{Type: "generation", WSID: "00-054-03"},
		{Type: "verification", WSID: "00-054-03"},
	}

	b := NewBrowser(events)

	// AC2: Filter by type=generation
	filtered := b.FilterByType("generation")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 generation events, got %d", len(filtered))
	}
	for _, e := range filtered {
		if e.Type != "generation" {
			t.Errorf("Expected all events to be generation, got %s", e.Type)
		}
	}
}

func TestBrowser_FilterByModel(t *testing.T) {
	events := []Event{
		{
			Type: "generation",
			WSID: "00-054-01",
			Data: map[string]interface{}{"model_id": "claude-sonnet-4"},
		},
		{
			Type: "generation",
			WSID: "00-054-02",
			Data: map[string]interface{}{"model_id": "claude-opus-4"},
		},
		{
			Type: "decision",
			WSID: "00-054-03",
			Data: map[string]interface{}{"choice": "use Go"},
		},
	}

	b := NewBrowser(events)

	// AC3: Filter by model=claude-sonnet-4
	filtered := b.FilterByModel("claude-sonnet-4")
	if len(filtered) != 1 {
		t.Errorf("Expected 1 claude-sonnet-4 event, got %d", len(filtered))
	}
	if filtered[0].WSID != "00-054-01" {
		t.Errorf("Expected WS 00-054-01, got %s", filtered[0].WSID)
	}
}

func TestBrowser_FilterBySince(t *testing.T) {
	now := time.Now()
	events := []Event{
		{Timestamp: now.Add(-48 * time.Hour).Format(time.RFC3339), WSID: "00-054-01"},
		{Timestamp: now.Add(-24 * time.Hour).Format(time.RFC3339), WSID: "00-054-02"},
		{Timestamp: now.Format(time.RFC3339), WSID: "00-054-03"},
	}

	b := NewBrowser(events)

	// AC4: Filter by date since 2026-02-01 (using a relative cutoff)
	since := now.Add(-36 * time.Hour).Format(time.RFC3339)
	filtered := b.FilterBySince(since)
	if len(filtered) != 2 {
		t.Errorf("Expected 2 events since cutoff, got %d", len(filtered))
	}
}

func TestBrowser_FilterByWS(t *testing.T) {
	events := []Event{
		{WSID: "00-054-01", Type: "generation"},
		{WSID: "00-054-02", Type: "decision"},
		{WSID: "00-054-03", Type: "generation"},
		{WSID: "00-054-03", Type: "verification"},
	}

	b := NewBrowser(events)

	// AC5: Filter by ws=00-054-03
	filtered := b.FilterByWS("00-054-03")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 events for WS 00-054-03, got %d", len(filtered))
	}
	for _, e := range filtered {
		if e.WSID != "00-054-03" {
			t.Errorf("Expected all events to be from WS 00-054-03, got %s", e.WSID)
		}
	}
}
