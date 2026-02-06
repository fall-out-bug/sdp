"""Tests for sdp.github.status_mapper."""

import pytest

from sdp.github.status_mapper import StatusMapper


def test_ws_to_github_label_backlog() -> None:
    """Test mapping backlog status to GitHub label."""
    result = StatusMapper.ws_to_github_label("backlog")
    assert result == "status/backlog"


def test_ws_to_github_label_active() -> None:
    """Test mapping active status to GitHub label."""
    result = StatusMapper.ws_to_github_label("active")
    assert result == "status/in-progress"


def test_ws_to_github_label_completed() -> None:
    """Test mapping completed status to GitHub label."""
    result = StatusMapper.ws_to_github_label("completed")
    assert result == "status/completed"


def test_ws_to_github_label_blocked() -> None:
    """Test mapping blocked status to GitHub label."""
    result = StatusMapper.ws_to_github_label("blocked")
    assert result == "status/blocked"


def test_ws_to_github_label_unknown_raises() -> None:
    """Test unknown WS status raises ValueError."""
    with pytest.raises(ValueError, match="Unknown WS status: invalid"):
        StatusMapper.ws_to_github_label("invalid")


def test_ws_to_github_state_backlog() -> None:
    """Test mapping backlog status to GitHub state."""
    result = StatusMapper.ws_to_github_state("backlog")
    assert result == "open"


def test_ws_to_github_state_active() -> None:
    """Test mapping active status to GitHub state."""
    result = StatusMapper.ws_to_github_state("active")
    assert result == "open"


def test_ws_to_github_state_completed() -> None:
    """Test mapping completed status to GitHub state."""
    result = StatusMapper.ws_to_github_state("completed")
    assert result == "closed"


def test_ws_to_github_state_blocked() -> None:
    """Test mapping blocked status to GitHub state."""
    result = StatusMapper.ws_to_github_state("blocked")
    assert result == "open"


def test_ws_to_github_state_unknown_raises() -> None:
    """Test unknown WS status raises ValueError."""
    with pytest.raises(ValueError, match="Unknown WS status: invalid"):
        StatusMapper.ws_to_github_state("invalid")


def test_github_label_to_ws_backlog() -> None:
    """Test mapping GitHub backlog label to WS status."""
    result = StatusMapper.github_label_to_ws("status/backlog")
    assert result == "backlog"


def test_github_label_to_ws_in_progress() -> None:
    """Test mapping GitHub in-progress label to WS status."""
    result = StatusMapper.github_label_to_ws("status/in-progress")
    assert result == "active"


def test_github_label_to_ws_completed() -> None:
    """Test mapping GitHub completed label to WS status."""
    result = StatusMapper.github_label_to_ws("status/completed")
    assert result == "completed"


def test_github_label_to_ws_blocked() -> None:
    """Test mapping GitHub blocked label to WS status."""
    result = StatusMapper.github_label_to_ws("status/blocked")
    assert result == "blocked"


def test_github_label_to_ws_non_status_label() -> None:
    """Test non-status labels return None."""
    result = StatusMapper.github_label_to_ws("bug")
    assert result is None
    
    result = StatusMapper.github_label_to_ws("enhancement")
    assert result is None


def test_detect_conflict_no_conflict() -> None:
    """Test detect_conflict returns False when WS and GitHub match."""
    # Backlog → status/backlog + open
    result = StatusMapper.detect_conflict("backlog", "status/backlog", "open")
    assert result is False
    
    # Active → status/in-progress + open
    result = StatusMapper.detect_conflict("active", "status/in-progress", "open")
    assert result is False
    
    # Completed → status/completed + closed
    result = StatusMapper.detect_conflict("completed", "status/completed", "closed")
    assert result is False


def test_detect_conflict_label_mismatch() -> None:
    """Test detect_conflict returns True when GitHub label differs."""
    # WS says backlog, GitHub says in-progress
    result = StatusMapper.detect_conflict("backlog", "status/in-progress", "open")
    assert result is True


def test_detect_conflict_state_mismatch() -> None:
    """Test detect_conflict returns True when GitHub state differs."""
    # WS says backlog (should be open), GitHub is closed
    result = StatusMapper.detect_conflict("backlog", "status/backlog", "closed")
    assert result is True


def test_detect_conflict_both_mismatch() -> None:
    """Test detect_conflict returns True when both label and state differ."""
    # WS says backlog, GitHub says completed + closed
    result = StatusMapper.detect_conflict("backlog", "status/completed", "closed")
    assert result is True


def test_detect_conflict_unknown_ws_status() -> None:
    """Test detect_conflict raises ValueError for unknown WS status."""
    with pytest.raises(ValueError, match="Unknown WS status"):
        StatusMapper.detect_conflict("invalid", "status/backlog", "open")
