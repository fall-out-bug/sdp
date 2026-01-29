# Quality Gate Quick Reference

## For Developers

### Basic Usage

```python
from sdp.quality import QualityGateValidator

# Validate your code
validator = QualityGateValidator()
violations = validator.validate_file("path/to/your/file.py")

# Check results
if violations:
    for v in violations:
        print(f"{v.file_path}:{v.line_number} - {v.message}")
else:
    print("All quality gates passed!")
```

### Common Violations and Fixes

#### 1. File Too Large
**Error:** `File too large: 250 lines (max: 200)`

**Fix:** Split the file into smaller modules by responsibility.

#### 2. Missing Type Hints
**Error:** `Function 'process_data' missing return type annotation`

**Fix:**
```python
# Before
def process_data(data):
    return data.strip()

# After
def process_data(data: str) -> str:
    return data.strip()
```

#### 3. Bare Except
**Error:** `Bare except clause detected (forbidden)`

**Fix:**
```python
# Before
try:
    risky_operation()
except:
    pass

# After
try:
    risky_operation()
except ValueError as e:
    logger.error(f"Invalid value: {e}")
    raise
```

#### 4. High Complexity
**Error:** `Function has complexity 15 (max: 10)`

**Fix:** Extract helper methods to reduce nesting.

```python
# Before
def complex_function(data):
    if data:
        if data.get('field1'):
            if data.get('field2'):
                # ... 10 more nested ifs
                pass

# After
def complex_function(data):
    if not data:
        return None
    if not has_required_fields(data):
        return None
    # ... simpler logic
```

#### 5. Hardcoded Secrets
**Error:** `Possible hardcoded secret detected`

**Fix:** Use environment variables.

```python
# Before
password = "SuperSecret123"

# After
import os
password = os.getenv("DB_PASSWORD")
assert password, "DB_PASSWORD not set"
```

### Configuration Tips

#### Tighten Standards Over Time
```toml
# Start lenient for legacy code
[complexity]
max_cc = 20

# Gradually tighten
[complexity]
max_cc = 15  # After 1 month
max_cc = 10  # After 3 months
```

#### Exclude Test Files from Coverage
```toml
[coverage]
exclude_patterns = [
    "*/tests/*",
    "*/test_*.py",
    "*/conftest.py"
]
```

#### Disable Specific Checks
```toml
[naming]
enabled = false  # If not needed
```

### CI/CD Integration

#### Fail Build on Errors
```yaml
- name: Quality Gate
  run: |
    python -c "
    from sdp.quality import QualityGateValidator
    validator = QualityGateValidator()
    violations = validator.validate_directory('src/')
    errors = [v for v in violations if v.severity == 'error']
    if errors:
        print('Quality gate failed!')
        exit(1)
    "
```

### Pre-commit Hook

Create `.git/hooks/pre-commit`:
```bash
#!/bin/bash
python -c "
from sdp.quality import QualityGateValidator
import sys

validator = QualityGateValidator()
violations = validator.validate_directory('src/')
errors = [v for v in violations if v.severity == 'error']

if errors:
    print(f'❌ Quality gate failed with {len(errors)} errors')
    for e in errors[:5]:
        print(f'  {e}')
    sys.exit(1)
else:
    print('✅ Quality gate passed')
"
```

### Troubleshooting

#### Q: How do I temporarily disable a check?
A: Set `enabled = false` in `quality-gate.toml`:
```toml
[complexity]
enabled = false  # Temporary disable
```

#### Q: What if my legacy code fails?
A: Use `exclude_patterns` or introduce standards gradually:
```toml
[coverage]
exclude_patterns = [
    "*/legacy/*"  # Exclude old code
]
```

#### Q: Can I have different rules per module?
A: Yes, place `quality-gate.toml` in subdirectories:
```
project/
├── quality-gate.toml       # Default rules
└── critical_module/
    └── quality-gate.toml   # Stricter rules
```

### Best Practices

1. **Start Early:** Implement quality gates from day one
2. **Be Realistic:** Set achievable thresholds for your team
3. **Iterate:** Tighten standards gradually over time
4. **Automate:** Integrate into CI/CD and pre-commit hooks
5. **Educate:** Help team understand why standards matter
6. **Review:** Regularly review and update configuration

### Resources

- Full documentation: `docs/quality-gate-schema.md`
- Examples: `examples/quality_gate_example.py`
- Tests: `tests/test_quality_gate.py`

---

**Need Help?** Check the full documentation or ask the team!
