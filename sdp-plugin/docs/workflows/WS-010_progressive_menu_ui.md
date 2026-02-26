# WS-010: Progressive Menu UI

> **Workstream ID:** WS-010  
> **Feature:** Unified Workflow  
> **Status:** Ready to Spec  
> **Dependencies:** WS-009

## Goal

Implement progressive menu system for @feature skill that allows users to:
- Skip phases with flags (--vision-only, --no-interview)
- Start from existing specifications (--spec PATH)
- Track progress with real-time updates
- Make interactive choices during execution

## Acceptance Criteria

### AC1: Phase Skipping
- [ ] --vision-only flag stops after Phase 2 (PRODUCT_VISION.md)
- [ ] --no-interview skips AskUserQuestion calls
- [ ] Flags validated before execution starts

### AC2: Existing Spec Import
- [ ] --spec PATH loads existing draft from docs/drafts/
- [ ] Skips vision and requirements phases
- [ ] Validates spec format before proceeding

### AC3: Progress Display
- [ ] Real-time updates: "[HH:MM] Executing WS-XXX..."
- [ ] Shows current phase (Vision → Requirements → Planning → Execution)
- [ ] Displays checkpoints reached

### AC4: Menu Logging
- [ ] User choices logged via `sdp decisions log`
- [ ] Flags and options recorded for reproducibility

## Scope Files

### Files to Modify
**prompts/skills/feature.md** - Add flag handling and menu logic

**cmd/sdp/feature.go** (NEW if doesn't exist) - CLI entry point
```go
func ExecuteFeature(args []string) error {
    // Parse flags
    visionOnly := flag.Bool("vision-only", false, "")
    noInterview := flag.Bool("no-interview", false, "")
    specPath := flag.String("spec", "", "")
    
    // Execute based on flags
}
```

## Implementation Steps

### Step 1: Add Flag Parsing to @feature
```markdown
## Power User Flags

- `--vision-only` -- Only create vision, skip planning
- `--no-interview` -- Skip questions, use defaults
- `--update-vision` -- Update existing PRODUCT_VISION.md
- `--spec PATH` -- Start from existing spec

## Progressive Menu

If user wants interactive mode:

1. **Phase Selection:**
   - Full workflow (vision → requirements → planning → execution)
   - Vision only (skip to design)
   - From existing spec (skip to execution)

2. **Progress Display:**
   ```
   [15:23] Phase 1: Vision Interview...
   [15:45] Phase 2: Generating PRODUCT_VISION.md...
   [15:50] Phase 3: Technical Interview...
   [16:20] Phase 4: Generating intent.json...
   [16:30] Phase 5: Creating requirements draft...
   [16:45] Phase 6: Calling @design...
   [17:00] Phase 7: Orchestrator executing...
   [17:05] → Executing WS-009...
   [17:27] → WS-009 complete (22m)
   ```

3. **Checkpoint Display:**
   ```
   📊 Progress: 1/26 workstreams (4%)
   ⏱️  Elapsed: 1h 23m
   💾 Last checkpoint: 5m ago
   ```
```

### Step 2: Implement Flag Logic
```go
func handleFeatureFlags(args []string) error {
    visionOnly := flag.Bool("vision-only", false, "")
    specPath := flag.String("spec", "", "")
    
    if *specPath != "" {
        // Load existing spec
        return executeFromSpec(*specPath)
    }
    
    if *visionOnly {
        // Run only vision phases
        return runVisionPhases()
    }
    
    // Full workflow
    return runFullWorkflow()
}
```

### Step 3: Add Progress Tracking
```go
type ProgressTracker struct {
    currentPhase string
    workstreamsCompleted int
    workstreamsTotal int
    startTime time.Time
}

func (pt *ProgressTracker) Display() {
    elapsed := time.Since(pt.startTime)
    
    fmt.Printf("📊 Phase: %s\n", pt.currentPhase)
    fmt.Printf("⏱️  Elapsed: %v\n", elapsed.Round(time.Minute))
    
    if pt.workstreamsTotal > 0 {
        pct := float64(pt.workstreamsCompleted) / float64(pt.workstreamsTotal) * 100
        fmt.Printf("📊 Progress: %d/%d workstreams (%.1f%%)\n",
            pt.workstreamsCompleted, pt.workstreamsTotal, pct)
    }
}
```

## Definition of Done

- [ ] All 4 acceptance criteria met
- [ ] Flag parsing implemented
- [ ] Progress display working
- [ ] Menu choices logged
- [ ] Tests ≥80% coverage
- [ ] Documentation updated

## Estimated Scope

- ~200 LOC implementation
- ~100 LOC tests
- Duration: 2 hours
- Size: SMALL

## Success Metrics

- Users can skip phases with flags
- Progress visible in real-time
- Existing specs load correctly
- Menu choices reproducible via decision log
