"""Tests for feature orchestrator edge cases."""

import pytest
from unittest.mock import Mock, patch

from sdp.unified.feature.orchestrator import FeatureOrchestrator
from sdp.unified.feature.models import (
    FeatureExecution,
    FeaturePhase,
    SkipFlags,
    StepResult,
    SkillResult,
)


@pytest.fixture
def orchestrator():
    """Create orchestrator instance."""
    return FeatureOrchestrator()


@pytest.fixture
def execution():
    """Create feature execution instance."""
    return FeatureExecution(
        feature_id="F001",
        feature_name="Test Feature",
        skip_flags=SkipFlags(),
    )


class TestGenerateProgressiveMenu:
    """Tests for generate_progressive_menu method."""

    def test_generate_menu_step_1(self, orchestrator, execution):
        """Test generates menu for step 1."""
        menu = orchestrator.generate_progressive_menu(execution, 1)

        assert "Test Feature" in menu
        assert "Requirements Gathering" in menu
        assert "Progress:" in menu
        assert "Step 1 of 5" in menu

    def test_generate_menu_with_skip_requirements(self, orchestrator):
        """Test shows skip warning for requirements."""
        execution = FeatureExecution(
            feature_id="F001",
            feature_name="Test Feature",
            skip_flags=SkipFlags(skip_requirements=True),
        )

        menu = orchestrator.generate_progressive_menu(execution, 1)

        assert "[SKIP] Requirements phase will be skipped" in menu

    def test_generate_menu_with_skip_architecture(self, orchestrator):
        """Test shows skip warning for architecture."""
        execution = FeatureExecution(
            feature_id="F001",
            feature_name="Test Feature",
            skip_flags=SkipFlags(skip_architecture=True),
        )

        menu = orchestrator.generate_progressive_menu(execution, 2)

        assert "[SKIP] Architecture phase will be skipped" in menu

    def test_generate_menu_progress_50_percent(self, orchestrator, execution):
        """Test calculates 50% progress correctly."""
        execution.completed_phases = [
            FeaturePhase.REQUIREMENTS,
            FeaturePhase.ARCHITECTURE,
        ]

        menu = orchestrator.generate_progressive_menu(execution, 3)

        assert "40%" in menu  # 2/5 * 100 = 40%

    def test_generate_menu_highlights_current_step(self, orchestrator, execution):
        """Test highlights current step."""
        menu = orchestrator.generate_progressive_menu(execution, 2)

        assert "(current)" in menu


class TestExecuteStep:
    """Tests for execute_step method."""

    def test_execute_step_1_skip_requirements(self, orchestrator):
        """Test skips step 1 when skip_requirements is set."""
        execution = FeatureExecution(
            feature_id="F001",
            feature_name="Test Feature",
            skip_flags=SkipFlags(skip_requirements=True),
        )

        result = orchestrator.execute_step(execution, 1)

        assert result.success is True
        assert result.phase == FeaturePhase.ARCHITECTURE

    @patch('sdp.unified.feature.orchestrator.invoke_idea_skill')
    def test_execute_step_1_success(self, mock_invoke, orchestrator, execution):
        """Test executes step 1 successfully."""
        mock_invoke.return_value = SkillResult(success=True, artifacts={})

        result = orchestrator.execute_step(execution, 1)

        assert result.success is True
        assert result.phase == FeaturePhase.REQUIREMENTS
        assert FeaturePhase.REQUIREMENTS in execution.completed_phases

    @patch('sdp.unified.feature.orchestrator.invoke_idea_skill')
    def test_execute_step_1_failure(self, mock_invoke, orchestrator, execution):
        """Test handles step 1 failure."""
        mock_invoke.return_value = SkillResult(success=False, error="Test error")

        result = orchestrator.execute_step(execution, 1)

        assert result.success is False
        assert result.phase is None

    @patch('sdp.unified.feature.orchestrator.invoke_design_skill')
    def test_execute_step_2_success(self, mock_invoke, orchestrator, execution):
        """Test executes step 2 successfully."""
        mock_invoke.return_value = SkillResult(success=True, artifacts={})

        result = orchestrator.execute_step(execution, 2)

        assert result.success is True
        assert result.phase == FeaturePhase.ARCHITECTURE
        assert FeaturePhase.ARCHITECTURE in execution.completed_phases

    @patch('sdp.unified.feature.orchestrator.invoke_design_skill')
    def test_execute_step_2_failure(self, mock_invoke, orchestrator, execution):
        """Test handles step 2 failure."""
        mock_invoke.return_value = SkillResult(success=False, error="Design error")

        result = orchestrator.execute_step(execution, 2)

        assert result.success is False
        assert result.phase is None

    @patch('sdp.unified.feature.orchestrator.invoke_oneshot_skill')
    def test_execute_step_3_success(self, mock_invoke, orchestrator, execution):
        """Test executes step 3 successfully."""
        mock_invoke.return_value = SkillResult(success=True, artifacts={})

        result = orchestrator.execute_step(execution, 3)

        assert result.success is True
        assert result.phase == FeaturePhase.EXECUTION
        assert FeaturePhase.EXECUTION in execution.completed_phases

    @patch('sdp.unified.feature.orchestrator.invoke_oneshot_skill')
    def test_execute_step_3_failure(self, mock_invoke, orchestrator, execution):
        """Test handles step 3 failure."""
        mock_invoke.return_value = SkillResult(success=False, error="Execution error")

        result = orchestrator.execute_step(execution, 3)

        assert result.success is False
        assert result.phase is None

    @patch('sdp.unified.feature.orchestrator.invoke_idea_skill')
    def test_execute_step_exception_handling(self, mock_invoke, orchestrator, execution):
        """Test handles exceptions during step execution."""
        mock_invoke.side_effect = RuntimeError("Unexpected error")

        result = orchestrator.execute_step(execution, 1)

        assert result.success is False
        assert result.phase is None


class TestExecuteFeature:
    """Tests for execute_feature method."""

    def test_execute_feature_all_skipped(self, orchestrator):
        """Test executes feature with all phases skipped."""
        execution = FeatureExecution(
            feature_id="F001",
            feature_name="Test Feature",
            skip_flags=SkipFlags(
                skip_requirements=True,
                skip_architecture=True,
                skip_uat=True,
            ),
        )

        result = orchestrator.execute_feature(execution)

        assert result.success is True
        assert result.phase == FeaturePhase.EXECUTION
