import pytest
from pathlib import Path
from sdp.schema.validator import IntentValidator, ValidationError


def test_validate_valid_intent():
    validator = IntentValidator()
    intent = {
        "problem": "Users need secure account access with social login",
        "users": ["end_users", "admins"],
        "success_criteria": [
            {"criterion": "Login success rate", "measurement": "Greater than 95%"},
            {"criterion": "Login latency", "measurement": "Under 500ms p95"}
        ]
    }

    # Should not raise
    validator.validate(intent)


def test_validate_missing_required_field():
    validator = IntentValidator()
    intent = {
        "users": ["end_users"]
        # Missing: problem, success_criteria
    }

    with pytest.raises(ValidationError) as exc_info:
        validator.validate(intent)

    # Error message contains at least one missing required field
    assert "problem" in str(exc_info.value) or "success_criteria" in str(exc_info.value)


def test_validate_invalid_user_enum():
    validator = IntentValidator()
    intent = {
        "problem": "Test problem" * 10,  # >50 chars
        "users": ["invalid_user_type"],
        "success_criteria": [{"criterion": "Test criterion", "measurement": "12345"}]
    }

    with pytest.raises(ValidationError) as exc_info:
        validator.validate(intent)

    assert "invalid_user_type" in str(exc_info.value)


def test_validate_min_length_problem():
    validator = IntentValidator()
    intent = {
        "problem": "Too short",  # <50 chars
        "users": ["end_users"],
        "success_criteria": [{"criterion": "Test criterion", "measurement": "12345"}]
    }

    with pytest.raises(ValidationError) as exc_info:
        validator.validate(intent)

    # Error mentions it's too short
    assert "too short" in str(exc_info.value)
