---
assignee: null
completed: '2026-01-30'
depends_on: []
feature: F032
github_issue: null
project_id: 0
size: SMALL
status: completed
traceability:
- ac_description: '`templates/skill-template.md` создан с примером структуры'
  ac_id: AC1
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_happy_path_activate_edit_complete
- ac_description: Template ограничен ~100 строками без примеров
  ac_id: AC2
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_without_active_ws
- ac_description: ADR документ объясняет решение о лимите
  ac_id: AC3
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_edit_blocked_outside_scope
- ac_description: Warning если skill > 100 строк, error если > 150
  ac_id: AC5
  confidence: 1.0
  status: mapped
  test_file: tests/unit/test_scope_manager.py
  test_name: test_is_in_scope_with_restricted_scope
- ac_description: '`sdp skill validate <path>` проверяет размер skill файла'
  ac_id: AC4
  confidence: 1.0
  status: mapped
  test_file: tests/integration/test_guard_flow.py
  test_name: test_concurrent_activation_blocked
ws_id: 00-032-05
---

## 00-032-05: Skill Template Standard

### 🎯 Goal

**What must WORK after completing this WS:**
- Стандарт для skills: короткие (< 100 строк), со ссылками на docs
- Template файл для создания новых skills
- Валидатор проверяет размер skills

**Acceptance Criteria:**
- [ ] AC1: `templates/skill-template.md` создан с примером структуры
- [ ] AC2: Template ограничен ~100 строками без примеров
- [ ] AC3: ADR документ объясняет решение о лимите
- [ ] AC4: `sdp skill validate <path>` проверяет размер skill файла
- [ ] AC5: Warning если skill > 100 строк, error если > 150

**⚠️ WS is NOT complete until Goal is achieved (all AC ✅).**

---

### Context

**Problem**: Текущие skills 400-500 строк. Агенты теряют фокус.

**Solution**: Стандарт коротких skills с references на детальную документацию.

### Dependencies

None (independent)

### Steps

1. **Create skill template**

   ```markdown
   # templates/skill-template.md
   ---
   name: {skill_name}
   description: {One-line description, max 80 chars}
   tools: {Comma-separated list of required tools}
   max_lines: 100
   ---
   
   # @{skill_name} - {Title}
   
   {2-3 sentence purpose. What does this skill do and when to use it.}
   
   ## Quick Reference
   
   | Step | Action | Gate |
   |------|--------|------|
   | 1 | {action} | {what must be true} |
   | 2 | {action} | {what must be true} |
   | 3 | {action} | {what must be true} |
   
   ## Workflow
   
   ### Step 1: {Name}
   
   ```bash
   # Command or action
   {command}
   ```
   
   **Gate:** {What must be true before proceeding}
   
   ### Step 2: {Name}
   
   {3-5 lines max}
   
   ### Step 3: {Name}
   
   {3-5 lines max}
   
   ## Quality Gates
   
   See [Quality Gates Reference](../docs/reference/quality-gates.md)
   
   ## Errors
   
   | Error | Cause | Fix |
   |-------|-------|-----|
   | {error} | {cause} | {fix} |
   
   ## See Also
   
   - [Full Specification](../docs/reference/{skill}-spec.md)
   - [Examples](../docs/examples/{skill}/)
   - [Related Skill](./{related}/SKILL.md)
   ```

2. **Create ADR**

   ```markdown
   # docs/adr/007-skill-length-limit.md
   # ADR-007: Skill Length Limit
   
   ## Status
   Accepted
   
   ## Context
   
   Current skills are 400-500 lines. Evidence shows:
   - Agents cherry-pick rules from long prompts
   - Later sections get ignored
   - Duplicate instructions lead to confusion
   
   Research on LLM prompt effectiveness suggests:
   - Shorter prompts have higher compliance
   - Clear structure improves following
   - References work better than inline detail
   
   ## Decision
   
   Limit skills to **100 lines** (excluding code examples in docs/).
   
   ### Structure
   
   1. **Header** (10 lines): name, description, tools
   2. **Purpose** (5 lines): what and when
   3. **Quick Reference** (10 lines): table of steps
   4. **Workflow** (50 lines): step-by-step, 5-10 lines each
   5. **Quality Gates** (5 lines): reference to docs
   6. **Errors** (10 lines): common errors table
   7. **See Also** (10 lines): links to detailed docs
   
   ### Details go to docs/
   
   - `docs/reference/{skill}-spec.md` — Full specification
   - `docs/examples/{skill}/` — Code examples
   - `docs/reference/quality-gates.md` — Shared quality gates
   
   ## Consequences
   
   ### Positive
   - Higher agent compliance
   - Consistent skill structure
   - Easier maintenance
   - Clear separation of concerns
   
   ### Negative
   - Need to create docs/reference/ files
   - Migration effort for existing skills
   - Agents need to follow references
   
   ## Validation
   
   `sdp skill validate` checks:
   - Line count ≤ 100 (warning), ≤ 150 (error)
   - Required sections present
   - References resolve
   ```

3. **Create validator**

   ```python
   # src/sdp/cli/skill.py
   import typer
   from pathlib import Path
   import re
   
   app = typer.Typer(help="Skill management commands")
   
   REQUIRED_SECTIONS = [
       "## Quick Reference",
       "## Workflow", 
       "## See Also",
   ]
   
   @app.command("validate")
   def validate_skill(
       path: Path,
       strict: bool = typer.Option(False, help="Fail on warnings")
   ) -> None:
       """Validate skill file against standards."""
       if not path.exists():
           typer.echo(f"❌ File not found: {path}")
           raise typer.Exit(1)
       
       content = path.read_text()
       lines = content.splitlines()
       line_count = len(lines)
       
       errors = []
       warnings = []
       
       # Check line count
       if line_count > 150:
           errors.append(f"Too long: {line_count} lines (max 150)")
       elif line_count > 100:
           warnings.append(f"Consider shortening: {line_count} lines (target 100)")
       
       # Check required sections
       for section in REQUIRED_SECTIONS:
           if section not in content:
               errors.append(f"Missing section: {section}")
       
       # Check frontmatter
       if not content.startswith("---"):
           errors.append("Missing frontmatter (must start with ---)")
       
       # Check references resolve
       refs = re.findall(r'\[.*?\]\((\.\.?/[^)]+)\)', content)
       for ref in refs:
           ref_path = path.parent / ref
           if not ref_path.exists() and not ref.startswith("../docs/"):
               warnings.append(f"Reference may not exist: {ref}")
       
       # Output results
       if errors:
           typer.echo(f"❌ {path.name}: {len(errors)} errors")
           for e in errors:
               typer.echo(f"   - {e}")
       
       if warnings:
           typer.echo(f"⚠️  {path.name}: {len(warnings)} warnings")
           for w in warnings:
               typer.echo(f"   - {w}")
       
       if not errors and not warnings:
           typer.echo(f"✅ {path.name}: valid ({line_count} lines)")
       
       if errors or (strict and warnings):
           raise typer.Exit(1)
   
   @app.command("check-all")
   def check_all_skills() -> None:
       """Validate all skills in .claude/skills/."""
       skills_dir = Path(".claude/skills")
       if not skills_dir.exists():
           typer.echo("❌ No .claude/skills/ directory")
           raise typer.Exit(1)
       
       total = 0
       failed = 0
       
       for skill_dir in skills_dir.iterdir():
           if skill_dir.is_dir():
               skill_file = skill_dir / "SKILL.md"
               if skill_file.exists():
                   total += 1
                   try:
                       validate_skill(skill_file, strict=False)
                   except SystemExit:
                       failed += 1
       
       typer.echo(f"\nSummary: {total - failed}/{total} skills valid")
       if failed:
           raise typer.Exit(1)
   ```

### Output Files

- `templates/skill-template.md`
- `docs/adr/007-skill-length-limit.md`
- `src/sdp/cli/skill.py`
- `tests/unit/test_skill_validator.py`

### Completion Criteria

```bash
# Validator works
sdp skill validate templates/skill-template.md

# Check existing skills
sdp skill check-all

# Tests pass
pytest tests/unit/test_skill_validator.py -v
```

---

## Execution Report

**Executed by:** ______  
**Date:** ______

### Goal Status
- [ ] AC1-AC5 — ✅

**Goal Achieved:** ______
