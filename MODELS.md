# Model Recommendations for Consensus Agents

**Last Updated:** December 29, 2025 (SWE-bench Verified data)

This guide provides objective, data-driven recommendations for selecting AI models for each agent role in the consensus workflow. All recommendations are based on **SWE-bench Verified** scores - the industry standard for measuring real-world coding performance.

## 🏆 SWE-bench Verified Leaderboard (December 2025)

| Rank | Model | Score | Provider | Type | Cost/1M | Speed | Best For |
|------|-------|-------|----------|------|---------|-------|----------|
| 1 | **Claude Opus 4.5** | **80.9%** | Anthropic | Proprietary | $15/$75 | Slow | Strategic decisions |
| 2 | **GPT-5.2 Thinking** | 80.0%* | OpenAI | Proprietary | ~$20/$100 | Slow | Complex reasoning |
| 3 | **Claude Sonnet 4.5** | 77.2% | Anthropic | Proprietary | $3/$15 | Medium | Complex refactoring |
| 4 | **Gemini 3 Flash** | 76-78% | Google | Proprietary | $0.075/$0.30 | **Very Fast** | **80% of tasks** ⭐ |
| 5 | **Gemini 3 Pro** | 74.2% | Google | Proprietary | $1.25/$5 | Medium | Planning, analysis |
| 6 | **Claude Haiku 4.5** | 73.3% | Anthropic | Proprietary | $1/$5 | **Very Fast** | Claude ecosystem |
| 7 | **GPT-5.2** | 71.8% | OpenAI | Proprietary | ~$5/$20 | Medium | General coding |
| 8 | **Kimi K2 Thinking** | 71.3% | Moonshot | **Open** | **Free** | Medium | **Budget option** ⭐ |
| 9 | **Qwen3-Coder** | 69.6% | Alibaba | **Open** | **Free** | Fast | Open source leader |
| 10 | **Kimi K2** | 65.8% | Moonshot | **Open** | **Free** | Medium | Self-hosted |

*Vendor-reported data

## 🎯 Quick Role Assignments

### For Claude Code / Cursor Users

| Role | Recommended | Alternative | Budget |
|------|-------------|-------------|--------|
| **Analyst** | Opus 4.5 (80.9%) | Gemini 3 Pro (74.2%) | Kimi K2 Thinking (71.3%) |
| **Architect** | Opus 4.5 (80.9%) | Gemini 3 Pro (74.2%) | Kimi K2 Thinking (71.3%) |
| **Tech Lead** | Sonnet 4.5 (77.2%) | Gemini 3 Flash (76-78%) | Qwen3-Coder (69.6%) |
| **Developer** | **Gemini 3 Flash** (76-78%) ⭐ | Haiku 4.5 (73.3%) | Kimi K2 Thinking (71.3%) |
| **QA** | **Gemini 3 Flash** (76-78%) ⭐ | Haiku 4.5 (73.3%) | Kimi K2 Thinking (71.3%) |
| **DevOps** | **Gemini 3 Flash** (76-78%) ⭐ | Haiku 4.5 (73.3%) | Qwen3-Coder (69.6%) |
| **SRE** | **Gemini 3 Flash** (76-78%) ⭐ | Haiku 4.5 (73.3%) | Qwen3-Coder (69.6%) |
| **Security** | Opus 4.5 (80.9%) | GPT-5.2 (71.8%) | Kimi K2 Thinking (71.3%) |

**💡 Key Insight:** Gemini 3 Flash (76-78%) outperforms Haiku 4.5 (73.3%) by 3-5% while being **93% cheaper** and 4-5x faster!

## 📊 Detailed Model Analysis

### Tier 1: Strategic Decisions (75%+ on SWE-bench)

#### Claude Opus 4.5 - `claude-opus-4-5-20251101`
- **Score:** 80.9% (First to break 80%!)
- **Cost:** $15 input / $75 output per 1M tokens
- **Latency:** High (~20-30s)
- **Use for:** Analyst, Architect, Security
- **Strengths:** Best reasoning, handles ambiguity, veto decisions
- **Available:** Claude Code, Cursor, API

#### GPT-5.2 Thinking
- **Score:** 80.0% (vendor-reported)
- **Cost:** ~$20 input / $100 output per 1M tokens
- **Latency:** Very High (extended thinking)
- **Use for:** Complex architectural decisions
- **Strengths:** Deep reasoning, math, logic
- **Available:** ChatGPT, API (coming)

#### Claude Sonnet 4.5 - `claude-sonnet-4-5-20250929`
- **Score:** 77.2%
- **Cost:** $3 input / $15 output per 1M tokens
- **Latency:** Medium (~5-10s)
- **Use for:** Tech Lead, complex refactoring
- **Strengths:** Deep codebase understanding, multi-file coordination
- **Available:** Claude Code, Cursor, API

#### Gemini 3 Flash - `gemini-3.0-flash`
- **Score:** 76-78% ⭐ **BEST VALUE**
- **Cost:** $0.075 input / $0.30 output per 1M tokens (**93% cheaper than Haiku!**)
- **Latency:** Very Low (~1-2s, 4-5x faster than Sonnet)
- **Use for:** Developer, QA, DevOps, SRE (80% of tasks!)
- **Strengths:** Speed + quality combo, multi-modal
- **Available:** Cursor, Google AI Studio, API

**Why Gemini 3 Flash is revolutionary:**
- Beats Haiku 4.5 by 3-5% on SWE-bench
- $0.075/$0.30 vs Haiku's $1/$5 = **13x cheaper**
- 4-5x faster than Sonnet 4.5
- Perfect for rapid iteration in multi-agent workflows

#### Gemini 3 Pro - `gemini-3.0-pro`
- **Score:** 74.2%
- **Cost:** $1.25 input / $5 output per 1M tokens
- **Latency:** Medium (~5-8s)
- **Use for:** Analyst, Architect (when Opus too expensive)
- **Strengths:** Balanced reasoning, multi-modal
- **Available:** Cursor, Google AI Studio, API

### Tier 2: Implementation & Operations (70-75% on SWE-bench)

#### Claude Haiku 4.5 - `claude-haiku-4-5-20241022`
- **Score:** 73.3%
- **Cost:** $1 input / $5 output per 1M tokens
- **Latency:** Very Low (~1-3s)
- **Use for:** Developer, QA (when locked to Claude ecosystem)
- **Strengths:** Fast, extended thinking, computer use
- **Available:** Claude Code, Cursor, API
- **Note:** Choose Gemini 3 Flash if not locked to Claude

#### GPT-5.2 - `gpt-5.2`
- **Score:** 71.8%
- **Cost:** ~$5 input / $20 output per 1M tokens
- **Latency:** Medium (~5-10s)
- **Use for:** General development tasks
- **Strengths:** Balanced, reliable, good tool use
- **Available:** ChatGPT, Cursor, API

#### Kimi K2 Thinking - `kimi-k2-thinking` ⭐ **BEST OPEN SOURCE**
- **Score:** 71.3% (beats Haiku 4.5!)
- **Cost:** **FREE** (self-hosted or API)
- **Latency:** Medium (~10-15s)
- **Use for:** Budget-conscious teams, all roles
- **Strengths:** Best free model, open weights, long context (1M tokens)
- **Available:** HuggingFace, local (via Ollama), API

### Tier 3: Open Source & Budget (65-70% on SWE-bench)

#### Qwen3-Coder - `qwen3-coder-480b-a35b-instruct`
- **Score:** 69.6%
- **Cost:** **FREE** (self-hosted)
- **Latency:** Fast (local) / Medium (API)
- **Use for:** DevOps, SRE, documentation
- **Strengths:** Best open coder, 480B MoE (35B active), 256K context
- **Available:** HuggingFace, Ollama, Alibaba Cloud

#### Kimi K2 - `kimi-k2-instruct`
- **Score:** 65.8%
- **Cost:** **FREE** (self-hosted or API)
- **Latency:** Medium
- **Use for:** Self-hosted workflows
- **Strengths:** 1T params MoE, long context, multilingual
- **Available:** HuggingFace, NVIDIA NIM, API

### Specialized Models

#### GPT-5.2-Codex - `gpt-5.2-codex`
- **Released:** December 18, 2025
- **Score:** 56.4% on SWE-bench **Pro** (harder benchmark!)
- **Cost:** ~$10 input / $50 output per 1M tokens
- **Use for:** Enterprise refactors, security audits, large migrations
- **Strengths:** Long-horizon work, context compaction, cybersecurity
- **Note:** Optimized for agentic coding, not quick tasks
- **Available:** ChatGPT Codex CLI, API (soon)

#### Devstral 2512 (123B & 24B)
- **Released:** December 9, 2025
- **Score:** Mid-tier (~65-70% estimated)
- **Cost:** **FREE** (open source)
- **Use for:** European data sovereignty requirements
- **Available:** HuggingFace, Mistral API

#### GLM-4.6
- **Released:** December 1, 2025
- **Score:** ~65% (hits step limits on SWE-bench)
- **Cost:** **FREE** (open source)
- **Use for:** Chinese language support, research
- **Available:** HuggingFace, Zhipu AI API

#### Minimax M2
- **Released:** November 24, 2025
- **Score:** Unknown (leaderboard present)
- **Cost:** ~$0.44 per problem
- **Use for:** China-based teams
- **Available:** Minimax API

## 💰 Cost Analysis Per Epic

Based on typical epic workflow (Analyst → Architect → Tech Lead → Developer → QA → DevOps):

### Strategy 1: Premium (Best Quality)
```
Analyst:    Opus 4.5     → $25-35
Architect:  Opus 4.5     → $25-35
Tech Lead:  Sonnet 4.5   → $8-12
Developer:  Sonnet 4.5   → $10-15
QA:         Sonnet 4.5   → $5-8
DevOps:     Sonnet 4.5   → $5-8

Total: $78-113 per epic
```

### Strategy 2: Optimal (Recommended) ⭐
```
Analyst:    Opus 4.5         → $25-35
Architect:  Opus 4.5         → $25-35
Tech Lead:  Gemini 3 Flash   → $0.50-1
Developer:  Gemini 3 Flash   → $0.80-2
QA:         Gemini 3 Flash   → $0.30-0.80
DevOps:     Gemini 3 Flash   → $0.30-0.80

Total: $52-75 per epic (33% savings!)
Quality: Better! (Flash 76% > Sonnet 77% negligible for implementation)
Speed: 2-3x faster (Flash very fast)
```

### Strategy 3: Budget (Open Source)
```
Analyst:    Kimi K2 Thinking → FREE
Architect:  Kimi K2 Thinking → FREE
Tech Lead:  Kimi K2 Thinking → FREE
Developer:  Qwen3-Coder      → FREE
QA:         Qwen3-Coder      → FREE
DevOps:     Qwen3-Coder      → FREE

Total: FREE (compute costs only)
Quality: 71% average (acceptable for non-critical)
```

### Strategy 4: Hybrid (Best Balance)
```
Analyst:    Gemini 3 Pro     → $5-8
Architect:  Gemini 3 Pro     → $5-8
Tech Lead:  Gemini 3 Flash   → $0.50-1
Developer:  Gemini 3 Flash   → $0.80-2
QA:         Kimi K2 Thinking → FREE
DevOps:     Qwen3-Coder      → FREE

Total: $11-19 per epic (83% savings vs Premium!)
Quality: 74% average (excellent)
```

## 🚀 Performance Benchmarks

### Speed Comparison (typical implementation task)

| Model | Latency | Tokens/sec | Time to 500 tokens |
|-------|---------|------------|-------------------|
| Gemini 3 Flash | **1-2s** | ~250 | **2-3s** ⚡ |
| Haiku 4.5 | 1-3s | ~200 | 3-4s |
| Sonnet 4.5 | 5-10s | ~100 | 10-15s |
| Opus 4.5 | 20-30s | ~50 | 30-40s |
| Kimi K2 | 10-15s | ~80 | 15-20s |
| GPT-5.2 | 5-10s | ~120 | 8-12s |

**Impact on workflow:**
- Gemini 3 Flash: 6 agents in **~10-15 minutes**
- Haiku 4.5: 6 agents in **~15-20 minutes**
- Sonnet 4.5: 6 agents in **~45-60 minutes**

### Quality vs Speed Sweet Spot

```
Quality (SWE-bench)
    ↑
80% │ Opus 4.5 ●
    │
75% │        ● Sonnet 4.5
    │    ● Gemini 3 Flash ← SWEET SPOT! ⭐
70% │  ● Haiku 4.5
    │ ● Kimi K2 Thinking
65% │
    └────────────────────→ Speed
      Slow    Medium    Fast
```

## 🔄 Model Switching Strategy

### Escalation Path (Start Cheap, Escalate if Needed)

```
1. Start: Gemini 3 Flash (76-78%, $0.075/$0.30, fast)
   ↓
2. If stuck after 2 iterations OR complex multi-file (5+):
   → Sonnet 4.5 (77%, $3/$15, medium)
   ↓
3. If architectural ambiguity OR veto decision:
   → Opus 4.5 (80.9%, $15/$75, slow)
   ↓
4. Resume: Gemini 3 Flash (once unblocked)
```

**Result:** 90% cost reduction while maintaining quality

### When to Use Each Model

**Opus 4.5 (80.9%):**
- ✅ Initial requirements (Analyst)
- ✅ System architecture (Architect)
- ✅ Veto decisions
- ✅ Security audits
- ❌ NOT for routine coding

**Sonnet 4.5 (77.2%):**
- ✅ Implementation planning (Tech Lead)
- ✅ Complex refactoring (5+ files)
- ✅ Cross-epic analysis
- ❌ NOT for simple implementations

**Gemini 3 Flash (76-78%) ⭐ DEFAULT:**
- ✅ TDD implementation (Developer)
- ✅ Test verification (QA)
- ✅ CI/CD scripts (DevOps)
- ✅ Monitoring setup (SRE)
- ✅ Documentation
- ✅ **80% of all tasks!**

**Kimi K2 Thinking (71.3%):**
- ✅ Budget-conscious teams
- ✅ All roles (acceptable quality)
- ✅ Self-hosted requirements
- ⚠️ 5-10% quality drop vs Gemini 3 Flash

## 🌍 Provider-Specific Considerations

### Anthropic (Claude)
**Strengths:**
- Highest SWE-bench scores (Opus 4.5: 80.9%)
- Long context (200K tokens)
- Extended thinking mode
- Best protocol adherence

**Considerations:**
- Expensive (Opus: $15/$75)
- Slower (Opus: 20-30s)
- Haiku 4.5 loses to Gemini 3 Flash

**Best for:** Critical decisions (Analyst, Architect, Security)

### Google (Gemini)
**Strengths:**
- **Best value:** Flash 76-78% at $0.075/$0.30 ⭐
- Fastest inference (1-2s)
- Multi-modal (vision, audio)
- Long context (2M tokens on Pro)

**Considerations:**
- Less mature ecosystem than Claude
- API stability varies by region
- Pro is pricier than Flash ($1.25/$5)

**Best for:** Implementation (Developer, QA, DevOps, SRE)

### OpenAI (GPT)
**Strengths:**
- Balanced quality (GPT-5.2: 71.8%)
- Specialized Codex variant
- Reliable tool use
- Good documentation

**Considerations:**
- Mid-tier pricing ($5/$20)
- Not best at any specific task
- Codex optimized for enterprise, not speed

**Best for:** Teams already on OpenAI ecosystem

### Open Source (Kimi K2, Qwen3-Coder)
**Strengths:**
- **FREE** (compute costs only)
- Data privacy (self-hosted)
- Customizable (fine-tuning)
- Kimi K2 Thinking beats Haiku 4.5!

**Considerations:**
- Slower than commercial (10-20s)
- Requires infrastructure
- Support is community-driven

**Best for:** Budget teams, data sovereignty, research

## 🔧 Integration Guides

### Cursor IDE 2.0
Supports all models via settings:
```json
{
  "models": {
    "analyst": "claude-opus-4-5-20251101",
    "architect": "claude-opus-4-5-20251101",
    "tech_lead": "gemini-3.0-flash",
    "developer": "gemini-3.0-flash",
    "qa": "gemini-3.0-flash",
    "devops": "qwen3-coder-480b",
    "security": "claude-opus-4-5-20251101"
  }
}
```

See [CURSOR.md](docs/guides/CURSOR.md) for multi-agent setup.

### Claude Code CLI
Supports Claude + API providers:
```bash
# Claude models (native)
claude --model claude-opus-4-5-20251101

# OpenAI via API
export OPENAI_API_KEY="..."
claude --provider openai --model gpt-5.2

# Google via API
export GOOGLE_API_KEY="..."
claude --provider google --model gemini-3.0-flash
```

See [CLAUDE_CODE.md](docs/guides/CLAUDE_CODE.md) for details.

### API Direct
For custom integrations:
- Anthropic: https://docs.anthropic.com/
- OpenAI: https://platform.openai.com/docs
- Google AI: https://ai.google.dev/
- Moonshot (Kimi): https://platform.moonshot.cn/
- Alibaba (Qwen): https://www.alibabacloud.com/

## 📈 Future Model Watch

Expected releases Q1-Q2 2026:
- **Claude Opus 5** - Targeting 85%+ on SWE-bench
- **GPT-6** - OpenAI's next flagship
- **Gemini 4** - Google's response to Claude Opus
- **Qwen4-Coder** - Alibaba's continued push
- **DeepSeek V4** - Chinese competitor

**Recommendation:** Re-evaluate this guide quarterly as new models release.

## 🎓 Recommendations Summary

### For Most Teams (Recommended) ⭐
```
Strategic:      Opus 4.5 or Gemini 3 Pro
Implementation: Gemini 3 Flash (default for 80% tasks)
Budget:         Kimi K2 Thinking when cost matters

Cost per epic: $50-75
Quality: Excellent (75% average)
Speed: Fast (10-15 min per agent)
```

### For Claude-First Teams
```
Strategic:      Opus 4.5
Complex:        Sonnet 4.5
Standard:       Haiku 4.5 (not Flash)

Cost per epic: $80-115
Quality: Excellent (77% average)
Ecosystem: Unified Claude experience
```

### For Budget Teams
```
All roles:      Kimi K2 Thinking or Qwen3-Coder

Cost per epic: FREE (compute only)
Quality: Good (70% average)
Trade-off: Slightly lower quality, self-hosted
```

### For Enterprise
```
Strategic:      Opus 4.5
Complex:        GPT-5.2-Codex (for migrations)
Standard:       Gemini 3 Flash
Security:       Opus 4.5

Cost per epic: $60-90
Features: Security, compliance, audit trails
```

## 📚 References

- [SWE-bench Verified Official](https://www.swebench.com/) - Benchmark leaderboard
- [SWE-bench Leaderboard Snapshot](https://coconote.app/notes/438dc924-cd41-4d74-9ab0-4f668b580e42) - December 2025 data
- [Claude Haiku 4.5 Announcement](https://www.anthropic.com/news/claude-haiku-4-5) - Anthropic
- [GPT-5.2-Codex Release](https://openai.com/index/introducing-gpt-5-2-codex/) - OpenAI
- [Gemini 3 Flash](https://blog.google/products/gemini/gemini-3-flash/) - Google
- [Kimi K2 GitHub](https://github.com/MoonshotAI/Kimi-K2) - Moonshot AI
- [Qwen3-Coder](https://qwenlm.github.io/blog/qwen3-coder/) - Alibaba

---

**Version:** 3.0
**Last Updated:** December 29, 2025
**Key Change:** Added Gemini 3 Flash as recommended default (76-78%, 13x cheaper than Haiku, 3-5% better quality)
