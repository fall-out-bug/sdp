"""Auto-detect AC coverage from test files."""

import ast
import re
from dataclasses import dataclass
from pathlib import Path


@dataclass
class DetectedMapping:
    """A detected AC→Test mapping."""

    ac_id: str
    test_file: str
    test_name: str
    confidence: float
    source: str  # "docstring", "name", "keyword"


class ACDetector:
    """Auto-detect AC coverage from test files."""

    def detect_all(
        self, test_dir: Path, ac_descriptions: dict[str, str]
    ) -> list[DetectedMapping]:
        """Detect mappings from all test files.

        Args:
            test_dir: Directory containing test files
            ac_descriptions: Dict of AC IDs to descriptions

        Returns:
            List of detected mappings
        """
        mappings = []

        for test_file in test_dir.rglob("test_*.py"):
            file_mappings = self.detect_from_file(test_file, ac_descriptions)
            mappings.extend(file_mappings)

        return mappings

    def detect_from_file(
        self, test_file: Path, ac_descriptions: dict[str, str]
    ) -> list[DetectedMapping]:
        """Detect mappings from single test file.

        Args:
            test_file: Path to test file
            ac_descriptions: Dict of AC IDs to descriptions

        Returns:
            List of detected mappings from this file
        """
        try:
            tree = ast.parse(test_file.read_text())
        except SyntaxError:
            return []

        mappings = []

        for node in ast.walk(tree):
            if isinstance(node, ast.FunctionDef) and node.name.startswith("test_"):
                # Check docstring
                doc_mappings = self._detect_from_docstring(node, str(test_file))
                mappings.extend(doc_mappings)

                # Check name
                name_mappings = self._detect_from_name(node.name, str(test_file))
                mappings.extend(name_mappings)

                # Check keywords against AC descriptions
                keyword_mappings = self._detect_from_keywords(
                    node, str(test_file), ac_descriptions
                )
                mappings.extend(keyword_mappings)

        return mappings

    def _detect_from_docstring(
        self, func: ast.FunctionDef, test_file: str
    ) -> list[DetectedMapping]:
        """Extract AC references from docstring.

        Patterns:
            '''Tests AC1: User can login'''
            '''Covers: AC1, AC2'''
            '''AC1'''

        Args:
            func: AST function node
            test_file: Path to test file

        Returns:
            List of detected mappings
        """
        docstring = ast.get_docstring(func)
        if not docstring:
            return []

        # Pattern: AC1, AC2, etc.
        ac_refs = re.findall(r"\bAC(\d+)\b", docstring, re.IGNORECASE)

        return [
            DetectedMapping(
                ac_id=f"AC{ref}",
                test_file=test_file,
                test_name=func.name,
                confidence=0.95,  # High confidence from docstring
                source="docstring",
            )
            for ref in ac_refs
        ]

    def _detect_from_name(
        self, test_name: str, test_file: str
    ) -> list[DetectedMapping]:
        """Extract AC from test name.

        Patterns:
            test_ac1_user_login -> AC1
            test_acceptance_criterion_2 -> AC2
            test_ac_1_something -> AC1

        Args:
            test_name: Test function name
            test_file: Path to test file

        Returns:
            List of detected mappings
        """
        patterns = [
            (r"test_ac(\d+)", 0.90),
            (r"test_ac_(\d+)", 0.85),
            (r"test_acceptance_criterion_(\d+)", 0.90),
        ]

        mappings = []
        for pattern, confidence in patterns:
            match = re.search(pattern, test_name, re.IGNORECASE)
            if match:
                mappings.append(
                    DetectedMapping(
                        ac_id=f"AC{match.group(1)}",
                        test_file=test_file,
                        test_name=test_name,
                        confidence=confidence,
                        source="name",
                    )
                )

        return mappings

    def _detect_from_keywords(
        self,
        func: ast.FunctionDef,
        test_file: str,
        ac_descriptions: dict[str, str],
    ) -> list[DetectedMapping]:
        """Match test to AC by keyword similarity.

        Lower confidence - heuristic matching.

        Args:
            func: AST function node
            test_file: Path to test file
            ac_descriptions: Dict of AC IDs to descriptions

        Returns:
            List of detected mappings
        """
        # Get test name words
        test_words = set(re.findall(r"\w+", func.name.lower()))

        # Get docstring words
        docstring = ast.get_docstring(func) or ""
        test_words.update(re.findall(r"\w+", docstring.lower()))

        mappings = []

        for ac_id, description in ac_descriptions.items():
            desc_words = set(re.findall(r"\w+", description.lower()))

            # Remove common words
            common = {
                "the",
                "a",
                "an",
                "is",
                "are",
                "can",
                "should",
                "must",
                "test",
            }
            test_words_filtered = test_words - common
            desc_words_filtered = desc_words - common

            # Calculate overlap
            if not desc_words_filtered:
                continue

            overlap = len(test_words_filtered & desc_words_filtered) / len(
                desc_words_filtered
            )

            if overlap >= 0.5:  # At least 50% word match
                mappings.append(
                    DetectedMapping(
                        ac_id=ac_id,
                        test_file=test_file,
                        test_name=func.name,
                        confidence=min(0.7, overlap),  # Cap at 0.7 for keywords
                        source="keyword",
                    )
                )

        return mappings
