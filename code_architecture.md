# Code Architecture

## Overview

PG GoAPI follows a pipeline architecture: Parse → Analyze → Extract → Transform → Generate. Each stage is independent, testable, and focused on a single responsibility.

## Core Components

### 1. Parser (`internal/parser`)
Reads and validates OpenAPI specifications.

```go
type Parser interface {
    Parse(source io.Reader) (*openapi3.T, error)
    ParseFile(path string) (*openapi3.T, error)
    ParseURL(url string) (*openapi3.T, error)
}
```

**Responsibilities:**
- Load OpenAPI from file, URL, or stdin
- Validate specification structure
- Handle both JSON and YAML formats
- Resolve external references

### 2. Analyzer (`internal/analyzer`)
Extracts semantic information from raw OpenAPI data.

```go
type Analyzer struct {
    spec *openapi3.T
}

func (a *Analyzer) ExtractResources() []*models.Resource
func (a *Analyzer) DetectRelationships() []*models.Relationship
func (a *Analyzer) IdentifyPatterns() *models.Patterns
```

**Key Algorithms:**
- **Resource Detection**: Groups endpoints by path segments and tags
- **Relationship Mapping**: Identifies parent-child and references
- **Pattern Recognition**: Finds pagination, filtering, auth patterns

### 3. Extractor (`internal/extractor`)
Builds high-level models from analysis results.

```go
type ResourceExtractor interface {
    ExtractFromPaths(paths openapi3.Paths) []*models.Resource
    ExtractFromSchemas(schemas openapi3.Schemas) []*models.Schema
    MergeResourceData(resources []*models.Resource, schemas []*models.Schema)
}
```

**Processing Steps:**
1. Parse path templates (`/users/{id}/orders`)
2. Group by resource noun
3. Identify CRUD operations
4. Extract path parameters
5. Map request/response schemas

### 4. Reporter (`internal/reporter`)
Generates documentation in various formats.

```go
type Reporter interface {
    Generate(analysis *models.APIAnalysis) (string, error)
}

type MarkdownReporter struct{}
type JSONReporter struct{}
type AIReporter struct{}
```

## Data Models

### Resource Model
```go
type Resource struct {
    Name         string
    Description  string
    Path         string
    Operations   []Operation
    Fields       []Field
    Children     []*Resource
    Parent       *Resource
    Relationships []Relationship
}
```

### Operation Model
```go
type Operation struct {
    Method      string
    Path        string
    OperationID string
    Summary     string
    Parameters  []Parameter
    RequestBody *Schema
    Responses   map[int]*Schema
    Security    []SecurityRequirement
}
```

### Relationship Model
```go
type Relationship struct {
    Type        RelationType // HasMany, BelongsTo, HasOne
    Source      *Resource
    Target      *Resource
    Through     string       // Path or field name
    Cardinality string       // 1:1, 1:N, N:N
}
```

## Processing Pipeline

```
OpenAPI Spec
    ↓
[Parser] → Validated Spec
    ↓
[Analyzer] → Raw Analysis
    ↓
[Extractor] → Resource Model
    ↓
[Transformer] → Enhanced Model
    ↓
[Reporter] → Documentation
```

## Key Algorithms

### Resource Detection Algorithm
```go
func detectResourceFromPath(path string) *Resource {
    // 1. Split path into segments
    // 2. Identify resource nouns (skip parameters)
    // 3. Build resource hierarchy
    // 4. Handle special cases (actions, subresources)
}
```

**Example:**
- `/users` → Resource: "users"
- `/users/{id}` → Resource: "users"
- `/users/{id}/orders` → Resource: "orders", Parent: "users"
- `/users/{id}/activate` → Resource: "users", Action: "activate"

### Relationship Detection
```go
func detectRelationships(resources []*Resource, schemas map[string]*Schema) {
    // 1. Analyze paths for parent-child
    // 2. Check schema references (user_id, customer_ref)
    // 3. Identify array fields pointing to other resources
    // 4. Detect many-to-many through naming conventions
}
```

### Pattern Recognition
```go
func recognizePatterns(operations []Operation) *Patterns {
    // 1. Pagination: page, limit, offset parameters
    // 2. Filtering: filter[field] patterns
    // 3. Sorting: sort, order_by parameters
    // 4. Expansion: include, expand parameters
    // 5. Versioning: headers, path, query params
}
```

## Extension Points

### Custom Analyzers
```go
type AnalyzerPlugin interface {
    Name() string
    Analyze(spec *openapi3.T) (interface{}, error)
}

// Example: Security analyzer
type SecurityAnalyzer struct{}
func (s *SecurityAnalyzer) Analyze(spec *openapi3.T) (interface{}, error) {
    // Analyze authentication schemes
    // Detect authorization patterns
    // Find security vulnerabilities
}
```

### Custom Reporters
```go
// Example: GraphQL schema generator
type GraphQLReporter struct{}
func (g *GraphQLReporter) Generate(analysis *models.APIAnalysis) (string, error) {
    // Convert REST resources to GraphQL types
    // Map operations to queries/mutations
    // Generate schema definition
}
```

## Performance Considerations

### Memory Usage
- Stream large specifications
- Lazy-load external references
- Release parsed data after extraction

### Processing Speed
- Concurrent resource analysis
- Cache computed relationships
- Skip unchanged sections in watch mode

### Optimization Strategies
```go
// Parallel processing for independent resources
func analyzeResourcesConcurrent(resources []*Resource) {
    var wg sync.WaitGroup
    for _, resource := range resources {
        wg.Add(1)
        go func(r *Resource) {
            defer wg.Done()
            analyzeResource(r)
        }(resource)
    }
    wg.Wait()
}
```

## Error Handling

### Error Types
```go
type ParseError struct {
    Line   int
    Column int
    Cause  error
}

type AnalysisError struct {
    Resource string
    Operation string
    Cause    error
}

type ValidationError struct {
    Field   string
    Value   interface{}
    Rule    string
}
```

### Error Strategy
1. Parse errors: Fail fast with clear location
2. Analysis errors: Collect and report all issues
3. Generation errors: Provide partial output with warnings

## Testing Strategy

### Unit Tests
- Parser: Valid/invalid specifications
- Analyzer: Resource detection accuracy
- Extractor: Relationship identification
- Reporter: Output format validation

### Integration Tests
- Real-world API specifications
- End-to-end pipeline testing
- Performance benchmarks
- Regression test suite

### Property-Based Tests
```go
func TestResourceDetection(t *testing.T) {
    quick.Check(func(path string) bool {
        resource := detectResourceFromPath(path)
        // Properties that should always hold
        return resource != nil && resource.Name != ""
    }, nil)
}
```

## Configuration

### Runtime Configuration
```go
type Config struct {
    MaxDepth        int
    IncludeSchemas  bool
    ResourceFilter  *regexp.Regexp
    OutputFormat    Format
    Concurrency     int
}
```

### Feature Flags
```go
type Features struct {
    EnableAIOptimization   bool
    EnableGraphQLExport    bool
    EnableAsyncOperations  bool
    EnableWebhookAnalysis  bool
}
```
