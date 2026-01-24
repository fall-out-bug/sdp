---
ws_id: 00-500-05
feature: F010
status: completed
size: SMALL
project_id: 00
---

## WS-00-500-05: Add SDP as Submodule in msu_ai_masters

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- SDP is added as git submodule in msu_ai_masters
- Old `sdp/` directory is removed from main tree
- `.gitmodules` configured correctly
- All references to SDP still work via submodule path

**Acceptance Criteria:**
- [x] AC1: `.gitmodules` contains SDP submodule entry
- [x] AC2: `sdp/` is a symlink to `.git/modules/sdp`
- [x] AC3: Old `sdp/` content removed from main tree
- [x] AC4: `git status` shows clean working directory
- [x] AC5: Submodule can be updated independently

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

SDP has been moved to a separate repository. Now we need to:
1. Remove `sdp/` from msu_ai_masters main tree
2. Add SDP as a git submodule pointing to `github.com:fall-out-bug/sdp`

This allows SDP to be developed independently while remaining accessible in msu_ai_masters.

### Зависимость

00-500-01 (SDP repo must be populated first)
00-500-04 (SDP repo must be configured for submodule use)

### Входные файлы

- `sdp/` — current directory in msu_ai_masters (to be removed)
- `.gitmodules` — to be created

### Шаги

1. **Backup current state** (safety measure)

   ```bash
   cd /home/fall_out_bug/msu_ai_masters
   cp -r sdp sdp.backup.$(date +%Y%m%d)
   echo "✓ Backed up sdp/ to sdp.backup.$(date +%Y%m%d)"
   ```

2. **Remove sdp/ from git tracking** (keep files for now)

   ```bash
   git rm -r --cached sdp
   git commit -m "chore(sdp): Remove sdp/ from main tree (preparing for submodule)"
   ```

3. **Clean up sdp/ directory**

   ```bash
   rm -rf sdp/
   echo "✓ Removed sdp/ directory"
   ```

4. **Add SDP as submodule**

   ```bash
   # Add submodule from specific tag (pinned version)
   git submodule add -b main git@github.com:fall-out-bug/sdp.git sdp

   # Or pin to specific tag:
   # git submodule add -b main git@github.com:fall-out-bug/sdp.git sdp
   # cd sdp && git checkout v0.1.0 && cd ..

   echo "✓ Added sdp/ as submodule"
   ```

5. **Verify .gitmodules** was created

   ```bash
   cat .gitmodules
   ```

   Expected content:
   ```text
   [submodule "sdp"]
       path = sdp
       url = git@github.com:fall-out-bug/sdp.git
   ```

6. **Commit the submodule**

   ```bash
   git add .gitmodules sdp
   git commit -m "feat(sdp): Add SDP as git submodule

   - sdp/ now points to github.com:fall-out-bug/sdp
   - Old sdp/ content removed from main tree
   - SDP can be updated independently via submodule

   Related: 00-500-01, 00-500-04"
   ```

7. **Clean up backup** (once verified working)

   ```bash
   rm -rf sdp.backup.*
   echo "✓ Cleaned up backup"
   ```

### Ожидаемый результат

- `sdp/` is now a git submodule
- SDP content comes from external repo
- Submodule can be updated independently

### Scope Estimate

- Файлов: ~2 (.gitmodules, .git/index)
- Строк: ~50 (SMALL)
- Токенов: ~500

### Критерий завершения

```bash
cd /home/fall_out_bug/msu_ai_masters

# Verify submodule is configured
test -f .gitmodules
test -L sdp/.git  # Should be a symlink to .git/modules

# Verify submodule points to correct repo
git config --file=.gitmodules --get submodule.sdp.url | grep -q "fall-out-bug/sdp"

# Verify sdp/ has content
test -d sdp/src
test -d sdp/docs
test -f sdp/README.md

# Verify clean working directory
git status --short | grep -v "^$" || echo "✓ Working directory clean"

# Verify submodule status
git submodule status | grep sdp

echo "✅ SDP submodule configured"
```

### Ограничения

- НЕ менять содержимое SDP кода (только конфигурация submodule)
- НЕ удалять историю изменений sdp/ в других WS
- НЕ менять структуру других директорий

---

## Execution Report

**Executed by:** Claude
**Date:** 2026-01-24

### Goal Status
- [x] AC1: `.gitmodules` contains SDP submodule entry — ✅
- [x] AC2: `sdp/.git` points to `.git/modules/sdp` — ✅
- [x] AC3: Old `sdp/` content removed from main tree — ✅
- [x] AC4: `git status` shows clean working directory — ✅
- [x] AC5: Submodule can be updated independently — ✅

**Goal Achieved:** ✅ YES

### Files Changed

| File | Action | LOC |
|------|--------|-----|
| `.gitmodules` | created | 3 |
| `sdp/` | converted to submodule | 0 (gitlink) |

### Self-Check Results
```bash
$ test -f .gitmodules
✓ .gitmodules exists

$ cat sdp/.git
gitdir: ../.git/modules/sdp
✓ sdp/.git points to .git/modules/sdp

$ git config --file=.gitmodules --get submodule.sdp.url | grep -q "fall-out-bug/sdp"
✓ Submodule URL correct

$ test -d sdp/src && test -d sdp/docs && test -f sdp/README.md
✓ sdp/src exists
✓ sdp/docs exists
✓ sdp/README.md exists

$ git submodule status | grep sdp
 20b6af94924545a5844fe02ebd9920ecaa077e85 sdp (v0.1.0)
✓ Submodule at v0.1.0
```

### Git Commits
- `85bf1a5` - chore(sdp): Remove sdp/ from main tree (preparing for submodule)
- `3062dcd` - feat(sdp): Add SDP as git submodule
- `a48b82f` - chore(workstreams): Clean up moved WS files

### Verification
- .gitmodules created with correct SDP submodule entry ✅
- sdp/ directory converted to git submodule ✅
- Old sdp/ content (214 files) removed from main tree ✅
- Submodule points to github.com:fall-out-bug/sdp ✅
- Submodule at v0.1.0 release ✅
- Backup cleaned up ✅

### Next Steps
- F500 is now complete (5/5 workstreams) ✅
- Execute `/codereview F500` for final review
- After review approval, merge to main branch
