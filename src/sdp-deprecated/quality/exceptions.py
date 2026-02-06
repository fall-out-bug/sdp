"""Quality gate configuration exceptions."""



class ConfigValidationError(Exception):
    """Raised when quality-gate.toml validation fails."""

    def __init__(self, errors: list[str]) -> None:
        self.errors = errors
        super().__init__("\n".join(errors))
