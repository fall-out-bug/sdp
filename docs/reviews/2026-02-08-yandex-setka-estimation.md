# Demo AdServer → Яндекс Сетка: Оценка масштаба работ

> **Status:** Research complete
> **Date:** 2026-02-08
> **Goal:** Оценить количество фич и агентского времени для превращения demo-adserver в систему, конкурентоспособную с Яндекс Сетка

---

## Table of Contents

1. [Overview](#overview)
2. [Current State](#current-state)
3. [Expert Analysis Results](#expert-analysis-results)
4. [Implementation Plan](#implementation-plan)
5. [Timeline Scenarios](#timeline-scenarios)
6. [Agent-Time Estimation](#agent-time-estimation)

---

## Overview

### Goals

1. **Production-Ready Ad Server** — Система, способная работать в промышленной эксплуатации
2. **RTB Capability** — Поддержка Real-Time Bidding для программатик-покупки
3. **Advanced Targeting** — Поведенческий таргетинг, ретаргетинг, lookalike-аудитории
4. **Real-Time Analytics** — Мгновенная отчётность и статистика
5. **Yandex Scale** — Масштабируемость до миллионов запросов в секунду

### Key Decisions

| Aspect | Decision | Approach |
|--------|----------|-----------|
| RTB Protocol | Minimal OpenRTB 2.5 | First-price auction, single DSP initially |
| Targeting | Incremental layered | User profiles → Frequency caps → Retargeting → Behavioral |
| Analytics | Kafka + ClickHouse | Event streaming + columnar storage |
| Scalability | Service decomposition | Separate delivery/tracking/management services |
| Data Storage | Polyglot persistence | PostgreSQL (transactional) + ClickHouse (analytics) + S3 (creatives) |
| Security | Progressive implementation | Phase 1: Rule-based → Phase 2: Heuristics → Phase 3: ML/API |
| UI/Dashboard | Incremental MVP | Extend existing portals, extract shared components |
| Billing | Wallet-first (pre-payment) | Russian market standard |
| CI/CD | GitOps with ArgoCD | OpenTelemetry observability stack |
| Migration | Phased feature-based | Strangler pattern with SDP @vision/@feature/@oneshot |

---

## Current State

**Уже реализовано:**
- ✅ Clean Architecture (Domain → Application → Infrastructure → Presentation)
- ✅ Delivery Engine (banner selection with cache-first strategy)
- ✅ Basic Targeting (geo, devices, OS, browsers, time of day)
- ✅ Tracking (impressions, clicks)
- ✅ JWT Authentication
- ✅ Rate Limiting (Redis-based)
- ✅ Publisher Portal (Next.js 14, dashboard)
- ✅ Advertiser Portal (Next.js 14, campaign wizard)
- ✅ Web SDK (<5KB)

**Критические пробелы:**
- ❌ Budget enforcement (campaign budget fields exist but never checked)
- ❌ Frequency capping
- ❌ RTB capability
- ❌ Real-time analytics
- ❌ Behavioral targeting
- ❌ Billing system
- ❌ Fraud detection (beyond basic sanitization)
- ❌ Production CI/CD

---

## Expert Analysis Results

### 1. Ядро Ad Server (190-252 agent-hours)

| Component | Status | Estimate | Priority |
|-----------|--------|----------|----------|
| Budget Enforcement | Missing | 24-32h | P0 |
| Frequency Capping | Missing | 24-32h | P0 |
| Campaign Scheduling | Partial | 24-32h | P0 |
| Waterfall Algorithm | Missing | 10-12h | P1 |
| Reporting Dashboards | Missing | 38-50h | P1 |

### 2. RTB Protocol (80-120 agent-hours)

**Решение:** Minimal OpenRTB 2.5 Bidder with First-Price Auction

| Phase | Workstreams | Hours |
|-------|-------------|-------|
| Foundation | 4 WS | 40-50h |
| Integration | 4 WS | 20-30h |
| Quality | 4 WS | 20-40h |

### 3. Таргетинг и профилирование (122-222 agent-hours)

**Решение:** Incremental Layered Approach

| Phase | Feature | Hours |
|-------|---------|-------|
| 1 | User Profile Store (Redis) | 24-32h |
| 2 | Frequency Capping Service | 18-24h |
| 3 | Retargeting Pixel | 20-28h |
| 4 | Segment Engine | 28-36h |
| 5 | Behavioral Targeting | 32-42h |
| 6 | Lookalike Audiences | 24-32h |

**Minimum viable Яндекс Сетка targeting:** Phases 1-5 = **122-162 hours**

### 4. Отчётность и аналитика (160-220 agent-hours)

**Решение:** Kafka + ClickHouse

| Phase | Hours | Deliverables |
|-------|-------|--------------|
| 1: Event Streaming | 40-50h | Kafka producers/consumers |
| 2: ClickHouse | 50-70h | Real-time aggregation |
| 3: Analytics API | 30-40h | RESTful reports |
| 4: Real-Time Stats | 20-30h | WebSocket live stats |
| 5: BI Integration | 20-30h | Grafana/Metabase |

### 5. Масштабируемость (240-280 agent-hours)

**Решение:** Service Decomposition

| Phase | Hours | Focus |
|-------|-------|-------|
| 1: Foundation | 80h | K8s, Kafka, read replicas |
| 2: Service Extraction | 120h | Separate delivery/tracking |
| 3: Optimization | 80h | Multi-level caching, auto-scaling |

**Scaling Roadmap:** 1K → 10K → 50K → 200K → 500K → 1M QPS

### 6. Хранение данных (82-108 agent-hours)

**Решение:** Evolutionary Polyglot Persistence

| Task | Hours |
|------|-------|
| ClickHouse schema design | 8-12h |
| ETL pipeline (dual-write) | 12-16h |
| Data migration scripts | 16-20h |
| S3 integration | 8-10h |
| Repository abstraction | 10-14h |
| Monitoring & testing | 28-36h |

### 7. Security & Fraud Detection (152-192 agent-hours)

**Решение:** Progressive Implementation

| Phase | Hours | Protection |
|-------|-------|------------|
| 1: Rule-based | 40-48h | 60-70% fraud detection |
| 2: Heuristics | 48-64h | 75-85% fraud detection |
| 3: Advanced | 64-80h | 85-95% fraud detection |

### 8. UI/Dashboard (428-720 agent-hours)

**Решение:** Incremental MVP Enhancement

| Portal | MVP | Full-Featured |
|--------|-----|---------------|
| Publisher | 60h | 112h |
| Advertiser | 104h | 212h |
| Admin | 148h | 212h |
| Shared Library | 116h | 184h |

**Recommended:** MVP path = **~428 hours** (~11 weeks)

### 9. Биллинг и интеграции (160 agent-hours)

**Решение:** Wallet-First Architecture

| Feature | Hours |
|---------|-------|
| Wallet domain + repo | 16h |
| Yandex.Kassa integration | 48h |
| Auto-recharge service | 24h |
| Invoice generation | 16h |
| Yandex Metrica integration | 8h |
| CRM integration | 12h |
| Testing | 24h |
| Notification service | 12h |

### 10. CI/CD & Observability (80-102 agent-hours)

**Решение:** GitOps with ArgoCD + OpenTelemetry

| Phase | Hours |
|-------|-------|
| Kubernetes manifests | 16-20h |
| Helm charts | 12-16h |
| CI/CD pipeline | 16-20h |
| Observability stack | 24-30h |
| Alert rules & runbooks | 12-16h |

### 11. SDP Migration (Фасилитация, не дополнительная работа)

**Решение:** Phased Feature-Based Migration

Использовать существующую SDP инфраструктуру:
- `@vision` → PRD с roadmap
- `@feature` → Workstreams (5-30 WS per feature)
- `@oneshot` → Parallel execution
- `@review` → Quality gates
- `@deploy` → GitFlow merges

**Формат миграции:**
```
F100-F105: Foundation (Infrastructure)
F110-F115: Core Domain (Campaign, Banner, Tracking)
F116-F118: API Layer (Delivery, Management, Auth)
F119-F120: Frontend (Portals)
```

---

## Implementation Plan

### Phase 1: Foundation (2-3 months)

**Target:** Базовая инфраструктура для production

- [x] Existing: Clean Architecture, Delivery Engine, Basic Targeting
- [ ] Budget Enforcement (24-32h)
- [ ] Frequency Capping (24-32h)
- [ ] User Profile Store (24-32h)
- [ ] Kafka Event Streaming (40-50h)
- [ ] ClickHouse Integration (50-70h)
- [ ] Kubernetes manifests (16-20h)

**Deliverables:**
- Бюджет enforcement работает
- Frequency caps防止 ad fatigue
- Event pipeline operational
- Analytics database ready

### Phase 2: Core Features (2-4 months)

**Target:** Функциональность уровня "Яндекс Сетка"

- [ ] RTB Protocol (80-120h)
- [ ] Retargeting Pixel (20-28h)
- [ ] Segment Engine (28-36h)
- [ ] Behavioral Targeting (32-42h)
- [ ] Service Decomposition (120h)
- [ ] Billing System (160h)
- [ ] Progressive Security (88-112h)

**Deliverables:**
- RTB integration с first-price auction
- Behavioral targeting operational
- Payment processing via Yandex.Kassa
- Anti-fraud basics deployed

### Phase 3: Advanced Features (2-3 months)

**Target:** Продвинутые возможности для конкурентного преимущества

- [ ] Lookalike Audiences (24-32h)
- [ ] Advanced Analytics (30-40h)
- [ ] Real-Time Stats WebSocket (20-30h)
- [ ] Multi-DSP RTB (120-180h)
- [ ] Advanced Fraud Detection (64-80h)

**Deliverables:**
- ML-based lookalike modeling
- Real-time dashboards
- Multiple DSP integrations
- Sophisticated fraud detection

### Phase 4: Scale & Optimize (1-2 months)

**Target:** Enterprise-scale performance

- [ ] Performance Optimization (120h)
- [ ] Full Observability Stack (60-80h)
- [ ] Auto-scaling policies (40h)
- [ ] CDN Integration (40h)

**Deliverables:**
- Sub-100ms RTB latency p95
- 1M+ QPS capability
- Production-ready monitoring
- Global CDN deployment

---

## Timeline Scenarios

### Scenario 1: Aggressive (Single Developer + AI)

| Metric | Value |
|--------|-------|
| **Timeline** | 5.5-8.5 months |
| **Confidence** | 60% |
| **Team** | 1 FT developer + AI agents |
| **Parallelism** | 5 agents simultaneously |
| **Risk** | Burnout, integration issues |

### Scenario 2: Realistic (Small Team + AI) ⭐ Recommended

| Metric | Value |
|--------|-------|
| **Timeline** | **4.5-7.5 months** |
| **Confidence** | 80% |
| **Team** | 2-3 FT developers (backend + frontend + devOps) |
| **Parallelism** | 5 agents + specialized human roles |
| **Risk** | Managed, peer review available |

### Scenario 3: Conservative (Enterprise Pace)

| Metric | Value |
|--------|-------|
| **Timeline** | 7-12 months |
| **Confidence** | 95% |
| **Team** | 3-5 developers + QA + DevOps + Product |
| **Parallelism** | 3 agents (cautious) |
| **Risk** | Low quality, extensive process overhead |

---

## Agent-Time Estimation

### Summary by Phase

| Phase | Agent-Hours | Human-Hours | Calendar Weeks |
|-------|-------------|-------------|----------------|
| **1: Foundation** | 200-280 | 60-84 | 8-12 weeks |
| **2: Core Features** | 600-800 | 180-240 | 16-24 weeks |
| **3: Advanced** | 300-400 | 90-120 | 8-12 weeks |
| **4: Scale & Optimize** | 200-280 | 60-84 | 4-8 weeks |
| **TOTAL** | **1,300-1,760** | **390-528** | **36-56 weeks** |

### Agent Productivity Assumptions

- **1 agent-turn** = ~5-10 minutes (average 7.5 min)
- **Small WS** = 10-15 turns = 1.25-1.9 hours
- **Medium WS** = 20-30 turns = 2.5-3.8 hours
- **Large WS** = 40-60 turns = 5-7.5 hours
- **Parallel execution** = 3-5 agents simultaneously = **3x speedup**
- **Human overhead** = ~30% for review, verification, decisions

### Net Multiplier

| Factor | Impact | Multiplier |
|--------|--------|------------|
| AI acceleration (Claude + SDP) | + | 3-5x |
| Parallel execution | + | 3x |
| Clean Architecture (existing) | + | 1.2x |
| Test coverage (>80%) | + | 1.1x |
| Ad tech domain complexity | - | 0.7x |
| RTB integration uncertainty | - | 0.8x |
| Performance (<100ms) | - | 0.8x |
| **Net** | | **~3-4x speedup** |

### Final Estimate

**Recommended approach:** Realistic Scenario (Small Team + AI)

| Metric | Value |
|--------|-------|
| **Expected completion** | **6 months** |
| **80% confidence interval** | **4.5-7.5 months** |
| **Total Agent-Hours** | ~1,500 hours |
| **Total Human-Hours** | ~450 hours |
| **Key risks** | RTB integration, performance bottlenecks, fraud detection |

---

## Success Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| **QPS** | ~1K | 1M+ |
| **RTB Latency p95** | N/A | <100ms |
| **Test Coverage** | 88% (backend) | ≥80% (all) |
| **Uptime** | N/A | 99.9% |
| **Fraud Detection** | ~0% | 85%+ |
| **Time to Market** | N/A | 6 months |

---

## Key Takeaways

1. **Feasible:** With SDP + AI agents, 6-month timeline is achievable for "Яндекс Сетка" parity
2. **Parallelism is critical:** 3-5 agents running simultaneously = 3x speedup
3. **Phased approach reduces risk:** Each phase delivers value incrementally
4. **AI acceleration is unprecedented:** Traditional estimates would be 18-24 months
5. **Small team viable:** 2-3 developers with AI assistance can compete with large teams

---

## Next Steps

1. **Review this document** — Decide on timeline and team composition
2. **Run `@vision`** — Generate comprehensive PRD with roadmap
3. **Create features** — Break down into SDP features (F100-F120 series)
4. **Start execution** — Use `@oneshot` for parallel workstream execution

---

**Sources:**
- All 12 expert analyses conducted via think-through:expert agents
- Current codebase analysis of demo-adserver
- Industry benchmarks for ad server development
- SDP v0.6.0 framework capabilities
