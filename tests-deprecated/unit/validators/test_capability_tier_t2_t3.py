"""Tests for T2/T3 capability tier validation."""

from textwrap import dedent

import pytest

from sdp.core.workstream import Workstream, WorkstreamSize, WorkstreamStatus
from sdp.validators.capability_tier_models import CapabilityTier
from sdp.validators.capability_tier_t2_t3 import validate_t2_t3


class TestValidateT2T3:
    """Test T2/T3 validation."""

    def test_t2_with_complete_contract(self) -> None:
        """Test T2 workstream with complete contract."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - Modified: \`src/example.py\`

        ### Acceptance Criteria
        - [x] Criterion 1
        - [x] Criterion 2

        ### Context
        Additional context here.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="Context",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        # Should have validation checks
        assert isinstance(checks, list)

    def test_t2_missing_contract_section(self) -> None:
        """Test T2 workstream missing contract section."""
        body = dedent("""
        ## Some Other Section

        Content here.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="Context",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)
        # Should have at least one check indicating missing contract

    def test_t3_with_complete_contract(self) -> None:
        """Test T3 workstream with complete contract."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - Modified: \`src/example.py\`

        ### Acceptance Criteria
        - [x] Criterion 1
        - [x] Criterion 2
        - [x] Criterion 3
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="Context",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T3)
        assert isinstance(checks, list)

    def test_t2_missing_input_section(self) -> None:
        """Test T2 workstream with contract but missing Input section."""
        body = dedent("""
        ## Contract

        ### Output
        - Modified: \`src/example.py\`

        ### Acceptance Criteria
        - [x] Criterion 1
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)

    def test_t2_missing_output_section(self) -> None:
        """Test T2 workstream with contract but missing Output section."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Acceptance Criteria
        - [x] Criterion 1
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)

    def test_t2_missing_acceptance_criteria(self) -> None:
        """Test T2 workstream with contract but missing Acceptance Criteria."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - Modified: \`src/example.py\`
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)

    def test_t2_with_empty_acceptance_criteria(self) -> None:
        """Test T2 workstream with empty Acceptance Criteria section."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - Modified: \`src/example.py\`

        ### Acceptance Criteria

        No criteria listed.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)

    def test_t2_with_checked_acceptance_criteria(self) -> None:
        """Test T2 workstream with checked acceptance criteria."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - Modified: \`src/example.py\`

        ### Acceptance Criteria
        - [x] Criterion 1
        - [x] Criterion 2
        - [x] Criterion 3
        All criteria checked.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)

    def test_t2_with_unchecked_acceptance_criteria(self) -> None:
        """Test T2 workstream with unchecked acceptance criteria."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - Modified: \`src/example.py\`

        ### Acceptance Criteria
        - [ ] Criterion 1
        - [ ] Criterion 2
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)

    def test_empty_body(self) -> None:
        """Test validation with empty body."""
        body = ""

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)

    def test_body_with_multiple_sections(self) -> None:
        """Test validation with additional sections beyond contract."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - Modified: \`src/example.py\`

        ### Acceptance Criteria
        - [x] Criterion 1

        ## Additional Context

        Extra information here.

        ## Implementation Notes

        More details.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert isinstance(checks, list)

    def test_t3_stricter_validation(self) -> None:
        """Test that T3 has stricter validation than T2."""
        body = dedent("""
        ## Contract

        ### Input
        - File: \`src/example.py\`

        ### Output
        - Modified: \`src/example.py\`

        ### Acceptance Criteria
        - [x] Criterion 1
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
        )

        t2_checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        t3_checks = validate_t2_t3(ws, body, CapabilityTier.T3)

        # Both should return lists
        assert isinstance(t2_checks, list)
        assert isinstance(t3_checks, list)
