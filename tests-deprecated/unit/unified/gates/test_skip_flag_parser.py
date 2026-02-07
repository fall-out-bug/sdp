"""Tests for SkipFlagParser."""

import pytest

from sdp.unified.gates.parser import SkipFlagParser
from sdp.unified.gates.models import GateType


class TestSkipFlagParser:
    """Tests for SkipFlagParser."""

    def test_parse_empty_args(self):
        """Test parsing empty command line arguments."""
        parser = SkipFlagParser([])
        skips = parser.parse()

        assert skips == set()

    def test_parse_no_skip_flags(self):
        """Test parsing args with no skip flags."""
        parser = SkipFlagParser(["--verbose", "--output", "file.txt"])
        skips = parser.parse()

        assert skips == set()

    def test_parse_skip_requirements_flag(self):
        """Test parsing --skip-requirements flag."""
        parser = SkipFlagParser(["--skip-requirements"])
        skips = parser.parse()

        assert skips == {GateType.REQUIREMENTS}

    def test_parse_skip_architecture_flag(self):
        """Test parsing --skip-architecture flag."""
        parser = SkipFlagParser(["--skip-architecture"])
        skips = parser.parse()

        assert skips == {GateType.ARCHITECTURE}

    def test_parse_skip_uat_flag(self):
        """Test parsing --skip-uat flag."""
        parser = SkipFlagParser(["--skip-uat"])
        skips = parser.parse()

        assert skips == {GateType.UAT}

    def test_parse_multiple_skip_flags(self):
        """Test parsing multiple skip flags."""
        parser = SkipFlagParser([
            "--skip-requirements",
            "--skip-architecture",
            "--skip-uat",
        ])
        skips = parser.parse()

        assert skips == {
            GateType.REQUIREMENTS,
            GateType.ARCHITECTURE,
            GateType.UAT,
        }

    def test_parse_partial_skip_flags(self):
        """Test parsing some skip flags."""
        parser = SkipFlagParser([
            "--skip-requirements",
            "--verbose",
            "--skip-uat",
        ])
        skips = parser.parse()

        assert skips == {
            GateType.REQUIREMENTS,
            GateType.UAT,
        }

    def test_is_skip_required_for_requirements(self):
        """Test checking if requirements gate should be skipped."""
        parser = SkipFlagParser(["--skip-requirements"])

        assert parser.is_skip_required(GateType.REQUIREMENTS) is True
        assert parser.is_skip_required(GateType.ARCHITECTURE) is False
        assert parser.is_skip_required(GateType.UAT) is False

    def test_is_skip_required_for_architecture(self):
        """Test checking if architecture gate should be skipped."""
        parser = SkipFlagParser(["--skip-architecture"])

        assert parser.is_skip_required(GateType.REQUIREMENTS) is False
        assert parser.is_skip_required(GateType.ARCHITECTURE) is True
        assert parser.is_skip_required(GateType.UAT) is False

    def test_is_skip_required_for_uat(self):
        """Test checking if UAT gate should be skipped."""
        parser = SkipFlagParser(["--skip-uat"])

        assert parser.is_skip_required(GateType.REQUIREMENTS) is False
        assert parser.is_skip_required(GateType.ARCHITECTURE) is False
        assert parser.is_skip_required(GateType.UAT) is True

    def test_parse_is_idempotent(self):
        """Test that parse can be called multiple times."""
        parser = SkipFlagParser(["--skip-requirements"])

        skips1 = parser.parse()
        skips2 = parser.parse()

        assert skips1 == skips2

    def test_unknown_flags_are_ignored(self):
        """Test that unknown flags are ignored."""
        parser = SkipFlagParser([
            "--unknown-flag",
            "--another-unknown",
            "--skip-requirements",
        ])
        skips = parser.parse()

        assert skips == {GateType.REQUIREMENTS}
