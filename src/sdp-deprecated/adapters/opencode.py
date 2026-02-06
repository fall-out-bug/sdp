"""OpenCode platform adapter implementation.

This module implements the PlatformAdapter interface for OpenCode,
managing .opencode/ directory structure and JavaScript plugin wrapper.
"""

from __future__ import annotations

import json
import re
from pathlib import Path
from typing import Any

from sdp.adapters.base import PlatformAdapter


class OpenCodeAdapter(PlatformAdapter):
    """OpenCode platform adapter.

    Project structure:
    - .opencode/plugin/sdp.js (JavaScript plugin wrapper)
    - .opencode/skills/ (project-level skills)

    User-level structure:
    - ~/.config/opencode/skills/ (persistent skills via XDG)
    """

    def install(self, target_dir: Path) -> None:
        """Install OpenCode directory structure.

        Creates:
        - .opencode/ directory
        - .opencode/plugin/ directory
        - .opencode/plugin/sdp.js (JavaScript wrapper)
        - .opencode/skills/ directory

        Args:
            target_dir: Directory where .opencode/ should be created
        """
        opencode_dir = target_dir / ".opencode"
        opencode_dir.mkdir(exist_ok=True)

        plugin_dir = opencode_dir / "plugin"
        plugin_dir.mkdir(exist_ok=True)

        (opencode_dir / "skills").mkdir(exist_ok=True)

        js_file = plugin_dir / "sdp.js"
        if not js_file.exists():
            js_file.write_text(self._default_plugin_wrapper())

    def configure_hooks(self, hooks: list[str], base_path: Path | None = None) -> None:
        """Configure hooks for OpenCode.

        OpenCode does not support hook configuration; this is a no-op.

        Args:
            hooks: List of hook names (ignored)
            base_path: Base path for .opencode/ directory (defaults to cwd)

        Raises:
            FileNotFoundError: If .opencode directory does not exist
        """
        if base_path is None:
            base_path = Path.cwd()

        opencode_dir = base_path / ".opencode"
        if not opencode_dir.exists():
            raise FileNotFoundError(".opencode directory not found. Run install() first.")

        _ = hooks

    def load_skill(
        self,
        skill_name: str,
        base_path: Path | None = None,
        xdg_config_home: Path | None = None,
    ) -> dict[str, Any]:
        """Load and copy skill to XDG config directory.

        Reads skill from .opencode/skills/{skill_name}/SKILL.md and copies it to
        ~/.config/opencode/skills/{skill_name}/SKILL.md.

        Args:
            skill_name: Name of skill to load
            base_path: Base path for .opencode/ directory (defaults to cwd)
            xdg_config_home: XDG config home (defaults to ~/.config)

        Returns:
            Skill configuration dictionary with name, description, tools, prompt

        Raises:
            FileNotFoundError: If skill file does not exist
            ValueError: If skill format is invalid
        """
        if base_path is None:
            base_path = Path.cwd()
        if xdg_config_home is None:
            xdg_config_home = Path.home() / ".config"

        skill_file = base_path / ".opencode" / "skills" / skill_name / "SKILL.md"
        if not skill_file.exists():
            raise FileNotFoundError(f"Skill '{skill_name}' not found at {skill_file}")

        content = skill_file.read_text()
        frontmatter, prompt = self._parse_frontmatter(content, skill_name)

        user_skill_dir = xdg_config_home / "opencode" / "skills" / skill_name
        user_skill_dir.mkdir(parents=True, exist_ok=True)
        user_skill_file = user_skill_dir / "SKILL.md"
        user_skill_file.write_text(content)

        return {
            "name": str(frontmatter.get("name", skill_name)),
            "description": str(frontmatter.get("description", "")),
            "tools": (
                frontmatter.get("tools", [])
                if isinstance(frontmatter.get("tools"), list)
                else []
            ),
            "prompt": prompt,
        }

    def get_settings(self, base_path: Path | None = None) -> dict[str, Any]:
        """Get OpenCode settings from .opencode/opencode.json if present.

        Args:
            base_path: Base path for .opencode/ directory (defaults to cwd)

        Returns:
            Parsed settings dictionary or empty dict if missing

        Raises:
            FileNotFoundError: If .opencode directory does not exist
            ValueError: If opencode.json is invalid
        """
        if base_path is None:
            base_path = Path.cwd()

        opencode_dir = base_path / ".opencode"
        if not opencode_dir.exists():
            raise FileNotFoundError(".opencode directory not found. Run install() first.")

        config_file = opencode_dir / "opencode.json"
        if not config_file.exists():
            return {}

        try:
            data: dict[str, Any] = json.loads(config_file.read_text())
            return data
        except json.JSONDecodeError as exc:
            raise ValueError(f"Invalid JSON in opencode.json: {exc}") from exc

    def _parse_frontmatter(
        self, content: str, skill_name: str
    ) -> tuple[dict[str, str | list[str]], str]:
        """Parse SKILL.md frontmatter and prompt content."""
        match = re.match(r"^---\n(.*?)\n---\n(.*)$", content, re.DOTALL)
        if not match:
            raise ValueError(f"Invalid skill format: {skill_name}")

        frontmatter_text = match.group(1)
        prompt_content = match.group(2).strip()

        frontmatter: dict[str, str | list[str]] = {}
        for line in frontmatter_text.split("\n"):
            if ":" in line:
                key, value = line.split(":", 1)
                key = key.strip()
                value = value.strip()

                if key == "tools" and "," in value:
                    frontmatter[key] = [t.strip() for t in value.split(",")]
                else:
                    frontmatter[key] = value

        return frontmatter, prompt_content

    def _default_plugin_wrapper(self) -> str:
        """Return JavaScript plugin wrapper for OpenCode."""
        return """// SDP (Spec-Driven Protocol) Plugin for OpenCode
// Generated by sdp.adapters.opencode.OpenCodeAdapter

export default {
  name: 'sdp',
  version: '0.3.0',
  description: 'Spec-Driven Protocol workstream automation',

  commands: {
    build: async (args) => {
      // Delegate to SDP CLI
      return await executeCommand('sdp build ' + args.join(' '));
    },
    design: async (args) => {
      return await executeCommand('sdp design ' + args.join(' '));
    },
    codereview: async (args) => {
      return await executeCommand('sdp codereview ' + args.join(' '));
    }
  }
};
"""
