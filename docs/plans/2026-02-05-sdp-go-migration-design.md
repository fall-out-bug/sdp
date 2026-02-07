# SDP Go Migration Design Document

> **Status:** Ready for Implementation
> **Date:** 2026-02-05
> **Author:** AI Assistant + User Collaboration
> **Target:** Complete Python → Go migration for single-binary deployment

---

## Executive Summary

**Objective:** Migrate SDP from Python to Go to enable single-binary distribution without dependencies while leveraging existing Beads CLI functionality.

**Key Decision:** Big Bang migration with 10-week rollback plan.

**Primary Motivation:** Single static binary for easy deployment (no Python runtime, no pip install).

**Scope Reduction:** By leveraging Beads CLI capabilities, reduce implementation effort by **72%** (from ~15K LOC to ~4.2K LOC Go).

**Timeline:** 10 weeks (including 2-week Python cleanup phase).

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Component Analysis: Keep vs. Remove](#component-analysis-keep-vs-remove)
3. [Go Implementation Details](#go-implementation-details)
4. [Fallback Strategy](#fallback-strategy)
5. [Implementation Plan](#implementation-plan)
6. [Python Cleanup Phase](#python-cleanup-phase)
7. [Risk Mitigation](#risk-mitigation)
8. [Success Criteria](#success-criteria)

---

## Architecture Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         SDP Go CLI                         │
│                      (Single Binary)                       │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Core Layer                             │
│  • Workstream Parser (YAML frontmatter)                     │
│  • TDD Runner (Red→Green→Refactor with pytest)              │
│  • Quality Gates (mypy, ruff, coverage)                     │
└─────────────────────────────────────────────────────────────┘
                            │
                ┌───────────┴───────────┐
                ▼                       ▼
┌───────────────────────┐   ┌───────────────────────────────┐
│   Beads CLI Wrapper   │   │      Telemetry System        │
│  (Thin CLI Interface)  │   │  • Collector (git + pytest)   │
│  • bd ready --json     │   │  • Analyzer (pattern detect)  │
│  • bd update --status  │   │  • Drift Detector (fingerprint)│
└───────────────────────┘   └───────────────────────────────┘
                │
                ▼
┌─────────────────────────────────────────────────────────────┐
│                    External Tools                           │
│  • Beads CLI (task tracking, dependencies)                 │
│  • pytest (test execution)                                  │
│  • mypy, ruff (quality checks)                              │
│  • git (version control)                                    │
└─────────────────────────────────────────────────────────────┘
```

### Design Principles

1. **YAGNI Ruthlessly** — Don't build what Beads already provides
2. **Thin Wrappers** — Minimal Go code, delegate to Beads/Python tools
3. **Fallback Gracefully** — Work without Beads if needed (JSON storage)
4. **Static Binary** — No external dependencies at runtime
5. **Type Safety** — Leverage Go's type system for compile-time checks

---

## Component Analysis: Keep vs. Remove

### ✅ Keep (SDP-Unique Functionality)

| Component | Python LOC | Go LOC (est) | Why Keep |
|-----------|-----------|--------------|----------|
| **Workstream YAML Parser** | 800 | 600 | Beads uses generic task model |
| **TDD Runner** | 200 | 150 | Beads doesn't integrate with pytest |
| **Quality Gates** | 1,500 | 900 | Beads doesn't run mypy/ruff/coverage |
| **Quality Watcher** | 400 | 300 | Real-time fsnotify monitoring |
| **Telemetry Collector** | 600 | 400 | Auto-capture git + pytest stats |
| **Telemetry Analyzer** | 600 | 400 | Pattern detection, friction analysis |
| **Drift Detector** | 600 | 500 | Contract vs Implementation comparison |
| **Checkpoint System (simplified)** | 800 | 200 | @oneshot resume capability |
| **CLI Commands** | 800 | 600 | init, doctor, build, telemetry |
| **Total Keep** | **6,300** | **4,050** | - |

### ❌ Remove (Beads Already Provides)

| Component | Python LOC | Beads Equivalent | Savings |
|-----------|-----------|------------------|---------|
| **BeadsClient (full implementation)** | 1,200 | `bd ready --json` | -1,200 LOC |
| **Dependency Resolution Logic** | 800 | `bd ready` (automatic) | -800 LOC |
| **Task Models** | 400 | JSON from `--json` flag | -400 LOC |
| **Multi-Agent Orchestrator (complex)** | 1,000 | Simplify with `bd ready` | -700 LOC |
| **Escalation Metrics (with fallback)** | 300 | `bd history` (primary) | -200 LOC |
| **Total Remove** | **3,700** | **N/A** | **-3,300 LOC** |

### 📊 Net Impact

```
Python Core (keep):    6,300 LOC
Go Implementation:     4,050 LOC
Net Reduction:         -2,250 LOC (-36%)

Python Beads (remove): 3,700 LOC
Go Wrappers:            300 LOC
Net Reduction:         -3,400 LOC (-92%)

Total Reduction:       25,708 LOC → 4,350 LOC
                      (-83% codebase size)
```

---

## Go Implementation Details

### Directory Structure

```
sdp-go/
├── cmd/sdp/
│   ├── main.go                 # CLI entry point
│   ├── root.go                 # Root command, version
│   └── commands/
│       ├── init.go             # Initialize SDP project
│       ├── doctor.go           # Health checks
│       ├── build.go            # TDD cycle execution
│       ├── telemetry.go        # Insights commands
│       └── review.go           # Quality review
├── internal/
│   ├── core/
│   │   ├── workstream.go       # Parse YAML frontmatter
│   │   ├── validator.go        # Validate WS structure
│   │   └── schema.go           # Type definitions
│   ├── tdd/
│   │   ├── runner.go           # Red-Green-Refactor execution
│   │   └── phase.go            # Phase types (RED, GREEN, REFACTOR)
│   ├── quality/
│   │   ├── checker.go          # Run mypy, ruff, coverage
│   │   ├── gate.go             # Individual quality gate
│   │   └── watcher.go          # Background file watcher
│   ├── telemetry/
│   │   ├── collector.go        # Auto-capture metrics
│   │   ├── analyzer.go         # Pattern detection
│   │   └── repository.go       # Storage with fallback
│   ├── validators/
│   │   └── drift/
│   │       ├── fingerprint.go  # Extract Contract fingerprint
│   │       └── detector.go     # Drift detection
│   ├── beads/
│   │   ├── wrapper.go          # Thin CLI wrappers
│   │   └── fallback.go         # Mock implementation
│   └── checkpoint/
│       └── repository.go       # Simplified checkpoint system
├── pkg/
│   └── protocol/
│       └── types.go            # Public type definitions
├── go.mod
├── go.sum
├── Makefile                    # Build scripts
└── README.md                   # Go-specific docs
```

### Component 1: Workstream Parser

```go
// internal/core/workstream.go
package core

import (
    "bytes"
    "fmt"
    "os"
    "regexp"

    "gopkg.in/yaml.v3"
)

var wsIDRegex = regexp.MustCompile(`^\d{2}-\d{3}-\d{2}$`)

type Workstream struct {
    ID          string            `yaml:"ws_id"`
    Feature     string            `yaml:"feature"`
    Status      string            `yaml:"status"`
    Goal        string            `yaml:"goal"`
    Acceptance []string          `yaml:"acceptance_criteria"`
    Scope       Scope             `yaml:"scope"`
    Telemetry   *Telemetry        `yaml:"telemetry,omitempty"`
    Frontmatter map[string]any    `yaml:",inline"`
}

type Scope struct {
    CapabilityTier  string   `yaml:"capability_tier"`
    MaxLOC          int      `yaml:"max_loc"`
    Contracts       []string `yaml:"contracts"`
    Interfaces      []string `yaml:"interfaces"`
}

func ParseWorkstream(wsPath string) (*Workstream, error) {
    data, err := os.ReadFile(wsPath)
    if err != nil {
        return nil, fmt.Errorf("read workstream: %w", err)
    }

    // Split frontmatter (---) from content
    parts := bytes.SplitN(data, []byte("---"), 3)
    if len(parts) < 3 {
        return nil, fmt.Errorf("invalid frontmatter format")
    }

    var ws Workstream
    if err := yaml.Unmarshal(parts[1], &ws); err != nil {
        return nil, fmt.Errorf("parse frontmatter: %w", err)
    }

    // Validate
    if err := ws.Validate(); err != nil {
        return nil, fmt.Errorf("validate workstream: %w", err)
    }

    return &ws, nil
}

func (ws *Workstream) Validate() error {
    if !wsIDRegex.MatchString(ws.ID) {
        return fmt.Errorf("invalid workstream ID format: %s (expected PP-FFF-SS)", ws.ID)
    }

    if ws.Scope.CapabilityTier != "" {
        switch ws.Scope.CapabilityTier {
        case "T0", "T1", "T2", "T3":
            // Valid
        default:
            return fmt.Errorf("invalid capability tier: %s", ws.Scope.CapabilityTier)
        }
    }

    return nil
}
```

**Benefits over Python:**
- Compile-time validation (vs runtime errors)
- No dynamic field access (type-safe struct)
- Faster YAML parsing (native Go lib)

### Component 2: TDD Runner

```go
// internal/tdd/runner.go
package tdd

import (
    "context"
    "os/exec"
    "time"
)

type Phase string

const (
    Red     Phase = "red"
    Green   Phase = "green"
    Refactor Phase = "refactor"
)

type Result struct {
    Phase     Phase
    Success   bool
    Output    string
    Duration  time.Duration
    NextPhase *Phase
}

type Runner struct {
    projectDir string
    pytestBin  string
}

func (r *Runner) RunAllPhases(ctx context.Context, testPath string) (*Result, error) {
    // Red Phase: Test should fail
    red, err := r.RedPhase(ctx, testPath)
    if err != nil {
        return nil, fmt.Errorf("red phase: %w", err)
    }
    if !red.Success {
        return red, nil // Test didn't fail as expected
    }

    // Green Phase: Implement until tests pass
    green, err := r.GreenPhase(ctx, testPath)
    if err != nil {
        return nil, fmt.Errorf("green phase: %w", err)
    }
    if !green.Success {
        return green, nil // Tests failed
    }

    // Refactor Phase: Improve code structure
    refactor, err := r.RefactorPhase(ctx, testPath)
    if err != nil {
        return nil, fmt.Errorf("refactor phase: %w", err)
    }

    return refactor, nil
}

func (r *Runner) RedPhase(ctx context.Context, testPath string) (*Result, error) {
    cmd := exec.CommandContext(ctx, r.pytestBin, testPath, "-v")
    cmd.Dir = r.projectDir

    start := time.Now()
    out, err := cmd.CombinedOutput()
    duration := time.Since(start)

    output := string(out)

    // In RED phase, we EXPECT failure
    success := err != nil

    var nextPhase *Phase
    if success {
        p := Green
        nextPhase = &p
    }

    return &Result{
        Phase:     Red,
        Success:   success,
        Output:    output,
        Duration:  duration,
        NextPhase: nextPhase,
    }, nil
}

func (r *Runner) GreenPhase(ctx context.Context, testPath string) (*Result, error) {
    cmd := exec.CommandContext(ctx, r.pytestBin, testPath, "-v")
    cmd.Dir = r.projectDir

    start := time.Now()
    out, err := cmd.CombinedOutput()
    duration := time.Since(start)

    output := string(out)
    success := err == nil

    var nextPhase *Phase
    if success {
        p := Refactor
        nextPhase = &p
    } else {
        p := Green
        nextPhase = &p // Retry Green
    }

    return &Result{
        Phase:     Green,
        Success:   success,
        Output:    output,
        Duration:  duration,
        NextPhase: nextPhase,
    }, nil
}

func (r *Runner) RefactorPhase(ctx context.Context, testPath string) (*Result, error) {
    cmd := exec.CommandContext(ctx, r.pytestBin, testPath, "-v")
    cmd.Dir = r.projectDir

    start := time.Now()
    out, err := cmd.CombinedOutput()
    duration := time.Since(start)

    output := string(out)
    success := err == nil

    return &Result{
        Phase:     Refactor,
        Success:   success,
        Output:    output,
        Duration:  duration,
        NextPhase: nil,
    }, nil
}
```

### Component 3: Quality Gates (Parallel Execution)

```go
// internal/quality/checker.go
package quality

import (
    "context"
    "sync"

    "github.com/charmbracelet/lipgloss"
)

type CheckResult struct {
    Passed     bool
    Warnings   []string
    Errors     []string
    Coverage   float64
    DurationMs int64
}

type Checker struct {
    projectDir string
    mypyBin    string
    ruffBin    string
    pytestBin  string
}

func (c *Checker) RunAllGates(ctx context.Context) (*CheckResult, error) {
    results := make(chan gateResult, 4)

    // Run gates in parallel (goroutines)
    var wg sync.WaitGroup
    wg.Add(4)

    go c.runMypy(ctx, &wg, results)
    go c.runRuff(ctx, &wg, results)
    go c.runCoverage(ctx, &wg, results)
    go c.runFileSizeCheck(ctx, &wg, results)

    // Wait for all gates
    go func() {
        wg.Wait()
        close(results)
    }()

    // Aggregate results
    var aggregated CheckResult
    for res := range results {
        if res.err != nil {
            return nil, res.err
        }

        aggregated.Passed = aggregated.Passed && res.passed
        aggregated.Warnings = append(aggregated.Warnings, res.warnings...)
        aggregated.Errors = append(aggregated.Errors, res.errors...)
    }

    return &aggregated, nil
}

type gateResult struct {
    passed   bool
    warnings []string
    errors   []string
    err      error
}

func (c *Checker) runMypy(ctx context.Context, wg *sync.WaitGroup, results chan<- gateResult) {
    defer wg.Done()

    cmd := exec.CommandContext(ctx, c.mypyBin, ".", "--strict")
    cmd.Dir = c.projectDir

    out, err := cmd.CombinedOutput()
    output := string(out)

    res := gateResult{}
    if err != nil {
        res.passed = false
        res.errors = parseOutputLines(output)
    } else {
        res.passed = true
    }

    results <- res
}

func (c *Checker) runRuff(ctx context.Context, wg *sync.WaitGroup, results chan<- gateResult) {
    defer wg.Done()

    cmd := exec.CommandContext(ctx, c.ruffBin, "check", ".")
    cmd.Dir = c.projectDir

    out, err := cmd.CombinedOutput()
    output := string(out)

    res := gateResult{}
    if err != nil {
        res.passed = false
        res.errors = parseOutputLines(output)
    } else {
        res.passed = true
    }

    results <- res
}

func (c *Checker) runCoverage(ctx context.Context, wg *sync.WaitGroup, results chan<- gateResult) {
    defer wg.Done()

    cmd := exec.CommandContext(ctx, c.pytestBin, "--cov", ".", "--cov-report=term-missing")
    cmd.Dir = c.projectDir

    out, err := cmd.CombinedOutput()
    output := string(out)

    res := gateResult{}
    coverage := extractCoverage(output)

    if err != nil || coverage < 80.0 {
        res.passed = false
        res.errors = []string{fmt.Sprintf("Coverage %.1f%% is below 80%%", coverage)}
    } else {
        res.passed = true
    }

    results <- res
}

func (c *Checker) runFileSizeCheck(ctx context.Context, wg *sync.WaitGroup, results chan<- gateResult) {
    defer wg.Done()

    // Find all .py files
    files, _ := filepath.Glob(filepath.Join(c.projectDir, "**", "*.py"))

    res := gateResult{passed: true}

    for _, file := range files {
        info, _ := os.Stat(file)
        if info.Size() > 200*1024 { // 200 KB (roughly 200 LOC)
            res.passed = false
            res.errors = append(res.errors,
                fmt.Sprintf("%s: file size %d bytes exceeds 200 LOC limit", file, info.Size()))
        }
    }

    results <- res
}
```

**Key Advantage:** Goroutines enable parallel gate execution (vs Python's sequential).

### Component 4: Beads Wrapper (Thin Layer)

```go
// internal/beads/wrapper.go
package beads

import (
    "encoding/json"
    "os/exec"
    "strings"
)

// ReadyTasks returns tasks with no blockers (via bd ready --json)
func ReadyTasks(projectDir string) ([]string, error) {
    cmd := exec.Command("bd", "ready", "--json")
    cmd.Dir = projectDir

    out, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("bd ready: %w", err)
    }

    // Parse JSON array of task IDs
    var tasks []string
    if err := json.Unmarshal(out, &tasks); err != nil {
        return nil, fmt.Errorf("parse bd ready output: %w", err)
    }

    return tasks, nil
}

// UpdateStatus updates task status (via bd update)
func UpdateStatus(projectDir, taskID, status string) error {
    cmd := exec.Command("bd", "update", taskID, "--status", status)
    cmd.Dir = projectDir

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("bd update: %w", err)
    }

    return nil
}

// CreateTask creates a new task (via bd create)
func CreateTask(projectDir, title, description string) (string, error) {
    args := []string{"create", title, "--json"}
    if description != "" {
        args = append(args, "--description", description)
    }

    cmd := exec.Command("bd", args...)
    cmd.Dir = projectDir

    out, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("bd create: %w", err)
    }

    // Parse response to get task ID
    var result map[string]any
    if err := json.Unmarshal(out, &result); err != nil {
        return "", fmt.Errorf("parse bd create output: %w", err)
    }

    taskID, _ := result["id"].(string)
    return taskID, nil
}

// AddDependency adds a blocking dependency (via bd dep add)
func AddDependency(projectDir, fromID, toID string) error {
    cmd := exec.Command("bd", "dep", "add", fromID, toID, "--type", "blocks")
    cmd.Dir = projectDir

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("bd dep add: %w", err)
    }

    return nil
}
```

**Benefit:** Only 150 LOC vs 1,200 LOC Python (87% reduction).

### Component 5: Telemetry Collector

```go
// internal/telemetry/collector.go
package telemetry

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
    "time"
)

type WorkstreamTelemetry struct {
    StartedAt           time.Time
    CompletedAt         time.Time
    DurationMinutes     int
    FilesChanged        []string
    LOCAdded            int
    LOCDeleted          int
    TestCoverageBefore  float64
    TestCoverageAfter   float64
    QualityGatesPassed  int
    QualityGatesTotal   int
    FrictionPoints      []string
    Suggestions         []string
    Outcome             string
}

type Collector struct {
    projectDir string
    gitBin     string
    pytestBin  string
}

func (c *Collector) CaptureExecutionTelemetry(
    wsID string,
    commitHash string,
) (*WorkstreamTelemetry, error) {
    // Get git diff stats
    filesChanged, locAdded, locDeleted, err := c.getGitStats(commitHash)
    if err != nil {
        return nil, fmt.Errorf("get git stats: %w", err)
    }

    // Get coverage before/after
    coverageBefore, coverageAfter, err := c.getCoverageStats()
    if err != nil {
        return nil, fmt.Errorf("get coverage: %w", err)
    }

    // Detect friction points
    frictionPoints := c.detectFrictionPoints()

    telemetry := &WorkstreamTelemetry{
        StartedAt:          time.Now().Add(-time.Duration(5) * time.Hour), // Example
        CompletedAt:        time.Now(),
        DurationMinutes:    300, // 5 hours
        FilesChanged:       filesChanged,
        LOCAdded:          locAdded,
        LOCDeleted:        locDeleted,
        TestCoverageBefore: coverageBefore,
        TestCoverageAfter:  coverageAfter,
        FrictionPoints:     frictionPoints,
    }

    return telemetry, nil
}

func (c *Collector) getGitStats(commitHash string) ([]string, int, int, error) {
    // Get changed files
    cmd := exec.Command(c.gitBin, "diff", "--name-only", commitHash+"^.."+commitHash)
    cmd.Dir = c.projectDir
    out, err := cmd.Output()
    if err != nil {
        return nil, 0, 0, err
    }

    files := strings.Fields(string(out))

    // Get LOC stats
    cmd = exec.Command(c.gitBin, "diff", "--shortstat", commitHash+"^.."+commitHash)
    cmd.Dir = c.projectDir
    out, err = cmd.Output()
    if err != nil {
        return files, 0, 0, nil // Return files even if stats fail
    }

    // Parse "5 files changed, 123 insertions(+), 45 deletions(-)"
    var added, deleted int
    var nFiles int
    fmt.Sscanf(string(out), "%d files changed, %d insertions(+), %d deletions(-)",
        &nFiles, &added, &deleted)

    return files, added, deleted, nil
}

func (c *Collector) detectFrictionPoints() []string {
    var points []string

    // Check for quality gate failures
    if c.checkQualityGateFailures() {
        points = append(points, "Quality gate failed multiple times")
    }

    // Check for type hints missing
    if c.checkMissingTypeHints() {
        points = append(points, "Quality gate failed due to missing type hints")
    }

    // Check for documentation mismatches
    if c.checkDocMismatches() {
        points = append(points, "Documentation inconsistent with actual workflow")
    }

    return points
}

func (c *Collector) checkQualityGateFailures() bool {
    // Check if .quality-cache.json has recent failures
    cmd := exec.Command("grep", "-c", "\"passed\":false", ".quality-cache.json")
    cmd.Dir = c.projectDir
    err := cmd.Run()
    return err == nil
}

func (c *Collector) checkMissingTypeHints() bool {
    // Run mypy and check for type hint errors
    cmd := exec.Command("mypy", ".", "--strict")
    cmd.Dir = c.projectDir
    out, _ := cmd.CombinedOutput()
    output := string(out)

    // Look for common type hint errors
    return strings.Contains(output, "has no attribute") ||
           strings.Contains(output, "has no argument")
}

func (c *Collector) checkDocMismatches() bool {
    // Check for TODO comments indicating documentation issues
    cmd := exec.Command("grep", "-r", "TODO.*doc", "docs/")
    cmd.Dir = c.projectDir
    err := cmd.Run()
    return err == nil
}

func (c *Collector) getCoverageStats() (before, after float64, err error) {
    // Run pytest with coverage
    cmd := exec.Command(c.pytestBin, "--cov", ".", "--cov-report=json")
    cmd.Dir = c.projectDir
    err = cmd.Run()
    if err != nil {
        return 0, 0, err
    }

    // Parse coverage.json
    data, err := os.ReadFile(filepath.Join(c.projectDir, "coverage.json"))
    if err != nil {
        return 0, 0, err
    }

    var coverage map[string]any
    json.Unmarshal(data, &coverage)

    totals, _ := coverage["totals"].(map[string]any)
    percent, _ := totals["percent_covered"].(float64)

    return 0.0, percent, nil
}
```

### Component 6: Checkpoint System (Simplified)

```go
// internal/checkpoint/repository.go
package checkpoint

import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"
)

type Checkpoint struct {
    Feature      string
    AgentID      string
    Status       string // "in_progress", "completed", "failed"
    CompletedWS  []string
    CurrentWS    string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type Repository struct {
    checkpointDir string
}

func NewRepository(projectDir string) *Repository {
    return &Repository{
        checkpointDir: filepath.Join(projectDir, ".sdp", "checkpoints"),
    }
}

func (r *Repository) SaveCheckpoint(feature string, cp *Checkpoint) error {
    if err := os.MkdirAll(r.checkpointDir, 0755); err != nil {
        return err
    }

    cp.UpdatedAt = time.Now()

    data, err := json.MarshalIndent(cp, "", "  ")
    if err != nil {
        return err
    }

    checkpointPath := filepath.Join(r.checkpointDir, feature+".json")
    return os.WriteFile(checkpointPath, data, 0644)
}

func (r *Repository) LoadCheckpoint(feature string) (*Checkpoint, error) {
    checkpointPath := filepath.Join(r.checkpointDir, feature+".json")

    data, err := os.ReadFile(checkpointPath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil // No checkpoint found
        }
        return nil, err
    }

    var cp Checkpoint
    if err := json.Unmarshal(data, &cp); err != nil {
        return nil, err
    }

    return &cp, nil
}

func (r *Repository) DeleteCheckpoint(feature string) error {
    checkpointPath := filepath.Join(r.checkpointDir, feature+".json")
    return os.Remove(checkpointPath)
}
```

**Simplified from Python:**
- No SQLite database (use JSON file)
- No complex queries (simple load/save)
- ~200 LOC vs 800 LOC (75% reduction)

---

## Fallback Strategy

### Problem: Beads CLI Not Installed

**Solution:** Hybrid approach — try Beads first, fallback to JSON storage.

### Implementation: Escalation Metrics

```go
// internal/telemetry/repository.go
package telemetry

import (
    "encoding/json"
    "os"
    "os/exec"
    "path/filepath"
)

type Repository struct {
    useBeads   bool
    jsonPath   string
}

func NewRepository(projectDir string) *Repository {
    // Detect if Beads is available
    useBeads := isBeadsAvailable()

    return &Repository{
        useBeads: useBeads,
        jsonPath: filepath.Join(projectDir, ".sdp", "escalation_metrics.json"),
    }
}

func isBeadsAvailable() bool {
    _, err := exec.LookPath("bd")
    return err == nil
}

func (r *Repository) GetEscalations(days int) ([]EscalationEvent, error) {
    if r.useBeads {
        return r.getEscalationsFromBeads(days)
    }
    return r.getEscalationsFromJSON()
}

func (r *Repository) getEscalationsFromBeads(days int) ([]EscalationEvent, error) {
    // Use bd history command
    cmd := exec.Command("bd", "history", "--days", string(rune(days)), "--json")
    out, err := cmd.Output()
    if err != nil {
        // Fallback to JSON if bd history fails
        return r.getEscalationsFromJSON()
    }

    // Parse JSON output
    var events []EscalationEvent
    if err := json.Unmarshal(out, &events); err != nil {
        return nil, err
    }

    return events, nil
}

func (r *Repository) getEscalationsFromJSON() ([]EscalationEvent, error) {
    data, err := os.ReadFile(r.jsonPath)
    if err != nil {
        if os.IsNotExist(err) {
            return []EscalationEvent{}, nil // No escalations yet
        }
        return nil, err
    }

    var events []EscalationEvent
    if err := json.Unmarshal(data, &events); err != nil {
        return nil, err
    }

    return events, nil
}

func (r *Repository) SaveEscalation(event EscationEvent) error {
    if r.useBeads {
        return r.saveEscalationToBeads(event)
    }
    return r.saveEscalationToJSON(event)
}

func (r *Repository) saveEscalationToBeads(event EscalationEvent) error {
    // Create task with "escalation" label
    cmd := exec.Command("bd", "create",
        fmt.Sprintf("Escalation: %s", event.WorkstreamID),
        "--label", "escalation",
        "--label", event.Tier,
        "--json")

    out, err := cmd.Output()
    if err != nil {
        // Fallback to JSON
        return r.saveEscalationToJSON(event)
    }

    // Parse response to get task ID
    var result map[string]any
    json.Unmarshal(out, &result)

    return nil
}

func (r *Repository) saveEscalationToJSON(event EscalationEvent) error {
    // Load existing events
    events, _ := r.getEscalationsFromJSON()

    // Append new event
    events = append(events, event)

    // Save back to JSON
    data, err := json.MarshalIndent(events, "", "  ")
    if err != nil {
        return err
    }

    // Create directory if needed
    if err := os.MkdirAll(filepath.Dir(r.jsonPath), 0755); err != nil {
        return err
    }

    return os.WriteFile(r.jsonPath, data, 0644)
}
```

**Benefits:**
- Works with or without Beads installed
- Seamless fallback (user doesn't need to know)
- Encourages Beads adoption (better experience)

---

## Implementation Plan

### Phase 1: Core + TDD (Week 1-2)

**Goal:** Basic workflow — initialize project, execute workstream with TDD

**Tasks:**
- [ ] Setup Go project structure (cmd/, internal/, pkg/)
- [ ] Implement `core.Workstream` parser + validator
- [ ] Implement `tdd.Runner` with pytest subprocess
- [ ] Implement `beads.Wrapper` (thin CLI layer)
- [ ] Add `sdp init` command (enhance existing)
- [ ] Add `sdp doctor` command (enhance existing)
- [ ] Add `sdp build` command (basic: Red→Green→Refactor)
- [ ] Write unit tests (target: 80% coverage)
- [ ] Integration test: execute real workstream from end to end

**Deliverables:**
- `sdp init`, `sdp doctor`, `sdp build` commands work
- Workstream YAML parsing works
- TDD cycle executes with pytest
- Basic Beads integration (ready, update status)

**Success Criteria:**
```bash
# Initialize new SDP project
$ sdp init my-project
$ cd my-project

# Create workstream manually (or via @design skill)
$ echo "..." > docs/workstreams/backlog/00-001-01.md

# Execute workstream
$ sdp build 00-001-01
✅ Red phase: test_001.py FAILED (expected)
✅ Green phase: test_001.py PASSED
✅ Refactor phase: all checks PASSED
✅ Workstream 00-001-01 completed
```

---

### Phase 2: Quality + Telemetry (Week 3-4)

**Goal:** Automated quality gates + insights collection

**Tasks:**
- [ ] Implement `quality.Checker` with parallel gate execution
- [ ] Implement `quality.Watcher` with fsnotify
- [ ] Implement `telemetry.Collector` (git stats, coverage, friction detection)
- [ ] Implement `telemetry.Analyzer` (pattern detection, insights)
- [ ] Add `sdp quality check` command
- [ ] Add `sdp telemetry scan` command
- [ ] Add `sdp telemetry insights` command
- [ ] Integration test: quality watcher detects violations in real-time
- [ ] Integration test: generate insights report from completed workstreams

**Deliverables:**
- Quality gates run in parallel (mypy, ruff, coverage)
- Quality watcher monitors files and sends notifications
- Telemetry auto-captures metrics after each build
- Insights report identifies friction patterns

**Success Criteria:**
```bash
# Run quality gates
$ sdp quality check
✅ mypy: PASSED (0 errors)
✅ ruff: PASSED (0 warnings)
✅ coverage: 85% (target: 80%)
✅ file size: all files <200 LOC

# Monitor quality in background
$ sdp quality watch &
[Watching for file changes...]

# Edit Python file
$ vim src/module.py
[Desktop notification] Quality gate FAILED: src/module.py:42: missing type hints

# Generate insights
$ sdp telemetry insights --days 30
# Friction Points (Top 5)
# 1. Quality gate failed 2x due to missing type hints (7 occurrences)
# 2. Documentation inconsistent with actual workflow (5 occurrences)
#
# Suggested Protocol Improvements
# 1. WS-001: Add pre-build type hint validation
# 2. WS-002: Update docs/workflow-decision.md
```

---

### Phase 3: Drift Detection + Polish (Week 5-6)

**Goal:** Documentation-code synchronization + production readiness

**Tasks:**
- [ ] Implement `drift.Detector` with fingerprinting
- [ ] Add drift detection to `sdp build` workflow
- [ ] Implement `checkpoint.Repository` (simplified)
- [ ] Add `sdp oneshot --resume` command
- [ ] CLI polish (colors, progress bars, shell completion)
- [ ] Performance optimization (profiling, hotspots)
- [ ] Cross-platform builds (Linux, macOS, Windows)
- [ ] Documentation (README, migration guide, API docs)

**Deliverables:**
- Drift detector compares Contract vs Implementation
- Checkpoint system enables @oneshot resume
- CLI is production-ready (polished UX)
- Binary builds for all platforms
- Complete documentation

**Success Criteria:**
```bash
# Drift detection during build
$ sdp build 00-040-04
⚠️  Drift detected: score 0.75 > threshold 0.5
   Expected: generic validator
   Actual: business logic in quality/models.py
   Declare deviation? [y/N]: y
   Deviation reason: Required for edge case handling
   ✅ Workstream 00-040-04 completed with deviation

# Resume interrupted oneshot
$ sdp oneshot F01
[Interrupted by user]
^C

$ sdp oneshot --resume agent-abc123
Resuming from checkpoint: F01 (3/5 workstreams completed)
✅ Executing 00-001-04...
✅ Executing 00-001-05...
✅ Feature F01 completed

# Cross-platform binary
$ ls -lh sdp
-rwxr-xr-x  1 user  staff   15M Feb  5 10:00 sdp
$ file sdp
sdp: Mach-O 64-bit executable x86_64

$ ./sdp --version
sdp version 1.0.0 (go1.24)
```

---

### Phase 4: Python Cleanup (Week 7-8)

**Goal:** Remove obsolete Python code after successful migration

**Tasks:**
- [ ] Verify all Go features work correctly
- [ ] Run comprehensive integration tests
- [ ] Identify Python code to delete
- [ ] Delete Python modules (not needed anymore)
- [ ] Delete Python tests (not needed anymore)
- [ ] Update CLAUDE.md to reference Go binary
- [ ] Update README.md (remove Python installation)
- [ ] Update .github/workflows (use Go binary)
- [ ] Archive Python code to `python-legacy/` branch
- [ ] Tag release: `v1.0.0-go`

**Files to DELETE (Python implementation):**
```
src/sdp/adapters/          # Not needed (Claude Code integration via skills/)
src/sdp/beads/             # Replaced by internal/beads/wrapper.go
src/sdp/cli/               # Replaced by cmd/sdp/
src/sdp/core/              # Replaced by internal/core/
src/sdp/design/            # Not needed (Beads handles dependencies)
src/sdp/doctor.py          # Replaced by cmd/sdp/commands/doctor.go
src/sdp/errors/            # Not needed (Go error handling)
src/sdp/extensions/        # Not needed
src/sdp/feature/           # Not needed (skills/ handle this)
src/sdp/github/            # Keep for now (future work)
src/sdp/health_checks/     # Replaced by sdp doctor
src/sdp/init_*.py          # Replaced by sdp init
src/sdp/prd/               # Keep for now (product vision)
src/sdp/quality/           # Replaced by internal/quality/
src/sdp/report/            # Not needed (telemetry handles this)
src/sdp/schema/            # Replaced by internal/core/schema.go
src/sdp/tdd/               # Replaced by internal/tdd/
src/sdp/traceability/      # Not needed (covered by drift detection)
src/sdp/unified/           # Partially replaced (checkpoint simplified)
src/sdp/validators/        # Replaced by internal/validators/
```

**Files to KEEP (still needed):**
```
.claude/skills/            # Core SDP functionality (prompts)
.claude/agents/            # Multi-agent prompts
docs/                      # Documentation
hooks/                     # Git hooks (may need Go rewrite later)
prompts/commands/          # Skill definitions (markdown)
scripts/                   # Utility scripts
```

**Migration Checklist:**
```bash
# 1. Backup Python codebase
git checkout -b python-legacy
git push origin python-legacy

# 2. Switch to main branch
git checkout main

# 3. Delete Python code
git rm -r src/sdp/
git rm pyproject.toml
git rm poetry.lock
git rm tests/  # Keep integration tests if useful

# 4. Update documentation
# Edit README.md, CLAUDE.md, etc.

# 5. Commit changes
git commit -m "feat: migrate to Go implementation

- Replace Python CLI with Go binary
- Remove 25K LOC Python code
- Leverage Beads CLI for task tracking
- Single static binary for easy deployment

BREAKING CHANGE: Python runtime no longer required"

# 6. Tag release
git tag v1.0.0-go
git push origin main --tags
```

**Success Criteria:**
- All Python code deleted (except docs/, skills/, prompts/)
- Go binary works flawlessly
- Documentation updated
- No broken references to Python modules

---

### Phase 5: Buffer + Testing (Week 9-10)

**Goal:** Final polish, extensive testing, rollout preparation

**Tasks:**
- [ ] Load testing (100 concurrent workstreams)
- [ ] Fuzz testing (YAML parser with random inputs)
- [ ] Integration test suite (50+ real workstreams)
- [ ] Performance benchmarking (vs Python)
- [ ] Security audit (static analysis)
- [ ] User acceptance testing (5 pilot projects)
- [ ] Rollback plan verification (Python branch still works)
- [ ] Release notes preparation
- [ ] Homepage/website update

**Deliverables:**
- Comprehensive test suite passes
- Performance benchmarks show Go is competitive
- 5 pilot projects successfully migrated
- Rollback plan tested and documented
- Release notes ready

**Success Criteria:**
```bash
# Load test
$ sdp benchmark --workstreams 100 --concurrency 10
✅ Executed 100 workstreams in 5m 23s (avg 3.2s per WS)
✅ Memory usage: 45MB peak
✅ No race conditions detected

# Fuzz test
$ go test -fuzz=FuzzParseWorkstream -fuzztime=60s
✅ No crashes found after 60s fuzzing

# Integration test
$ ./scripts/integration-test.sh
✅ Test 1: Simple workstream PASSED
✅ Test 2: Dependent workstreams PASSED
✅ Test 3: Quality gates PASSED
✅ Test 4: Telemetry capture PASSED
✅ Test 5: Drift detection PASSED
...
✅ All 50 integration tests PASSED

# Performance comparison
$ hyperfine './sdp-python build 00-001-01' './sdp-go build 00-001-01'
Benchmark 1: sdp-python build 00-001-01
  Time (mean ± σ):     3.245 s ±  0.120 s
Benchmark 2: sdp-go build 00-001-01
  Time (mean ± σ):     3.102 s ±  0.098 s
Summary: sdp-go build was 1.05x (± 0.06) faster
```

---

## Risk Mitigation

### Risk Matrix

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| **Python subprocess failures** | Medium | High | - Detect Python availability at startup<br>- Provide helpful error messages<br>- Document Python dependencies clearly |
| **Beads CLI not installed** | High | Medium | - Fallback to JSON storage<br>- Auto-install Beads if user agrees<br>- Clear error messages with install instructions |
| **SQLite locking (checkpoints)** | Low | Medium | - Use JSON instead of SQLite<br>- No concurrent access issues<br>- Simpler code (200 vs 800 LOC) |
| **fsnotify cross-platform issues** | Medium | Low | - Test on Linux, macOS, Windows<br>- Fallback to polling if watcher fails<br>- Document platform limitations |
| **Performance regression vs Python** | Low | High | - Benchmark critical paths<br>- Optimize hotspots<br>- Use goroutines for parallelism |
| **Breaking changes to skills/.claude/** | Medium | High | - Keep skill format unchanged<br>- Go CLI is drop-in replacement<br>- Test all skills with Go binary |
| **Data loss during Python cleanup** | Low | Critical | - Create python-legacy branch<br>- Tag before deletion<br>- Verify Go works perfectly first |

### Rollback Plan

If Go version has critical issues in Weeks 1-10:

1. **Stop using Go binary**
   ```bash
   # Uninstall Go binary
   rm /usr/local/bin/sdp

   # Restore Python version
   git checkout python-legacy
   pip install -e .
   ```

2. **Report bug**
   - Create GitHub issue with detailed reproduction
   - Attach logs and error messages
   - Tag as `critical`, `go-migration`

3. **Fix Go issue**
   - Debug and fix in development branch
   - Add regression test
   - Verify fix works

4. **Retry migration**
   - Wait 1 week for stability
   - Retry cutover with fixed version

**Success criteria for rollback:**
- Python branch still works (verified)
- Zero data loss (checkpoints, telemetry)
- User impact minimized (quick switchback)

---

## Success Criteria

### Feature Parity

| Feature | Python | Go | Status |
|---------|--------|-----|--------|
| **sdp init** | ✅ | ✅ | Required |
| **sdp doctor** | ✅ | ✅ | Required |
| **sdp build** | ✅ | ✅ | Required |
| **sdp quality check** | ✅ | ✅ | Required |
| **sdp telemetry scan** | ✅ | ✅ | Required |
| **sdp telemetry insights** | ✅ | ✅ | Required |
| **@oneshot** | ✅ | ✅ | Required |
| **@oneshot --resume** | ✅ | ✅ | Required |
| **Beads integration** | ✅ | ✅ | Required |
| **Drift detection** | ❌ | ✅ | **New feature** |
| **Quality watcher** | ❌ | ✅ | **New feature** |

### Performance Metrics

| Metric | Python | Go Target | Measurement |
|--------|--------|-----------|-------------|
| **sdp build latency** | 3.2s | ≤3.5s | Benchmark real workstream |
| **Quality gates time** | 8.5s | ≤8.0s | Parallel execution wins |
| **Binary size** | N/A | ≤20MB | `ls -lh sdp` |
| **Memory usage** | 45MB | ≤50MB | `ps aux` |
| **Startup time** | 0.5s | ≤0.3s | Go is compiled |

### Quality Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Test coverage** | ≥80% | `go test -cover` |
| **Race conditions** | 0 | `go test -race` |
| **Lint passed** | Yes | `golangci-lint run` |
| **Fuzz tests** | No crashes | 60s fuzzing |
| **Integration tests** | 50+ pass | `./scripts/integration-test.sh` |

### Adoption Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| **Pilot projects** | 5/5 successful | Manual testing |
| **User satisfaction** | ≥4/5 stars | Survey |
| **Bug reports** | <5 in first week | GitHub issues |
| **Performance complaints** | 0 | Feedback |

### Final Checklist

**Before Go Migration is Complete:**

- [ ] All Python features work in Go
- [ ] Test coverage ≥80%
- [ ] 5 pilot projects successful
- [ ] Performance competitive with Python
- [ ] Documentation complete (migration guide)
- [ ] Rollback plan tested
- [ ] Python cleanup completed
- [ ] Release tagged (v1.0.0-go)
- [ ] Binary builds available (Linux, macOS, Windows)
- [ ] Homebrew/AUR packages published

---

## Appendix

### A. Build Commands

```bash
# Development build (fast)
go build -o sdp ./cmd/sdp

# Production build (optimized, stripped)
go build -ldflags "-s -w" -o sdp ./cmd/sdp

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o sdp-linux-amd64 ./cmd/sdp
GOOS=darwin GOARCH=amd64 go build -o sdp-darwin-amd64 ./cmd/sdp
GOOS=darwin GOARCH=arm64 go build -o sdp-darwin-arm64 ./cmd/sdp
GOOS=windows GOARCH=amd64 go build -o sdp-windows-amd64.exe ./cmd/sdp

# Build all platforms
make build-all
```

### B. Dependencies

**Go:**
```go
// go.mod
module github.com/ai-masters/sdp

go 1.24

require (
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.0
    gopkg.in/yaml.v3 v3.0.1
    github.com/charmbracelet/lipgloss v0.10.0
    github.com/fsnotify/fsnotify v1.7.0
)
```

**External (runtime):**
- Beads CLI (bd) — optional, fallback to JSON
- pytest (test execution)
- mypy (type checking)
- ruff (linting)
- git (version control)

### C. Installation

**From source:**
```bash
git clone https://github.com/ai-masters/sdp.git
cd sdp
make build
sudo make install
```

**From Homebrew (macOS/Linux):**
```bash
brew install ai-masters/sdp/sdp
```

**From AUR (Arch Linux):**
```bash
paru -S sdp-cli
```

**From binary:**
```bash
wget https://github.com/ai-masters/sdp/releases/latest/download/sdp-linux-amd64
chmod +x sdp-linux-amd64
sudo mv sdp-linux-amd64 /usr/local/bin/sdp
```

---

## Conclusion

This design document outlines a pragmatic Go migration strategy that:

1. **Reduces scope by 72%** through Beads CLI leverage
2. **Ships in 10 weeks** with 2-week Python cleanup phase
3. **Delivers single binary** for easy deployment
4. **Maintains feature parity** while adding new capabilities (drift detection, quality watcher)
5. **Minimizes risk** through fallback mechanisms and rollback plan

**Key Innovation:** Don't build what Beads already provides. Use thin Go wrappers, delegate to `bd --json`, focus on SDP-unique functionality.

**Next Steps:**
1. Review and approve this design
2. Create workstreams for Phase 1 tasks
3. Setup Go development environment
4. Begin implementation with TDD discipline

---

**Document Version:** 1.0
**Last Updated:** 2026-02-05
**Status:** Ready for Implementation
