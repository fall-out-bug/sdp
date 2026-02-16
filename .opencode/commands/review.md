# /review — Review Feature

When calling `/review {feature}`:

1. Load full prompt: `@.claude/skills/review.md`
2. Find all feature WS in INDEX.md
3. Check each WS against checklist (Check 0-11)
4. Perform cross-WS checks
5. Append Review Results to each WS file
6. Output Feature Summary

## Quick Reference

**Input:** All feature WS
**Output:** Review Results in each WS + Feature Summary
**Verdict:** APPROVED or CHANGES REQUESTED
**Next:** `/deploy F{XX}` (if APPROVED)
