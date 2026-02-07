# Platform Adapters

Multi-platform support for different AI code assistants.

## Overview

The adapter pattern allows SDP to work seamlessly across different AI platforms:

- **Claude Code** (`.claude/`) - Anthropic
- **Codex** (`.codex/`) - OpenAI  
- **OpenCode** (`.opencode/`) - Open source

## Architecture

```
PlatformAdapter (abstract)
├── install() - Create platform directory structure
├── configure_hooks() - Set up git/platform hooks
├── load_skill() - Read skill configuration
└── get_settings() - Parse platform settings

detect_platform() - Auto-detect current platform
```

## Usage

### Detect Platform

```python
from sdp.adapters import detect_platform, PlatformType

platform = detect_platform()

if platform == PlatformType.CLAUDE_CODE:
    print("Running in Claude Code")
elif platform == PlatformType.CODEX:
    print("Running in Codex")
elif platform is None:
    print("No platform detected")
```

### Implement Adapter

```python
from pathlib import Path
from typing import Any
from sdp.adapters import PlatformAdapter

class ClaudeCodeAdapter(PlatformAdapter):
    """Claude Code platform adapter."""

    def install(self, target_dir: Path) -> None:
        """Create .claude/ directory structure."""
        claude_dir = target_dir / ".claude"
        claude_dir.mkdir(exist_ok=True)
        (claude_dir / "skills").mkdir(exist_ok=True)
        (claude_dir / "settings.json").write_text("{}")

    def configure_hooks(self, hooks: list[str]) -> None:
        """Update settings.json with hooks."""
        # Implementation...
        pass

    def load_skill(self, skill_name: str) -> dict[str, Any]:
        """Read from .claude/skills/{skill_name}/SKILL.md."""
        # Implementation...
        return {"name": skill_name}

    def get_settings(self) -> dict[str, Any]:
        """Parse .claude/settings.json."""
        # Implementation...
        return {}
```

## Platform Detection

Detection searches upward from current directory until:

1. Platform directory found (`.claude/`, `.codex/`, `.opencode/`)
2. `.git` directory reached (stops search)

**Priority:** Claude Code > Codex > OpenCode

## Codex Setup

Codex uses a project-level `.codex/` directory and a user-level `~/.codex/` directory.

Project layout:
```
.codex/
├── INSTALL.md       # Setup instructions (read by Codex)
└── skills/          # Project-level skills
```

User-level layout:
```
~/.codex/
└── skills/          # Persistent skills (copied from project)
```

The `CodexAdapter` creates `.codex/INSTALL.md` and copies skills to the
user directory when `load_skill()` is called.

## OpenCode Setup

OpenCode uses a project-level `.opencode/` directory and XDG config directory.

Project layout:
```
.opencode/
├── plugin/
│   └── sdp.js       # JavaScript plugin wrapper
└── skills/          # Project-level skills
```

User-level layout (XDG):
```
~/.config/opencode/
└── skills/          # Persistent skills (copied from project)
```

The `OpenCodeAdapter` generates `.opencode/plugin/sdp.js` and copies skills
to the XDG config directory when `load_skill()` is called.

## File Structure

```
sdp/src/sdp/adapters/
├── __init__.py - Public API
├── base.py - Abstract interface + detection
├── claude_code.py - Claude Code adapter (WS-192-02)
├── codex.py - Codex adapter (WS-192-03)
└── opencode.py - OpenCode adapter (WS-192-04)
```

## Tests

```bash
# Run adapter tests
cd sdp
poetry run pytest tests/unit/adapters/ -v

# Check coverage
poetry run pytest tests/unit/adapters/ --cov=sdp.adapters --cov-report=term-missing
```

## Related Workstreams

- **WS-192-01** - Platform Adapter Interface (this workstream)
- **WS-192-02** - Claude Code Adapter Implementation
- **WS-192-03** - Codex Adapter Implementation
- **WS-192-04** - OpenCode Adapter Implementation
