"""Skill management CLI commands."""
import re
from pathlib import Path

import click


@click.group()
def skill() -> None:
    """Skill management commands."""
    pass


REQUIRED_SECTIONS = [
    "## Quick Reference",
    "## Workflow",
    "## See Also",
]


def _check_line_count(lines: list[str]) -> tuple[list[str], list[str]]:
    """Check line count. Returns (errors, warnings)."""
    errors, warnings = [], []
    line_count = len(lines)
    if line_count > 150:
        errors.append(f"Too long: {line_count} lines (max 150)")
    elif line_count > 100:
        warnings.append(f"Consider shortening: {line_count} lines (target 100)")
    return errors, warnings


def _check_sections(content: str) -> list[str]:
    """Check required sections present."""
    errors = []
    for section in REQUIRED_SECTIONS:
        if section not in content:
            errors.append(f"Missing section: {section}")
    return errors


def _check_refs(content: str, path: Path) -> list[str]:
    """Check references resolve."""
    warnings = []
    refs = re.findall(r"\[.*?\]\((\.\.?/[^)]+)\)", content)
    for ref in refs:
        ref_path = path.parent / ref
        if not ref_path.exists() and not ref.startswith("../docs/"):
            warnings.append(f"Reference may not exist: {ref}")
    return warnings


def _output_validation(path: Path, errors: list[str], warnings: list[str], line_count: int) -> None:
    """Output validation results."""
    if errors:
        click.echo(f"❌ {path.name}: {len(errors)} errors")
        for e in errors:
            click.echo(f"   - {e}")
    if warnings:
        click.echo(f"⚠️  {path.name}: {len(warnings)} warnings")
        for w in warnings:
            click.echo(f"   - {w}")
    if not errors and not warnings:
        click.echo(f"✅ {path.name}: valid ({line_count} lines)")


@skill.command("validate")
@click.argument("path", type=click.Path(exists=True, path_type=Path))
@click.option("--strict", is_flag=True, help="Fail on warnings")
def validate_skill(path: Path, strict: bool) -> None:
    """Validate skill file against standards."""
    content = path.read_text()
    lines = content.splitlines()
    line_count = len(lines)

    errors, warnings = _check_line_count(lines)
    errors.extend(_check_sections(content))
    if not content.startswith("---"):
        errors.append("Missing frontmatter (must start with ---)")
    warnings.extend(_check_refs(content, path))

    _output_validation(path, errors, warnings, line_count)
    if errors or (strict and warnings):
        raise click.Abort()


@skill.command("check-all")
def check_all_skills() -> None:
    """Validate all skills in .claude/skills/."""
    skills_dir = Path(".claude/skills")
    if not skills_dir.exists():
        click.echo("❌ No .claude/skills/ directory")
        raise click.Abort()

    total = 0
    failed = 0

    for skill_dir in skills_dir.iterdir():
        if skill_dir.is_dir():
            skill_file = skill_dir / "SKILL.md"
            if skill_file.exists():
                total += 1
                try:
                    ctx = click.Context(validate_skill)
                    ctx.invoke(validate_skill, path=skill_file, strict=False)
                except click.Abort:
                    failed += 1

    click.echo(f"\nSummary: {total - failed}/{total} skills valid")
    if failed:
        raise click.Abort()
