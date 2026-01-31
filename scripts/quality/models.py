"""Data models for quality gate violations."""

from collections import namedtuple

Violation = namedtuple("Violation", ["category", "file_path", "line_no", "message", "severity"])
