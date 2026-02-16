# Release v{X.Y.Z}

**Date:** {YYYY-MM-DD}
**Feature:** {Feature ID} - {Feature Name}

---

## Overview

{Brief description of what was added in this release - 2-3 sentences}

---

## New Features

### {Feature Name}

{Description of functionality for users}

**What's new:**
- {Item 1}
- {Item 2}
- {Item 3}

**Usage:**

```bash
# Example command or usage
hwc {command} {args}
```

**API (if applicable):**

```bash
# Example API request
curl -X POST http://localhost:8000/api/endpoint \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"key": "value"}'
```

---

## Improvements

- {Improvement 1}
- {Improvement 2}

---

## Bug Fixes

- {Fix 1}
- {Fix 2}

---

## Breaking Changes

{If no breaking changes, write "None"}

### {Breaking Change 1}

**Before:**
```python
# Old way
old_function(arg1, arg2)
```

**After:**
```python
# New way
new_function(arg1, arg2, arg3)
```

**Migration:**
1. Replace `old_function` with `new_function`
2. Add the third argument

---

## Migration Guide

{If no migration required, write "No migration required"}

### Database Migrations

```bash
# Run migrations
cd tools/hw_checker
alembic upgrade head
```

### Configuration Changes

{If configuration format changed}

```yaml
# Before
old_config: value

# After
new_config:
  nested: value
```

---

## Known Issues

{If no known issues, write "None"}

- {Issue 1}: {description} - workaround: {how to work around}
- {Issue 2}: {description}

---

## Dependencies

### Updated
- {Library 1}: v{old} -> v{new}
- {Library 2}: v{old} -> v{new}

### Added
- {New library}: v{version} - {purpose}

### Removed
- {Removed library} - {reason for removal}

---

## Contributors

- {Contributor 1}
- {Contributor 2}

---

## Full Changelog

See [CHANGELOG.md](../CHANGELOG.md) for full history.

**Workstreams in this release:**
- WS-{ID1}: {title}
- WS-{ID2}: {title}
- WS-{ID3}: {title}
