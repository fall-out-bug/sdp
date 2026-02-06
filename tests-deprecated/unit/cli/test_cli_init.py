"""Tests for CLI module."""

from click.testing import CliRunner


class TestCliMain:
    """Test main CLI entry point."""

    def test_main_group_exists(self) -> None:
        """Test that main CLI group is accessible."""
        from sdp.cli.main import main
        runner = CliRunner()
        result = runner.invoke(main, ["--help"])
        assert result.exit_code == 0
        assert "SDP (Spec-Driven Protocol)" in result.output

    def test_version_command(self) -> None:
        """Test version command."""
        from sdp.cli.main import main
        runner = CliRunner()
        result = runner.invoke(main, ["version"])
        assert result.exit_code == 0
        assert "sdp version" in result.output


class TestCliQuality:
    """Test quality gate commands."""

    def test_quality_group_exists(self) -> None:
        """Test that quality command group exists."""
        from sdp.cli.quality import quality
        runner = CliRunner()
        result = runner.invoke(quality, ["--help"])
        assert result.exit_code == 0


class TestCliWorkstream:
    """Test workstream commands."""

    def test_workstream_group_exists(self) -> None:
        """Test that workstream command group exists."""
        from sdp.cli.workstream import workstream
        runner = CliRunner()
        result = runner.invoke(workstream, ["--help"])
        assert result.exit_code == 0

    def test_workstream_parse_help(self) -> None:
        """Test workstream parse command help."""
        from sdp.cli.workstream import workstream
        runner = CliRunner()
        result = runner.invoke(workstream, ["parse", "--help"])
        assert result.exit_code == 0
        assert "Parse a workstream markdown file" in result.output

    def test_workstream_validate_help(self) -> None:
        """Test workstream validate command help."""
        from sdp.cli.workstream import workstream
        runner = CliRunner()
        result = runner.invoke(workstream, ["validate", "--help"])
        assert result.exit_code == 0
        assert "Validate workstream against capability tier" in result.output
