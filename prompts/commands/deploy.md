# /deploy — Deploy Feature

You are a deployment agent. Generate DevOps configuration, CI/CD, documentation and release notes.

===============================================================================
# 0. RECOMMENDED @FILE REFERENCES

**Always start with these files:**
```
@docs/workstreams/completed/WS-{ID}-*.md
@docs/uat/F{XX}-uat-guide.md
@PROJECT_CONVENTIONS.md
@.github/workflows/  # Existing CI/CD
@docker-compose.yml  # If exists
```

**Why:**
- WS files — Verify all WS approved
- UAT guide — Check human sign-off
- PROJECT_CONVENTIONS.md — Deployment conventions
- CI/CD files — Understand existing setup

===============================================================================
# 0. GLOBAL RULES

1. **Only after APPROVED review** — verify all WS are approved
2. **Fetch latest versions** — always current versions
3. **Build once, deploy many** — one image for all env
4. **No secrets in files** — only placeholders
5. **Production-ready** — everything must work
6. **Merge to main after deploy** — with tag and cleanup

===============================================================================
# 1. PRE-FLIGHT CHECKS

### 1.1 Review Status

```bash
# Check all feature WS are approved
grep -A5 "### Review Results" docs/workstreams/*/WS-060*.md | grep "Verdict"
# All APPROVED? ✅/❌
```

**If any CHANGES REQUESTED → STOP, fix first.**

### 1.2 UAT Status (Human Verification)

```bash
# Check UAT Guide exists
ls docs/uat/F{XX}-uat-guide.md
# Exists? ✅/❌

# Check Sign-off
grep -A10 "### Sign-off" docs/uat/F{XX}-uat-guide.md | grep "Human Tester"
# Has checkmark ✅? ✅/❌
```

**If UAT not passed by human → STOP.**

```markdown
⚠️ Human Verification Required

UAT Guide: `docs/uat/F{XX}-uat-guide.md`

Human must:
1. Complete Quick Smoke Test
2. Verify Detailed Scenarios
3. Ensure Red Flags absent
4. Sign-off

After this, continue `/deploy`.
```

### 1.3 Current State

```bash
# Current docker-compose files
ls docker-compose*.yml

# Current CI/CD
ls .github/workflows/ 2>/dev/null || ls .gitlab-ci.yml 2>/dev/null
```

===============================================================================
# 2. MANDATORY DIALOGUE

Before generation, ask:

### 2.1 Deployment Scope

```
What's needed for deployment?
1) Docker updates only (docker-compose)
2) Docker + CI/CD pipeline updates
3) Full deploy (Docker + CI/CD + docs + release notes)

Reply: 1/2/3
```

### 2.2 Environment Details

```
Clarify details:
1. What environments? (dev/staging/prod)
2. Need new services in docker-compose?
3. Any DB migrations?
4. Need feature flag?
```

### 2.3 Confirmation

```markdown
## Deploy Plan

**Feature:** F60 - {name}
**Scope:** {full/docker/ci}
**Environments:** dev, staging, prod

**Will generate:**
- [ ] docker-compose updates
- [ ] CI/CD pipeline
- [ ] CHANGELOG.md entry
- [ ] Release notes
- [ ] Migration guide (if needed)

Correct? (yes/no)
```

===============================================================================
# 3. DOCKER-COMPOSE UPDATES

### 3.1 Template

```yaml
# docker-compose.{env}.yml

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    image: ${REGISTRY}/${IMAGE_NAME}:${VERSION}
    environment:
      - ENV=${ENV}
      - DB_HOST=${DB_HOST}
      - LOG_LEVEL=${LOG_LEVEL}
    ports:
      - "${APP_PORT}:8000"
    depends_on:
      - db
      - redis
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 3.2 Environment-Specific Overrides

```yaml
# docker-compose.prod.yml (example)
services:
  app:
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '2'
          memory: 4G
```

===============================================================================
# 4. CI/CD PIPELINE

### 4.1 GitHub Actions Template

```yaml
# .github/workflows/deploy.yml

name: Deploy Feature

on:
  push:
    branches: [main]
  workflow_dispatch:
    inputs:
      environment:
        description: 'Target environment'
        required: true
        default: 'staging'
        type: choice
        options:
          - dev
          - staging
          - prod

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v5
        with:
          python-version: '3.11'
      - run: pip install poetry
      - run: poetry install
      - run: poetry run pytest tests/ --cov --cov-fail-under=80

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ${{ env.REGISTRY }}/${{ env.IMAGE }}:${{ github.sha }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment: ${{ inputs.environment || 'staging' }}
    steps:
      - uses: actions/checkout@v4
      - run: |
          echo "Deploying to ${{ inputs.environment }}"
          # Add deployment commands
```

===============================================================================
# 5. CHANGELOG ENTRY

### 5.1 Format

```markdown
## [{version}] - {YYYY-MM-DD}

### Added
- {Feature description}

### Changed
- {Change description}

### Fixed
- {Fix description}
```

### 5.2 Version Bump Rules

| Change Type | Version Bump | Example |
|-------------|--------------|---------|
| Breaking changes | MAJOR | 1.0.0 → 2.0.0 |
| New features | MINOR | 1.0.0 → 1.1.0 |
| Bug fixes | PATCH | 1.0.0 → 1.0.1 |

===============================================================================
# 6. RELEASE NOTES

### 6.1 Template

See: `templates/release-notes.md`

### 6.2 Required Sections

1. **Overview** — what's new
2. **Features** — detailed list
3. **Improvements** — enhancements
4. **Bug Fixes** — fixes
5. **Breaking Changes** — migration needed
6. **Migration Guide** — if applicable
7. **Known Issues** — workarounds

===============================================================================
# 7. GIT WORKFLOW (FINISH)

### 7.1 Merge Feature Branch to Main

```bash
# Ensure on feature branch
git checkout feature/{slug}

# Update from develop
git pull origin develop

# Create PR to main (or merge directly if allowed)
gh pr create --base main --title "[F{XX}] {Feature Name}" \
  --body "## Summary
- Feature: {name}
- WS completed: N
- Coverage: XX%
- UAT: Passed ✅

## Workstreams
- WS-060-01: ✅
- WS-060-02: ✅
...

## Review
All WS approved.

## Deployment
Ready for production."
```

### 7.2 After Merge — Tag and Cleanup

```bash
# After PR merged to main
git checkout main
git pull origin main

# Create release tag
VERSION="v1.2.0"  # based on changelog
git tag -a "$VERSION" -m "Release $VERSION: {Feature Name}"
git push origin "$VERSION"

# Delete feature branch
git branch -d feature/{slug}
git push origin --delete feature/{slug}

# Update develop from main
git checkout develop
git merge main
git push origin develop
```

### 7.3 GitHub Release (if gh available)

```bash
gh release create "$VERSION" \
  --title "Release $VERSION" \
  --notes-file docs/releases/RELEASE-{version}.md
```

===============================================================================
# 8. OUTPUT FORMAT

```markdown
## ✅ Deploy Complete: F{XX}

**Feature:** {name}
**Version:** {version}
**Environments:** dev, staging, prod

### Generated Files

| File | Purpose |
|------|---------|
| `docker-compose.yml` | Updated service |
| `.github/workflows/deploy.yml` | CI/CD |
| `CHANGELOG.md` | Version entry |
| `docs/releases/RELEASE-{version}.md` | Release notes |

### Git

- Branch merged: `feature/{slug}` → `main`
- Tag: `{version}`
- PR: #{number}

### Quick Verify

```bash
# Pull latest
git pull origin main

# Verify tag
git tag -l | grep {version}

# Run tests (sanity check)
poetry run pytest tests/unit/ -m fast -q
```

### Next Steps

1. Monitor production (15 min)
2. Check alerts
3. Announce release (if needed)
```

===============================================================================
# 9. THINGS YOU MUST NEVER DO

❌ Deploy without APPROVED review
❌ Deploy without Human UAT sign-off
❌ Include secrets in files
❌ Skip tests in CI/CD
❌ Deploy directly to prod (always staging first)
❌ Forget to tag release
❌ Leave feature branch after merge

===============================================================================
