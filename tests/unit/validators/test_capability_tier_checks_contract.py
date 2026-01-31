"""Tests for capability_tier_checks_contract validator."""

import pytest

from sdp.validators.capability_tier_checks_contract import (
    _check_contract_is_read_only,
    _check_no_vague_language,
    _check_verification_commands,
)


class TestCheckVerificationCommands:
    """Tests for _check_verification_commands."""

    def test_has_bash_commands_passes(self) -> None:
        """Verification section with bash commands passes."""
        section = "```bash\npytest tests/\n```"
        result = _check_verification_commands(section)
        assert result.passed

    def test_no_bash_blocks_fails(self) -> None:
        """Verification without bash blocks fails."""
        result = _check_verification_commands("Just text")
        assert not result.passed
        assert "bash" in result.message.lower()

    def test_empty_bash_block_fails(self) -> None:
        """Bash block with only comments fails."""
        section = "```bash\n# comment only\n```"
        result = _check_verification_commands(section)
        assert not result.passed


class TestCheckContractIsReadOnly:
    """Tests for _check_contract_is_read_only."""

    def test_no_contract_passes(self) -> None:
        """Body without contract section passes."""
        result = _check_contract_is_read_only("## Goal\nContent")
        assert result.passed

    def test_contract_spec_only_passes(self) -> None:
        """Contract with only specification (no class/def) passes."""
        body = "## Contract\n\nInterface: Foo\nRequired: bar(), baz()"
        result = _check_contract_is_read_only(body)
        assert result.passed

    def test_contract_with_implementation_fails(self) -> None:
        """Contract with full implementation fails."""
        body = "## Contract\n```python\nclass Foo:\n    def bar(self):\n        return 1\n```"
        result = _check_contract_is_read_only(body)
        assert not result.passed


class TestCheckNoVagueLanguage:
    """Tests for _check_no_vague_language."""

    def test_specific_language_passes(self) -> None:
        """Specific language passes."""
        result = _check_no_vague_language("Use pytest for testing")
        assert result.passed

    def test_vague_appropriate_fails(self) -> None:
        """Word 'appropriate' fails."""
        result = _check_no_vague_language("Use appropriate values")
        assert not result.passed
        assert "appropriate" in str(result.details or [])

    def test_vague_perhaps_fails(self) -> None:
        """Word 'perhaps' fails."""
        result = _check_no_vague_language("Perhaps we should")
        assert not result.passed
