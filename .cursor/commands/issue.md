# /issue — Analyze & Route Issues

When calling `/issue "description"`:

1. Load full prompt: `@prompts/commands/issue.md`
2. Systematic debugging (5 phases):
   - Symptom analysis
   - Hypothesis formation
   - Systematic elimination
   - Root cause isolation
   - Impact chain analysis
3. Classify severity (P0/P1/P2/P3)
4. Route to appropriate fix
5. Create issue file
6. Create GitHub issue (if gh available)

## Quick Reference

**Input:** Bug description
**Output:** Issue file + Routing recommendation

**Routing:**
- P0 CRITICAL → `/hotfix`
- P1 HIGH → `/bugfix`
- P2 MEDIUM → New WS
- P3 LOW → Defer

**Next:** `/hotfix` or `/bugfix` depending on severity
