"""Targeted tests for uncovered lines in core/model/parser.py."""

from sdp.core.model.parser import parse_context_window, parse_models_table


def test_parse_models_table_no_table() -> None:
    """Test parse_models_table with no table (line 25)."""
    section = """
# Tier 1

Some content without a table.
"""
    result = parse_models_table(section)
    assert result == []


def test_parse_models_table_old_format() -> None:
    """Test parse_models_table with old format (lines 43-45)."""
    section = """
| Provider | Model | Context | Tool Use | Notes |
|----------|-------|---------|----------|-------|
| OpenAI | gpt-4 | 128K | ✅ Full | Great |
"""
    result = parse_models_table(section)
    assert len(result) == 1
    assert result[0].provider == "OpenAI"
    assert result[0].model == "gpt-4"
    assert result[0].cost_per_1m_tokens == 0.0  # Default
    assert result[0].availability_pct == 0.99  # Default


def test_parse_models_table_skip_empty_rows() -> None:
    """Test parse_models_table skips empty rows (lines 52-53)."""
    section = """
| Provider | Model | Context | Tool Use | Notes |
|----------|-------|---------|----------|-------|
|  |  | 128K | ✅ | Empty |
| OpenAI | gpt-4 | 128K | ✅ | Valid |
"""
    result = parse_models_table(section)
    assert len(result) == 1
    assert result[0].provider == "OpenAI"


def test_parse_models_table_cost_parse_error() -> None:
    """Test parse_models_table with invalid cost (lines 66-67)."""
    section = """
| Provider | Model | Cost ($/1M) | Availability | Context | Tool Use | Notes |
|----------|-------|-------------|--------------|---------|----------|-------|
| OpenAI | gpt-4 | invalid | 99.9% | 128K | ✅ | Great |
"""
    result = parse_models_table(section)
    assert len(result) == 1
    assert result[0].cost_per_1m_tokens == 0.0  # Default on parse error


def test_parse_models_table_availability_parse_error() -> None:
    """Test parse_models_table with invalid availability (lines 72-73)."""
    section = """
| Provider | Model | Cost ($/1M) | Availability | Context | Tool Use | Notes |
|----------|-------|-------------|--------------|---------|----------|-------|
| OpenAI | gpt-4 | 3.00 | invalid | 128K | ✅ | Great |
"""
    result = parse_models_table(section)
    assert len(result) == 1
    assert result[0].availability_pct == 0.99  # Default on parse error


def test_parse_context_window_m_plus() -> None:
    """Test parse_context_window with M+ suffix (lines 119-120)."""
    assert parse_context_window("1M+") == 1_000_000
    assert parse_context_window("2M+") == 2_000_000
    assert parse_context_window("1.5M+") == 1_500_000


def test_parse_context_window_m() -> None:
    """Test parse_context_window with M suffix (lines 121-122)."""
    assert parse_context_window("1M") == 1_000_000
    assert parse_context_window("2M") == 2_000_000


def test_parse_context_window_k() -> None:
    """Test parse_context_window with K suffix (lines 125-126)."""
    assert parse_context_window("200K") == 200_000
    assert parse_context_window("128K") == 128_000


def test_parse_context_window_plain_number() -> None:
    """Test parse_context_window with plain number (lines 129-130)."""
    assert parse_context_window("1000") == 1000
    assert parse_context_window("5000") == 5000


def test_parse_context_window_invalid() -> None:
    """Test parse_context_window with invalid input (lines 131-132)."""
    assert parse_context_window("invalid") == 128_000  # Default fallback
    assert parse_context_window("") == 128_000
    assert parse_context_window("abc") == 128_000
