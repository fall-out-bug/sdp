"""Extension CLI commands for SDP.

Commands:
- list: Show all discovered extensions
- show: Display extension details
- create: Scaffold new extension
"""

from pathlib import Path

import click

from sdp.extensions import ExtensionLoader


@click.group()
def extension() -> None:
    """Manage SDP extensions."""
    pass


@extension.command("list")
def list_extensions() -> None:
    """List all discovered extensions.
    
    Example:
        $ sdp extension list
        ╭──────────────┬─────────┬──────────────────────────────────╮
        │ Name         │ Version │ Description                      │
        ├──────────────┼─────────┼──────────────────────────────────┤
        │ hw_checker   │ 1.0.0   │ Clean Architecture validation    │
        │ github       │ 1.0.0   │ GitHub integration config        │
        ╰──────────────┴─────────┴──────────────────────────────────╯
    """
    loader = ExtensionLoader()
    extensions = loader.discover_extensions()
    
    if not extensions:
        click.echo("No extensions found.")
        click.echo("")
        click.echo("Extensions are discovered from:")
        click.echo("  - Project-local: sdp.local/")
        click.echo("  - User-global: ~/.sdp/extensions/")
        return
    
    # Print table header
    click.echo("╭" + "─" * 14 + "┬" + "─" * 9 + "┬" + "─" * 50 + "╮")
    click.echo("│ Name         │ Version │ Description                                      │")
    click.echo("├" + "─" * 14 + "┼" + "─" * 9 + "┼" + "─" * 50 + "┤")
    
    # Print extensions
    for ext in sorted(extensions, key=lambda e: e.manifest.name):
        name = ext.manifest.name[:12].ljust(12)
        version = ext.manifest.version[:7].ljust(7)
        description = ext.manifest.description[:48].ljust(48)
        click.echo(f"│ {name} │ {version} │ {description} │")
    
    # Print table footer
    click.echo("╰" + "─" * 14 + "┴" + "─" * 9 + "┴" + "─" * 50 + "╯")
    click.echo(f"\nTotal: {len(extensions)} extension(s)")


@extension.command("show")
@click.argument("name")
def show_extension(name: str) -> None:
    """Show detailed information about an extension.
    
    Args:
        name: Extension name
    
    Example:
        $ sdp extension show hw_checker
        Name: hw_checker
        Version: 1.0.0
        Description: Clean Architecture validation
        Author: SDP Team
        Location: /home/user/project/sdp.local/hw-checker
        
        Directories:
          ✓ hooks: 1 file
          ✓ patterns: 1 file
          ✓ skills: 0 files
          ✓ integrations: 0 files
    """
    loader = ExtensionLoader()
    extensions = loader.discover_extensions()
    
    # Find extension by name
    ext = next((e for e in extensions if e.manifest.name == name), None)
    
    if not ext:
        click.echo(f"Error: Extension '{name}' not found.", err=True)
        click.echo(f"\nAvailable extensions:", err=True)
        for e in extensions:
            click.echo(f"  - {e.manifest.name}", err=True)
        raise click.Abort()
    
    # Display extension details
    click.echo(f"Name: {ext.manifest.name}")
    click.echo(f"Version: {ext.manifest.version}")
    click.echo(f"Description: {ext.manifest.description}")
    click.echo(f"Author: {ext.manifest.author}")
    click.echo(f"Location: {ext.root_path}")
    click.echo("")
    click.echo("Directories:")
    
    # Count files in each directory
    def count_files(path: Path | None) -> int:
        if path is None or not path.exists():
            return 0
        return sum(1 for _ in path.rglob("*") if _.is_file())
    
    hooks_path = ext.get_hooks_path()
    hooks_count = count_files(hooks_path)
    hooks_status = "✓" if hooks_path else "✗"
    click.echo(f"  {hooks_status} hooks: {hooks_count} file(s)")
    
    patterns_path = ext.get_patterns_path()
    patterns_count = count_files(patterns_path)
    patterns_status = "✓" if patterns_path else "✗"
    click.echo(f"  {patterns_status} patterns: {patterns_count} file(s)")
    
    skills_path = ext.get_skills_path()
    skills_count = count_files(skills_path)
    skills_status = "✓" if skills_path else "✗"
    click.echo(f"  {skills_status} skills: {skills_count} file(s)")
    
    integrations_path = ext.get_integrations_path()
    integrations_count = count_files(integrations_path)
    integrations_status = "✓" if integrations_path else "✗"
    click.echo(f"  {integrations_status} integrations: {integrations_count} file(s)")


@extension.command("create")
@click.argument("name")
@click.option(
    "--location",
    type=click.Choice(["local", "global"], case_sensitive=False),
    default="local",
    help="Extension location (local: sdp.local/, global: ~/.sdp/extensions/)",
)
def create_extension(name: str, location: str) -> None:
    """Create a new extension scaffold.
    
    Args:
        name: Extension name (alphanumeric + underscores)
        location: Where to create extension (local or global)
    
    Example:
        $ sdp extension create my-extension
        Created: sdp.local/my-extension/
          ├── extension.yaml
          ├── hooks/
          ├── patterns/
          ├── skills/
          └── integrations/
    """
    # Validate name
    if not name.replace("_", "").replace("-", "").isalnum():
        click.echo(
            f"Error: Extension name must be alphanumeric (with _ or -): {name}",
            err=True,
        )
        raise click.Abort()
    
    # Normalize name (replace - with _)
    normalized_name = name.replace("-", "_")
    
    # Determine extension root
    if location == "local":
        ext_root = Path.cwd() / "sdp.local" / normalized_name
    else:
        ext_root = Path.home() / ".sdp" / "extensions" / normalized_name
    
    if ext_root.exists():
        click.echo(f"Error: Extension already exists: {ext_root}", err=True)
        raise click.Abort()
    
    # Create directory structure
    ext_root.mkdir(parents=True)
    (ext_root / "hooks").mkdir()
    (ext_root / "patterns").mkdir()
    (ext_root / "skills").mkdir()
    (ext_root / "integrations").mkdir()
    
    # Create extension.yaml
    manifest_content = f"""name: {normalized_name}
version: 0.1.0
description: Add your extension description here
author: Your Name

# Directory paths (relative to extension root)
hooks_dir: hooks
patterns_dir: patterns
skills_dir: skills
integrations_dir: integrations
"""
    
    (ext_root / "extension.yaml").write_text(manifest_content)
    
    # Create README
    readme_content = f"""# {normalized_name.title().replace('_', ' ')} Extension

Brief description of what this extension provides.

## Features

- Feature 1: [Describe feature]
- Feature 2: [Describe feature]
- Feature 3: [Describe feature]

## Usage

[Add usage instructions here]

## Configuration

[Document configuration options here]
"""
    
    (ext_root / "patterns" / "README.md").write_text(readme_content)
    
    # Success message
    click.echo(f"Created: {ext_root}/")
    click.echo("  ├── extension.yaml")
    click.echo("  ├── hooks/")
    click.echo("  ├── patterns/")
    click.echo("  │   └── README.md")
    click.echo("  ├── skills/")
    click.echo("  └── integrations/")
    click.echo("")
    click.echo("Next steps:")
    click.echo(f"  1. Edit {ext_root}/extension.yaml")
    click.echo(f"  2. Add content to directories")
    click.echo(f"  3. Test: sdp extension show {normalized_name}")
