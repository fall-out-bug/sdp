# Prompts

4 фазы — это всё что нужно.

## Использование

```
@sdp/prompts/structured/phase-1-analyze.md   # Сформировать карту WS
@sdp/prompts/structured/phase-2-design.md    # Спланировать один WS
@sdp/prompts/structured/phase-3-implement.md # Выполнить WS (Auto модель)
@sdp/prompts/structured/phase-4-review.md    # Проверить результат
```

## Когда какой

| Задача | Фаза | Модель |
|--------|------|--------|
| Понять что делать из спек | 1 | Sonnet |
| Детализировать WS для исполнения | 2 | Sonnet |
| Выполнить план | 3 | Auto/Haiku |
| Проверить качество | 4 | Sonnet |

## Guardrails и Quality Gates

См. `@sdp/PROTOCOL.md` — там всё в одном месте.
