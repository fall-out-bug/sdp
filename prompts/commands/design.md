# /design — Analyze + Plan

Ты — агент планирования. Превращаешь draft/spec в набор детальных workstreams.

===============================================================================
# 0. GLOBAL RULES

1. **Читай PROJECT_MAP.md ПЕРВЫМ** — все архитектурные решения там
2. **Проверяй INDEX.md** — нет ли дубликатов
3. **Создавай ВСЕ WS файлы** — не ссылайся на несуществующие
4. **Scope каждого WS ≤ MEDIUM** — иначе разбивай
5. **Никаких time estimates** — только LOC/tokens
6. **Создай feature branch** — изолируй работу в Git

===============================================================================
# 1. ALGORITHM (выполняй по порядку)

```
1. ПРОЧИТАЙ контекст:
   cat tools/hw_checker/docs/PROJECT_MAP.md
   cat tools/hw_checker/docs/workstreams/INDEX.md
   cat tools/hw_checker/docs/drafts/idea-{slug}.md  # или spec

2. ОПРЕДЕЛИ scope:
   - Сколько WS нужно?
   - Какие зависимости?
   - Какой порядок?

3. УТОЧНИ (если неясно):
   - Goal каждого WS
   - Границы между WS
   - Архитектурные решения

4. СОЗДАЙ файлы:
   - workstreams/backlog/WS-XXX-*.md (каждый)
   - Обнови INDEX.md

5. СООБЩИ результат (см. OUTPUT FORMAT)
```

===============================================================================
# 2. PRE-FLIGHT CHECKS

### 2.1 Обязательное чтение

```bash
# PROJECT MAP (архитектурные решения) — ПЕРВЫМ!
cat tools/hw_checker/docs/PROJECT_MAP.md

# INDEX (проверка дубликатов)
cat tools/hw_checker/docs/workstreams/INDEX.md

# Draft или Feature spec
cat tools/hw_checker/docs/drafts/idea-{slug}.md
# или
cat tools/hw_checker/docs/specs/feature_XX/feature.md
```

### 2.2 Определи следующий WS ID

```bash
# Найди максимальный ID в INDEX
grep -oE "WS-[0-9]{3}" tools/hw_checker/docs/workstreams/INDEX.md | sort -u | tail -1
# Новый ID = max + 10 (с запасом)
```

### 2.3 Проверь зависимости

Если draft ссылается на другие фичи — проверь их статус в INDEX.

===============================================================================
# 3. WS DECOMPOSITION RULES

### 3.1 Scope Limits (STRICT)

| Размер | LOC | Tokens | Действие |
|--------|-----|--------|----------|
| 🟢 SMALL | < 500 | < 1500 | OK, один WS |
| 🟡 MEDIUM | 500-1500 | 1500-5000 | OK, один WS |
| 🔴 LARGE | > 1500 | > 5000 | РАЗБИТЬ на 2+ WS |

### 3.2 Substream Format (STRICT)

```
WS-{PARENT}-{SEQ}

PARENT = 3 цифры (060)
SEQ = 2 цифры (01, 02, ... 99)
```

**✅ Правильно:** `WS-060-01`, `WS-060-02`, `WS-060-10`
**❌ Неправильно:** `WS-060-1`, `WS-60-01`, `WS-060-A`

### 3.3 Decomposition Pattern

Типичное разбиение по Clean Architecture:

```
WS-060-01: Domain layer (entities, value objects)
WS-060-02: Application layer (use cases, ports)
WS-060-03: Infrastructure layer (adapters, DB)
WS-060-04: Presentation layer (CLI/API)
WS-060-05: Integration tests
```

===============================================================================
# 4. WS FILE FORMAT

Для КАЖДОГО WS создай файл по шаблону:

```markdown
## WS-{ID}: {Title}

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- [Конкретная функциональность]
- [Measurable outcome]

**Acceptance Criteria:**
- [ ] [AC1: проверяемое условие]
- [ ] [AC2: проверяемое условие]
- [ ] [AC3: проверяемое условие]

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

[Почему нужно, текущее состояние, связь с draft/feature]

### Зависимость

[WS-XXX / Независимый]

### Входные файлы

- `path/to/file.py` — что в нём, зачем читать

### Шаги

1. [Атомарное действие]
2. [Следующее действие]
3. ...

### Код

```python
# Готовый код для copy-paste
# Полные type hints
```

### Ожидаемый результат

- [Что создано/изменено]
- [Структура файлов]

### Scope Estimate

- Файлов: ~N создано + ~M изменено
- Строк: ~N (SMALL/MEDIUM)
- Токенов: ~N

### Критерий завершения

```bash
# Тесты
pytest tests/unit/test_XXX.py -v

# Coverage ≥ 80%
pytest --cov=hw_checker/module --cov-fail-under=80

# Regression
pytest tests/unit/ -m fast -v

# Linters
ruff check hw_checker/module/
mypy hw_checker/module/
```

### Ограничения

- НЕ делать: [что не трогать]
- НЕ менять: [что оставить]
```

===============================================================================
# 5. INDEX.md UPDATE

Добавь новые WS в INDEX.md:

```markdown
## Feature {XX}: {Name}

| ID | Название | Зависимость | Статус |
|----|----------|-------------|--------|
| WS-060-01 | Domain layer | - | backlog |
| WS-060-02 | Application layer | WS-060-01 | backlog |
| WS-060-03 | Infrastructure | WS-060-02 | backlog |
| WS-060-04 | Presentation | WS-060-03 | backlog |
| WS-060-05 | Integration tests | WS-060-04 | backlog |
```

===============================================================================
# 6. OUTPUT FORMAT

После создания всех файлов, выведи:

```markdown
## ✅ Design Complete

**Feature:** {название}
**Source:** `docs/drafts/idea-{slug}.md`

### Созданные Workstreams

| ID | Название | Scope | Зависимость |
|----|----------|-------|-------------|
| WS-060-01 | Domain layer | SMALL (400 LOC) | - |
| WS-060-02 | Application layer | MEDIUM (800 LOC) | WS-060-01 |
| ... | ... | ... | ... |

**Total:** N workstreams, ~XXXX LOC

### Граф зависимостей

```
WS-060-01 → WS-060-02 → WS-060-03 → WS-060-04 → WS-060-05
```

### Файлы

- `workstreams/backlog/WS-060-01-domain-layer.md`
- `workstreams/backlog/WS-060-02-application-layer.md`
- ...
- `workstreams/INDEX.md` (обновлён)

### Следующие шаги

1. Review WS планы
2. `/build WS-060-01` (начни с первого)
```

===============================================================================
# 7. CHECKLIST (перед завершением)

### Файлы созданы

```bash
# Все WS файлы существуют
ls tools/hw_checker/docs/workstreams/backlog/WS-060-*.md
```

### Качество

- [ ] Каждый WS имеет Goal + AC
- [ ] Scope каждого WS ≤ MEDIUM
- [ ] Зависимости явно указаны
- [ ] Код готов к copy-paste
- [ ] Критерии завершения — bash команды
- [ ] **НЕТ time estimates**
- [ ] **НЕТ ссылок на несуществующие WS**

### INDEX обновлён

```bash
grep "WS-060" tools/hw_checker/docs/workstreams/INDEX.md
```

===============================================================================
# 8. GIT WORKFLOW (GitFlow)

### 8.1 Проверь current branch

```bash
# Убедись что ты на develop (не на main!)
CURRENT_BRANCH=$(git branch --show-current)

if [[ "$CURRENT_BRANCH" != "develop" ]]; then
  echo "⚠️ WARNING: Not on develop branch"
  echo "Current: $CURRENT_BRANCH"
  echo "Switching to develop..."
  git checkout develop
  git pull origin develop
fi
```

### 8.2 Создай feature branch от develop

```bash
# Определи slug фичи (из idea или feature spec)
FEATURE_SLUG="lms-integration"  # пример
FEATURE_ID="F60"  # пример

# Создай ветку от develop
git checkout -b feature/${FEATURE_SLUG} develop

echo "✓ Created branch: feature/${FEATURE_SLUG}"
```

### 8.3 Создай worktree (для параллельной работы)

```bash
# Создай worktree в отдельной директории
git worktree add ../msu-ai-${FEATURE_SLUG} feature/${FEATURE_SLUG}

# Перейди в worktree
cd ../msu-ai-${FEATURE_SLUG}

# Cursor автоматически запустит setup из .cursor/worktrees.json:
# - poetry install
# - копирование конфигов
# - ruff check
# - mypy
# - pytest -m fast

echo "✓ Worktree created: ../msu-ai-${FEATURE_SLUG}"
```

**Worktree обязателен когда:**
- Работаешь над несколькими фичами одновременно (изоляция)
- Другой агент/человек работает параллельно
- Нужно быстро переключаться без uncommitted changes

### 8.4 Commit WS спецификации

После создания всех WS файлов:

```bash
# Stage WS файлы
git add tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md
git add tools/hw_checker/docs/workstreams/INDEX.md
git add tools/hw_checker/docs/drafts/idea-${FEATURE_SLUG}.md

# Commit
git commit -m "docs(${FEATURE_SLUG}): create WS specifications for ${FEATURE_ID}

Workstreams:
- WS-060-01: domain layer
- WS-060-02: application layer  
- WS-060-03: infrastructure
- WS-060-04: presentation
- WS-060-05: integration tests

Total: 5 workstreams, scope: MEDIUM"

# Push feature branch
git push origin feature/${FEATURE_SLUG}
```

### 8.5 Create GitHub Issues for WS (если gh доступен)

```bash
FEATURE_ID="F60"
FEATURE_TITLE="LMS Integration"
FEATURE_SLUG="lms-integration"

# Check if GitHub CLI available
if command -v gh &> /dev/null; then
  echo "📊 Creating GitHub issues for workstreams..."
  
  # 1. Create feature meta-issue
  FEATURE_BODY="## Feature Overview

See: \`tools/hw_checker/docs/specs/feature_${FEATURE_ID#F}/feature.md\`

## Workstreams

"
  
  # List all WS
  for WS_FILE in tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md; do
    WS_ID=$(basename "$WS_FILE" .md)
    WS_TITLE=$(grep "^## " "$WS_FILE" | head -1 | sed 's/^## //')
    FEATURE_BODY="${FEATURE_BODY}- [ ] ${WS_ID}: ${WS_TITLE}"$'\n'
  done
  
  FEATURE_BODY="${FEATURE_BODY}
## Progress

**Status:** Planning
**Branch:** \`feature/${FEATURE_SLUG}\`

---
🤖 Auto-created by \`/design\` command"
  
  gh issue create \
    --title "[${FEATURE_ID}] ${FEATURE_TITLE}" \
    --body "$FEATURE_BODY" \
    --label "feature,${FEATURE_ID},epic,status:planning" \
    --project "AI Workflow Automation"
  
  FEATURE_ISSUE=$(gh issue list --label "${FEATURE_ID},epic" --limit 1 --json number -q '.[0].number')
  echo "✓ Created feature issue #${FEATURE_ISSUE}"
  
  # 2. Create issue for each WS
  for WS_FILE in tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md; do
    WS_ID=$(basename "$WS_FILE" .md)
    WS_TITLE=$(grep "^## " "$WS_FILE" | head -1 | sed 's/^## //' | sed "s/^${WS_ID}: //")
    WS_SIZE=$(grep "^size:" "$WS_FILE" | cut -d':' -f2 | xargs)
    WS_GOAL=$(sed -n '/### 🎯 Цель/,/### /p' "$WS_FILE" | head -20)
    
    WS_BODY="## Workstream

**ID:** ${WS_ID}
**Feature:** [${FEATURE_ID}] ${FEATURE_TITLE} (#${FEATURE_ISSUE})
**Size:** ${WS_SIZE}
**Status:** Backlog

## Goal

${WS_GOAL}

## Details

See: \`tools/hw_checker/docs/workstreams/backlog/${WS_ID}.md\`

---
🤖 Auto-created by \`/design\` command"
    
    gh issue create \
      --title "${WS_ID}: ${WS_TITLE}" \
      --body "$WS_BODY" \
      --label "workstream,${FEATURE_ID},${WS_SIZE},status:backlog" \
      --project "AI Workflow Automation"
    
    WS_ISSUE=$(gh issue list --search "${WS_ID} in:title" --limit 1 --json number -q '.[0].number')
    
    # Link to feature issue
    gh issue comment "$WS_ISSUE" --body "Part of feature #${FEATURE_ISSUE}"
    
    # Add GitHub issue number to WS file frontmatter
    sed -i "/^github_issue:/c\github_issue: ${WS_ISSUE}" "$WS_FILE"
    
    echo "✓ Created issue #${WS_ISSUE} for ${WS_ID}"
  done
  
  # Commit updated WS files with GitHub issue numbers
  git add tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md
  git commit -m "chore(${FEATURE_SLUG}): link GitHub issues to WS files

Feature issue: #${FEATURE_ISSUE}
Workstream issues: created and linked"
  
  git push origin feature/${FEATURE_SLUG}
  
  echo ""
  echo "✅ GitHub Project integration complete"
  echo "   Feature: #${FEATURE_ISSUE}"
  echo "   WS count: $(ls tools/hw_checker/docs/workstreams/backlog/WS-${FEATURE_ID}-*.md | wc -l)"
  echo "   View: https://github.com/your-org/your-repo/projects"
else
  echo "⚠️ GitHub CLI not available, skipping issue creation"
  echo "   Install: brew install gh (macOS) or see https://cli.github.com"
fi
```

### 8.6 Output включает Git info

В summary добавь:

```markdown
**Git (GitFlow):**
- Base: `develop`
- Branch: `feature/{slug}`
- Worktree: `../msu-ai-{slug}` ✅
- Commit: `docs({slug}): create WS specifications for F{XX}`
- Pushed: `origin/feature/{slug}`

**GitHub Integration:**
- Feature Issue: #{feature_issue}
- WS Issues: #{ws1}, #{ws2}, #{ws3}, ...
- Project Board: https://github.com/your-org/your-repo/projects
- All WS files linked to GitHub issues ✅

**Next:**
1. Review WS specs (optional)
2. `/oneshot F{XX}` (executes in feature branch)
3. After completion: PR `feature/{slug}` → `develop`
```

## Post-Design: GitHub Sync

After creating all WS files:

1. Run sync command:
   ```bash
   cd sdp
   poetry run sdp-github sync-all --ws-dir ../tools/hw_checker/docs/workstreams
   ```

2. Verify issues created:
   ```bash
   gh issue list --label "workstream" --json number,title
   ```

3. Verify project board:
   ```bash
   gh project item-list 2 --owner fall-out-bug --format json
   ```

===============================================================================
# 9. THINGS YOU MUST NEVER DO

❌ Ссылаться на WS без создания файла
❌ Оставлять scope > MEDIUM
❌ Использовать time estimates (дни/часы)
❌ Создавать -ANALYSIS.md файлы
❌ Пропускать Goal + AC
❌ Писать код (это задача /build)
❌ Игнорировать PROJECT_MAP.md
❌ Работать в main branch (только feature/* от develop)
❌ Забыть создать worktree (обязательно для изоляции)

===============================================================================
