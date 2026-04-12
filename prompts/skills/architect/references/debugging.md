# Debugging the Architect Module

## Empty C4 Diagram

If L2/L3 diagrams are empty:
1. Check `--tier` level (Tier 1 may omit detail)
2. Run with `--verbose` to see extractor timing — a zero-duration extractor may have skipped
3. Check if the repo has recognizable infra files (Dockerfile, compose, k8s YAML, pom.xml)
4. For L3: verify import graph clusters exist — components come from cluster→container mapping

## Wrong Container Detection

Container detection follows priority order:
1. Docker/Compose services
2. Kubernetes workloads
3. Service dependencies (depends_on)
4. Module boundaries (Maven, SBT, npm, Go cmd/)
5. Import graph clusters (fallback)
6. Single "Application" (last resort)

If wrong containers appear:
1. Check if CI-only containers leak through (filtered by `isCIContainer`)
2. Verify module boundary detection matches the build system
3. Run with `--extractors infra` to isolate infrastructure extraction

## Missing Language Analysis

If a language isn't detected:
1. Use `--language` flag to force detection
2. Verify the language adapter exists in `extract/adapters.go`
3. Check file extension mapping in `assembler.go` → `extToLanguage`

## Large JSON Output

Use `--section` to avoid huge blobs:
- `--section summary` — compact text (<5KB)
- `--section model` — C4 model only (~90KB vs ~600KB full)
- `--section profile` — codebase profile
- `--section report` — patterns and risks
- `--section diagrams` — mermaid code only

Never parse full JSON with Python one-liners. The `--section` flag exists for this reason.

## Timeout Issues

Default timeout is 5 minutes. For very large repos (>50K files):
- Use `--timeout 10m`
- Use `--skip-git` to skip git history (often the slowest extractor)
- Use `--extractors filetree,deps,infra` to run only fast extractors first
