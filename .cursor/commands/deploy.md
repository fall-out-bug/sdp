# /deploy — Deploy Feature

When calling `/deploy {feature}`:

1. Load full prompt: `@prompts/commands/deploy.md`
2. Verify all WS are APPROVED
3. Execute Mandatory Dialogue (scope, environments)
4. Generate:
   - docker-compose updates
   - CI/CD pipeline updates
   - CHANGELOG.md entry
   - Release notes
   - Deployment plan

## Quick Reference

**Input:** APPROVED feature
**Output:** DevOps configs + docs + release notes
**Next:** Execute deployment plan
