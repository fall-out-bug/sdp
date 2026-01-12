#!/usr/bin/env python3
"""
SDP (Spec-Driven Protocol) Interactive Setup

Helps integrate SDP into your project with guided setup and reminders.
"""

import os
import sys
from pathlib import Path
from typing import Optional


class Color:
    """ANSI color codes"""
    BLUE = '\033[94m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    BOLD = '\033[1m'
    END = '\033[0m'


def print_header(text: str) -> None:
    """Print formatted header"""
    print(f"\n{Color.BOLD}{Color.BLUE}{'='*60}{Color.END}")
    print(f"{Color.BOLD}{Color.BLUE}{text:^60}{Color.END}")
    print(f"{Color.BOLD}{Color.BLUE}{'='*60}{Color.END}\n")


def print_success(text: str) -> None:
    """Print success message"""
    print(f"{Color.GREEN}✓ {text}{Color.END}")


def print_warning(text: str) -> None:
    """Print warning message"""
    print(f"{Color.YELLOW}⚠ {text}{Color.END}")


def print_error(text: str) -> None:
    """Print error message"""
    print(f"{Color.RED}✗ {text}{Color.END}")


def print_info(text: str) -> None:
    """Print info message"""
    print(f"{Color.BOLD}{text}{Color.END}")


def ask_yes_no(question: str, default: bool = True) -> bool:
    """Ask yes/no question"""
    default_str = "Y/n" if default else "y/N"
    while True:
        response = input(f"{question} [{default_str}]: ").strip().lower()
        if not response:
            return default
        if response in ('y', 'yes'):
            return True
        if response in ('n', 'no'):
            return False
        print_warning("Please answer 'y' or 'n'")


def ask_text(question: str, default: str = "") -> str:
    """Ask text question"""
    default_str = f" [{default}]" if default else ""
    response = input(f"{question}{default_str}: ").strip()
    return response or default


def ask_choice(question: str, choices: list[str], default: int = 0) -> str:
    """Ask multiple choice question"""
    print(f"\n{question}")
    for i, choice in enumerate(choices, 1):
        default_mark = " (default)" if i-1 == default else ""
        print(f"  {i}. {choice}{default_mark}")
    
    while True:
        response = input(f"Enter choice [1-{len(choices)}]: ").strip()
        if not response:
            return choices[default]
        try:
            idx = int(response) - 1
            if 0 <= idx < len(choices):
                return choices[idx]
        except ValueError:
            pass
        print_warning(f"Please enter a number between 1 and {len(choices)}")


def check_existing_files(project_root: Path) -> dict[str, bool]:
    """Check which SDP files already exist"""
    checks = {
        '.cursorrules': project_root / '.cursorrules',
        '.cursor/': project_root / '.cursor',
        '.claudecode/': project_root / '.claudecode',
        'CLAUDE.md': project_root / 'CLAUDE.md',
        'prompts/': project_root / 'prompts',
        'schema/': project_root / 'schema',
        'hooks/': project_root / 'hooks',
        'PROJECT_CONVENTIONS.md': project_root / 'PROJECT_CONVENTIONS.md',
    }
    return {name: path.exists() for name, path in checks.items()}


def welcome() -> None:
    """Print welcome message"""
    print_header("SDP Interactive Setup")
    print("This wizard will help you integrate Spec-Driven Protocol (SDP)")
    print("into your project with guided setup and reminders.")
    print("\nYou can cancel at any time with Ctrl+C")


def detect_ide() -> str:
    """Detect which IDE to configure"""
    print_header("IDE Detection")
    
    if ask_yes_no("Are you using Cursor IDE?", default=True):
        return "cursor"
    elif ask_yes_no("Are you using Claude Code?", default=False):
        return "claude"
    else:
        return ask_choice(
            "Which IDE will you use?",
            ["Cursor IDE", "Claude Code", "Both", "Other/Skip"],
            default=0
        ).lower().split()[0]


def configure_project_conventions() -> dict:
    """Gather project-specific conventions"""
    print_header("Project Conventions")
    
    print("\nThese settings will be saved to PROJECT_CONVENTIONS.md")
    print("AI agents will follow these rules when working on your project.")
    
    return {
        'language': ask_choice(
            "Primary code language",
            ["English", "Russian", "Spanish", "Other"],
            default=0
        ),
        'python_line_length': ask_text("Python max line length", "88"),
        'test_coverage': ask_text("Minimum test coverage (%)", "80"),
        'max_file_loc': ask_text("Maximum lines per file", "200"),
        'max_complexity': ask_text("Maximum cyclomatic complexity", "10"),
    }


def create_directories(project_root: Path, ide: str) -> None:
    """Create necessary directory structure"""
    print_header("Creating Directories")
    
    dirs_to_create = [
        'docs/drafts',
        'docs/workstreams/backlog',
        'docs/workstreams/in_progress',
        'docs/workstreams/completed',
        'docs/specs',
    ]
    
    if ide in ('cursor', 'both'):
        dirs_to_create.append('.cursor/commands')
    
    if ide in ('claude', 'both'):
        dirs_to_create.extend([
            '.claudecode/skills',
            '.claudecode/agents',
        ])
    
    for dir_path in dirs_to_create:
        full_path = project_root / dir_path
        if not full_path.exists():
            full_path.mkdir(parents=True, exist_ok=True)
            print_success(f"Created {dir_path}/")
        else:
            print_info(f"Already exists: {dir_path}/")


def copy_templates(project_root: Path, sdp_root: Path) -> None:
    """Copy SDP files to project"""
    print_header("Copying SDP Files")
    
    import shutil
    
    files_to_copy = [
        ('prompts/', 'prompts/'),
        ('schema/', 'schema/'),
        ('hooks/', 'hooks/'),
        ('templates/PROJECT_CONVENTIONS.md', 'PROJECT_CONVENTIONS.md'),
        ('templates/workstream.md', 'templates/workstream.md'),
    ]
    
    for src, dst in files_to_copy:
        src_path = sdp_root / src
        dst_path = project_root / dst
        
        if src_path.is_dir():
            if not dst_path.exists():
                shutil.copytree(src_path, dst_path)
                print_success(f"Copied {src}")
            else:
                print_info(f"Already exists: {dst}")
        else:
            if not dst_path.exists():
                dst_path.parent.mkdir(parents=True, exist_ok=True)
                shutil.copy2(src_path, dst_path)
                print_success(f"Copied {dst}")
            else:
                print_info(f"Already exists: {dst}")


def setup_git_hooks(project_root: Path) -> None:
    """Set up Git hooks"""
    print_header("Git Hooks Setup")
    
    if not (project_root / '.git').exists():
        print_warning("Not a Git repository. Skipping Git hooks.")
        return
    
    if not ask_yes_no("Install Git hooks for validation?", default=True):
        print_info("Skipped Git hooks installation")
        return
    
    hooks = [
        'pre-commit',
        'commit-msg',
        'pre-build',
        'post-build',
        'pre-deploy',
        'post-oneshot',
    ]
    
    for hook in hooks:
        src = project_root / 'hooks' / f'{hook}.sh'
        dst = project_root / '.git' / 'hooks' / hook
        
        if src.exists():
            dst.parent.mkdir(parents=True, exist_ok=True)
            # Create symlink instead of copy for easier updates
            if dst.exists():
                dst.unlink()
            os.symlink(src.resolve(), dst)
            dst.chmod(0o755)
            print_success(f"Installed {hook} hook")


def print_next_steps(project_root: Path, ide: str, conventions: dict) -> None:
    """Print post-installation instructions"""
    print_header("Setup Complete!")
    
    print_success("SDP has been integrated into your project\n")
    
    print_info("📋 NEXT STEPS:\n")
    
    print("1. Fill out PROJECT_CONVENTIONS.md:")
    print(f"   {Color.YELLOW}nano PROJECT_CONVENTIONS.md{Color.END}")
    print("   → Add project-specific DO/DON'T rules")
    print("   → Define naming conventions")
    print("   → Specify code style preferences\n")
    
    if ide == 'cursor':
        print("2. Review .cursorrules:")
        print(f"   {Color.YELLOW}cat .cursorrules{Color.END}")
        print("   → Cursor will auto-load this file\n")
        
        print("3. Try your first command:")
        print(f"   {Color.YELLOW}/idea \"Your feature description\"{Color.END}\n")
    
    elif ide == 'claude':
        print("2. Review CLAUDE.md:")
        print(f"   {Color.YELLOW}cat CLAUDE.md{Color.END}")
        print("   → Claude Code will auto-load this file\n")
        
        print("3. Try your first skill:")
        print(f"   {Color.YELLOW}@idea \"Your feature description\"{Color.END}\n")
    
    elif ide == 'both':
        print("2. Review IDE configs:")
        print(f"   {Color.YELLOW}cat .cursorrules{Color.END}")
        print(f"   {Color.YELLOW}cat CLAUDE.md{Color.END}\n")
        
        print("3. Try commands:")
        print(f"   Cursor: {Color.YELLOW}/idea \"Feature description\"{Color.END}")
        print(f"   Claude: {Color.YELLOW}@idea \"Feature description\"{Color.END}\n")
    
    print("4. Read core documentation:")
    print(f"   {Color.YELLOW}README.md{Color.END} - Quick start")
    print(f"   {Color.YELLOW}PROTOCOL.md{Color.END} - Full specification")
    print(f"   {Color.YELLOW}docs/PRINCIPLES.md{Color.END} - Core principles\n")
    
    print("5. Install Git hooks (if skipped):")
    print(f"   {Color.YELLOW}python scripts/init.py --install-hooks{Color.END}\n")
    
    print_info("📚 RESOURCES:\n")
    print("  • Cursor guide: docs/guides/CURSOR.md")
    print("  • Claude Code guide: docs/guides/CLAUDE_CODE.md")
    print("  • Code patterns: CODE_PATTERNS.md")
    print("  • Model recommendations: MODELS.md")
    
    print(f"\n{Color.GREEN}Happy coding with SDP! 🚀{Color.END}\n")


def main() -> None:
    """Main setup wizard"""
    try:
        welcome()
        
        # Detect project root
        project_root = Path.cwd()
        print(f"Project root: {project_root}\n")
        
        # Find SDP root (where this script is)
        sdp_root = Path(__file__).parent.parent.resolve()
        
        # Check existing files
        existing = check_existing_files(project_root)
        if any(existing.values()):
            print_warning("Some SDP files already exist:")
            for name, exists in existing.items():
                if exists:
                    print(f"  • {name}")
            print()
            if not ask_yes_no("Continue anyway?", default=True):
                print("Setup cancelled.")
                return
        
        # Detect IDE
        ide = detect_ide()
        
        # Gather conventions
        conventions = configure_project_conventions()
        
        # Create directories
        create_directories(project_root, ide)
        
        # Copy templates
        copy_templates(project_root, sdp_root)
        
        # Setup Git hooks
        setup_git_hooks(project_root)
        
        # Print next steps
        print_next_steps(project_root, ide, conventions)
        
    except KeyboardInterrupt:
        print("\n\nSetup cancelled by user.")
        sys.exit(1)
    except Exception as e:
        print_error(f"Setup failed: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
