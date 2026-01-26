"""Dataclass models for intent specification."""

from dataclasses import dataclass, field
from typing import Literal, cast


@dataclass
class SuccessCriterion:
    """A single success criterion with measurement."""
    criterion: str
    measurement: str


@dataclass
class TechnicalApproach:
    """Technical approach preferences."""
    architecture: str | None = None
    storage: str | None = None
    failure_mode: str | None = None
    auth_method: str | None = None


@dataclass
class Tradeoffs:
    """Strategic tradeoff preferences."""
    security: Literal["prioritize", "accept", "reject"] | None = None
    performance: Literal["prioritize", "accept", "reject"] | None = None
    complexity: Literal["prioritize", "accept", "reject"] | None = None
    time_to_market: Literal["prioritize", "accept", "reject"] | None = None


@dataclass
class Intent:
    """Structured intent specification."""

    problem: str
    users: list[Literal["end_users", "admins", "developers", "api_consumers", "operators"]]
    success_criteria: list[SuccessCriterion]
    tradeoffs: Tradeoffs | None = None
    technical_approach: TechnicalApproach | None = None

    @classmethod
    def from_dict(cls, data: dict[str, object]) -> "Intent":
        """Create from dict, converting nested structures."""
        success_criteria_data = cast(list[object], data["success_criteria"])
        success_criteria: list[SuccessCriterion] = []
        for sc in success_criteria_data:
            sc_dict = cast(dict[str, object], sc)
            success_criteria.append(
                SuccessCriterion(
                    criterion=cast(str, sc_dict["criterion"]),
                    measurement=cast(str, sc_dict["measurement"]),
                )
            )

        tradeoffs = None
        if data.get("tradeoffs"):
            tradeoffs_dict = cast(dict[str, object], data["tradeoffs"])
            tradeoffs = Tradeoffs(
                security=cast(Literal["prioritize", "accept", "reject"] | None, tradeoffs_dict.get("security")),
                performance=cast(Literal["prioritize", "accept", "reject"] | None, tradeoffs_dict.get("performance")),
                complexity=cast(Literal["prioritize", "accept", "reject"] | None, tradeoffs_dict.get("complexity")),
                time_to_market=cast(Literal["prioritize", "accept", "reject"] | None, tradeoffs_dict.get("time_to_market")),
            )

        technical_approach = None
        if data.get("technical_approach"):
            ta_dict = cast(dict[str, object], data["technical_approach"])
            technical_approach = TechnicalApproach(
                architecture=cast(str | None, ta_dict.get("architecture")),
                storage=cast(str | None, ta_dict.get("storage")),
                failure_mode=cast(str | None, ta_dict.get("failure_mode")),
                auth_method=cast(str | None, ta_dict.get("auth_method")),
            )

        return cls(
            problem=cast(str, data["problem"]),
            users=cast(list[Literal["end_users", "admins", "developers", "api_consumers", "operators"]], data["users"]),
            success_criteria=success_criteria,
            tradeoffs=tradeoffs,
            technical_approach=technical_approach,
        )

    def to_dict(self) -> dict[str, object]:
        """Convert to dict for JSON serialization."""
        return {
            "problem": self.problem,
            "users": self.users,
            "success_criteria": [
                {"criterion": sc.criterion, "measurement": sc.measurement}
                for sc in self.success_criteria
            ],
            "tradeoffs": self.tradeoffs.__dict__ if self.tradeoffs else None,
            "technical_approach": (
                self.technical_approach.__dict__ if self.technical_approach else None
            ),
        }
