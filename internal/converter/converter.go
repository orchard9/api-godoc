// Package converter handles conversion from Swagger 2.0 to OpenAPI 3.x format
package converter

import (
	"encoding/json"
	"fmt"
)

// Converter defines the interface for Swagger to OpenAPI conversion
type Converter interface {
	// Convert transforms a Swagger 2.0 spec to OpenAPI 3.x format
	Convert(swagger []byte) ([]byte, error)
}

// New creates a new converter instance
func New() Converter {
	return &converter{}
}

type converter struct{}

// Convert implements the Converter interface
func (c *converter) Convert(swagger []byte) ([]byte, error) {
	// Parse Swagger 2.0 spec
	var swaggerSpec map[string]interface{}
	if err := json.Unmarshal(swagger, &swaggerSpec); err != nil {
		return nil, fmt.Errorf("failed to parse swagger spec: %w", err)
	}

	// Validate Swagger version
	version, ok := swaggerSpec["swagger"].(string)
	if !ok || version != "2.0" {
		return nil, fmt.Errorf("invalid swagger version: expected 2.0, got %v", swaggerSpec["swagger"])
	}

	// Validate required fields
	if _, ok := swaggerSpec["info"]; !ok {
		return nil, fmt.Errorf("missing required field: info")
	}

	// Create OpenAPI 3.x spec
	openAPISpec := map[string]interface{}{
		"openapi": "3.0.3",
	}

	// Convert info section
	if info, ok := swaggerSpec["info"].(map[string]interface{}); ok {
		openAPISpec["info"] = info
	}

	// Convert servers from host, basePath, and schemes
	servers := c.convertServers(swaggerSpec)
	if len(servers) > 0 {
		openAPISpec["servers"] = servers
	}

	// Convert paths
	if paths, ok := swaggerSpec["paths"].(map[string]interface{}); ok {
		openAPISpec["paths"] = c.convertPaths(paths)
	} else {
		openAPISpec["paths"] = map[string]interface{}{}
	}

	// Convert components (definitions, securityDefinitions, etc.)
	components := c.convertComponents(swaggerSpec)
	if len(components) > 0 {
		openAPISpec["components"] = components
	}

	// Convert security
	if security, ok := swaggerSpec["security"]; ok {
		openAPISpec["security"] = security
	}

	// Convert tags
	if tags, ok := swaggerSpec["tags"]; ok {
		openAPISpec["tags"] = tags
	}

	// Convert externalDocs
	if docs, ok := swaggerSpec["externalDocs"]; ok {
		openAPISpec["externalDocs"] = docs
	}

	// Marshal to JSON
	result, err := json.MarshalIndent(openAPISpec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal openapi spec: %w", err)
	}

	return result, nil
}
