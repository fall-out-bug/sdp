"""Security checks for quality gates.

Checks for:
- Hardcoded secrets (passwords, API keys, tokens)
- Use of eval() function
"""

import ast
import re
from pathlib import Path

from scripts.quality.models import Violation


class SecurityChecker:
    """Security issue checker."""

    def __init__(
        self,
        forbid_hardcoded_secrets: bool = True,
        forbid_eval_usage: bool = True,
    ) -> None:
        """Initialize security checker.

        Args:
            forbid_hardcoded_secrets: Check for hardcoded secrets
            forbid_eval_usage: Check for eval() usage
        """
        self.forbid_hardcoded_secrets = forbid_hardcoded_secrets
        self.forbid_eval_usage = forbid_eval_usage
        self._violations: list[Violation] = []

    def check(self, path: Path, source_code: str, tree: ast.AST) -> list[Violation]:
        """Run all security checks.

        Args:
            path: File path
            source_code: Source code string
            tree: AST tree

        Returns:
            List of violations
        """
        self._violations.clear()

        if self.forbid_hardcoded_secrets:
            self._check_hardcoded_secrets(path, source_code)

        if self.forbid_eval_usage:
            self._check_eval_usage(path, tree)

        return self._violations

    def _check_hardcoded_secrets(self, path: Path, source_code: str) -> None:
        """Check for hardcoded secrets.

        Detects:
        - password, passwd, pwd
        - api_key, apikey, api-key
        - secret, secret_key, secret-key
        - token, auth_token, auth-token
        - private_key, private-key, privatekey

        Excludes obvious test/example values.
        """
        secret_patterns = [
            r'(?:password|passwd|pwd)\s*=\s*["\']([^"\']{8,})["\']',
            r'(?:api_key|apikey|api-key)\s*=\s*["\']([^"\']{8,})["\']',
            r'(?:secret|secret_key|secret-key)\s*=\s*["\']([^"\']{8,})["\']',
            r'(?:token|auth_token|auth-token)\s*=\s*["\']([^"\']{8,})["\']',
            r'(?:private_key|private-key|privatekey)\s*=\s*["\']([^"\']{8,})["\']',
        ]

        for pattern in secret_patterns:
            matches = re.finditer(pattern, source_code, re.IGNORECASE)
            for match in matches:
                # Exclude obvious test/example values
                value = match.group(1)
                if not re.search(
                    r'^(test|example|mock|dummy|xxx|xxx+|\*+)$', value, re.IGNORECASE
                ):
                    line_num = source_code[: match.start()].count("\n") + 1
                    self._violations.append(
                        Violation(
                            "security",
                            str(path),
                            line_num,
                            f"Possible hardcoded secret: {match.group(1)[:10]}...",
                            "error",
                        )
                    )

    def _check_eval_usage(self, path: Path, tree: ast.AST) -> None:
        """Check for eval() usage using AST parsing.

        Uses AST instead of string matching to avoid false positives.
        """
        for node in ast.walk(tree):
            if isinstance(node, ast.Call):
                # Check if calling eval()
                if isinstance(node.func, ast.Name) and node.func.id == "eval":
                    self._violations.append(
                        Violation(
                            "security",
                            str(path),
                            node.lineno,
                            "Use of eval() detected (security risk)",
                            "error",
                        )
                    )
