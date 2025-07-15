# Contributing to API GoDoc

## Development Setup

### Prerequisites
- Go 1.24.x (or 1.23.x for compatibility)
- Make
- Git
- Docker (for local CI testing with ACT)

### Quick Start
```bash
git clone https://github.com/orchard9/api-godoc
cd api-godoc
make ci        # Runs full CI pipeline
make build
make uat       # User acceptance testing
```

### Setting Up Local CI Testing
```bash
# Install ACT for local GitHub Actions testing
curl -q https://raw.githubusercontent.com/nektos/act/master/install.sh | bash

# Test workflows locally before pushing
make act-test  # Run CI workflow locally
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
cmd/api-godoc/      # CLI entry point
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
2. Test locally with `make act-test` (optional but recommended)
3. Add tests for new functionality
4. Update relevant documentation
5. Ensure backward compatibility

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

## CI/CD Workflow

### Continuous Integration
Our CI pipeline runs on every push and pull request:

- **Multi-platform testing**: Linux, Windows, macOS
- **Multi-version testing**: Go 1.24.x and 1.23.x
- **Automated linting**: golangci-lint, go vet, gofmt
- **Security scanning**: gosec for security vulnerabilities
- **Coverage reporting**: CodeCov integration
- **Real-world testing**: UAT with production API specs

### Local CI Testing
```bash
# Install ACT for local testing
make ci-setup

# Test CI workflow locally
make act-test

# Build release binaries locally
make release-local
```

### GitHub Actions Workflows

#### CI Pipeline (`.github/workflows/ci.yml`)
- Triggers on push to main/develop, PRs to main
- Runs tests, linting, security scans
- Matrix testing across OS and Go versions
- Uploads coverage reports

#### Release Pipeline (`.github/workflows/release.yml`)
- Triggers on version tags (`v*.*.*`)
- Builds cross-platform binaries
- Creates GitHub releases with artifacts
- Generates release notes automatically

#### Dependabot Configuration
- Weekly dependency updates
- Automated security updates
- Separate updates for Go modules and GitHub Actions

## Release Process

### Automated Releases
Releases are triggered by pushing version tags:
```bash
git tag v1.0.0
git push origin v1.0.0
```

This automatically:
- Runs full CI pipeline
- Builds binaries for all platforms
- Creates GitHub release with artifacts
- Generates checksums and release notes

### Manual Release Testing
```bash
# Test release process locally
make release-local

# Validate binaries
./dist/api-godoc-linux-amd64 --version
./dist/api-godoc-darwin-arm64 --version
```

### Release Checklist
- [ ] All tests passing in CI
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Tag created and pushed
- [ ] GitHub release artifacts generated
- [ ] Release notes reviewed

## Getting Help

### Communication
- GitHub Issues for bugs/features
- Discussions for questions
- PR comments for code review

### Resources
- [OpenAPI Specification](https://spec.openapis.org/oas/latest.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
