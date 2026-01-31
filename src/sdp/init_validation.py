"""Validation and installation for SDP init wizard."""

from pathlib import Path


def install_git_hooks(target_dir: Path) -> bool:
    """Install git hooks.

    Args:
        target_dir: Target directory

    Returns:
        True if hooks were installed
    """
    git_dir = target_dir / ".git"
    if not git_dir.exists():
        return False

    hooks_dir = git_dir / "hooks"
    pre_commit_hook = hooks_dir / "pre-commit"

    # Check if SDP is installed
    try:
        # Try to find SDP hooks
        sdp_hooks = Path(__file__).parent.parent.parent / "hooks" / "pre-commit.sh"
        if not sdp_hooks.exists():
            # Fallback to installed location
            import sdp
            sdp_dir = Path(sdp.__file__).parent
            sdp_hooks = sdp_dir / "hooks" / "pre-commit.sh"

        if not sdp_hooks.exists():
            return False

        # Copy hook
        import shutil
        shutil.copy(sdp_hooks, pre_commit_hook)
        pre_commit_hook.chmod(0o755)

        return True
    except Exception:
        return False


def run_doctor(target_dir: Path) -> bool:
    """Run sdp doctor for validation.

    Args:
        target_dir: Target directory

    Returns:
        True if all critical checks passed
    """
    import click

    try:
        # Import doctor command
        from sdp.health_checks import get_health_checks

        # Run health checks
        checks = get_health_checks()
        passed = 0
        failed = 0

        for check in checks:
            try:
                result = check.run()
                if result.passed:
                    passed += 1
                else:
                    failed += 1
            except Exception:
                failed += 1

        # Display summary
        if failed == 0:
            click.echo(click.style("✓ All health checks passed", fg="green"))
            return True
        else:
            # Check if any critical checks failed
            critical_failed = 0
            for check in checks:
                if check.critical:
                    try:
                        result = check.run()
                        if not result.passed:
                            critical_failed += 1
                    except Exception:
                        critical_failed += 1

            click.echo(
                click.style(
                    f"⚠ {passed} passed, {failed} failed",
                    fg="yellow"
                )
            )
            return critical_failed == 0  # True only if no critical failures

    except Exception as e:
        click.echo(click.style(f"⊘ Could not run doctor: {e}", fg="yellow"))
        return True  # Don't fail setup if doctor has issues
