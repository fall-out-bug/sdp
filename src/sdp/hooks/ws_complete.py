"""Post-WS-Complete Hook - Verifies completion before status change."""

import sys
from dataclasses import dataclass
from pathlib import Path

import click

from sdp.validators.ws_completion import WSCompletionVerifier


@dataclass
class HookResult:
    """Result of hook execution."""

    passed: bool
    ws_id: str
    messages: list[str]
    bypass_used: bool
    bypass_reason: str | None


class PostWSCompleteHook:
    """Hook to verify WS completion before status change."""

    def __init__(self, verifier: WSCompletionVerifier):
        """Initialize hook.

        Args:
            verifier: Completion verifier instance
        """
        self.verifier = verifier

    def _handle_bypass(self, ws_id: str, reason: str) -> HookResult:
        """Handle bypass case."""
        if not reason:
            return HookResult(
                passed=False,
                ws_id=ws_id,
                messages=["❌ Bypass requires --reason argument"],
                bypass_used=False,
                bypass_reason=None,
            )
        messages = [
            f"⚠️  Bypassing verification for {ws_id}",
            f"   Reason: {reason}",
        ]
        self._log_bypass(ws_id, reason)
        return HookResult(
            passed=True, ws_id=ws_id, messages=messages, bypass_used=True, bypass_reason=reason
        )

    def _handle_verification_passed(
        self, ws_id: str, result
    ) -> HookResult:
        """Build result for passed verification."""
        messages = [
            f"🔍 Verifying {ws_id} completion...",
            f"✅ Verification PASSED for {ws_id}",
            f"   - Checks: {len(result.checks)}/{len(result.checks)} passed",
        ]
        if result.coverage_actual:
            messages.append(f"   - Coverage: {result.coverage_actual:.1f}%")
        return HookResult(
            passed=True, ws_id=ws_id, messages=messages, bypass_used=False, bypass_reason=None
        )

    def _handle_verification_failed(self, ws_id: str, result) -> HookResult:
        """Build result for failed verification."""
        failed_checks = [c for c in result.checks if not c.passed]
        messages = [
            f"🔍 Verifying {ws_id} completion...",
            f"❌ Verification FAILED for {ws_id}",
            f"   - Failed checks: {len(failed_checks)}/{len(result.checks)}",
        ]
        for check in failed_checks[:5]:
            messages.append(f"     • {check.name}: {check.message}")
        if len(failed_checks) > 5:
            messages.append(f"     ... and {len(failed_checks) - 5} more")
        messages.append("\n📋 Remediation steps:")
        if result.missing_files:
            messages.append("   1. Create missing files:")
            for f in result.missing_files[:3]:
                messages.append(f"      - {f}")
        if result.failed_commands:
            messages.append("   2. Fix failing commands:")
            for cmd in result.failed_commands[:3]:
                messages.append(f"      - {cmd}")
        if result.coverage_actual and result.coverage_actual < 80.0:
            messages.append("   3. Increase test coverage:")
            messages.append(f"      Current: {result.coverage_actual:.1f}%, Target: 80%")
        messages.append("\n💡 Or use --bypass with --reason if verification is incorrect")
        return HookResult(
            passed=False, ws_id=ws_id, messages=messages, bypass_used=False, bypass_reason=None
        )

    def run(self, ws_id: str, bypass: bool = False, reason: str = "") -> HookResult:
        """Run verification and return result.

        Args:
            ws_id: Workstream to verify
            bypass: Skip verification (requires reason)
            reason: Bypass justification (logged)

        Returns:
            HookResult with pass/fail and messages
        """
        if bypass:
            return self._handle_bypass(ws_id, reason)

        result = self.verifier.verify(ws_id)
        if result.passed:
            return self._handle_verification_passed(ws_id, result)
        return self._handle_verification_failed(ws_id, result)

    def _log_bypass(self, ws_id: str, reason: str) -> None:
        """Log bypass event for auditing.

        Args:
            ws_id: Workstream ID
            reason: Bypass reason
        """
        from datetime import datetime, timezone

        log_file = Path(".sdp") / "bypass_log.txt"
        log_file.parent.mkdir(parents=True, exist_ok=True)

        timestamp = datetime.now(timezone.utc).isoformat()
        log_entry = f"{timestamp} | {ws_id} | BYPASS | {reason}\n"

        with log_file.open("a", encoding="utf-8") as f:
            f.write(log_entry)


@click.command()
@click.argument("ws_id")
@click.option("--bypass", is_flag=True, default=False, help="Skip verification (requires --reason)")
@click.option("--reason", default="", help="Bypass justification (logged for audit)")
def main(ws_id: str, bypass: bool, reason: str) -> None:
    """Post-WS-Complete Hook - Verify completion before status change.

    Args:
        ws_id: Workstream ID to verify
        bypass: Skip verification flag
        reason: Bypass reason
    """
    verifier = WSCompletionVerifier()
    hook = PostWSCompleteHook(verifier)

    result = hook.run(ws_id, bypass=bypass, reason=reason)

    # Print messages
    for msg in result.messages:
        click.echo(msg)

    # Exit with appropriate code
    sys.exit(0 if result.passed else 1)


if __name__ == "__main__":
    main()
