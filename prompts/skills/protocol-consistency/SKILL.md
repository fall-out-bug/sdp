---
name: protocol-consistency
description: Audit consistency across workstream docs, CLI capabilities, and CI workflows.
---

# Protocol Consistency

Run this skill when you suspect process drift between documentation, CLI commands, and automation workflows.

## Workflow

### 1) Verify declared commands vs available CLI

```bash
sdp --help
sdp <command> --help
```

Check that commands referenced in docs/workstreams/hooks actually exist.

### 2) Validate workstream schema compatibility

For target workstreams:

```bash
sdp parse ws <ws-id>
sdp drift detect <ws-id>
```

Identify schema mismatches (e.g. `feature` vs `feature_id`).

### 3) Validate CI/workflow command paths

```bash
rg -n "sdp .*" .github/workflows hooks scripts -S
```

Confirm every referenced command is valid in current CLI and has expected flags.

### 4) Report mismatches

For each mismatch, include:

- Source file + line
- Observed behavior
- Expected behavior
- Risk (blocking/non-blocking)
- Suggested minimal fix

### 5) Track in Beads

For blocking or repeat issues:

```bash
bd create --title="Protocol drift: <summary>" --type=task --priority=2
bd sync
```

## Output Template

```markdown
## Protocol Consistency Report

- Scope: ...
- Blocking mismatches: N
- Non-blocking mismatches: N

### Findings
1. ...
2. ...

### Recommended fixes
1. ...
2. ...
```
