---
name: deploy
description: Deployment orchestration. Generates docker-compose, CI/CD, release notes. Merges feature to main with tagging.
tools: Read, Write, Edit, Bash, Glob, Grep
---

# /deploy - Deployment Orchestration

Generate deployment artifacts and execute GitFlow merge.

## When to Use

- After `/review` APPROVED
- After human UAT passes
- Ready for production

## Invocation

```bash
/deploy F60
```

## Master Prompt

📄 **sdp/prompts/commands/deploy.md** (480+ lines)

**Contains:**
- Pre-deployment checks (e2e tests)
- Version resolution (semantic versioning)
- Artifact generation (docker-compose, CI/CD)
- GitFlow merge (feature → develop → main)
- Tag creation
- Release notes
- GitHub notification

## Workflow

1. Pre-flight: e2e tests (via pre-deploy.sh)
2. Version resolution dialogue
3. Generate artifacts:
   - docker-compose.yml
   - .github/workflows/
   - Release notes
4. GitFlow merge:
   - feature → develop (fast-forward)
   - develop → main (--no-ff)
5. Tag: v{X.Y.Z}
6. Push + cleanup branches
7. Update INDEX (WS to completed/)

## Generated Artifacts

- `docker-compose.yml` (production-ready)
- `.github/workflows/ci.yml`
- `.github/workflows/deploy.yml`
- `docs/releases/v{X.Y.Z}.md`

## GitFlow

```
feature/{slug} → develop → main
                            ↓
                        tag v{X.Y.Z}
```

## Output

Feature deployed to main + tagged + release notes + cleanup

## Quick Reference

**Input:** Feature ID (APPROVED + UAT passed)  
**Output:** Production deployment + v{X.Y.Z} tag  
**Next:** Monitor production
