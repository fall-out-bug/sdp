"""Dataclass models for quality gate configuration."""

from dataclasses import dataclass


@dataclass
class CoverageConfig:
    """Test coverage configuration."""

    enabled: bool = True
    minimum: int = 80
    fail_under: int = 80
    exclude_patterns: list[str] | None = None

    def __post_init__(self) -> None:
        if self.exclude_patterns is None:
            self.exclude_patterns = []


@dataclass
class ComplexityConfig:
    """Cyclomatic complexity configuration."""

    enabled: bool = True
    max_cc: int = 10
    max_average_cc: int = 5


@dataclass
class FileSizeConfig:
    """File size limits configuration."""

    enabled: bool = True
    max_lines: int = 200
    max_imports: int = 20
    max_functions: int = 15


@dataclass
class TypeHintsConfig:
    """Type hinting requirements configuration."""

    enabled: bool = True
    require_return_types: bool = True
    require_param_types: bool = True
    strict_mode: bool = True
    allow_implicit_any: bool = False


@dataclass
class ErrorHandlingConfig:
    """Error handling patterns configuration."""

    enabled: bool = True
    forbid_bare_except: bool = True
    forbid_bare_raise: bool = True
    forbid_pass_with_except: bool = True
    require_explicit_errors: bool = True


@dataclass
class LayerPattern:
    """Pattern for detecting architectural layers."""

    name: str
    path_regex: str  # Regex pattern to match file paths
    module_regex: str | None = None  # Optional regex for import statements


@dataclass
class ArchitectureConfig:
    """Clean architecture boundaries configuration."""

    enabled: bool = True
    enforce_layer_separation: bool = True
    allowed_layer_imports: list[str] | None = None
    forbid_violations: list[str] | None = None
    layer_patterns: list[LayerPattern] | None = None  # Configurable layer detection

    def __post_init__(self) -> None:
        if self.allowed_layer_imports is None:
            self.allowed_layer_imports = []
        if self.forbid_violations is None:
            # Default clean architecture rules
            # Domain is innermost - cannot depend on anyone
            # Application can use domain, but not infrastructure/presentation
            # Infrastructure can use domain and application, but not presentation
            self.forbid_violations = [
                "domain -> application",
                "domain -> infrastructure",
                "domain -> presentation",
                "application -> infrastructure",
                "application -> presentation",
                "infrastructure -> presentation",
            ]
        # Note: layer_patterns defaults handled in ArchitectureChecker


@dataclass
class DocumentationConfig:
    """Documentation requirements configuration."""

    enabled: bool = True
    require_docstrings: bool = False
    min_docstring_coverage: float = 0.5
    require_module_docstrings: bool = True
    require_class_docstrings: bool = False
    require_function_docstrings: bool = False


@dataclass
class TestingConfig:
    """Test quality requirements configuration."""

    enabled: bool = True
    require_test_for_new_code: bool = True
    min_test_to_code_ratio: float = 0.8
    require_fast_marker: bool = True
    forbid_print_statements: bool = True


@dataclass
class NamingConfig:
    """Naming conventions configuration."""

    enabled: bool = True
    enforce_pep8: bool = True
    allow_single_letter: bool = False
    min_variable_name_length: int = 3
    max_variable_name_length: int = 50


@dataclass
class SecurityConfig:
    """Security checks configuration."""

    enabled: bool = True
    forbid_hardcoded_secrets: bool = True
    forbid_sql_injection_patterns: bool = True
    forbid_eval_usage: bool = True
    require_https_urls: bool = True


@dataclass
class PerformanceConfig:
    """Performance checks configuration."""

    enabled: bool = True
    forbid_sql_queries_in_loops: bool = True
    max_nesting_depth: int = 5
    warn_large_string_concatenation: bool = True


@dataclass
class QualityGateConfig:
    """Complete quality gate configuration."""

    coverage: CoverageConfig
    complexity: ComplexityConfig
    file_size: FileSizeConfig
    type_hints: TypeHintsConfig
    error_handling: ErrorHandlingConfig
    architecture: ArchitectureConfig
    documentation: DocumentationConfig | None = None
    testing: TestingConfig | None = None
    naming: NamingConfig | None = None
    security: SecurityConfig | None = None
    performance: PerformanceConfig | None = None
