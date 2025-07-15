# UAT Artifacts

This directory contains example API specifications for testing pg-goapi.

## Available Specifications

### warden.v1.swagger.json
- **Format**: Swagger 2.0
- **Description**: Authentication and authorization API with comprehensive user, organization, and OAuth2 management
- **Resources**: 
  - Accounts (users, service accounts)
  - Organizations (with membership management)
  - Authentication (login, register, API keys)
  - OAuth2 (clients, providers, tokens)
  - Admin operations
  - Health checks

This specification demonstrates a complex, real-world API with:
- Multiple resource types with relationships
- Authentication and authorization patterns
- Standard REST operations (CRUD)
- Specialized operations (login, logout, token refresh)
- Pagination patterns
- Rich schema definitions

## Usage

During development and testing:
```bash
pg-goapi uat/artifacts/warden.v1.swagger.json
```

## Notes

- The warden specification is in Swagger 2.0 format, which will test our ability to handle legacy specifications
- Future artifacts may include OpenAPI 3.x specifications for comparison
- Additional test cases can be added as needed
