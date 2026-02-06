"""Unit tests for AC detector."""

import tempfile
from pathlib import Path

import pytest

from sdp.traceability.detector import ACDetector, DetectedMapping


@pytest.fixture
def temp_test_dir():
    """Create temporary test directory."""
    with tempfile.TemporaryDirectory() as tmpdir:
        yield Path(tmpdir)


class TestACDetector:
    """Tests for ACDetector."""

    def test_detect_from_docstring_explicit(self, temp_test_dir):
        """AC1: Detect AC from docstring with explicit reference."""
        test_file = temp_test_dir / "test_auth.py"
        test_file.write_text(
            '''
def test_user_login():
    """Tests AC1: User can login."""
    pass
'''
        )

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, {})

        assert len(mappings) == 1
        assert mappings[0].ac_id == "AC1"
        assert mappings[0].test_name == "test_user_login"
        assert mappings[0].source == "docstring"
        assert mappings[0].confidence == 0.95

    def test_detect_from_docstring_multiple(self, temp_test_dir):
        """AC1: Detect multiple ACs from docstring."""
        test_file = temp_test_dir / "test_auth.py"
        test_file.write_text(
            '''
def test_auth_flow():
    """Tests AC1 and AC2 together."""
    pass
'''
        )

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, {})

        assert len(mappings) == 2
        ac_ids = {m.ac_id for m in mappings}
        assert "AC1" in ac_ids
        assert "AC2" in ac_ids

    def test_detect_from_name_ac_prefix(self, temp_test_dir):
        """AC2: Detect AC from test name with ac prefix."""
        test_file = temp_test_dir / "test_auth.py"
        test_file.write_text(
            '''
def test_ac1_user_login():
    pass

def test_ac_2_user_logout():
    pass
'''
        )

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, {})

        assert len(mappings) == 2
        ac1 = next(m for m in mappings if m.ac_id == "AC1")
        ac2 = next(m for m in mappings if m.ac_id == "AC2")

        assert ac1.source == "name"
        assert ac1.confidence == 0.90
        assert ac2.source == "name"
        assert ac2.confidence == 0.85  # Lower for ac_ pattern

    def test_detect_from_name_acceptance_criterion(self, temp_test_dir):
        """AC2: Detect AC from test name with acceptance_criterion."""
        test_file = temp_test_dir / "test_auth.py"
        test_file.write_text(
            '''
def test_acceptance_criterion_1():
    pass
'''
        )

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, {})

        assert len(mappings) == 1
        assert mappings[0].ac_id == "AC1"
        assert mappings[0].confidence == 0.90

    def test_detect_from_keywords(self, temp_test_dir):
        """AC3: Detect AC from keyword matching."""
        test_file = temp_test_dir / "test_auth.py"
        test_file.write_text(
            '''
def test_user_login_success():
    """Test successful user login."""
    pass
'''
        )

        ac_descriptions = {"AC1": "User can login successfully"}

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, ac_descriptions)

        # Should detect from keywords
        keyword_mappings = [m for m in mappings if m.source == "keyword"]
        assert len(keyword_mappings) > 0

        ac1 = keyword_mappings[0]
        assert ac1.ac_id == "AC1"
        assert ac1.confidence <= 0.7  # Capped for keywords

    def test_detect_all_multiple_files(self, temp_test_dir):
        """AC5: Detect from multiple test files."""
        # Create multiple test files
        (temp_test_dir / "test_auth.py").write_text(
            '''
def test_ac1_login():
    pass
'''
        )

        (temp_test_dir / "test_session.py").write_text(
            '''
def test_ac2_session():
    pass
'''
        )

        detector = ACDetector()
        mappings = detector.detect_all(temp_test_dir, {})

        assert len(mappings) == 2
        ac_ids = {m.ac_id for m in mappings}
        assert "AC1" in ac_ids
        assert "AC2" in ac_ids

    def test_detect_all_nested_directories(self, temp_test_dir):
        """AC5: Detect from nested test directories."""
        # Create nested structure
        nested = temp_test_dir / "unit"
        nested.mkdir()

        (nested / "test_auth.py").write_text(
            '''
def test_ac1_login():
    pass
'''
        )

        detector = ACDetector()
        mappings = detector.detect_all(temp_test_dir, {})

        assert len(mappings) == 1
        assert mappings[0].ac_id == "AC1"

    def test_detect_ignores_syntax_errors(self, temp_test_dir):
        """Detector should ignore files with syntax errors."""
        test_file = temp_test_dir / "test_broken.py"
        test_file.write_text("def broken syntax(")

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, {})

        assert len(mappings) == 0

    def test_detect_ignores_non_test_functions(self, temp_test_dir):
        """Detector should only process test functions."""
        test_file = temp_test_dir / "test_auth.py"
        test_file.write_text(
            '''
def helper_function():
    """AC1: Helper"""
    pass

def test_ac2_actual_test():
    pass
'''
        )

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, {})

        # Should only detect AC2 from test function
        ac_ids = {m.ac_id for m in mappings}
        assert "AC2" in ac_ids
        # Should NOT detect AC1 from helper
        assert "AC1" not in ac_ids

    def test_detect_case_insensitive(self, temp_test_dir):
        """AC1, AC2: Detection should be case-insensitive."""
        test_file = temp_test_dir / "test_auth.py"
        test_file.write_text(
            '''
def test_ac1_login():
    """Tests ac2 lowercase."""
    pass
'''
        )

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, {})

        ac_ids = {m.ac_id for m in mappings}
        assert "AC1" in ac_ids
        assert "AC2" in ac_ids

    def test_confidence_scores(self, temp_test_dir):
        """AC4: Confidence scores are assigned correctly."""
        test_file = temp_test_dir / "test_auth.py"
        test_file.write_text(
            '''
def test_ac1_docstring():
    """Tests AC10 explicitly."""
    pass

def test_ac2_name():
    pass

def test_keyword_match():
    """Test user login."""
    pass
'''
        )

        ac_descriptions = {"AC20": "User can login"}

        detector = ACDetector()
        mappings = detector.detect_from_file(test_file, ac_descriptions)

        # Check confidence ranges
        docstring_mappings = [m for m in mappings if m.source == "docstring"]
        name_mappings = [m for m in mappings if m.source == "name"]
        keyword_mappings = [m for m in mappings if m.source == "keyword"]

        assert all(m.confidence == 0.95 for m in docstring_mappings)
        assert all(0.85 <= m.confidence <= 0.90 for m in name_mappings)
        assert all(m.confidence <= 0.7 for m in keyword_mappings)
