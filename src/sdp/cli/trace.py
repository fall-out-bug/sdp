"""Traceability CLI commands."""

import json
import sys

import click

from sdp.beads import create_beads_client
from sdp.traceability.service import TraceabilityService


@click.group()
def trace() -> None:
    """AC→Test traceability commands."""
    pass


@trace.command("check")
@click.argument("ws_id")
@click.option("--json", "json_output", is_flag=True, help="JSON output")
def check_traceability(ws_id: str, json_output: bool) -> None:
    """Check AC→Test traceability for workstream.

    Args:
        ws_id: Workstream ID (e.g., "00-032-01")
        json_output: Output as JSON instead of human-readable table

    Example:
        sdp trace check 00-032-01
        sdp trace check 00-032-01 --json
    """
    service = TraceabilityService(create_beads_client())

    try:
        report = service.check_traceability(ws_id)
    except ValueError as e:
        click.echo(f"❌ {e}", err=True)
        sys.exit(1)

    if json_output:
        click.echo(json.dumps(report.to_dict(), indent=2))
    else:
        click.echo(f"\nTraceability Report: {ws_id}")
        click.echo("=" * 50)
        click.echo(report.to_markdown_table())
        click.echo("")
        click.echo(
            f"Coverage: {report.coverage_pct:.0f}% "
            f"({report.mapped_acs}/{report.total_acs} ACs mapped)"
        )

        if report.is_complete:
            click.echo("Status: ✅ COMPLETE")
        else:
            click.echo(f"Status: ❌ INCOMPLETE ({report.missing_acs} unmapped)")

    # Exit 1 if incomplete
    if not report.is_complete:
        sys.exit(1)


@trace.command("add")
@click.argument("ws_id")
@click.option("--ac", required=True, help="AC ID (e.g., AC1)")
@click.option("--test", required=True, help="Test function name")
@click.option("--file", "file_path", default="", help="Test file path")
def add_mapping(ws_id: str, ac: str, test: str, file_path: str) -> None:
    """Add AC→Test mapping.

    Args:
        ws_id: Workstream ID
        ac: Acceptance criterion ID (e.g., "AC1")
        test: Test function name
        file_path: Test file path (optional)

    Example:
        sdp trace add 00-032-01 --ac AC1 --test test_login --file tests/test_auth.py
    """
    service = TraceabilityService(create_beads_client())

    try:
        service.add_mapping(ws_id, ac.upper(), file_path, test)
        click.echo(f"✅ Mapped {ac} → {test}")
    except ValueError as e:
        click.echo(f"❌ {e}", err=True)
        sys.exit(1)


@trace.command("auto")
@click.argument("ws_id")
@click.option("--test-dir", default="tests/", help="Test directory")
@click.option("--apply", is_flag=True, help="Apply detected mappings")
def auto_detect(ws_id: str, test_dir: str, apply: bool) -> None:
    """Auto-detect AC→Test mappings using AST analysis.

    Detects mappings by:
    - Parsing test docstrings for AC references
    - Matching test function names (test_ac1_*, etc.)
    - Keyword similarity between AC descriptions and test names

    Args:
        ws_id: Workstream ID
        test_dir: Directory containing test files
        apply: Apply high-confidence mappings (≥80%)

    Example:
        sdp trace auto 00-032-01
        sdp trace auto 00-032-01 --apply
        sdp trace auto 00-032-01 --test-dir tests/unit/
    """
    from pathlib import Path

    from sdp.traceability.detector import ACDetector

    service = TraceabilityService(create_beads_client())
    detector = ACDetector()

    # Get WS and extract ACs
    try:
        report = service.check_traceability(ws_id)
    except ValueError as e:
        click.echo(f"❌ {e}", err=True)
        sys.exit(1)

    ac_descriptions = {m.ac_id: m.ac_description for m in report.mappings}

    # Detect mappings
    test_path = Path(test_dir)
    if not test_path.exists():
        click.echo(f"❌ Test directory not found: {test_dir}", err=True)
        sys.exit(1)

    detected = detector.detect_all(test_path, ac_descriptions)

    if not detected:
        click.echo("No mappings detected")
        return

    click.echo(f"\nDetected {len(detected)} potential mappings:\n")

    # Sort by confidence descending, then AC ID
    for d in sorted(detected, key=lambda x: (-x.confidence, x.ac_id)):
        # Confidence bar
        bar_len = int(d.confidence * 10)
        conf_bar = "█" * bar_len + "░" * (10 - bar_len)

        click.echo(f"  {d.ac_id} → {d.test_name}")
        click.echo(f"    File: {d.test_file}")
        click.echo(f"    Confidence: [{conf_bar}] {d.confidence:.0%} ({d.source})")
        click.echo("")

    if apply:
        applied = 0
        for d in detected:
            if d.confidence >= 0.8:  # Only apply high-confidence
                service.add_mapping(ws_id, d.ac_id, d.test_file, d.test_name)
                click.echo(f"✅ Applied: {d.ac_id} → {d.test_name}")
                applied += 1

        click.echo(f"\nApplied {applied} high-confidence mappings (≥80%)")
    else:
        high_conf = sum(1 for d in detected if d.confidence >= 0.8)
        click.echo(f"Run with --apply to save {high_conf} high-confidence mappings")
