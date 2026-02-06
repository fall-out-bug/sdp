"""Shared pytest fixtures."""

import sys
from pathlib import Path
from typing import Any

import pytest

# Add tests directory to path for fixture imports (before importing fixtures)
sys.path.insert(0, str(Path(__file__).parent))

from fixtures.github_responses import (  # noqa: E402 # type: ignore
    create_mock_github_client,
)


@pytest.fixture
def mock_github() -> tuple[Any, Any]:
    """Mock GitHub client and repo.

    Returns:
        Tuple of (client, repo) mocks

    """
    return create_mock_github_client()


@pytest.fixture
def temp_ws_file(tmp_path: Path) -> Path:
    """Create temporary WS file for testing.

    Args:
        tmp_path: Temporary directory

    Returns:
        Path to temporary WS file

    """
    ws_file = tmp_path / "WS-160-01-test.md"
    ws_file.write_text(
        """---
ws_id: WS-160-01
feature: F160
status: backlog
size: SMALL
github_issue: null
---

## WS-160-01: Test Workstream

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Test functionality

**Acceptance Criteria:**
- [ ] AC1: Test criterion 1
- [ ] AC2: Test criterion 2

### Зависимость

Независимый
""",
        encoding="utf-8",
    )
    return ws_file


@pytest.fixture
def mock_env(monkeypatch):  # type: ignore
    """Mock environment variables.

    Args:
        monkeypatch: Pytest monkeypatch fixture

    """
    monkeypatch.setenv("GITHUB_TOKEN", "test_token")
    monkeypatch.setenv("GITHUB_REPO", "test-org/test-repo")
    monkeypatch.setenv("GITHUB_ORG", "test-org")


@pytest.fixture(autouse=True)
def use_mock_beads(monkeypatch):  # type: ignore
    """Force mock Beads in all tests by default.

    Args:
        monkeypatch: Pytest monkeypatch fixture

    This ensures tests don't accidentally try to use real Beads CLI.
    Override with real_beads fixture when testing real integration.
    """
    monkeypatch.setenv("BEADS_USE_MOCK", "true")


@pytest.fixture
def real_beads(monkeypatch):  # type: ignore
    """Use real Beads for specific tests.

    Args:
        monkeypatch: Pytest monkeypatch fixture

    Apply this fixture to tests that need real Beads CLI integration.
    """
    monkeypatch.setenv("BEADS_USE_MOCK", "false")
