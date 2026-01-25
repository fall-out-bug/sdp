"""Escalation metrics tracking and analysis for T2/T3 workstreams."""

import json
from dataclasses import dataclass
from datetime import datetime, timedelta
from pathlib import Path
from typing import Optional


@dataclass
class EscalationEvent:
    """Record of human escalation event.

    Args:
        ws_id: Workstream ID that escalated
        tier: Capability tier (T2 or T3)
        attempts: Number of failed attempts before escalation
        timestamp: When escalation occurred
        diagnostics: Diagnostic info for human
        feature_id: Optional feature ID
    """

    ws_id: str
    tier: str
    attempts: int
    timestamp: datetime
    diagnostics: str
    feature_id: Optional[str] = None

    def __post_init__(self) -> None:
        """Validate escalation event fields."""
        if self.tier not in ("T2", "T3"):
            raise ValueError(f"Invalid tier for escalation: {self.tier}")
        if self.attempts <= 0:
            raise ValueError(f"Attempts must be positive: {self.attempts}")


class EscalationMetricsStore:
    """Store and analyze escalation metrics.

    Provides persistent storage for tracking human escalation events
    from T2/T3 workstreams, with analysis capabilities for monitoring.
    """

    def __init__(self, storage_path: Path = Path(".sdp/escalation_metrics.json")) -> None:
        """Initialize escalation metrics store.

        Args:
            storage_path: Path to JSON file for persistence
        """
        self.storage_path = storage_path
        self.storage_path.parent.mkdir(parents=True, exist_ok=True)
        self._events: list[EscalationEvent] = []
        self._load()

    def _load(self) -> None:
        """Load events from storage file."""
        if not self.storage_path.exists():
            self._events = []
            return

        try:
            with open(self.storage_path, encoding="utf-8") as f:
                data = json.load(f)
                self._events = []
                for event_data in data:
                    # Parse datetime from ISO format
                    event_data["timestamp"] = datetime.fromisoformat(event_data["timestamp"])
                    self._events.append(EscalationEvent(**event_data))
        except (json.JSONDecodeError, KeyError, ValueError):
            # If file is corrupted, start fresh
            self._events = []

    def _save(self) -> None:
        """Save events to storage file."""
        data = [
            {
                "ws_id": event.ws_id,
                "tier": event.tier,
                "attempts": event.attempts,
                "timestamp": event.timestamp.isoformat(),
                "diagnostics": event.diagnostics,
                "feature_id": event.feature_id,
            }
            for event in self._events
        ]
        with open(self.storage_path, "w", encoding="utf-8") as f:
            json.dump(data, f, indent=2)

    def record_escalation(self, event: EscalationEvent) -> None:
        """Record escalation event.

        Args:
            event: Escalation event to record
        """
        self._events.append(event)
        self._save()

    def get_events(
        self, tier: Optional[str] = None, days: int = 7
    ) -> list[EscalationEvent]:
        """Get escalation events filtered by tier and time window.

        Args:
            tier: Filter by tier (None = all tiers)
            days: Time window in days

        Returns:
            List of escalation events
        """
        cutoff = datetime.now() - timedelta(days=days)

        events = self._events
        if tier:
            events = [e for e in events if e.tier == tier]
        events = [e for e in events if e.timestamp >= cutoff]

        return events

    def get_escalation_count(
        self, tier: Optional[str] = None, days: int = 7
    ) -> int:
        """Get number of escalations in time window.

        Args:
            tier: Filter by tier (None = all tiers)
            days: Time window in days

        Returns:
            Number of escalation events
        """
        return len(self.get_events(tier, days))

    def get_escalation_rate(
        self, tier: Optional[str] = None, days: int = 7, total_builds: int = 20
    ) -> float:
        """Calculate escalation rate.

        Args:
            tier: Filter by tier (None = all tiers)
            days: Time window in days
            total_builds: Total number of builds in period

        Returns:
            Escalation rate as fraction (0.0 - 1.0)
        """
        escalation_count = self.get_escalation_count(tier, days)
        if total_builds == 0:
            return 0.0
        return escalation_count / total_builds

    def get_top_escalating_ws(self, limit: int = 10, days: int = 7) -> list[tuple[str, int]]:
        """Get workstreams with most escalations.

        Args:
            limit: Max number of results
            days: Time window in days

        Returns:
            List of (ws_id, escalation_count) tuples
        """
        events = self.get_events(days=days)

        # Count escalations per workstream
        counts: dict[str, int] = {}
        for event in events:
            counts[event.ws_id] = counts.get(event.ws_id, 0) + 1

        # Sort by count descending
        sorted_counts = sorted(counts.items(), key=lambda x: x[1], reverse=True)

        return sorted_counts[:limit]

    def get_average_attempts(self, tier: Optional[str] = None, days: int = 7) -> float:
        """Get average number of attempts before escalation.

        Args:
            tier: Filter by tier (None = all tiers)
            days: Time window in days

        Returns:
            Average attempts, or 0.0 if no escalations
        """
        events = self.get_events(tier, days)
        if not events:
            return 0.0

        total_attempts = sum(event.attempts for event in events)
        return total_attempts / len(events)
