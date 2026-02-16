# /build — Execute Workstream

When calling `/build {WS-ID}`:

1. Load skill: `@.claude/skills/build/SKILL.md`
2. Run pre-build hook: `hooks/pre-build.sh {WS-ID}`
3. Read WS plan
4. Execute steps using TDD
5. Run post-build hook: `hooks/post-build.sh {WS-ID}`
6. Append Execution Report to WS file

## Quick Reference

**Input:** `workstreams/backlog/WS-XXX-*.md`
**Output:** Code + tests + Execution Report
**Next:** `/build WS-XXX-02` or `/review F{XX}`
