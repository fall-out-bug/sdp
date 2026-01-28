"""
Tests for @feature skill orchestrator.

Tests the unified workflow entry point that orchestrates
@idea → @design → @oneshot with progressive menu.
"""

import pytest
from pathlib import Path
from datetime import datetime, timezone
from unittest.mock import Mock, patch

from sdp.unified.feature.orchestrator import FeatureOrchestrator
from sdp.unified.feature.models import FeatureExecution, FeaturePhase, SkipFlags


class TestFeatureOrchestratorInit:
    """Test orchestrator initialization."""

    def test_creates_orchestrator_with_config(self):
        """Should initialize with default configuration."""
        orchestrator = FeatureOrchestrator()

        assert orchestrator is not None
        assert hasattr(orchestrator, 'current_phase')


class TestProgressiveMenu:
    """Test progressive menu generation."""

    def test_generate_initial_menu_shows_step_1(self):
        """Should show Step 1 of 5 for requirements gathering."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            skip_flags=SkipFlags(),
        )

        menu = orchestrator.generate_progressive_menu(execution, step=1)

        assert "Step 1 of 5" in menu
        assert "Requirements Gathering" in menu

    def test_menu_shows_skip_option(self):
        """Should include skip option in menu."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            skip_flags=SkipFlags(skip_requirements=True),
        )

        menu = orchestrator.generate_progressive_menu(execution, step=1)

        assert "[skip]" in menu.lower()

    def test_menu_shows_current_progress(self):
        """Should display current progress percentage."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            completed_phases=[],
        )

        menu = orchestrator.generate_progressive_menu(execution, step=1)

        assert "0%" in menu

    def test_menu_shows_visual_progress_bar(self):
        """Should display visual progress bar with blocks."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            completed_phases=[],
        )

        menu = orchestrator.generate_progressive_menu(execution, step=1)

        # Should have progress bar with ░ characters (empty when 0% complete)
        assert "░" in menu
        # Should have progress bar format [blocks]
        assert "[" in menu and "]" in menu

    def test_menu_shows_step_emoji(self):
        """Should show emoji indicator for current step."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
        )

        # Step 1 (requirements) should show 📋
        menu = orchestrator.generate_progressive_menu(execution, step=1)
        assert "📋" in menu

        # Step 2 (architecture) should show 🏗️
        menu = orchestrator.generate_progressive_menu(execution, step=2)
        assert "🏗" in menu

    def test_menu_shows_skip_badge_when_flag_set(self):
        """Should show [SKIP] badge when skip flag is set."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            skip_flags=SkipFlags(skip_requirements=True),
        )

        menu = orchestrator.generate_progressive_menu(execution, step=1)

        # Should show [SKIP] badge (uppercase, prominent)
        assert "[SKIP]" in menu

    def test_menu_shows_remaining_steps(self):
        """Should list remaining steps with brief descriptions."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            completed_phases=[],
        )

        menu = orchestrator.generate_progressive_menu(execution, step=1)

        # Should show remaining steps section
        assert "Remaining" in menu or "Next steps" in menu or "Upcoming" in menu

    def test_menu_shows_section_separators(self):
        """Should use visual separators between sections."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
        )

        menu = orchestrator.generate_progressive_menu(execution, step=1)

        # Should have separator lines (─ or ─ or similar)
        assert "─" in menu or "=" in menu or "-" in menu

    def test_menu_shows_feature_name(self):
        """Should display feature name in header."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
        )

        menu = orchestrator.generate_progressive_menu(execution, step=1)

        # Should show feature name
        assert "Test Feature" in menu

    def test_menu_progress_bar_accuracy(self):
        """Should show accurate progress bar based on completed phases."""
        orchestrator = FeatureOrchestrator()

        # 0 phases complete = 0%
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            completed_phases=[],
        )
        menu = orchestrator.generate_progressive_menu(execution, step=1)
        assert "0%" in menu or "[░░░" in menu

        # 1 phase complete = 20%
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            completed_phases=[FeaturePhase.REQUIREMENTS],
        )
        menu = orchestrator.generate_progressive_menu(execution, step=2)
        assert "20%" in menu or "[█░░" in menu

        # 2 phases complete = 40%
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            completed_phases=[
                FeaturePhase.REQUIREMENTS,
                FeaturePhase.ARCHITECTURE,
            ],
        )
        menu = orchestrator.generate_progressive_menu(execution, step=3)
        assert "40%" in menu or "[██░" in menu


class TestSkillInvocation:
    """Test skill invocation (@idea/@design/@oneshot)."""

    @patch('sdp.unified.feature.orchestrator.call_skill')
    def test_invokes_idea_skill_step1(self, mock_skill):
        """Should invoke @idea skill for requirements gathering."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
        )

        result = orchestrator.execute_step(execution, step=1)

        assert result.success
        mock_skill.assert_called_once_with("idea")

    @patch('sdp.unified.feature.orchestrator.call_skill')
    def test_invokes_design_skill_step2(self, mock_skill):
        """Should invoke @design skill after requirements approved."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
        )

        # Step 1 complete
        execution.completed_phases = [FeaturePhase.REQUIREMENTS]

        result = orchestrator.execute_step(execution, step=2)

        assert result.success
        mock_skill.assert_called_once_with("design")

    @patch('sdp.unified.feature.orchestrator.call_skill')
    def test_invokes_oneshot_skill_step3(self, mock_skill):
        """Should invoke @oneshot skill after architecture approved."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
        )

        # Steps 1-2 complete
        execution.completed_phases = [
            FeaturePhase.REQUIREMENTS,
            FeaturePhase.ARCHITECTURE,
        ]

        result = orchestrator.execute_step(execution, step=3)

        assert result.success
        mock_skill.assert_called_once_with("oneshot")


class TestSkipFlags:
    """Test skip flags functionality."""

    def test_skip_requirements_bypasses_step1(self):
        """Should skip requirements phase when flag is set."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            skip_flags=SkipFlags(skip_requirements=True),
        )

        result = orchestrator.execute_step(execution, step=1)

        # Should skip to step 2 (architecture)
        assert result.phase == FeaturePhase.ARCHITECTURE

    def test_skip_all_flags_goes_direct_to_oneshot(self):
        """Should go directly to @oneshot when all flags set."""
        orchestrator = FeatureOrchestrator()
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            skip_flags=SkipFlags(
                skip_requirements=True,
                skip_architecture=True,
                skip_uat=True,
            ),
        )

        result = orchestrator.execute_feature(execution)

        assert result.phase == FeaturePhase.EXECUTION


class TestApprovalGates:
    """Test approval gates between phases."""

    def test_requirements_gate_waits_for_approval(self):
        """Should wait for human approval before @design."""
        orchestrator = FeatureOrchestrator()

        # Mock user input
        with patch('builtins.input', return_value='y'):
            result = orchestrator.request_approval(
                phase=FeaturePhase.REQUIREMENTS,
                artifacts_path="docs/drafts/test.md",
            )

        assert result == "approved"

    def test_gate_reject_stops_execution(self):
        """Rejected gate should stop execution and report error."""
        orchestrator = FeatureOrchestrator()

        # Mock user input
        with patch('builtins.input', return_value='n'):
            result = orchestrator.request_approval(
                phase=FeaturePhase.REQUIREMENTS,
                artifacts_path="docs/drafts/test.md",
            )

        assert result == "rejected"

    def test_skip_flag_auto_approves_gate(self):
        """Skip flag should auto-approve gate without asking."""
        orchestrator = FeatureOrchestrator()

        # Create execution with skip flag
        execution = FeatureExecution(
            feature_id="test-feature",
            feature_name="Test Feature",
            skip_flags=SkipFlags(skip_requirements=True),
        )

        with patch('builtins.input') as mock_input:
            result = orchestrator.request_approval(
                phase=FeaturePhase.REQUIREMENTS,
                execution=execution,
                artifacts_path="docs/drafts/test.md",
            )

        # Should NOT call input (auto-approved)
        mock_input.assert_not_called()
        assert result == "skipped"


class TestCheckpointIntegration:
    """Test checkpoint integration throughout execution."""

    @patch('sdp.unified.checkpoint.repository.CheckpointRepository')
    def test_saves_checkpoint_after_each_phase(self, mock_repo):
        """Should save checkpoint after each phase completion."""
        orchestrator = FeatureOrchestrator(repository=mock_repo)
        execution = FeatureExecution(
            feature_id="sdp-118",
            feature_name="Test Feature",
        )

        # Execute phase 1
        execution.completed_phases = [FeaturePhase.REQUIREMENTS]
        orchestrator.after_phase(execution, phase=FeaturePhase.REQUIREMENTS)

        # Verify checkpoint saved
        mock_repo().save_checkpoint.assert_called_once()

    @patch('sdp.unified.checkpoint.repository.CheckpointRepository')
    def test_checkpoint_tracks_completed_phases(self, mock_repo):
        """Checkpoint should track which phases completed."""
        orchestrator = FeatureOrchestrator(repository=mock_repo)
        execution = FeatureExecution(
            feature_id="sdp-118",
            feature_name="Test Feature",
        )

        # Complete 2 phases
        execution.completed_phases = [
            FeaturePhase.REQUIREMENTS,
            FeaturePhase.ARCHITECTURE,
        ]
        orchestrator.after_phase(execution, phase=FeaturePhase.ARCHITECTURE)

        # Verify checkpoint saved with both phases
        saved_checkpoint = mock_repo().save_checkpoint.call_args[0][0]
        assert "requirements" in saved_checkpoint.completed_phases
        assert "architecture" in saved_checkpoint.completed_phases
