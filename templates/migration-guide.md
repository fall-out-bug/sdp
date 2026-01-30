# Migration Guide

**From:** {old_version}  
**To:** {new_version}  
**Date:** {YYYY-MM-DD}

---

## Overview

This guide helps you upgrade from v{old} to v{new} which contains breaking changes.

**⚠️ Important:** Read BREAKING_CHANGES.md first for summary.

---

## Prerequisites

Before starting migration:

- [ ] Backup production database
- [ ] Review all breaking changes
- [ ] Test on staging environment
- [ ] Plan maintenance window
- [ ] Prepare rollback plan

---

## Step 1: Database Migration

### 1.1 Backup

```bash
# Backup database
pg_dump -h localhost -U postgres hw_checker > backup-$(date +%Y%m%d).sql
```

### 1.2 Run Migrations

```bash
cd tools/hw_checker

# Check current version
poetry run alembic current

# Run migrations
poetry run alembic upgrade head

# Verify
poetry run alembic current
# Should show: {new_revision}
```

### 1.3 Data Migration (if needed)

```bash
# Run data migration script
poetry run python scripts/migrate_data_v{version}.py

# Verify data integrity
poetry run python scripts/verify_migration.py
```

---

## Step 2: Update Application Code

### 2.1 API Client Updates

**Change 1: Updated endpoint**

```python
# Before (v{old})
response = client.post("/old-endpoint", json={
    "old_field": value
})

# After (v{new})
response = client.post("/new-endpoint", json={
    "new_field": value  # Field renamed
})
```

**Change 2: Response format**

```python
# Before (v{old})
result = response.json()
data = result["data"]

# After (v{new})
result = response.json()
data = result["items"]  # Key renamed
```

### 2.2 CLI Script Updates

```bash
# Before (v{old})
hwc grading run --repo-url https://github.com/user/repo

# After (v{new})
hwc grading run --repo https://github.com/user/repo
```

**Update scripts:**
```bash
# Find all uses of old argument
grep -r "--repo-url" scripts/

# Replace with --repo
sed -i 's/--repo-url/--repo/g' scripts/*.sh
```

---

## Step 3: Update Configuration

### 3.1 Config File Format

**hw_checker.yaml:**

```yaml
# Before (v{old})
old_config:
  setting: value
  
# After (v{new})
new_config:
  setting: value
  additional_setting: default  # New required field
```

### 3.2 Environment Variables

```bash
# Before (v{old})
export OLD_VAR_NAME=value

# After (v{new})
export NEW_VAR_NAME=value
```

---

## Step 4: Backward Compatibility (Transition Period)

If you need to support both old and new versions:

```python
# Wrapper for backward compatibility
def submit_work(repo: str, **kwargs):
    """Support both old and new API."""
    try:
        # Try new endpoint first
        return client.post("/new-endpoint", json={"repo": repo})
    except Exception:
        # Fallback to old endpoint
        return client.post("/old-endpoint", json={"repo_url": repo})
```

**Remove this wrapper after:** {YYYY-MM-DD + 30d}

---

## Step 5: Testing

### 5.1 Smoke Tests

```bash
cd tools/hw_checker

# Run smoke tests
poetry run pytest tests/smoke/ -v
```

### 5.2 Integration Tests

```bash
# Run full integration suite
poetry run pytest tests/integration/ -v
```

### 5.3 Manual Verification

1. Submit test work via API: ✅/❌
2. Check status via CLI: ✅/❌
3. Verify database records: ✅/❌
4. Check worker processing: ✅/❌

---

## Step 6: Deployment

### 6.1 Staging

```bash
# Deploy to staging
docker-compose -f docker-compose.staging.yml down
docker-compose -f docker-compose.staging.yml pull
docker-compose -f docker-compose.staging.yml up -d

# Verify
curl https://staging.api/health
# Should return: {"status": "ok", "version": "{new_version}"}
```

### 6.2 Production

**Maintenance window:** {start_time} - {end_time}

```bash
# 1. Enable maintenance mode
curl -X POST https://api/admin/maintenance/enable

# 2. Stop services
docker-compose -f docker-compose.prod.yml down

# 3. Backup (again, just before upgrade)
pg_dump -h localhost -U postgres hw_checker > backup-final-$(date +%Y%m%d-%H%M).sql

# 4. Deploy new version
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d

# 5. Run migrations
docker-compose -f docker-compose.prod.yml exec api alembic upgrade head

# 6. Verify health
curl https://api/health

# 7. Disable maintenance mode
curl -X POST https://api/admin/maintenance/disable
```

### 6.3 Monitor

**Monitor continuously:**
- Error rate: should be < 0.1%
- Response time: should be < 200ms
- Database connections: should be stable

```bash
# Watch logs
docker-compose -f docker-compose.prod.yml logs -f api

# Watch metrics
# Check Grafana: {grafana_url}
```

---

## Rollback Plan

If issues occur:

### Immediate Rollback

```bash
# 1. Stop new version
docker-compose -f docker-compose.prod.yml down

# 2. Restore database
psql -h localhost -U postgres hw_checker < backup-final-{timestamp}.sql

# 3. Deploy old version
docker tag hw-checker:{new_version} hw-checker:{new_version}-broken
docker tag hw-checker:{old_version} hw-checker:latest

docker-compose -f docker-compose.prod.yml up -d

# 4. Verify
curl https://api/health
# Should return: {"status": "ok", "version": "{old_version}"}
```

### Rollback Timeline

- **< 30 min:** Issues detected
- **30-45 min:** Decision to rollback
- **45-60 min:** Rollback complete
- **60-90 min:** Verify old version stable

---

## Common Issues

### Issue 1: Database migration fails

**Symptom:** `alembic upgrade` returns error

**Solution:**
```bash
# Check current version
alembic current

# Check migration log
alembic history

# If stuck, manually fix
psql hw_checker -c "UPDATE alembic_version SET version_num = '{target}';"
```

### Issue 2: API returns 500 errors

**Symptom:** All API requests fail

**Possible causes:**
1. Database not migrated
2. Config file not updated
3. Environment variables missing

**Debug:**
```bash
# Check logs
docker-compose logs api | tail -100

# Check config
docker-compose exec api cat /app/config/hw_checker.yaml

# Check env vars
docker-compose exec api env | grep HW_CHECKER
```

---

## Post-Migration Checklist

- [ ] All smoke tests pass
- [ ] Integration tests pass
- [ ] Production deployed successfully
- [ ] No error spikes in monitoring
- [ ] Rollback plan tested on staging
- [ ] Team notified of changes
- [ ] Documentation updated
- [ ] Old code/configs archived

---

## Support

**Need help?**
- Slack: #{channel}
- Email: {support_email}
- Issues: {github_issues_url}

**Emergency contact:** {phone_number}
