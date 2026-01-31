# PRODUCT_VISION.md

> **Last Updated:** 2026-01-26
> **Version:** 1.0

## Mission

SDP is a spec-driven development framework that enables AI agents and humans to collaborate on building reliable, maintainable software.

## Users

1. **Solo developers** building AI-assisted projects
2. **Small teams** (2-5 engineers) with AI workflows
3. **DevOps engineers** integrating SDP into CI/CD

## Success Metrics

- [ ] Time from idea to running code: <1 hour
- [ ] New user onboarding: <30 min to first WS
- [ ] Test feedback latency: <2 seconds (watch mode)
- [ ] Feature completion rate: >85%

## Strategic Tradeoffs

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| DX vs Control | Prioritize DX | Friction kills adoption |
| Speed vs Quality | Both | TDD + quality gates |
| Simple vs Expressive | Progressive | Simple entry, power when needed |

## Non-Goals

- Real-time collaboration (multiplayer)
- Enterprise SSO (out of scope)
- Language-agnostic (Python-first, extensible)
