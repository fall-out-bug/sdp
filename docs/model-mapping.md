# Model Mapping Registry

**Purpose:** Updatable registry mapping capability tiers (T0–T3) to specific models and requirements.

**Key Principle:** This registry is **separate from SDP protocol and WS templates**. Updating the model list does **not** require changes to SDP specifications or workstream definitions.

---

## Capability Tier → Model Mapping

### T0 — Architect (Contract Writer)

**Required Capabilities:**
- Strong reasoning and architectural decision-making
- Ability to create/modify contracts and interfaces
- Test generation and validation
- Workstream decomposition and dependency resolution
- Context: 32K+ tokens recommended

**Recommended Models:**

| Provider | Model | Cost ($/1M) | Availability | Context | Tool Use | Notes |
|----------|-------|-------------|--------------|---------|----------|-------|
| Anthropic | Claude Sonnet 4.5 | 3.00 | 99.9% | 200K | ✅ Full | Primary choice for /design |
| Anthropic | Claude Opus 4 | 15.00 | 99.9% | 200K | ✅ Full | Highest quality, higher cost |
| OpenAI | GPT-4 Turbo | 10.00 | 99.5% | 128K | ✅ Full | Alternative to Sonnet |
| OpenAI | GPT-4o | 5.00 | 99.5% | 128K | ✅ Full | Good balance of quality/cost |
| Google | Gemini Pro 1.5 | 3.50 | 99.0% | 1M+ | ✅ Full | Very large context support |
| Local | Llama 3.1 70B (4-bit) | 0.00 | 100% | 128K | ⚠️ Limited | Requires 40GB+ VRAM |
| Local | Qwen 2.5 72B (4-bit) | 0.00 | 100% | 32K | ⚠️ Limited | Good reasoning, Chinese support |

**Minimum Requirements:**
- Context window: ≥ 32K tokens
- Tool use: Full support (function calling)
- Code quality: High (architectural decisions)

---

### T1 — Integrator (Complex Build)

**Required Capabilities:**
- Multi-file modifications within given architecture
- Algorithm and state machine implementation
- Component integration
- Context: 16K+ tokens recommended

**Recommended Models:**

| Provider | Model | Cost ($/1M) | Availability | Context | Tool Use | Notes |
|----------|-------|-------------|--------------|---------|----------|-------|
| Anthropic | Claude Sonnet 4.5 | 3.00 | 99.9% | 200K | ✅ Full | Can handle T1, but overkill |
| Anthropic | Claude Haiku 4 | 0.25 | 99.9% | 200K | ✅ Full | **Recommended** for T1 |
| OpenAI | GPT-4o-mini | 0.15 | 99.5% | 128K | ✅ Full | Cost-effective for T1 |
| OpenAI | GPT-4 Turbo | 10.00 | 99.5% | 128K | ✅ Full | Higher quality, higher cost |
| Google | Gemini Flash 1.5 | 0.075 | 99.0% | 1M+ | ✅ Full | Fast, large context |
| Local | Llama 3.1 13B (4-bit) | 0.00 | 100% | 128K | ⚠️ Limited | Requires 8GB+ VRAM |
| Local | DeepSeek-Coder 33B (4-bit) | 0.00 | 100% | 16K | ⚠️ Limited | Code-focused, requires 20GB+ VRAM |
| Local | Qwen 2.5 14B (4-bit) | 0.00 | 100% | 32K | ⚠️ Limited | Requires 8GB+ VRAM |

**Minimum Requirements:**
- Context window: ≥ 16K tokens
- Tool use: Full support preferred, limited acceptable
- Code quality: Medium-High (complex logic)

---

### T2 — Implementer (Contract-Driven Build)

**Required Capabilities:**
- Function body implementation from signatures
- Simple data transformations
- Following strict contracts (read-only Interface/Tests)
- Context: 8K+ tokens sufficient

**Recommended Models:**

| Provider | Model | Cost ($/1M) | Availability | Context | Tool Use | Notes |
|----------|-------|-------------|--------------|---------|----------|-------|
| Anthropic | Claude Haiku 4 | 0.25 | 99.9% | 200K | ✅ Full | **Primary choice** for T2 |
| OpenAI | GPT-4o-mini | 0.15 | 99.5% | 128K | ✅ Full | **Recommended** alternative |
| OpenAI | GPT-3.5 Turbo | 0.50 | 99.0% | 16K | ✅ Full | Cost-effective, lower quality |
| Google | Gemini Flash 1.5 | 0.075 | 99.0% | 1M+ | ✅ Full | Fast, very cost-effective |
| Local | Llama 3.1 8B (4-bit) | 0.00 | 100% | 128K | ❌ None | **Baseline** for local T2 |
| Local | CodeLlama 7B (4-bit) | 0.00 | 100% | 16K | ❌ None | Code-focused baseline |
| Local | DeepSeek-Coder 6.7B (4-bit) | 0.00 | 100% | 16K | ❌ None | Code-focused, 8GB VRAM |
| Local | Qwen 2.5 7B (4-bit) | 0.00 | 100% | 32K | ❌ None | Multilingual support |
| Local | Mistral 7B (4-bit) | 0.00 | 100% | 8K | ❌ None | Smallest viable local model |

**Minimum Requirements:**
- Context window: ≥ 8K tokens
- Tool use: Not required (contract-driven)
- Code quality: Medium (straightforward implementations)
- **Local models:** Minimum 7B parameters, 4-bit quantization, 8GB+ VRAM

**Local Model Baseline:**
- **Practical minimum:** 7B code model in 4-bit quantization (e.g., CodeLlama 7B, Qwen 2.5 7B)
- **Comfortable:** 8B-13B models (e.g., Llama 3.1 8B, Qwen 2.5 14B)
- **VRAM requirement:** 8GB minimum for 7B 4-bit, 12GB+ for 13B 4-bit

---

### T3 — Autocomplete (Micro Build)

**Required Capabilities:**
- Boilerplate generation
- Type conversions
- Getters/setters
- Short, obvious function bodies
- Context: 4K+ tokens sufficient

**Recommended Models:**

| Provider | Model | Cost ($/1M) | Availability | Context | Tool Use | Notes |
|----------|-------|-------------|--------------|---------|----------|-------|
| Anthropic | Claude Haiku 4 | 0.25 | 99.9% | 200K | ✅ Full | Overkill but works |
| OpenAI | GPT-4o-mini | 0.15 | 99.5% | 128K | ✅ Full | **Recommended** for T3 |
| OpenAI | GPT-3.5 Turbo | 0.50 | 99.0% | 16K | ✅ Full | Very cost-effective |
| Google | Gemini Flash 1.5 | 0.075 | 99.0% | 1M+ | ✅ Full | Fastest, cheapest cloud option |
| Local | Llama 3.1 8B (4-bit) | 0.00 | 100% | 128K | ❌ None | Works for T3 |
| Local | CodeLlama 7B (4-bit) | 0.00 | 100% | 16K | ❌ None | Code-focused |
| Local | TinyLlama 1.1B (4-bit) | 0.00 | 100% | 2K | ❌ None | **Minimum viable** (4GB VRAM) |
| Local | Phi-2 2.7B (4-bit) | 0.00 | 100% | 2K | ❌ None | Small, efficient (4GB VRAM) |

**Minimum Requirements:**
- Context window: ≥ 4K tokens
- Tool use: Not required
- Code quality: Low-Medium (obvious implementations)
- **Local models:** Can use smaller models (1B-3B) for T3, 4GB+ VRAM acceptable

---

## Model Selection Criteria

### Context Window

| Tier | Minimum | Recommended | Notes |
|------|---------|-------------|-------|
| T0 | 32K | 128K+ | Large context for architecture decisions |
| T1 | 16K | 32K+ | Multi-file modifications |
| T2 | 8K | 16K+ | Single function focus |
| T3 | 4K | 8K+ | Micro-tasks |

### Tool Use Requirements

| Tier | Required | Notes |
|------|----------|-------|
| T0 | ✅ Full | Function calling for /design |
| T1 | ✅ Preferred | Full support recommended |
| T2 | ❌ Not required | Contract-driven, no tool use needed |
| T3 | ❌ Not required | Simple completions |

### Code Quality Expectations

| Tier | Quality Level | Examples |
|------|---------------|----------|
| T0 | High | Architectural patterns, design decisions |
| T1 | Medium-High | Complex algorithms, state machines |
| T2 | Medium | Straightforward implementations |
| T3 | Low-Medium | Boilerplate, obvious code |

### Local Model Requirements (VRAM)

| Tier | Minimum Model | VRAM (4-bit) | VRAM (8-bit) | Notes |
|------|---------------|--------------|---------------|-------|
| T0 | 70B | 40GB | 80GB | Rarely used locally |
| T1 | 13B | 8GB | 16GB | Comfortable for T1 |
| T2 | 7B | 8GB | 14GB | **Baseline** for local T2 |
| T3 | 1B-3B | 4GB | 6GB | Can use smaller models |

**Quantization Notes:**
- **4-bit (Q4):** Recommended for local models, ~50% quality retention
- **8-bit (Q8):** Better quality, ~75% quality retention, 2x VRAM
- **16-bit (FP16):** Full quality, 2x VRAM of 8-bit (rarely used)

---

## Update Process

### When to Update This Registry

Update this document when:

1. **New models are released** that better fit a tier
2. **Model capabilities change** (context window, tool use support)
3. **Cost structures change** (new pricing tiers)
4. **Local model benchmarks** show better options
5. **Provider deprecates models** (remove outdated entries)

### How to Update

1. **No SDP/WS changes required** — This registry is independent
2. **Update the relevant tier table** — Add/remove/modify model entries
3. **Update criteria if needed** — Adjust minimum requirements based on benchmarks
4. **Add notes** — Document why a model was added/removed
5. **Version the change** — Add a changelog entry (see below)

### Changelog Format

```markdown
## Changelog

### 2026-01-XX
- Added: GPT-4o-mini to T2 (cost-effective alternative to Haiku)
- Removed: GPT-3.5 Turbo from T1 (upgraded to T2-only)
- Updated: Local T2 baseline to 8GB VRAM (was 7GB, accounting for overhead)
```

---

## Provider-Specific Notes

### Anthropic (Claude)

- **Sonnet 4.5:** Best for T0, can handle T1-T3 but overkill
- **Haiku 4:** Primary choice for T2/T3, cost-effective
- **Opus 4:** Highest quality, use only for complex T0 tasks

### OpenAI

- **GPT-4 Turbo:** Good T0 alternative, higher cost than Sonnet
- **GPT-4o:** Balanced quality/cost for T0-T1
- **GPT-4o-mini:** Recommended for T2/T3, very cost-effective
- **GPT-3.5 Turbo:** T3 only, lowest cost cloud option

### Google

- **Gemini Pro 1.5:** T0 option with massive context (1M+ tokens)
- **Gemini Flash 1.5:** Fast, cheap, good for T1-T3

### Local Models (Ollama, llama.cpp, vLLM)

- **Quantization:** Always use 4-bit for VRAM efficiency
- **Code-focused models:** CodeLlama, DeepSeek-Coder preferred for T2/T3
- **General models:** Llama, Qwen, Mistral work but code models perform better
- **VRAM planning:** Add 2GB overhead for system + quantization

---

## Cost Considerations

### Cloud Models (per 1M tokens, approximate)

| Tier | Primary Model | Cost | Alternative | Cost |
|------|---------------|------|-------------|------|
| T0 | Sonnet 4.5 | $3 | GPT-4 Turbo | $10 |
| T1 | Haiku 4 | $0.25 | GPT-4o-mini | $0.15 |
| T2 | Haiku 4 | $0.25 | GPT-4o-mini | $0.15 |
| T3 | GPT-4o-mini | $0.15 | GPT-3.5 Turbo | $0.50 |

### Local Models

- **Electricity cost:** ~$0.10-0.50 per 1M tokens (depends on hardware)
- **Hardware cost:** One-time (GPU purchase)
- **Best for:** High-volume T2/T3 execution

**Cost Optimization Strategy:**
- Use cloud for T0 (low volume, high value)
- Use cloud for T1 (medium volume, complex)
- Use local for T2/T3 (high volume, simple)

---

## Testing & Validation

### Model Compatibility Testing

When adding a new model to this registry:

1. **Run T2-ready WS** — Test with 5-10 validated T2 workstreams
2. **Measure success rate** — Target ≥ 90% for cloud, ≥ 80% for local
3. **Check context limits** — Verify model handles typical WS size
4. **Validate tool use** — For T0/T1, test function calling if required
5. **Benchmark code quality** — Review generated code for correctness

### Success Rate Targets

| Tier | Cloud Models | Local Models |
|------|--------------|--------------|
| T0 | ≥ 95% | N/A (rarely used) |
| T1 | ≥ 90% | ≥ 75% |
| T2 | ≥ 90% | ≥ 80% |
| T3 | ≥ 95% | ≥ 85% |

---

## References

- **Feature Spec:** `tools/hw_checker/docs/specs/feature_194/feature.md`
- **Idea Draft:** `docs/drafts/idea-model-agnostic-ws-protocol.md`
- **SDP Protocol:** `sdp/PROTOCOL.md`
- **WS Template:** `tools/hw_checker/docs/workstreams/TEMPLATE.md`

---

**Last Updated:** 2026-01-21  
**Version:** 1.0  
**Maintainer:** SDP Team
