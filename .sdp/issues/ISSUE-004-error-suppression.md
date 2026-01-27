# ISSUE-004: Error Suppression with `pass`

**Severity:** 🟡 MEDIUM
**Status:** Open

## Problem

Silent exception handlers make debugging difficult:

1. `src/sdp/dashboard/dashboard_app.py:83` - Hides test runner errors
2. `src/sdp/dashboard/dashboard_app.py:107` - Hides refresh errors
3. `src/sdp/queue/state.py:41` - Hides queue state errors

## Acceptance Criteria

- [ ] No `except: pass` in F012 code
- [ ] All exceptions logged with context
- [ ] Error handling follows project patterns

## Solution

Replace `pass` with proper logging:

```python
# Before
except Exception:
    pass

# After
except Exception as e:
    logger.warning("Test runner failed: %s", e)
```

## Steps to Fix

1. Fix `dashboard_app.py:83` - Log test runner errors
2. Fix `dashboard_app.py:107` - Log refresh errors
3. Fix `queue/state.py:41` - Log queue state errors
4. Verify no `except: pass` remains
