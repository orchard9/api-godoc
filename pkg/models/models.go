// Package models defines the core data structures for API analysis
package models

import "time"

// APIAnalysis represents the complete analysis of an OpenAPI specification
type APIAnalysis struct {
	Title         string       `json:"title"`
	Version       string       `json:"version"`
	Description   string       `json:"description,omitempty"`
	BaseURL       string       `json:"baseUrl,omitempty"`
	Resources     []Resource   `json:"resources"`
	Patterns      []Pattern    `json:"patterns,omitempty"`
	Summary       AnalysisStat `json:"summary"`
	GeneratedAt   time.Time    `json:"generatedAt"`
	SpecType      string       `json:"specType"` // OpenAPI 3.x, Swagger 2.0
	OriginalPaths int          `json:"originalPaths"`
}

// AnalysisStat provides high-level statistics about the API
type AnalysisStat struct {
	TotalResources   int `json:"totalResources"`
	TotalOperations  int `json:"totalOperations"`
	TotalEndpoints   int `json:"totalEndpoints"`
	ResourceCoverage int `json:"resourceCoverage"` // percentage of paths that map to resources
}

// Resource represents a business resource extracted from the API
type Resource struct {
	Name          string         `json:"name"`
	Description   string         `json:"description,omitempty"`
	Operations    []Operation    `json:"operations"`
	Relationships []Relationship `json:"relationships,omitempty"`
	Fields        []Field        `json:"fields,omitempty"`
	Category      string         `json:"category,omitempty"` // core, admin, utility, etc.
	IsCollection  bool           `json:"isCollection"`       // true if this represents a collection resource
}

// Operation represents an API operation (HTTP method + path)
type Operation struct {
	Method       string       `json:"method"`
	Path         string       `json:"path"`
	Summary      string       `json:"summary,omitempty"`
	Description  string       `json:"description,omitempty"`
	OperationID  string       `json:"operationId,omitempty"`
	Tags         []string     `json:"tags,omitempty"`
	Parameters   []Parameter  `json:"parameters,omitempty"`
	RequestBody  *RequestBody `json:"requestBody,omitempty"`
	Responses    []Response   `json:"responses,omitempty"`
	Security     []string     `json:"security,omitempty"`
	Deprecated   bool         `json:"deprecated,omitempty"`
	IsResourceOp bool         `json:"isResourceOp"` // true if this is a standard CRUD operation
}

// Parameter represents operation parameters
type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"` // query, path, header, cookie
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Example     string `json:"example,omitempty"`
}

// RequestBody represents request body information
type RequestBody struct {
	Description string     `json:"description,omitempty"`
	Required    bool       `json:"required"`
	ContentType string     `json:"contentType"`
	Schema      *FieldType `json:"schema,omitempty"`
}

// Response represents operation responses
type Response struct {
	StatusCode  string     `json:"statusCode"`
	Description string     `json:"description,omitempty"`
	ContentType string     `json:"contentType,omitempty"`
	Schema      *FieldType `json:"schema,omitempty"`
}

// Field represents a schema field or property
type Field struct {
	Name        string    `json:"name"`
	Type        FieldType `json:"type"`
	Description string    `json:"description,omitempty"`
	Required    bool      `json:"required"`
	Example     string    `json:"example,omitempty"`
	Deprecated  bool      `json:"deprecated,omitempty"`
}

// FieldType represents the type information for a field
type FieldType struct {
	Type       string     `json:"type"` // string, integer, object, array, etc.
	Format     string     `json:"format,omitempty"`
	Items      *FieldType `json:"items,omitempty"`      // for array types
	Properties []Field    `json:"properties,omitempty"` // for object types
	Reference  string     `json:"reference,omitempty"`  // $ref value
	Enum       []string   `json:"enum,omitempty"`
	Pattern    string     `json:"pattern,omitempty"`
	MinLength  *int       `json:"minLength,omitempty"`
	MaxLength  *int       `json:"maxLength,omitempty"`
	Minimum    *float64   `json:"minimum,omitempty"`
	Maximum    *float64   `json:"maximum,omitempty"`
}

// Relationship represents a connection between resources
type Relationship struct {
	Resource    string `json:"resource"`              // target resource name
	Type        string `json:"type"`                  // has_many, belongs_to, references, etc.
	Via         string `json:"via,omitempty"`         // field or parameter that creates the relationship
	Description string `json:"description,omitempty"` // human-readable relationship description
	Strength    string `json:"strength"`              // strong, weak, inferred
}

// Pattern represents a detected API pattern
type Pattern struct {
	Type        string   `json:"type"`        // pagination, authentication, versioning, etc.
	Description string   `json:"description"` // human-readable description
	Examples    []string `json:"examples,omitempty"`
	Confidence  string   `json:"confidence"` // high, medium, low
	Impact      string   `json:"impact"`     // affects how the API should be consumed
}
