# Security Guide

## Secret Management

### Principles

1. **NEVER commit secrets to git** - This includes API keys, tokens, passwords
2. **Use environment variables** - Load secrets from environment at runtime
3. **Use GitHub Secrets** - For CI/CD pipeline secrets
4. **Rotate compromised secrets** - If a secret is exposed, rotate it immediately

### Environment Variables

The `.env` file is **gitignored** and should **NEVER** be committed.

**Required for development:**
- `TELEGRAM_BOT_TOKEN` - Telegram bot token for notifications (optional)

**Setup:**
```bash
# Copy template
cp .env.example .env

# Edit with your values
nano .env

# Load in your shell
source .env  # or export TELEGRAM_BOT_TOKEN="your_token_here"
```

### CI/CD Secrets

**GitHub Secrets (used in CI/CD):**
- None currently needed for this CLI tool

**To add secrets:**
1. Go to: https://github.com/fall-out-bug/sdp/settings/secrets/actions
2. Click: "New repository secret"
3. Name: `TELEGRAM_BOT_TOKEN` (or appropriate)
4. Value: Your secret value
5. Enable: "Required for workflow"

### Secret Rotation

**If a secret is exposed (committed, leaked, etc.):**

1. **Immediately rotate the secret**
   ```bash
   # Telegram bot token
   # 1. Message @BotFather on Telegram
   # 2: Select your bot → Revoke old token
   # 3. Generate new token
   # 4. Update .env locally
   ```

2. **Verify it's not in git history**
   ```bash
   git log --all --full-history -- .env
   # Should return: "No history found"
   ```

3. **If it WAS committed:**
   ```bash
   # Remove from all commits
   git filter-branch --force --index-filter \
     "git rm --cached --ignore-unmatch .env"
   
   # Force push
   git push origin --force --all
   ```

### Current Status (2026-02-07)

✅ **SAFE** - `.env` is gitignored and not in git history
⚠️ **Token exposed locally** - This is acceptable for development
📋 **Template** - `.env.example` provided for setup

### Best Practices

1. **Never commit `.env`** - It's in `.gitignore`
2. **Use `.env.example`** - Template for required variables
3. **Document secrets** - Keep this SECURITY.md up to date
4. **Review regularly** - Audit git log for accidental commits
5. **Use `.env.local`** - For local overrides (also gitignored)

### Monitoring

**Check for exposed secrets:**
```bash
# Search git history for sensitive patterns
git log --all --oneline -S | grep -i "token\|secret\|password\|api[_-]key"

# Search all tracked files for secrets
grep -r "TELEGRAM_BOT_TOKEN\|password\|secret" --include="*.go" --exclude-dir=vendor
```

---

**Last Updated:** 2026-02-07  
**Version:** 1.0
