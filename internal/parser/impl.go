// Package parser provides OpenAPI specification parsing functionality
package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/orchard9/api-godoc/internal/converter"
	"gopkg.in/yaml.v3"
)

// New creates a new parser instance with enhanced validation
func New() Parser {
	return NewEnhanced()
}

// NewBasic creates a basic parser instance (for fallback scenarios)
func NewBasic() Parser {
	return &parser{
		converter: converter.New(),
	}
}

type parser struct {
	converter converter.Converter
}

// ParseFile parses an OpenAPI spec from a file path
func (p *parser) ParseFile(path string) (*OpenAPISpec, error) {
	data, err := os.ReadFile(path) // #nosec G304 - CLI tool, user controls file path
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return p.Parse(data)
}

// ParseURL parses an OpenAPI spec from a URL
func (p *parser) ParseURL(url string) (*OpenAPISpec, error) {
	resp, err := http.Get(url) // #nosec G107 - CLI tool, user controls URL
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return p.Parse(data)
}

// ParseStdin parses an OpenAPI spec from standard input
func (p *parser) ParseStdin() (*OpenAPISpec, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %w", err)
	}

	return p.Parse(data)
}

// Parse parses an OpenAPI spec from raw bytes
func (p *parser) Parse(data []byte) (*OpenAPISpec, error) {
	// Detect format and version
	format, version, err := p.detectFormat(data)
	if err != nil {
		return nil, fmt.Errorf("failed to detect format: %w", err)
	}

	// Handle version-specific logic
	if strings.HasPrefix(version, "2.") {
		// Convert Swagger 2.0 to OpenAPI 3.x if needed
		if version != "2.0" {
			return nil, fmt.Errorf("unsupported Swagger version: %s", version)
		}
		// Convert YAML to JSON if needed before conversion
		if format == "yaml" {
			var yamlDoc interface{}
			if err := yaml.Unmarshal(data, &yamlDoc); err != nil {
				return nil, fmt.Errorf("failed to parse YAML for conversion: %w", err)
			}
			data, err = json.Marshal(yamlDoc)
			if err != nil {
				return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
			}
		}

		data, err = p.converter.Convert(data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert Swagger 2.0: %w", err)
		}
		// Re-detect format after conversion (always JSON)
		format = "json"
	} else if !strings.HasPrefix(version, "3.") {
		return nil, fmt.Errorf("unsupported OpenAPI version: %s", version)
	}

	// Parse based on format
	var spec OpenAPISpec
	switch format {
	case "json":
		if err := json.Unmarshal(data, &spec); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	case "yaml":
		if err := yaml.Unmarshal(data, &spec); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	// Validate spec
	if err := p.validateSpec(&spec); err != nil {
		return nil, fmt.Errorf("invalid spec: %w", err)
	}

	return &spec, nil
}

// detectFormat detects the format (JSON/YAML) and version of the spec
func (p *parser) detectFormat(data []byte) (format string, version string, err error) {
	// Try JSON first
	var jsonDoc map[string]interface{}
	if err := json.Unmarshal(data, &jsonDoc); err == nil {
		format = "json"

		// Check for version
		if v, ok := jsonDoc["openapi"].(string); ok {
			version = v
		} else if v, ok := jsonDoc["swagger"].(string); ok {
			version = v
		} else {
			return "", "", fmt.Errorf("missing version field (openapi or swagger)")
		}

		return format, version, nil
	}

	// Try YAML
	var yamlDoc map[string]interface{}
	if err := yaml.Unmarshal(data, &yamlDoc); err == nil {
		format = "yaml"

		// Check for version
		if v, ok := yamlDoc["openapi"].(string); ok {
			version = v
		} else if v, ok := yamlDoc["swagger"].(string); ok {
			version = v
		} else {
			return "", "", fmt.Errorf("missing version field (openapi or swagger)")
		}

		return format, version, nil
	}

	return "", "", fmt.Errorf("unable to parse as JSON or YAML")
}

// validateSpec performs basic validation on the parsed spec
func (p *parser) validateSpec(spec *OpenAPISpec) error {
	if spec.OpenAPI == "" {
		return fmt.Errorf("missing openapi version")
	}

	if !strings.HasPrefix(spec.OpenAPI, "3.") {
		return fmt.Errorf("unsupported OpenAPI version: %s", spec.OpenAPI)
	}

	if spec.Info.Title == "" {
		return fmt.Errorf("missing info.title")
	}

	if spec.Info.Version == "" {
		return fmt.Errorf("missing info.version")
	}

	if spec.Paths == nil {
		spec.Paths = make(map[string]PathItem)
	}

	return nil
}
