# PRODUCT_VISION.md

> **Last Updated:** 2026-02-02
> **Version:** 2.0

## Mission

Transform SDP from a Python-centric CLI tool into a **language-agnostic Claude Plugin** that enforces development protocols through AI-driven validation and portable prompts, with optional tooling for automation.

## Users

1. **Non-Python Development Teams**
   - Java teams using Maven/Gradle
   - Go teams using go.mod
   - Polyglot teams working across multiple languages

2. **Solo Developers on Small Projects**
   - Want protocol enforcement without heavy Python dependency
   - Prefer lightweight setup via Claude Plugin Marketplace

3. **DevOps Engineers**
   - Need language-agnostic protocol enforcement
   - Want Git integration without runtime dependencies

## Success Metrics

- [ ] Zero Python dependency for basic SDP protocol usage
- [ ] Installation via Claude Plugin Marketplace (single click)
- [ ] Support for Python, Java, Go projects with language-agnostic validation
- [ ] Backward compatibility with existing Python SDP users

## Strategic Tradeoffs

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| **Validation Mechanism** | AI-based via prompts | Universal across languages, no language-specific parsers needed |
| **Distribution** | Claude Plugin Marketplace | Native integration with Claude Code/CLI, zero-install friction |
| **Binary Component** | Optional, not required | Prompts work standalone; binary adds convenience (CLI, git hooks) |
| **Language Support** | Language-agnostic first | Universal rules (file size, structure) over language-specific details |
| **Architecture** | Split protocol + tools | Protocol (prompts) is core; Tools (binary) is optional automation |

## Non-Goals

- ~~Language-specific code parsing~~ (AI validation instead)
- ~~Enforcing Python as primary runtime~~ (language-agnostic)
- ~~Breaking change for existing users~~ (Python SDP becomes reference implementation)
- ~~Requiring binary for basic usage~~ (prompts-only should work)
