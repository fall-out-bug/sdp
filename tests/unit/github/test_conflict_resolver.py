"""Tests for ConflictResolver."""

import pytest

from sdp.github.conflict_resolver import Conflict, ConflictResolver


def test_detect_conflict_when_status_mismatch() -> None:
    resolver = ConflictResolver()
    ws_state = {"status": "backlog"}
    gh_state = {"status": "In Progress"}

    conflict = resolver.detect(ws_state, gh_state)
    assert conflict is not None
    assert conflict.field == "status"


def test_no_conflict_when_status_matches() -> None:
    resolver = ConflictResolver()
    ws_state = {"status": "backlog"}
    gh_state = {"status": "Backlog"}

    conflict = resolver.detect(ws_state, gh_state)
    assert conflict is None


def test_no_conflict_when_both_in_progress() -> None:
    resolver = ConflictResolver()
    ws_state = {"status": "in-progress"}
    gh_state = {"status": "In Progress"}

    conflict = resolver.detect(ws_state, gh_state)
    assert conflict is None


def test_resolve_conflict_ws_wins() -> None:
    resolver = ConflictResolver()
    conflict = Conflict(field="status", ws_value="backlog", gh_value="In Progress")

    resolved = resolver.resolve(conflict)
    assert resolved == "backlog"


def test_resolve_conflict_completed() -> None:
    resolver = ConflictResolver()
    conflict = Conflict(field="status", ws_value="completed", gh_value="Backlog")

    resolved = resolver.resolve(conflict)
    assert resolved == "completed"


def test_no_conflict_with_missing_ws_status() -> None:
    resolver = ConflictResolver()
    ws_state = {}
    gh_state = {"status": "In Progress"}

    conflict = resolver.detect(ws_state, gh_state)
    assert conflict is None


def test_no_conflict_with_missing_gh_status() -> None:
    resolver = ConflictResolver()
    ws_state = {"status": "backlog"}
    gh_state = {}

    conflict = resolver.detect(ws_state, gh_state)
    assert conflict is None
