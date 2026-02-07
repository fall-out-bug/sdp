"""Tests for quality gate configuration loader."""

from pathlib import Path
import pytest
from unittest.mock import Mock, patch, mock_open

from sdp.quality.config import QualityGateConfigLoader
from sdp.quality.exceptions import ConfigValidationError


class TestConfigPathResolution:
    """Tests for configuration path resolution."""

    def test_resolve_explicit_path_exists(self, tmp_path):
        """Test resolves explicit path when it exists."""
        config_file = tmp_path / "custom.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 80\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))

        assert loader._config_path == config_file
        assert loader.config.coverage.enabled is True

    def test_resolve_explicit_path_not_exists(self):
        """Test raises error when explicit path doesn't exist."""
        with pytest.raises(ConfigValidationError, match="Config file not found"):
            QualityGateConfigLoader(config_path="/nonexistent/path.toml")

    def test_resolve_default_paths_first_match(self, tmp_path, monkeypatch):
        """Test finds first matching default path."""
        monkeypatch.chdir(tmp_path)
        config_file = tmp_path / "quality-gate.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 80\n')

        loader = QualityGateConfigLoader()

        assert loader._config_path == Path("quality-gate.toml")

    def test_resolve_default_paths_second_match(self, tmp_path, monkeypatch):
        """Test finds second default path if first doesn't exist."""
        monkeypatch.chdir(tmp_path)
        config_file = tmp_path / ".quality-gate.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 80\n')

        loader = QualityGateConfigLoader()

        assert loader._config_path == Path(".quality-gate.toml")

    def test_resolve_no_defaults_use_defaults(self, tmp_path, monkeypatch):
        """Test uses default config when no files exist."""
        monkeypatch.chdir(tmp_path)

        loader = QualityGateConfigLoader()

        assert loader._config_path is None
        assert loader.config.coverage.minimum == 80


class TestConfigLoading:
    """Tests for TOML loading."""

    def test_load_invalid_toml(self, tmp_path):
        """Test raises error for invalid TOML syntax."""
        config_file = tmp_path / "bad.toml"
        config_file.write_text('invalid [[ toml syntax')

        with pytest.raises(ConfigValidationError, match="Invalid TOML"):
            QualityGateConfigLoader(config_path=str(config_file))

    def test_load_io_error(self, tmp_path):
        """Test raises error on file read error."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\n')

        with patch('builtins.open', side_effect=OSError("Permission denied")):
            with pytest.raises(ConfigValidationError, match="Error reading"):
                QualityGateConfigLoader(config_path=str(config_file))

    def test_load_default_config(self):
        """Test returns default config when no file exists."""
        loader = QualityGateConfigLoader()

        raw = loader.raw

        assert "coverage" in raw
        assert raw["coverage"]["enabled"] is True
        assert raw["coverage"]["minimum"] == 80


class TestConfigParsing:
    """Tests for configuration parsing."""

    def test_parse_config_value_error(self, tmp_path):
        """Test handles ValueError during parsing."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 80\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))

        # Mock parsers to raise ValueError
        with patch.object(loader._parsers, 'parse_coverage', side_effect=ValueError("Bad value")):
            with pytest.raises(ConfigValidationError, match="Configuration error"):
                loader._parse_config()

    def test_parse_config_type_error(self, tmp_path):
        """Test handles TypeError during parsing."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 80\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))

        # Mock parsers to raise TypeError
        with patch.object(loader._parsers, 'parse_coverage', side_effect=TypeError("Bad type")):
            with pytest.raises(ConfigValidationError, match="Configuration error"):
                loader._parse_config()


class TestConfigValidation:
    """Tests for configuration validation."""

    def test_validate_coverage_minimum_negative(self, tmp_path):
        """Test validates coverage.minimum >= 0."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = -10\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert any("coverage.minimum must be between 0 and 100" in e for e in errors)

    def test_validate_coverage_minimum_over_100(self, tmp_path):
        """Test validates coverage.minimum <= 100."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 150\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert any("coverage.minimum must be between 0 and 100" in e for e in errors)

    def test_validate_coverage_fail_under_negative(self, tmp_path):
        """Test validates coverage.fail_under >= 0."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 80\nfail_under = -5\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert any("coverage.fail_under must be between 0 and 100" in e for e in errors)

    def test_validate_coverage_fail_under_over_100(self, tmp_path):
        """Test validates coverage.fail_under <= 100."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 80\nfail_under = 120\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert any("coverage.fail_under must be between 0 and 100" in e for e in errors)

    def test_validate_complexity_max_cc_too_low(self, tmp_path):
        """Test validates complexity.max_cc >= 1."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[complexity]\nenabled = true\nmax_cc = 0\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert any("complexity.max_cc must be at least 1" in e for e in errors)

    def test_validate_complexity_max_average_cc_too_low(self, tmp_path):
        """Test validates complexity.max_average_cc >= 1."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[complexity]\nenabled = true\nmax_cc = 10\nmax_average_cc = 0\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert any("complexity.max_average_cc must be at least 1" in e for e in errors)

    def test_validate_file_size_max_lines_too_low(self, tmp_path):
        """Test validates file_size.max_lines >= 10."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[file_size]\nenabled = true\nmax_lines = 5\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert any("file_size.max_lines must be at least 10" in e for e in errors)

    def test_validate_file_size_max_imports_too_low(self, tmp_path):
        """Test validates file_size.max_imports >= 1."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[file_size]\nenabled = true\nmax_lines = 200\nmax_imports = 0\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert any("file_size.max_imports must be at least 1" in e for e in errors)

    def test_validate_valid_config(self, tmp_path):
        """Test validation passes for valid config."""
        config_file = tmp_path / "test.toml"
        config_file.write_text("""
[coverage]
enabled = true
minimum = 80
fail_under = 80

[complexity]
enabled = true
max_cc = 10
max_average_cc = 5

[file_size]
enabled = true
max_lines = 200
max_imports = 20
""")

        loader = QualityGateConfigLoader(config_path=str(config_file))
        errors = loader.validate()

        assert errors == []


class TestConfigProperties:
    """Tests for configuration properties."""

    def test_config_property(self, tmp_path):
        """Test config property returns parsed config."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 85\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))

        assert loader.config.coverage.minimum == 85

    def test_raw_property(self, tmp_path):
        """Test raw property returns raw TOML dict."""
        config_file = tmp_path / "test.toml"
        config_file.write_text('[coverage]\nenabled = true\nminimum = 85\n')

        loader = QualityGateConfigLoader(config_path=str(config_file))

        assert loader.raw["coverage"]["minimum"] == 85
