# Minimal Go Project Template

This template demonstrates a minimal SDP-compatible Go project structure.

## Quick Start

```bash
# 1. Initialize SDP
sdp init --guided

# 2. Verify environment
sdp doctor

# 3. Plan a feature
sdp plan "Add hello world function"

# 4. Execute workstream
sdp apply --ws 00-001-01

# 5. Check status
sdp status --text
```

## Project Structure

```
minimal-go/
├── .claude/           # SDP prompts (created by sdp init)
├── .sdp/              # SDP config (created by sdp init)
├── docs/
│   └── workstreams/   # Workstream tracking
├── go.mod             # Go module definition
├── main.go            # Entry point
└── README.md          # This file
```

## Requirements

- Go 1.21 or later
- Git
- SDP CLI

## What This Template Demonstrates

1. **Minimal Structure**: Just enough to be a valid SDP project
2. **Go Project Detection**: Has go.mod for auto-detection
3. **Git Ready**: Initialize with `git init` before `sdp init`
4. **Quality Gates**: Go test integration

## Next Steps

1. Copy this template: `cp -r templates/minimal-go my-project`
2. Initialize: `cd my-project && git init && sdp init --guided`
3. Start building: `sdp plan "your first feature"`
