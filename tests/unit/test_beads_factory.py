"""Tests for Beads client factory with real as default."""

from pathlib import Path
from unittest.mock import patch
import warnings

import pytest

from sdp.beads.client import create_beads_client
from sdp.beads.cli import CLIBeadsClient
from sdp.beads.mock import MockBeadsClient


class TestBeadsClientFactory:
    """Test suite for create_beads_client factory function."""

    @patch.dict("os.environ", {"BEADS_USE_MOCK": "true"}, clear=False)
    def test_explicit_mock_via_env_var(self) -> None:
        """AC1: BEADS_USE_MOCK=true forces mock client."""
        # Act
        client = create_beads_client()

        # Assert
        assert isinstance(client, MockBeadsClient)

    def test_explicit_mock_via_parameter(self) -> None:
        """AC1: use_mock=True forces mock client."""
        # Act
        client = create_beads_client(use_mock=True)

        # Assert
        assert isinstance(client, MockBeadsClient)

    @patch("shutil.which")
    @patch("subprocess.run")
    @patch.dict("os.environ", {}, clear=True)
    def test_real_client_when_bd_installed(
        self, mock_run, mock_which
    ) -> None:
        """AC1: Returns real client when bd CLI is installed (default)."""
        # Arrange
        mock_which.return_value = "/usr/local/bin/bd"
        mock_run.return_value.returncode = 0
        mock_run.return_value.stdout = "bd version 0.1.0"

        # Act
        client = create_beads_client()

        # Assert
        assert isinstance(client, CLIBeadsClient)

    @patch("shutil.which")
    @patch.dict("os.environ", {}, clear=True)
    def test_falls_back_to_mock_with_warning(self, mock_which) -> None:
        """AC1: Falls back to mock with warning when bd not installed."""
        # Arrange
        mock_which.return_value = None  # bd not found

        # Act & Assert - should warn
        with warnings.catch_warnings(record=True) as w:
            warnings.simplefilter("always")
            client = create_beads_client()

            # Should return mock
            assert isinstance(client, MockBeadsClient)

            # Should have warning
            assert len(w) == 1
            assert "Beads CLI (bd) not found" in str(w[0].message)
            assert "go install" in str(w[0].message)

    @patch("shutil.which")
    @patch("subprocess.run")
    @patch.dict("os.environ", {}, clear=True)
    def test_project_dir_passed_to_cli_client(
        self, mock_run, mock_which
    ) -> None:
        """CLI client receives project_dir parameter."""
        # Arrange
        mock_which.return_value = "/usr/local/bin/bd"
        mock_run.return_value.returncode = 0
        mock_run.return_value.stdout = "bd version 0.1.0"

        project_dir = Path("/tmp/test-project")

        # Act
        client = create_beads_client(project_dir=project_dir)

        # Assert
        assert isinstance(client, CLIBeadsClient)
        assert client.project_dir == project_dir

    @patch("shutil.which")
    @patch.dict("os.environ", {"BEADS_USE_MOCK": "false"}, clear=False)
    def test_env_var_false_uses_real_if_available(self, mock_which) -> None:
        """AC1: BEADS_USE_MOCK=false explicitly requests real client."""
        # Arrange
        mock_which.return_value = "/usr/local/bin/bd"

        with patch("subprocess.run") as mock_run:
            mock_run.return_value.returncode = 0
            mock_run.return_value.stdout = "bd version 0.1.0"

            # Act
            client = create_beads_client()

            # Assert
            assert isinstance(client, CLIBeadsClient)

    @patch("shutil.which")
    @patch.dict("os.environ", {"BEADS_USE_MOCK": "FALSE"}, clear=False)
    def test_env_var_case_insensitive(self, mock_which) -> None:
        """Environment variable handling is case-insensitive."""
        # Arrange
        mock_which.return_value = "/usr/local/bin/bd"

        with patch("subprocess.run") as mock_run:
            mock_run.return_value.returncode = 0
            mock_run.return_value.stdout = "bd version 0.1.0"

            # Act
            client = create_beads_client()

            # Assert
            assert isinstance(client, CLIBeadsClient)

    def test_use_mock_parameter_overrides_env_var(self) -> None:
        """AC1: use_mock parameter takes precedence over env var."""
        # Arrange - env says real, parameter says mock
        with patch.dict("os.environ", {"BEADS_USE_MOCK": "false"}):
            # Act
            client = create_beads_client(use_mock=True)

            # Assert
            assert isinstance(client, MockBeadsClient)
