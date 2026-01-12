# Deployment Plan: User Profile API

## Overview

This document outlines the deployment strategy for the User Profile API feature.

## Pre-Deployment Checklist

- [ ] All tests passing (unit, integration, E2E)
- [ ] Code review completed and approved
- [ ] Architecture audit passed
- [ ] Security review passed
- [ ] Documentation updated
- [ ] Database migrations tested
- [ ] Rollback plan verified

## Deployment Stages

### Stage 1: Staging Environment

**Actions:**
1. Deploy to staging environment
2. Run smoke tests
3. Verify API responses
4. Check logs for errors
5. Monitor metrics (response time, error rate)

**Verification:**
```bash
# Smoke test
curl -X GET https://staging.example.com/api/users/test/profile \
  -H "Authorization: Bearer $TOKEN"

# Expected: 200 OK with profile data
```

**Rollback trigger:**
- Any smoke test failure
- Error rate > 1%
- Response time > 500ms (p95)

### Stage 2: Production Canary (10%)

**Actions:**
1. Deploy to 10% of production traffic
2. Monitor for 15 minutes
3. Compare metrics with baseline

**Metrics to monitor:**
| Metric | Baseline | Threshold |
|--------|----------|-----------|
| Error rate | 0.05% | < 0.5% |
| Response time (p95) | 150ms | < 250ms |
| CPU usage | 40% | < 60% |
| Memory usage | 50% | < 70% |

**Rollback trigger:**
- Error rate > 0.5%
- Response time > 250ms
- Any 5xx errors

### Stage 3: Production Full Rollout

**Actions:**
1. Gradually increase to 100% traffic
2. Monitor for 30 minutes
3. Verify all endpoints functional
4. Update status page

**Verification:**
- All health checks passing
- No alerts triggered
- Customer-facing functionality working

## Database Migrations

### Migration Script
```sql
-- 001_add_profile_fields.sql
-- Adds bio and updated_at fields to users table

ALTER TABLE users
  ADD COLUMN bio VARCHAR(500) DEFAULT '',
  ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Backfill updated_at
UPDATE users SET updated_at = created_at WHERE updated_at IS NULL;
```

### Migration Procedure
1. Backup database before migration
2. Run migration in staging, verify
3. Run migration in production during low-traffic window
4. Verify data integrity

### Rollback Migration
```sql
-- 001_add_profile_fields_rollback.sql
ALTER TABLE users
  DROP COLUMN bio,
  DROP COLUMN updated_at;
```

## Rollback Plan

### Immediate Rollback (< 5 min)
1. Revert to previous container image
2. Verify old version is serving traffic
3. Check health endpoints

```bash
# Kubernetes rollback
kubectl rollout undo deployment/api-server

# Docker rollback
docker-compose up -d --force-recreate api-server:previous
```

### Data Rollback
If database changes need reverting:
1. Stop new deployments
2. Run rollback migration
3. Restore from backup if needed
4. Verify data integrity

## Configuration Changes

### Environment Variables
| Variable | Value | Description |
|----------|-------|-------------|
| PROFILE_BIO_MAX_LENGTH | 500 | Max bio length |
| PROFILE_NAME_MAX_LENGTH | 100 | Max name length |
| PROFILE_CACHE_TTL | 300 | Cache TTL in seconds |

### Feature Flags
| Flag | Default | Description |
|------|---------|-------------|
| profile_api_enabled | true | Enable profile endpoints |
| profile_caching | false | Enable response caching |

## Monitoring

### Dashboards
- API response times by endpoint
- Error rates by status code
- Database query performance
- Cache hit/miss ratio

### Alerts
| Alert | Condition | Severity |
|-------|-----------|----------|
| High error rate | > 1% for 5 min | Critical |
| Slow response | p95 > 500ms for 5 min | Warning |
| Database connection | Pool exhausted | Critical |
| Auth service down | No response | Critical |

## Post-Deployment

### Verification Steps
1. Verify all endpoints responding
2. Check logs for errors
3. Verify metrics within normal range
4. Test a few real requests manually

### Documentation Updates
- [ ] API documentation updated
- [ ] Changelog updated
- [ ] Runbook updated (if needed)
- [ ] Status page updated

### Communication
- Notify stakeholders of successful deployment
- Update ticket/feature status to "Deployed"
- Schedule post-deployment review if needed

---

**Status**: Approved by DevOps
**Version**: 1.0
**Last Updated**: 2025-12-27
