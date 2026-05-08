# SDP OpenCode Integration

This directory contains configuration for integrating SDP with OpenCode (OhMyOpenCode).

## OpenCode-First Approach

SDP is now **OpenCode-first**. Configuration lives in `.opencode/` directory:
- `.opencode/hooks/` - Hook configurations (this directory)
- `.opencode/agents` - Symlink to `prompts/agents`
- `.opencode/commands` - Symlink to `prompts/commands`

Legacy Claude config (`.claude/`) is maintained for compatibility but OpenCode is the primary environment.

## sdp-omc-guard

Pre-tool-call guard that enforces scope before edit/write operations.

### Usage

```bash
# Build
go build ./cmd/sdp-omc-guard

# Run manually (for testing)
echo '{"tool_name":"edit","tool_input":{"file_path":"test.go"},"cwd":"."}' | \
  sdp-omc-guard --ws 00-059-01 --session-id test-001

# Exit codes:
#   0 = allow (files in scope)
#   1 = ask (not used)
#   2 = deny (files out of scope)
```

### OpenCode Hook Configuration

Add to `~/.config/opencode/opencode.json` or use the pre-configured `.opencode/hooks/pre-tool-use.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "edit*",
        "hooks": [
          {
            "type": "command",
            "command": "sdp-omc-guard --ws ${SDP_WORKSTREAM} --emit-evidence"
          }
        ]
      }
    ]
  }
}
```

### Environment Variables

- `SDP_WORKSTREAM`: Workstream ID (e.g., `00-059-01`). Required for scope checking.
- `SDP_SESSION_ID`: Session ID for evidence logging.

## sdp-ready

Ready queue bridge for Beads issue tracking.

```bash
# Build
go build ./cmd/sdp-ready

# Run
sdp ready              # Text output
sdp ready --format json  # JSON output
sdp ready --no-cache   # Skip cache
```
