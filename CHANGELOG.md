# Changelog

All notable changes to the Spec-Driven Protocol (SDP).

## [0.9.3] - 2026-02-17

### Patch Release

**Fixes:**
- Fixed TempDir cleanup in prototype tests (CI stability)

---

## [0.9.2] - 2026-02-17

### Patch Release

**New:**
- IDE selection in installer: `SDP_IDE=claude|cursor|opencode`
- OpenCode/Windsurf integration via `.opencode/` directory

**Fixes:**
- Fixed nil pointer panic in `NewWizard` (flaky TestInitCommand)

---

## [0.9.1] - 2026-02-16

### Patch Release

**Improvements:**
- **One-liner installer:** `curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | bash`
- **OpenCode integration:** `.opencode/` directory with skills, agents, commands
- **Cross-platform sync:** Commands available for Claude Code, Cursor, OpenCode, Windsurf

**Fixes:**
- Fixed Go version mismatch in `sdp-verify-dogfood.yml` (1.24 → 1.26)

**Dependencies:**
- Bump `github.com/spf13/cobra` from 1.8.0 to 1.10.2
- Bump `actions/upload-artifact` from 4 to 6

---

## [0.9.0] - 2026-02-16

### M1 Milestone - UX Excellence & Intelligent Assistance

**Theme:** Enhanced Developer Experience with Smart Recovery and Guidance

This release focuses on UX improvements, intelligent next-step recommendations, structured error handling, and guided onboarding.

### Highlights

- **Next-Step Engine:** Intelligent recommendations with confidence scoring
- **Error Taxonomy:** 38 structured error codes with recovery hints
- **Guided Onboarding:** Interactive wizard and headless mode
- **Self-Healing Doctor:** Automatic environment repair
- **Enhanced Evidence:** Full skills instrumentation
