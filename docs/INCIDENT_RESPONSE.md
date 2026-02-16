# Incident Response Procedures

**Scope:** SDP CLI tool (developer library, not hosted service)

## Severity Definitions

### P0 - Critical (Immediate Response Required)
- **Data loss/corruption:** CLI deletes or corrupts user data
- **Security vulnerability:** Known exploit in wild
- **Build failure:** Binary compilation broken, no installable version
- **Distribution down:** Homebrew, GitHub releases unavailable for >1 hour

**Response Time:** < 4 hours
**Examples:**
- `sdp parse` deletes workstream files
- CVE published with exploit
- `go build` fails on main branch

### P1 - High (Next Release)
- **Breaking change:** CLI behavior change without migration path
- **Feature regression:** Previously working feature now broken
- **Documentation gap:** Critical documentation missing for new feature

**Response Time:** < 48 hours
**Examples:**
- Workstream ID format change breaks existing WS files
- `sdp telemetry` crashes on valid input
- New `@oneshot` skill not documented

### P2 - Medium (Backlog)
- **UX friction:** Confusing error messages, poor defaults
- **Performance:** CLI operation slower than expected
- **Code quality:** Technical debt, refactoring needs

**Response Time:** Next sprint
**Examples:**
- Error message "failed" without details
- `sdp parse` takes >10s on large workstreams
- Coverage below 80%

---

## Hotfix Procedure

**When:** P0 bugs only

### Process

```bash
# 1. Identify severity
bd create --title="P0: Data loss in parse command" --priority=0 --type=bug

# 2. Create hotfix branch from main
git checkout main
git pull
git checkout -b hotfix/parse-data-loss

# 3. Implement fix (TDD required)
# Add failing test first, then fix
@tdd red ./cmd/sdp
@tdd green ./cmd/sdp
@tdd refactor ./cmd/sdp

# 4. Quick quality check (skip full review)
go test ./...
go build ./cmd/sdp

# 5. Tag and release (bypass PR)
git commit -m "hotfix: fix data loss in parse command (P0)"
git tag v0.8.1
git push origin main --tags

# 6. Trigger Goreleaser
gh release create v0.8.1 --generate-notes
```

### Post-Hotfix
- Update CHANGELOG.md
- Close Beads issue
- Announce in GitHub Release notes
- Create follow-up task for root cause analysis

**Time Budget:** < 4 hours from detection to release

---

## Security Vulnerability Response

**When:** Security report received (GitHub Security Advisory, email, issue)

### Process

```bash
# 1. Triage (within 1 hour)
- Confirm vulnerability exists
- Assess severity (CVSS score)
- Check if exploit in wild

# 2. Coordinate disclosure
- Create private security advisory
- Notify maintainers privately
- Set disclosure deadline (typically 7 days)

# 3. Develop fix (private branch)
git checkout -b security/cve-2025-xxxxx
# Implement fix without public commits

# 4. Release fix
git checkout main
git merge security/cve-2025-xxxxx
git tag v0.8.2-security
git push origin main --tags

# 5. Publish advisory
gh security advisory publish
```

### Prevention
- Run `go vet ./...` in CI
- Use `gosec` security scanner
- Audit dependencies with `go list -json -m all | nancy sleuth`

---

## Rollback Procedure

**When:** Bad release deployed (P0 bug in production)

### Scenario 1: GitHub Release
```bash
# 1. Delete release (keep tag)
gh release delete v0.8.0 --yes
# Note: Tag remains for git history

# 2. Create new release (patch version)
git checkout v0.7.9  # Last known good
git checkout -b hotfix/rollback-v0.8.1
git tag v0.8.1
gh release create v0.8.1 --notes "Rollback of v0.8.0 due to P0 bug"

# 3. Update Homebrew
brew bump-formula-pr --url=... sdp
```

### Scenario 2: Homebrew Formula Issue
```bash
# 1. Uninstall broken version
brew uninstall sdp

# 2. Install previous version
brew install sdp@0.7.9
brew link sdp@0.7.9

# 3. Pin until fix released
brew pin sdp
```

### Verification
```bash
# Test rollback locally
sdp --version  # Should show rolled-back version
sdp parse 00-050-01  # Should work without P0 bug
```

---

## Distribution Failures

### GitHub Releases Down
**Symptom:** `brew install sdp` fails with 404

**Workaround:**
```bash
# Install from source
git clone https://github.com/fall-out-bug/sdp
cd sdp/sdp-plugin
go build ./cmd/sdp
sudo mv sdp /usr/local/bin/
```

### Homebrew Formula Outdated
**Symptom:** `brew upgrade sdp` installs old version

**Fix:**
```bash
# Update formula manually
brew edit sdp
# Change url/sha256 to latest release

# Submit PR
brew bump-formula-pr --url=... sdp
```

---

## Communication Templates

### P0 Incident Announcement
```
🚨 INCIDENT: P0 Bug in v0.8.0

**Issue:** [brief description]
**Impact:** [who affected, what breaks]
**Workaround:** [immediate mitigation]
**Fix:** v0.8.1 released
**Upgrade:** brew upgrade sdp

We apologize for the disruption.
```

### Security Disclosure
```
🔒 SECURITY: CVE-2025-XXXXX

**Severity:** [CVSS score, Critical/High/Medium/Low]
**Affected Versions:** v0.7.0 - v0.8.0
**Fixed In:** v0.8.2
**Action Required:** Upgrade immediately

Upgrade: brew upgrade sdp
Report: https://github.com/fall-out-bug/sdp/security/advisories/GHSA-xxxxx
```

---

## Related Documentation

- [PROTOCOL.md](PROTOCOL.md) - SDP quality gates
- [development/release.md](development/release.md) - Release procedure
- [GitHub Security Advisories](https://github.com/fall-out-bug/sdp/security/advisories)

---

**Last Updated:** 2026-02-06
**Version:** 1.0
