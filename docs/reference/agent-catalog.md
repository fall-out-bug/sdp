# Agent Catalog

Complete catalog of all 19 agents in the SDP multi-agent system.

## Phase 1A: @vision Agents (7)

### 1. Product Strategist
**Role:** Analyzes product vision and market fit
**When to use:** Starting new product or major pivot
**Input:** Product idea
**Output:** Product positioning, market analysis

### 2. Market Analyst
**Role:** Analyzes market size, competition, trends
**When to use:** Need market intelligence
**Input:** Product domain
**Output:** Market size, competitive landscape

### 3. Technical Architect
**Role:** Evaluates technical feasibility
**When to use:** Technical decisions needed
**Input:** Product requirements
**Output:** Architecture recommendations

### 4. UX Designer
**Role:** Analyzes user experience needs
**When to use:** UX-critical features
**Input:** User scenarios
**Output:** UX recommendations

### 5. Business Analyst
**Role:** Analyzes business model and ROI
**When to use:** Need business case
**Input:** Product concept
**Output:** Business model, revenue projections

### 6. Growth Strategist
**Role:** Plans user acquisition and growth
**When to use:** Need growth strategy
**Input:** Product concept
**Output:** Growth channels, metrics

### 7. Risk Analyst
**Role:** Identifies and mitigates risks
**When to use:** Risk assessment needed
**Input:** Product plan
**Output:** Risk matrix, mitigation strategies

## Phase 1B: @reality Agents (8)

### 8. Architecture Reviewer
**Role:** Reviews codebase architecture
**When to use:** Architecture audit
**Input:** Codebase
**Output:** Architecture assessment

### 9. Quality Analyst
**Role:** Analyzes code quality and technical debt
**When to use:** Quality assessment
**Input:** Codebase
**Output:** Quality metrics, debt analysis

### 10. Testing Specialist
**Role:** Reviews test coverage and strategy
**When to use:** Testing audit
**Input:** Codebase
**Output:** Test coverage, gaps

### 11. Security Analyst
**Role:** Analyzes security posture
**When to use:** Security assessment
**Input:** Codebase
**Output:** Security findings, recommendations

### 12. Performance Analyst
**Role:** Analyzes performance characteristics
**When to use:** Performance audit
**Input:** Codebase
**Output:** Performance bottlenecks

### 13. Documentation Specialist
**Role:** Reviews documentation quality
**When to use:** Docs assessment
**Input:** Codebase + docs
**Output:** Docs gaps, quality score

### 14. Technical Debt Analyst
**Role:** Identifies and prioritizes technical debt
**When to use:** Debt assessment
**Input:** Codebase
**Output:** Debt inventory, prioritization

### 15. Standards Compliance Analyst
**Role:** Checks coding standards compliance
**When to use:** Standards audit
**Input:** Codebase
**Output:** Compliance violations, recommendations

## Phase 2: Review Agents (2)

### 16. Implementer Agent
**Role:** Executes workstreams with TDD
**When to use:** Executing workstream
**Input:** Workstream specification
**Output:** Implemented code, test results

### 17. Spec Compliance Reviewer
**Role:** Verifies implementation matches specification
**When to use:** After implementation
**Input:** Spec + code
**Output:** Compliance report

## Phase 3: Orchestrator (1)

### 18. Parallel Orchestrator
**Role:** Coordinates parallel execution of workstreams
**When to use:** Multi-workstream features
**Input:** Dependency graph
**Output:** Coordinated execution

## Phase 4: Synthesis Agents (3)

### 19. Agent Synthesizer
**Role:** Resolves conflicts between agent proposals
**When to use:** Multiple agents disagree
**Input:** Agent proposals
**Output:** Synthesized solution

### 20. Rules Engine
**Role:** Applies synthesis rules in priority order
**When to use:** Synthesizing proposals
**Input:** Proposals
**Output:** Synthesis result

### 21. Hierarchical Supervisor
**Role:** Coordinates specialist agents
**When to use:** Multi-agent coordination
**Input:** Task
**Output:** Coordinated decision

## Integration Examples

### @vision Usage
```bash
@vision "AI-powered task manager"
# Launches 7 expert agents in parallel
# Generates: PRODUCT_VISION.md, PRD.md, ROADMAP.md
```

### @reality Usage
```bash
@reality --quick
# Launches 8 expert agents
# Generates: Reality report (health, gaps, debt)
```

### Synthesis Usage
```go
supervisor := synthesis.NewSupervisor(engine, 5)
supervisor.RegisterAgent(agent1)
supervisor.RegisterAgent(agent2)
decision := supervisor.MakeDecision(task)
```

---

**Generated:** 2026-02-07  
**SDP Version:** 4.0.0
