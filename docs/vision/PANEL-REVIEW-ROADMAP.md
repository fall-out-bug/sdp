# Panel Review: Roadmap + Features

> **Панель:** PG, Collison (Stripe), Karpathy, Masad (Replit), Hashimoto (Terraform), Charity Majors (Honeycomb)
> **Дата:** 2026-02-08

## Конкретные изменения (TOP 3 от каждого)

### PG
1. **Убрать `sdp drive` из Phase 1.** Шипнуть в Phase 2 после того, как `sdp ship` доказал себя
2. **Добавить `sdp plan` как first-class команду.** Украсть plan/apply у Terraform. Это trust-building вместо drive mode
3. **Переделать evidence strategy.** Убить N=10. Инструментировать каждый ран и публиковать "AI Code Quality Benchmark" квартально

### Collison (Stripe)
1. **Compliance export + verification certificates → P1.** Enterprise нужно подписать. P2 слишком поздно
2. **Billing infrastructure в Phase 1.** Нельзя брать деньги без metering, invoicing, payment processing
3. **Usage-based pricing с первого дня.** Не "free OSS / paid private". Free tier 50 runs/month, paid выше

### Karpathy
1. **Model provenance tracking → P0.** Каждый код: model name, version, prompt hash, temperature, timestamp. Non-negotiable для enterprise
2. **Model selection policy → P1.** Не просто multi-provider — policy engine для routing по risk level
3. **Moat = verification + evidence, не decomposition.** Decomposition будет commodity. Инвестировать в evidence

### Masad (Replit)
1. **JSON API рядом с CLI в Phase 1.** `sdp ship --output=json`. Replit мог бы интегрировать не дожидаясь Phase 3
2. **Team/collaboration features → P1.** Shared templates, conflict detection. Single-player не продаётся в enterprise
3. **Убрать Phase 4 "standard".** Стандарты не создают намеренно. Создаются через 10K repos

### Hashimoto (Terraform)
1. **`sdp plan` / `sdp apply` как primary UX.** `ship` = `plan --auto-apply`. `drive` = `plan --interactive`. Одна ментальная модель
2. **Usage-based pricing, не OSS/private split.** "Free then restrict" → community riot. Charge от первого дня
3. **Open-source verification engine, proprietary orchestration.** Правильный open-core split

### Charity Majors (Honeycomb)
1. **Production rollback / feature flags.** Feature map заканчивается на "merge PR". Production дальше. `sdp rollback` или хотя бы feature flag integration
2. **Определить ЧТО В audit trail — сейчас.** Model provenance, prompt hash, verification output, approval chain, tamper-evident hashing. Phase 1 или Phase 2 compliance export экспортирует мусор
3. **Добавить `sdp incident`.** Prod сломался → "это AI-сгенерировано?" + "какая верификация прошла?" Forensic trace от git commit назад через всю SDP цепочку. Enterprise killer feature

## Критические проблемы с текущим роадмапом

### Phase 1 перегружен
- **3 продукта за 6 недель** (ship + drive + GitHub Action) — нереалистично
- **`sdp ship` и `sdp drive` не существуют в коде.** 22 CLI-команды, ни одна не является целевым продуктом
- **Решение:** убрать drive из Phase 1. Добавить `sdp plan` вместо него

### Приоритеты перепутаны
- Compliance export P2 → должен быть **P1** (enterprise нужно подписать)
- Model provenance отсутствует → должен быть **P0**
- Team features P3 → должны быть **P1** (enterprise = teams)
- Billing infrastructure отсутствует → нужен в **Phase 1**

### Evidence strategy слаба
- N=10 — не исследование, а анecdote
- **Решение:** инструментировать каждый ран, собирать observational data, публиковать quarterly benchmark
- **Kill criteria уточнить:** "500 ранов, catch rate <5%, post-merge defect rate = baseline → убить продукт"

### Revenue targets
- $50K ARR (month 4) — реалистично при 2 enterprise контрактах
- $500K ARR (month 8) — stretch, но возможно с GitHub Action viralность
- $2M ARR (month 12) — aspirational. Нет sales motion в роадмапе. Нет salesperson

### Отсутствующие features
| Feature | Кто предложил | Priority |
|---------|--------------|----------|
| `sdp plan` (standalone) | Hashimoto, PG | P0 |
| `sdp incident` (forensic trace) | Charity | P1 |
| `sdp rollback` (production) | Charity | P1 |
| Model provenance tracking | Karpathy, Charity | P0 |
| JSON API (`--output=json`) | Masad | P0 |
| Billing/metering infrastructure | Collison | P1 |
| Team collaboration | Masad | P1 |
| Model selection policy | Karpathy | P1 |
| Failure mode UX | PG | P1 |

## Hashimoto's open-core формула

> Open-source verification engine. Proprietary orchestration + evidence platform.
> Community строит интеграции вокруг verification. Ты продаёшь orchestration + evidence сверху.

## Collison's kill criteria

> "Если после 500 SDP ранов catch rate < 5% и post-merge defect rate неотличим от baseline — убить продукт. Конкретно. Тестируемо. Честно."
