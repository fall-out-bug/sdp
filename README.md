# SDP: Spec-Driven Protocol

**Структурированный протокол для AI-разработки.**

SDP превращает AI-ассистента в предсказуемый процесс: Discovery → Delivery → Evidence.

## Что это

SDP — набор prompts (skills), которые загружаются в Claude Code, Cursor или другой AI-инструмент. Skills определяют:

- **Discovery**: Как собирать требования и планировать features
- **Delivery**: Как писать код с TDD и quality gates
- **Evidence**: Как фиксировать решения в audit log

**Всё работает через skills.** CLI и Beads — опциональные дополнения.

## Быстрый старт

```bash
# Добавить SDP как submodule
git submodule add https://github.com/fall-out-bug/sdp.git sdp

# Добавить в .gitignore
echo "sdp/.git" >> .gitignore
```

В Claude Code: skills автоматически загрузятся из `sdp/.claude/skills/`.

## Основной workflow

### Discovery (планирование)

```
@vision "AI task manager"     → VISION.md, PRD.md, ROADMAP.md
@reality --quick              → Анализ кодовой базы
@feature "Add authentication" → Workstreams для feature
```

### Delivery (реализация)

```
@oneshot F001                 → Автономное выполнение всех workstreams
@review F001                  → Multi-agent quality review
@deploy F001                  → Merge в main
```

### Ручной режим

```
@build 00-001-01              → Один workstream с TDD
@build 00-001-02
@review F001
@deploy F001
```

### Debug

```
@debug "Test fails"           → Systematic debugging
@hotfix "API down"            → Emergency fix (P0)
@bugfix "Wrong totals"        → Quality fix (P1/P2)
```

## Skills

| Skill | Назначение | Фаза |
|-------|------------|------|
| `@vision` | Стратегическое планирование (7 агентов) | Discovery |
| `@reality` | Анализ кодовой базы (8 агентов) | Discovery |
| `@feature` | Планирование feature (@idea + @design) | Discovery |
| `@idea` | Сбор требований | Discovery |
| `@design` | Декомпозиция на workstreams | Discovery |
| `@oneshot` | Автономное выполнение | Delivery |
| `@build` | TDD для одного workstream | Delivery |
| `@review` | Quality review (6 агентов) | Delivery |
| `@deploy` | Deploy в main | Delivery |
| `@debug` | Systematic debugging | Debug |
| `@hotfix` | Emergency fix | Debug |
| `@bugfix` | Quality fix | Debug |

## Protocol Flow

```
@oneshot F001  →  @review F001  →  @deploy F001
     │                 │                │
     ▼                 ▼                ▼
Execute WS       APPROVED?         Merge PR
                    │
                    ├─ YES → proceed
                    └─ NO → fix loop
```

**Done = @review APPROVED + @deploy completed**, не просто "PR merged".

## Quality Gates

| Gate | Requirement |
|------|-------------|
| TDD | Tests first |
| Coverage | >= 80% |
| File size | < 200 LOC |
| Architecture | No layer violations |

## Workstream ID

Формат: `PP-FFF-SS`

- `PP` — Project (00 = SDP itself)
- `FFF` — Feature number
- `SS` — Step number

Пример: `00-024-03` = SDP, feature 24, step 3

## Структура проекта

```
your-project/
├── sdp/                      # SDP submodule
│   ├── prompts/skills/       # Skills (source of truth)
│   ├── prompts/agents/       # Agent definitions
│   ├── .claude/              # Claude Code integration
│   ├── docs/                 # Документация
│   └── CLAUDE.md             # Quick reference
└── docs/workstreams/         # Ваши workstreams
```

---

# Опциональные компоненты

## Go CLI (экспериментально)

CLI предоставляет helper-команды. **Не обязателен для работы протокола.**

```bash
# Установка
cd sdp/sdp-plugin && go build -o sdp ./cmd/sdp

# Команды
sdp doctor              # Health check
sdp status              # Project state
sdp guard activate WS   # Edit scope enforcement
sdp log show            # Evidence log
```

## Beads (экспериментально)

Task tracking для multi-session работы. **Не обязателен.**

```bash
brew tap beads-dev/tap && brew install beads
bd ready                # Найти доступные задачи
bd create --title="..." # Создать задачу
bd close <id>           # Закрыть задачу
```

## Evidence Layer (экспериментально)

Audit log в `.sdp/log/events.jsonl` с hash-chain.

```bash
sdp log show            # Показать события
sdp log trace           # Trace по commit/workstream
```

---

## Документация

| Файл | Содержимое |
|------|------------|
| [CLAUDE.md](CLAUDE.md) | Quick reference для Claude Code |
| [docs/PROTOCOL.md](docs/PROTOCOL.md) | Полная спецификация |
| [docs/vision/ROADMAP.md](docs/vision/ROADMAP.md) | Roadmap и milestones |
| [CHANGELOG.md](CHANGELOG.md) | История версий |

## License

MIT

---

**GitHub:** [fall-out-bug/sdp](https://github.com/fall-out-bug/sdp)
