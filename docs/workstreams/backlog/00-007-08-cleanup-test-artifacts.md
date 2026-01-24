---
id: WS-201-08
title: Cleanup F201 test artifacts and update INDEX
feature: F007
status: backlog
size: SMALL
github_issue: TBD
dependencies:
  - WS-201-01 # Validate /oneshot in Cursor and OpenCode (UAT completed)
---

## 02-201-08: Cleanup F201 test artifacts and update INDEX

### 🎯 Цель (Goal)

**Что должно РАБОТАТЬ после завершения WS:**
- Тестовые WS (02-201-TEST-01/02/03) перемещены в `completed/2026-01/`
- `.oneshot/` directory очищен от тестовых checkpoint и progress файлов
- INDEX.md обновлен (удалена F201-TEST секция)
- test-oneshot-validation.md обновлен с финальными результатами UAT

**Acceptance Criteria:**
- [ ] 02-201-TEST-01 перемещен в `completed/2026-01/`
- [ ] 02-201-TEST-01 перемещен в `completed/2026-01/`
- [ ] 02-201-TEST-01 перемещен в `completed/2026-01/`
- [ ] `.oneshot/F201-TEST-checkpoint.json` удален (если существует)
- [ ] `.oneshot/F201-TEST-progress.json` удален (если существует)
- [ ] INDEX.md обновлен (удалена F201-TEST секция)
- [ ] `tools/hw_checker/docs/test-oneshot-validation.md` обновлен с UAT результатами

**⚠️ WS НЕ завершён, пока Goal не достигнута (все AC ✅).**

---

### Контекст

WS-201-01 успешно прошел UAT и все тестовые артефакты можно удалить:

**Созданные тестовые артефакты:**
- 3 тестовых WS (02-201-TEST-01/02/03) в backlog/
- Checkpoint файл `.oneshot/F201-TEST-checkpoint.json`
- Progress JSON `.oneshot/F201-TEST-progress.json`
- F201-TEST секция в INDEX.md

**Проблема:**
- Тестовые WS остаются в backlog/ (не соответствует статусу "completed")
- `.oneshot/` содержит тестовые файлы
- INDEX.md содержит секцию для временной тестовой фичи

**Решение:**
- Переместить тестовые WS в completed/2026-01/
- Очистить `.oneshot/` directory
- Обновить INDEX.md
- Задокументировать финальные результаты UAT

---

### Зависимость

WS-201-01 (UAT completed successfully)

---

### Входные файлы

- `tools/hw_checker/docs/workstreams/backlog/02-201-TEST-01.md`
- `tools/hw_checker/docs/workstreams/backlog/02-201-TEST-01.md`
- `tools/hw_checker/docs/workstreams/backlog/02-201-TEST-01.md`
- `.oneshot/F201-TEST-checkpoint.json` (если существует)
- `.oneshot/F201-TEST-progress.json` (если существует)
- `tools/hw_checker/docs/test-oneshot-validation.md`

---

### Шаги

1. **Переместить тестовые WS:**
   - Создать directory `tools/hw_checker/docs/workstreams/completed/2026-01/` (если не существует)
   - Переместить `02-201-TEST-01.md` → `completed/2026-01/`
   - Переместить `02-201-TEST-01.md` → `completed/2026-01/`
   - Переместить `02-201-TEST-01.md` → `completed/2026-01/`

2. **Очистить `.oneshot/` directory:**
   - Удалить `.oneshot/F201-TEST-checkpoint.json` (если существует)
   - Удалить `.oneshot/F201-TEST-progress.json` (если существует)
   - Оставить `.oneshot/` directory пустым (или удалить если нет других файлов)

3. **Обновить INDEX.md:**
   - Найти секцию "### P0: F201-TEST: /oneshot Validation Test"
   - Удалить эту секцию (включая таблицу с 02-201-TEST-01/02/03)
   - Убедиться что INDEX.md валиден (markdown)

4. **Обновить test-oneshot-validation.md:**
   - Добавить секцию "## Final UAT Results (2026-01-23)"
   - Задокументировать успешные результаты UAT
   - Указать что обе IDE (Cursor и OpenCode) прошли тест
   - Добавить информацию о cleanup (дата, какие файлы удалены)

5. **Верификация:**
   - Проверить что тестовые WS больше нет в backlog/
   - Проверить что `.oneshot/` очищен
   - Проверить что INDEX.md не содержит F201-TEST
   - Проверить что test-oneshot-validation.md обновлен

---

### Ожидаемый результат

- `completed/2026-01/02-201-TEST-01.md`
- `completed/2026-01/02-201-TEST-01.md`
- `completed/2026-01/02-201-TEST-01.md`
- `.oneshot/` directory: очищен
- INDEX.md: обновлен (нет F201-TEST секции)
- test-oneshot-validation.md: обновлен с UAT результатами

### Scope Estimate

- Файлов: 3 перемещено + 2 изменено
- Строк: ~50 (SMALL)
- Токенов: ~150

---

### Критерий завершения

```bash
# Тестовые WS перемещены
ls -la tools/hw_checker/docs/workstreams/completed/2026-01/WS-F201-TEST-*.md

# Тестовые WS удалены из backlog
! ls -la tools/hw_checker/docs/workstreams/backlog/WS-F201-TEST-*.md 2>/dev/null

# .oneshot/ очищен
! test -f .oneshot/F201-TEST-checkpoint.json
! test -f .oneshot/F201-TEST-progress.json

# INDEX.md обновлен
! grep -q "F201-TEST" tools/hw_checker/docs/workstreams/INDEX.md

# test-oneshot-validation.md обновлен
grep -q "Final UAT Results" tools/hw_checker/docs/test-oneshot-validation.md
```

---

### Ограничения

- НЕ трогать: реальные workstreams (WS-201-01..06)
- НЕ удалять: master-промпты
- НЕ менять: команду `/oneshot` (только cleanup тестов)

---

### Шаблон для обновления test-oneshot-validation.md

```markdown
## Final UAT Results (2026-01-23)

### Test Summary

**Feature:** F201-TEST (/oneshot Validation)
**Test Date:** 2026-01-23
**Test Duration:** ~10 min
**Test Environments:** Cursor IDE, OpenCode IDE

### Test Execution

**Cursor IDE:**
- Command: `/oneshot F201-TEST`
- Status: ✅ PASSED
- Workstreams executed: 3/3 (02-201-TEST-01, 02-201-TEST-01, 02-201-TEST-01)
- Checkpoint: ✅ Created
- Progress tracking: ✅ Working
- Error handling: ✅ Tested (CRITICAL/HIGH/MEDIUM)

**OpenCode IDE:**
- Command: `/oneshot-simple F201-TEST`
- Status: ✅ PASSED
- Workstreams executed: 3/3
- Checkpoint: ✅ Created
- Progress tracking: ✅ Working
- Error handling: ✅ Tested

### Verification

| Check | Cursor | OpenCode |
|-------|---------|----------|
| All WS executed | ✅ | ✅ |
| Checkpoint created | ✅ | ✅ |
| Progress updated | ✅ | ✅ |
| Error handling | ✅ | ✅ |
| PR approval gate | ✅ (when gh available) | ✅ (when gh available) |

### Cleanup (2026-01-23)

**Completed:**
- Test workstreams moved to `completed/2026-01/`
- `.oneshot/F201-TEST-checkpoint.json` deleted
- `.oneshot/F201-TEST-progress.json` deleted
- INDEX.md F201-TEST section removed

### Conclusion

**Overall Status:** ✅ PRODUCTION READY

The `/oneshot` command has been successfully validated in both Cursor and OpenCode IDEs. All acceptance criteria met. Feature is ready for production use.

**Next Steps:**
- Feature F201 code review complete
- All workstreams (WS-201-01..06) APPROVED
- Proceed to deployment: `/deploy F201`
```
