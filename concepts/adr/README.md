# Architecture Decision Records (ADR)

ADRs document important architectural decisions along with their context and consequences.

## Why ADRs?

Six months from now, you (or your teammates) will ask:
- "Why did we use PostgreSQL instead of MongoDB?"
- "Why is authentication done this way?"
- "What were we thinking?"

ADRs answer these questions.

## When to Write an ADR

Write an ADR when you:
- Choose between alternatives (library A vs B)
- Make trade-offs (security vs convenience)
- Deviate from common patterns
- Introduce new technology
- Make decisions that are hard to reverse

**Don't write an ADR for:**
- Obvious choices with no real alternatives
- Easily reversible decisions
- Implementation details

## ADR Format

### Template

```markdown
# ADR-NNNN: [Short Title]

## Status

[Proposed | Accepted | Deprecated | Superseded by ADR-XXXX]

## Context

[What is the issue that we're seeing that motivates this decision?
What forces are at play? Include technical, business, and social factors.]

## Decision

[What is the change that we're proposing or have agreed to implement?]

## Consequences

[What becomes easier or more difficult to do because of this change?
Include both positive and negative consequences.]
```

### Example: Database Choice

```markdown
# ADR-0001: Use PostgreSQL for Primary Database

## Status

Accepted

## Context

We need a database for our user management system. Requirements:
- ACID compliance for financial transactions
- Complex queries with joins
- JSON support for flexible schema parts
- Team familiarity

Options considered:
1. PostgreSQL - Relational, ACID, JSON support
2. MongoDB - Document store, flexible schema
3. MySQL - Relational, simpler than PostgreSQL

## Decision

Use PostgreSQL as our primary database.

## Consequences

### Positive
- Strong ACID guarantees for transactions
- Rich query capabilities with SQL
- JSON columns for semi-structured data
- Excellent tooling ecosystem
- Team has experience

### Negative
- Requires more upfront schema design
- Horizontal scaling is more complex than MongoDB
- Need to manage migrations carefully
```

### Example: Authentication Approach

```markdown
# ADR-0002: Use JWT for API Authentication

## Status

Accepted

## Context

Our API needs authentication. The system is stateless and may scale
to multiple instances. We need to support:
- Web frontend
- Mobile apps
- Third-party integrations

Options:
1. Session-based auth (server-side sessions)
2. JWT (JSON Web Tokens)
3. OAuth2 with opaque tokens

## Decision

Use JWT (JSON Web Tokens) for API authentication:
- Access tokens: 15-minute expiry
- Refresh tokens: 7-day expiry, stored in httpOnly cookies
- Tokens contain user ID and roles

## Consequences

### Positive
- Stateless - easy to scale horizontally
- Self-contained - no database lookup per request
- Works well across services
- Standard format, good library support

### Negative
- Cannot revoke individual tokens (use short expiry + refresh)
- Token size larger than session ID
- Must handle token refresh logic on client
- Security requires careful implementation (HTTPS, secure storage)
```

### Example: Reversing a Decision

```markdown
# ADR-0003: Replace Redis Cache with In-Memory Cache

## Status

Accepted (supersedes ADR-0002 caching section)

## Context

ADR-0002 introduced Redis for caching. After 6 months:
- We've never needed distributed caching
- Redis adds operational complexity
- Single-instance deployment is sufficient

## Decision

Replace Redis with in-memory caching (Python's functools.lru_cache
or similar) for:
- User session data
- Configuration caching
- Frequently accessed lookups

Keep Redis only for rate limiting (needs to be shared).

## Consequences

### Positive
- Simpler deployment (one less service)
- Lower latency (no network hop)
- Reduced operational burden

### Negative
- Cache lost on restart
- Cannot scale to multiple instances without re-adding Redis
- Rate limiting still needs Redis
```

## ADR Numbering

Use sequential numbers with leading zeros:
```
docs/adr/
├── 0001-use-postgresql.md
├── 0002-jwt-authentication.md
├── 0003-replace-redis.md
└── 0004-adopt-clean-architecture.md
```

## Status Values

| Status | Meaning |
|--------|---------|
| **Proposed** | Under discussion, not yet decided |
| **Accepted** | Decision made, being implemented |
| **Deprecated** | No longer recommended, but may exist |
| **Superseded** | Replaced by another ADR |

## Tips

### Keep It Short
ADRs should be 1-2 pages maximum. If longer, you're over-explaining.

### Focus on Why
The decision itself is usually simple. The value is in the context and trade-offs.

### Include Alternatives
Explain what you didn't choose and why. This prevents relitigating decisions.

### Date Your Decisions
Context changes. What made sense in 2023 might not in 2025. Include dates.

### Make Them Findable
- Use clear titles
- Keep in standard location (docs/adr/)
- Reference from code comments when relevant

## Using ADRs with AI

### Creating ADRs
```
"We're choosing between Redis and Memcached for caching.
Create an ADR documenting this decision with trade-offs.
We chose Redis because [reasons]."
```

### Referencing ADRs
```
"Before implementing caching, check docs/adr/ for existing decisions
about caching strategy."
```

### Updating ADRs
```
"Our caching needs have changed. Create a new ADR that supersedes
ADR-0003. We now need distributed caching because [reasons]."
```

## Further Reading

- [Michael Nygard's original ADR article](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)
- [ADR GitHub organization](https://adr.github.io/)
