"""Tests for Python 3.10 compatibility (tomllib/tomli)."""

import sys

import pytest


def test_tomllib_import_works():
    """Test that TOML parsing works on Python 3.10."""
    # This should work on both Python 3.10 (via tomli) and 3.11+ (via tomllib)
    if sys.version_info >= (3, 11):
        import tomllib
    else:
        import tomli as tomllib

    # Test basic functionality
    toml_content = """
[tool.sdp]
version = "0.7.0"
"""
    data = tomllib.loads(toml_content)
    assert data["tool"]["sdp"]["version"] == "0.7.0"


def test_hooks_common_imports():
    """Test that hooks.common module imports successfully."""
    from sdp.hooks import common

    # Verify find_project_root works
    # (assuming we're in SDP repo)
    root = common.find_project_root()
    assert root.exists()
    assert (root / "docs" / "workstreams").exists()


def test_quality_config_imports():
    """Test that quality.config module imports successfully."""
    from sdp.quality import config

    # Verify QualityGateConfigLoader works
    loader = config.QualityGateConfigLoader()
    assert loader.config is not None
