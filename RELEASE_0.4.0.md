# SDP v0.4.0 — PRD Command + Two-Stage Review 🚀

**TL;DR:** Выкатили SDP v0.4.0 — two-stage review, multi-IDE parity, extension system, PRD command с авто-графиками. 48 workstreams в багаже, 83% готовности.

---

## Что нового в v0.4.0

### 📝 F011: PRD Command — Ваша PRD, автогенерируемые диаграммы

**Проблема:** PRD documentation — боль. Написание, обновление, синхронизация с кодом — всё ручками, легко устаревает.

**Решение:** `/prd` команда делает всю грязную работу:
- **Авто-диаграммы** — генерирует architecture diagrams из annotations в коде
- **Annotations** — парсит `@prd` комментарии, обновляет документацию
- **Валидация** — проверяет line limits, ссылочную целостность
- **Профили** — разные форматы для разных проектов

**6 workstreams:**
- PRD Command + Profiles — базовая команда
- Line Limits Validator — проверяет длину строк в документах
- Annotation Parser — парсит `@prd:` комментарии из кода
- Diagram Generator — строит Mermaid diagrams из кода
- CodeReview Hook Integration — интеграция с review process
- SDP PRD Migration — миграция существующей PRD

**Результат:** Пишешь код с `@prd:` аннотациями → SDP сам обновляет PRD с новыми диаграммами. Магия.

---

### 🔍 F003: Two-Stage Code Review — "не полируй неправильное"

**Проблема:** Традиционный review пропускает "well-written but wrong" баги. Код чистый, тесты есть, но не соответствует spec.

**Решение:** Двухэтапный review:
1. **Stage 1 (Spec Compliance):** Соответствует ли код спецификации?
   - Goal Achievement (все AC выполнены?)
   - Specification Alignment (все фичи на месте?)
   - AC Coverage (есть тесты на всё?)
   - No Over/Under-Engineering

2. **Stage 2 (Code Quality):** Код качественный?
   - Coverage ≥80%, mypy strict
   - AI-Readiness (файлы <200 LOC)
   - Clean Architecture, Security, No Tech Debt

**Ключевой инсайт:** Stage 2 запускается ТОЛЬКО если Stage 1 прошёл. Не трать время на polishing incorrect code.

**5 workstreams** → `sdp/prompts/skills/two-stage-review.md`, `/codereview` skill обновлён

---

### 🔌 F004: Platform Adapters — Claude, Cursor, OpenCode

**Проблема:** У каждого AI-IDE свой формат настроек, skills, hooks. Дублируем логику для каждой платформы.

**Решение:** Единый адаптер для всех платформ:
- `PlatformAdapter` interface — единый API
- `detect_platform()` — автоопределение IDE (ищет `.claude/`, `.codex/`, `.opencode/`)
- Общие операции: install skills, configure hooks, load settings

**4 workstreams:**
- Interface definition + base implementation
- Claude Code adapter (`.claude/` support)
- Codex adapter (`.codex/` support)
- OpenCode adapter (`.opencode/` support)

**Результат:** SDP работает везде одинаково. Switch IDE — не теряй навыки.

---

### 🧩 F005: Extension System — кастомизация без форка

**Проблема:** Хочется project-specific настройки (hooks, patterns, skills), но forking core — боль.

**Решение:** Extension system:
- `sdp.local/` или `~/.sdp/extensions/{name}/` — папка расширения
- `extension.yaml` — манифест (name, version, author)
- `hooks/`, `patterns/`, `skills/`, `integrations/` — компоненты
- `ExtensionLoader` — автообнаружение и загрузка

**3 workstreams:**
- Extension interface + Protocol-based design
- Manifest parser + validator
- Extension loader с двумя search paths

**Результат:** Добавляешь свои hooks/patterns без изменения core SDP.任何人 может contribute extensions.

---

### 🏗️ F006: Core SDP — фреймворк

**6 workstreams:** базовые компоненты SDP
- Workstream parser — парсит YAML frontmatter из markdown
- Feature decomposition — разбивает фичи на workstreams
- Project map parser — читает `PROJECT_MAP.md`
- Pip package — `pip install sdp`
- File size reduction — держит файлы <200 LOC
- Integration tests — тесты интеграции

---

### ⚡ F007: Oneshot & Hooks — автономное выполнение

**10 workstreams:** автономное исполнение фич + git hooks
- Oneshot validation — проверяет перед запуском
- Git hooks (pre-commit, post-commit, pre-push) — quality gates
- Cursor agents integration
- Debug command implementation
- Test command implementation
- Documentation cleanup
- `/idea` и `/design` skills
- EP30 misclassification fix
- Debug title fix

**Результат:** `/oneshot F060` — executes all workstreams for feature F60 autonomously. Ты пьёшь кофе — AI делает свою работу.

---

### 📏 F008: Contract-Driven WS Tiers — уровни сложности

**9 workstreams:** система tiers для workstreams
- Contract-driven WS spec — yaml schema для WS
- Capability tier validator — проверяет tier
- Model mapping registry — регистр моделей
- Test command workflow — workflow для тестов
- Model agnostic builder router — маршрутизация builder
- Model selection optimization — оптимизация выбора модели
- Tier auto-promotion — авто-повышение tier
- Escalation metrics — метрики для эскалации
- Runtime contract validation — валидация в runtime

**Результат:** Starter → Standard → Advanced tiers. Новички видят только Starter, эксперты — все уровни.

---

### 🛠️ F010: SDP Infrastructure — инфраструктура

**5 workstreams:** базовая инфраструктура
- Sync SDP content — синхронизация контента
- PP-FFF-SS naming migration — миграция имен
- Update SDP documentation
- Configure SDP as submodule
- Add SDP submodule

**Результат:** SDP как submodule в проектах. `git submodule update` — и всё актуально.

---

## Что внутри

| Feature | Workstreams | Статус | Описание |
|---------|-------------|--------|-------------|
| F003: Two-Stage Review | 5 | ✅ | Spec → Quality, не полирим неправильное |
| F004: Platform Adapters | 4 | ✅ | Claude Code + Cursor + OpenCode |
| F005: Extension System | 3 | ✅ | Кастомизация без fork |
| F006: Core SDP | 6 | ✅ | Базовый фреймворк |
| F007: Oneshot & Hooks | 10 | ✅ | Автономное выполнение |
| F008: WS Tiers | 9 | ✅ | Уровни сложности |
| F010: Infrastructure | 5 | ✅ | Инфраструктура |
| **F011: PRD Command** | **6** | **✅** | **PRD + диаграммы** |

**Total:** 48/58 workstreams completed (83%)

---

## Quick Start

```bash
git clone https://github.com/fall-out-bug/sdp.git
cd sdp
poetry install

# Verify
sdp --version  # v0.4.0
```

---

## Контрибуция

Open source, Pull Requests welcome!
GitHub: https://github.com/fall-out-bug/sdp

---

**v0.4.0 — available now.**
