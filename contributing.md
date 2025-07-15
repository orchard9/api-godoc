# Contributing to PG GoAPI

## Development Setup

### Prerequisites
- Go 1.21+
- Make
- Git

### Quick Start
```bash
git clone https://github.com/orchard9/pg-goapi
cd pg-goapi
make test
make build
```

## Development Workflow

### 1. Check Issues
```bash
gh issue list
gh issue create
```

### 2. Create Branch
```bash
git checkout -b feature/resource-detection
```

### 3. Run Tests Continuously
```bash
make watch  # Runs tests on file changes
```

### 4. Submit PR
- Reference issue number
- Include test coverage
- Update documentation

## Code Standards

### Philosophy
- Correctness over cleverness
- Clarity over brevity
- Performance when needed, not by default

### Guidelines
- Single responsibility per function
- Table-driven tests preferred
- Error wrapping with context
- No panic in library code

### Testing
- Unit tests for analyzers
- Integration tests for real OpenAPI specs
- Property-based tests for parsers
- Minimum 80% coverage

## Project Structure
```
cmd/pg-goapi/      # CLI entry point
internal/
  analyzer/        # OpenAPI parsing and analysis
  extractor/       # Resource and relationship extraction
  reporter/        # Documentation generation
pkg/
  models/          # Core data structures
  openapi/         # OpenAPI helpers
tests/
  fixtures/        # Sample OpenAPI specs
  integration/     # End-to-end tests
```

## Pull Request Process

### Before Submitting
1. Run `make ci` - all checks must pass
2. Add tests for new functionality
3. Update relevant documentation
4. Ensure backward compatibility

### PR Description Template
```
## Summary
Brief description of changes

## Related Issue
Fixes #123

## Changes
- Added resource grouping algorithm
- Improved relationship detection
- Fixed edge case in path parsing

## Testing
- Added unit tests for X
- Tested against Stripe, GitHub, Kubernetes APIs
```

## Testing Real APIs

### Adding Test Fixtures
```bash
# Add new OpenAPI spec to tests/fixtures/
curl https://api.example.com/openapi.json > tests/fixtures/example-api.json

# Create corresponding test
touch internal/analyzer/example_test.go
```

### Integration Testing
```go
func TestRealWorldAPIs(t *testing.T) {
    specs := []string{
        "stripe-api.json",
        "github-api.json",
        "kubernetes-api.json",
    }
    
    for _, spec := range specs {
        t.Run(spec, func(t *testing.T) {
            // Test resource extraction
            // Test relationship detection
            // Test report generation
        })
    }
}
```

## Performance Considerations

### Benchmarking
```bash
make bench
make bench-compare  # Compare with previous commit
```

### Optimization Guidelines
- Profile before optimizing
- Focus on algorithmic improvements
- Minimize allocations in hot paths
- Use streaming for large specs

## Documentation

### Code Comments
- Package-level documentation required
- Public API must be documented
- Complex algorithms need explanation
- Include examples for non-obvious usage

### User Documentation
- Update usage.md for new features
- Add examples to README
- Document breaking changes
- Include migration guides

## Release Process

### Version Bumping
```bash
make version-patch  # 0.1.0 -> 0.1.1
make version-minor  # 0.1.0 -> 0.2.0
make version-major  # 0.1.0 -> 1.0.0
```

### Release Checklist
- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Binary builds successfully
- [ ] Tagged and pushed

## Getting Help

### Communication
- GitHub Issues for bugs/features
- Discussions for questions
- PR comments for code review

### Resources
- [OpenAPI Specification](https://spec.openapis.org/oas/latest.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
