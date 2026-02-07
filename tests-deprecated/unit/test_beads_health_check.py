"""Tests for Beads CLI health check."""

import subprocess
from unittest.mock import Mock, patch

import pytest

from sdp.health_checks.checks import BeadsCLICheck


class TestBeadsCLICheck:
    """Test suite for BeadsCLICheck."""

    @patch("shutil.which")
    @patch("subprocess.run")
    def test_beads_installed_with_version(self, mock_run: Mock, mock_which: Mock) -> None:
        """AC1: sdp doctor shows Beads CLI status when installed."""
        # Arrange
        mock_which.side_effect = lambda cmd: {
            "go": "/usr/local/bin/go",
            "bd": "/Users/test/go/bin/bd",
        }.get(cmd)
        
        mock_run.return_value = Mock(
            returncode=0,
            stdout="bd version 0.1.0\n",
            stderr="",
        )
        
        check = BeadsCLICheck()
        
        # Act
        result = check.run()
        
        # Assert
        assert result.passed is True
        assert "0.1.0" in result.message
        assert "/Users/test/go/bin/bd" in result.message
        assert result.remediation is None

    @patch("shutil.which")
    def test_beads_not_installed_with_go(self, mock_which: Mock) -> None:
        """AC2: Shows installation instructions when bd absent but Go present."""
        # Arrange
        mock_which.side_effect = lambda cmd: {
            "go": "/usr/local/bin/go",
            "bd": None,
        }.get(cmd)
        
        check = BeadsCLICheck()
        
        # Act
        result = check.run()
        
        # Assert
        assert result.passed is True  # Optional check, so still passes
        assert "not installed" in result.message.lower()
        assert result.remediation is not None
        assert "go install github.com/steveyegge/beads/cmd/bd@latest" in result.remediation
        # Should NOT include Go installation step since Go is present
        assert "Install Go:" not in result.remediation

    @patch("shutil.which")
    def test_beads_not_installed_without_go(self, mock_which: Mock) -> None:
        """AC2: Shows Go + Beads installation when both absent."""
        # Arrange
        mock_which.return_value = None  # Neither go nor bd found
        
        check = BeadsCLICheck()
        
        # Act
        result = check.run()
        
        # Assert
        assert result.passed is True  # Optional check
        assert "not installed" in result.message.lower()
        assert result.remediation is not None
        # Should include BOTH Go and Beads installation
        assert "Install Go:" in result.remediation
        assert "brew install go" in result.remediation
        assert "go install github.com/steveyegge/beads/cmd/bd@latest" in result.remediation
        assert "docs/setup/beads-installation.md" in result.remediation

    @patch("shutil.which")
    @patch("subprocess.run")
    def test_version_parsing_various_formats(self, mock_run: Mock, mock_which: Mock) -> None:
        """Test version parsing handles different output formats."""
        # Arrange
        mock_which.side_effect = lambda cmd: {
            "go": "/usr/local/bin/go",
            "bd": "/usr/local/bin/bd",
        }.get(cmd)
        
        # Test various version output formats
        test_cases = [
            ("bd version 0.1.0", "0.1.0"),
            ("0.2.5", "0.2.5"),
            ("v1.0.0-beta", "1.0.0-beta"),
            ("", "unknown"),
        ]
        
        check = BeadsCLICheck()
        
        for output, expected_version in test_cases:
            mock_run.return_value = Mock(
                returncode=0,
                stdout=output,
                stderr="",
            )
            
            # Act
            result = check.run()
            
            # Assert
            assert result.passed is True
            if expected_version != "unknown":
                assert expected_version in result.message

    @patch("shutil.which")
    @patch("subprocess.run")
    def test_version_command_fails(self, mock_run: Mock, mock_which: Mock) -> None:
        """Test graceful handling when version command fails."""
        # Arrange
        mock_which.side_effect = lambda cmd: {
            "go": "/usr/local/bin/go",
            "bd": "/usr/local/bin/bd",
        }.get(cmd)
        
        mock_run.side_effect = subprocess.TimeoutExpired("bd --version", 5)
        
        check = BeadsCLICheck()
        
        # Act
        result = check.run()
        
        # Assert
        assert result.passed is True
        assert "unknown" in result.message.lower()

    @patch("shutil.which")
    @patch("subprocess.run")
    def test_remediation_includes_path_setup(self, mock_run: Mock, mock_which: Mock) -> None:
        """AC3: Remediation includes PATH setup instructions."""
        # Arrange
        mock_which.side_effect = lambda cmd: {
            "go": "/usr/local/bin/go",
            "bd": None,
        }.get(cmd)
        
        check = BeadsCLICheck()
        
        # Act
        result = check.run()
        
        # Assert
        assert result.remediation is not None
        assert "export PATH=$PATH:$(go env GOPATH)/bin" in result.remediation
        assert "bd --version" in result.remediation

    def test_check_is_optional(self) -> None:
        """AC4: Check is marked as optional (not critical)."""
        check = BeadsCLICheck()
        assert check.critical is False
        assert check.name == "Beads CLI"
