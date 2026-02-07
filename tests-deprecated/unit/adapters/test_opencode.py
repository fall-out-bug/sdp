"""Tests for OpenCode platform adapter."""

import json
import pytest
from pathlib import Path

from sdp.adapters.opencode import OpenCodeAdapter


@pytest.fixture
def opencode_adapter() -> OpenCodeAdapter:
    """Create OpenCode adapter instance."""
    return OpenCodeAdapter()


class TestOpenCodeInstall:
    """Test OpenCode adapter install."""

    def test_install_creates_directory_structure(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify install creates .opencode structure."""
        opencode_adapter.install(tmp_path)

        assert (tmp_path / ".opencode").exists()
        assert (tmp_path / ".opencode" / "plugin").exists()
        assert (tmp_path / ".opencode" / "skills").exists()
        assert (tmp_path / ".opencode" / "plugin" / "sdp.js").exists()

    def test_install_creates_plugin_wrapper(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify install creates sdp.js plugin."""
        opencode_adapter.install(tmp_path)

        js_content = (tmp_path / ".opencode" / "plugin" / "sdp.js").read_text()
        assert "sdp" in js_content
        assert "build" in js_content


class TestOpenCodeConfigureHooks:
    """Test OpenCode hook configuration."""

    def test_configure_hooks_requires_install(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify configure_hooks raises when .opencode missing."""
        with pytest.raises(FileNotFoundError, match=".opencode directory not found"):
            opencode_adapter.configure_hooks(["pre-commit"], base_path=tmp_path)

    def test_configure_hooks_noop_when_installed(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify configure_hooks does not raise when installed (no-op)."""
        opencode_adapter.install(tmp_path)

        opencode_adapter.configure_hooks(["pre-commit"], base_path=tmp_path)
        # No exception - OpenCode hooks are no-op


class TestOpenCodeLoadSkill:
    """Test OpenCode skill loading."""

    def test_load_skill_raises_when_missing(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify load_skill raises when skill file does not exist."""
        opencode_adapter.install(tmp_path)

        with pytest.raises(FileNotFoundError, match="Skill 'build' not found"):
            opencode_adapter.load_skill("build", base_path=tmp_path)

    def test_load_skill_parses_and_copies(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify load_skill parses skill and copies to XDG."""
        (tmp_path / ".opencode" / "skills" / "build").mkdir(parents=True)
        skill_content = "---\nname: build\ndescription: Build WS\n---\n\n# Build"
        (tmp_path / ".opencode" / "skills" / "build" / "SKILL.md").write_text(skill_content)

        xdg_home = tmp_path / "xdg_config"
        xdg_home.mkdir()

        result = opencode_adapter.load_skill("build", base_path=tmp_path, xdg_config_home=xdg_home)

        assert result["name"] == "build"
        assert result["description"] == "Build WS"
        user_skill = xdg_home / "opencode" / "skills" / "build" / "SKILL.md"
        assert user_skill.exists()
        assert user_skill.read_text() == skill_content


class TestOpenCodeGetSettings:
    """Test OpenCode settings retrieval."""

    def test_get_settings_raises_when_not_installed(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify get_settings raises when .opencode missing."""
        with pytest.raises(FileNotFoundError, match=".opencode directory not found"):
            opencode_adapter.get_settings(base_path=tmp_path)

    def test_get_settings_returns_empty_when_no_config(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify get_settings returns empty dict when opencode.json missing."""
        opencode_adapter.install(tmp_path)

        result = opencode_adapter.get_settings(base_path=tmp_path)

        assert result == {}

    def test_get_settings_returns_parsed_json(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify get_settings returns parsed opencode.json."""
        opencode_adapter.install(tmp_path)
        (tmp_path / ".opencode" / "opencode.json").write_text('{"model": "opencode-beta"}')

        result = opencode_adapter.get_settings(base_path=tmp_path)

        assert result["model"] == "opencode-beta"

    def test_get_settings_raises_invalid_json(self, opencode_adapter: OpenCodeAdapter, tmp_path: Path) -> None:
        """Verify get_settings raises for invalid JSON."""
        opencode_adapter.install(tmp_path)
        (tmp_path / ".opencode" / "opencode.json").write_text("{invalid}")

        with pytest.raises(ValueError, match="Invalid JSON"):
            opencode_adapter.get_settings(base_path=tmp_path)
