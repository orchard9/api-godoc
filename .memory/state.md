# Project State

## Documentation Files
- README.md: Main project documentation and quick start guide
- CLAUDE.md: Development guide for AI-assisted coding with memory management
- LICENSE.md: MIT license text
- code_architecture.md: System design with pipeline architecture and component overview
- coding_guidelines.md: Go coding standards, patterns, and best practices
- contributing.md: Development setup, workflow, and contribution guidelines
- usage.md: Usage guide with basic CLI options (updated to reflect current implementation)
- why.md: Project rationale explaining the problem and solution
- example-output.md: Concrete sample outputs showing markdown, JSON, and AI-optimized formats

## Memory Files
- .memory/vision.md: Long-term project vision and design principles
- .memory/tasks.md: Comprehensive task list with pending and completed items
- .memory/semantic-memory.md: Project facts and core characteristics
- .memory/working-memory.md: Current work progress tracking
- .memory/state.md: This file - current project file listing

## Directories
- .memory/: Memory management files for development context
- uat/: User acceptance testing directory
- uat/artifacts/: Test API specifications for UAT

## Source Files
- internal/analyzer/patterns.go: Pattern detection implementation for common API patterns
- internal/analyzer/schema_reducer.go: Schema reduction implementation for filtering non-essential fields
- internal/analyzer/resource_filter.go: Resource filtering implementation for targeted documentation

## Test Files
- uat/artifacts/warden.v1.swagger.json: Example Swagger 2.0 authentication/authorization API specification
- uat/artifacts/forge.swagger.json: Example Swagger 2.0 AI-driven development platform API specification
- uat/artifacts/README.md: Documentation explaining the UAT artifacts and their usage
- test/fixtures/: Real-world API fixtures directory with production specifications
- test/fixtures/stripe-openapi3.json: Stripe API specification (6.8MB)
- test/fixtures/github-openapi3.json: GitHub API specification (10.8MB)
- test/fixtures/kubernetes-simplified.json: Simplified Kubernetes API specification
- test/fixtures_test.go: Integration tests for real-world API fixtures
- test/benchmark_test.go: Performance benchmarks for production APIs
- test/PERFORMANCE.md: Performance analysis and benchmark results
- internal/reporter/mermaid_test.go: Tests for Mermaid diagram generation functionality
- internal/analyzer/patterns_test.go: Tests for API pattern detection functionality
- internal/analyzer/schema_reducer_test.go: Tests for schema reduction functionality
- internal/analyzer/resource_filter_test.go: Tests for resource filtering functionality
- pkg/models/models_test.go: Tests for data model structures
- cmd/pg-goapi/main_test.go: Tests for CLI functionality
- internal/basic_test.go: Basic tests for analyzer and reporter components

## Project Status
- Pre-Alpha Development phase
- Documentation framework complete and aligned with vision
- Tasks reorganized by priority (26 total tasks)
- Ready to begin Foundation phase: Makefile, Go structure, Swagger converter
- Critical finding: Swagger 2.0 converter needed immediately (UAT examples require it)
