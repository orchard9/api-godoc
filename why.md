# Why PG GoAPI?

## The Problem
API documentation is overwhelming. OpenAPI specs for modern APIs span thousands of lines, mixing critical business logic with boilerplate. Developers waste hours hunting through endpoints to understand resource relationships. AI systems choke on verbose specifications. Existing tools focus on pretty rendering, not comprehension.

## The Solution
PG GoAPI automatically generates concise, resource-centric documentation by analyzing OpenAPI specifications. One command transforms sprawling specs into focused documentation that captures what actually matters: resources, relationships, and capabilities.

## Key Benefits
- **Resource-Centric**: Groups by business resources, not HTTP endpoints
- **Relationship Mapping**: Visualizes how resources connect
- **90/10 Rule**: Captures 90% of value in 10% of the content
- **AI-Optimized**: Perfect context size for LLM consumption
- **Zero Dependencies**: Single Go binary, no Node.js or Python
- **CI/CD Ready**: Automate documentation in your pipeline

## Use Cases
- Rapid API onboarding for new developers
- Architecture reviews and API audits
- AI integration preparation
- API change impact analysis
- Client SDK planning
- Microservice dependency mapping

## Why Not Alternatives?

### Swagger UI / Redoc
- Endpoint-centric, not resource-centric
- No relationship visualization
- Overwhelming for large APIs
- Poor for offline/CLI use

### Postman Documentation
- Requires collection maintenance
- Cloud-based (privacy concerns)
- Not git-friendly
- No semantic analysis

### ReadMe.io / Stoplight
- Commercial services
- Manual maintenance required
- No automated analysis
- Focused on presentation, not comprehension

### Custom Scripts
- Maintenance burden
- No standard approach
- Lacks sophisticated analysis
- Reinventing the wheel

## Our Philosophy
APIs represent business capabilities, not just HTTP endpoints. PG GoAPI understands this distinction and generates documentation that reflects how developers actually think about and use APIs. We extract signal from noise, relationships from routes, and meaning from methods.
