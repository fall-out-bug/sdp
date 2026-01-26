"""Affected test detection for watch mode."""

from pathlib import Path


def detect_affected_tests(changed_file: str, project_dir: str | Path) -> list[str]:
    """Detect tests affected by a changed file.

    Args:
        changed_file: Path to changed file
        project_dir: Project root directory

    Returns:
        List of test file paths to run
    """
    project_path = Path(project_dir)
    changed_path = Path(changed_file).relative_to(project_path)

    affected = []

    # If a test file changed, run just that file
    if "tests" in changed_path.parts:
        return [str(changed_path)]

    # If a source file changed, find corresponding test files
    if "src" in changed_path.parts:
        module_name = changed_path.stem
        tests_dir = project_path / "tests"

        if not tests_dir.exists():
            return []

        # Look for test files with matching module name
        for test_file in tests_dir.glob(f"**/test_{module_name}.py"):
            affected.append(str(test_file))

        # Also look for test files in same directory structure
        src_parts = list(changed_path.parts)
        if "src" in src_parts:
            src_index = src_parts.index("src")
            relative_parts = src_parts[src_index + 1:]

            # Convert src/path/to/module.py -> tests/path/to/test_module.py
            test_path_parts = ["tests"] + list(relative_parts[:-1])
            test_path_parts[0] = "tests"  # Ensure tests dir

            test_path = project_path / Path(*test_path_parts) / f"test_{module_name}.py"
            if test_path.exists():
                affected.append(str(test_path))

    return affected
