# Welcome to SDP! 🚀

**Spec-Driven Protocol** - Workstream-driven development for AI agents with multi-agent coordination.

---

## 🎯 What is SDP?

SDP turns AI coding tools (Claude Code, Cursor, OpenCode) into a **structured software development process** with:

- ✅ **Atomic workstreams** - Small, focused tasks (500-1500 LOC)
- ✅ **Quality gates** - TDD, 80%+ coverage, type hints enforced
- ✅ **Task tracking** - Beads CLI integration for dependency management
- ✅ **Multi-agent coordination** - Specialized agents for planning, building, reviewing
- ✅ **Autonomous execution** - Execute entire features with checkpoint/resume

**Perfect for:** Solo developers and small teams using AI-IDEs (5-500 workstreams)

---

## 🚀 Quick Start (30 seconds)

```bash
# 1. Install
git clone https://github.com/fall-out-bug/sdp.git
cd sdp
pip install -e .

# 2. Create your first feature
@feature "Add user authentication"

# 3. Plan workstreams
@design idea-user-auth

# 4. Execute
@build WS-001-01

# Or execute all autonomously
@oneshot F001
```

**What happens:** AI interviews you → Creates workstreams → Executes with TDD → Validates quality → Tracks progress

---

## 📚 Where to Start?

### New to SDP? Start Here:

**Beginner Path (Progressive Learning):**

1. **[Quick Start (5 min)](docs/beginner/00-quick-start.md)** ⭐ Start here
   - What is SDP and why it matters
   - Key concepts at a glance
   - See it in action

2. **[Hands-on Tutorial (15 min)](docs/beginner/01-first-feature.md)** 🔧 Learn by doing
   - Create your first feature
   - Interactive walkthrough
   - Real examples

3. **[Common Tasks Guide](docs/beginner/02-common-tasks.md)** 📋 Daily reference
   - Feature development workflow
   - Bug fixes and debugging
   - Quality checks and deployment

4. **[Troubleshooting Guide](docs/beginner/03-troubleshooting.md)** 🔧 When stuck
   - Common issues and solutions
   - Error patterns
   - Debugging tips

**Quick References:**

- **[Glossary](docs/reference/GLOSSARY.md)** 📖 150+ term reference
- **[README](README.md)** 📄 Project overview

### Ready to Dive Deeper?

**By Role:**

- **Team Leads** → See [docs/SITEMAP.md](docs/SITEMAP.md#guides-by-role)
- **Engineers** → [reference/](docs/reference/) lookup docs
- **Maintainers** → [internals/](docs/internals/) architecture & extending

**By Topic:**

- **[Protocol Specification](PROTOCOL.md)** 📋 Complete SDP specification
- **[Code Patterns](CODE_PATTERNS.md)** 🔧 Implementation patterns
- **[Reference Docs](docs/reference/)** 📚 Command & config reference
- **[Internals](docs/internals/)** 🏗️ Architecture & extending
- **[Site Map](docs/SITEMAP.md)** 🗺️ Full documentation index

---

## 🎓 Learning Paths

### For Team Leads & Managers

**Goal:** Understand SDP workflow and capabilities

1. [README](README.md) - Project overview (5 min)
2. [docs/overview-for-leads.md](docs/overview-for-leads.md) - Executive summary (10 min)
3. [docs/TUTORIAL.md](docs/TUTORIAL.md) - Hands-on intro (15 min)
4. [PROTOCOL.md](PROTOCOL.md) - Full specification (30 min)

**Outcome:** Know what SDP does, how it works, when to use it

---

### For Engineers & Developers

**Goal:** Use SDP daily for feature development

1. [docs/TUTORIAL.md](docs/TUTORIAL.md) - Learn by doing (15 min)
2. [docs/GLOSSARY.md](docs/GLOSSARY.md) - Reference as needed (ongoing)
3. [CODE_PATTERNS.md](CODE_PATTERNS.md) - Implementation patterns (20 min)
4. [docs/PRINCIPLES.md](docs/PRINCIPLES.md) - Quality standards (30 min)
5. [PROTOCOL.md](PROTOCOL.md) - Complete reference (as needed)

**Outcome:** Ready to use @feature, @design, @build, @review, @deploy

---

### For AI Tool Users (Claude Code, Cursor)

**Goal:** Integrate SDP into your AI-IDE workflow

1. [CLAUDE.md](CLAUDE.md) - Claude Code integration (10 min)
2. [docs/TUTORIAL.md](docs/TUTORIAL.md) - Skill system walkthrough (15 min)
3. [docs/guides/CLAUDE_CODE.md](docs/guides/CLAUDE_CODE.md) - Claude Code guide (20 min)
4. [docs/guides/CURSOR.md](docs/guides/CURSOR.md) - Cursor guide (20 min)
5. [docs/multi-ide-parity.md](docs/multi-ide-parity.md) - Tool comparisons (10 min)

**Outcome:** Productive with SDP skills in your AI-IDE

---

## 🛠️ Common Tasks

### "I want to..."

**...add a new feature**
```bash
@feature "Feature description"
@design idea-{slug}
@build WS-001-01
```
→ See: [docs/TUTORIAL.md](docs/TUTORIAL.md), [PROTOCOL.md](PROTOCOL.md#feature-development-flow)

**...fix a bug**
```bash
@issue "Bug description"
# Routes to @hotfix (P0) or @bugfix (P1/P2)
```
→ See: [PROTOCOL.md](PROTOCOL.md#quick-reference), [docs/GLOSSARY.md](docs/GLOSSARY.md#severity-classification)

**...debug a failing test**
```bash
/debug "Test fails unexpectedly"
```
→ See: [docs/runbooks/debug-runbook.md](docs/runbooks/debug-runbook.md)

**...review code quality**
```bash
@review F001
```
→ See: [PROTOCOL.md](PROTOCOL.md#quality-review), [docs/two-stage-review.md](docs/two-stage-review.md)

**...deploy to production**
```bash
@deploy F001
```
→ See: [PROTOCOL.md](PROTOCOL.md#deployment), [docs/github-integration/SETUP.md](docs/github-integration/SETUP.md)

**...understand terminology**
→ See: [docs/GLOSSARY.md](docs/GLOSSARY.md) (150+ terms)

**...find specific documentation**
→ See: [docs/SITEMAP.md](docs/SITEMAP.md) (full index)

---

## 📂 Key Files

| File | Purpose | When to Read |
|------|---------|--------------|
| **START_HERE.md** | You are here! 👋 | First time visiting |
| **README.md** | Project overview | Understanding SDP at high level |
| **PROTOCOL.md** | Full specification | Deep dive into SDP workflow |
| **CODE_PATTERNS.md** | Implementation patterns | Writing code with SDP |
| **CLAUDE.md** | Claude Code integration | Using SDP with Claude Code |
| **docs/GLOSSARY.md** | 150+ term reference | Looking up unfamiliar terms |
| **docs/SITEMAP.md** | Documentation index | Finding specific docs |
| **docs/TUTORIAL.md** | 15-minute hands-on | Learning by doing |
| **docs/PRINCIPLES.md** | Engineering principles | Understanding quality standards |

---

## 🎯 Core Concepts (Cheat Sheet)

### Hierarchy

```
Release (product milestone)
  └─ Feature (5-30 workstreams)
      └─ Workstream (atomic task, one-shot)
```

**Example:** R1 → F24 → WS-AUTH-01

---

### Workflow

```
Idea → @feature → @design → @build → @review → @deploy
  ↓       ↓          ↓         ↓        ↓        ↓
Beads  Agents    Plan      TDD     Quality  Production
```

---

### Quality Gates

Every workstream must pass:
- ✅ TDD (tests first)
- ✅ Coverage ≥ 80%
- ✅ mypy --strict
- ✅ ruff clean
- ✅ Files < 200 LOC
- ✅ No `except: pass`

---

### Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `@feature` | Create feature | `@feature "Add user auth"` |
| `@design` | Plan workstreams | `@design idea-user-auth` |
| `@build` | Execute workstream | `@build WS-001-01` |
| `@review` | Quality check | `@review F001` |
| `@deploy` | Deploy to production | `@deploy F001` |
| `@oneshot` | Autonomous execution | `@oneshot F001` |
| `/debug` | Debug issues | `/debug "Test fails"` |
| `@issue` | Route bugs | `@issue "Login fails"` |

---

## 🆘 Need Help?

### Stuck? Try These:

1. **[Glossary](docs/GLOSSARY.md)** - Look up unfamiliar terms
2. **[Site Map](docs/SITEMAP.md)** - Find relevant documentation
3. **[Runbooks](docs/runbooks/)** - Step-by-step guides for common tasks
4. **[GitHub Issues](https://github.com/fall-out-bug/sdp/issues)** - Report bugs or request features

### Quick Questions:

**"What does {term} mean?"** → [docs/GLOSSARY.md](docs/GLOSSARY.md)

**"How do I {action}?"** → [docs/TUTORIAL.md](docs/TUTORIAL.md) or [docs/SITEMAP.md](docs/SITEMAP.md)

**"Where is {topic} documented?"** → [docs/SITEMAP.md](docs/SITEMAP.md)

**"Why does SDP {action}?"** → [PROTOCOL.md](PROTOCOL.md) or [docs/PRINCIPLES.md](docs/PRINCIPLES.md)

---

## 🎉 Next Steps

1. ✅ **Read** [15-Minute Tutorial](docs/TUTORIAL.md) (15 min)
2. ✅ **Try** creating your first feature with `@feature`
3. ✅ **Reference** [Glossary](docs/GLOSSARY.md) as needed
4. ✅ **Deep dive** into [PROTOCOL.md](PROTOCOL.md) when ready

**Welcome to workstream-driven development! 🚀**

---

**Version:** SDP v0.5.0
**Updated:** 2026-01-29
**Maintained by:** SDP Protocol Team

---

**Quick Links:**
- [Tutorial](docs/TUTORIAL.md) - Learn by doing
- [Glossary](docs/GLOSSARY.md) - 150+ terms
- [Site Map](docs/SITEMAP.md) - Full documentation index
- [Protocol](PROTOCOL.md) - Complete specification
- [GitHub](https://github.com/fall-out-bug/sdp) - Source code
