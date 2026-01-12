# Spec Driven Protocol (SDP)

Протокол разработки на основе воркстримов для AI-агентов с one-shot выполнением.

[English version](README.md)

---

## Основная идея

**Workstream** = атомарная задача, которую AI выполняет за один проход, без итеративных циклов.

```
Feature → Workstreams → One-shot выполнение → Готово
```

## Терминология

| Термин | Scope | Размер | Пример |
|--------|-------|--------|--------|
| **Release** | Продуктовая веха | 10-30 Features | R1: MVP |
| **Feature** | Крупная фича | 5-30 Workstreams | F1: User Auth |
| **Workstream** | Атомарная задача | SMALL/MEDIUM/LARGE | WS-001: Domain entities |

**Метрики scope:**
- **SMALL**: < 500 LOC, < 1500 токенов
- **MEDIUM**: 500-1500 LOC, 1500-5000 токенов
- **LARGE**: > 1500 LOC → **разбить на 2+ WS**

**НЕ используем временные оценки.** Только LOC/токены.

## Workflow

Используйте slash-команды для упрощённого выполнения:

```bash
/idea "Аутентификация пользователей"  # 1. Сбор требований
/design idea-user-auth                # 2. Создание воркстримов
/build WS-001-01                      # 3. Реализация воркстрима
/review F01                           # 4. Проверка качества
/deploy F01                           # 5. Деплой в продакшн
```

## Быстрый старт

### 1. Собрать требования

```bash
/idea "Добавить аутентификацию по email/паролю"
```

**Результат:** `docs/drafts/idea-user-auth.md`

### 2. Создать воркстримы

```bash
/design idea-user-auth
```

**Результат:**
- `docs/workstreams/backlog/WS-001-01-domain.md`
- `docs/workstreams/backlog/WS-001-02-repository.md`
- `docs/workstreams/backlog/WS-001-03-service.md`
- `docs/workstreams/backlog/WS-001-04-api.md`
- `docs/workstreams/backlog/WS-001-05-tests.md`

### 3. Реализовать воркстримы

```bash
/build WS-001-01  # Domain layer
/build WS-001-02  # Repository
/build WS-001-03  # Service
# ... и т.д.
```

Или автономное выполнение:

```bash
/oneshot F01  # Выполнит все WS автоматически
```

### 4. Проверить качество

```bash
/review F01
```

Проверяет:
- ✅ Все критерии приёмки выполнены
- ✅ Покрытие ≥ 80%
- ✅ Нет TODO/FIXME
- ✅ Clean Architecture соблюдена

### 5. Задеплоить

```bash
/deploy F01
```

Генерирует:
- Docker конфиги
- CI/CD пайплайны
- Release notes
- План деплоя

## Справка по командам

| Команда | Назначение | Когда использовать |
|---------|------------|-------------------|
| `/idea` | Сбор требований | Начало новой фичи |
| `/design` | Создание воркстримов | После прояснения требований |
| `/build` | Реализация воркстрима | Выполнение одного WS |
| `/review` | Проверка качества | После завершения всех WS |
| `/deploy` | Деплой в продакшн | После APPROVED review |
| `/issue` | Отладка и маршрутизация | Анализ багов |
| `/hotfix` | Экстренное исправление | P0 проблемы в продакшне |
| `/bugfix` | Качественное исправление | P1/P2 баги |
| `/oneshot` | Автономное выполнение | Выполнить все WS без участия |

## Quality Gates

| Gate | Требования |
|------|------------|
| **AI-Readiness** | Файлы < 200 LOC, CC < 10, type hints |
| **Clean Architecture** | Нет нарушений слоёв |
| **Error Handling** | Нет `except: pass` |
| **Test Coverage** | ≥ 80% |
| **No TODOs** | Все выполнены или новый WS |

## Базовые принципы

| Принцип | Суть |
|---------|------|
| **SOLID** | SRP, OCP, LSP, ISP, DIP |
| **DRY** | Don't Repeat Yourself |
| **KISS** | Keep It Simple |
| **YAGNI** | Строй только нужное |
| **TDD** | Сначала тесты (Red → Green → Refactor) |
| **Clean Code** | Читаемый, поддерживаемый |
| **Clean Architecture** | Зависимости направлены внутрь |

Подробнее: [docs/PRINCIPLES.md](docs/PRINCIPLES.md)

## Структура файлов

```
sdp/
├── PROTOCOL.md              # Полная спецификация
├── CODE_PATTERNS.md         # Паттерны реализации
├── RULES_COMMON.md          # Общие правила
├── docs/
│   ├── PRINCIPLES.md        # SOLID, DRY, KISS, YAGNI
│   └── concepts/            # Clean Architecture, Artifacts, Roles
├── prompts/
│   └── commands/            # Slash-команды (/idea, /design, и т.д.)
├── schema/                  # JSON валидация
├── scripts/                 # Утилиты
├── hooks/                   # Git hooks
└── templates/               # Шаблоны документов
```

## Ресурсы

| Ресурс | Назначение |
|--------|------------|
| [PROTOCOL.md](PROTOCOL.md) | Полная спецификация |
| [docs/PRINCIPLES.md](docs/PRINCIPLES.md) | SOLID, DRY, KISS, YAGNI |
| [docs/concepts/](docs/concepts/) | Архитектурные концепции |
| [CODE_PATTERNS.md](CODE_PATTERNS.md) | Паттерны кода |
| [MODELS.md](MODELS.md) | Рекомендации по моделям |
| [CLAUDE.md](CLAUDE.md) | Интеграция с Claude Code |

## Интеграция

### Интерактивная установка (рекомендуется)

Запустите мастер установки:

```bash
python scripts/init.py
```

Мастер:
- ✅ Определит вашу IDE (Cursor/Claude Code)
- ✅ Создаст необходимые директории
- ✅ Скопирует файлы SDP в проект
- ✅ Установит Git hooks для валидации
- ✅ Проведёт по настройкам проекта
- ✅ Покажет следующие шаги

### Ручная установка

#### Для Claude Code

1. Скопировать файлы в проект:
```bash
cp -r prompts/ your-project/
cp -r schema/ your-project/
cp -r .claude/ your-project/
cp CLAUDE.md your-project/
```

2. Заполнить `PROJECT_CONVENTIONS.md` правилами проекта

3. Использовать skills: `@idea`, `@design`, `@build`, и т.д.

#### Для Cursor

1. Скопировать файлы в проект:
```bash
cp -r prompts/ your-project/
cp -r schema/ your-project/
cp -r .cursor/ your-project/
cp .cursorrules your-project/
```

2. Заполнить `PROJECT_CONVENTIONS.md` правилами проекта

3. Использовать slash-команды: `/idea`, `/design`, `/build`, и т.д.

### После интеграции

1. **Заполнить PROJECT_CONVENTIONS.md** — Добавить специфичные для проекта правила DO/DON'T
2. **Проверить конфиг IDE** — Посмотреть `.cursorrules` или `CLAUDE.md`
3. **Установить Git hooks** — Запустить `scripts/init.py --install-hooks`
4. **Прочитать документацию** — См. [PROTOCOL.md](PROTOCOL.md) и [docs/PRINCIPLES.md](docs/PRINCIPLES.md)

---

**Версия:** 0.3.0 | **Статус:** Активен
