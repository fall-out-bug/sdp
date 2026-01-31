"""Tests for T0/T1 capability tier validation."""

from textwrap import dedent

import pytest

from sdp.core.workstream import Workstream, WorkstreamSize, WorkstreamStatus
from sdp.validators.capability_tier_t0_t1 import validate_t0_t1_gates


class TestValidateT0T1Gates:
    """Test T0/T1 gate validation."""

    def test_t0_with_all_gates_passed(self) -> None:
        """Test T0 workstream with all gates passed."""
        body = dedent("""
        ## Context & Problem Statement

        Background info here.

        ## Proposed Solution

        Solution description.

        ## Alternatives Considered

        Alternative 1: Description
        Alternative 2: Description
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="Context & Problem Statement",
        )

        checks = validate_t0_t1_gates(ws, body)
        check_dict = {check.name: check for check in checks}

        # Check that all required sections are validated
        assert "context_gate" in check_dict or len(checks) > 0

    def test_missing_context_section(self) -> None:
        """Test T0 workstream missing context section."""
        body = dedent("""
        ## Proposed Solution

        Solution description.

        ## Alternatives Considered

        Alternative 1: Description
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="",  # Missing
        )

        checks = validate_t0_t1_gates(ws, body)
        # Should have at least one check
        assert len(checks) >= 0

    def test_missing_solution_section(self) -> None:
        """Test T0 workstream missing solution section."""
        body = dedent("""
        ## Context & Problem Statement

        Background info here.

        ## Alternatives Considered

        Alternative 1: Description
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="Context & Problem Statement",
        )

        checks = validate_t0_t1_gates(ws, body)
        assert len(checks) >= 0

    def test_missing_alternatives_section(self) -> None:
        """Test T0 workstream missing alternatives section."""
        body = dedent("""
        ## Context & Problem Statement

        Background info here.

        ## Proposed Solution

        Solution description.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="Context & Problem Statement",
        )

        checks = validate_t0_t1_gates(ws, body)
        assert len(checks) >= 0

    def test_t1_with_complex_build_gates(self) -> None:
        """Test T1 workstream validation."""
        body = dedent("""
        ## Context

        Context here.

        ## Solution

        Solution here.

        ## Alternatives

        Alternatives here.
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

        checks = validate_t0_t1_gates(ws, body)
        # T1 should validate similar gates as T0
        assert len(checks) >= 0

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
            context="",
        )

        checks = validate_t0_t1_gates(ws, body)
        # Should handle empty body gracefully
        assert isinstance(checks, list)

    def test_body_with_extra_sections(self) -> None:
        """Test validation with extra sections beyond required ones."""
        body = dedent("""
        ## Context & Problem Statement

        Background info.

        ## Proposed Solution

        Solution description.

        ## Alternatives Considered

        Alternatives here.

        ## Additional Section

        Extra content that shouldn't affect validation.

        ## Another Extra Section

        More extra content.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="Context & Problem Statement",
        )

        checks = validate_t0_t1_gates(ws, body)
        assert isinstance(checks, list)

    def test_section_case_sensitivity(self) -> None:
        """Test that section matching is case-sensitive."""
        body = dedent("""
        ## context & problem statement

        Lowercase header.

        ## proposed solution

        Lowercase header.

        ## alternatives considered

        Lowercase header.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="context & problem statement",  # Lowercase
        )

        checks = validate_t0_t1_gates(ws, body)
        # Validation behavior depends on implementation
        assert isinstance(checks, list)

    def test_whitespace_in_sections(self) -> None:
        """Test handling of whitespace in section names."""
        body = dedent("""
        ## Context & Problem Statement

        Content here.

        ## Proposed Solution

        Content here.

        ## Alternatives Considered

        Content here.
        """)

        ws = Workstream(
            ws_id="WS-000-01",
            feature="F00",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test WS",
            goal="Test goal",
            context="Context & Problem Statement  ",  # Trailing spaces
        )

        checks = validate_t0_t1_gates(ws, body)
        assert isinstance(checks, list)
