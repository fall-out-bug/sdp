# Extending SDP

Guide to extending SDP with custom skills, validators, and integrations.

---

## Table of Contents

- [Custom Skills](#custom-skills)
- [Custom Validators](#custom-validators)
- [Custom Quality Gates](#custom-quality-gates)
- [Custom Integrations](#custom-integrations)
- [Examples](#examples)

---

## Custom Skills

### Skill Structure

**Directory:**
```
.claude/skills/{skill_name}/
├── SKILL.md          # Required: Skill definition
├── prompt.md         # Optional: Additional prompts
└── examples/         # Optional: Usage examples
```

### SKILL.md Template

```markdown
# @{skill-name}

One-line description of what this skill does.

## Usage
\`\`\`bash
@{skill-name} "argument"
\`\`\`

## Process
1. Step one
2. Step two
3. Step three

## Output
What the skill produces

## Examples
### Basic Usage
\`\`\`bash
@{skill-name} "example"
\`\`\`

Result: What happens

### Advanced Usage
\`\`\`bash
@{skill-name} "complex example" --option=value
\`\`\`

Result: What happens with options

## See Also
- [Related skill](../other-skill/)
- [Documentation](../../../docs/)
```

### Skill Invocation

**From Claude Code:**
```bash
@{skill-name} "argument"
```

**Programmatic:**
```python
from sdp.cli import main

# Invoke skill
main(['@skill-name', 'argument'])
```

### Skill Context

**Available Context:**
- `CLAUDE_CODE_TEAM_NAME` - Current team (if in swarm)
- `CLAUDE_CODE_AGENT_ID` - Agent ID
- `CLAUDE_CODE_PLAN_MODE_REQUIRED` - Whether plan mode is required

**Access in Skill:**
```markdown
## Environment
- Team: {{CLAUDE_CODE_TEAM_NAME}}
- Agent: {{CLAUDE_CODE_AGENT_ID}}
- Plan Mode: {{CLAUDE_CODE_PLAN_MODE_REQUIRED}}
```

---

## Custom Validators

### Validator Structure

**Location:** `src/sdp/validators/{name}.py`

**Template:**
```python
"""Custom validator for SDP."""

from pathlib import Path
from dataclasses import dataclass
from typing import List

@dataclass
class ValidationResult:
    """Result of validation."""
    passed: bool
    issues: List[str]

class CustomValidator:
    """Custom validator."""

    def __init__(self, config: dict):
        """Initialize validator.

        Args:
            config: Validator configuration
        """
        self.config = config

    def validate(self, file_path: Path) -> ValidationResult:
        """Validate file.

        Args:
            file_path: Path to file

        Returns:
            ValidationResult with findings
        """
        issues = []

        # Custom validation logic
        if self._check_something(file_path):
            issues.append("Issue found")

        return ValidationResult(
            passed=len(issues) == 0,
            issues=issues
        )

    def _check_something(self, file_path: Path) -> bool:
        """Check specific condition.

        Args:
            file_path: Path to file

        Returns:
            True if condition fails
        """
        # Implementation
        return False
```

### Register Validator

**Location:** `src/sdp/validators/__init__.py`

```python
from .custom_validator import CustomValidator

__all__ = ['CustomValidator']
```

### Use Validator

**In Quality Gate:**
```python
from sdp.validators import CustomValidator

validator = CustomValidator(config)
result = validator.validate(file_path)

if not result.passed:
    for issue in result.issues:
        print(f"Issue: {issue}")
```

---

## Custom Quality Gates

### Configuration

**Location:** `quality-gate.toml`

```toml
[my_custom_check]
enabled = true
threshold = 100
severity = "error"
```

### Quality Gate Implementation

**Location:** `src/sdp/quality/custom_checks.py`

```python
"""Custom quality gate checks."""

from pathlib import Path
from sdp.quality.config import QualityGateConfig

class CustomQualityCheck:
    """Custom quality gate check."""

    def __init__(self, config: QualityGateConfig):
        """Initialize check.

        Args:
            config: Quality gate configuration
        """
        self.config = config
        self.enabled = config.get('my_custom_check', {}).get('enabled', False)
        self.threshold = config.get('my_custom_check', {}).get('threshold', 100)

    def check(self, file_path: Path) -> bool:
        """Run custom check.

        Args:
            file_path: Path to check

        Returns:
            True if check passes
        """
        if not self.enabled:
            return True

        # Custom check logic
        lines = file_path.read_text().splitlines()
        return len(lines) <= self.threshold
```

### Register Check

**Location:** `src/sdp/quality/validator.py`

```python
from .custom_checks import CustomQualityCheck

class QualityGateValidator:
    """Quality gate validator."""

    def __init__(self, config: QualityGateConfig):
        self.custom_check = CustomQualityCheck(config)
```

---

## Custom Integrations

### Integration Template

**Location:** `src/sdp/integrations/{name}.py`

```python
"""Custom integration for SDP."""

from dataclasses import dataclass
from typing import Any, Dict

@dataclass
class IntegrationConfig:
    """Integration configuration."""
    api_key: str
    endpoint: str
    options: Dict[str, Any]

class CustomIntegration:
    """Custom integration."""

    def __init__(self, config: IntegrationConfig):
        """Initialize integration.

        Args:
            config: Integration configuration
        """
        self.config = config
        self.client = self._create_client()

    def _create_client(self):
        """Create API client."""
        # Implementation
        pass

    def notify(self, message: str) -> bool:
        """Send notification.

        Args:
            message: Message to send

        Returns:
            True if successful
        """
        # Implementation
        return True

    def sync(self, data: Dict[str, Any]) -> bool:
        """Sync data.

        Args:
            data: Data to sync

        Returns:
            True if successful
        """
        # Implementation
        return True
```

### Register Integration

**Location:** `src/sdp/__init__.py`

```python
from sdp.integrations.custom import CustomIntegration

__all__ = ['CustomIntegration']
```

---

## Examples

### Example 1: Custom Skill

**Skill:** `@review-pr`

**Purpose:** Review pull request

**File:** `.claude/skills/review-pr/SKILL.md`

```markdown
# @review-pr

Review pull request for quality.

## Usage
\`\`\`bash
@review-pr PR-123
\`\`\`

## Process
1. Fetch PR diff
2. Check quality gates
3. Review changes
4. Provide feedback

## Output
Review report with:
- Quality gate status
- Code review findings
- Approval/rejection

## Examples
\`\`\`bash
@review-pr PR-456
\`\`\`

Result:
✅ Quality gates: PASSED
✅ Code review: APPROVED
✅ Coverage: 85%
```

---

### Example 2: Custom Validator

**Validator:** `CommentDensityValidator`

**Purpose:** Ensure adequate code comments

**File:** `src/sdp/validators/comment_density.py`

```python
"""Comment density validator."""

from pathlib import Path
from typing import List

class CommentDensityValidator:
    """Validate comment density."""

    def __init__(self, min_ratio: float = 0.1):
        """Initialize validator.

        Args:
            min_ratio: Minimum comment ratio (comments / lines)
        """
        self.min_ratio = min_ratio

    def validate(self, file_path: Path) -> bool:
        """Validate comment density.

        Args:
            file_path: Path to Python file

        Returns:
            True if comment density adequate
        """
        content = file_path.read_text()
        lines = content.splitlines()

        code_lines = 0
        comment_lines = 0

        for line in lines:
            stripped = line.strip()
            if not stripped or stripped.startswith('#'):
                if stripped.startswith('#'):
                    comment_lines += 1
            else:
                code_lines += 1

        if code_lines == 0:
            return True

        ratio = comment_lines / code_lines
        return ratio >= self.min_ratio
```

---

### Example 3: Custom Integration

**Integration:** Slack notifications

**Purpose:** Send notifications to Slack

**File:** `src/sdp/integrations/slack.py`

```python
"""Slack integration for SDP."""

import requests
from dataclasses import dataclass

@dataclass
class SlackConfig:
    """Slack configuration."""
    webhook_url: str
    channel: str
    username: str = "SDP Bot"

class SlackIntegration:
    """Slack integration."""

    def __init__(self, config: SlackConfig):
        """Initialize integration.

        Args:
            config: Slack configuration
        """
        self.config = config

    def send_message(self, message: str) -> bool:
        """Send message to Slack.

        Args:
            message: Message to send

        Returns:
            True if successful
        """
        payload = {
            "channel": self.config.channel,
            "username": self.config.username,
            "text": message
        }

        response = requests.post(
            self.config.webhook_url,
            json=payload
        )

        return response.status_code == 200

    def notify_workstream_complete(self, ws_id: str) -> bool:
        """Notify workstream completion.

        Args:
            ws_id: Workstream ID

        Returns:
            True if successful
        """
        message = f"✅ Workstream {ws_id} completed"
        return self.send_message(message)
```

---

## Best Practices

### DO ✅

1. **Follow conventions** - Use existing patterns
2. **Document well** - Clear docs and examples
3. **Test thoroughly** - ≥80% coverage
4. **Handle errors** - Graceful failure
5. **Type hints** - Full mypy compliance

### DON'T ❌

1. **Don't break compatibility** - Maintain backwards compat
2. **Don't skip validation** - Always validate inputs
3. **Don't hardcode config** - Use config files
4. **Don't ignore errors** - Handle all exceptions
5. **Don't duplicate code** - Reuse existing utilities

---

## Testing Custom Code

### Unit Tests

**File:** `tests/unit/integrations/test_slack.py`

```python
"""Tests for Slack integration."""

import pytest
from sdp.integrations.slack import SlackIntegration, SlackConfig

def test_send_message(mocker):
    """Test sending message."""
    config = SlackConfig(
        webhook_url="https://hooks.slack.com/test",
        channel="#test"
    )
    integration = SlackIntegration(config)

    mock_post = mocker.patch('requests.post')
    mock_post.return_value.status_code = 200

    result = integration.send_message("Test message")

    assert result is True
    mock_post.assert_called_once()

def test_notify_workstream_complete(mocker):
    """Test workstream completion notification."""
    config = SlackConfig(
        webhook_url="https://hooks.slack.com/test",
        channel="#test"
    )
    integration = SlackIntegration(config)

    mock_send = mocker.patch.object(integration, 'send_message')
    mock_send.return_value = True

    result = integration.notify_workstream_complete("WS-001-01")

    assert result is True
    mock_send.assert_called_once_with("✅ Workstream WS-001-01 completed")
```

---

## See Also

- [architecture.md](architecture.md) - System architecture
- [contributing.md](contributing.md) - Contribution guide
- [development.md](development.md) - Development setup
- [reference/skills.md](../reference/skills.md) - Skill reference

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
