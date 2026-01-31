---
description: Quality bug fixes (P1/P2). Full TDD cycle, branch from dev, no production deploy.
agent: builder
---

# /bugfix — Quality Bug Fixes

При вызове `/bugfix issue NNN`:

1. **Прочитай issue** — Загрузи `docs/issues/{NNN}-*.md`
2. **Создай ветку** — `git checkout -b bugfix/{NNN}-{slug}` от dev
3. **TDD цикл** — Напиши падающий тест → реализуй фикс → рефактор
4. **Quality gates** — pytest, coverage ≥80%, mypy --strict, ruff
5. **Коммит** — `fix(scope): description (issue NNN)`
6. **Закрой issue** — Обнови статус в файле issue
7. **MERGE И PUSH** — Выполни сам, не давай инструкции!

## КРИТИЧНО: Ты ДОЛЖЕН завершить

```bash
git checkout dev
git merge bugfix/{branch} --no-edit
git push
git status  # ДОЛЖЕН показать "up to date with origin"
```

**Работа НЕ завершена пока `git push` не выполнен.**

## Quick Reference

**Input:** P1/P2 issue  
**Output:** Баг исправлен + тесты + запушено в origin

| Aspect | Hotfix | Bugfix |
|--------|--------|--------|
| Severity | P0 | P1/P2 |
| Branch from | main | dev |
| Testing | Fast | Full |
