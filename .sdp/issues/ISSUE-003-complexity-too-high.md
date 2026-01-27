# ISSUE-003: Cyclomatic Complexity > 10

**Severity:** 🟡 MEDIUM
**File:** src/sdp/webhook/server.py:53
**Status:** Open

## Problem

```
src/sdp/webhook/server.py:53:9: C901 `start` is too complex (13 > 10)
```

The `start()` function contains a nested `WebhookHandler` class definition with multiple methods, inflating complexity.

## Acceptance Criteria

- [ ] All F012 functions have CC < 10
- [ ] ruff check passes with no C901 errors
- [ ] Functionality unchanged

## Solution

Extract `WebhookHandler` class to module level:

```python
# Before (nested class, CC=13)
def start(self) -> None:
    class WebhookHandler(BaseHTTPRequestHandler):
        def do_POST(self): ...
        def do_GET(self): ...
        def log_message(self): ...

# After (module-level class, CC<10 per function)
class WebhookHandler(BaseHTTPRequestHandler):
    def do_POST(self): ...
    def do_GET(self): ...
    def log_message(self): ...

def start(self) -> None:
    server = HTTPServer(...)
```

## Steps to Fix

1. Extract `WebhookHandler` to module level
2. Pass `server_instance` via constructor or closure
3. Update tests
4. Verify ruff check passes
