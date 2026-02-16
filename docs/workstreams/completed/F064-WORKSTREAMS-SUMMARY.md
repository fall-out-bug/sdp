# F064: Unified Task Resolver — Workstream Summary

> Feature: sdp-t01g | Priority: P0 | 4 workstreams

## Goal

Implement unified task resolver and registration system in Go CLI, enabling:
- ID resolution (workstream, beads, issue)
- @review artifact creation
- /issue skill backend
- Full fallback without beads

## Design Reference

- [Integration Design](../plans/2026-02-12-review-issue-execution-integration-design.md)

## Dependency Graph

```
00-064-01 ID Resolver Core
    ├─► 00-064-02 Task Registration (review/issue)
    └─► 00-064-03 Beads Integration

00-064-04 Skill Updates (depends on 01, 02, 03)
```

## Workstream Registry

| WS ID | Title | Est. | Depends On |
|-------|-------|------|------------|
| 00-064-01 | ID Resolver Core | 2h | - |
| 00-064-02 | Task Registration API | 2h | 00-064-01 |
| 00-064-03 | Beads Integration | 1.5h | 00-064-01 |
| 00-064-04 | Skill Updates | 1h | 00-064-02, 00-064-03 |

## Acceptance Criteria (Feature-level)

- [ ] AC1: `sdp resolve <id>` returns task info (WS/beads/issue)
- [ ] AC2: `sdp task create --type=bug|task` creates workstream + beads
- [ ] AC3: All skills work with any ID format
- [ ] AC4: Full functionality without beads installed
- [ ] AC5: Bidirectional linking maintained

## Key Files

```
sdp-plugin/
├── internal/
│   ├── resolver/
│   │   ├── resolver.go       # ID detection and resolution
│   │   ├── resolver_test.go
│   │   ├── workstream.go     # WS file parsing
│   │   ├── beads.go          # Beads integration
│   │   └── issues.go         # Issue file handling
│   └── task/
│       ├── creator.go        # Task creation logic
│       ├── creator_test.go
│       └── types.go          # Task types
├── cmd/sdp/
│   ├── resolve.go            # sdp resolve command
│   └── task.go               # sdp task command
└── docs/
    └── issues/               # Fallback issue storage
```
