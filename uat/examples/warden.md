# warden/v1/common.proto

## Overview

- **API Version**: version not set
- **Specification Type**: OpenAPI 3.0.3
- **Generated**: 2025-07-16 11:01:05

## API Statistics

- **Total Resources**: 2
- **Total Operations**: 8
- **Total Endpoints**: 6
- **Resource Coverage**: 133%

## Resources

This section groups API endpoints by business resources for better understanding.

### Api-keys

Api-keys resource operations

**Operations**: 4

| Method | Path | Summary |
|--------|------|----------|
| GET | `/v1/auth/api-keys` |  |
| GET | `/v1/auth/api-keys/{id}` |  |
| POST | `/v1/auth/api-keys` | API key management |
| DELETE | `/v1/auth/api-keys/{id}` |  |

### Auth

Auth resource operations

**Operations**: 4

| Method | Path | Summary |
|--------|------|----------|
| POST | `/v1/auth/login` |  |
| POST | `/v1/auth/logout` |  |
| POST | `/v1/auth/refresh` |  |
| POST | `/v1/auth/register` | Authentication methods |

## Detected Patterns

### Versioning

**Confidence**: high  
**Impact**: Clients should be aware of API version compatibility

API uses URL path versioning. Versions found: v1

**Examples**:
- /v1/auth/api-keys

