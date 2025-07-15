package parser

// Parser defines the interface for OpenAPI specification parsing
type Parser interface {
	// ParseFile parses an OpenAPI spec from a file path
	ParseFile(path string) (*OpenAPISpec, error)

	// ParseURL parses an OpenAPI spec from a URL
	ParseURL(url string) (*OpenAPISpec, error)

	// ParseStdin parses an OpenAPI spec from standard input
	ParseStdin() (*OpenAPISpec, error)

	// Parse parses an OpenAPI spec from raw bytes
	Parse(data []byte) (*OpenAPISpec, error)
}
