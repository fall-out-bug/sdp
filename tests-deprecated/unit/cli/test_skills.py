"""Tests for skills discovery CLI."""
from click.testing import CliRunner

# Import to ensure coverage
from sdp.cli.skill import skill

runner = CliRunner()


def test_skills_list():
    """Test skill listing."""
    result = runner.invoke(skill, ["list"])
    assert result.exit_code == 0
    assert "@build" in result.output
    assert "@feature" in result.output


def test_skills_show():
    """Test skill details."""
    result = runner.invoke(skill, ["show", "build"])
    assert result.exit_code == 0
    assert "Execute single workstream" in result.output
    assert "Example:" in result.output


def test_skills_show_with_prefix():
    """Test skill lookup with @ prefix."""
    result = runner.invoke(skill, ["show", "@build"])
    assert result.exit_code == 0
    assert "@build" in result.output


def test_skills_show_not_found():
    """Test unknown skill."""
    result = runner.invoke(skill, ["show", "unknown"])
    assert result.exit_code != 0
    assert "not found" in result.output


def test_skills_filter_by_category():
    """Test category filter."""
    result = runner.invoke(skill, ["list", "--category", "workflow"])
    assert result.exit_code == 0
    assert "@build" in result.output
    assert "@hotfix" not in result.output  # fixes category


def test_all_12_skills_in_registry():
    """Test that all 12 skills are present."""
    from sdp.cli.skills.registry import SKILLS_REGISTRY

    expected_skills = {
        "feature", "idea", "design", "build", "oneshot",  # workflow
        "debug", "issue",  # debugging
        "hotfix", "bugfix",  # fixes
        "review", "deploy",  # deployment
        "help",  # utility
    }

    assert set(SKILLS_REGISTRY.keys()) == expected_skills
    assert len(SKILLS_REGISTRY) == 12


def test_skill_categories():
    """Test that skills are properly categorized."""
    from sdp.cli.skills.registry import SKILLS_REGISTRY, SkillCategory

    workflow_skills = [s.name for s in SKILLS_REGISTRY.values()
                       if s.category == SkillCategory.WORKFLOW]
    assert len(workflow_skills) == 5
    assert "@feature" in workflow_skills

    debugging_skills = [s.name for s in SKILLS_REGISTRY.values()
                        if s.category == SkillCategory.DEBUGGING]
    assert len(debugging_skills) == 2

    fixes_skills = [s.name for s in SKILLS_REGISTRY.values()
                    if s.category == SkillCategory.FIXES]
    assert len(fixes_skills) == 2

    deployment_skills = [s.name for s in SKILLS_REGISTRY.values()
                         if s.category == SkillCategory.DEPLOYMENT]
    assert len(deployment_skills) == 2

    utility_skills = [s.name for s in SKILLS_REGISTRY.values()
                      if s.category == SkillCategory.UTILITY]
    assert len(utility_skills) == 1

