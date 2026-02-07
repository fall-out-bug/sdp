"""Feature loading utilities."""

from pathlib import Path

from sdp.core.feature.models import Feature
from sdp.core.workstream import Workstream, WorkstreamParseError, parse_workstream


def load_feature_from_directory(
    feature_id: str, directory: Path, pattern: str = "WS-*.md"
) -> Feature:
    """Load feature from directory containing workstream files.

    Args:
        feature_id: Feature identifier (FXXX)
        directory: Directory containing workstream markdown files
        pattern: Glob pattern for workstream files (default: "WS-*.md")

    Returns:
        Feature instance with loaded workstreams

    Raises:
        WorkstreamParseError: If any workstream file fails to parse
        ValueError: If no workstreams found or feature_id mismatch
    """
    ws_files = sorted(directory.glob(pattern))
    if not ws_files:
        raise ValueError(f"No workstream files found in {directory} matching {pattern}")

    workstreams: list[Workstream] = []
    for ws_file in ws_files:
        try:
            ws = parse_workstream(ws_file)
            if ws.feature != feature_id:
                raise ValueError(
                    f"Workstream {ws.ws_id} has feature {ws.feature}, expected {feature_id}"
                )
            workstreams.append(ws)
        except WorkstreamParseError as e:
            raise WorkstreamParseError(
                message=f"Failed to parse {ws_file}: {e}",
                file_path=ws_file,
                parse_error=str(e),
            ) from e

    return Feature(feature_id=feature_id, workstreams=workstreams)


def load_feature_from_spec(feature_id: str, spec_file: Path) -> Feature:
    """Load feature from spec file (future: parse feature.md and find WS).

    For now, this is a placeholder that loads from the workstreams directory
    based on feature_id.

    Args:
        feature_id: Feature identifier (FXXX)
        spec_file: Path to feature spec file (e.g., feature.md)

    Returns:
        Feature instance

    Raises:
        NotImplementedError: This is a placeholder for future implementation
    """
    # Future: Parse feature.md to find workstream directory
    # For now, infer from spec_file location
    spec_dir = spec_file.parent
    workstreams_dir = spec_dir.parent.parent / "workstreams" / "backlog"

    if not workstreams_dir.exists():
        raise ValueError(f"Workstreams directory not found: {workstreams_dir}")

    return load_feature_from_directory(feature_id, workstreams_dir)
