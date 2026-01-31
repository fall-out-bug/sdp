# WS-DOC-02: Documentation Reorganization Plan

## Current State
- `docs/beginner/` exists with 4 files
- Need to create `docs/reference/` and `docs/internals/`
- Need to reorganize existing docs by role

## New Structure

```
docs/
├── beginner/           # Progressive learning paths
│   ├── README.md              # Overview
│   ├── 00-quick-start.md      # Already exists
│   ├── 01-first-feature.md    # Hands-on tutorial
│   ├── 02-common-tasks.md     # Common workflows
│   ├── 03-troubleshooting.md  # From troubleshooting.md
│   └── tutorial-*.py          # Tutorial files (keep)
│
├── reference/          # Lookup documentation
│   ├── README.md              # Overview
│   ├── commands.md            # All @ and / commands
│   ├── skills.md              # Skill reference
│   ├── quality-gates.md       # From quality-gates.md
│   ├── configuration.md       # Config files
│   └── error-handling.md      # From error_patterns.md
│
├── internals/          # Maintainer docs
│   ├── README.md              # Overview
│   ├── architecture.md        # From CODE_PATTERNS.md
│   ├── extending.md           # How to extend SDP
│   ├── contributing.md        # From CONTRIBUTING.md
│   └── development.md         # Development setup
│
├── runbooks/           # Preserved (no changes)
├── workstreams/        # Preserved (no changes)
└── [existing dirs]     # Preserved (adr, guides, etc.)
```

## Tasks

### Phase 1: Create Beginner Docs ✅ (exists)
- [x] beginner/ directory exists
- [ ] Move TUTORIAL.md → beginner/01-first-feature.md
- [ ] Create beginner/02-common-tasks.md
- [ ] Move troubleshooting.md → beginner/03-troubleshooting.md

### Phase 2: Create Reference Docs
- [ ] Create reference/README.md
- [ ] Create reference/commands.md
- [ ] Move quality-gates.md → reference/quality-gates.md
- [ ] Move error_patterns.md → reference/error-handling.md
- [ ] Create reference/configuration.md
- [ ] Create reference/skills.md

### Phase 3: Create Internals Docs
- [ ] Create internals/README.md
- [ ] Create internals/architecture.md
- [ ] Create internals/extending.md
- [ ] Create internals/contributing.md
- [ ] Create internals/development.md

### Phase 4: Update Navigation
- [ ] Update START_HERE.md links
- [ ] Update SITEMAP.md
- [ ] Update README.md
- [ ] Create redirect files

## Quality Gates
- [ ] All links work (markdown-link-check)
- [ ] No broken references
- [ ] All files <200 LOC
