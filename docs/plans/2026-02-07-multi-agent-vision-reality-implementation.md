# Multi-Agent SDP + @vision + @reality Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Transform SDP into multi-agent ecosystem with strategic planning (@vision) and universal project analysis (@reality)

**Architecture:**
- **Strategic Level:** @vision (7 experts) creates product vision/PRD, @reality (8 experts) analyzes any project
- **Tactical Level:** @feature orchestrates @idea (10 experts) + @design (3-5 experts) with progressive disclosure
- **Execution Level:** @oneshot orchestrates parallel @build (two-stage review: implementer → spec → quality)
- **Quality Level:** @review spawns 6 parallel agents + synthesizer for unified decision
- **Supervision:** Circuit breaker prevents cascade failures, checkpoint/resume for long features

**Tech Stack:** Python 3.10+, Click, PyYAML, Git, Claude Code skills (markdown-based), Task tool for agent spawning

**Timeline:** 13 weeks (3 months)
**Approach:** Quality Lock-in (Phase 1 first), then parallel tracks (Quality/Speed/Synthesis/UX/Documentation)

---

## Phase 0: Preparation & Backup (1 day)

### Task 00-0: Create Safety Net

**Files:**
- Create: Git tag `before-multi-agent-vision`

**Step 1: Create backup tag**

```bash
git tag before-multi-agent-vision
git push origin before-multi-agent-vision
```

**Step 2: Verify tag created**

```bash
git tag -l before-multi-agent-vision
git log --oneline -1 before-multi-agent-vision
```

Expected: Tag shows current commit

**Step 3: Create worktree for implementation**

```bash
cd /Users/fall_out_bug/projects/vibe_coding
git worktree add sdp-multi-agent sdp/dev
cd sdp-multi-agent
```

**Step 4: Verify worktree**

```bash
pwd
git status
```

Expected: `/Users/fall_out_bug/projects/vibe_coding/sdp-multi-agent`, on `dev` branch

---

## Phase 1: Strategic Skills - @vision + @reality (2 weeks, parallel with Phase 2)

### Track 1A: @vision Skill (1 week)

### Task 01-0: Create @vision Skill Structure

**Files:**
- Create: `.claude/skills/vision/SKILL.md`
- Create: `src/sdp/vision/__init__.py`
- Create: `src/sdp/vision/extractor.py`
- Create: `tests/sdp/vision/test_extractor.py`

**Step 1: Create skill directory**

```bash
mkdir -p .claude/skills/vision
mkdir -p src/sdp/vision
mkdir -p tests/sdp/vision
```

**Step 2: Write @vision skill specification**

Create `.claude/skills/vision/SKILL.md`:

```markdown
---
name: vision
description: Strategic product planning - vision, PRD, roadmap from expert analysis
tools: Read, Write, Edit, AskUserQuestion, Task, Skill
version: 1.0.0
---

# @vision - Strategic Product Planning

**Transform project ideas into product vision, PRD, and roadmap.**

## When to Use
- **Initial project setup** - "Что мы строим?"
- **Quarterly review** - `@vision --review` - update vision based on progress
- **Major pivot** - "Меняется направление?"
- **New market entry** - "Выходим на новый рынок?"

## Workflow

### Step 1: Quick Interview (3-5 questions)

Use AskUserQuestion to understand:
- What problem are you solving?
- Who are your target users?
- What defines success in 1 year?

### Step 2: Deep-Thinking Analysis (7 Expert Agents)

Spawn parallel expert agents via Task tool:
- Product expert (product-market fit)
- Market expert (competitive landscape)
- Technical expert (feasibility)
- UX expert (user experience)
- Business expert (business model)
- Growth expert (growth strategy)
- Risk expert (risks & mitigation)

### Step 3: Synthesize & Generate Artifacts

Use synthesizer to combine expert recommendations, generate:
1. `PRODUCT_VISION.md` - why, goals, target users, success metrics
2. `docs/prd/PRD.md` - requirements, constraints, prioritized features
3. `docs/roadmap/ROADMAP.md` - Q1-Q4 milestones

### Step 4: Extract Features

Use `src/sdp/vision/extractor.py` to extract features from PRD, create drafts in `docs/drafts/`

## Outputs
- PRODUCT_VISION.md (project root)
- docs/prd/PRD.md
- docs/roadmap/ROADMAP.md
- docs/drafts/feature-*.md (5-10 drafts)

## Example

```bash
@vision "AI-powered task manager"
→ Interview (3-5 questions)
→ Deep-thinking (7 experts)
→ Artifacts generated
→ 8 feature drafts created
```

## See Also
- `.claude/skills/idea/SKILL.md` - Feature-level requirements
- `.claude/skills/reality/SKILL.md` - Reality check for completed projects
```

**Step 3: Commit skill structure**

```bash
git add .claude/skills/vision/SKILL.md
git commit -m "feat: add @vision skill specification"
```

---

### Task 01-1: Implement Vision Extractor

**Files:**
- Create: `src/sdp/vision/extractor.py`
- Test: `tests/sdp/vision/test_extractor.py`

**Step 1: Write failing test for feature extraction**

Create `tests/sdp/vision/test_extractor.py`:

```python
import pytest
from pathlib import Path
from sdp.vision.extractor import extract_features_from_prd, FeatureDraft

def test_extract_features_from_prd(tmp_path):
    """Test extracting P0 and P1 features from PRD."""
    # Create test PRD
    prd_path = tmp_path / "PRD.md"
    prd_path.write_text("""
# Product Requirements Document

## Features

### P0 (Must Have)
- Feature 1: User authentication
- Feature 2: Task creation

### P1 (Should Have)
- Feature 3: Calendar integration
- Feature 4: Notifications

### P2 (Nice to Have)
- Feature 5: Analytics
""")

    # Extract features
    features = extract_features_from_prd(str(prd_path))

    # Should extract only P0 and P1
    assert len(features) == 4
    assert features[0] == FeatureDraft(
        title="User authentication",
        description="Feature 1: User authentication",
        priority="P0"
    )
    assert features[2].priority == "P1"

def test_feature_draft_slugify():
    """Test slugifying feature titles."""
    draft = FeatureDraft(title="Calendar Integration", description="Test", priority="P1")
    assert draft.slug == "calendar-integration"
```

**Step 2: Run test to verify it fails**

```bash
pytest tests/sdp/vision/test_extractor.py -v
```

Expected: FAIL - "Module sdp.vision.extractor not found"

**Step 3: Implement minimal extractor**

Create `src/sdp/vision/extractor.py`:

```python
from dataclasses import dataclass
from pathlib import Path
import re
from typing import List

@dataclass
class FeatureDraft:
    """Feature draft extracted from PRD."""
    title: str
    description: str
    priority: str

    @property
    def slug(self) -> str:
        """Convert title to slug."""
        return re.sub(r'[^\w\s-]', '', self.title.lower().strip().replace(' ', '-'))

def extract_features_from_prd(prd_path: str) -> List[FeatureDraft]:
    """
    Extract P0 and P1 features from PRD.

    Args:
        prd_path: Path to PRD.md

    Returns:
        List of FeatureDraft
    """
    prd_content = Path(prd_path).read_text()
    features = []
    current_priority = None

    for line in prd_content.split('\n'):
        # Match "### P0 (Must Have)" headers
        priority_match = re.match(r'### (P0|P1)\s+\(', line)
        if priority_match:
            current_priority = priority_match.group(1)
            continue

        # Match feature items "- Feature N: Title"
        if current_priority in ['P0', 'P1']:
            feature_match = re.match(r'-\s+\d+:\s+(.+)', line)
            if feature_match:
                title = feature_match.group(1).strip()
                features.append(FeatureDraft(
                    title=title,
                    description=f"Feature: {title}",
                    priority=current_priority
                ))

    return features
```

**Step 4: Run test to verify it passes**

```bash
pytest tests/sdp/vision/test_extractor.py -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add src/sdp/vision/extractor.py tests/sdp/vision/test_extractor.py
git commit -m "feat: implement PRD feature extractor"
```

---

### Task 01-2: Update CLAUDE.md with @vision

**Files:**
- Modify: `CLAUDE.md`

**Step 1: Add @vision to decision tree**

Add after "Quick Reference" section:

```markdown
### Strategic Planning: @vision

**Before features, start with vision:**

```bash
@vision "AI-powered task manager"
→ Product vision, PRD, roadmap
→ 5-10 feature drafts extracted
```

**When to use @vision:**
- Starting new project
- Quarterly strategic review
- Major pivot or direction change

**Output artifacts:**
- PRODUCT_VISION.md (why, goals, metrics)
- docs/prd/PRD.md (requirements, features)
- docs/roadmap/ROADMAP.md (Q1-Q4 milestones)
```

**Step 2: Update main decision tree**

Modify the decision tree section to include:

```markdown
## Decision Tree: Strategic vs Tactical

**Start with @vision** (strategic), then @feature (tactical):

```
@vision → @feature → @idea → @design → @oneshot/@build → @review → @deploy
(strategic)  (tactical)  (requirements) (workstreams)  (execution)  (quality)   (shipping)
```
```

**Step 3: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: add @vision to CLAUDE.md decision tree"
```

---

### Track 1B: @reality Skill (1 week, parallel with 1A)

### Task 02-0: Create @reality Skill Structure

**Files:**
- Create: `.claude/skills/reality/SKILL.md`
- Create: `src/sdp/reality/__init__.py`
- Create: `src/sdp/reality/scanner.py`
- Create: `src/sdp/reality/detectors.py`
- Create: `tests/sdp/reality/test_scanner.py`

**Step 1: Create skill directory**

```bash
mkdir -p .claude/skills/reality
mkdir -p src/sdp/reality
mkdir -p tests/sdp/reality
```

**Step 2: Write @reality skill specification**

Create `.claude/skills/reality/SKILL.md`:

```markdown
---
name: reality
description: Universal project analyzer - works with any codebase, SDP or not
tools: Read, Glob, Grep, Task, Write, Edit, Bash
version: 1.0.0
---

# @reality - Universal Project Analyzer

**Analyze any project and discover what you actually built.**

## When to Use
- **Reality check** - "Что я вообще построил?"
- **Drift analysis** - Compare PRODUCT_VISION.md vs actual code
- **Legacy audit** - Understand unfamiliar codebase
- **Quick scan** - `@reality --quick` (5-10 min)
- **Deep analysis** - `@reality --deep` (30-60 min, 8 expert agents)

## Modes

### --quick (5-10 minutes)
Detect: language, frameworks, architecture pattern, LOC
Output: REALITY.md with overview

### --deep (30-60 minutes)
Spawn 8 parallel expert agents:
- Architect (architecture patterns)
- Quality (code quality, debt)
- Security (vulnerabilities, OWASP Top 10)
- Performance (bottlenecks)
- Dependencies (outdated, vulnerable)
- Documentation (completeness)
- Operations (deployment, monitoring)
- Testing (coverage, strategy)

Output: REALITY.md + DRIFT_ANALYSIS.md (if PRODUCT_VISION.md exists)

### --focus=TOPIC
Analyze single aspect: `@reality --focus=security,performance`

## Workflow

### Step 1: Detect Project Type

Scan for markers:
- package.json → Node.js
- requirements.txt/pyproject.toml → Python
- pom.xml/build.gradle → Java
- go.mod → Go

### Step 2: Scan Structure

Count LOC, identify components, detect frameworks

### Step 3: Analyze (deep mode)

Spawn expert agents in parallel via Task tool

### Step 4: Generate Artifacts

- REALITY.md - what you actually built
- DRIFT_ANALYSIS.md - vision vs reality gaps (if vision exists)

## Outputs
- REALITY.md
- DRIFT_ANALYSIS.md (optional)

## Example

```bash
# Any project, even without SDP
@reality --quick
→ "Python Django project, 12 apps, ~15K LOC"

@reality --deep
→ 8 expert agents analyze
→ REALITY.md + DRIFT_ANALYSIS.md

@reality --focus=security
→ Security expert only
```

## See Also
- `.claude/skills/vision/SKILL.md` - Create vision (before building)
- `.claude/skills/review/SKILL.md` - Quality review (after building)
```

**Step 3: Commit skill structure**

```bash
git add .claude/skills/reality/SKILL.md
git commit -m "feat: add @reality skill specification"
```

---

### Task 02-1: Implement Project Scanner

**Files:**
- Create: `src/sdp/reality/scanner.py`
- Create: `src/sdp/reality/detectors.py`
- Test: `tests/sdp/reality/test_scanner.py`

**Step 1: Write failing test for project scanning**

Create `tests/sdp/reality/test_scanner.py`:

```python
import pytest
from pathlib import Path
from sdp.reality.scanner import ProjectScanner, ProjectInfo

def test_scan_python_project(tmp_path):
    """Test scanning Python project."""
    # Create test project
    (tmp_path / "requirements.txt").write_text("django\n")
    (tmp_path / "pyproject.toml").write_text("[tool.poetry]\n")
    (tmp_path / "src").mkdir()
    (tmp_path / "src" / "app.py").write_text("# App\n" + "\n".join(["# Line"] * 50))

    # Scan
    scanner = ProjectScanner(str(tmp_path))
    info = scanner.scan()

    # Verify
    assert info.language == "python"
    assert "django" in info.frameworks
    assert info.loc > 0

def test_detect_nodejs_project(tmp_path):
    """Test detecting Node.js project."""
    (tmp_path / "package.json").write_text('{"name": "test"}\n')

    scanner = ProjectScanner(str(tmp_path))
    info = scanner.scan()

    assert info.language == "nodejs"
```

**Step 2: Run test to verify it fails**

```bash
pytest tests/sdp/reality/test_scanner.py -v
```

Expected: FAIL - "Module not found"

**Step 3: Implement scanner**

Create `src/sdp/reality/scanner.py`:

```python
from pathlib import Path
from dataclasses import dataclass
from typing import List, Dict
from sdp.reality.detectors import FrameworkDetector

@dataclass
class ProjectInfo:
    """Basic project information."""
    language: str
    frameworks: List[str]
    structure: Dict[str, str]
    loc: int

class ProjectScanner:
    """Universal project scanner."""

    MARKERS = {
        "package.json": "nodejs",
        "requirements.txt": "python",
        "pyproject.toml": "python",
        "pom.xml": "java",
        "build.gradle": "java",
        "go.mod": "go",
        "Cargo.toml": "rust",
    }

    def __init__(self, root_dir: str):
        self.root = Path(root_dir)

    def scan(self) -> ProjectInfo:
        """Scan project and return basic info."""
        language = self._detect_language()
        frameworks = self._detect_frameworks(language) if language != "unknown" else []
        structure = self._scan_structure()
        loc = self._count_loc(structure)

        return ProjectInfo(
            language=language,
            frameworks=frameworks,
            structure=structure,
            loc=loc
        )

    def _detect_language(self) -> str:
        """Detect primary language from markers."""
        for marker, lang in self.MARKERS.items():
            if (self.root / marker).exists():
                return lang
        return "unknown"

    def _detect_frameworks(self, language: str) -> List[str]:
        """Detect frameworks for given language."""
        detector = FrameworkDetector(self.root, language)
        return detector.detect()

    def _scan_structure(self) -> Dict[str, str]:
        """Scan directory structure."""
        exclusions = {".git", "node_modules", "__pycache__", "venv", ".venv"}
        structure = {}

        for path in self.root.rglob("*"):
            rel_path = str(path.relative_to(self.root))
            if any(excl in rel_path for excl in exclusions):
                continue

            if path.is_dir():
                structure[rel_path] = "dir"
            elif path.is_file():
                structure[rel_path] = "file"

        return structure

    def _count_loc(self, structure: Dict[str, str]) -> int:
        """Count lines of code."""
        code_extensions = {".py", ".js", ".ts", ".java", ".go", ".rs"}
        loc = 0

        for path, type_ in structure.items():
            if type_ == "file":
                ext = Path(path).suffix
                if ext in code_extensions:
                    full_path = self.root / path
                    try:
                        loc += len(full_path.read_text().split('\n'))
                    except:
                        pass  # Binary or unreadable

        return loc
```

Create `src/sdp/reality/detectors.py`:

```python
from pathlib import Path
from typing import List

class FrameworkDetector:
    """Detect frameworks for a given language."""

    PYTHON_FRAMEWORKS = {
        "django": ["django", "Django"],
        "flask": ["flask", "Flask"],
        "fastapi": ["fastapi", "FastAPI"],
        "celery": ["celery", "Celery"],
    }

    def __init__(self, root: Path, language: str):
        self.root = root
        self.language = language

    def detect(self) -> List[str]:
        """Detect frameworks."""
        if self.language == "python":
            return self._detect_python()
        elif self.language == "nodejs":
            return self._detect_nodejs()
        return []

    def _detect_python(self) -> List[str]:
        """Detect Python frameworks."""
        frameworks = []

        # Check requirements.txt
        req_file = self.root / "requirements.txt"
        if req_file.exists():
            content = req_file.read_text()
            for framework, markers in self.PYTHON_FRAMEWORKS.items():
                if any(marker in content for marker in markers):
                    frameworks.append(framework.capitalize())

        # Check pyproject.toml
        pyproject = self.root / "pyproject.toml"
        if pyproject.exists():
            content = pyproject.read_text()
            for framework, markers in self.PYTHON_FRAMEWORKS.items():
                if any(marker in content for marker in markers):
                    if framework.capitalize() not in frameworks:
                        frameworks.append(framework.capitalize())

        return frameworks

    def _detect_nodejs(self) -> List[str]:
        """Detect Node.js frameworks."""
        package_json = self.root / "package.json"
        if not package_json.exists():
            return []

        import json
        content = json.loads(package_json.read_text())
        deps = content.get("dependencies", {})
        dev_deps = content.get("devDependencies", {})

        frameworks = []
        for package in list(deps.keys()) + list(dev_deps.keys()):
            if package in ["express", "fastify"]:
                frameworks.append(package.capitalize())
            elif package == "react":
                frameworks.append("React")
            elif package == "vue":
                frameworks.append("Vue")

        return frameworks
```

**Step 4: Run tests to verify they pass**

```bash
pytest tests/sdp/reality/test_scanner.py -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add src/sdp/reality/ tests/sdp/reality/
git commit -m "feat: implement universal project scanner"
```

---

### Task 02-2: Update CLAUDE.md with @reality

**Files:**
- Modify: `CLAUDE.md`

**Step 1: Add @reality section**

Add after @vision section:

```markdown
### Reality Check: @reality

**After building, verify what you actually shipped:**

```bash
@reality --quick
→ Language, frameworks, architecture overview

@reality --deep
→ 8 expert agents analyze project
→ REALITY.md + DRIFT_ANALYSIS.md
```

**When to use @reality:**
- Quarterly: compare vision vs reality
- Legacy projects: understand unfamiliar code
- Any project: works without SDP
```

**Step 2: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: add @reality to CLAUDE.md"
```

---

## Phase 2: Multi-Agent Foundation - Two-Stage Review (2 weeks, QUALITY LOCK-IN)

**CRITICAL:** This phase MUST complete first. All other phases depend on two-stage review.

### Task 03-0: Create Implementer Agent

**Files:**
- Create: `.claude/agents/implementer.md`

**Step 1: Write implementer agent specification**

Create `.claude/agents/implementer.md`:

```markdown
# Implementer Subagent

You are a TDD implementer. Your job: execute workstream with discipline.

## Input
- WS file: `docs/workstreams/backlog/{ws_id}-*.md`
- Goal, Acceptance Criteria, Scope files

## Your Workflow

### Step 1: Read WS File
Use Read tool to read the workstream file. Extract:
- Goal
- Acceptance Criteria (AC)
- Scope files (if listed)
- Steps to execute

### Step 2: Activate Guard
Run: `sdp guard activate {ws_id}`

If fails, WS is not ready. Report error and stop.

### Step 3: TDD Cycle for Each AC

For each Acceptance Criteria:

**3a. Write Failing Test (Red)**
- Create test file or add test
- Run pytest
- Verify it FAILS (expected)

**3b. Implement Minimum Code (Green)**
- Write minimal code to pass test
- Run pytest
- Verify it PASSES

**3c. Refactor**
- Improve code quality
- Run pytest again
- Verify still PASSES

### Step 4: Quality Check
Run: `sdp quality check --module {module}`

Must pass:
- Coverage ≥80%
- mypy --strict
- ruff (no errors)
- Files <200 LOC

If fails, fix and re-run.

### Step 5: Commit
```bash
git add .
git commit -m "feat({scope}): {ws_id} - {title}"
```

### Step 6: Self-Review
Report what you implemented.

## Output Format

Report in this EXACT format:

```markdown
## Implementation Report: {ws_id}

### Acceptance Criteria Status
- AC1: {description} → ✅ PASS / ❌ FAIL
- AC2: {description} → ✅ PASS / ❌ FAIL

### Files Changed
- {path} (NEW/MODIFIED, {LOC} LOC)

### Test Results
- pytest: {passing}/{total} passing
- Coverage: {percentage}%

### Deviations from Spec
- None / OR: Describe any deviations

### Commit SHA
{commit_hash}
```

## Quality Standards
- TDD discipline: Red → Green → Refactor
- Coverage ≥80%
- Type hints everywhere
- No TODO/FIXME without new WS
- Clean Architecture compliance

## Example
```
Input: WS 00-001-01 "Domain entities"
AC1: User entity with id, email
AC2: Order entity with id, user_id

→ TDD cycle for AC1
→ TDD cycle for AC2
→ Quality check passes
→ Commit: abc123

Report:
## Implementation Report: 00-001-01

### Acceptance Criteria Status
- AC1: User entity → ✅ PASS
- AC2: Order entity → ✅ PASS

### Files Changed
- src/domain/user.py (NEW, 45 LOC)
- src/domain/order.py (NEW, 38 LOC)
- tests/test_domain.py (NEW, 120 LOC)

### Test Results
- pytest: 47/47 passing
- Coverage: 87%

### Deviations from Spec
- None

### Commit SHA
abc123def
```

## CRITICAL RULES
1. **DO NOT skip TDD** - Tests first, always
2. **DO NOT commit if quality check fails**
3. **DO report deviations** - If you deviate from spec, explain why
4. **DO ask questions** - If WS is unclear, ask before implementing
```

**Step 2: Commit**

```bash
git add .claude/agents/implementer.md
git commit -m "feat: add implementer agent specification"
```

---

### Task 03-1: Create Spec Compliance Reviewer Agent

**Files:**
- Create: `.claude/agents/spec-reviewer.md`

**Step 1: Write spec reviewer specification**

Create `.claude/agents/spec-reviewer.md`:

```markdown
# Spec Compliance Reviewer

You are a SPEC REVIEWER. Your job: **Verify implementation matches requirements.**

## Golden Rule
**DO NOT TRUST THE IMPLEMENTER REPORT.**

Read the actual code. Compare to actual requirements.

## Input
- WS file: `docs/workstreams/backlog/{ws_id}-*.md`
- Implementer report (claims)
- Git commit SHA (to inspect actual code)

## Your Workflow

### Step 1: Read WS File
Extract Acceptance Criteria (ground truth).

### Step 2: Read Implementer Report
Understand what implementer claims.

### Step 3: **READ ACTUAL CODE**
Use Read tool to inspect ALL scope files listed in WS.
Do NOT trust implementer report. Verify everything.

### Step 4: Compare Requirements vs Reality

For each AC:
- What does WS require?
- What does implementer claim?
- What ACTUALLY exists in code?

### Step 5: Check Drift
- All scope_files exist?
- All functions/classes from docs exist in code?
- File purpose matches documentation?
- No TODO/FIXME/HACK in production code?

## Output Format

Report in this EXACT format:

```markdown
## Spec Compliance Report: {ws_id}

### Verdict: ✅ PASS / ❌ FAIL

### Acceptance Criteria Status

#### AC1: {description}
- **Required:** {what WS asks for}
- **Claimed:** {what implementer says}
- **Reality:** {what actually exists in code}
- **Status:** ✅ MATCH / ❌ MISMATCH

[Repeat for all AC]

### Issues Found
- ❌ {file}:{line} - {specific issue}
OR
- None

### Files Checked
- {file1} - ✅ EXISTS / ❌ MISSING
- {file2} - ✅ EXISTS / ❌ MISSING

### Recommendation
- **PASS** - Spec compliant, proceed to quality review
- **FAIL** - Issues found, implementer must fix:
  - {list of specific fixes}

### Severity
- **Critical** - Block quality review, must fix
- **Minor** - Can proceed, note for future
```

## Review Checklist
- [ ] All scope_files exist?
- [ ] All functions/classes documented actually exist?
- [ ] File purpose matches documentation?
- [ ] No TODO/FIXME/HACK in production code?
- [ ] All Acceptance Criteria satisfied?
- [ ] No unapproved deviations from spec?

## Examples

### Example 1: PASS
```
AC1: User entity with id, email
Required: src/domain/user.py with User(id, email)
Claimed: "Created User entity"
Reality: src/domain/user.py exists, has User(dataclass) with id: int, email: str
Status: ✅ MATCH

Verdict: ✅ PASS
```

### Example 2: FAIL
```
AC1: User entity with id, email
Required: src/domain/user.py with User(id, email)
Claimed: "Created User entity"
Reality: src/domain/user.py exists, but User only has id, missing email
Status: ❌ MISMATCH - Missing email field

Verdict: ❌ FAIL
Recommendation: Implementer must add email field to User entity
```

## CRITICAL RULES
1. **NEVER trust implementer report** - Always read actual code
2. **BE specific** - Point to exact files and lines
3. **CHECK drift** - Look for TODO/FIXME, missing files
4. **REQUIRE all AC** - All acceptance criteria must be satisfied
```

**Step 2: Commit**

```bash
git add .claude/agents/spec-reviewer.md
git commit -m "feat: add spec compliance reviewer agent"
```

---

### Task 03-2: Update @build Skill for Two-Stage Review

**Files:**
- Modify: `.claude/skills/build/SKILL.md`

**Step 1: Read current @build skill**

```bash
Read .claude/skills/build/SKILL.md
```

**Step 2: Modify workflow section**

Replace existing workflow with two-stage review:

```markdown
## Workflow

### Step 0: Resolve Task ID
```bash
# Input: ws_id (00-001-01) OR beads_id (sdp-xxx)
# Resolve ws_id ↔ beads_id via .beads-sdp-mapping.jsonl
```

### Step 1: Beads IN_PROGRESS (when Beads enabled)
```bash
[ -n "$beads_id" ] && bd update "$beads_id" --status in_progress
```

### Step 2: Activate Guard
```bash
sdp guard activate {ws_id}
```

### Step 3: Read Workstream
```bash
Read("docs/workstreams/backlog/{WS-ID}-*.md")
```

### Step 4: Two-Stage Review

#### Stage 1: Implementer
```python
Task(
    subagent_type="general-purpose",
    prompt=f"""
    Read .claude/agents/implementer.md

    WORKSTREAM: {ws_id}
    GOAL: {goal}
    ACCEPTANCE CRITERIA: {AC}

    Execute TDD cycle and report.
    """,
    description=f"Implement {ws_id}"
)
```

**Gate:** Implementer must provide implementation report.

#### Stage 2: Spec Compliance Reviewer
```python
spec_report = Task(
    subagent_type="general-purpose",
    prompt=f"""
    Read .claude/agents/spec-reviewer.md

    WS: {ws_id}
    REQUIREMENTS: {ws_file['acceptance_criteria']}
    IMPLEMENTER CLAIMS: {implementer_report}
    COMMIT SHA: {commit_sha}

    CRITICAL: Do not trust report. Read actual code.
    """,
    description=f"Spec review {ws_id}"
)
```

**Gate:** If spec reviewer reports FAIL → Stage 1 fixes → Re-review (max 2 retries)

#### Stage 3: Code Quality Reviewer
```python
quality_report = Task(
    subagent_type="general-purpose",
    prompt=f"""
    Read .claude/agents/reviewer.md

    BASE_SHA: {commit_before_ws}
    HEAD_SHA: {current_commit}

    Review quality: coverage, mypy, ruff, architecture
    """,
    description=f"Quality review {ws_id}"
)
```

**Gate:** If quality reviewer reports FAIL → Stage 1 fixes → Re-review (max 2 retries)

### Step 5: All Stages Pass
```bash
# When Beads enabled: sync before commit
[ -d .beads ] && bd sync

sdp guard complete {ws_id}
git add .
git commit -m "feat({scope}): {ws_id} - {title}"
```

### Step 6: Beads CLOSED (when Beads enabled)
```bash
[ -n "$beads_id" ] && bd close "$beads_id" --reason "WS completed via two-stage review"
```
```

**Step 3: Commit**

```bash
git add .claude/skills/build/SKILL.md
git commit -m "feat: @build now uses two-stage review (implementer → spec → quality)"
```

---

[План продолжается в следующей части из-за размера...]

---

**Это начало implementation plan. План включает:**

✅ **Phase 0:** Backup & worktree setup
✅ **Phase 1:** @vision + @reality skills (2 weeks, parallel)
✅ **Phase 2:** Two-stage review foundation (2 weeks) — QUALITY LOCK-IN
📝 **Phase 3-5:** Speed track, Synthesis track, UX track (remaining weeks)

**Plan incomplete — need to continue with:**
- Phase 3: Parallel execution (@oneshot enhancement)
- Phase 4: Agent synthesizer
- Phase 5: Progressive disclosure for @idea/@design
- Phase 6: Documentation (agent catalog, migration guide)

**Хотите чтобы я:**
1. **Продолжил план** (дописать остальные фазы)?
2. **Начать execution** с текущего плана (Phases 0-2 сначала)?
3. **Создать отдельные планы** для каждой фазы?
4. **Что-то другое**?
