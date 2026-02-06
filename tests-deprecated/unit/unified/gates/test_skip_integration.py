"""Tests for SkipFlagParser integration with ApprovalGateManager."""

import pytest

from sdp.unified.gates.manager import ApprovalGateManager
from sdp.unified.gates.parser import SkipFlagParser
from sdp.unified.gates.models import ApprovalStatus, GateType
from sdp.unified.checkpoint.repository import CheckpointRepository
from sdp.unified.checkpoint.schema import Checkpoint, CheckpointStatus
from datetime import datetime


@pytest.fixture
def temp_db_path(tmp_path):
    """Create temporary database path."""
    return str(tmp_path / "test.db")


@pytest.fixture
def checkpoint_repo(temp_db_path):
    """Create checkpoint repository for testing."""
    repo = CheckpointRepository(temp_db_path)
    repo.initialize()
    yield repo
    repo.close()


@pytest.fixture
def sample_checkpoint():
    """Create sample checkpoint for testing."""
    return Checkpoint(
        feature="F01",
        agent_id="agent-001",
        status=CheckpointStatus.IN_PROGRESS,
        completed_ws=[],
        execution_order=["WS-001", "WS-002"],
        started_at=datetime.now(),
        current_ws="WS-001",
    )


class TestSkipFlagParserIntegration:
    """Tests for SkipFlagParser integration with ApprovalGateManager."""

    def test_manager_with_skip_parser(self, checkpoint_repo):
        """Test creating manager with skip parser."""
        skip_parser = SkipFlagParser(["--skip-requirements"])
        manager = ApprovalGateManager(
            checkpoint_repo,
            skip_parser=skip_parser
        )

        assert manager.checkpoint_repo == checkpoint_repo
        assert manager.skip_parser == skip_parser

    def test_manager_without_skip_parser_is_none(self, checkpoint_repo):
        """Test that manager without skip parser has None."""
        manager = ApprovalGateManager(checkpoint_repo)

        assert manager.skip_parser is None

    def test_request_approval_with_skip_flag_skips_gate(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that request_approval skips gate when flag is set."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        skip_parser = SkipFlagParser(["--skip-requirements"])
        manager = ApprovalGateManager(
            checkpoint_repo,
            skip_parser=skip_parser
        )

        # Request approval (should auto-skip)
        manager.request_approval(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            requestor="agent-001"
        )

        # Verify gate is skipped
        is_skipped = manager.is_skipped(
            feature="F01",
            gate_type=GateType.REQUIREMENTS
        )
        assert is_skipped is True

    def test_request_approval_without_skip_flag_does_not_skip(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that request_approval doesn't skip when flag not set."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        skip_parser = SkipFlagParser([])  # No skip flags
        manager = ApprovalGateManager(
            checkpoint_repo,
            skip_parser=skip_parser
        )

        # Request approval (should not skip)
        manager.request_approval(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            requestor="agent-001"
        )

        # Verify gate is NOT skipped
        is_skipped = manager.is_skipped(
            feature="F01",
            gate_type=GateType.REQUIREMENTS
        )
        assert is_skipped is False

    def test_auto_skip_gates_skips_all_flagged_gates(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that auto_skip_gates skips all gates with flags."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        skip_parser = SkipFlagParser([
            "--skip-requirements",
            "--skip-architecture"
        ])
        manager = ApprovalGateManager(
            checkpoint_repo,
            skip_parser=skip_parser
        )

        # Auto-skip all gates
        manager.auto_skip_gates(feature="F01")

        # Verify both gates are skipped
        assert manager.is_skipped("F01", GateType.REQUIREMENTS) is True
        assert manager.is_skipped("F01", GateType.ARCHITECTURE) is True
        assert manager.is_skipped("F01", GateType.UAT) is False

    def test_auto_skip_gates_with_no_flags_skips_none(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that auto_skip_gates with no flags skips none."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        skip_parser = SkipFlagParser([])  # No flags
        manager = ApprovalGateManager(
            checkpoint_repo,
            skip_parser=skip_parser
        )

        # Auto-skip all gates
        manager.auto_skip_gates(feature="F01")

        # Verify no gates are skipped
        assert manager.is_skipped("F01", GateType.REQUIREMENTS) is False
        assert manager.is_skipped("F01", GateType.ARCHITECTURE) is False
        assert manager.is_skipped("F01", GateType.UAT) is False

    def test_auto_skip_gates_without_parser_is_noop(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that auto_skip_gates without parser does nothing."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)  # No skip_parser

        # Should not raise error, just do nothing
        manager.auto_skip_gates(feature="F01")

        # Verify no gates are skipped
        assert manager.is_skipped("F01", GateType.REQUIREMENTS) is False
        assert manager.is_skipped("F01", GateType.ARCHITECTURE) is False
        assert manager.is_skipped("F01", GateType.UAT) is False

    def test_request_approval_without_parser_is_noop(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that request_approval without parser does nothing."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        manager = ApprovalGateManager(checkpoint_repo)  # No skip_parser

        # Should not raise error, just do nothing
        manager.request_approval(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            requestor="agent-001"
        )

        # Verify gate is NOT skipped
        assert manager.is_skipped("F01", GateType.REQUIREMENTS) is False

    def test_approve_respects_existing_skip_status(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that approve doesn't change already-skipped gates."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        skip_parser = SkipFlagParser(["--skip-requirements"])
        manager = ApprovalGateManager(
            checkpoint_repo,
            skip_parser=skip_parser
        )

        # Skip the gate first
        manager.skip(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            reason="Auto-skipped"
        )

        # Try to approve (should respect skip status)
        manager.approve(
            feature="F01",
            gate_type=GateType.REQUIREMENTS,
            approved_by="user123"
        )

        # Verify still skipped (approval ignored for skipped gates)
        status = manager.get_gate_status("F01", GateType.REQUIREMENTS)
        assert status == ApprovalStatus.SKIPPED

    def test_multiple_skip_flags_all_respected(
        self, checkpoint_repo, sample_checkpoint
    ):
        """Test that all skip flags are respected."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        skip_parser = SkipFlagParser([
            "--skip-requirements",
            "--skip-architecture",
            "--skip-uat"
        ])
        manager = ApprovalGateManager(
            checkpoint_repo,
            skip_parser=skip_parser
        )

        # Auto-skip all gates
        manager.auto_skip_gates(feature="F01")

        # Verify all gates are skipped
        assert manager.is_skipped("F01", GateType.REQUIREMENTS) is True
        assert manager.is_skipped("F01", GateType.ARCHITECTURE) is True
        assert manager.is_skipped("F01", GateType.UAT) is True

    def test_request_approval_logs_skip_action(
        self, checkpoint_repo, sample_checkpoint, caplog
    ):
        """Test that request_approval logs when skipping gate."""
        checkpoint_id = checkpoint_repo.save_checkpoint(sample_checkpoint)
        skip_parser = SkipFlagParser(["--skip-requirements"])
        manager = ApprovalGateManager(
            checkpoint_repo,
            skip_parser=skip_parser
        )

        # Enable log capture
        with caplog.at_level("INFO"):
            # Request approval (should auto-skip)
            manager.request_approval(
                feature="F01",
                gate_type=GateType.REQUIREMENTS,
                requestor="agent-001"
            )

        # Verify log message
        assert "skip" in caplog.text.lower()
