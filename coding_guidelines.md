# Coding Guidelines

## Core Principles

1. **Correctness**: Accurate analysis over feature count
2. **Clarity**: Obvious code over clever code  
3. **Efficiency**: Fast enough, not fastest possible
4. **Reliability**: Graceful degradation over hard failures

## Go Standards

### Package Structure
```go
// Package analyzer extracts semantic information from OpenAPI specifications.
// It identifies resources, relationships, and patterns to build a high-level
// understanding of the API structure.
package analyzer
```

### Naming Conventions
```go
// Good
type ResourceAnalyzer struct{}
func (r *ResourceAnalyzer) ExtractResources() []*Resource

// Bad
type Analyzer struct{}  // Too generic
func (r *ResourceAnalyzer) GetRes() []*Resource  // Unclear abbreviation
```

### Error Handling
```go
// Always wrap errors with context
resource, err := analyzeResource(path)
if err != nil {
    return nil, fmt.Errorf("failed to analyze resource %s: %w", path, err)
}

// Use sentinel errors for known conditions
var (
    ErrInvalidSpec = errors.New("invalid OpenAPI specification")
    ErrNoResources = errors.New("no resources found in specification")
)
```

### Interface Design
```go
// Small, focused interfaces
type ResourceExtractor interface {
    ExtractFromPath(path string) (*Resource, error)
}

// Not large, kitchen-sink interfaces
type BadAnalyzer interface {
    Parse()
    Validate() 
    Extract()
    Transform()
    Generate()
    // ... 20 more methods
}
```

## Code Organization

### Function Length
- Target: 20-30 lines
- Maximum: 50 lines
- Extract complex logic into well-named helpers

### Dependency Injection
```go
// Good: Testable and flexible
type Analyzer struct {
    parser    Parser
    extractor Extractor
}

func NewAnalyzer(parser Parser, extractor Extractor) *Analyzer {
    return &Analyzer{
        parser:    parser,
        extractor: extractor,
    }
}

// Bad: Hard dependencies
type BadAnalyzer struct{}

func (a *BadAnalyzer) Analyze() {
    parser := NewParser()  // Hard to test
}
```

### Table-Driven Tests
```go
func TestResourceExtraction(t *testing.T) {
    tests := []struct {
        name     string
        path     string
        expected *Resource
    }{
        {
            name: "simple resource",
            path: "/users",
            expected: &Resource{Name: "users"},
        },
        {
            name: "nested resource",
            path: "/users/{id}/orders",
            expected: &Resource{Name: "orders", Parent: "users"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ExtractResource(tt.path)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Performance Guidelines

### Premature Optimization
```go
// Bad: Optimizing before measuring
func processResources(resources []Resource) {
    // Complex concurrent processing for 10 items
}

// Good: Simple first, optimize if needed
func processResources(resources []Resource) {
    for _, r := range resources {
        processResource(r)
    }
}
```

### Memory Allocation
```go
// Good: Preallocate when size is known
resources := make([]*Resource, 0, len(paths))

// Good: Reuse buffers
var buf strings.Builder
for _, item := range items {
    buf.Reset()
    buf.WriteString(item.Name)
    // Use buf.String()
}
```

## OpenAPI Specific Guidelines

### Schema Handling
```go
// Always check for nil references
if param.Schema != nil && param.Schema.Value != nil {
    paramType := param.Schema.Value.Type
}

// Handle both inline and referenced schemas
func resolveSchema(schemaRef *openapi3.SchemaRef) *openapi3.Schema {
    if schemaRef.Value != nil {
        return schemaRef.Value
    }
    if schemaRef.Ref != "" {
        // Resolve reference
    }
    return nil
}
```

### Path Analysis
```go
// Extract clean resource names
func cleanResourceName(segment string) string {
    // Remove common prefixes/suffixes
    segment = strings.TrimPrefix(segment, "v1/")
    segment = strings.TrimPrefix(segment, "api/")
    segment = strings.TrimSuffix(segment, ".json")
    return segment
}
```

## Documentation Standards

### Function Documentation
```go
// ExtractResources analyzes OpenAPI paths and groups them into logical resources.
// It identifies resource hierarchies based on path structure and returns a slice
// of resources with their relationships mapped.
//
// The algorithm:
//  1. Parses each path into segments
//  2. Identifies resource nouns (skipping parameters)
//  3. Builds parent-child relationships
//  4. Groups operations by resource
//
// Example:
//   /users/{id}/orders -> Resource "orders" with parent "users"
func ExtractResources(spec *openapi3.T) ([]*Resource, error) {
    // Implementation
}
```

### Inline Comments
```go
// Good: Explain why, not what
// Skip internal resources to focus on public API
if strings.HasPrefix(resource.Name, "internal_") {
    continue
}

// Bad: Redundant comment
// Increment counter
counter++
```

## Testing Guidelines

### Test Names
```go
// Good: Descriptive test names
func TestExtractResources_WithNestedPaths_ReturnsHierarchy(t *testing.T)
func TestAnalyzer_InvalidSpec_ReturnsError(t *testing.T)

// Bad: Generic names
func TestExtract(t *testing.T)
func TestError(t *testing.T)
```

### Test Coverage
- Minimum: 80% coverage
- Focus on edge cases
- Test error paths
- Include integration tests

### Benchmarks
```go
func BenchmarkResourceExtraction(b *testing.B) {
    spec := loadLargeSpec()
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, _ = ExtractResources(spec)
    }
}
```

## Common Patterns

### Options Pattern
```go
type AnalyzerOption func(*Analyzer)

func WithMaxDepth(depth int) AnalyzerOption {
    return func(a *Analyzer) {
        a.maxDepth = depth
    }
}

func NewAnalyzer(opts ...AnalyzerOption) *Analyzer {
    a := &Analyzer{
        maxDepth: 3,  // default
    }
    for _, opt := range opts {
        opt(a)
    }
    return a
}
```

### Builder Pattern
```go
type ResourceBuilder struct {
    resource *Resource
}

func (b *ResourceBuilder) WithName(name string) *ResourceBuilder {
    b.resource.Name = name
    return b
}

func (b *ResourceBuilder) WithParent(parent string) *ResourceBuilder {
    b.resource.Parent = parent
    return b
}

func (b *ResourceBuilder) Build() *Resource {
    return b.resource
}
```

## Anti-Patterns to Avoid

### Global State
```go
// Bad
var currentSpec *openapi3.T  // Global variable

// Good
type Analyzer struct {
    spec *openapi3.T  // Encapsulated state
}
```

### Panic in Library Code
```go
// Bad
if resource == nil {
    panic("resource cannot be nil")
}

// Good
if resource == nil {
    return fmt.Errorf("resource cannot be nil")
}
```

### Ignored Errors
```go
// Bad
result, _ := processResource(r)  // Error ignored

// Good
result, err := processResource(r)
if err != nil {
    // Handle or propagate error
}
```

## Code Review Checklist

- [ ] No commented-out code
- [ ] No TODO comments without issue numbers
- [ ] Error messages provide context
- [ ] Tests cover happy path and edge cases
- [ ] Documentation explains why, not what
- [ ] No magic numbers (use constants)
- [ ] Interfaces are small and focused
- [ ] Dependencies are injected
- [ ] Performance is measured, not assumed
