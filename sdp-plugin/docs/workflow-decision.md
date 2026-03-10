# Workflow Decision Guide

Use this guide when choosing between direct CLI execution and Beads-backed task tracking.

## Beads-first workflow

Choose this path when work spans multiple sessions, contributors, or explicit dependencies.

- `sdp beads ready` - inspect available tasks
- `sdp beads show <id>` - inspect task detail
- `sdp beads update <id>` - update task state
- `sdp beads sync` - persist tracker state to the repo
- `sdp build <ws-id>` or `@build 00-001-01` - execute the claimed workstream

Best for:

- multi-session work
- team handoffs
- roadmap execution with dependencies
- queues that need explicit claiming and sync

## Direct workflow

Choose this path when you already know the target workstream or need a tight local loop.

- `sdp next` - get the next recommended action from the current repo state
- `sdp parse <ws-id>` - inspect a workstream
- `sdp build <ws-id>` - execute one workstream
- `sdp verify <ws-id>` - validate completion and evidence
- `sdp log show` - inspect emitted evidence events

Best for:

- short focused fixes
- local verification
- rapid prototyping
- single-developer flows without tracker overhead

## How they fit together

- Beads tracks what should happen next.
- `sdp build`, `sdp verify`, `sdp log`, and related CLI commands execute and validate the work.
- Skills like `@feature`, `@design`, `@build`, and `@review` wrap the same protocol concepts in agent-friendly workflows.
