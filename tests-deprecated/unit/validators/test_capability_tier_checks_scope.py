"""Tests for capability_tier_checks_scope validator."""

from unittest.mock import MagicMock

import pytest

from sdp.validators.capability_tier_checks_scope import _check_scope_tiny


class TestCheckScopeTiny:
    """Tests for _check_scope_tiny."""

    def test_tiny_scope_passes(self) -> None:
        """TINY scope passes."""
        ws = MagicMock()
        ws.size.value = "TINY"
        result = _check_scope_tiny(ws, "")
        assert result.passed
        assert "TINY" in result.message

    def test_large_scope_fails(self) -> None:
        """LARGE scope fails."""
        ws = MagicMock()
        ws.size.value = "LARGE"
        result = _check_scope_tiny(ws, "")
        assert not result.passed
        assert "TINY" in result.message

    def test_xlarge_scope_fails(self) -> None:
        """XLARGE scope fails."""
        ws = MagicMock()
        ws.size.value = "XLARGE"
        result = _check_scope_tiny(ws, "")
        assert not result.passed

    def test_small_scope_fails(self) -> None:
        """SMALL scope fails (T3 requires TINY)."""
        ws = MagicMock()
        ws.size.value = "SMALL"
        result = _check_scope_tiny(ws, "")
        assert not result.passed
        assert "TINY" in result.message
