# Decision Log Backup and Restore Guide

## Overview

Decision logs are stored in `docs/decisions/decisions.jsonl` in JSONL (newline-delimited JSON) format. This guide covers backup and restore procedures.

## Backup Methods

### Method 1: Git Version Control (Recommended)

Decisions are tracked in git by default (after removing from `.gitignore`):

```bash
# Add decisions to git
git add docs/decisions/decisions.jsonl
git commit -m "docs: add decision log"

# Push to remote
git push
```

**Pros:** Automatic versioning, distributed backup  
**Cons:** Commits needed, large binary files

### Method 2: Export to Markdown (Manual)

Export decisions to human-readable markdown:

```bash
# Export to default location
sdp decisions export

# Export to custom location
sdp decisions export /path/to/backup/DECISIONS.md
```

**Pros:** Human-readable, format for ADR  
**Cons:** Manual, one-way

### Method 3: JSONL File Copy (Simple)

Direct file copy:

```bash
# Timestamped backup
cp docs/decisions/decisions.jsonl \
   "backups/decisions-$(date +%Y%m%d-%H%M%S).jsonl"

# Compressed backup
gzip -c docs/decisions/decisions.jsonl > \
   "backups/decisions-$(date +%Y%m%d-%H%M%S).jsonl.gz"
```

**Pros:** Simple, fast  
**Cons:** Manual, no versioning

## Restore Procedures

### Scenario 1: Recover from Git History

```bash
# Find commit containing decisions
git log --all --full-history -- "docs/decisions/decisions.jsonl"

# Restore from specific commit
git show <commit-hash>:docs/decisions/decisions.jsonl > \
   docs/decisions/decisions.jsonl
```

### Scenario 2: Restore from JSONL Backup

```bash
# Decompress if needed
gunzip -c backups/decisions-20250206-120000.jsonl.gz > \
   docs/decisions/decisions.jsonl

# Or direct copy
cp backups/decisions-20250206-120000.jsonl \
   docs/decisions/decisions.jsonl
```

### Scenario 3: Corrupted decisions.jsonl

**Symptoms:** `sdp decisions list` shows parse errors

**Diagnosis:**
```bash
# Check for malformed JSON
python3 -m json.tool docs/decisions/decisions.jsonl > /dev/null
# Look for "Expecting" errors

# Count lines vs valid decisions
wc -l docs/decisions/decisions.jsonl
sdp decisions list | grep "Found.*decision"
```

**Recovery Options:**

**Option A: Extract valid decisions only**
```bash
# Create recovery script
python3 << 'PYTHON'
import json

valid = []
with open('docs/decisions/decisions.jsonl', 'r') as f:
    for line in f:
        try:
            valid.append(json.loads(line))
        except:
            pass  # Skip malformed lines

# Write cleaned file
with open('docs/decisions/decisions.jsonl', 'w') as f:
    for d in valid:
        f.write(json.dumps(d) + '\n')

print(f"Recovered {len(valid)} valid decisions")
PYTHON
```

**Option B: Restore from latest export**
```bash
# If markdown export exists, recreate from it
# Note: This loses the original JSON format
# You'll need to manually convert back to JSONL
```

## Disaster Recovery

### Total Data Loss

If `docs/decisions/` is completely lost:

1. **Check git history** (most recent first):
   ```bash
   git log --all -- docs/decisions/decisions.jsonl
   git show <commit-hash>:docs/decisions/decisions.jsonl
   ```

2. **Check remote backups**:
   - GitHub/GitLab repository
   - Time Machine (macOS)
   - System backups

3. **Recreate from memory** (last resort):
   - Search codebase for decisions made
   - Check git commits for rationale
   - Interview team members

### Partial Data Loss

If only some decisions are lost:

1. **Identify missing decisions** from git commits:
   ```bash
   git log --oneline --all | grep -i "decision"
   git show <commit> # Look for decisions in commit messages
   ```

2. **Manual re-entry**:
   ```bash
   sdp decisions log \
     --question "What was decided?" \
     --decision "The decision" \
     --rationale "Because..." \
     --type technical
   ```

## Maintenance

### Regular Backups (Recommended)

Add to cron or launchd:

```bash
# Daily backup at 2 AM
0 2 * * * cp docs/decisions/decisions.jsonl \
   "backups/decisions-$(date +\%Y\%m\%d-\%H\%M\%S).jsonl"
```

### Retention Policy

- **Daily backups:** Keep 30 days
- **Weekly backups:** Keep 12 weeks
- **Monthly backups:** Keep 12 months
- **Git history:** Keep forever

### Monitor File Size

```bash
# Check decision log size
ls -lh docs/decisions/decisions.jsonl

# Alert if > 10MB
if [ $(stat -f%z docs/decisions/decisions.jsonl) -gt 10485760 ]; then
    echo "WARNING: Decision log exceeds 10MB"
fi
```

## Migration Guide

### Python SDP → Go Plugin

If migrating from Python SDP:

1. **Export from Python** (if still available):
   ```bash
   # Python SDP stored decisions in YAML
   # Convert to JSONL manually
   ```

2. **Import to Go**:
   ```bash
   # Create decisions.jsonl in correct format
   sdp decisions log --question="..." --decision="..."
   ```

3. **Verify**:
   ```bash
   sdp decisions list
   sdp decisions export
   ```

## Troubleshooting

### Issue: "No decisions found" but file exists

**Cause:** File is empty or malformed

**Fix:**
```bash
# Check file contents
cat docs/decisions/decisions.jsonl

# If empty, that's expected (no decisions yet)
# If malformed, follow "Corrupted decisions.jsonl" above
```

### Issue: Permission denied writing decisions

**Cause:** File ownership or permissions

**Fix:**
```bash
# Fix permissions
chmod 644 docs/decisions/decisions.jsonl
chown $USER:$USER docs/decisions/decisions.jsonl

# Fix directory
chmod 755 docs/decisions/
```

## Best Practices

1. **Commit to git regularly** after important decisions
2. **Export to markdown** before major releases
3. **Monitor file size** to prevent unbounded growth
4. **Test backups** by restoring periodically
5. **Document critical decisions** in ADR format

## Related Documentation

- [SRE SLOs](../../SRE_SLOS.md) - Performance targets
- [Quality Gates](../../quality-gates.md) - Decision audit requirements
- [ADR Format](https://github.com/joelparkerhenderson/adr) - Architecture Decision Records
