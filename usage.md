# Usage Guide

## Quick Start

```bash
# Analyze a local OpenAPI spec
api-godoc api-spec.json

# Analyze from URL
api-godoc https://api.example.com/openapi.json

# Output to specific file
api-godoc -o api-docs.md api-spec.json

# Generate AI-optimized format
api-godoc -f ai api-spec.json
```

## Command Line Options

```bash
api-godoc [flags] <spec-file-or-url>
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path | `api-docs.md` |
| `--format` | `-f` | Output format (markdown, json, ai) | `markdown` |
| `--verbose` | `-v` | Enable verbose logging | `false` |
| `--version` | | Show version information | |

### Planned Features

The following features are planned for future releases:
- Schema inclusion options (`--include-schemas`)
- Resource filtering (`--filter`)
- Authentication headers (`-H`)
- Validation options
- Configuration file support

## Output Formats

### Markdown (Default)
Human-readable documentation with:
- Resource hierarchy diagram
- Grouped endpoints by resource
- Relationship mapping
- Common patterns

```bash
api-godoc -f markdown api-spec.json
```

### JSON
Machine-readable format for tooling integration:
```bash
api-godoc -f json -o api-analysis.json api-spec.json
```

### AI-Optimized
Condensed format optimized for LLM context windows:
```bash
api-godoc -f ai -o api-context.txt api-spec.json
```

## Examples

### Basic Usage

```bash
# Analyze a local file
api-godoc my-api.json

# Analyze from URL
api-godoc https://petstore.swagger.io/v2/swagger.json

# Specify output format and file
api-godoc -f json -o analysis.json my-api.json
```

### Working with UAT Examples

```bash
# Analyze the warden authentication API
api-godoc uat/artifacts/warden.v1.swagger.json

# Generate AI-optimized output for forge API
api-godoc -f ai uat/artifacts/forge.swagger.json
```

## Understanding Output

### Resource Grouping
The tool groups API endpoints by their business resources rather than listing them sequentially. For example:

```
Users Resource:
- GET /users (list users)
- POST /users (create user)
- GET /users/{id} (get user)
- PUT /users/{id} (update user)
- DELETE /users/{id} (delete user)
```

### Capability Matrix
Shows available operations for each resource:
- ✓ Operation available
- ✗ Operation not available
- Partial: Limited functionality

### Relationship Detection
Identifies how resources connect to each other based on:
- URL path analysis
- Schema references
- Common patterns

## Best Practices

1. **Start Simple**: Use default settings first, then customize as needed
2. **Version Control**: Commit generated documentation alongside your API specs
3. **Regular Updates**: Regenerate documentation when your API changes
4. **Review Output**: Always review generated documentation for accuracy

## Troubleshooting

### Common Issues

#### File Not Found
Ensure the path to your OpenAPI specification is correct:
```bash
# Use absolute path if relative path fails
api-godoc /full/path/to/api-spec.json
```

#### Invalid OpenAPI Format
The tool currently supports:
- OpenAPI 3.x (native support)
- Swagger 2.0 (requires conversion)

#### Network Errors
When analyzing remote specifications:
```bash
# Download first if network is unreliable
curl https://api.example.com/openapi > local-spec.json
api-godoc local-spec.json
```

### Debug Mode
Use verbose flag for detailed logging:
```bash
api-godoc -v api-spec.json
```

## Next Steps

1. Sample outputs will be generated through UAT process once core functionality is implemented
2. Check [why.md](why.md) to understand the tool's philosophy
3. Read [contributing.md](contributing.md) to help improve the tool