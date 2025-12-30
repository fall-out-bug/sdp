# Model Recommendations for Consensus Agents

This guide provides recommendations for selecting AI models for each agent role in the consensus workflow.

**Note**: AI model capabilities change rapidly. Check official provider documentation for current offerings and benchmarks.

## Quick Role Assignments

| Role | Recommended | Why |
|------|-------------|-----|
| **Analyst** | Most capable model available | Requirements need deep understanding |
| **Architect** | Most capable model available | Architecture decisions are critical |
| **Tech Lead** | Medium-tier model | Planning doesn't need highest capability |
| **Developer** | Fast model | Rapid iteration is key |
| **QA** | Fast model | Test verification is straightforward |
| **DevOps** | Fast model | Config generation is routine |
| **Security** | Most capable model available | Security requires careful analysis |

## General Principles

1. **Strategic roles** (Analyst, Architect, Security) benefit from more capable models
2. **Implementation roles** (Developer, QA, DevOps) work well with faster, cheaper models
3. **Start cheap, escalate if needed** - try faster model first, use capable model only when stuck

## For Claude Code Users

Claude Code supports Claude models only. Use `/model` command to switch:

```
/model opus    # Most capable - for Analyst, Architect, Security
/model sonnet  # Balanced - for Tech Lead
/model haiku   # Fastest - for Developer, QA, DevOps
```

### Recommended Assignment

| Role | Model | Command |
|------|-------|---------|
| Analyst | Opus | `/model opus` |
| Architect | Opus | `/model opus` |
| Tech Lead | Sonnet | `/model sonnet` |
| Developer | Haiku or Sonnet | `/model haiku` |
| QA | Haiku | `/model haiku` |
| DevOps | Haiku | `/model haiku` |
| Security | Opus | `/model opus` |

## For Cursor Users

Cursor supports multiple providers. Configure in Settings → Models.

### Suggested Strategy

- **Strategic decisions**: Use most capable model (e.g., Claude Opus, GPT-4)
- **Implementation**: Use fast model (e.g., Claude Haiku, fast variants)
- **Check Cursor's current offerings** - they change frequently

## Cost Optimization

### Strategy 1: Quality-First
Use capable models for all roles. Higher cost, best results.

### Strategy 2: Balanced (Recommended)
- Strategic roles (Analyst, Architect, Security): Capable model
- Implementation roles (Developer, QA, DevOps): Fast model

### Strategy 3: Budget
Use fast models for all roles. Lowest cost, adequate for many projects.

## Model Selection Tips

1. **Match model to task complexity**
   - Simple tasks → Fast model
   - Complex reasoning → Capable model

2. **Escalate when stuck**
   - Start with fast model
   - If 2+ iterations without progress, switch to capable model

3. **Strategic decisions matter most**
   - Analyst and Architect set the foundation
   - Don't skimp on these roles

4. **Implementation is repetitive**
   - Developer, QA, DevOps do similar tasks repeatedly
   - Fast models work well here

## Current Model Landscape

Check these sources for current benchmarks:
- [SWE-bench](https://www.swebench.com/) - Coding benchmark leaderboard
- [Anthropic Models](https://docs.anthropic.com/claude/docs/models-overview) - Claude model specifications
- [OpenAI Models](https://platform.openai.com/docs/models) - GPT model specifications
- [Google AI](https://ai.google.dev/models) - Gemini model specifications

---

**Note**: This guide focuses on principles rather than specific model versions, as the AI landscape evolves rapidly. Always verify current model capabilities with official documentation.
