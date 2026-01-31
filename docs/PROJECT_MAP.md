---
project_type: cli
prd_version: "2.0"
last_updated: 2026-01-26
diagrams_hash: sdp-prd-v1
---

# PROJECT_MAP: SDP

## 1. Назначение

Spec-Driven Protocol (SDP) - фреймворк для создания надежного, AI-ready кода через структурированные workstreams. Обеспечивает качество через TDD, Clean Architecture, и обязательный code review.

**Ключевые возможности:**
- Автоматическое создание workstream спецификаций
- TDD-цикл с обязательным покрытием ≥80%
- Clean Architecture разделение слоев
- Интеграция с Claude Code через skills
- PRD-документация с авто-генерируемыми диаграммами

## 2. Глоссарий

| Термин | Описание |
|--------|----------|
| Workstream (WS) | Атомарная единица работы, выполняемая за один проход |
| Feature | Группа связанных workstreams (5-30 WS) |
| PRD | Product Requirements Document - документация проекта |
| TDD | Test-Driven Development (Red-Green-Refactor) |
| Clean Architecture | Многослойная архитектура (Domain-Application-Infrastructure-Presentation) |
| AC | Acceptance Criteria - критерии завершения workstream |

## 3. Command Reference

### build

Выполняет один workstream через TDD-цикл.

```bash
sdp core build <ws-file>
```

**Options:**
  --skip-tests        Пропустить запуск тестов (не рекомендуется)
  --no-validation     Отключить проверку quality gates

### design

Создает дизайн для feature через интерактивный диалог.

```bash
sdp design <feature-id>
```

### oneshot

Автономное выполнение всех workstreams feature.

```bash
sdp oneshot <feature-id>
```

**Options:**
  --background        Выполнить в фоновом режиме
  --resume <id>      Продолжить после прерывания

### review

Code review с проверкой quality gates.

```bash
sdp review <feature-id>
```

### prd

Генерация и валидация PRD документации.

```bash
sdp prd validate <project-map>
sdp prd detect-type <project-path>
```

## 4. Configuration

### Config File

`~/.sdp/config.toml`:

```toml
# SDP configuration

[defaults]
project_id = "00"
max_loc = 200
min_coverage = 80

[quality]
enable_linters = true
enable_type_check = true
```

### Environment Variables

- `SDP_PROJECT_ID` - ID проекта (default: "00")
- `SDP_SKIP_HOOKS` - Отключить git hooks (1 = skip)

## 5. Usage Examples

### Basic Usage

```bash
# Создать новый workstream
sdp design F001

# Выполнить workstream
sdp build docs/workstreams/backlog/WS-001-01.md

# Автономное выполнение feature
sdp oneshot F001
```

### Interactive Design

```bash
sdp design F012
> Enter feature title: GitHub Agent Orchestrator
> Enter scope: MEDIUM
> [Dialog continues...]
```

## 6. Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Quality gate failure |
| 2 | Invalid workstream format |
| 3 | Dependency not satisfied |

## 7. Error Handling

### Error Messages

```
Error: Workstream validation failed

Cause: Acceptance criteria not met
Solution: Run tests with pytest --cov
```

### Recovery Strategy

При ошибке валидации:
1. Проверь вывод тестов (`pytest -v`)
2. Добавь недостающие тесты
3. Перезапусти `sdp build`
