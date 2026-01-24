---
ws_id: 00-500-01
feature: F010
status: completed
size: MEDIUM
project_id: 00
---

## WS-00-500-01: Sync SDP Content to Separate Repo

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- All SDP content from `msu_ai_masters/sdp/` is copied to `github.com:fall-out-bug/sdp` repo
- Git history is preserved in the new SDP repo
- SDP repo has proper initial commit structure

**Acceptance Criteria:**
- [x] AC1: All files from `msu_ai_masters/sdp/` exist in SDP repo
- [x] AC2: Git history is preserved (or clean import if preferred)
- [x] AC3: SDP repo has README.md with project description
- [x] AC4: `.gitignore` is properly configured for Python project
- [x] AC5: `pyproject.toml` exists with proper package metadata

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

SDP (Spec-Driven Protocol) currently lives in `msu_ai_masters/sdp/` but needs to move to a separate repository at `github.com:fall-out-bug/sdp`. This workstream copies all content to the new repository while preserving git history.

The SDP repo already exists but is empty. This WS populates it with content.

### Зависимость

Независимый (первый WS в F500)

### Входные файлы

- `sdp/` — entire directory tree in msu_ai_masters
- `sdp/pyproject.toml` — existing package configuration
- `sdp/README.md` or `sdp/docs/PROTOCOL.md` — documentation

### Шаги

1. **Clone the empty SDP repo locally**
   ```bash
   cd /tmp
   git clone git@github.com:fall-out-bug/sdp.git
   cd sdp
   ```

2. **Copy content from msu_ai_masters/sdp/**
   ```bash
   # Option A: Copy with rsync (preserves structure)
   rsync -av --exclude='.git' --exclude='__pycache__' --exclude='*.pyc' \
       /home/fall_out_bug/msu_ai_masters/sdp/ /tmp/sdp/

   # Option B: Use git filter-branch (preserves history)
   # From msu_ai_masters:
   git filter-branch --subdirectory-filter sdp --prune-empty --index-filter \
       "git rm -rf --cached __pycache__ '*.pyc' .DS_Store" HEAD
   git push git@github.com:fall-out-bug/sdp.git HEAD:main
   ```

3. **Create/update README.md** in SDP repo root
   ```markdown
   # SDP - Spec-Driven Protocol

   Universal meta-protocol for agent-driven development.

   ## What is SDP?

   SDP is a workstream-driven development framework designed for AI agent
   collaboration. It provides structured workflows, validation, and tooling
   for building software with AI assistance.

   ## Installation

   SDP is typically used as a git submodule:

   ```bash
   git submodule add git@github.com:fall-out-bug/sdp.git sdp
   ```

   ## Project ID Registry

   | ID | Project | Description |
   |----|---------|-------------|
   | 00 | SDP Protocol | Universal meta-protocol |
   | 02 | hw_checker | Homework validation system |
   | 03 | mlsd | ML System Design course |
   | 04 | bdde | Big Data course |
   | 05 | msu_ai_masters | Meta-repo configuration |

   ## Documentation

   See `docs/PROTOCOL.md` for complete protocol specification.
   ```

4. **Verify `.gitignore`** exists with Python patterns:
   ```python
   # Python
   __pycache__/
   *.py[cod]
   *.so
   .Python
   build/
   develop-eggs/
   dist/
   downloads/
   eggs/
   .eggs/
   lib/
   lib64/
   parts/
   sdist/
   var/
   wheels/

   # Testing
   .pytest_cache/
   .coverage
   htmlcov/

   # IDEs
   .vscode/
   .cursor/
   .idea/

   # OS
   .DS_Store
   Thumbs.db
   ```

5. **Create initial commit**
   ```bash
   git add .
   git commit -m "feat(sdp): Initial SDP repository content

   - Copy all content from msu_ai_masters/sdp/
   - Add README.md with project description
   - Configure .gitignore for Python project
   - Establish Project ID registry (00-05)

   This establishes SDP as a standalone repository for the
   universal meta-protocol used across all projects."
   ```

6. **Push to origin**
   ```bash
   git push origin main
   ```

### Ожидаемый результат

- SDP repo at `github.com:fall-out-bug/sdp` contains all SDP code/docs
- Repo structure follows Python package best practices
- README.md clearly explains SDP purpose and installation

### Scope Estimate

- Файлов: ~100+ (all SDP content)
- Строк: ~5000 (MEDIUM)
- Токенов: ~2500

### Критерий завершения

```bash
# Verify repo exists and is accessible
git ls-remote git@github.com:fall-out-bug/sdp.git

# Clone and verify content
cd /tmp
rm -rf sdp-verify
git clone git@github.com:fall-out-bug/sdp.git sdp-verify
cd sdp-verify

# Check key files exist
test -f README.md
test -f .gitignore
test -f pyproject.toml
test -d docs
test -d src

# Verify no .pyc files in repo
! find . -name "*.pyc" | grep -q "."

# Check file count (should be substantial)
file_count=$(find . -type f | wc -l)
test $file_count -gt 50

echo "✅ SDP repo content sync verified: $file_count files"
```

### Ограничения

- НЕ изменять содержимое файлов, только копировать
- НЕ менять структуру директорий SDP
- НЕ добавлять новые зависимости в pyproject.toml
- НЕ удалять существующую документацию
