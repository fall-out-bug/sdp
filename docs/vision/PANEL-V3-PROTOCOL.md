# Panel v3: Protocol + Execution

> **Панель:** Hashimoto (Terraform), Masad (Replit), Collison (Stripe), Dario (Anthropic)
> **Дата:** 2026-02-08
> **Фокус:** protocol design, adoption, enterprise readiness, evidence format

---

## Вердикты: "Адоптировали бы?"

| Эксперт | Вердикт | Условие |
|---------|---------|---------|
| **Hashimoto** | Conditional YES | Schema-first protocol + independent implementation |
| **Masad** | YES на протокол, NO на методологию | Clean separation protocol/methodology |
| **Collison** | NOT YET | Compliance design doc, data residency, RBAC, per-seat pricing |
| **Dario** | Interested | Multi-vendor standard, would participate in working group |

---

## #1 от каждого: что делать в следующие 30 дней

### Единогласно: SHIP THE SCHEMA

> **Hashimoto:** "Ship the schema."
> **Masad:** "Ship the schema."
> **Dario:** "Ship the schema."
> **Collison:** "Write the compliance design doc."

**3 к 1. Приоритет: JSON Schema для четырёх примитивов (plan, apply, evidence, incident). Machine-readable. Versioned. Public repo.**

---

## Конкретные изменения в роадмапе

### Hashimoto: schema > spec

> "JSON Schema / protobuf definition в week 2, не prose spec. Если другой тул не может имплементировать протокол по одной schema — schema неполная."

**plan/apply UX:** Правильный паттерн, но `plan` у SDP — не план, а estimation. Настоящий plan должен показывать интерфейсы и контракты между компонентами, не "vibes of what we'll build."

### Masad: отдели протокол от методологии

> "Replit evaluate evidence format и model provenance. Мы НЕ adopt workstream hierarchy, multi-agent orchestration, progressive disclosure. Это SDP product opinions, не protocol requirements."

**Предложение:** milestone "Protocol Extraction" между Phase 1 и 2. Убрать из протокола всё SDP-специфичное.

> "Я пришлю двух инженеров если schema чистая. Серьёзно."

### Collison: compliance design doc → Phase 1

> "Каждый месяц задержки с compliance doc = 3 месяца задержки enterprise sales cycle."

Конкретно: data residency, retention policies, RBAC, audit trail immutability. Не build — document the plan.

**Ценообразование:** Per-seat/per-repo annual для enterprise, не usage-based. "Compliance — не фича, которую включают по ситуации. Это checkbox."

### Dario: third-party verification

> "Evidence chain only as trustworthy as the verifier. SDP generates evidence AND verifies evidence = conflict of interest. Design for independent verification."

---

## Protocol design: уроки Terraform

### Что Hashimoto сделал правильно:
- `plan` output = deterministic, diffable artifact = контракт между тулом и оператором

### Что Hashimoto сделал неправильно:
- Plan output слишком verbose → все скроллили мимо. **Опасное выделять, рутинное сворачивать.**
- Не версионировал format рано → экосистема зависела от parsing output

### Совет для SDP:
- Publish schema, NOT spec (JSON Schema / protobuf)
- Draft в week 6, mark it 0.1, commit to breaking at least twice before 1.0

---

## Open-core: уроки HashiCorp

> "Don't start Apache 2.0 if you think you might change later. BSL transition was painful."

**Решение Hashimoto:** Moat = не hash function, а corpus of evidence (failure patterns, risk profiles). Evidence DATA — proprietary by nature.

**Формула:** Open ФОРМАТ evidence, proprietary ХРАНЕНИЕ и АНАЛИЗ. Docker model: OCI spec open, Docker Desktop commercial.

**Dario:** "Evidence format MUST be open. If there's a real AI incident, any auditor/regulator/court must be able to read it."

---

## Evidence format: стандарт или понты?

**Masad:** Forget courts for now. "Engineering managers who want to know 'what did the AI do to my codebase?' — that's the audience."

**Collison:** "Don't call it 'tamper-evident' unless it actually is. Call it 'AI Activity Log v0.1'."

**Hashimoto:** "Ship evidence format in Phase 1. Explicitly say 'NOT a compliance artifact yet.' Then earn it."

---

## Phase 1: что резать?

**Dario:** Не резать — scoping ruthlessly.

> "Protocol spec — не 50-page RFC, а one-page JSON Schema для четырёх примитивов. Model provenance — не solving AI traceability, а тройка (model_id, model_version, timestamp). Evidence chain — не crypto hash linking, а log of what happened."

**Консенсус:** Ship MVP scope of all five, not complete version of three.

Hashimoto: "Write down what's MVP and what's v2 for each one BEFORE you start coding."
