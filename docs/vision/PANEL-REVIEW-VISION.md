# Panel Review: Vision + Manifesto

> **Панель:** PG, Collison (Stripe), Karpathy, Masad (Replit), Hashimoto (Terraform), Charity Majors (Honeycomb)
> **Дата:** 2026-02-08

## Ключевые рекомендации

| Эксперт | Рекомендация |
|---------|-------------|
| **Hashimoto** | Выбери ОДИН продукт (агент-фреймворк ИЛИ CLI) и шипни. GitHub Action — точка входа |
| **Karpathy** | Перепиши аргумент декомпозиции: это permanent из-за ЛЮДЕЙ, не из-за context windows |
| **Charity** | Лиди с forensics, не verification. Добавь OpenTelemetry. Никто не строит forensic chain для AI-кода |
| **Collison** | Сделай compliance brief: SOC2, GDPR, DORA. Compliance buyer требует конкретный чеклист |
| **Masad** | Убей дуальность ship/drive. Один режим, один флаг. Или хотя бы `sdp auto` / `sdp step` |
| **PG** | Лиди с УЖАСОМ ответственности, а не с механизмом. "Какие у вас доказательства в суде?" |

## Консенсус

1. **Accountability thesis — сильнейший инсайт.** Но он закопан. Должен быть ПЕРВЫМ, что видит пользователь
2. **"Evidence not hope" — маркетинг или инженерия?** Нужно уточнить: "verified build" в заголовке oversells. Это heuristic, не proof
3. **Decomposition permanent — ДА, но по другой причине.** Не context windows (temporary), а human cognition (permanent). Переписать аргумент
4. **Forensics > Verification.** Verification — crowded claim. Forensic chain для AI-кода — уникальное предложение
5. **Два продукта в одном.** Codebase = agent orchestration framework. Manifesto = CLI tool. Выбрать одно
6. **IP attribution отсутствует.** Enterprise спрашивает "AI скопировал GPL код?" раньше, чем "код корректный?"

## Aha-момент

**Charity:** "Код более интересен, чем манифест. Это лучшая возможная проблема."

**Karpathy:** Провенанс-цепочка (какой модель, какой промпт, какая спека) = ценность и для верификации, и для IP-аудита. Инфраструктура уже есть, позиционирование нет.
