# PG GoAPI

Automatically generate resource-centric documentation from OpenAPI specifications with relationship diagrams, capability matrices, and semantic analysis.

*This is an AI-written project, built with modern software engineering practices.*

## Current Status

🚧 **Pre-Alpha Development** - This project is in early development. Core functionality is being implemented.

**Implemented:**
- Project structure and documentation
- Vision and architecture design

**In Progress:**
- Makefile and build system
- OpenAPI/Swagger parsing
- Resource analysis engine

**Planned:**
- CLI interface
- Multiple output formats (Markdown, JSON, AI-optimized)
- Relationship detection and visualization

See [tasks.md](.memory/tasks.md) for detailed development progress.

## Features
- Resource-centric documentation
- Relationship visualization
- Pattern recognition
- AI-optimized output
- Single binary distribution

## Quick Start
```bash
go install github.com/orchard9/pg-goapi@latest
pg-goapi api-spec.json
```

## Output Example
PG GoAPI generates documentation containing:
- Resource hierarchy with relationships
- Grouped endpoints by business capability
- Capability matrix showing available operations
- Common patterns (pagination, filtering, auth)
- Mermaid relationship diagrams

📄 **Sample Output** - Will be generated through UAT process once core functionality is implemented

## Installation

### From Source
```bash
git clone https://github.com/orchard9/pg-goapi
cd pg-goapi
go build -o pg-goapi cmd/pg-goapi/main.go
```

### Using Go
```bash
go install github.com/orchard9/pg-goapi@latest
```

## Requirements
- OpenAPI 3.x specifications
- Go 1.21+ (for building from source)

## Documentation
- [Usage Guide](usage.md)
- [Architecture](code_architecture.md)
- [Contributing](contributing.md)
- [Why PG GoAPI?](why.md)

## License
MIT License. See [LICENSE.md](LICENSE.md) for details.

## Support
Report issues at https://github.com/orchard9/pg-goapi/issues
