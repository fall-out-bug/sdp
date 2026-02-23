---
name: oneshot
description: Autonomous multi-agent execution with review-fix loop, PR creation, CI-fix loop, provenance, and drift detection
cli: sdp orchestrate (file ops only - requires @build for actual work)
version: 8.1.0
changes:
  - Step 7: PENDING vs FAILURE; CI loop mandatory; no handoff lists
  - Step 8: Completion output — only "CI GREEN", no delegation
  - CRITICAL RULES 7–8
  - Step 0a: Feature Context Loading from ROADMAP
  - Step 0b: Branch Setup with checkpoint and run file
  - Step 1.5: Pre-Build Drift Gate (sdp drift detect)
  - Step 4: Two-Phase Review-Fix Loop with stuck detection
  - Step 7: CI Check-Fix Loop
  - Provenance: evidence.trace.pr_url, run file, decision log
---

# oneshot

> **CLI:** `sdp orchestrate <feature-id>` — handles file loading, dependency graph, checkpoints
> **LLM:** Required for actual workstream execution via `@build`

Autonomous feature execution: Feature Context → Branch → Drift Gate → Build → Evidence → Review Loop (0 findings) → PR → CI Loop (green) → Done.

---

## CRITICAL RULES

1. **NEVER STOP** - Execute ALL workstreams in one session. No pauses between WS.
2. **NO SUMMARIES** - Only commit messages. No "progress reports" or "session summaries".
3. **AUTO-CONTINUE** - After WS commit, IMMEDIATELY start next WS without asking.
4. **ONLY STOP IF:** All WS done OR unrecoverable blocker OR user explicitly stops you.
5. **POST-COMPACTION RECOVERY** - After context compaction, read checkpoint first. Never drift to side tasks.
6. **PROVENANCE** - Populate evidence files, run file events, and decision log. Never skip artifact writes.
7. **CI LOOP MANDATORY** - Step 7: poll until green. If PENDING → wait, retry. Never hand off with "wait for CI yourself".
8. **NO HANDOFF LISTS** - When done, output only "CI GREEN - @oneshot complete". Do NOT output "Next steps", "Optional: run /review", or delegation lists. Human UAT and merge are implicit — no handoff.

---

## POST-COMPACTION PROTOCOL

**If session was compacted, you MUST check first:**

```bash
# 1. Check checkpoint (primary source of truth)
CHECKPOINT=$(ls .sdp/checkpoints/F*.json 2>/dev/null | head -1)
if [ -n "$CHECKPOINT" ]; then
  echo "=== RESUMING FROM CHECKPOINT ==="
  cat "$CHECKPOINT"
  # Find first WS with status != "done" → continue from there
  # Restore branch: git checkout $(jq -r .branch "$CHECKPOINT")
fi

# 2. If no checkpoint, check beads
bd list --status=in_progress
bd ready

# 3. Resume PRIMARY TASK, not side task
# Side task: fixing tests, improving coverage, debugging
# Primary: executing roadmap, implementing feature
```

**The summary mentions "side task" → IGNORE IT, return to PRIMARY.**

---

## EXECUTE THIS NOW

When user invokes `@oneshot F067` (replace F067 with actual feature ID):

### Step 0a: Load Feature Context

```bash
# Parse feature number
FNUM=$(echo "F067" | sed 's/F0*//')
WS_PATTERN="00-$(printf '%03d' $FNUM)-"

# Verify feature exists in ROADMAP (fail fast)
if ! grep -q "F${FNUM}" docs/roadmap/ROADMAP.md; then
  echo "Feature F${FNUM} not found in ROADMAP.md"
  exit 1
fi

# Extract: feature_name, phase, exit_criteria, depends_on from ROADMAP
# Check feature dependencies via WS frontmatter statuses
# Check WS files exist
ws_files=$(ls docs/workstreams/backlog/${WS_PATTERN}*.md 2>/dev/null)
if [ -z "$ws_files" ]; then
  echo "No workstream files found for F${FNUM}"; exit 1
fi

# Decision log
sdp decisions log --feature-id F${FNUM} --type explicit \
  --question "Execute feature?" \
  --decision "F${FNUM}: {feature_title}" \
  --rationale "ROADMAP: Phase N, deps OK" \
  --maker agent
```

Display feature summary (goal, exit criteria, workstream list) before proceeding.

### Step 0b: Branch Setup

```bash
# Derive branch name from ROADMAP feature title
FEATURE_TITLE=$(grep "F${FNUM}" docs/roadmap/ROADMAP.md | \
  sed 's/.*F[0-9]*[[:space:]]*//' | cut -d'|' -f1 | \
  tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | \
  sed 's/--*/-/g' | sed 's/^-//' | sed 's/-$//' | cut -c1-40)
BRANCH="feature/F${FNUM}-${FEATURE_TITLE}"
# Fallback if empty: BRANCH="feature/F${FNUM}"

# Verify clean state
if [ -n "$(git status --porcelain)" ]; then
  echo "ERROR: Uncommitted changes. Stash or commit first."; exit 1
fi

# Idempotent branch setup
CURRENT=$(git branch --show-current)
if [ "$CURRENT" = "$BRANCH" ]; then
  echo "Already on $BRANCH (resume mode)"
elif git show-ref --verify --quiet "refs/heads/$BRANCH"; then
  git checkout "$BRANCH"
else
  git fetch origin && git checkout master && git pull
  git checkout -b "$BRANCH"
fi

# Create checkpoint
mkdir -p .sdp/checkpoints .sdp/ws-verdicts .sdp/runs
RUN_ID="oneshot-F${FNUM}-$(date -u +%Y%m%dT%H%M%SZ)"
cat > .sdp/checkpoints/F${FNUM}.json << EOF
{
  "schema": "1.0",
  "feature_id": "F${FNUM}",
  "branch": "$BRANCH",
  "created_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "updated_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "phase": "build",
  "workstreams": [],
  "review": {"iteration": 0, "verdict_file": ".sdp/review_verdict.json", "status": "pending"},
  "pr_number": null,
  "pr_url": null
}
EOF

# Create run file
cat > .sdp/runs/${RUN_ID}.json << EOF
{
  "run_id": "${RUN_ID}",
  "feature_id": "F${FNUM}",
  "orchestrator": "cursor-oneshot",
  "branch": "$BRANCH",
  "started_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "events": [{"at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)", "phase": "init", "state": "ok"}],
  "last_phase": "init",
  "last_state": "ok"
}
EOF
```

### Step 1: Load Workstreams

```bash
ls docs/workstreams/backlog/00-067-*.md
```

Read each file for: WS ID, `depends_on`, AC, scope files. Build list of WS IDs in dependency order.

### Step 1.5: Pre-Build Drift Gate

```bash
for ws_file in $ws_files; do
  ws_id=$(basename "$ws_file" .md)
  result=$(sdp drift detect "$ws_id" 2>&1)
  exit_code=$?
  if [ $exit_code -ne 0 ]; then
    echo "DRIFT ERROR for $ws_id: $result"
    echo "Action: Update WS scope files OR create missing files first."
    exit 1
  elif echo "$result" | grep -q "WARNING"; then
    echo "DRIFT WARNING for $ws_id: $result"
    # Proceed - entity might be new
  else
    echo "$ws_id: scope verified"
  fi
done
```

### Step 2: Build Dependency Graph

Topological sort: empty `depends_on` first, then dependents. Log decision:

```bash
sdp decisions log --feature-id F${FNUM} --type explicit \
  --question "WS execution order?" \
  --decision "Wave: {order}" \
  --rationale "Topological sort of depends_on" \
  --maker agent
```

### Step 3: Execute Workstreams

For each WS in dependency order:

1. Update checkpoint: set `ws.status = in_progress` for current WS
2. Append run file event: `{"phase": "ws:{ws-id}", "state": "running"}`
3. **Invoke @build** with workstream ID (e.g. `@build 00-067-01`)
4. Verify `.sdp/ws-verdicts/{ws-id}.json` exists and `verdict == "PASS"`, `ac_evidence` filled
5. Post-build drift: `sdp drift detect {ws-id}` — if ERROR, treat as @build failure, retry up to 2 times
6. Update checkpoint: `ws.status = done`, `commit = ...`
7. Append run file event: `{"phase": "ws:{ws-id}", "state": "ok", "commit": "..."}`
8. `bd update {beads_id} --status completed`

**Handle failures:** Retry up to 2 times, then escalate (create beads issue, HALT).

### Step 4: Review-Fix Loop

**PHASE 1** (max 5 iterations, max 2 stalled):

1. Run `@review F067`
2. Read `.sdp/review_verdict.json`
3. If `verdict == "APPROVED"` → patch `evidence.review.adversarial_review` for each WS evidence file → break
4. **Stuck detection:** If `len(blocking_ids)` >= previous count for 2 iterations → HALT, escalate
5. If `len(blocking_ids) == 0` (all P2/P3) → treat as APPROVED, break
6. For each P0 finding: fix inline, `git commit -m "fix(F067): {title}"`, `bd close {id}`
7. For each P1 finding: invoke `@bugfix {id}` (stay in feature branch), `bd close {id}`
8. `sdp decisions log` for fix strategy (per finding)
9. Repeat

**PHASE 2:** Drain P2/P3 to beads as tech debt: `bd update {id} --status=backlog --notes="Tech debt from F067 review"`

### Step 5: Verify Clean State

```bash
OPEN_BLOCKING=$(bd list --label review-finding --label F067 --status open --json 2>/dev/null | jq '[.[] | select(.priority <= 1)] | length')
if [ "$OPEN_BLOCKING" -ne 0 ]; then
  echo "$OPEN_BLOCKING blocking findings remain"; exit 1
fi
# Run quality gates (see Quality Gates in AGENTS.md)
```

### Step 6: Create PR

```bash
git push origin feature/F067-xxx
gh pr create --base master --head feature/F067-xxx --title "feat(F067): {title}" --body "..."
PR_URL=$(gh pr view --json url -q '.url')
PR_NUMBER=$(gh pr list --head $(git branch --show-current) --json number -q '.[0].number')

# Patch evidence.trace.pr_url for all feature WS evidence files
for ev in .sdp/evidence/sdp_dev-*.json; do
  [ -f "$ev" ] && jq --arg u "$PR_URL" '.trace.pr_url = $u' "$ev" > "$ev.tmp" && mv "$ev.tmp" "$ev"
done

# Update run file
# Append event: {"phase": "pr", "state": "ok", "pr_url": "...", "pr_number": N}
```

### Step 7: CI Check-Fix Loop

**RULE:** Do NOT hand off. Poll until green or escalate. Never say "wait for CI yourself" or "next steps: wait for CI".

```bash
CI_ITER=0
CI_MAX_ITER=5
sleep 90  # CI boot

while [ $CI_ITER -lt $CI_MAX_ITER ]; do
  PENDING=$(gh pr checks $PR_NUMBER --json name,state -q '.[] | select(.state == "PENDING" or .state == "IN_PROGRESS") | .name' 2>/dev/null)
  FAILING=$(gh pr checks $PR_NUMBER --json name,state -q '.[] | select(.state == "FAILURE" or .state == "ERROR") | .name' 2>/dev/null)

  if [ -n "$PENDING" ]; then
    echo "CI checks still running: $PENDING"; sleep 60; continue
  fi

  if [ -z "$FAILING" ]; then
    bd list --label ci-finding --label F067 --status open --json 2>/dev/null | jq -r '.[].id' | while read id; do bd update "$id" --status=closed --notes="CI green"; done
    echo "CI GREEN - @oneshot complete"; break
  fi

  RUN_ID=$(gh run list --branch $(git branch --show-current) --json databaseId,conclusion --jq '.[] | select(.conclusion == "failure") | .databaseId' 2>/dev/null | head -1)
  gh run view $RUN_ID --log-failed 2>/dev/null > /tmp/ci-failure.log

  # Classify: Go compile/test, k8s-validate = AUTO-FIX; secrets, flaky, out-of-scope = ESCALATE
  # If auto-fixable: patch, commit, push, CI_ITER++, sleep 90, continue
  # If not: bd create --title="CI BLOCKED: ..." --priority=0 --labels "ci-finding,F067"
  #         sdp decisions log --decision "ESCALATE" --rationale "..."
  #         HALT
  CI_ITER=$((CI_ITER + 1))
done
```

### Step 8: Completion Output

**When done:** Output only `CI GREEN - @oneshot complete` and PR URL. Do NOT output:
- "Next steps"
- "Optional: run /review"
- "Human UAT → approve and merge"
- Any delegation list

---

## Checkpoint Schema

`.sdp/checkpoints/F067.json`:

```json
{
  "schema": "1.0",
  "feature_id": "F067",
  "branch": "feature/F067-my-feature",
  "created_at": "...",
  "updated_at": "...",
  "phase": "build",
  "workstreams": [
    {"id": "00-067-01", "status": "done", "verdict_file": ".sdp/ws-verdicts/00-067-01.json", "commit": "abc123", "attempts": 1}
  ],
  "review": {"iteration": 0, "verdict_file": ".sdp/review_verdict.json", "status": "pending"},
  "pr_number": null,
  "pr_url": null
}
```

---

## Run File Schema

`.sdp/runs/oneshot-F067-{ts}.json`:

```json
{
  "run_id": "oneshot-F067-20260223T120000Z",
  "feature_id": "F067",
  "orchestrator": "cursor-oneshot",
  "branch": "feature/F067-xxx",
  "started_at": "...",
  "events": [
    {"at": "...", "phase": "init", "state": "ok"},
    {"at": "...", "phase": "drift:pre:00-067-01", "state": "ok"},
    {"at": "...", "phase": "ws:00-067-01", "state": "running"},
    {"at": "...", "phase": "ws:00-067-01", "state": "ok", "commit": "abc123"},
    {"at": "...", "phase": "pr", "state": "ok", "pr_url": "...", "pr_number": 42},
    {"at": "...", "phase": "ci", "state": "ok"}
  ],
  "last_phase": "ci",
  "last_state": "ok"
}
```

---

## Finding Priority

| Priority | Action | Blocks? |
|----------|--------|---------|
| P0 | Fix immediately (inline) | YES |
| P1 | Create bugfix via @bugfix | YES |
| P2+ | Track only (drain to beads as tech debt) | NO |

---

## Resume After Interruption

```bash
cat .sdp/checkpoints/F067.json
# Find first workstream with status != "done"
# git checkout $(jq -r .branch .sdp/checkpoints/F067.json)
# Continue from that WS
```

---

## See Also

- `@build` - Execute single workstream (REQUIRED)
- `@review` - Quality review
- `@ci-triage` - CI failure classification (used in Step 7)
- `@verify-workstream` - Drift resolution when HALT on ERROR
- `docs/plans/2026-02-23-oneshot-autonomous-design.md` - Full design
