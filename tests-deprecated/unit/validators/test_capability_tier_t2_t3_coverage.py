"""Additional tests for validators/capability_tier_t2_t3.py

Increase coverage from 64% to 80%+ by covering error paths and edge cases.
"""

from textwrap import dedent

import pytest

from sdp.core.workstream import Workstream, WorkstreamSize, WorkstreamStatus
from sdp.validators.capability_tier_models import CapabilityTier
from sdp.validators.capability_tier_t2_t3 import validate_t2_t3


class TestValidateT2T3EdgeCases:
    """Test edge cases and error paths in T2/T3 validation."""

    def test_missing_interface_section(self) -> None:
        """T2/T3 requires Interface section in Contract."""
        body = dedent("""
        ## Contract
        
        ### Tests
        ```python
        def test_example():
            pass
        ```
        """)
        
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test",
            goal="Goal",
        )
        
        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        
        # Should fail interface_section_exists check
        assert any(c.name == "interface_section_exists" and not c.passed for c in checks)

    def test_interface_missing_code_block(self) -> None:
        """Interface section must contain Python code block."""
        body = dedent("""
        ## Contract
        
        ### Interface
        No code block here, just text.
        """)
        
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test",
            goal="Goal",
        )
        
        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        
        # Should fail interface_code_block check
        assert any(c.name == "interface_code_block" and not c.passed for c in checks)

    def test_tests_section_missing(self) -> None:
        """T2/T3 requires Tests section in Contract."""
        body = dedent("""
        ## Contract
        
        ### Interface
        ```python
        def example():
            pass
        ```
        """)
        
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test",
            goal="Goal",
        )
        
        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        
        # Should fail tests_section_exists check
        assert any(c.name == "tests_section_exists" and not c.passed for c in checks)

    def test_tests_missing_code_block(self) -> None:
        """Tests section must contain Python code block."""
        body = dedent("""
        ## Contract
        
        ### Interface
        ```python
        def example():
            pass
        ```
        
        ### Tests
        No code block here.
        """)
        
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test",
            goal="Goal",
        )
        
        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        
        # Should fail tests_code_block check
        assert any(c.name == "tests_code_block" and not c.passed for c in checks)

    def test_verification_section_missing(self) -> None:
        """T2/T3 requires Verification section."""
        body = dedent("""
        ## Contract
        
        ### Interface
        ```python
        def example():
            raise NotImplementedError
        ```
        
        ### Tests
        ```python
        def test_example():
            pass
        ```
        """)
        
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test",
            goal="Goal",
        )
        
        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        
        # Should fail verification_section_exists check
        assert any(c.name == "verification_section_exists" and not c.passed for c in checks)

    def test_t3_specific_scope_check(self) -> None:
        """T3 has additional scope check (TINY workstreams only)."""
        body = dedent("""
        ## Contract
        
        ### Interface
        ```python
        def example():
            raise NotImplementedError
        ```
        
        ### Tests
        ```python
        def test_example():
            pass
        ```
        
        ### Verification
        ```bash
        pytest tests/
        ```
        """)
        
        # Large workstream
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.LARGE,  # T3 requires SMALL
            title="Test",
            goal="Goal",
        )
        
        # T2 doesn't check scope
        t2_checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        assert not any(c.name == "scope_tiny" for c in t2_checks)
        
        # T3 checks scope
        t3_checks = validate_t2_t3(ws, body, CapabilityTier.T3)
        assert any(c.name == "scope_tiny" for c in t3_checks)

    def test_all_checks_pass_for_valid_t2(self) -> None:
        """Valid T2 workstream passes all checks."""
        body = dedent("""
        ## Contract
        
        ### Interface
        ```python
        def example():
            raise NotImplementedError
        ```
        
        ### Tests
        ```python
        def test_example():
            assert example() is not None
        ```
        
        ## Verification
        ```bash
        pytest tests/
        ```
        
        No vague language here.
        """)
        
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test",
            goal="Goal",
        )
        
        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        
        # All checks should pass
        failed_checks = [c for c in checks if not c.passed]
        assert len(failed_checks) == 0, f"Failed checks: {[c.name for c in failed_checks]}"

    def test_early_return_when_contract_missing(self) -> None:
        """Return early when Contract section missing (can't continue validation)."""
        body = "## Some Other Section\n\nNo contract here."
        
        ws = Workstream(
            ws_id="00-001-01",
            feature="F001",
            status=WorkstreamStatus.BACKLOG,
            size=WorkstreamSize.SMALL,
            title="Test",
            goal="Goal",
        )
        
        checks = validate_t2_t3(ws, body, CapabilityTier.T2)
        
        # Should return early with contract_section_exists failure
        assert len(checks) == 1
        assert checks[0].name == "contract_section_exists"
        assert not checks[0].passed
