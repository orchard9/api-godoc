package models

import (
	"testing"
	"time"
)

func TestFieldValidation(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  string
	}{
		{
			name: "Required string field",
			field: Field{
				Name: "email",
				Type: FieldType{
					Type:   "string",
					Format: "email",
				},
				Required: true,
			},
			want: "email",
		},
		{
			name: "Optional integer with constraints",
			field: Field{
				Name: "age",
				Type: FieldType{
					Type:    "integer",
					Minimum: floatPtr(0),
					Maximum: floatPtr(150),
				},
			},
			want: "age",
		},
		{
			name: "Array field",
			field: Field{
				Name: "tags",
				Type: FieldType{
					Type: "array",
					Items: &FieldType{
						Type: "string",
					},
				},
			},
			want: "tags",
		},
		{
			name: "Enum field",
			field: Field{
				Name: "status",
				Type: FieldType{
					Type: "string",
					Enum: []string{"active", "inactive", "pending"},
				},
			},
			want: "status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.Name; got != tt.want {
				t.Errorf("Field.Name = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationshipTypes(t *testing.T) {
	tests := []struct {
		name string
		rel  Relationship
		want string
	}{
		{
			name: "Has many relationship",
			rel: Relationship{
				Resource: "Order",
				Type:     "has_many",
				Strength: "strong",
			},
			want: "has_many",
		},
		{
			name: "Belongs to relationship",
			rel: Relationship{
				Resource: "User",
				Type:     "belongs_to",
				Strength: "strong",
			},
			want: "belongs_to",
		},
		{
			name: "References relationship",
			rel: Relationship{
				Resource: "Product",
				Type:     "references",
				Strength: "medium",
			},
			want: "references",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rel.Type; got != tt.want {
				t.Errorf("Relationship.Type = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatternConfidence(t *testing.T) {
	tests := []struct {
		name    string
		pattern Pattern
		want    string
	}{
		{
			name: "High confidence pagination",
			pattern: Pattern{
				Type:       "pagination",
				Confidence: "high",
				Examples:   []string{"/users?page=1", "/orders?limit=10", "/products?offset=20"},
			},
			want: "high",
		},
		{
			name: "Medium confidence filtering",
			pattern: Pattern{
				Type:       "filtering",
				Confidence: "medium",
				Examples:   []string{"/users?status=active", "/orders?date_from=2024-01-01"},
			},
			want: "medium",
		},
		{
			name: "Low confidence versioning",
			pattern: Pattern{
				Type:       "versioning",
				Confidence: "low",
				Examples:   []string{"/v1/users"},
			},
			want: "low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pattern.Confidence; got != tt.want {
				t.Errorf("Pattern.Confidence = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIAnalysisSummary(t *testing.T) {
	analysis := &APIAnalysis{
		Title:       "Test API",
		Version:     "1.0.0",
		Description: "Test API for unit testing",
		BaseURL:     "https://api.example.com",
		SpecType:    "OpenAPI 3.0.0",
		GeneratedAt: time.Now(),
		Resources: []Resource{
			{
				Name: "User",
				Operations: []Operation{
					{Method: "GET", Path: "/users"},
					{Method: "POST", Path: "/users"},
				},
			},
			{
				Name: "Order",
				Operations: []Operation{
					{Method: "GET", Path: "/orders"},
				},
			},
		},
		Summary: AnalysisStat{
			TotalResources:   2,
			TotalOperations:  3,
			TotalEndpoints:   10,
			ResourceCoverage: 30,
		},
	}

	if analysis.Summary.TotalResources != 2 {
		t.Errorf("Expected 2 resources, got %d", analysis.Summary.TotalResources)
	}

	if analysis.Summary.TotalOperations != 3 {
		t.Errorf("Expected 3 operations, got %d", analysis.Summary.TotalOperations)
	}

	if len(analysis.Resources) != 2 {
		t.Errorf("Expected 2 resources in array, got %d", len(analysis.Resources))
	}
}

func TestOperationSecurity(t *testing.T) {
	op := Operation{
		Method:      "POST",
		Path:        "/orders",
		Summary:     "Create order",
		Description: "Creates a new order",
		Security: []string{
			"BearerAuth",
			"ApiKeyAuth",
		},
	}

	if len(op.Security) != 2 {
		t.Errorf("Expected 2 security schemes, got %d", len(op.Security))
	}

	if op.Security[0] != "BearerAuth" {
		t.Errorf("Expected BearerAuth security scheme, got %s", op.Security[0])
	}
}

// Helper function
func floatPtr(f float64) *float64 {
	return &f
}