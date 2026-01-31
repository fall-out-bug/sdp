"""Tests for ws_completion verifier."""

from pathlib import Path
from unittest.mock import patch

import pytest

from sdp.validators.ws_completion import (
    CheckResult,
    VerificationResult,
    WSCompletionVerifier,
)


@pytest.fixture
def ws_dir(tmp_path: Path) -> Path:
    """Create temp workstream directory."""
    for sub in ("backlog", "completed"):
        (tmp_path / sub).mkdir()
    return tmp_path


@pytest.fixture
def verifier(ws_dir: Path) -> WSCompletionVerifier:
    """Create verifier with temp ws_dir."""
    return WSCompletionVerifier(ws_dir=ws_dir)


class TestWSCompletionVerifier:
    """Tests for WSCompletionVerifier."""

    def test_verify_ws_not_found_returns_failed(
        self, verifier: WSCompletionVerifier
    ) -> None:
        """verify when WS not found returns failed result."""
        result = verifier.verify("99-999-99")
        assert not result.passed
        assert "not found" in result.checks[0].message.lower()

    def test_verify_output_files_missing(
        self, verifier: WSCompletionVerifier, ws_dir: Path
    ) -> None:
        """verify_output_files reports missing scope files."""
        ws = ws_dir / "completed" / "00-032-01.md"
        ws.write_text("""---
ws_id: 00-032-01
scope_files:
  - src/missing.py
---
""")
        checks = verifier.verify_output_files({"scope_files": ["src/missing.py"]})
        assert len(checks) == 1
        assert not checks[0].passed

    def test_verify_output_files_exists(
        self, verifier: WSCompletionVerifier, tmp_path: Path
    ) -> None:
        """verify_output_files passes when file exists."""
        existing = tmp_path / "exists.py"
        existing.write_text("# exists")
        checks = verifier.verify_output_files(
            {"scope_files": [str(existing)]}
        )
        assert len(checks) == 1
        assert checks[0].passed

    def test_verify_commands_empty_list(self, verifier: WSCompletionVerifier) -> None:
        """verify_commands with empty list returns empty."""
        checks = verifier.verify_commands({"verification_commands": []})
        assert checks == []

    @patch("sdp.validators.ws_completion.subprocess.run")
    def test_verify_commands_success(
        self, mock_run: object, verifier: WSCompletionVerifier
    ) -> None:
        """verify_commands passes when command succeeds."""
        from unittest.mock import MagicMock

        mock_run.return_value = MagicMock(returncode=0, stdout="", stderr="")
        checks = verifier.verify_commands(
            {"verification_commands": ["echo ok"]}
        )
        assert len(checks) == 1
        assert checks[0].passed

    def test_verify_coverage_no_scope_returns_none(
        self, verifier: WSCompletionVerifier
    ) -> None:
        """verify_coverage with no scope returns None."""
        assert verifier.verify_coverage({}) is None
        assert verifier.verify_coverage({"scope_files": []}) is None

    def test_verify_coverage_no_py_files_returns_none(
        self, verifier: WSCompletionVerifier
    ) -> None:
        """verify_coverage with no Python files returns None."""
        assert verifier.verify_coverage({"scope_files": ["docs/readme.md"]}) is None

    def test_find_ws_file_returns_path(
        self, verifier: WSCompletionVerifier, ws_dir: Path
    ) -> None:
        """_find_ws_file finds existing WS."""
        ws = ws_dir / "completed" / "00-032-01-guard.md"
        ws.write_text("---\nws_id: 00-032-01\n---")
        path = verifier._find_ws_file("00-032-01")
        assert path is not None
        assert "00-032-01" in str(path)

    def test_find_ws_file_not_found(self, verifier: WSCompletionVerifier) -> None:
        """_find_ws_file returns None for missing WS."""
        assert verifier._find_ws_file("99-999-99") is None

    def test_parse_coverage_from_output(self, verifier: WSCompletionVerifier) -> None:
        """_parse_coverage_from_output extracts percentage."""
        output = "TOTAL    100    20    80%"
        assert verifier._parse_coverage_from_output(output) == 80.0

    def test_parse_coverage_not_found(self, verifier: WSCompletionVerifier) -> None:
        """_parse_coverage_from_output returns None when not found."""
        assert verifier._parse_coverage_from_output("no coverage here") is None

    def test_extract_coverage_from_check(self, verifier: WSCompletionVerifier) -> None:
        """_extract_coverage parses check message."""
        check = CheckResult("cov", True, "Coverage: 85.5%", None)
        assert verifier._extract_coverage(check) == 85.5

    def test_extract_coverage_none_check(self, verifier: WSCompletionVerifier) -> None:
        """_extract_coverage returns None for None check."""
        assert verifier._extract_coverage(None) is None

    def test_parse_ws_file_extracts_scope(
        self, verifier: WSCompletionVerifier, ws_dir: Path
    ) -> None:
        """_parse_ws_file extracts scope_files from frontmatter."""
        ws = ws_dir / "completed" / "00-032-01.md"
        ws.write_text("""---
ws_id: 00-032-01
scope_files:
  - src/foo.py
  - src/bar.py
---

## Verification
""")
        data = verifier._parse_ws_file(ws)
        assert "src/foo.py" in data["scope_files"]
        assert "src/bar.py" in data["scope_files"]

    def test_verify_full_flow_passes(
        self, verifier: WSCompletionVerifier, ws_dir: Path, tmp_path: Path
    ) -> None:
        """verify passes when all checks pass."""
        scope_file = tmp_path / "exists.py"
        scope_file.write_text("# exists")
        ws = ws_dir / "completed" / "00-032-01.md"
        ws.write_text(f"""---
ws_id: 00-032-01
scope_files:
  - {scope_file}
---

## Verification
""")
        with patch("sdp.validators.ws_completion.subprocess.run") as mock_run:
            from unittest.mock import MagicMock

            mock_run.return_value = MagicMock(returncode=0, stdout="TOTAL 100 20 80%")
            result = verifier.verify("00-032-01")
            assert result.passed
