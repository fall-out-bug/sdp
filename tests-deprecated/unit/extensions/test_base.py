"""Tests for sdp.extensions.base."""

from pathlib import Path
from tempfile import TemporaryDirectory

from sdp.extensions.base import BaseExtension, ExtensionManifest


def test_extension_manifest_defaults() -> None:
    """Test ExtensionManifest uses default directory names."""
    manifest = ExtensionManifest(
        name="test-ext",
        version="1.0.0",
        description="Test extension",
        author="Test Author",
    )
    
    assert manifest.name == "test-ext"
    assert manifest.version == "1.0.0"
    assert manifest.description == "Test extension"
    assert manifest.author == "Test Author"
    assert manifest.hooks_dir == "hooks"
    assert manifest.patterns_dir == "patterns"
    assert manifest.skills_dir == "skills"
    assert manifest.integrations_dir == "integrations"


def test_extension_manifest_custom_dirs() -> None:
    """Test ExtensionManifest with custom directory names."""
    manifest = ExtensionManifest(
        name="test-ext",
        version="1.0.0",
        description="Test extension",
        author="Test Author",
        hooks_dir="custom_hooks",
        patterns_dir="custom_patterns",
        skills_dir="custom_skills",
        integrations_dir="custom_integrations",
    )
    
    assert manifest.hooks_dir == "custom_hooks"
    assert manifest.patterns_dir == "custom_patterns"
    assert manifest.skills_dir == "custom_skills"
    assert manifest.integrations_dir == "custom_integrations"


def test_base_extension_get_hooks_path_exists() -> None:
    """Test BaseExtension.get_hooks_path returns path when directory exists."""
    with TemporaryDirectory() as tmpdir:
        root = Path(tmpdir)
        hooks_dir = root / "hooks"
        hooks_dir.mkdir()
        
        manifest = ExtensionManifest(
            name="test-ext",
            version="1.0.0",
            description="Test",
            author="Test",
        )
        ext = BaseExtension(manifest=manifest, root_path=root)
        
        result = ext.get_hooks_path()
        assert result == hooks_dir
        assert result.is_dir()


def test_base_extension_get_hooks_path_not_exists() -> None:
    """Test BaseExtension.get_hooks_path returns None when directory missing."""
    with TemporaryDirectory() as tmpdir:
        root = Path(tmpdir)
        
        manifest = ExtensionManifest(
            name="test-ext",
            version="1.0.0",
            description="Test",
            author="Test",
        )
        ext = BaseExtension(manifest=manifest, root_path=root)
        
        result = ext.get_hooks_path()
        assert result is None


def test_base_extension_get_patterns_path_exists() -> None:
    """Test BaseExtension.get_patterns_path returns path when directory exists."""
    with TemporaryDirectory() as tmpdir:
        root = Path(tmpdir)
        patterns_dir = root / "patterns"
        patterns_dir.mkdir()
        
        manifest = ExtensionManifest(
            name="test-ext",
            version="1.0.0",
            description="Test",
            author="Test",
        )
        ext = BaseExtension(manifest=manifest, root_path=root)
        
        result = ext.get_patterns_path()
        assert result == patterns_dir


def test_base_extension_get_patterns_path_not_exists() -> None:
    """Test BaseExtension.get_patterns_path returns None when directory missing."""
    with TemporaryDirectory() as tmpdir:
        root = Path(tmpdir)
        
        manifest = ExtensionManifest(
            name="test-ext",
            version="1.0.0",
            description="Test",
            author="Test",
        )
        ext = BaseExtension(manifest=manifest, root_path=root)
        
        result = ext.get_patterns_path()
        assert result is None


def test_base_extension_get_skills_path_exists() -> None:
    """Test BaseExtension.get_skills_path returns path when directory exists."""
    with TemporaryDirectory() as tmpdir:
        root = Path(tmpdir)
        skills_dir = root / "skills"
        skills_dir.mkdir()
        
        manifest = ExtensionManifest(
            name="test-ext",
            version="1.0.0",
            description="Test",
            author="Test",
        )
        ext = BaseExtension(manifest=manifest, root_path=root)
        
        result = ext.get_skills_path()
        assert result == skills_dir


def test_base_extension_get_skills_path_not_exists() -> None:
    """Test BaseExtension.get_skills_path returns None when directory missing."""
    with TemporaryDirectory() as tmpdir:
        root = Path(tmpdir)
        
        manifest = ExtensionManifest(
            name="test-ext",
            version="1.0.0",
            description="Test",
            author="Test",
        )
        ext = BaseExtension(manifest=manifest, root_path=root)
        
        result = ext.get_skills_path()
        assert result is None


def test_base_extension_get_integrations_path_exists() -> None:
    """Test BaseExtension.get_integrations_path returns path when directory exists."""
    with TemporaryDirectory() as tmpdir:
        root = Path(tmpdir)
        integrations_dir = root / "integrations"
        integrations_dir.mkdir()
        
        manifest = ExtensionManifest(
            name="test-ext",
            version="1.0.0",
            description="Test",
            author="Test",
        )
        ext = BaseExtension(manifest=manifest, root_path=root)
        
        result = ext.get_integrations_path()
        assert result == integrations_dir


def test_base_extension_get_integrations_path_not_exists() -> None:
    """Test BaseExtension.get_integrations_path returns None when directory missing."""
    with TemporaryDirectory() as tmpdir:
        root = Path(tmpdir)
        
        manifest = ExtensionManifest(
            name="test-ext",
            version="1.0.0",
            description="Test",
            author="Test",
        )
        ext = BaseExtension(manifest=manifest, root_path=root)
        
        result = ext.get_integrations_path()
        assert result is None
