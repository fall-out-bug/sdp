"""Tests for intent dataclass models."""

from sdp.schema.models import Intent, SuccessCriterion, Tradeoffs, TechnicalApproach


def test_intent_from_dict():
    data = {
        "problem": "Users need secure login",
        "users": ["end_users"],
        "success_criteria": [{"criterion": "Success rate", "measurement": ">95%"}],
        "tradeoffs": {"security": "prioritize", "performance": "accept"},
    }

    intent = Intent.from_dict(data)

    assert intent.problem == "Users need secure login"
    assert intent.users == ["end_users"]
    assert len(intent.success_criteria) == 1
    assert intent.tradeoffs.security == "prioritize"


def test_intent_to_dict():
    intent = Intent(
        problem="Test", users=["developers"], success_criteria=[SuccessCriterion("A", "B")]
    )

    data = intent.to_dict()

    assert data["problem"] == "Test"
    assert data["users"] == ["developers"]
    assert "success_criteria" in data
