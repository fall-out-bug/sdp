# OpenCode + GLM Quick Start

SDP works with OpenCode (Windsurf) and GLM models. This guide covers a minimal setup.

## 1. Install SDP

```bash
# In your project root
curl -sSL https://raw.githubusercontent.com/fall-out-bug/sdp/main/install.sh | SDP_IDE=opencode sh
```

This creates `.opencode/skills` → `sdp/prompts/skills` and `.opencode/agents` → `sdp/prompts/agents`.

## 2. OpenCode (Windsurf) Setup

- Install [Windsurf](https://codeium.com/windsurf) or OpenCode CLI
- SDP skills load from `.opencode/skills/` (symlinked to sdp repo)

## 3. GLM Model

**Option A: GLM via API (Novita, OpenRouter, etc.)**

Configure your IDE to use an OpenAI-compatible endpoint for GLM-4:

```json
{
  "baseUrl": "https://api.openrouter.ai/api/v1",
  "model": "zhipu/glm-4-flash"
}
```

**Option B: Local Docker Model Runner**

```bash
# Pull GLM or compatible model
docker run -d -p 12434:12434 ...

# Configure OpenCode to use http://localhost:12434/v1
```

## 4. Verify Install

```bash
# Check skills are linked
ls -la .opencode/skills
# Should show: .opencode/skills -> ../sdp/prompts/skills

# Check agents
ls -la .opencode/agents
```

## 5. Run @oneshot

For OpenCode, use `sdp-orchestrate` as the outer loop (OpenCode lacks Stop hooks):

```bash
sdp-orchestrate --feature F028 --runtime opencode
```

Requires `sdp-orchestrate` from [sdp_lab](https://github.com/fall-out-bug/sdp_lab).

## See Also

- [SDP README](../../README.md)
- [docs/attestation/coding-workflow-v1.md](../attestation/coding-workflow-v1.md) — evidence format
- [CHANGELOG](../../CHANGELOG.md) — v0.9.7+ for @feature --auto, @design
