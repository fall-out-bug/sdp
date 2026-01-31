"""Data models for AC→Test traceability."""

from dataclasses import dataclass, field
from enum import Enum


class MappingStatus(Enum):
    """Status of AC→Test mapping."""

    MAPPED = "mapped"  # AC has test
    MISSING = "missing"  # AC has no test
    FAILED = "failed"  # Test exists but fails


@dataclass
class ACTestMapping:
    """Maps Acceptance Criterion to Test."""

    ac_id: str  # "AC1", "AC2", etc.
    ac_description: str  # "User can login"
    test_file: str | None  # "tests/unit/test_auth.py"
    test_name: str | None  # "test_user_login"
    status: MappingStatus
    confidence: float = 1.0  # 0.0-1.0 for auto-detected

    def to_dict(self) -> dict:
        """Serialize to dictionary."""
        return {
            "ac_id": self.ac_id,
            "ac_description": self.ac_description,
            "test_file": self.test_file,
            "test_name": self.test_name,
            "status": self.status.value,
            "confidence": self.confidence,
        }

    @classmethod
    def from_dict(cls, data: dict) -> "ACTestMapping":
        """Deserialize from dictionary."""
        return cls(
            ac_id=data["ac_id"],
            ac_description=data["ac_description"],
            test_file=data.get("test_file"),
            test_name=data.get("test_name"),
            status=MappingStatus(data["status"]),
            confidence=data.get("confidence", 1.0),
        )


@dataclass
class TraceabilityReport:
    """Report of AC→Test traceability for a workstream."""

    ws_id: str
    mappings: list[ACTestMapping] = field(default_factory=list)

    @property
    def total_acs(self) -> int:
        """Total number of ACs."""
        return len(self.mappings)

    @property
    def mapped_acs(self) -> int:
        """Number of ACs with tests."""
        return sum(1 for m in self.mappings if m.status == MappingStatus.MAPPED)

    @property
    def missing_acs(self) -> int:
        """Number of ACs without tests."""
        return sum(1 for m in self.mappings if m.status == MappingStatus.MISSING)

    @property
    def failed_acs(self) -> int:
        """Number of ACs with failing tests."""
        return sum(1 for m in self.mappings if m.status == MappingStatus.FAILED)

    @property
    def coverage_pct(self) -> float:
        """Percentage of ACs with tests."""
        if self.total_acs == 0:
            return 100.0
        return (self.mapped_acs / self.total_acs) * 100

    @property
    def is_complete(self) -> bool:
        """Whether all ACs have tests."""
        return self.missing_acs == 0

    def to_dict(self) -> dict:
        """Serialize to dictionary."""
        return {
            "ws_id": self.ws_id,
            "total_acs": self.total_acs,
            "mapped_acs": self.mapped_acs,
            "missing_acs": self.missing_acs,
            "coverage_pct": self.coverage_pct,
            "mappings": [m.to_dict() for m in self.mappings],
        }

    def to_markdown_table(self) -> str:
        """Generate markdown table for report."""
        lines = [
            "| AC | Description | Test | Status |",
            "|----|-------------|------|--------|",
        ]

        for m in self.mappings:
            test = f"`{m.test_name}`" if m.test_name else "-"
            status = "✅" if m.status == MappingStatus.MAPPED else "❌"
            # Truncate description to 30 chars
            desc = m.ac_description[:30]
            if len(m.ac_description) > 30:
                desc += "..."
            lines.append(f"| {m.ac_id} | {desc} | {test} | {status} |")

        return "\n".join(lines)
