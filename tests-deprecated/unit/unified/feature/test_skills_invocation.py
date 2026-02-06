"""Tests for unified feature skills invocation."""

import pytest
from unittest.mock import patch

from sdp.unified.feature.skills import (
    invoke_idea_skill,
    invoke_design_skill,
    invoke_oneshot_skill,
)


class TestInvokeIdeaSkill:
    """Tests for invoke_idea_skill function."""

    def test_invoke_idea_skill_success(self):
        """Test successful invocation of @idea skill."""
        result = invoke_idea_skill("F001", "Test Feature")

        assert result.success is True
        assert "requirements" in result.artifacts
        assert "intent" in result.artifacts
        assert "docs/drafts/F001-requirements.md" in result.artifacts["requirements"]

    @patch('sdp.unified.feature.skills.logger')
    def test_invoke_idea_skill_logs_info(self, mock_logger):
        """Test logs info message on invocation."""
        invoke_idea_skill("F001", "Test Feature")

        mock_logger.info.assert_called_once()
        assert "F001" in mock_logger.info.call_args[0][0]

    @patch('sdp.unified.feature.skills.logger')
    def test_invoke_idea_skill_exception_handling(self, mock_logger):
        """Test handles exceptions and returns error result."""
        # Patch logger.info to raise exception
        mock_logger.info.side_effect = RuntimeError("Test error")

        result = invoke_idea_skill("F001", "Test Feature")

        assert result.success is False
        assert result.error == "Test error"
        mock_logger.error.assert_called_once()


class TestInvokeDesignSkill:
    """Tests for invoke_design_skill function."""

    def test_invoke_design_skill_success(self):
        """Test successful invocation of @design skill."""
        result = invoke_design_skill("F001", "docs/drafts/F001-requirements.md")

        assert result.success is True
        assert "architecture" in result.artifacts
        assert "workstreams" in result.artifacts
        assert "docs/drafts/F001-architecture.md" in result.artifacts["architecture"]

    @patch('sdp.unified.feature.skills.logger')
    def test_invoke_design_skill_logs_info(self, mock_logger):
        """Test logs info message on invocation."""
        invoke_design_skill("F001", "docs/drafts/F001-requirements.md")

        mock_logger.info.assert_called_once()
        assert "F001" in mock_logger.info.call_args[0][0]

    @patch('sdp.unified.feature.skills.logger')
    def test_invoke_design_skill_exception_handling(self, mock_logger):
        """Test handles exceptions and returns error result."""
        # Patch logger.info to raise exception
        mock_logger.info.side_effect = RuntimeError("Design error")

        result = invoke_design_skill("F001", "docs/drafts/F001-requirements.md")

        assert result.success is False
        assert result.error == "Design error"
        mock_logger.error.assert_called_once()


class TestInvokeOneshotSkill:
    """Tests for invoke_oneshot_skill function."""

    def test_invoke_oneshot_skill_success(self):
        """Test successful invocation of @oneshot skill."""
        result = invoke_oneshot_skill("F001", "docs/drafts/F001-architecture.md")

        assert result.success is True
        assert "execution_plan" in result.artifacts
        assert "checkpoint" in result.artifacts
        assert "docs/plans/F001-execution.md" in result.artifacts["execution_plan"]

    @patch('sdp.unified.feature.skills.logger')
    def test_invoke_oneshot_skill_logs_info(self, mock_logger):
        """Test logs info message on invocation."""
        invoke_oneshot_skill("F001", "docs/drafts/F001-architecture.md")

        mock_logger.info.assert_called_once()
        assert "F001" in mock_logger.info.call_args[0][0]

    @patch('sdp.unified.feature.skills.logger')
    def test_invoke_oneshot_skill_exception_handling(self, mock_logger):
        """Test handles exceptions and returns error result."""
        # Patch logger.info to raise exception
        mock_logger.info.side_effect = RuntimeError("Oneshot error")

        result = invoke_oneshot_skill("F001", "docs/drafts/F001-architecture.md")

        assert result.success is False
        assert result.error == "Oneshot error"
        mock_logger.error.assert_called_once()
