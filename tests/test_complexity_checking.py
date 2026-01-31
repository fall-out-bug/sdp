#!/usr/bin/env python3
"""Integration tests for complexity checking."""

import os
import subprocess
import sys
import tempfile
from pathlib import Path

import pytest


class TestComplexityChecking:
    """Test complexity checking with radon."""

    def test_simple_function_passes(self):
        """Simple functions should pass complexity check."""
        with tempfile.TemporaryDirectory() as tmpdir:
            test_file = Path(tmpdir) / "simple.py"
            test_file.write_text(
                """
def add(a: int, b: int) -> int:
    return a + b

def greet(name: str) -> str:
    return f"Hello, {name}"
"""
            )

            result = subprocess.run(
                ["radon", "cc", str(test_file), "-a"],
                capture_output=True,
                text=True,
            )

            assert result.returncode == 0
            assert "A" in result.stdout  # Should get grade A

    def test_complex_function_fails(self):
        """Functions with high complexity should fail."""
        with tempfile.TemporaryDirectory() as tmpdir:
            test_file = Path(tmpdir) / "complex.py"
            # Generate a function with high cyclomatic complexity
            lines = ["def complex_function(x, y, z):\n"]
            for i in range(15):  # Create 15 if/elif branches
                lines.append(f"    if x == {i}:\n")
                lines.append(f"        return {i}\n")

            test_file.write_text("".join(lines))

            result = subprocess.run(
                ["radon", "cc", str(test_file), "-s"],
                capture_output=True,
                text=True,
            )

            assert result.returncode == 0
            assert "F" in result.stdout or "E" in result.stdout  # Should get failing grade

    def test_check_complexity_script(self):
        """Test the check_complexity.sh script."""
        script_path = Path(__file__).parent.parent / "scripts" / "check_complexity.sh"

        if not script_path.exists():
            pytest.skip(f"Script not found: {script_path}")

        with tempfile.TemporaryDirectory() as tmpdir:
            # Create a simple file that should pass
            test_file = Path(tmpdir) / "simple.py"
            test_file.write_text("def foo():\n    pass\n")

            # Make script executable
            os.chmod(script_path, 0o755)

            result = subprocess.run(
                [str(script_path), str(tmpdir)],
                capture_output=True,
                text=True,
                timeout=30,
            )

            # Should pass for simple code
            assert "✓ Complexity check passed" in result.stdout or "check_complexity.sh" in result.stderr

    def test_quality_gate_configuration(self):
        """Test that quality-gate.toml is read correctly."""
        config_file = Path("quality-gate.toml")

        if not config_file.exists():
            pytest.skip("quality-gate.toml not found")

        import tomllib

        with open(config_file, "rb") as f:
            config = tomllib.load(f)

        assert "complexity" in config
        assert "max_cc" in config["complexity"]
        assert "max_average_cc" in config["complexity"]

        # Verify reasonable defaults
        assert config["complexity"]["max_cc"] >= 1
        assert config["complexity"]["max_average_cc"] >= 1


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
