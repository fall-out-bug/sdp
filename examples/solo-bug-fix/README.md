# Example: Solo Mode Bug Fix

A simple example of using Solo mode to fix a bug.

## Scenario

Users report they can't login with emails containing a plus sign (e.g., `user+test@gmail.com`).

## The Conversation

### Step 1: Describe the Bug

```
User:
Fix bug: users can't login with emails containing plus sign.
Example: user+test@gmail.com gets "Invalid email" error.
Look at src/auth/validators.py
```

### Step 2: AI Investigates

AI reads the file and finds:
```python
# src/auth/validators.py
import re

def is_valid_email(email: str) -> bool:
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return bool(re.match(pattern, email))
```

AI notices the regex already includes `+`, so the issue is elsewhere.

### Step 3: AI Digs Deeper

```
AI:
The regex looks correct. Let me check where validation happens.
Looking at src/auth/login.py...
```

Finds:
```python
# src/auth/login.py
def login(email: str, password: str):
    # Bug: sanitizing removes + before validation
    email = email.replace("+", "")
    if not is_valid_email(email):
        raise InvalidEmailError("Invalid email")
    ...
```

### Step 4: AI Fixes and Tests

```
AI:
Found the bug! In login.py line 4, the + is being stripped before validation.
This was probably intended to sanitize input but breaks valid emails.

Fix: Remove the sanitization line, as + is a valid email character.
```

AI creates the fix:
```python
def login(email: str, password: str):
    # Removed: email = email.replace("+", "")
    email = email.strip().lower()  # Only normalize whitespace and case
    if not is_valid_email(email):
        raise InvalidEmailError("Invalid email")
    ...
```

And adds a test:
```python
# tests/test_login.py
def test_login_accepts_plus_in_email():
    # Create test user with plus email
    user = create_user("user+test@gmail.com", "password123")

    # Should not raise
    result = login("user+test@gmail.com", "password123")
    assert result.user_id == user.id
```

### Step 5: Verify

```
User:
Run the tests

AI:
pytest tests/test_login.py -v
...
tests/test_login.py::test_login_accepts_plus_in_email PASSED
...
All tests pass.
```

## Complete Conversation

```
You: Fix bug: users can't login with emails containing plus sign.
     Example: user+test@gmail.com gets "Invalid email" error.
     Look at src/auth/validators.py

AI:  [reads validators.py] Regex looks correct, checking login.py...
     [reads login.py] Found it! Line 4 strips the + character.

     The fix is to remove that line. Here's the change:
     [shows diff]

     I'll also add a regression test.
     [shows test]

You: Run the tests

AI:  pytest tests/test_login.py -v
     All 12 tests pass including the new one.

You: Looks good, commit it

AI:  git add src/auth/login.py tests/test_login.py
     git commit -m "Fix: allow plus sign in email addresses

     The login function was stripping + characters before validation,
     breaking valid email addresses like user+tag@gmail.com.

     Added regression test."
```

## Key Points

1. **Single conversation** - no mode switching needed
2. **Iterative investigation** - AI dug deeper when first guess was wrong
3. **Test included** - prevents regression
4. **Small scope** - one bug, one fix

## Time

~10 minutes from bug report to committed fix.
