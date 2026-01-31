# Breaking Change: ```

> **Migration Guide**: See [Breaking Changes Summary](breaking-changes.md) for all migrations

---

## Overview

**Versions**: v0.5.0

```

---

### 7. QualityGateValidator Removal

#### What Changed

The `QualityGateValidator` class was removed as dead code (P1-04).

**OLD (QualityGateValidator):**
```python
from sdp.quality import QualityGateValidator

validator = QualityGateValidator()
violations = validator.validate_file("src/module.py")
```

**NEW (Direct Validation):**
```bash
# Use hooks or CLI directly
python -m sdp.quality validate path/to/file.py

# Or use pre-commit hook
git commit  # Automatically runs hooks/pre-commit.sh
```

#### Why It Changed

| Problem | Solution |
|---------|----------|
| Unused class in codebase | Remove dead code (YAGNI) |
| Validation logic duplicated | Single source of truth in hooks |
| Unclear when to use class vs CLI | Use CLI or hooks (clear interface) |

#### Migration Steps

**Step 1: Find Usages of QualityGateValidator**

```bash
grep -r "QualityGateValidator" --include="*.py" .
```

**Step 2: Replace with Direct Validation**

```python
# OLD
from sdp.quality import QualityGateValidator

validator = QualityGateValidator()
violations = validator.validate_file("src/module.py")
if violations:
    print(f"Found {len(violations)} violations")

# NEW
import subprocess

result = subprocess.run(
    ["python", "-m", "sdp.quality", "validate", "src/module.py"],
    capture_output=True,
    text=True
)
if result.returncode != 0:
    print(f"Validation failed:\n{result.stdout}")
```

**Step 3: Use Git Hooks (Recommended)**

```bash
# Copy hook to .git/hooks/
cp scripts/hooks/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# Now validation runs automatically on commit
git commit  # Validates all changed files
```

**Step 4: Update Tests**

```python
# OLD
def test_quality_gate():
    validator = QualityGateValidator()
    violations = validator.validate_file("test_file.py")
    assert len(violations) == 0

# NEW
def test_quality_gate():
    result = subprocess.run(
        ["python", "-m", "sdp.quality", "validate", "test_file.py"],
        capture_output=True,
    )
    assert result.returncode == 0
```

#### Before/After Comparison

**OLD (QualityGateValidator):**
```python
from sdp.quality import QualityGateValidator

# Initialize
validator = QualityGateValidator()

# Validate
violations = validator.validate_directory("src/")
if violations:
    validator.print_report()
    exit(1)
```

**NEW (CLI + Hooks):**
```bash
# Option 1: Direct CLI
python -m sdp.quality validate src/

# Option 2: Pre-commit hook (automatic)
git commit  # Runs validation automatically

# Option 3: Make target
make quality  # Defined in Makefile
```

#### Timeline

- **Deprecated:** 2026-01-15 (v0.4.9)
- **Removed:** 2026-01-30 (v0.5.0)
- **Migration Support:** Ended 2026-02-15

#### Impact Analysis

**Files Removed:**
- `src/sdp/quality/validator.py` (QualityGateValidator class)

**Files Kept:**
- `src/sdp/quality/models.py` (QualityGateConfig, ValidationIssue)
- `src/sdp/quality/config.py` (TOML configuration parser)
- `scripts/hooks/pre-commit.sh` (Git hook)
- `scripts/hooks/post-commit.sh` (Git hook)

**Migration Effort:** Low (replace with CLI or hooks)

---

## Troubleshooting

### Issue: Legacy Workstream IDs Still Present

**Symptom:**
```bash
grep -r "ws_id: WS-" docs/workstreams/
