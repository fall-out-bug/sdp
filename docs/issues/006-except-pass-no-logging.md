# 006: except Exception: pass without logging

**Source:** F020 Review (2026-01-31)  
**Status:** Closed (00-020-04)  
**Priority:** P3  
**Route:** @bugfix (tech debt)

## Problem

`src/sdp/hooks/common.py` contains two instances of `except Exception: pass` without logging, violating quality gate "except Exception only with logging".

**Locations:**
- Line 41-42: `find_project_root()` — TOML parse fallback
- Line 79-80: `find_workstream_dir()` — config parse fallback

## Code

```python
# Line 41-42
try:
    config = tomllib.loads((path / "pyproject.toml").read_text())
    if "tool" in config and "sdp" in config.get("tool", {}):
        return path
except Exception:
    pass  # ← No logging

# Line 79-80
try:
    config = tomllib.loads(config_file.read_text())
    # ...
except Exception:
    pass  # ← No logging
```

## Fix Options

1. **Add logging:**
   ```python
   except Exception as e:
       logger.debug("Failed to parse TOML: %s", e)
   ```

2. **Catch specific exception:**
   ```python
   except tomllib.TOMLDecodeError:
       pass  # TOML parsing failed, try next marker
   ```

## Related

- F020 Review Report: `docs/reports/2026-01-31-F020-review.md`
- Quality gate: Error handling (no bare except)
