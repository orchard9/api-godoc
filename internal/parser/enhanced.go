package parser

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/orchard9/pg-goapi/internal/converter"
)

// enhancedParser provides library-backed OpenAPI parsing with validation
type enhancedParser struct {
	converter converter.Converter
}

// NewEnhanced creates a new enhanced parser with library validation
func NewEnhanced() Parser {
	return &enhancedParser{
		converter: converter.New(),
	}
}

// ParseFile parses an OpenAPI spec from a file path with enhanced validation
func (p *enhancedParser) ParseFile(path string) (*OpenAPISpec, error) {
	// First try to load with go-openapi for validation
	doc, err := loads.Spec(path)
	if err != nil {
		// If it fails, fall back to basic parser (handles Swagger 2.0 conversion)
		return p.fallbackParser().ParseFile(path)
	}

	// Validate that it's a supported OpenAPI version
	spec := doc.Spec()
	if spec.Swagger != "" && spec.Swagger != "2.0" {
		return nil, fmt.Errorf("unsupported Swagger version: %s", spec.Swagger)
	}

	// Convert go-openapi spec to our format
	return p.convertFromGoOpenAPI(spec), nil
}

// ParseURL parses an OpenAPI spec from a URL with enhanced validation
func (p *enhancedParser) ParseURL(url string) (*OpenAPISpec, error) {
	// Try to load with go-openapi first
	doc, err := loads.Spec(url)
	if err != nil {
		// Fall back to manual parsing with conversion
		return p.fallbackParser().ParseURL(url)
	}

	return p.convertFromGoOpenAPI(doc.Spec()), nil
}

// ParseStdin parses an OpenAPI spec from stdin with enhanced validation
func (p *enhancedParser) ParseStdin() (*OpenAPISpec, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %w", err)
	}

	return p.Parse(data)
}

// Parse parses an OpenAPI spec from raw bytes with enhanced validation
func (p *enhancedParser) Parse(data []byte) (*OpenAPISpec, error) {
	// Try to load with go-openapi first
	doc, err := loads.Analyzed(data, "")
	if err != nil {
		// If it fails, fall back to basic parser (handles Swagger 2.0 conversion)
		return p.fallbackParser().Parse(data)
	}

	// Validate that it's a supported OpenAPI version
	spec := doc.Spec()
	if spec.Swagger != "" && spec.Swagger != "2.0" {
		return nil, fmt.Errorf("unsupported Swagger version: %s", spec.Swagger)
	}

	// Note: Basic validation is done by loads.Analyzed
	// Additional validation could be added here if needed

	return p.convertFromGoOpenAPI(spec), nil
}

// fallbackParser returns the basic parser for fallback scenarios
func (p *enhancedParser) fallbackParser() Parser {
	return &parser{
		converter: p.converter,
	}
}

// convertFromGoOpenAPI converts go-openapi spec to our internal format
func (p *enhancedParser) convertFromGoOpenAPI(s *spec.Swagger) *OpenAPISpec {
	result := &OpenAPISpec{
		OpenAPI: s.Swagger,
		Info: Info{
			Title:       s.Info.Title,
			Description: s.Info.Description,
			Version:     s.Info.Version,
		},
		Paths: make(map[string]PathItem),
	}

	// Handle OpenAPI 3.x version
	if s.Swagger == "" || s.Swagger == "2.0" {
		result.OpenAPI = "3.0.3" // Default for converted specs or when using go-openapi
	}

	// Convert info
	if s.Info != nil {
		if s.Info.Contact != nil {
			result.Info.Contact = &Contact{
				Name:  s.Info.Contact.Name,
				URL:   s.Info.Contact.URL,
				Email: s.Info.Contact.Email,
			}
		}

		if s.Info.License != nil {
			result.Info.License = &License{
				Name: s.Info.License.Name,
				URL:  s.Info.License.URL,
			}
		}
	}

	// Convert servers (from host, basePath, schemes)
	if s.Host != "" || s.BasePath != "" {
		servers := p.convertServers(s)
		result.Servers = servers
	}

	// Convert paths
	if s.Paths != nil {
		for path, pathItem := range s.Paths.Paths {
			result.Paths[path] = p.convertPathItem(pathItem)
		}
	}

	// Convert components/definitions
	if s.Definitions != nil || s.Parameters != nil || s.Responses != nil {
		result.Components = p.convertComponents(s)
	}

	// Convert tags
	if len(s.Tags) > 0 {
		result.Tags = make([]Tag, len(s.Tags))
		for i, tag := range s.Tags {
			result.Tags[i] = Tag{
				Name:        tag.Name,
				Description: tag.Description,
			}
			if tag.ExternalDocs != nil {
				result.Tags[i].ExternalDocs = &ExternalDocs{
					Description: tag.ExternalDocs.Description,
					URL:         tag.ExternalDocs.URL,
				}
			}
		}
	}

	return result
}

// convertServers converts host/basePath/schemes to OpenAPI 3.x servers
func (p *enhancedParser) convertServers(s *spec.Swagger) []Server {
	var servers []Server

	host := s.Host
	basePath := s.BasePath
	schemes := s.Schemes

	if host == "" && basePath == "" {
		return servers
	}

	if host == "" {
		host = "localhost"
	}
	if basePath == "" {
		basePath = ""
	}
	if len(schemes) == 0 {
		schemes = []string{"https"}
	}

	for _, scheme := range schemes {
		servers = append(servers, Server{
			URL: fmt.Sprintf("%s://%s%s", scheme, host, basePath),
		})
	}

	return servers
}

// convertPathItem converts spec.PathItem to our PathItem
func (p *enhancedParser) convertPathItem(pathItem spec.PathItem) PathItem {
	result := PathItem{
		// Note: PathItem in go-openapi doesn't have Summary/Description at path level
	}

	if pathItem.Get != nil {
		result.Get = p.convertOperation(pathItem.Get)
	}
	if pathItem.Post != nil {
		result.Post = p.convertOperation(pathItem.Post)
	}
	if pathItem.Put != nil {
		result.Put = p.convertOperation(pathItem.Put)
	}
	if pathItem.Delete != nil {
		result.Delete = p.convertOperation(pathItem.Delete)
	}
	if pathItem.Options != nil {
		result.Options = p.convertOperation(pathItem.Options)
	}
	if pathItem.Head != nil {
		result.Head = p.convertOperation(pathItem.Head)
	}
	if pathItem.Patch != nil {
		result.Patch = p.convertOperation(pathItem.Patch)
	}

	return result
}

// convertOperation converts spec.Operation to our Operation
func (p *enhancedParser) convertOperation(op *spec.Operation) *Operation {
	result := &Operation{
		Tags:        op.Tags,
		Summary:     op.Summary,
		Description: op.Description,
		OperationID: op.ID,
		Deprecated:  op.Deprecated,
		Responses:   make(map[string]Response),
	}

	// Convert responses
	if op.Responses != nil {
		for code, response := range op.Responses.StatusCodeResponses {
			result.Responses[fmt.Sprintf("%d", code)] = p.convertResponse(&response)
		}
		if op.Responses.Default != nil {
			result.Responses["default"] = p.convertResponse(op.Responses.Default)
		}
	}

	return result
}

// convertResponse converts spec.Response to our Response
func (p *enhancedParser) convertResponse(resp *spec.Response) Response {
	result := Response{
		Description: resp.Description,
	}

	// Convert schema if present
	if resp.Schema != nil {
		result.Content = Content{
			"application/json": MediaType{
				Schema: p.convertSchema(resp.Schema),
			},
		}
	}

	return result
}

// convertSchema converts spec.Schema to our Schema
func (p *enhancedParser) convertSchema(s *spec.Schema) *Schema {
	if s == nil {
		return nil
	}

	result := &Schema{
		Format:      s.Format,
		Title:       s.Title,
		Description: s.Description,
		Default:     s.Default,
		Example:     s.Example,
		ReadOnly:    s.ReadOnly,
		Required:    s.Required,
		Enum:        s.Enum,
	}

	// Handle type (spec.Schema.Type is []string)
	if len(s.Type) > 0 {
		result.Type = s.Type[0]
	}

	// Handle $ref
	if s.Ref.String() != "" {
		result.Ref = s.Ref.String()
		// Convert references to OpenAPI 3.x format
		if strings.HasPrefix(result.Ref, "#/definitions/") {
			result.Ref = strings.Replace(result.Ref, "#/definitions/", "#/components/schemas/", 1)
		}
	}

	// Handle properties
	if s.Properties != nil {
		result.Properties = make(map[string]Schema)
		for name, prop := range s.Properties {
			result.Properties[name] = *p.convertSchema(&prop)
		}
	}

	// Handle items
	if s.Items != nil && s.Items.Schema != nil {
		result.Items = p.convertSchema(s.Items.Schema)
	}

	return result
}

// convertComponents converts definitions and other components
func (p *enhancedParser) convertComponents(s *spec.Swagger) *Components {
	result := &Components{
		Schemas: make(map[string]Schema),
	}

	// Convert definitions to schemas
	if s.Definitions != nil {
		for name, schema := range s.Definitions {
			result.Schemas[name] = *p.convertSchema(&schema)
		}
	}

	return result
}
