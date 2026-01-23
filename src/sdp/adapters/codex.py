"""Codex platform adapter implementation.

This module implements the PlatformAdapter interface for Codex,
managing .codex/ directory structure and user-level skills.
"""

from __future__ import annotations

import re
from pathlib import Path
from typing import Any

import yaml  # type: ignore[import-untyped]

from sdp.adapters.base import PlatformAdapter


class CodexAdapter(PlatformAdapter):
    """Codex platform adapter.

    Project structure:
    - .codex/INSTALL.md (setup instructions)
    - .codex/skills/ (project-level skills)

    User-level structure:
    - ~/.codex/skills/ (persistent skills)
    """

    def install(self, target_dir: Path) -> None:
        """Install Codex directory structure.

        Creates:
        - .codex/ directory
        - .codex/skills/ directory
        - .codex/INSTALL.md with setup instructions

        Args:
            target_dir: Directory where .codex/ should be created
        """
        codex_dir = target_dir / ".codex"
        codex_dir.mkdir(exist_ok=True)

        (codex_dir / "skills").mkdir(exist_ok=True)

        install_file = codex_dir / "INSTALL.md"
        if not install_file.exists():
            install_file.write_text(self._default_install_instructions())

    def configure_hooks(self, hooks: list[str], base_path: Path | None = None) -> None:
        """Configure hooks for Codex.

        Codex does not support hook configuration; this is a no-op.

        Args:
            hooks: List of hook names (ignored)
            base_path: Base path for .codex/ directory (defaults to cwd)

        Raises:
            FileNotFoundError: If .codex directory does not exist
        """
        if base_path is None:
            base_path = Path.cwd()

        codex_dir = base_path / ".codex"
        if not codex_dir.exists():
            raise FileNotFoundError(".codex directory not found. Run install() first.")

        _ = hooks

    def load_skill(
        self,
        skill_name: str,
        base_path: Path | None = None,
        home_path: Path | None = None,
    ) -> dict[str, Any]:
        """Load and copy skill to user-level directory.

        Reads skill from .codex/skills/{skill_name}/SKILL.md and copies it to
        ~/.codex/skills/{skill_name}/SKILL.md.

        Args:
            skill_name: Name of skill to load
            base_path: Base path for .codex/ directory (defaults to cwd)
            home_path: Base path for user home (defaults to Path.home())

        Returns:
            Skill configuration dictionary with name, description, tools, prompt

        Raises:
            FileNotFoundError: If skill file does not exist
            ValueError: If skill format is invalid
        """
        if base_path is None:
            base_path = Path.cwd()
        if home_path is None:
            home_path = Path.home()

        skill_file = base_path / ".codex" / "skills" / skill_name / "SKILL.md"
        if not skill_file.exists():
            raise FileNotFoundError(f"Skill '{skill_name}' not found at {skill_file}")

        content = skill_file.read_text()
        frontmatter, prompt = self._parse_frontmatter(content, skill_name)

        user_skill_dir = home_path / ".codex" / "skills" / skill_name
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
        """Get Codex settings from .codex/config.yaml if present.

        Args:
            base_path: Base path for .codex/ directory (defaults to cwd)

        Returns:
            Parsed settings dictionary or empty dict if missing

        Raises:
            FileNotFoundError: If .codex directory does not exist
            ValueError: If config.yaml is invalid
        """
        if base_path is None:
            base_path = Path.cwd()

        codex_dir = base_path / ".codex"
        if not codex_dir.exists():
            raise FileNotFoundError(".codex directory not found. Run install() first.")

        config_file = codex_dir / "config.yaml"
        if not config_file.exists():
            return {}

        try:
            data = yaml.safe_load(config_file.read_text())
        except yaml.YAMLError as exc:
            raise ValueError(f"Invalid YAML in config.yaml: {exc}") from exc

        if data is None:
            return {}

        if not isinstance(data, dict):
            raise ValueError("config.yaml must define a mapping at top level")

        return data

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

    def _default_install_instructions(self) -> str:
        """Return INSTALL.md content for Codex users."""
        return (
            "# Codex Setup\n\n"
            "This project uses Codex skills for the Spec-Driven Protocol (SDP).\n\n"
            "## Steps\n\n"
            "1. Ensure Codex recognizes this repo by keeping `.codex/INSTALL.md`.\n"
            "2. Project skills live in `.codex/skills/`.\n"
            "3. User-level skills are copied to `~/.codex/skills/` when loaded.\n\n"
            "## Usage\n\n"
            "- Run `sdp` commands normally in this repo.\n"
            "- Skills are loaded from project `.codex/skills/` and copied to your\n"
            "  user directory for persistence.\n"
        )
