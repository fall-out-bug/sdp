# Sample Feature: User Profile API

This is a complete example of a feature that has gone through the full consensus workflow.

## Summary
Add user profile endpoints to the REST API, allowing users to view and update their profile information.

## Goals
- Create GET /api/users/{id}/profile endpoint
- Create PUT /api/users/{id}/profile endpoint
- Validate input data (email format, name length)
- Return appropriate HTTP status codes
- Add comprehensive test coverage (≥80%)

## Non-Goals (v1)
- Profile picture upload
- Social media links
- Profile visibility settings
- Activity history

## Users
- **Authenticated User**: Can view and edit their own profile
- **Admin**: Can view any user's profile

## API Specification

### GET /api/users/{id}/profile
**Response 200:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "John Doe",
  "bio": "Software developer",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-15T12:30:00Z"
}
```

**Response 404:** User not found
**Response 403:** Not authorized to view this profile

### PUT /api/users/{id}/profile
**Request:**
```json
{
  "name": "John Doe",
  "bio": "Updated bio"
}
```

**Response 200:** Updated profile
**Response 400:** Validation error
**Response 403:** Not authorized to edit this profile

## Success Metrics
- Response time < 200ms (p95)
- Test coverage ≥ 80%
- Zero security vulnerabilities
- API documentation complete

## Dependencies
- Authentication service (existing)
- User database (existing)

## Timeline
- Feature ID: F-SAMPLE
- Iteration: 1
