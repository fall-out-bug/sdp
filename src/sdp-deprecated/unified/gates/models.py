"""Approval gate data models."""

from dataclasses import dataclass
from datetime import datetime
from enum import Enum


class GateType(Enum):
    """Types of approval gates."""

    REQUIREMENTS = "requirements"
    ARCHITECTURE = "architecture"
    UAT = "uat"


class ApprovalStatus(Enum):
    """Status of an approval gate."""

    PENDING = "pending"
    APPROVED = "approved"
    REJECTED = "rejected"
    SKIPPED = "skipped"


@dataclass
class ApprovalGate:
    """Approval gate data model."""

    gate_type: GateType
    status: ApprovalStatus
    approved_by: str | None = None
    approved_at: datetime | None = None
    comments: str | None = None
