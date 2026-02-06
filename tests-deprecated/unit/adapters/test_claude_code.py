"""Tests for Claude Code platform adapter."""

import json
import pytest
from pathlib import Path

from sdp.adapters.claude_code import ClaudeCodeAdapter


@pytest.fixture
def claude_adapter() -> ClaudeCodeAdapter:
    """Create Claude Code adapter instance."""
    return ClaudeCodeAdapter()


class TestClaudeCodeInstall:
    """Test Claude Code adapter install."""

    def test_install_creates_directory_structure(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify install creates .claude structure."""
        claude_adapter.install(tmp_path)

        assert (tmp_path / ".claude").exists()
        assert (tmp_path / ".claude" / "skills").exists()
        assert (tmp_path / ".claude" / "agents").exists()
        assert (tmp_path / ".claude" / "settings.json").exists()

    def test_install_creates_default_settings(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify install creates default settings.json when missing."""
        claude_adapter.install(tmp_path)

        content = json.loads((tmp_path / ".claude" / "settings.json").read_text())
        assert "permissions" in content
        assert "allow" in content["permissions"]
        assert "deny" in content["permissions"]

    def test_install_idempotent(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify install does not overwrite existing settings."""
        claude_adapter.install(tmp_path)
        settings_file = tmp_path / ".claude" / "settings.json"
        original_content = settings_file.read_text()

        claude_adapter.install(tmp_path)

        assert settings_file.read_text() == original_content


class TestClaudeCodeConfigureHooks:
    """Test Claude Code hook configuration."""

    def test_configure_hooks_requires_install(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify configure_hooks raises when settings.json missing."""
        with pytest.raises(FileNotFoundError, match="settings.json not found"):
            claude_adapter.configure_hooks(["pre-commit"], base_path=tmp_path)

    def test_configure_hooks_updates_settings(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify configure_hooks adds hooks to settings."""
        claude_adapter.install(tmp_path)
        claude_adapter.configure_hooks(["pre-edit-check"], base_path=tmp_path)

        content = json.loads((tmp_path / ".claude" / "settings.json").read_text())
        assert "hooks" in content
        assert "PreToolUse" in content["hooks"]


class TestClaudeCodeLoadSkill:
    """Test Claude Code skill loading."""

    def test_load_skill_raises_when_missing(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify load_skill raises when skill file does not exist."""
        claude_adapter.install(tmp_path)

        with pytest.raises(FileNotFoundError, match="Skill 'build' not found"):
            claude_adapter.load_skill("build", base_path=tmp_path)

    def test_load_skill_parses_valid_skill(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify load_skill parses skill with frontmatter."""
        (tmp_path / ".claude" / "skills" / "build").mkdir(parents=True)
        skill_file = tmp_path / ".claude" / "skills" / "build" / "SKILL.md"
        skill_file.write_text(
            "---\nname: build\ndescription: Execute workstream\ntools: Read, Write\n---\n\n# Build\n\nExecute WS."
        )

        result = claude_adapter.load_skill("build", base_path=tmp_path)

        assert result["name"] == "build"
        assert result["description"] == "Execute workstream"
        assert "Execute WS" in result["prompt"]

    def test_load_skill_raises_invalid_format(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify load_skill raises for invalid frontmatter format."""
        (tmp_path / ".claude" / "skills" / "bad").mkdir(parents=True)
        (tmp_path / ".claude" / "skills" / "bad" / "SKILL.md").write_text("no frontmatter here")

        with pytest.raises(ValueError, match="Invalid skill format"):
            claude_adapter.load_skill("bad", base_path=tmp_path)


class TestClaudeCodeGetSettings:
    """Test Claude Code settings retrieval."""

    def test_get_settings_raises_when_missing(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify get_settings raises when settings.json missing."""
        with pytest.raises(FileNotFoundError, match="settings.json not found"):
            claude_adapter.get_settings(base_path=tmp_path)

    def test_get_settings_returns_parsed_json(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify get_settings returns parsed settings."""
        claude_adapter.install(tmp_path)

        result = claude_adapter.get_settings(base_path=tmp_path)

        assert isinstance(result, dict)
        assert "permissions" in result

    def test_get_settings_raises_invalid_json(self, claude_adapter: ClaudeCodeAdapter, tmp_path: Path) -> None:
        """Verify get_settings raises for invalid JSON."""
        (tmp_path / ".claude").mkdir()
        (tmp_path / ".claude" / "settings.json").write_text("{invalid json}")

        with pytest.raises(ValueError, match="Invalid JSON"):
            claude_adapter.get_settings(base_path=tmp_path)
