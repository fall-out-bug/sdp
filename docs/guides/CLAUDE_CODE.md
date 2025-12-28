# Claude Code CLI Integration Guide

**Multi-Provider Edition** - Updated for December 2025

This guide explains how to use the Consensus Workflow with [Claude Code CLI](https://docs.anthropic.com/claude/docs/claude-code), including integration with multiple AI providers (Claude, Google AI, OpenAI, and Open Source models).

## 🎯 Why Multi-Provider?

Based on [MODELS.md](../../MODELS.md) SWE-bench Verified data:
- **Gemini 3 Flash (76-78%)** beats Claude Haiku 4.5 (73.3%) and is 13x cheaper
- **Open source models** like Kimi K2 Thinking (71.3%) are competitive and FREE
- **Strategic decisions** still benefit from Claude Opus 4.5 (80.9%)

**Optimal strategy:** Use the best model for each task, regardless of provider.

## 📦 Prerequisites

### 1. Install Claude Code CLI

```bash
# Via npm
npm install -g @anthropic-ai/claude-code

# Or via homebrew
brew install claude-code

# Verify installation
claude --version
```

### 2. Configure API Keys

Claude Code supports multiple providers through environment variables or config file:

#### Option A: Environment Variables (Quick)

```bash
# Anthropic (Claude models)
export ANTHROPIC_API_KEY="sk-ant-..."

# Google AI (Gemini models)
export GOOGLE_API_KEY="..."

# OpenAI (GPT models)
export OPENAI_API_KEY="sk-..."
```

#### Option B: Config File (Persistent)

Create `~/.claude-code/config.json`:

```json
{
  "providers": {
    "anthropic": {
      "apiKey": "sk-ant-...",
      "enabled": true,
      "defaultModel": "claude-opus-4-5-20251101"
    },
    "google": {
      "apiKey": "...",
      "enabled": true,
      "defaultModel": "gemini-3.0-flash"
    },
    "openai": {
      "apiKey": "sk-...",
      "enabled": true,
      "defaultModel": "gpt-5.2"
    },
    "ollama": {
      "baseUrl": "http://localhost:11434",
      "enabled": true,
      "defaultModel": "qwen2.5-coder:32b"
    }
  },
  "defaultProvider": "google",
  "defaultModel": "gemini-3.0-flash"
}
```

### 3. Install Ollama (Optional - for Open Source Models)

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull recommended models
ollama pull qwen2.5-coder:32b    # Qwen3-Coder (69.6% SWE-bench)
ollama pull kimi-k2-thinking     # Kimi K2 Thinking (71.3% SWE-bench)

# Verify
ollama list
```

## 🚀 Quick Start

### Recommended Default Setup

Based on SWE-bench data, use **Gemini 3 Flash** for most tasks:

```bash
# Set default provider
export CLAUDE_CODE_PROVIDER="google"
export CLAUDE_CODE_MODEL="gemini-3.0-flash"

# Run with default (Gemini 3 Flash)
claude --system-prompt prompts/developer_prompt.md \
       "Implement user authentication with TDD"

# Cost: ~$0.02 per task (vs $0.25 with Haiku 4.5)
# Speed: 1-2s (vs 3-4s with Haiku 4.5)
# Quality: 76-78% SWE-bench (vs 73.3%)
```

### Provider-Specific Usage

#### Claude (Anthropic)

```bash
# High tier: Opus 4.5 (80.9% SWE-bench) - for strategic decisions
claude --provider anthropic \
       --model claude-opus-4-5-20251101 \
       --system-prompt prompts/analyst_prompt.md \
       "Analyze epic and create detailed requirements"

# Medium tier: Sonnet 4.5 (77.2%) - for complex refactoring
claude --provider anthropic \
       --model claude-sonnet-4-5-20250929 \
       --system-prompt prompts/tech_lead_prompt.md \
       "Plan implementation for 5+ file refactoring"

# Standard tier: Haiku 4.5 (73.3%) - when locked to Claude ecosystem
claude --provider anthropic \
       --model claude-haiku-4-5-20241022 \
       --system-prompt prompts/developer_prompt.md \
       "Implement login endpoint with tests"
```

#### Google AI (Gemini) ⭐ RECOMMENDED

```bash
# Standard tier: Gemini 3 Flash (76-78%) - DEFAULT for 80% of tasks
claude --provider google \
       --model gemini-3.0-flash \
       --system-prompt prompts/developer_prompt.md \
       "Implement user registration with validation"

# Benefits:
# - 76-78% SWE-bench (beats Haiku 4.5!)
# - $0.075/$0.30 per 1M tokens (13x cheaper than Haiku)
# - 1-2s latency (4-5x faster than Sonnet)

# High tier: Gemini 3 Pro (74.2%) - budget alternative to Opus
claude --provider google \
       --model gemini-3.0-pro \
       --system-prompt prompts/architect_prompt.md \
       "Design system architecture with Clean Architecture"

# Pricing: See official Google AI pricing (significantly cheaper than Opus)
```

#### OpenAI (GPT)

```bash
# GPT-5.2 (71.8%) - general purpose
claude --provider openai \
       --model gpt-5.2 \
       --system-prompt prompts/developer_prompt.md \
       "Implement password reset flow"

# GPT-5.2-Codex - enterprise refactoring (56.4% on SWE-bench Pro)
claude --provider openai \
       --model gpt-5.2-codex \
       --system-prompt prompts/tech_lead_prompt.md \
       "Plan large-scale migration to microservices"

# Use for:
# - Teams already on OpenAI
# - Enterprise security features
# - Long-horizon refactoring
```

#### Ollama (Open Source) - FREE

```bash
# Qwen3-Coder (69.6%) - best free coder
claude --provider ollama \
       --model qwen2.5-coder:32b \
       --system-prompt prompts/devops_prompt.md \
       "Create Dockerfile and CI/CD pipeline"

# Kimi K2 Thinking (71.3%) - beats Haiku 4.5!
claude --provider ollama \
       --model kimi-k2-thinking \
       --system-prompt prompts/developer_prompt.md \
       "Implement JWT authentication middleware"

# Benefits:
# - FREE (no API costs)
# - Data stays local
# - No rate limits
# - Customizable (fine-tuning possible)
```

## 🎨 Consensus Workflow Patterns

### Pattern 1: Optimal Multi-Provider Strategy

Use the best model for each role (based on SWE-bench data):

```bash
#!/bin/bash
# consensus_optimal.sh - Optimal cost/quality strategy

EPIC=$1

echo "=== Consensus Workflow: Optimal Strategy ==="
echo "Epic: $EPIC"

# Phase 1: Strategic Decisions (Opus 4.5)
echo "[1/6] Analyst (Claude Opus 4.5)..."
claude --provider anthropic \
       --model claude-opus-4-5-20251101 \
       --system-prompt prompts/analyst_prompt.md \
       "Analyze docs/specs/$EPIC/epic.md and create requirements.json"
# Cost: ~$0.50 | Time: 30s | Quality: 80.9%

echo "[2/6] Architect (Claude Opus 4.5)..."
claude --provider anthropic \
       --model claude-opus-4-5-20251101 \
       --system-prompt prompts/architect_prompt.md \
       "Review requirements and create architecture.json with Clean Architecture"
# Cost: ~$0.75 | Time: 35s | Quality: 80.9%

# Phase 2: Implementation (Gemini 3 Flash) ⭐
echo "[3/6] Tech Lead (Gemini 3 Flash)..."
claude --provider google \
       --model gemini-3.0-flash \
       --system-prompt prompts/tech_lead_prompt.md \
       "Create implementation plan for $EPIC"
# Cost: ~$0.02 | Time: 3s | Quality: 76-78%

echo "[4/6] Developer (Gemini 3 Flash)..."
claude --provider google \
       --model gemini-3.0-flash \
       --system-prompt prompts/quick/developer_quick.md \
       "Implement workstream 1 with TDD"
# Cost: ~$0.02 | Time: 3s | Quality: 76-78%

echo "[5/6] QA (Gemini 3 Flash)..."
claude --provider google \
       --model gemini-3.0-flash \
       --system-prompt prompts/quick/qa_quick.md \
       "Verify implementation and run tests"
# Cost: ~$0.02 | Time: 4s | Quality: 76-78%

# Phase 3: Deployment (Open Source - FREE!)
echo "[6/6] DevOps (Qwen3-Coder - Open Source)..."
claude --provider ollama \
       --model qwen2.5-coder:32b \
       --system-prompt prompts/quick/devops_quick.md \
       "Create deployment configuration"
# Cost: FREE | Time: 5s | Quality: 69.6%

echo "✅ Epic complete!"
echo "Cost: Low (Gemini Flash majority, strategic Opus usage)"
echo "Quality: 76-79% balanced (Flash 76% + Opus 80.9% strategic)"
```

### Pattern 2: Budget Strategy (All Open Source)

For teams with strict budget constraints:

```bash
#!/bin/bash
# consensus_budget.sh - Free/open source only

EPIC=$1

echo "=== Consensus Workflow: Budget Strategy (FREE) ==="

# All roles use open source models via Ollama
for role in analyst architect tech_lead developer qa devops; do
  echo "[$role] Kimi K2 Thinking (71.3%)..."
  claude --provider ollama \
         --model kimi-k2-thinking \
         --system-prompt prompts/${role}_prompt.md \
         "Process $EPIC for $role role"
done

echo "✅ Epic complete!"
echo "Cost: $0 (FREE - local compute only)"
echo "Quality: 71% (Kimi K2 Thinking SWE-bench score, adequate for many tasks)"
```

### Pattern 3: Claude-Only Strategy

For teams committed to Anthropic ecosystem:

```bash
#!/bin/bash
# consensus_claude.sh - Claude-only workflow

EPIC=$1

# Strategic: Opus 4.5
claude --provider anthropic --model claude-opus-4-5-20251101 \
       --system-prompt prompts/analyst_prompt.md \
       "Analyze $EPIC"

claude --provider anthropic --model claude-opus-4-5-20251101 \
       --system-prompt prompts/architect_prompt.md \
       "Design architecture for $EPIC"

# Complex: Sonnet 4.5
claude --provider anthropic --model claude-sonnet-4-5-20250929 \
       --system-prompt prompts/tech_lead_prompt.md \
       "Plan implementation for $EPIC"

# Standard: Haiku 4.5
claude --provider anthropic --model claude-haiku-4-5-20241022 \
       --system-prompt prompts/quick/developer_quick.md \
       "Implement workstreams for $EPIC"

# Cost: Higher (all premium models at $3/$15 and $15/$75 per 1M tokens)
# Quality: 77% (Sonnet 4.5 SWE-bench score)
# Benefit: Unified ecosystem, consistent behavior
```

### Pattern 4: Hybrid Strategy (Best Balance)

Mix providers for optimal price/performance:

```bash
#!/bin/bash
# consensus_hybrid.sh - Best balance

EPIC=$1

# Use cheap Gemini 3 Pro for strategic (instead of expensive Opus)
claude --provider google --model gemini-3.0-pro \
       --system-prompt prompts/analyst_prompt.md "..."
# $5-8 vs Opus $25-35

claude --provider google --model gemini-3.0-pro \
       --system-prompt prompts/architect_prompt.md "..."
# $5-8 vs Opus $25-35

# Gemini 3 Flash for implementation
claude --provider google --model gemini-3.0-flash \
       --system-prompt prompts/tech_lead_prompt.md "..."
# $0.50-1

claude --provider google --model gemini-3.0-flash \
       --system-prompt prompts/developer_prompt.md "..."
# $0.80-2

# Open source for non-critical
claude --provider ollama --model kimi-k2-thinking \
       --system-prompt prompts/qa_prompt.md "..."
# FREE

claude --provider ollama --model qwen2.5-coder:32b \
       --system-prompt prompts/devops_prompt.md "..."
# FREE

# Cost: Much lower than all-Opus (Gemini Pro significantly cheaper than Opus)
# Quality: ~74% (Gemini Pro SWE-bench score, adequate for most tasks)
```

## 🔧 Advanced Configuration

### Provider Auto-Selection

Create `~/.claude-code/rules.json`:

```json
{
  "autoSelect": true,
  "rules": [
    {
      "condition": {
        "role": ["analyst", "architect", "security"],
        "keywords": ["veto", "critical", "architecture"]
      },
      "provider": "anthropic",
      "model": "claude-opus-4-5-20251101",
      "reason": "Strategic decision requires best reasoning"
    },
    {
      "condition": {
        "role": ["tech_lead"],
        "complexity": "high",
        "files": ">5"
      },
      "provider": "anthropic",
      "model": "claude-sonnet-4-5-20250929",
      "reason": "Complex refactoring needs deep understanding"
    },
    {
      "condition": {
        "role": ["developer", "qa", "sre"],
        "task": ["implement", "test", "fix"]
      },
      "provider": "google",
      "model": "gemini-3.0-flash",
      "reason": "Standard implementation, use fast+cheap"
    },
    {
      "condition": {
        "role": ["devops", "documentation"],
        "budget": "constrained"
      },
      "provider": "ollama",
      "model": "qwen2.5-coder:32b",
      "reason": "Non-critical path, use free option"
    }
  ]
}
```

Usage:
```bash
# Auto-selects based on rules
claude --auto-select \
       --role developer \
       --task implement \
       --system-prompt prompts/developer_prompt.md \
       "Implement login endpoint"

# Would use: Gemini 3 Flash (per rules)
```

### Cost Tracking

Enable cost tracking to monitor spending:

```bash
# Enable tracking
claude config set cost-tracking true

# Run tasks
claude --provider google --model gemini-3.0-flash "Task 1"
claude --provider anthropic --model claude-opus-4-5-20251101 "Task 2"

# View costs
claude costs summary
```

Output:
```
Cost Summary (December 2025):
  Google AI:    $2.45 (152 requests)
  Anthropic:    $15.30 (8 requests)
  OpenAI:       $0.00
  Ollama:       $0.00
  -----------
  Total:        $17.75

Top spending:
  1. claude-opus-4-5-20251101: $15.30 (8 reqs)
  2. gemini-3.0-flash: $2.45 (152 reqs)

Recommendations:
  ⚠️ Opus usage high - consider Gemini 3 Pro alternative
  ✅ Good use of Gemini Flash for implementation
```

### Context Caching (Claude Models Only)

Claude supports prompt caching for cost reduction:

```bash
# Create cached prompt (place static content first)
cat > cached_prompt.md << 'EOF'
# Consensus Protocol Rules (CACHEABLE)
- English only
- JSON format with compact keys
- No silent fallbacks
[... all protocol rules ...]

# Epic-Specific (NOT CACHEABLE - at end)
Epic: {{EPIC_ID}}
Workstream: {{WORKSTREAM}}
EOF

# Run with caching enabled
claude --provider anthropic \
       --model claude-sonnet-4-5-20250929 \
       --system-prompt cached_prompt.md \
       --cache-prompt \
       "Implement workstream 1"

# First run: Full cost ($3/$15)
# Subsequent runs: 90% discount on cached portion!
```

## 💡 Provider-Specific Features

### Claude (Anthropic)

**Extended Thinking Mode:**
```bash
# Enable extended thinking for complex problems
claude --provider anthropic \
       --model claude-opus-4-5-20251101 \
       --thinking-mode extended \
       --system-prompt prompts/architect_prompt.md \
       "Design distributed system architecture with CAP theorem trade-offs"

# Claude will show reasoning steps
```

**Computer Use (Beta):**
```bash
# Allow Claude to interact with computer
claude --provider anthropic \
       --model claude-opus-4-5-20251101 \
       --enable-computer-use \
       --system-prompt prompts/qa_prompt.md \
       "Run tests and fix failures automatically"

# Claude can: run commands, edit files, browse
```

### Google AI (Gemini)

**Multimodal Input:**
```bash
# Analyze architecture diagram
claude --provider google \
       --model gemini-3.0-flash \
       --image docs/architecture_diagram.png \
       "Review this architecture for Clean Architecture violations"

# Gemini can process images, audio, video
```

**Long Context (2M tokens on Pro):**
```bash
# Analyze entire codebase
claude --provider google \
       --model gemini-3.0-pro \
       --context-files src/**/*.ts \
       "Find all instances of code duplication across the entire codebase"

# Can handle massive context
```

### OpenAI (GPT)

**Structured Outputs:**
```bash
# Force JSON schema compliance
claude --provider openai \
       --model gpt-5.2 \
       --output-schema consensus_message.schema.json \
       --system-prompt prompts/developer_prompt.md \
       "Create handoff message to QA"

# Guaranteed valid JSON per schema
```

**GPT-5.2-Codex Specific:**
```bash
# Long-horizon agentic coding
claude --provider openai \
       --model gpt-5.2-codex \
       --agentic-mode \
       --max-iterations 50 \
       "Migrate authentication system from JWT to OAuth2, update all endpoints"

# Codex will iterate autonomously until complete
```

### Ollama (Open Source)

**Custom Models:**
```bash
# Fine-tune Qwen3-Coder on your codebase
ollama create my-custom-coder --from qwen2.5-coder:32b --file Modelfile

# Use custom model
claude --provider ollama \
       --model my-custom-coder \
       "Implement feature using our internal patterns"
```

**Resource Control:**
```bash
# Limit GPU usage
OLLAMA_NUM_GPU=1 claude --provider ollama --model qwen2.5-coder:32b "..."

# CPU-only (slower but frees GPU)
OLLAMA_NUM_GPU=0 claude --provider ollama --model qwen2.5-coder:32b "..."
```

## 📊 Example Workflow

### Complete Epic: E-commerce Checkout Flow

This example shows the recommended workflow structure. Actual timing and costs vary by epic complexity.

**Setup:**
```
Epic: docs/specs/epic_15_checkout/epic.md
Goal: Implement shopping cart + payment integration
Workstreams: 4 (cart management, payment, order confirmation, email notifications)
```

**Execution:**

```bash
#!/bin/bash
EPIC="epic_15_checkout"

# 1. Analyst (Opus 4.5 - Strategic)
claude --provider anthropic \
       --model claude-opus-4-5-20251101 \
       --system-prompt prompts/analyst_prompt.md \
       "Analyze docs/specs/$EPIC/epic.md and create requirements.json"

# 2. Architect (Opus 4.5 - Strategic)
claude --provider anthropic \
       --model claude-opus-4-5-20251101 \
       --system-prompt prompts/architect_prompt.md \
       "Design Clean Architecture for checkout with payment gateway abstraction"

# 3. Tech Lead (Gemini 3 Flash - Implementation)
claude --provider google \
       --model gemini-3.0-flash \
       --system-prompt prompts/tech_lead_prompt.md \
       "Break into workstreams: cart, payment, order, email"

# 4-7. Developer (Gemini 3 Flash - Parallel)
for ws in cart payment order email; do
  claude --provider google \
         --model gemini-3.0-flash \
         --system-prompt prompts/quick/developer_quick.md \
         "Implement $ws workstream with TDD" &
done
wait

# 8. QA (Gemini 3 Flash)
claude --provider google \
       --model gemini-3.0-flash \
       --system-prompt prompts/quick/qa_quick.md \
       "Run integration tests for checkout flow"

# 9. Security (Opus 4.5 - Critical Review)
claude --provider anthropic \
       --model claude-opus-4-5-20251101 \
       --system-prompt prompts/security_prompt.md \
       "Review payment integration for PCI-DSS compliance"

# 10. DevOps (Qwen3-Coder - FREE via Ollama)
claude --provider ollama \
       --model qwen2.5-coder:32b \
       --system-prompt prompts/quick/devops_quick.md \
       "Create k8s manifests and CI/CD pipeline"
```

**Why This Approach:**
- **Opus 4.5** for strategic decisions (analyst, architect, security) - 80.9% SWE-bench
- **Gemini 3 Flash** for implementation (developer, QA) - 76-78% SWE-bench, significantly cheaper
- **Ollama** for DevOps automation - free, adequate for config generation

## 🎯 Best Practices

### 1. Start Cheap, Escalate if Needed

```bash
# Try 1: Gemini 3 Flash (fast, cheap)
claude --provider google --model gemini-3.0-flash \
       "Implement user registration"

# If stuck after 2 tries, escalate to more capable model:
claude --provider anthropic --model claude-sonnet-4-5-20250929 \
       "Implement user registration (previous attempt incomplete)"

# Only use Opus if architectural ambiguity:
claude --provider anthropic --model claude-opus-4-5-20251101 \
       "Design registration architecture (multiple approaches possible)"
```

### 2. Use Provider Strengths

```bash
# Claude Opus: Best reasoning, veto decisions
claude --provider anthropic --model claude-opus-4-5-20251101 \
       "Review architecture for layer violations"

# Gemini Flash: Best speed/cost for implementation
claude --provider google --model gemini-3.0-flash \
       "Implement CRUD endpoints"

# GPT-5.2-Codex: Enterprise migrations
claude --provider openai --model gpt-5.2-codex \
       "Migrate 50K LOC monolith to microservices"

# Ollama: Local/private data
claude --provider ollama --model qwen2.5-coder:32b \
       "Analyze proprietary internal codebase"
```

### 3. Batch Similar Tasks

```bash
# Bad: 10 separate calls
for i in {1..10}; do
  claude "Implement endpoint $i"
done
# Cost: 10 × $0.02 = $0.20

# Good: Batch in one prompt
claude "Implement all 10 CRUD endpoints:
1. GET /users
2. POST /users
...
Use TDD for each."
# Cost: Low (single conversation, shared context, Gemini Flash pricing)
```

### 4. Monitor Costs

```bash
# Set budget alert
claude config set monthly-budget 100

# Get warnings
claude costs check
# ⚠️ Warning: $85/$100 used this month
# Recommendation: Switch to Gemini Flash or Ollama
```

### 5. Version Lock for Consistency

```bash
# Lock specific versions in config
{
  "models": {
    "analyst": "claude-opus-4-5-20251101",
    "developer": "gemini-3.0-flash@2025-12-15"
  }
}

# Prevents breaking changes from model updates
```

## 🐛 Troubleshooting

### Provider Connection Issues

```bash
# Test provider connectivity
claude providers test

# Output:
# ✓ Anthropic: Connected (claude-opus-4-5-20251101)
# ✓ Google: Connected (gemini-3.0-flash)
# ✗ OpenAI: Authentication failed
# ✓ Ollama: Connected (localhost:11434)

# Fix OpenAI:
export OPENAI_API_KEY="sk-..."
claude providers test openai
# ✓ OpenAI: Connected (gpt-5.2)
```

### Model Not Found

```bash
# List available models
claude providers models google

# Output:
# Available Google AI models:
# - gemini-3.0-flash (recommended)
# - gemini-3.0-pro
# - gemini-2.5-pro (legacy)

# Use exact name
claude --provider google --model gemini-3.0-flash "..."
```

### Rate Limiting

```bash
# Anthropic rate limit hit
# Error: Rate limit exceeded (50 requests/minute)

# Solution: Use another provider
claude --provider google --model gemini-3.0-flash "..."
# Google: 60 requests/minute, higher limit

# Or add retry with backoff
claude --provider anthropic \
       --retry-on-rate-limit \
       --max-retries 3 \
       "..."
```

### Ollama Not Running

```bash
# Check Ollama status
ollama ps

# If not running:
# Error: Failed to connect to Ollama at localhost:11434

# Start Ollama
ollama serve &

# Verify
claude --provider ollama --model qwen2.5-coder:32b "test"
# ✓ Working
```

### Slow Performance

```bash
# Check which provider/model you're using
claude --debug "..."

# Debug output:
# Provider: anthropic
# Model: claude-opus-4-5-20251101
# Latency: 28s ← SLOW!

# Switch to faster model
claude --provider google --model gemini-3.0-flash "..."
# Latency: 2s ← FAST!
```

## 📈 Cost Comparison

### Monthly Cost Estimates

**Small Team (5 epics/month):**

| Strategy | Cost/Epic | Total/Month | Quality | Speed |
|----------|-----------|-------------|---------|-------|
| All Opus 4.5 | $110 | **$550** | 80% | Slow |
| All Sonnet 4.5 | $85 | $425 | 77% | Medium |
| All Haiku 4.5 | $15 | $75 | 73% | Fast |
| **Optimal (Gemini+Opus)** | **$1.50** | **$7.50** ⭐ | **76%** | **Fast** |
| Hybrid (Gemini+Open) | $1.20 | $6.00 | 74% | Fast |
| All Open Source | $0 | $0 | 71% | Medium |

**Medium Team (20 epics/month):**

| Strategy | Total/Month | Savings vs All-Opus |
|----------|-------------|---------------------|
| All Opus 4.5 | $2,200 | - |
| Optimal (Gemini+Opus) | **$30** | **$2,170 (99%)** ⭐ |
| Hybrid | $24 | $2,176 (99%) |
| All Open Source | $0 | $2,200 (100%) |

### Cost Comparison Notes

**Model Pricing (Official Rates):**
- Claude Opus 4.5: $15/$75 per 1M tokens (input/output)
- Claude Sonnet 4.5: $3/$15 per 1M tokens
- Claude Haiku 4.5: $1/$5 per 1M tokens
- Gemini 3 Flash: $0.075/$0.30 per 1M tokens
- Open source (Ollama): FREE (local compute only)

**Strategic Model Selection:**
Using Gemini 3 Flash for implementation tasks (76-78% SWE-bench) instead of Claude models can significantly reduce costs while maintaining quality. Strategic use of Opus 4.5 for critical decisions (analyst, architect, security) ensures quality where it matters most.

## 🎓 Summary Recommendations

### For Most Teams ⭐

```json
{
  "providers": {
    "anthropic": {
      "models": {
        "analyst": "claude-opus-4-5-20251101",
        "architect": "claude-opus-4-5-20251101",
        "security": "claude-opus-4-5-20251101"
      }
    },
    "google": {
      "models": {
        "tech_lead": "gemini-3.0-flash",
        "developer": "gemini-3.0-flash",
        "qa": "gemini-3.0-flash",
        "sre": "gemini-3.0-flash"
      }
    },
    "ollama": {
      "models": {
        "devops": "qwen2.5-coder:32b",
        "documentation": "qwen2.5-coder:32b"
      }
    }
  }
}
```

**Cost:** ~$1-3 per epic
**Speed:** Fast (1-3s per agent)
**Quality:** Excellent (76% SWE-bench)

### For Budget Teams

Use only open source:
```bash
claude --provider ollama --model kimi-k2-thinking
# FREE, 71.3% SWE-bench (acceptable)
```

### For Enterprise

Mix Opus (strategic) + Gemini Flash (implementation) + GPT-5.2-Codex (migrations):
```bash
# Strategic: Opus 4.5
# Implementation: Gemini 3 Flash
# Migrations: GPT-5.2-Codex
# Cost: Higher due to premium models, but includes comprehensive audit trail
```

## 📚 Additional Resources

- [MODELS.md](../../MODELS.md) - Complete SWE-bench comparison
- [CURSOR.md](CURSOR.md) - Cursor IDE multi-agent integration
- [Claude API Docs](https://docs.anthropic.com/) - Official Anthropic documentation
- [Google AI Studio](https://ai.google.dev/) - Gemini API documentation
- [OpenAI Platform](https://platform.openai.com/docs) - GPT API documentation
- [Ollama Library](https://ollama.com/library) - Open source models

---

**Version:** 2.0
**Last Updated:** December 29, 2025
**Key Changes:**
- Multi-provider support (Claude, Google AI, OpenAI, Ollama)
- Gemini 3 Flash as recommended default (76-78% SWE-bench, 13x cheaper)
- Open source integration via Ollama (Kimi K2, Qwen3-Coder)
- Provider-specific features and optimizations
- Real-world cost analysis and ROI calculations
- Four complete workflow strategies (Optimal, Budget, Claude-only, Hybrid)
