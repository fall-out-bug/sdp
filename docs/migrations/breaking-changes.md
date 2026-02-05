# SDP Breaking Changes Migration Guide

**Version:** v0.5.0
**Last Updated:** 2026-01-30

## Table of Contents

- [Introduction](#introduction)
- [Breaking Changes Summary](#breaking-changes-summary)
- [Detailed Migration Guides](#detailed-migration-guides)
- [Troubleshooting](#troubleshooting)

---

## Introduction

This document helps you migrate between major versions of SDP by documenting all breaking changes, their rationale, and step-by-step migration instructions.

### What Are Breaking Changes?

A **breaking change** is any modification that breaks backward compatibility, requiring manual updates to your code, configuration, or workflow. Breaking changes occur when:

- APIs are removed or renamed
- File formats change structure
- Commands are deprecated or replaced
- Workflow steps are reordered

### Why We Document Breaking Changes

- **Transparency**: Clear communication about what changed and why
- **Migration Path**: Step-by-step instructions to upgrade safely
- **Timeline**: Deprecation warnings before removal
- **Rationale**: Understanding the "why" behind changes

---

## Breaking Changes Summary

| Change | Deprecated | Removed | Migration Effort | Impact |
|--------|------------|---------|------------------|--------|
| **1. Consensus → Slash Commands** | v1.2 | v0.3.0 | High | Complete workflow redesign |
| **2. WS-FFF-SS → PP-FFF-SS Format** | v0.2 | v0.3.0 | Medium | All workstream IDs |
| **3. 4-Phase → Slash Commands** | v0.1 | v0.3.0 | High | Agent coordination model |
| **4. State Machine → File-based** | v1.2 | v0.3.0 | High | `status.json` removal |
| **5. JSON → Message Router** | v0.4 | v0.5.0 | Medium | Agent messaging API |
| **6. Beads Integration** | N/A | v0.5.0 | Low | Optional feature |
| **7. QualityGateValidator Removal** | v0.4.9 | v0.5.0 | Low | Code validation |
| **8. Python SDP → SDP Plugin** | 2026-02-05 | 2026-08-03 | Medium | Complete replacement |

**⚠️ IMPORTANT: Python SDP Deprecation**

The Python SDP implementation is **deprecated** in favor of the [SDP Plugin](https://github.com/ai-masters/sdp-plugin).

**Migration timeline:**
- **2026-02-05:** Deprecation announced
- **2026-02-05 to 2026-08-03:** Maintenance period (bug fixes only)
- **2026-08-03:** Maintenance ends (community-supported)

**See [Python SDP Deprecation Guide](python-sdp-deprecation.md) for complete migration instructions.**

---

## Detailed Migration Guides

For detailed migration instructions, see:

1. [Consensus → Slash Commands](bc-001-consensus-to-slash.md)
2. [WS-FFF-SS → PP-FFF-SS Format](bc-002-workstream-id-format.md)
3. [4-Phase → Slash Commands](bc-003-four-phase-to-slash.md)
4. [State Machine → File-based](bc-004-state-to-file-based.md)
5. [JSON → Message Router](bc-005-json-to-message-router.md)
6. [Beads Integration](bc-006-beads-integration.md)
7. [QualityGateValidator Removal](bc-007-qualitygate-validator-removal.md)
8. **[Python SDP Deprecation](python-sdp-deprecation.md)** ⚠️ **IMPORTANT**

---

## ⚠️ Python SDP Deprecation (IMPORTANT)

**The Python SDP is deprecated and succeeded by the [SDP Plugin](https://github.com/ai-masters/sdp-plugin).**

### Why?

The Python implementation has fundamental limitations:
- Language lock-in (Python only)
- Heavy dependencies (Python, Poetry, pytest, mypy, ruff)
- Tool-based validation (fast but inflexible)

The SDP Plugin addresses these:
- Multi-language support (Python, Java, Go, any)
- Zero dependencies (optional Go binary)
- Language-agnostic validation (AI-based)

### Migration

**Your existing workstreams are compatible!** No conversion needed.

**Quick migration:**
```bash
# 1. Install plugin (no pip required)
git clone https://github.com/ai-masters/sdp-plugin.git ~/.claude/sdp
cp -r ~/.claude/sdp/prompts/* .claude/

# 2. Your workstreams work as-is
@build 00-001-01  # Same command!
```

**See the [complete migration guide](python-sdp-deprecation.md) for:**
- Step-by-step instructions
- Feature parity comparison
- Common questions
- Rollback plan

---

## Troubleshooting

### Common Issues

**Issue**: "ModuleNotFoundError after upgrade"
- **Solution**: Run `pip install -e .` to reinstall dependencies

**Issue**: "Tests fail with import errors"
- **Solution**: Update imports from old to new module paths (see individual guides)

**Issue**: "Workstream ID format rejected"
- **Solution**: Run migration script `scripts/migrate_workstream_ids.py`

### Getting Help

- GitHub Issues: [sdoproject/sdp/issues](https://github.com/sdoproject/sdp/issues)
- Documentation: [docs/](../README.md)
- Migration Support: Post in GitHub Discussions with "migration" tag
