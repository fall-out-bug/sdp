---
description: Emergency P0 fixes. Fast-track production deployment with minimal changes. Branch from main, deploy < 2h.
agent: builder
---

# /hotfix — Emergency Production Fixes

При вызове `/hotfix "description" --issue-id=001`:

1. **Создай ветку** — `git checkout -b hotfix/{id}-{slug}` от main
2. **Минимальный фикс** — Без рефакторинга, только баг
3. **Быстрое тестирование** — Smoke + critical path
4. **Коммит** — `fix(scope): description (issue NNN)`
5. **MERGE, TAG, PUSH** — Выполни сам!
6. **Backport** — Мерж в dev и feature ветки
7. **Закрой issue** — Обнови статус в файле

## КРИТИЧНО: Ты ДОЛЖЕН завершить

```bash
# Мерж в main и тег
git checkout main
git merge hotfix/{branch} --no-edit
git tag -a v{VERSION} -m "Hotfix: {description}"
git push origin main --tags

# Backport в dev
git checkout dev
git merge main --no-edit
git push origin dev
```

**Работа НЕ завершена пока все `git push` не выполнены.**

## Quick Reference

**Input:** P0 CRITICAL issue  
**Output:** Production fix + запушено в origin

**Key Rules:**
- Минимальные изменения
- Без рефакторинга
- Без новых фич
- Быстрое тестирование
- Backport обязателен
