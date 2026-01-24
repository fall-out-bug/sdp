---
ws_id: 00-500-04
feature: F010
status: completed
size: SMALL
project_id: 00
---

## WS-00-500-04: Configure SDP Repo for Submodule Use

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- SDP repo is properly configured as a git submodule
- README.md explains submodule usage
- Repo has proper git history for submodule reference
- `.gitattributes` configured for proper line endings

**Acceptance Criteria:**
- [x] AC1: README.md has submodule installation section
- [x] AC2: `.gitattributes` exists with text=auto
- [x] AC3: Repo has proper tags/releases structure
- [x] AC4: No unintended files in repo (no .DS_Store, etc.)

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

SDP will be used as a git submodule in msu_ai_masters and potentially other projects. This WS ensures the SDP repo is properly configured for submodule usage.

### Зависимость

00-500-01 (SDP repo must exist with content)

### Входные файлы

- `sdp/README.md` — from 00-500-01
- `sdp/.gitignore` — from 00-500-01

### Шаги

1. **Update README.md** with submodule section

   Add to `sdp/README.md`:

   ```markdown
   ## Installation as Submodule

   SDP is typically used as a git submodule in projects:

   ```bash
   # Add SDP as submodule
   git submodule add git@github.com:fall-out-bug/sdp.git sdp

   # Update submodule
   git submodule update --remote sdp

   # Initialize submodules in fresh clone
   git submodule update --init --recursive
   ```

   ## Project Structure

   ```
   sdp/
   ├── src/sdp/          # Source code
   ├── prompts/          # Command prompts
   ├── templates/        # WS templates
   ├── docs/             # Documentation
   ├── tests/            # Test suite
   └── scripts/          # Utility scripts
   ```

   ## Versioning

   SDP follows semantic versioning. When used as a submodule, pin to a specific tag:

   ```bash
   cd sdp
   git checkout v1.0.0  # or specific commit
   cd ..
   git add sdp
   git commit -m "Pin SDP to v1.0.0"
   ```

   ## Development

   For SDP development:

   ```bash
   # Clone directly (not as submodule)
   git clone git@github.com:fall-out-bug/sdp.git
   cd sdp
   poetry install
   pytest
   ```
   ```

2. **Create `.gitattributes`** for proper line endings

   Create `sdp/.gitattributes`:

   ```text
   # Auto detect text files and normalize line endings to LF
   * text=auto eol=lf

   # Explicitly declare text files
   *.md text eol=lf
   *.py text eol=lf
   *.yml text eol=lf
   *.yaml text eol=lf
   *.toml text eol=lf
   *.txt text eol=lf
   *.sh text eol=lf

   # Declare files that should always be checked out binary
   *.png binary
   *.jpg binary
   *.jpeg binary
   *.gif binary
   *.ico binary
   *.woff binary
   *.woff2 binary
   ```

3. **Create initial release/tag structure**

   ```bash
   cd /tmp/sdp  # the SDP repo

   # Create initial tag
   git tag -a v0.1.0 -m "Initial SDP release

   - PP-FFF-SS workstream naming
   - Project ID registry (00-05)
   - Universal meta-protocol for agent-driven development
   - Submodule-ready configuration"

   git push origin v0.1.0
   ```

4. **Verify clean repo state**

   ```bash
   # Check for unintended files
   git ls-files | grep -E "\\.DS_Store|Thumbs.db|\\.pyc$" || echo "✓ No unwanted files"

   # Verify .gitignore is working
   git status --short  # Should be empty
   ```

### Ожидаемый результат

- SDP repo is submodule-ready
- README explains submodule usage
- Proper gitattributes for cross-platform development
- Initial tag created

### Scope Estimate

- Файлов: ~3 (README.md, .gitattributes)
- Строк: ~300 (SMALL)
- Токенов: ~1500

### Критерий завершения

```bash
# Clone SDP repo fresh to verify
cd /tmp
rm -rf sdp-verify
git clone git@github.com:fall-out-bug/sdp.git sdp-verify
cd sdp-verify

# Check key files
test -f README.md
test -f .gitattributes
test -f .gitignore

# Verify submodule instructions in README
grep -q "submodule" README.md
grep -q "git submodule" README.md

# Verify tag exists
git tag | grep -q "v0.1.0"

echo "✅ SDP repo configured for submodule use"
```

### Ограничения

- НЕ менять структуру SDP кода
- НЕ добавлять PyPI publishing (не в scope этого WS)
- НЕ создавать GitHub Actions (будет в другом WS)

---

## Execution Report

**Executed by:** Claude
**Date:** 2026-01-24

### Goal Status
- [x] AC1: README.md has submodule installation section — ✅
- [x] AC2: `.gitattributes` exists with text=auto — ✅
- [x] AC3: Repo has proper tags/releases structure — ✅ (v0.1.0)
- [x] AC4: No unintended files in repo (no .DS_Store, etc.) — ✅

**Goal Achieved:** ✅ YES

### Files Changed

| File | Action | LOC |
|------|--------|-----|
| `sdp/README.md` | modified | +52 (submodule section) |
| `sdp/.gitattributes` | created | 17 |

### Self-Check Results
```bash
$ cd /tmp/sdp-verify && test -f README.md
✓ README.md exists

$ test -f .gitattributes
✓ .gitattributes exists

$ test -f .gitignore
✓ .gitignore exists

$ grep -q "submodule" README.md
✓ README has submodule section

$ grep -q "git submodule" README.md
✓ README has git submodule commands

$ git tag | grep -q "v0.1.0"
✓ Tag v0.1.0 exists

$ git ls-files | grep -E "\.DS_Store|Thumbs\.db|\.pyc$"
✓ No unwanted files
```

### SDP Repo Commits
- `20b6af9` - docs(sdp): WS-00-500-04 - configure SDP for submodule use
- Tag: `v0.1.0` - Initial SDP release

### Verification
- README.md updated with "Installation as Submodule" section ✅
- Project structure documented ✅
- Versioning and pinning instructions added ✅
- .gitattributes created with text=auto eol=lf ✅
- Initial tag v0.1.0 created and pushed ✅
- Fresh clone verification passed ✅

### Next Steps
- Continue with WS-00-500-05: Add SDP as Submodule in msu_ai_masters
- After F500 completion: `/codereview F500`
