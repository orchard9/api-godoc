package analyzer

import (
	"testing"

	"github.com/orchard9/api-godoc/internal/parser"
)

func TestResourceAnalysis(t *testing.T) {
	tests := []struct {
		name      string
		spec      *parser.OpenAPISpec
		wantErr   bool
		wantCount int
	}{
		{
			name: "basic resource detection",
			spec: &parser.OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: parser.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
				Paths: map[string]parser.PathItem{
					"/users": {
						Get: &parser.Operation{
							Summary: "List users",
						},
						Post: &parser.Operation{
							Summary: "Create user",
						},
					},
					"/users/{id}": {
						Get: &parser.Operation{
							Summary: "Get user",
						},
						Put: &parser.Operation{
							Summary: "Update user",
						},
						Delete: &parser.Operation{
							Summary: "Delete user",
						},
					},
				},
			},
			wantErr:   false,
			wantCount: 1, // Should identify "users" resource
		},
		{
			name: "multiple resources",
			spec: &parser.OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: parser.Info{
					Title:   "E-commerce API",
					Version: "1.0.0",
				},
				Paths: map[string]parser.PathItem{
					"/users": {
						Get: &parser.Operation{Summary: "List users"},
					},
					"/users/{id}": {
						Get: &parser.Operation{Summary: "Get user"},
					},
					"/products": {
						Get: &parser.Operation{Summary: "List products"},
					},
					"/products/{id}": {
						Get: &parser.Operation{Summary: "Get product"},
					},
					"/orders": {
						Get: &parser.Operation{Summary: "List orders"},
					},
					"/orders/{id}": {
						Get: &parser.Operation{Summary: "Get order"},
					},
				},
			},
			wantErr:   false,
			wantCount: 3, // Should identify users, products, orders
		},
		{
			name: "nested resources",
			spec: &parser.OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: parser.Info{
					Title:   "Blog API",
					Version: "1.0.0",
				},
				Paths: map[string]parser.PathItem{
					"/users": {
						Get: &parser.Operation{Summary: "List users"},
					},
					"/users/{id}": {
						Get: &parser.Operation{Summary: "Get user"},
					},
					"/users/{id}/posts": {
						Get: &parser.Operation{Summary: "List user posts"},
					},
					"/users/{id}/posts/{post_id}": {
						Get: &parser.Operation{Summary: "Get user post"},
					},
				},
			},
			wantErr:   false,
			wantCount: 2, // Should identify users and posts
		},
		{
			name: "empty spec",
			spec: &parser.OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: parser.Info{
					Title:   "Empty API",
					Version: "1.0.0",
				},
				Paths: map[string]parser.PathItem{},
			},
			wantErr:   false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New()
			result, err := a.Analyze(tt.spec)

			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(result.Resources) != tt.wantCount {
					t.Errorf("Expected %d resources, got %d", tt.wantCount, len(result.Resources))
				}

				// Verify spec info is preserved
				if result.Title != tt.spec.Info.Title {
					t.Errorf("Expected title %s, got %s", tt.spec.Info.Title, result.Title)
				}
				if result.Version != tt.spec.Info.Version {
					t.Errorf("Expected version %s, got %s", tt.spec.Info.Version, result.Version)
				}
			}
		})
	}
}

func TestResourceGrouping(t *testing.T) {
	spec := &parser.OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: parser.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: map[string]parser.PathItem{
			"/users": {
				Get: &parser.Operation{
					Summary: "List users",
				},
				Post: &parser.Operation{
					Summary: "Create user",
				},
			},
			"/users/{id}": {
				Get: &parser.Operation{
					Summary: "Get user",
				},
				Put: &parser.Operation{
					Summary: "Update user",
				},
				Delete: &parser.Operation{
					Summary: "Delete user",
				},
			},
		},
	}

	a := New()
	result, err := a.Analyze(spec)

	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	if len(result.Resources) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(result.Resources))
	}

	resource := result.Resources[0]
	if resource.Name != "users" {
		t.Errorf("Expected resource name 'users', got %s", resource.Name)
	}

	if len(resource.Operations) != 5 {
		t.Errorf("Expected 5 operations, got %d", len(resource.Operations))
	}

	// Verify operations are correctly grouped
	methods := make(map[string]bool)
	for _, op := range resource.Operations {
		methods[op.Method] = true
	}

	expectedMethods := []string{"GET", "POST", "PUT", "DELETE"}
	for _, method := range expectedMethods {
		if !methods[method] {
			t.Errorf("Expected method %s not found", method)
		}
	}
}

func TestResourceNaming(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		expected []string
	}{
		{
			name:     "simple plural",
			paths:    []string{"/users", "/users/{id}"},
			expected: []string{"users"},
		},
		{
			name:     "multiple resources",
			paths:    []string{"/users", "/products", "/orders"},
			expected: []string{"users", "products", "orders"},
		},
		{
			name:     "nested resources",
			paths:    []string{"/users/{id}/posts", "/users/{id}/posts/{post_id}"},
			expected: []string{"users", "posts"},
		},
		{
			name:     "complex nesting",
			paths:    []string{"/organizations/{org_id}/users/{user_id}/projects"},
			expected: []string{"organizations", "users", "projects"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &parser.OpenAPISpec{
				OpenAPI: "3.0.3",
				Info: parser.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
				Paths: make(map[string]parser.PathItem),
			}

			// Add paths to spec
			for _, path := range tt.paths {
				spec.Paths[path] = parser.PathItem{
					Get: &parser.Operation{
						Summary: "Test operation",
					},
				}
			}

			a := New()
			result, err := a.Analyze(spec)

			if err != nil {
				t.Fatalf("Analyze() error = %v", err)
			}

			if len(result.Resources) != len(tt.expected) {
				t.Errorf("Expected %d resources, got %d", len(tt.expected), len(result.Resources))
			}

			// Check resource names
			resourceNames := make(map[string]bool)
			for _, resource := range result.Resources {
				resourceNames[resource.Name] = true
			}

			for _, expected := range tt.expected {
				if !resourceNames[expected] {
					t.Errorf("Expected resource %s not found", expected)
				}
			}
		})
	}
}
