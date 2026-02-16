# Deep-Thinking Integration Guide

## Overview

This guide explains how to integrate the `/think` deep-structured thinking skill into `@idea` and `@design` for complex decision-making.

## What is Deep-Thinking?

The `/think` skill uses parallel expert agents to analyze problems from multiple perspectives:
- **Architecture**: System design, patterns, coupling
- **Security**: Threats, auth, data protection
- **Performance**: Latency, scalability, caching
- **UX**: User experience, workflows
- **Ops**: Deployability, monitoring, maintenance

## When to Use Deep-Thinking

### Use in @idea (Requirements Gathering)

Trigger deep-thinking when:

- **Complex tradeoffs**: Multiple valid approaches with different pros/cons
- **Unknown unknowns**: Requirements unclear, need exploration
- **High-risk features**: Security implications, performance critical
- **User experience ambiguity**: UX needs deep consideration
- **System-level impact**: Changes affect multiple components

### Use in @design (Workstream Design)

Trigger deep-thinking when:

- **Architectural decisions**: Component boundaries, patterns
- **Technology choices**: Language, framework, library selection
- **Integration patterns**: How components interact
- **Data modeling**: Schema design, relationships
- **Quality attributes**: Performance, security, reliability tradeoffs

## Integration Pattern

### @idea Integration (10-Expert Deep-Thinking)

**When to trigger:** After initial requirement gathering, when exploring complex tradeoffs.

**Expert roles for @idea:**

1. **Product Manager** - Business value, user needs, prioritization
2. **Security Expert** - Threats, compliance, data protection
3. **Performance Engineer** - Latency, scalability, efficiency
4. **UX Designer** - User experience, workflows, accessibility
5. **Backend Architect** - API design, data modeling, integration
6. **Frontend Architect** - Component design, state management
7. **DevOps Engineer** - Deployability, monitoring, maintenance
8. **Data Engineer** - Data flow, storage, analytics
9. **QA Engineer** - Testing strategy, edge cases, quality
10. **Technical Writer** - Documentation, onboarding, clarity

**Example workflow:**

```markdown
## @idea with Deep-Thinking Integration

### Phase 1: Initial Requirements (Standard @idea)
- Ask 3-5 core questions
- Gather basic requirements
- Identify feature scope

### Phase 2: Deep-Thinking Analysis (NEW)
Trigger when:
- Complex architectural tradeoffs identified
- Multiple valid implementation approaches
- Unclear requirements with business impact

**Execute:**
```
/think "Analyze: OAuth2 vs. API Keys for authentication in SDP"
```

**Deep-thinking spawns 10 parallel experts:**
- Product Manager: User experience impact, complexity
- Security Expert: OWASP compliance, token storage
- Performance Engineer: Latency impact, caching needs
- UX Designer: Login flows, user friction
- Backend Architect: Integration with existing auth
- Frontend Architect: Component design
- DevOps Engineer: Key rotation, certificate management
- Data Engineer: Audit logging, session storage
- QA Engineer: Testing complexity, edge cases
- Technical Writer: Documentation needs

**Output:**
- Expert analysis from all 10 perspectives
- Synthesis of recommendations
- Tradeoff analysis (pros/cons of each approach)
- Risk assessment by category
- Final recommendation with rationale

### Phase 3: Progressive Disclosure (Standard @idea)
- Use deep-thinking insights to refine questions
- Ask 3-5 follow-up questions per cycle
- Continue until requirements clear

### Phase 4: Workstream Generation (Standard @idea)
- Proceed to @design phase
```

### @design Integration (3-5 Expert Deep-Thinking)

**When to trigger:** During architecture design, when facing complex technical decisions.

**Expert roles for @design:**

1. **Architect** - System design, patterns, modularity
2. **Security** - Threat model, auth, data protection
3. **Performance** - Caching, indexing, optimization
4. **UX** - User flows, component design
5. **Ops** - Monitoring, deployment, maintenance

**Example workflow:**

```markdown
## @design with Deep-Thinking Integration

### Phase 1: Exploration Blocks (Standard @design)
- Explore codebase structure
- Identify integration points
- Check existing patterns

### Phase 2: Architecture Design (Standard @design)
- Propose system architecture
- Define component boundaries
- Specify integration patterns

### Phase 3: Deep-Thinking Analysis (NEW)
Trigger when:
- Multiple architectural patterns valid
- Complex component interactions
- Unclear best approach for quality attributes

**Execute:**
```
/think "Analyze: Monolith vs. Microservices for SDP contract validation"
```

**Deep-thinking spawns 5 parallel experts:**
- Architect: Coupling, complexity, team velocity
- Security: Attack surface, isolation, compliance
- Performance: Resource usage, latency, scaling
- UX: Developer experience, consistency
- Ops: Deployment complexity, monitoring, debugging

**Output:**
- Expert analysis from 5 perspectives
- Pattern recommendation with rationale
- Migration strategy (if applicable)
- Risk mitigation plan
- Implementation considerations

### Phase 4: Workstream Design (Standard @design)
- Generate workstreams based on architecture
- Define dependencies and execution order
- Create WS files
```

## Implementation

### Step 1: Update @idea Skill

Add to `.claude/skills/idea/SKILL.md`:

```markdown
## Deep-Thinking Integration

### When to Trigger Deep-Thinking

After initial requirement gathering, if any of these conditions are met:

1. **Complex Tradeoffs** - Multiple valid approaches with different pros/cons
2. **Unknown Unknowns** - Requirements unclear, need exploration
3. **High Risk** - Security/performance/UX implications
4. **User Experience Ambiguity** - UX needs deep consideration
5. **System Impact** - Changes affect multiple components

### Deep-Thinking Workflow

1. **Identify Analysis Topic**
   - Extract core question from requirements
   - Frame as: "Analyze: {topic} for {feature}"
   - Example: "Analyze: OAuth2 vs. API Keys for authentication"

2. **Execute Deep-Thinking**
   - Use `/think` skill with 10-expert parallel analysis
   - Spawn 10 expert agents simultaneously
   - Wait for all experts to complete

3. **Synthesize Results**
   - Collect expert analyses
   - Identify consensus vs. disagreement
   - Generate synthesis with tradeoffs
   - Provide recommendation with rationale

4. **Continue Progressive Disclosure**
   - Use deep-thinking insights to refine questions
   - Ask follow-up questions based on expert analysis
   - Continue until requirements clear

### Example Deep-Thinking Session

**Context:** User wants OAuth2 authentication for SDP

**Trigger:** Multiple authentication approaches valid (OAuth2, API Keys, JWT, SSO)

**Execute:**
```python
# Spawn 10 parallel experts
Task(subagent_type="general-purpose", prompt="You are the PRODUCT MANAGER expert...")
Task(subagent_type="general-purpose", prompt="You are the SECURITY EXPERT expert...")
# ... 8 more experts
```

**Results after 45 seconds:**
- Product Manager: OAuth2 best for user experience, higher complexity
- Security Expert: OAuth2 industry standard, need token rotation
- Performance Engineer: 50ms overhead per request, acceptable
- UX Designer: SSO preferred (single sign-on), reduces friction
- Backend Architect: Integration complexity high, need rate limiting
- Frontend Architect: Token refresh logic adds complexity
- DevOps Engineer: Certificate rotation, operational overhead
- Data Engineer: Audit logs required for compliance
- QA Engineer: Testing matrix increases 10x
- Technical Writer: Documentation needs increase significantly

**Synthesis:**
- Recommended: OAuth2 with PKCE extension
- Tradeoff: Higher complexity vs. better UX and security
- Migration: Start with OAuth2, add SSO later if needed
- Risks: Token storage, rate limiting, certificate rotation
- Timeline: 3-5 days longer than API Keys

**Follow-up Questions:**
- "OAuth2 provider: Google/GitHub or custom?" → Progressive question
- "Token revocation needed for compliance?" → Progressive question
- "Fallback mechanism if provider down?" → Progressive question
```

### Step 2: Update @design Skill

Add to `.claude/skills/design/SKILL.md`:

```markdown
## Deep-Thinking Integration

### When to Trigger Deep-Thinking

During architecture design, if any of these conditions are met:

1. **Multiple Patterns Valid** - No clear best approach
2. **Complex Interactions** - Components have intricate dependencies
3. **Quality Attribute Conflicts** - Performance vs. security vs. maintainability
4. **Technology Choices** - Language/framework/library selection
5. **Uncertain Approach** - Need exploration before committing

### Deep-Thinking Workflow

1. **Identify Design Decision**
   - Extract core architectural question
   - Frame as: "Analyze: {pattern A} vs. {pattern B} for {component}"
   - Example: "Analyze: Monolith vs. Microservices for contract validation"

2. **Execute Deep-Thinking**
   - Use `/think` skill with 5-expert parallel analysis
   - Spawn 5 expert agents simultaneously
   - Wait for all experts to complete

3. **Synthesize Results**
   - Collect expert analyses
   - Identify pattern recommendations
   - Generate synthesis with migration path
   - Provide implementation guidance

4. **Continue Workstream Design**
   - Use deep-thinking insights for architecture
   - Generate workstreams based on chosen pattern
   - Define dependencies accordingly

### Example Deep-Thinking Session

**Context:** Designing contract validation system

**Trigger:** Monolith vs. Microservices both valid

**Execute:**
```python
# Spawn 5 parallel experts
Task(subagent_type="general-purpose", prompt="You are the ARCHITECT expert...")
Task(subagent_type="general-purpose", prompt="You are the SECURITY expert...")
Task(subagent_type="general-purpose", prompt="You are the PERFORMANCE expert...")
Task(subagent_type="general-purpose", prompt="You are the UX expert...")
Task(subagent_type="general-purpose", prompt="You are the OPS expert...")
```

**Results:**
- Architect: Monolith simpler for now, microservices overkill
- Security: Monolith easier to secure, smaller attack surface
- Performance: Both sufficient, microservices scale better long-term
- UX: No impact (internal feature)
- Ops: Monolith easier to deploy and monitor

**Synthesis:**
- Recommended: Monolith with modular design
- Rationale: Simpler, faster to implement, easier to secure
- Migration: Extract microservices later if needed (pattern: Strangler Fig)
- Implementation: Clean architecture boundaries within monolith

**Workstreams Generated:**
- 00-053-01: Contract domain models
- 00-053-02: Validation logic
- 00-053-03: Report generation
- 00-053-04: API integration
```

## Expert Agent Templates

### @idea Expert Templates

**Product Manager Expert:**
```python
prompt="""You are the PRODUCT MANAGER expert.

QUESTION: {question}
CONTEXT: {context}

Analyze from product perspective:
1. User impact and experience
2. Business value and ROI
3. Implementation complexity
4. Time to market
5. Competitive landscape

Return 3-5 bullet points with rationale."""
```

**Security Expert:**
```python
prompt="""You are the SECURITY EXPERT expert.

QUESTION: {question}
CONTEXT: {context}

Analyze from security perspective:
1. OWASP Top 10 implications
2. Data protection requirements
3. Authentication/authorization needs
4. Compliance considerations (GDPR, SOC2, etc.)
5. Attack surface analysis

Return 3-5 bullet points with severity ratings."""
```

**Performance Engineer:**
```python
prompt="""You are the PERFORMANCE ENGINEER expert.

QUESTION: {question}
CONTEXT: {context}

Analyze from performance perspective:
1. Latency implications
2. Scalability concerns
3. Resource usage (CPU, memory, I/O)
4. Caching opportunities
5. Bottleneck identification

Return 3-5 bullet points with metrics where possible."""
```

### @design Expert Templates

**Architect Expert:**
```python
prompt="""You are the ARCHITECT expert.

DESIGN DECISION: {decision}
CONTEXT: {context}

Analyze from architecture perspective:
1. Coupling and cohesion
2. Modularity and extensibility
3. Pattern applicability
4. Team velocity implications
5. Long-term maintainability

Return 3-5 bullet points with pattern recommendations."""
```

**Security Expert:**
```python
prompt="""You are the SECURITY expert.

DESIGN DECISION: {decision}
CONTEXT: {context}

Analyze from security perspective:
1. Threat model implications
2. Attack surface changes
3. Security control integration
4. Compliance requirements
5. Security testing needs

Return 3-5 bullet points with risk ratings."""
```

## Quality Metrics

### Deep-Thinking Quality Gates

**For @idea (10-expert):**
- All 10 experts must provide analysis
- Synthesis must identify tradeoffs
- Recommendation must have clear rationale
- Open questions explicitly listed

**For @design (5-expert):**
- All 5 experts must provide analysis
- Pattern recommendation provided
- Migration strategy included (if applicable)
- Implementation guidance clear

### Performance Targets

- **10-expert analysis**: Complete in <60 seconds
- **5-expert analysis**: Complete in <30 seconds
- **Synthesis quality**: All perspectives covered
- **Recommendation clarity**: Actionable and unambiguous

## Common Pitfalls

### Don't Over-Use Deep-Thinking

❌ **Don't use for:**
- Simple, obvious decisions
- Questions with one clear answer
- Low-risk choices
- Well-established patterns

✅ **Do use for:**
- Complex tradeoffs
- Multiple valid approaches
- High-impact decisions
- Uncertain requirements

### Don't Ignore Expert Analysis

❌ **Don't:**
- Skip reading expert outputs
- Make decisions without synthesis
- Dismiss expert disagreements
- Proceed with unclear tradeoffs

✅ **Do:**
- Read all expert analyses carefully
- Synthesize conflicting views
- Document disagreements and rationale
- Use insights to refine requirements/design

## Testing

### Manual Testing

**Test @idea deep-thinking:**
1. Request feature with complex requirements
2. Verify deep-thinking triggers at appropriate point
3. Check 10 experts spawned in parallel
4. Verify synthesis quality
5. Verify follow-up questions use insights

**Test @design deep-thinking:**
1. Design feature with architectural choices
2. Verify deep-thinking triggers for complex decisions
3. Check 5 experts spawned in parallel
4. Verify pattern recommendation
5. Verify workstreams reflect chosen pattern

## Rollout Plan

### Phase 1: Documentation (CURRENT)
- Create integration guide
- Define expert templates
- Document when to trigger

### Phase 2: Skill Updates (FUTURE)
- Update @idea skill with deep-thinking workflow
- Update @design skill with deep-thinking workflow
- Add trigger conditions

### Phase 3: Testing (FUTURE)
- Manual testing with complex features
- Validate expert analysis quality
- Measure synthesis effectiveness

### Phase 4: Rollout (FUTURE)
- Enable deep-thinking in @idea
- Enable deep-thinking in @design
- Monitor quality and performance
- Iterate based on feedback

## References

- `/think` skill: `.claude/skills/think/SKILL.md`
- @idea skill: `.claude/skills/idea/SKILL.md`
- @design skill: `.claude/skills/design/SKILL.md`
- Progressive disclosure: `docs/feature/progressive-disclosure.md`

---

**Version:** 1.0.0
**Status:** Draft - Ready for implementation
**Last Updated:** 2026-02-08
