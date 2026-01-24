---
ws_id: 00-195-06
project_id: 00
feature: F011
status: backlog
size: SMALL
github_issue: 1038
assignee: null
started: null
completed: null
blocked_reason: null
---

## 02-195-06: hw-checker PRD Migration

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- hw-checker PROJECT_MAP.md обновлён до PRD v2.0 формата
- Добавлены все 7 секций для service profile
- Диаграммы сгенерированы из существующего кода
- diagrams_hash установлен в frontmatter

**Acceptance Criteria:**
- [ ] AC1: PROJECT_MAP.md содержит frontmatter с project_type: service, prd_version: "2.0"
- [ ] AC2: Все 7 секций service profile заполнены
- [ ] AC3: Sequence flow "submission-processing" задокументирован
- [ ] AC4: Модель БД содержит ключевые таблицы (submissions, runs, results)
- [ ] AC5: `/codereview` проходит PRD check (hash match)

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

Финальный WS фичи — применение всех инструментов к реальному проекту. Это валидация что F195 работает end-to-end.

### Зависимость

00--01..05 (все предыдущие WS фичи)

### Входные файлы

- `tools/hw_checker/docs/PROJECT_MAP.md` — текущий файл
- `tools/hw_checker/src/` — исходный код для аннотаций
- `sdp/src/sdp/prd/` — все модули F195

### Шаги

1. Бэкап текущего PROJECT_MAP.md
2. Запустить `/prd hw-checker` для scaffold через диалог
3. Заполнить 7 секций контентом
4. Добавить @prd аннотации в ключевые файлы (use_case.py, saga_orchestrator.py)
5. Сгенерировать диаграммы через `/prd hw-checker --update`
6. Проверить `/codereview` проходит

### Код

```markdown
# Пример обновлённого PROJECT_MAP.md (секция 1)

---
project_type: service
prd_version: "2.0"
last_updated: 2026-01-22
diagrams_hash: abc123def456
---

# PROJECT_MAP: hw-checker

## 1. Назначение

Автоматизированная система проверки домашних заданий для курсов ML System Design и Big Data.
Принимает submissions через API или CLI, запускает в изолированных Docker-контейнерах,
оценивает по rubric, публикует результаты в Google Sheets и отправляет уведомления.

## 2. Глоссарий

| Термин | Описание |
|--------|----------|
| Submission | Отправка домашней работы студентом (git URL или файл) |
| Run | Одна попытка выполнения submission в sandbox |
| SAGA | Распределённая транзакция с компенсациями для staged execution |
| DinD | Docker-in-Docker: изоляция sandbox от host |

## 3. Внешний API

### POST /api/v1/submissions

Создание новой submission для проверки.

```json
{"student_id": "ivan_petrov", "assignment_id": "hw1", "git_url": "https://..."}
```

Response: `202 Accepted` с `submission_id`

...
```

```python
# Пример аннотаций в use_case.py

from sdp.prd import prd_flow, prd_step

@prd_flow("submission-processing")
@prd_step(1, "Получение submission из очереди")
async def process_submission(self, job: Job) -> RunResult:
    """Process single submission through SAGA orchestrator."""
    ...

@prd_flow("submission-processing")
@prd_step(2, "Клонирование репозитория")
async def clone_repository(self, url: str) -> Path:
    ...

@prd_flow("submission-processing")
@prd_step(3, "Запуск в Docker sandbox")
async def run_in_sandbox(self, path: Path) -> ExecutionResult:
    ...
```

### Ожидаемый результат

```
tools/hw_checker/docs/
├── PROJECT_MAP.md              # PRD v2.0 format
└── diagrams/
    ├── sequence-submission-processing.mmd
    ├── sequence-submission-processing.puml
    ├── component-overview.mmd
    └── deployment-production.puml

tools/hw_checker/src/hw_checker/application/
├── run_homework/
│   └── use_case.py            # @prd annotations added
└── saga_orchestrator.py       # @prd annotations added
```

### Scope Estimate

- Файлов: ~3 изменено + 4 диаграммы
- Строк кода: ~300 (PROJECT_MAP: 200, аннотации: 50, диаграммы: 50)
- Токенов: ~1500

**Оценка размера:** SMALL

### Критерий завершения

```bash
# PRD validation
cd sdp && poetry run sdp-prd validate ../tools/hw_checker/docs/PROJECT_MAP.md

# Hash match
./sdp/hooks/post-codereview.sh F195
# Должен вывести: ✓ PRD diagrams up-to-date

# Diagrams exist
ls tools/hw_checker/docs/diagrams/*.mmd
ls tools/hw_checker/docs/diagrams/*.puml
```

### Ограничения

- НЕ менять существующую логику кода (только добавить аннотации)
- НЕ переписывать весь PROJECT_MAP.md (инкрементальное обновление)
- НЕ добавлять аннотации во все файлы (только ключевые flows)

---

### Human Verification (UAT)

#### 🚀 Quick Smoke Test (30 секунд)

```bash
# Проверить что PROJECT_MAP.md валиден
cd sdp && poetry run sdp-prd validate ../tools/hw_checker/docs/PROJECT_MAP.md
# Ожидаемый результат: ✅ Valid PRD v2.0
```

#### 📋 Manual Test Scenarios

| # | Сценарий | Шаги | Ожидание | ✅/❌ |
|---|----------|------|----------|------|
| 1 | PRD sections | Открыть PROJECT_MAP.md | 7 секций заполнены |  |
| 2 | Diagrams | ls docs/diagrams/ | 4 файла (.mmd, .puml) |  |
| 3 | Code review | /codereview F195 | PRD check passes |  |
