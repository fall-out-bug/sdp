"""SDP init command - Interactive project setup wizard.

Creates standard directory structure, configuration, and validates setup.
"""

from pathlib import Path

import click

from sdp.init_wizard import (
    collect_metadata,
    create_env_template,
    create_structure,
    detect_dependencies,
    generate_quality_gate,
    install_git_hooks,
    run_doctor,
    show_dependencies,
)


@click.command()
@click.option(
    "--force",
    is_flag=True,
    help="Overwrite existing files",
)
@click.option(
    "--path",
    type=click.Path(file_okay=False, dir_okay=True, path_type=Path),
    default=None,
    help="Target directory (defaults to current directory)",
)
@click.option(
    "--non-interactive",
    is_flag=True,
    help="Run non-interactively (use defaults)",
)
def init(force: bool, path: Path | None, non_interactive: bool) -> None:
    """Initialize SDP in current project (interactive wizard).

    Creates standard directory structure:
    - docs/workstreams/ (INDEX.md, TEMPLATE.md, backlog/)
    - docs/PROJECT_MAP.md
    - quality-gate.toml (configurable quality gates)
    - .env.template (environment variables)
    - .git/hooks/pre-commit (SDP validation hooks)

    Optional dependencies (auto-detected):
    - Beads CLI (task tracking)
    - GitHub CLI (GitHub integration)
    - Telegram (notifications)

    Example:
        $ sdp init
        Project name [my-project]: my-awesome-project
        Description [SDP project]: A project using Spec-Driven Protocol
        Author [Your Name]: John Doe
        ...
        ✓ Created docs/workstreams/INDEX.md
        ✓ Created quality-gate.toml
        ✓ Created .env.template
        ✓ Installed git hooks
        ✓ Ran sdp doctor - all checks passed

        SDP initialized! Next steps:
        1. Edit docs/PROJECT_MAP.md with your project info
        2. Run: sdp extension list
        3. Start: /idea "your first feature"
    """
    # Determine target directory
    target_dir = path.resolve() if path else Path.cwd()

    # Create target directory if it doesn't exist
    if not target_dir.exists():
        target_dir.mkdir(parents=True)

    click.echo(click.style("🚀 SDP Project Setup Wizard", fg="cyan", bold=True))
    click.echo("=" * 50)
    click.echo()

    # Step 1: Collect project metadata
    project_name, _description, _author = collect_metadata(
        target_dir, non_interactive
    )

    # Step 2: Detect optional dependencies
    deps = detect_dependencies()
    show_dependencies(deps)

    # Step 3: Create directory structure
    click.echo()
    click.echo(click.style("Step 3: Creating directory structure...", fg="cyan"))
    created_files, skipped_files = create_structure(
        target_dir, project_name, force
    )

    # Step 4: Generate quality gate configuration
    click.echo()
    click.echo(click.style("Step 4: Generating quality-gate.toml...", fg="cyan"))
    quality_gate_file = generate_quality_gate(target_dir, deps)
    if quality_gate_file:
        created_files.append(str(quality_gate_file))

    # Step 5: Create .env template
    click.echo()
    click.echo(click.style("Step 5: Creating .env.template...", fg="cyan"))
    env_template = create_env_template(target_dir, deps)
    if env_template:
        created_files.append(str(env_template))

    # Step 6: Install git hooks
    click.echo()
    click.echo(click.style("Step 6: Installing git hooks...", fg="cyan"))
    hooks_installed = install_git_hooks(target_dir)
    if hooks_installed:
        click.echo("✓ Git hooks installed")
    else:
        click.echo("⊘ Git hooks skipped (not a git repository)")

    # Step 7: Run sdp doctor
    click.echo()
    click.echo(click.style("Step 7: Running sdp doctor for validation...", fg="cyan"))
    doctor_passed = run_doctor(target_dir)

    # Display results
    click.echo()
    click.echo(click.style("=" * 50, bold=True))
    if created_files:
        for file in created_files:
            click.echo(click.style(f"✓ Created {file}", fg="green"))

    if skipped_files:
        click.echo()
        click.echo("Skipped (already exists, use --force to overwrite):")
        for file in skipped_files:
            click.echo(f"  ⊘ {file}")

    # Final summary
    click.echo()
    click.echo(click.style("Setup Summary", bold=True))
    click.echo(f"  Project: {project_name}")
    click.echo(f"  Location: {target_dir}")
    click.echo(f"  Files created: {len(created_files)}")
    click.echo(f"  Files skipped: {len(skipped_files)}")
    click.echo(f"  Dependencies: {len([d for d in deps if deps[d]])} detected")

    if doctor_passed:
        click.echo(click.style("  Health check: ✓ PASSED", fg="green"))
    else:
        click.echo(click.style("  Health check: ⚠ WARNING", fg="yellow"))

    click.echo()
    click.echo(click.style("SDP initialized!", fg="green", bold=True))
    click.echo("Next steps:")
    click.echo(f"  1. Edit {target_dir / 'docs' / 'PROJECT_MAP.md'} with your project info")
    click.echo("  2. Run: sdp extension list")
    click.echo("  3. Start: /idea \"your first feature\"")
