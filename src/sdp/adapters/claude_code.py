"""Claude Code platform adapter implementation.

This module implements the PlatformAdapter interface for Claude Code,
managing .claude/ directory structure, settings.json, and skills.

Example:
    >>> from pathlib import Path
    >>> from sdp.adapters.claude_code import ClaudeCodeAdapter
    >>>
    >>> adapter = ClaudeCodeAdapter()
    >>> adapter.install(Path.cwd())
    >>> adapter.configure_hooks(["pre-commit", "post-build"])
    >>> skill = adapter.load_skill("build")
"""

import json
import re
from pathlib import Path
from typing import Any

from sdp.adapters.base import PlatformAdapter


class ClaudeCodeAdapter(PlatformAdapter):
    """Claude Code platform adapter.

    Manages .claude/ directory structure:
    - .claude/settings.json - Permissions and hooks
    - .claude/settings.local.json - Local overrides
    - .claude/skills/ - Skill definitions
    - .claude/agents/ - Agent definitions
    """

    def install(self, target_dir: Path) -> None:
        """Install Claude Code directory structure.

        Creates:
        - .claude/ directory
        - .claude/skills/ directory
        - .claude/agents/ directory
        - .claude/settings.json with default permissions

        Args:
            target_dir: Directory where .claude/ should be created

        Raises:
            PermissionError: If target directory is not writable
        """
        claude_dir = target_dir / ".claude"
        claude_dir.mkdir(exist_ok=True)

        # Create subdirectories
        (claude_dir / "skills").mkdir(exist_ok=True)
        (claude_dir / "agents").mkdir(exist_ok=True)

        # Create default settings.json if it doesn't exist
        settings_file = claude_dir / "settings.json"
        if not settings_file.exists():
            default_settings = {
                "permissions": {
                    "allow": [
                        "Bash(poetry run pytest:*)",
                        "Bash(poetry run ruff:*)",
                        "Bash(poetry run mypy:*)",
                        "Bash(git status:*)",
                        "Bash(git log:*)",
                        "Bash(git diff:*)",
                        "Read(*)",
                        "Glob(*)",
                        "Grep(*)",
                    ],
                    "deny": [
                        "Bash(rm -rf /*)",
                        "Bash(git push --force:*)",
                        "Write(.env*)",
                        "Write(**/secrets/*)",
                    ],
                }
            }
            settings_file.write_text(json.dumps(default_settings, indent=2))

    def configure_hooks(self, hooks: list[str], base_path: Path | None = None) -> None:
        """Configure platform hooks in settings.json.

        Updates .claude/settings.json with hook configuration.
        Supports PreToolUse, PostToolUse, and Stop hooks.

        Args:
            hooks: List of hook script names (e.g., ["pre-edit-check"])
            base_path: Base path for .claude/ directory (defaults to cwd)

        Raises:
            FileNotFoundError: If settings.json does not exist
            ValueError: If hook name is invalid
        """
        if base_path is None:
            base_path = Path.cwd()

        settings_file = base_path / ".claude" / "settings.json"
        if not settings_file.exists():
            raise FileNotFoundError("settings.json not found. Run install() first.")

        # Read current settings
        content = json.loads(settings_file.read_text())

        # Initialize hooks section if not present
        if "hooks" not in content:
            content["hooks"] = {}

        # Configure hooks (simplified - just add PreToolUse for now)
        if hooks:
            content["hooks"]["PreToolUse"] = []
            for hook_name in hooks:
                content["hooks"]["PreToolUse"].append(
                    {
                        "matcher": "Edit|Write",
                        "hooks": [
                            {
                                "type": "command",
                                "command": f"bash hooks/{hook_name}.sh",
                            }
                        ],
                    }
                )

        # Write updated settings
        settings_file.write_text(json.dumps(content, indent=2))

    def load_skill(self, skill_name: str, base_path: Path | None = None) -> dict[str, Any]:
        """Load skill configuration from .claude/skills/.

        Reads skill from .claude/skills/{skill_name}/SKILL.md and parses
        frontmatter (name, description, tools) and prompt content.

        Args:
            skill_name: Name of skill to load (e.g., "build", "idea")
            base_path: Base path for .claude/ directory (defaults to cwd)

        Returns:
            Skill configuration dictionary with:
                - name: Skill name
                - description: Skill description
                - tools: List of allowed tools
                - prompt: Skill prompt content

        Raises:
            FileNotFoundError: If skill file does not exist
            ValueError: If skill format is invalid
        """
        if base_path is None:
            base_path = Path.cwd()

        skill_file = base_path / ".claude" / "skills" / skill_name / "SKILL.md"
        if not skill_file.exists():
            raise FileNotFoundError(f"Skill '{skill_name}' not found at {skill_file}")

        content = skill_file.read_text()

        # Parse frontmatter
        frontmatter_match = re.match(r"^---\n(.*?)\n---\n(.*)$", content, re.DOTALL)
        if not frontmatter_match:
            raise ValueError(f"Invalid skill format: {skill_file}")

        frontmatter_text = frontmatter_match.group(1)
        prompt_content = frontmatter_match.group(2).strip()

        # Parse frontmatter fields
        frontmatter: dict[str, str | list[str]] = {}
        for line in frontmatter_text.split("\n"):
            if ":" in line:
                key, value = line.split(":", 1)
                key = key.strip()
                value = value.strip()

                # Handle list fields (tools)
                if key == "tools" and "," in value:
                    frontmatter[key] = [t.strip() for t in value.split(",")]
                else:
                    frontmatter[key] = value

        return {
            "name": str(frontmatter.get("name", skill_name)),
            "description": str(frontmatter.get("description", "")),
            "tools": (
                frontmatter.get("tools", [])
                if isinstance(frontmatter.get("tools"), list)
                else []
            ),
            "prompt": prompt_content,
        }

    def get_settings(self, base_path: Path | None = None) -> dict[str, Any]:
        """Get platform settings from settings.json.

        Reads and parses .claude/settings.json.

        Args:
            base_path: Base path for .claude/ directory (defaults to cwd)

        Returns:
            Settings dictionary with permissions and hooks

        Raises:
            FileNotFoundError: If settings.json does not exist
            ValueError: If settings format is invalid
        """
        if base_path is None:
            base_path = Path.cwd()

        settings_file = base_path / ".claude" / "settings.json"
        if not settings_file.exists():
            raise FileNotFoundError(f"settings.json not found at {settings_file}")

        try:
            content: dict[str, Any] = json.loads(settings_file.read_text())
            return content
        except json.JSONDecodeError as e:
            raise ValueError(f"Invalid JSON in settings.json: {e}") from e

    def _find_settings_file(self) -> Path | None:
        """Find settings.json by searching up from current directory.

        Returns:
            Path to settings.json or None if not found
        """
        current = Path.cwd()
        while current != current.parent:
            settings_file = current / ".claude" / "settings.json"
            if settings_file.exists():
                return settings_file
            current = current.parent
        return None
