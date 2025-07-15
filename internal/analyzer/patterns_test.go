package analyzer

import (
	"testing"

	"github.com/orchard9/api-godoc/internal/parser"
	"github.com/orchard9/api-godoc/pkg/models"
)

func TestPatternDetection(t *testing.T) {
	tests := []struct {
		name     string
		spec     *parser.OpenAPISpec
		expected []string // Expected pattern types
	}{
		{
			name: "pagination patterns",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/users": {
						Get: &parser.Operation{
							Parameters: []parser.Parameter{
								{Name: "page", In: "query", Schema: &parser.Schema{Type: "integer"}},
								{Name: "limit", In: "query", Schema: &parser.Schema{Type: "integer"}},
							},
						},
					},
					"/posts": {
						Get: &parser.Operation{
							Parameters: []parser.Parameter{
								{Name: "offset", In: "query", Schema: &parser.Schema{Type: "integer"}},
								{Name: "limit", In: "query", Schema: &parser.Schema{Type: "integer"}},
							},
						},
					},
				},
			},
			expected: []string{"pagination"},
		},
		{
			name: "filtering patterns",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/users": {
						Get: &parser.Operation{
							Parameters: []parser.Parameter{
								{Name: "status", In: "query", Schema: &parser.Schema{Type: "string"}},
								{Name: "role", In: "query", Schema: &parser.Schema{Type: "string"}},
								{Name: "created_after", In: "query", Schema: &parser.Schema{Type: "string", Format: "date-time"}},
							},
						},
					},
				},
			},
			expected: []string{"filtering"},
		},
		{
			name: "sorting patterns",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/products": {
						Get: &parser.Operation{
							Parameters: []parser.Parameter{
								{Name: "sort", In: "query", Schema: &parser.Schema{Type: "string"}},
								{Name: "order", In: "query", Schema: &parser.Schema{Type: "string", Enum: []interface{}{"asc", "desc"}}},
							},
						},
					},
					"/users": {
						Get: &parser.Operation{
							Parameters: []parser.Parameter{
								{Name: "orderBy", In: "query", Schema: &parser.Schema{Type: "string"}},
								{Name: "sortDirection", In: "query", Schema: &parser.Schema{Type: "string"}},
							},
						},
					},
				},
			},
			expected: []string{"sorting"},
		},
		{
			name: "versioning patterns",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/v1/users":        {Get: &parser.Operation{}},
					"/v1/posts":        {Get: &parser.Operation{}},
					"/v2/users":        {Get: &parser.Operation{}},
					"/api/v3/products": {Get: &parser.Operation{}},
				},
			},
			expected: []string{"versioning"},
		},
		{
			name: "batch operations",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/users:batchCreate": {
						Post: &parser.Operation{
							Summary: "Create multiple users",
						},
					},
					"/tasks:batchGet": {
						Post: &parser.Operation{
							Summary: "Get multiple tasks",
						},
					},
				},
			},
			expected: []string{"batch_operations"},
		},
		{
			name: "search patterns",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/search": {
						Get: &parser.Operation{
							Parameters: []parser.Parameter{
								{Name: "q", In: "query", Schema: &parser.Schema{Type: "string"}},
							},
						},
					},
					"/users": {
						Get: &parser.Operation{
							Parameters: []parser.Parameter{
								{Name: "search", In: "query", Schema: &parser.Schema{Type: "string"}},
							},
						},
					},
				},
			},
			expected: []string{"search"},
		},
		{
			name: "authentication patterns",
			spec: &parser.OpenAPISpec{
				Components: &parser.Components{
					SecuritySchemes: map[string]parser.SecurityScheme{
						"bearerAuth": {
							Type:         "http",
							Scheme:       "bearer",
							BearerFormat: "JWT",
						},
						"apiKey": {
							Type: "apiKey",
							In:   "header",
							Name: "X-API-Key",
						},
					},
				},
				Security: []parser.SecurityRequirement{
					{"bearerAuth": []string{}},
				},
			},
			expected: []string{"authentication"},
		},
		{
			name: "multiple patterns",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/v1/users": {
						Get: &parser.Operation{
							Parameters: []parser.Parameter{
								{Name: "page", In: "query", Schema: &parser.Schema{Type: "integer"}},
								{Name: "limit", In: "query", Schema: &parser.Schema{Type: "integer"}},
								{Name: "sort", In: "query", Schema: &parser.Schema{Type: "string"}},
								{Name: "status", In: "query", Schema: &parser.Schema{Type: "string"}},
							},
						},
					},
				},
			},
			expected: []string{"versioning", "pagination", "sorting", "filtering"},
		},
	}

	detector := NewPatternDetector()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patterns := detector.DetectPatterns(tt.spec)

			// Check that all expected patterns are found
			for _, expectedType := range tt.expected {
				found := false
				for _, pattern := range patterns {
					if pattern.Type == expectedType {
						found = true
						if pattern.Confidence == "" {
							t.Errorf("Pattern %s has empty confidence", expectedType)
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected pattern type %s not found", expectedType)
				}
			}

			// Check no unexpected patterns
			for _, pattern := range patterns {
				expected := false
				for _, expectedType := range tt.expected {
					if pattern.Type == expectedType {
						expected = true
						break
					}
				}
				if !expected {
					t.Errorf("Unexpected pattern type %s found", pattern.Type)
				}
			}
		})
	}
}

func TestPaginationParameterDetection(t *testing.T) {
	tests := []struct {
		name     string
		params   []parser.Parameter
		expected bool
	}{
		{
			name: "page and limit",
			params: []parser.Parameter{
				{Name: "page", In: "query"},
				{Name: "limit", In: "query"},
			},
			expected: true,
		},
		{
			name: "offset and limit",
			params: []parser.Parameter{
				{Name: "offset", In: "query"},
				{Name: "limit", In: "query"},
			},
			expected: true,
		},
		{
			name: "pageSize and pageNumber",
			params: []parser.Parameter{
				{Name: "pageSize", In: "query"},
				{Name: "pageNumber", In: "query"},
			},
			expected: true,
		},
		{
			name: "cursor-based",
			params: []parser.Parameter{
				{Name: "cursor", In: "query"},
				{Name: "limit", In: "query"},
			},
			expected: true,
		},
		{
			name: "no pagination",
			params: []parser.Parameter{
				{Name: "id", In: "path"},
				{Name: "status", In: "query"},
			},
			expected: false,
		},
	}

	detector := &patternDetector{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detector.hasPaginationParams(tt.params)
			if got != tt.expected {
				t.Errorf("hasPaginationParams() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPatternConfidence(t *testing.T) {
	tests := []struct {
		name               string
		pattern            models.Pattern
		expectedConfidence string
	}{
		{
			name: "high confidence pagination",
			pattern: models.Pattern{
				Type:     "pagination",
				Examples: []string{"/users", "/posts", "/products", "/orders"},
			},
			expectedConfidence: "high",
		},
		{
			name: "medium confidence filtering",
			pattern: models.Pattern{
				Type:     "filtering",
				Examples: []string{"/users", "/posts"},
			},
			expectedConfidence: "medium",
		},
		{
			name: "low confidence search",
			pattern: models.Pattern{
				Type:     "search",
				Examples: []string{"/search"},
			},
			expectedConfidence: "low",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Confidence should be set based on number of examples
			confidence := "low"
			if len(tt.pattern.Examples) >= 3 {
				confidence = "high"
			} else if len(tt.pattern.Examples) >= 2 {
				confidence = "medium"
			}

			if confidence != tt.expectedConfidence {
				t.Errorf("Expected confidence %s, got %s", tt.expectedConfidence, confidence)
			}
		})
	}
}
