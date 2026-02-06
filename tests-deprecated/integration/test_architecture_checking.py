"""Integration tests for architecture checking."""

import tempfile
from pathlib import Path
from textwrap import dedent

import pytest

from sdp.quality.config import QualityGateConfigLoader
from sdp.quality.architecture import ArchitectureChecker


class TestArchitectureChecking:
    """Integration tests for portable architecture checks."""

    def test_clean_architecture_default_config(self, tmp_path: Path) -> None:
        """Test architecture checking with default clean architecture config."""
        # Create domain entity with violation
        domain_file = tmp_path / "domain" / "entity.py"
        domain_file.parent.mkdir(parents=True)
        domain_file.write_text(
            dedent("""
                from infrastructure.database import Database  # VIOLATION

                class UserEntity:
                    def __init__(self, name: str):
                        self.name = name
            """)
        )

        # Load default config (no quality-gate.toml)
        config_loader = QualityGateConfigLoader(None)

        violations: list = []
        checker = ArchitectureChecker(
            config=config_loader.config.architecture,
            violations=violations,
        )

        import ast

        tree = ast.parse(domain_file.read_text())
        checker.check_architecture(domain_file, tree)

        assert len(violations) == 1
        assert "domain" in violations[0].message
        assert "infrastructure" in violations[0].message

    def test_hexagonal_architecture_config(self, tmp_path: Path) -> None:
        """Test with hexagonal architecture configuration."""
        # Create quality-gate.toml with hexagonal architecture
        config_file = tmp_path / "quality-gate.toml"
        config_file.write_text(
            dedent("""
                [architecture]
                enabled = true
                forbid_violations = ["domain -> infrastructure"]

                [architecture.layers.domain]
                path_regex = "(^|/)domain/"
                module_regex = "(^|\\\\.)domain(\\\\.|$)"

                [architecture.layers.infrastructure]
                path_regex = "(^|/)infrastructure/"
                module_regex = "(^|\\\\.)infrastructure(\\\\.|$)"
            """)
        )

        # Create violating file
        domain_file = tmp_path / "domain" / "entity.py"
        domain_file.parent.mkdir(parents=True)
        domain_file.write_text(
            dedent("""
                from infrastructure.persistence import Repository

                class Entity:
                    pass
            """)
        )

        # Load config
        config_loader = QualityGateConfigLoader(config_file)

        violations: list = []
        checker = ArchitectureChecker(
            config=config_loader.config.architecture,
            violations=violations,
        )

        import ast

        tree = ast.parse(domain_file.read_text())
        checker.check_architecture(domain_file, tree)

        assert len(violations) == 1
        assert "domain" in violations[0].message
        assert "infrastructure" in violations[0].message

    def test_allowed_imports(self, tmp_path: Path) -> None:
        """Test that allowed imports don't trigger violations."""
        config_file = tmp_path / "quality-gate.toml"
        config_file.write_text(
            dedent("""
                [architecture]
                enabled = true
                forbid_violations = ["application -> domain"]

                [architecture.layers.domain]
                path_regex = "(^|/)domain/"

                [architecture.layers.application]
                path_regex = "(^|/)application/"
                # Note: domain -> application is NOT forbidden
            """)
        )

        # Domain importing from application (allowed by this config)
        domain_file = tmp_path / "domain" / "model.py"
        domain_file.parent.mkdir(parents=True)
        domain_file.write_text(
            dedent("""
                from application.service import Service

                class Model:
                    pass
            """)
        )

        config_loader = QualityGateConfigLoader(config_file)

        violations: list = []
        checker = ArchitectureChecker(
            config=config_loader.config.architecture,
            violations=violations,
        )

        import ast

        tree = ast.parse(domain_file.read_text())
        checker.check_architecture(domain_file, tree)

        # No violation because domain -> application is not in forbid_violations
        assert len(violations) == 0

    def test_custom_layer_patterns(self, tmp_path: Path) -> None:
        """Test with custom layer patterns."""
        config_file = tmp_path / "quality-gate.toml"
        config_file.write_text(
            dedent("""
                [architecture]
                enabled = true
                forbid_violations = ["core -> api"]

                [architecture.layers.core]
                path_regex = "(^|/)core/"

                [architecture.layers.api]
                path_regex = "(^|/)api/"
            """)
        )

        # Core importing from API (violation)
        core_file = tmp_path / "core" / "logic.py"
        core_file.parent.mkdir(parents=True)
        core_file.write_text(
            dedent("""
                from api.controller import Controller

                class Logic:
                    pass
            """)
        )

        config_loader = QualityGateConfigLoader(config_file)

        violations: list = []
        checker = ArchitectureChecker(
            config=config_loader.config.architecture,
            violations=violations,
        )

        import ast

        tree = ast.parse(core_file.read_text())
        checker.check_architecture(core_file, tree)

        assert len(violations) == 1
        assert "core" in violations[0].message
        assert "api" in violations[0].message

    def test_python_standard_library_imports(self, tmp_path: Path) -> None:
        """Test that standard library imports don't trigger violations."""
        domain_file = tmp_path / "domain" / "entity.py"
        domain_file.parent.mkdir(parents=True)
        domain_file.write_text(
            dedent("""
                from dataclasses import dataclass
                from typing import List
                from datetime import datetime

                @dataclass
                class Entity:
                    created: datetime
                    items: List[str]
            """)
        )

        config_loader = QualityGateConfigLoader(None)

        violations: list = []
        checker = ArchitectureChecker(
            config=config_loader.config.architecture,
            violations=violations,
        )

        import ast

        tree = ast.parse(domain_file.read_text())
        checker.check_architecture(domain_file, tree)

        # Standard library imports should not trigger violations
        assert len(violations) == 0

    def test_disabled_architecture_checks(self, tmp_path: Path) -> None:
        """Test that disabled config doesn't check architecture."""
        config_file = tmp_path / "quality-gate.toml"
        config_file.write_text(
            dedent("""
                [architecture]
                enabled = false
            """)
        )

        domain_file = tmp_path / "domain" / "entity.py"
        domain_file.parent.mkdir(parents=True)
        domain_file.write_text(
            dedent("""
                from infrastructure.db import Database  # Would be violation if enabled

                class Entity:
                    pass
            """)
        )

        config_loader = QualityGateConfigLoader(config_file)

        violations: list = []
        checker = ArchitectureChecker(
            config=config_loader.config.architecture,
            violations=violations,
        )

        import ast

        tree = ast.parse(domain_file.read_text())
        checker.check_architecture(domain_file, tree)

        # No violations because architecture checks are disabled
        assert len(violations) == 0


class TestConfigFileLoading:
    """Tests for loading architecture configuration from files."""

    def test_load_hexagonal_arch_config(self, tmp_path: Path) -> None:
        """Test loading hexagonal architecture example config."""
        examples_dir = Path(__file__).parent.parent.parent.parent / "docs" / "examples"
        config_path = examples_dir / "quality-gate-hexagonal.toml"

        if not config_path.exists():
            pytest.skip(f"Example config not found: {config_path}")

        config_loader = QualityGateConfigLoader(config_path)

        assert config_loader.config.architecture.enabled is True
        assert len(config_loader.config.architecture.forbid_violations) > 0
        assert "domain -> infrastructure" in config_loader.config.architecture.forbid_violations

    def test_load_onion_arch_config(self, tmp_path: Path) -> None:
        """Test loading onion architecture example config."""
        examples_dir = Path(__file__).parent.parent.parent.parent / "docs" / "examples"
        config_path = examples_dir / "quality-gate-onion.toml"

        if not config_path.exists():
            pytest.skip(f"Example config not found: {config_path}")

        config_loader = QualityGateConfigLoader(config_path)

        assert config_loader.config.architecture.enabled is True
        assert len(config_loader.config.architecture.layer_patterns) > 0

    def test_load_layered_arch_config(self, tmp_path: Path) -> None:
        """Test loading layered architecture example config."""
        examples_dir = Path(__file__).parent.parent.parent.parent / "docs" / "examples"
        config_path = examples_dir / "quality-gate-layered.toml"

        if not config_path.exists():
            pytest.skip(f"Example config not found: {config_path}")

        config_loader = QualityGateConfigLoader(config_path)

        assert config_loader.config.architecture.enabled is True
        # Layered architecture has custom layer names
        layer_names = {layer.name for layer in config_loader.config.architecture.layer_patterns}
        assert "business" in layer_names or "persistence" in layer_names
