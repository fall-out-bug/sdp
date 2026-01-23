# Systematic Debugging — 4-Phase Root Cause Analysis

**Purpose:** Replace trial-and-error with scientific method. Evidence-based, not assumption-based.

---

## Core Principles

1. **Evidence First** — Collect facts before guessing
2. **One Hypothesis** — Test one thing at a time
3. **Minimal Change** — Smallest possible fix
4. **Failsafe Rule** — 3 failed fixes → STOP, question architecture

---

## 4-Phase Process

### Phase 1: Evidence Collection

**Goal:** Gather all observable facts about the problem.

#### Checklist

- [ ] **Error Messages**
  ```bash
  # Collect all error logs
  grep -r "ERROR\|Exception\|Traceback" logs/ --include="*.log" | tail -50
  
  # Stack traces
  grep -A 20 "Traceback" logs/*.log | head -100
  ```

- [ ] **Reproduce the Issue**
  ```bash
  # Document exact steps
  # 1. What command/action triggers it?
  # 2. What environment? (OS, Python version, dependencies)
  # 3. How consistently? (always/sometimes/rarely)
  # 4. What should happen vs. what actually happens?
  ```

- [ ] **Recent Changes**
  ```bash
  # Git history
  git log --oneline --since="7 days ago" --all
  
  # Changed files
  git diff HEAD~5..HEAD --stat
  
  # Specific file changes
  git log -p --since="7 days ago" -- path/to/suspicious/file.py
  ```

- [ ] **Environment State**
  ```bash
  # Dependencies
  pip list | grep -E "package-name|version"
  
  # Environment variables
  env | grep -E "HW_CHECKER|SDP"
  
  # System info
  python --version
  uname -a
  ```

#### Output Format

```markdown
### Phase 1: Evidence Collection

**Error Messages:**
```
[Paste relevant error logs]
```

**Reproduction Steps:**
1. [Step 1]
2. [Step 2]
3. [Step 3]

**Expected:** [What should happen]
**Actual:** [What actually happens]
**Consistency:** [always/sometimes/rarely]

**Recent Changes:**
- [File 1]: [What changed]
- [File 2]: [What changed]

**Environment:**
- Python: [version]
- OS: [version]
- Key dependencies: [list]
```

---

### Phase 2: Pattern Analysis

**Goal:** Find working examples and compare with broken case.

#### Checklist

- [ ] **Find Working Examples**
  ```bash
  # Similar code that works
  grep -r "similar_pattern" src/ --include="*.py"
  
  # Test cases that pass
  pytest tests/ -k "similar_test" -v
  
  # Git history: when did it last work?
  git log --all --grep="fix\|working" --oneline
  ```

- [ ] **Compare Working vs. Broken**
  ```markdown
  | Aspect | Working Case | Broken Case | Difference |
  |--------|--------------|-------------|------------|
  | Input | [value] | [value] | [diff] |
  | Context | [context] | [context] | [diff] |
  | Code path | [path] | [path] | [diff] |
  ```

- [ ] **Identify Patterns**
  - What's common in working cases?
  - What's unique in broken case?
  - What changed between working and broken?

#### Output Format

```markdown
### Phase 2: Pattern Analysis

**Working Examples:**
- [Example 1]: [Why it works]
- [Example 2]: [Why it works]

**Comparison Table:**
| Aspect | Working | Broken | Difference |
|--------|---------|--------|------------|
| [Aspect 1] | [value] | [value] | [diff] |

**Pattern Identified:**
[What pattern emerges from comparison]
```

---

### Phase 3: Hypothesis Testing

**Goal:** Form one hypothesis, test with minimal change.

#### Rules

1. **ONE hypothesis at a time** — No parallel testing
2. **Minimal change** — Smallest possible modification
3. **Isolated test** — Test in isolation, not in full system
4. **Clear pass/fail** — Binary outcome, no ambiguity

#### Checklist

- [ ] **Form Hypothesis**
  ```markdown
  **Hypothesis:** [One clear statement]
  
  **Reasoning:**
  - Evidence from Phase 1: [fact]
  - Pattern from Phase 2: [pattern]
  - Logic: [why this hypothesis]
  
  **Prediction:** If hypothesis is correct, then [expected outcome]
  ```

- [ ] **Design Minimal Test**
  ```python
  # Minimal reproduction
  def test_hypothesis():
      # Setup: minimal context
      # Action: one change
      # Assert: clear pass/fail
      pass
  ```

- [ ] **Execute Test**
  ```bash
  # Run isolated test
  pytest tests/unit/test_hypothesis.py::test_hypothesis -v
  
  # Or manual verification
  python -c "from module import function; function(test_input)"
  ```

- [ ] **Record Result**
  ```markdown
  **Test Result:** [PASS/FAIL]
  
  **If PASS:**
  - Hypothesis confirmed → proceed to Phase 4
  
  **If FAIL:**
  - Hypothesis rejected → form new hypothesis
  - Update evidence with new test result
  ```

#### Output Format

```markdown
### Phase 3: Hypothesis Testing

**Hypothesis #1:**
[Clear statement]

**Reasoning:**
- Evidence: [from Phase 1]
- Pattern: [from Phase 2]

**Test:**
```python
[Minimal test code]
```

**Result:** [PASS/FAIL]

**Conclusion:** [Confirmed/Rejected]
```

---

### Phase 4: Implementation

**Goal:** Fix the root cause using TDD approach.

#### Checklist

- [ ] **Write Failing Test First**
  ```python
  def test_fix_issue():
      # Reproduce the bug
      result = broken_function(input)
      assert result == expected  # This should fail initially
  ```

- [ ] **Implement Minimal Fix**
  ```python
  def broken_function(input):
      # Minimal change to make test pass
      # No refactoring, no improvements
      # Just fix the root cause
      pass
  ```

- [ ] **Verify Fix**
  ```bash
  # Test passes
  pytest tests/unit/test_fix.py::test_fix_issue -v
  
  # Regression: existing tests still pass
  pytest tests/unit/ -m fast -q
  
  # Integration: full system works
  pytest tests/integration/ -k "related" -v
  ```

- [ ] **Document Root Cause**
  ```markdown
  **Root Cause:** [Clear explanation]
  
  **Fix:** [What changed]
  
  **Why This Fix Works:** [Explanation]
  ```

#### Output Format

```markdown
### Phase 4: Implementation

**Failing Test:**
```python
[Test code]
```

**Fix:**
```python
[Minimal fix code]
```

**Verification:**
- Unit test: ✅ PASS
- Regression: ✅ PASS
- Integration: ✅ PASS

**Root Cause:**
[Clear explanation of root cause]

**Fix Explanation:**
[Why this fix works]
```

---

## Root-Cause Tracing Technique

**Purpose:** Trace from symptom to root cause through the call stack.

### Method

1. **Start at Symptom**
   ```python
   # Error occurs here
   result = function_a()  # Line 42
   ```

2. **Trace Backward**
   ```python
   # Where did function_a get its input?
   def function_a():
       return function_b(data)  # Line 30
   
   # Where did function_b get its input?
   def function_b(data):
       return process(data)  # Line 15
   ```

3. **Find First Bad State**
   ```python
   # First place where data becomes incorrect
   def process(data):
       # data is already wrong here
       # Root cause is before this function
   ```

4. **Identify Root Cause**
   ```python
   # Root cause: incorrect data creation
   data = create_data(input)  # Line 5 - BUG HERE
   ```

### Visualization

```
Symptom (Error)
    ↓
Function A (receives bad data)
    ↓
Function B (passes bad data)
    ↓
Function C (creates bad data) ← ROOT CAUSE
```

### Tools

```bash
# Python debugger
python -m pdb script.py

# Stack trace analysis
python -c "
import traceback
try:
    broken_code()
except Exception as e:
    traceback.print_exc()
"

# Call graph
python -m pycallgraph graphviz -- ./script.py
```

---

## Failsafe Rule: 3 Strikes

**Purpose:** Prevent infinite fix loops. If 3 fixes fail, question the architecture.

### Rule

**After 3 failed fix attempts → STOP, escalate to architecture review.**

### Tracking

```markdown
### Fix Attempts

**Attempt #1:**
- Hypothesis: [hypothesis]
- Fix: [what was changed]
- Result: ❌ FAIL
- Reason: [why it failed]

**Attempt #2:**
- Hypothesis: [hypothesis]
- Fix: [what was changed]
- Result: ❌ FAIL
- Reason: [why it failed]

**Attempt #3:**
- Hypothesis: [hypothesis]
- Fix: [what was changed]
- Result: ❌ FAIL
- Reason: [why it failed]

**🚨 FAILSAFE TRIGGERED**

**Analysis:**
- All 3 attempts failed
- Common pattern: [what's common in failures]
- Architecture issue suspected: [what architecture problem]

**Next Steps:**
1. Review architecture: [what to review]
2. Consider refactoring: [what to refactor]
3. Create new WS: [WS-ID for architecture fix]
```

### When to Trigger

- ✅ 3 distinct fix attempts failed
- ✅ Each attempt was based on different hypothesis
- ✅ Each attempt was properly tested
- ✅ Root cause still not found

### Escalation

When failsafe triggers:

1. **Document the pattern** — What's common in all failures?
2. **Question architecture** — Is the design fundamentally flawed?
3. **Create architecture WS** — New workstream to fix design
4. **Do NOT continue fixing** — Stop debugging, start redesigning

---

## Complete Workflow Example

```markdown
# /debug "API returns 500 on /submissions endpoint"

## Phase 1: Evidence Collection

**Error Messages:**
```
2024-01-15 10:23:45 ERROR: Internal server error
Traceback (most recent call last):
  File "src/hw_checker/presentation/api/routers/submissions.py", line 42
    result = service.get_submission(id)
  File "src/hw_checker/application/submission_service.py", line 30
    return repository.find_by_id(id)
  File "src/hw_checker/infrastructure/repositories.py", line 15
    raise ValueError("ID not found")
```

**Reproduction:**
1. POST /submissions → returns 201 with ID=123
2. GET /submissions/123 → returns 500

**Recent Changes:**
- repository.py: Changed find_by_id logic (2 days ago)

## Phase 2: Pattern Analysis

**Working Examples:**
- GET /submissions/456 works (ID from old format)
- GET /submissions/123 fails (ID from new format)

**Comparison:**
| Aspect | Working (456) | Broken (123) | Difference |
|--------|---------------|--------------|------------|
| ID format | "sub-456" | "123" | Format changed |
| Created | 7 days ago | 2 days ago | After change |

**Pattern:** New ID format not handled by repository

## Phase 3: Hypothesis Testing

**Hypothesis:** Repository.find_by_id() doesn't handle new ID format

**Test:**
```python
def test_new_id_format():
    repo = Repository()
    result = repo.find_by_id("123")  # New format
    assert result is not None
```

**Result:** FAIL (raises ValueError)

**Conclusion:** Hypothesis confirmed

## Phase 4: Implementation

**Failing Test:**
```python
def test_find_by_id_new_format():
    repo = Repository()
    submission = repo.find_by_id("123")
    assert submission.id == "123"
```

**Fix:**
```python
def find_by_id(self, id: str) -> Submission:
    # Handle both old and new formats
    if id.startswith("sub-"):
        return self._find_old_format(id)
    else:
        return self._find_new_format(id)
```

**Verification:**
- Unit test: ✅ PASS
- Regression: ✅ PASS
- Integration: ✅ PASS

**Root Cause:** Repository only handled old ID format ("sub-XXX"), not new format ("XXX")

**Fix Explanation:** Added format detection to handle both formats
```

---

## Integration with /issue

The `/debug` command is a focused debugging tool. For full issue analysis (severity, routing, GitHub), use `/issue`.

**When to use `/debug`:**
- You have a specific bug to fix
- You need systematic root cause analysis
- You want to follow 4-phase process

**When to use `/issue`:**
- You need severity classification
- You need routing to /hotfix or /bugfix
- You want GitHub issue creation
- You need full issue lifecycle management

---

## Key Takeaways

1. **Evidence First** — Never guess, always collect facts
2. **One at a Time** — Test one hypothesis per iteration
3. **Minimal Change** — Smallest fix that works
4. **3 Strikes Rule** — After 3 failures, question architecture
5. **TDD in Phase 4** — Failing test first, then fix

---

**Version:** 1.0.0  
**Last Updated:** 2024-01-15  
**Related:** `/issue`, `/hotfix`, `/bugfix`
