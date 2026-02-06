# SDP Telemetry Privacy Policy

## Overview

SDP collects anonymized usage telemetry to improve the tool's reliability and performance. This document explains what data we collect, how it's used, and your privacy options.

**Last Updated:** 2026-02-06

## What Data We Collect

### Collected Data

SDP collects the following **anonymized** usage data:

| Data Type | Description | Example | PII? |
|-----------|-------------|---------|-----|
| **Command invocations** | Which SDP commands are run | `@build`, `@review`, `sdp doctor` | ❌ No |
| **Execution duration** | How long commands take to run | `45 seconds`, `2 minutes` | ❌ No |
| **Success/failure** | Whether commands completed successfully | `success`, `error: timeout` | ❌ No |
| **Workstream IDs** | Which workstreams are being executed | `00-001-01`, `00-050-13` | ❌ No |
| **Feature IDs** | Which features are being developed | `F01`, `F05` | ❌ No |
| **Quality gate results** | Pass/fail status of quality checks | `coverage: 85%`, `passed` | ❌ No |

### What We DON'T Collect

❌ **Personally Identifiable Information (PII):**
- Your name
- Your email address
- Your username
- Your hostname
- File paths (except relative paths within project)
- Project names
- Code content
- Environment variables
- Secrets or credentials

❌ **Sensitive project data:**
- Source code content
- Commit messages
- Branch names
- Repository URLs
- API keys or tokens

## How Data Is Used

Collected telemetry is used **exclusively** for:

1. **Performance Optimization:** Identifying slow commands or workflows
2. **Reliability Improvements:** Detecting and fixing common failure patterns
3. **Feature Usage:** Understanding which features are most/least used
4. **Quality Metrics:** Tracking quality gate pass rates

**We do NOT:**
- Sell data to third parties
- Use data for marketing purposes
- Share data with anyone outside the SDP development team
- Analyze code content or intellectual property

## Data Storage

### Location

All telemetry data is stored **locally** on your machine:

```
~/.sdp/telemetry.jsonl  # Event log
~/.sdp/telemetry.json   # Configuration
```

### File Format

Telemetry is stored in **JSONL format** (one JSON object per line):

```jsonl
{"type":"command_start","timestamp":"2026-02-06T10:30:00Z","data":{"command":"@build","ws_id":"00-001-01"}}
{"type":"command_complete","timestamp":"2026-02-06T10:30:45Z","data":{"command":"@build","duration":45,"status":"success"}}
```

### Data Retention

**Default retention period:** 90 days

Old telemetry events are automatically deleted after 90 days. You can configure this:

```json
// ~/.sdp/telemetry.json
{
  "enabled": true,
  "retention_days": 90
}
```

To disable auto-cleanup:

```json
{
  "retention_days": 0  // Keep forever
}
```

## Your Privacy Options

### Opt-In (Enable Telemetry)

**Default: DISABLED** 🔒

SDP telemetry is **opt-in by default**. This means:
- Telemetry is DISABLED until you explicitly enable it
- Your privacy is respected from the first run
- No data is collected without your consent

**To enable telemetry:**

```bash
sdp telemetry enable
```

You will see:
```
✓ Telemetry enabled
  Thank you for helping improve SDP!
  To disable: sdp telemetry disable
```

**To disable telemetry:**

**Method 1: Configuration file**
```bash
echo '{"enabled": false}' > ~/.sdp/telemetry.json
```

**Method 2: Environment variable**
```bash
export SDP_TELEMETRY_ENABLED=false
```

**Method 3: Command-line flag**
```bash
sdp --no-telemetry
```

### Verify Telemetry Status

```bash
sdp telemetry status
```

Output:
```
Telemetry Status
================
Enabled: No
Event Count: 0
File Path: /home/user/.sdp/telemetry.jsonl
```

### Clear Telemetry Data

```bash
sdp telemetry clear
```

**Warning:** This permanently deletes all locally stored telemetry data.

### Export Telemetry Data

```bash
sdp telemetry export /tmp/my-telemetry.jsonl
```

You can review your telemetry data before sharing or analyzing it.

## Data Sharing and Disclosure

### We NEVER Share Your Data

SDP telemetry is:
- ❌ Not sent to any remote server
- ❌ Not shared with third parties
- ❌ Not used for commercial purposes
- ❌ Not sold or monetized

### Voluntary Data Submission

If you choose to share telemetry data with the SDP team (e.g., for bug reports):

1. **Export your data:** `sdp telemetry export telemetry.jsonl`
2. **Review the file:** Ensure no sensitive information is included
3. **Submit voluntarily:** Attach to GitHub issue or email

**We will only use voluntarily submitted data** for the purpose stated in your submission (e.g., debugging a specific issue).

## Children's Privacy

SDP is a professional development tool and is not intended for use by children under 13. We do not knowingly collect information from children.

## Security

### Data Protection

- **File permissions:** Telemetry files are created with `0600` (owner read/write only)
- **No network transmission:** Data never leaves your machine unless you choose to export it
- **Encryption:** Not applicable (data stays local)

### Data Breach Notification

Since telemetry data is stored locally and never transmitted, there is no risk of remote data breaches. However, if your local machine is compromised:

- Attackers could access anonymized usage data
- Data does NOT contain PII, credentials, or sensitive content
- To completely remove data: `sdp telemetry clear`

## Compliance

### GDPR (EU General Data Protection Regulation)

Since SDP telemetry:
- Does NOT collect PII
- Does NOT transfer data outside the EU
- Does NOT use data for automated decision-making
- Stores data locally on your machine

**GDPR compliance status:** ✅ Compliant (no personal data processed)

### CCPA (California Consumer Privacy Act)

Under CCPA, SDP telemetry:
- Does NOT sell personal information
- Does NOT collect sensitive personal information
- Allows you to opt-out (see "Your Privacy Options" above)

**CCPA compliance status:** ✅ Compliant

## Changes to This Policy

We may update this privacy policy from time to time. We will notify users of significant changes via:

1. Update to the `PRIVACY.md` file in SDP releases
2. Announcement in release notes
3. In-product notification (for major changes)

**Continued use of SDP after policy changes constitutes acceptance of the new policy.**

## Contact and Questions

### Privacy Questions

If you have questions about SDP telemetry or privacy:

- **GitHub Issues:** https://github.com/ai-masters/sdp/issues
- **Email:** privacy@ai-masters.com (hypothetical)

### Data Access Requests

To access, export, or delete your telemetry data:

```bash
# Access/Export
sdp telemetry export my-telemetry.jsonl

# Delete
sdp telemetry clear
```

### Opt-Out Requests

See "Your Privacy Options" above for instructions on disabling telemetry.

## Transparency

### View Your Telemetry Data

```bash
# View raw telemetry log
cat ~/.sdp/telemetry.jsonl

# View summarized statistics
sdp telemetry analyze
```

### Audit Trail

All telemetry operations are logged locally:

```bash
# View what SDP has been tracking
tail -f ~/.sdp/telemetry.jsonl
```

## Summary

✅ **No PII collected**
✅ **No data transmission** (stays local)
✅ **Opt-out available**
✅ **Auto-cleanup after 90 days**
✅ **Clear documentation**
✅ **GDPR/CCPA compliant**

---

**This policy applies to SDP version 0.3.0 and later.**
