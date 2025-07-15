# Real-World API Fixtures

This directory contains real-world OpenAPI specifications used for testing api-godoc against production-quality APIs.

## Included APIs

### 1. Stripe API
- **File**: stripe-openapi3.json
- **Description**: Payment processing API
- **Characteristics**: 
  - Large specification with 100+ resources
  - Complex nested schemas
  - Extensive use of references
  - Well-documented operations

### 2. GitHub API
- **File**: github-openapi3.json
- **Description**: Developer platform API
- **Characteristics**:
  - RESTful design patterns
  - Clear resource hierarchy
  - Authentication patterns
  - Webhooks and events

### 3. Kubernetes API
- **File**: kubernetes-openapi3.json
- **Description**: Container orchestration API
- **Characteristics**:
  - Very large specification
  - Deep resource nesting
  - Complex schema definitions
  - Multiple API versions

## Usage

These fixtures are used in:
- Integration tests to validate parsing
- Performance benchmarks
- Feature development testing
- Documentation examples

## Updating Fixtures

To update these fixtures:
1. Download latest OpenAPI specs from official sources
2. Ensure they are OpenAPI 3.x format (convert if needed)
3. Validate the specifications
4. Update this README with any changes

## License

These specifications are provided by their respective owners:
- Stripe API spec is provided under Stripe's terms
- GitHub API spec is provided under GitHub's terms
- Kubernetes API spec is provided under CNCF's terms