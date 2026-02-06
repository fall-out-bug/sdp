"""Tests for metrics CLI commands."""
from pathlib import Path
from unittest.mock import MagicMock, patch

import pytest
from click.testing import CliRunner

from sdp.cli.metrics import metrics


@pytest.fixture
def runner():
    """Create CLI runner."""
    return CliRunner()


def test_metrics_group_help(runner):
    """Test metrics group help text."""
    result = runner.invoke(metrics, ["--help"])
    assert result.exit_code == 0
    assert "Metrics and monitoring commands" in result.output
    assert "escalations" in result.output


def test_escalations_default_options(runner):
    """Test escalations command with default options."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        # Mock return values
        mock_store.get_escalation_count.return_value = 5
        mock_store.get_escalation_rate.return_value = 0.25  # 25%
        mock_store.get_average_attempts.return_value = 2.5
        mock_store.get_top_escalating_ws.return_value = [
            ("00-001-01", 3),
            ("00-001-02", 2),
        ]
        
        result = runner.invoke(metrics, ["escalations"])
        
        assert result.exit_code == 0
        assert "Escalation Metrics (last 7 days)" in result.output
        assert "Total Escalations: 5" in result.output
        assert "Escalation Rate: 25.0%" in result.output
        assert "5/20 builds" in result.output
        assert "Avg Attempts Before Escalation: 2.5" in result.output
        assert "Top 2 Escalating Workstreams:" in result.output
        assert "00-001-01: 3 escalations" in result.output
        assert "00-001-02: 2 escalations" in result.output
        assert "ALERT" in result.output  # 25% > 20% threshold


def test_escalations_with_tier_filter(runner):
    """Test escalations command filtered by tier."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 3
        mock_store.get_escalation_rate.return_value = 0.15
        mock_store.get_average_attempts.return_value = 3.0
        mock_store.get_top_escalating_ws.return_value = []
        
        result = runner.invoke(metrics, ["escalations", "--tier", "T2"])
        
        assert result.exit_code == 0
        assert "Total Escalations: 3" in result.output
        assert "Escalation Rate: 15.0%" in result.output
        
        # Verify tier filter passed to store
        mock_store.get_escalation_count.assert_called_once_with(tier="T2", days=7)
        mock_store.get_escalation_rate.assert_called_once_with(tier="T2", days=7, total_builds=20)


def test_escalations_with_custom_days(runner):
    """Test escalations command with custom time window."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 10
        mock_store.get_escalation_rate.return_value = 0.10
        mock_store.get_average_attempts.return_value = 2.0
        mock_store.get_top_escalating_ws.return_value = []
        
        result = runner.invoke(metrics, ["escalations", "--days", "30"])
        
        assert result.exit_code == 0
        assert "Escalation Metrics (last 30 days)" in result.output
        assert "Total Escalations: 10" in result.output
        
        # Verify days parameter passed
        mock_store.get_escalation_count.assert_called_once_with(tier=None, days=30)


def test_escalations_with_custom_top(runner):
    """Test escalations command with custom top N limit."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 8
        mock_store.get_escalation_rate.return_value = 0.10
        mock_store.get_average_attempts.return_value = 2.5
        mock_store.get_top_escalating_ws.return_value = [
            ("00-001-01", 5),
            ("00-001-02", 3),
            ("00-001-03", 2),
        ]
        
        result = runner.invoke(metrics, ["escalations", "--top", "3"])
        
        assert result.exit_code == 0
        assert "Top 3 Escalating Workstreams:" in result.output
        
        # Verify top parameter passed
        mock_store.get_top_escalating_ws.assert_called_once_with(limit=3, days=7)


def test_escalations_with_custom_storage_path(runner, tmp_path):
    """Test escalations command with custom storage path."""
    storage_path = tmp_path / "custom_metrics.json"
    
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 0
        mock_store.get_escalation_rate.return_value = 0.0
        mock_store.get_average_attempts.return_value = 0.0
        mock_store.get_top_escalating_ws.return_value = []
        
        result = runner.invoke(metrics, ["escalations", "--storage", str(storage_path)])
        
        assert result.exit_code == 0
        
        # Verify custom storage path used
        mock_store_class.assert_called_once_with(storage_path)


def test_escalations_with_custom_total_builds(runner):
    """Test escalations command with custom total builds."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 15
        mock_store.get_escalation_rate.return_value = 0.30  # 30%
        mock_store.get_average_attempts.return_value = 2.0
        mock_store.get_top_escalating_ws.return_value = []
        
        result = runner.invoke(metrics, ["escalations", "--total-builds", "50"])
        
        assert result.exit_code == 0
        assert "15/50 builds" in result.output
        
        # Verify total_builds parameter passed
        mock_store.get_escalation_rate.assert_called_once_with(tier=None, days=7, total_builds=50)


def test_escalations_no_average_attempts(runner):
    """Test escalations command when average attempts is zero."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 0
        mock_store.get_escalation_rate.return_value = 0.0
        mock_store.get_average_attempts.return_value = 0.0  # No attempts
        mock_store.get_top_escalating_ws.return_value = []
        
        result = runner.invoke(metrics, ["escalations"])
        
        assert result.exit_code == 0
        # Should not show avg attempts when zero
        assert "Avg Attempts" not in result.output


def test_escalations_no_top_workstreams(runner):
    """Test escalations command with no top workstreams."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 2
        mock_store.get_escalation_rate.return_value = 0.10
        mock_store.get_average_attempts.return_value = 1.5
        mock_store.get_top_escalating_ws.return_value = []  # No top workstreams
        
        result = runner.invoke(metrics, ["escalations"])
        
        assert result.exit_code == 0
        # Should not show top workstreams section
        assert "Top" not in result.output or "Escalating Workstreams:" not in result.output


def test_escalations_high_rate_alert(runner):
    """Test escalations command shows alert for high rate."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 10
        mock_store.get_escalation_rate.return_value = 0.50  # 50% - very high!
        mock_store.get_average_attempts.return_value = 2.0
        mock_store.get_top_escalating_ws.return_value = []
        
        result = runner.invoke(metrics, ["escalations"])
        
        assert result.exit_code == 0
        assert "⚠️  ALERT: High escalation rate" in result.output
        assert "50.0% > 20.0%" in result.output
        assert "Consider reviewing workstream quality" in result.output


def test_escalations_low_rate_no_alert(runner):
    """Test escalations command with no alert for low rate."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 2
        mock_store.get_escalation_rate.return_value = 0.10  # 10% - acceptable
        mock_store.get_average_attempts.return_value = 2.0
        mock_store.get_top_escalating_ws.return_value = []
        
        result = runner.invoke(metrics, ["escalations"])
        
        assert result.exit_code == 0
        # Should not show alert
        assert "ALERT" not in result.output


def test_escalations_tier_t3(runner):
    """Test escalations command with T3 tier filter."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 1
        mock_store.get_escalation_rate.return_value = 0.05
        mock_store.get_average_attempts.return_value = 4.0
        mock_store.get_top_escalating_ws.return_value = []
        
        result = runner.invoke(metrics, ["escalations", "--tier", "T3"])
        
        assert result.exit_code == 0
        assert "Total Escalations: 1" in result.output
        
        # Verify T3 tier filter
        mock_store.get_escalation_count.assert_called_once_with(tier="T3", days=7)


def test_escalations_combined_options(runner):
    """Test escalations command with multiple options combined."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 8
        mock_store.get_escalation_rate.return_value = 0.16
        mock_store.get_average_attempts.return_value = 3.2
        mock_store.get_top_escalating_ws.return_value = [
            ("00-001-01", 4),
            ("00-001-02", 3),
            ("00-001-03", 1),
        ]
        
        result = runner.invoke(metrics, [
            "escalations",
            "--tier", "T2",
            "--days", "14",
            "--top", "5",
            "--total-builds", "50"
        ])
        
        assert result.exit_code == 0
        assert "Escalation Metrics (last 14 days)" in result.output
        assert "8/50 builds" in result.output
        
        # Verify all parameters passed correctly
        mock_store.get_escalation_count.assert_called_once_with(tier="T2", days=14)
        mock_store.get_escalation_rate.assert_called_once_with(tier="T2", days=14, total_builds=50)
        mock_store.get_average_attempts.assert_called_once_with(tier="T2", days=14)
        mock_store.get_top_escalating_ws.assert_called_once_with(limit=5, days=14)


def test_escalations_formatting(runner):
    """Test escalations command output formatting."""
    with patch("sdp.core.escalation_metrics.EscalationMetricsStore") as mock_store_class:
        mock_store = MagicMock()
        mock_store_class.return_value = mock_store
        
        mock_store.get_escalation_count.return_value = 7
        mock_store.get_escalation_rate.return_value = 0.175  # 17.5%
        mock_store.get_average_attempts.return_value = 2.75
        mock_store.get_top_escalating_ws.return_value = [
            ("00-001-01", 5),
        ]
        
        result = runner.invoke(metrics, ["escalations"])
        
        assert result.exit_code == 0
        # Check percentage formatting
        assert "17.5%" in result.output
        # Check average formatting
        assert "2.8" in result.output or "2.7" in result.output  # Rounded to 1 decimal
