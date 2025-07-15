package analyzer

import (
	"testing"

	"github.com/orchard9/pg-goapi/internal/parser"
	"github.com/orchard9/pg-goapi/pkg/models"
)

func TestRelationshipDetection(t *testing.T) {
	tests := []struct {
		name     string
		spec     *parser.OpenAPISpec
		expected map[string][]models.Relationship
	}{
		{
			name: "path hierarchy relationships",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/users/{userId}/posts": {
						Get: &parser.Operation{Summary: "Get user posts"},
					},
					"/users/{userId}/posts/{postId}": {
						Get: &parser.Operation{Summary: "Get specific post"},
					},
				},
			},
			expected: map[string][]models.Relationship{
				"users": {
					{Resource: "posts", Type: "has_many", Via: "path hierarchy", Strength: "strong"},
				},
				"posts": {
					{Resource: "users", Type: "belongs_to", Via: "path hierarchy", Strength: "strong"},
				},
			},
		},
		{
			name: "parameter-based relationships",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/orders/{orderId}": {
						Get: &parser.Operation{Summary: "Get order"},
					},
					"/products/{productId}": {
						Get: &parser.Operation{Summary: "Get product"},
					},
					"/orders/{orderId}/items/{itemId}": {
						Get: &parser.Operation{Summary: "Get order item"},
					},
				},
			},
			expected: map[string][]models.Relationship{
				"orders": {
					{Resource: "items", Type: "has_many", Via: "path hierarchy", Strength: "strong"},
				},
				"items": {
					{Resource: "orders", Type: "belongs_to", Via: "path hierarchy", Strength: "strong"},
				},
			},
		},
		{
			name: "foreign key parameter relationships",
			spec: &parser.OpenAPISpec{
				Paths: map[string]parser.PathItem{
					"/users/{userId}": {
						Get: &parser.Operation{Summary: "Get user"},
					},
					"/projects/{projectId}": {
						Get: &parser.Operation{Summary: "Get project"},
					},
					"/tasks/{taskId}": {
						Get: &parser.Operation{Summary: "Get task"},
					},
				},
			},
			expected: map[string][]models.Relationship{
				// Foreign key relationships would need parameter analysis
				// For this test, we'll expect no relationships since the paths don't show clear FK patterns
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up analyzer
			resourceAnalyzer := NewResourceAnalyzer()
			relationshipDetector := NewRelationshipDetector()

			// Extract resources first
			resources := resourceAnalyzer.ExtractResources(tt.spec)

			// Detect relationships
			relationshipDetector.DetectRelationships(resources, tt.spec)

			// Verify relationships
			for _, resource := range resources {
				expectedRels, hasExpected := tt.expected[resource.Name]
				if !hasExpected {
					if len(resource.Relationships) > 0 {
						t.Errorf("Resource %s: expected no relationships, got %d", resource.Name, len(resource.Relationships))
					}
					continue
				}

				if len(resource.Relationships) != len(expectedRels) {
					t.Errorf("Resource %s: expected %d relationships, got %d", resource.Name, len(expectedRels), len(resource.Relationships))
					continue
				}

				// Check each expected relationship
				for _, expectedRel := range expectedRels {
					found := false
					for _, actualRel := range resource.Relationships {
						if actualRel.Resource == expectedRel.Resource &&
							actualRel.Type == expectedRel.Type &&
							actualRel.Via == expectedRel.Via &&
							actualRel.Strength == expectedRel.Strength {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Resource %s: missing expected relationship %+v", resource.Name, expectedRel)
					}
				}
			}
		})
	}
}

func TestPathSegmentExtraction(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected []PathSegment
	}{
		{
			name: "simple resource path",
			path: "/users",
			expected: []PathSegment{
				{Value: "users", IsParameter: false, Position: 1},
			},
		},
		{
			name: "path with parameter",
			path: "/users/{userId}",
			expected: []PathSegment{
				{Value: "users", IsParameter: false, Position: 1},
				{Value: "{userId}", IsParameter: true, ParamName: "userId", Position: 2},
			},
		},
		{
			name: "nested resource path",
			path: "/users/{userId}/posts/{postId}",
			expected: []PathSegment{
				{Value: "users", IsParameter: false, Position: 1},
				{Value: "{userId}", IsParameter: true, ParamName: "userId", Position: 2},
				{Value: "posts", IsParameter: false, Position: 3},
				{Value: "{postId}", IsParameter: true, ParamName: "postId", Position: 4},
			},
		},
		{
			name: "complex nested path",
			path: "/api/v1/projects/{projectId}/tasks/{taskId}/comments",
			expected: []PathSegment{
				{Value: "api", IsParameter: false, Position: 1},
				{Value: "v1", IsParameter: false, Position: 2},
				{Value: "projects", IsParameter: false, Position: 3},
				{Value: "{projectId}", IsParameter: true, ParamName: "projectId", Position: 4},
				{Value: "tasks", IsParameter: false, Position: 5},
				{Value: "{taskId}", IsParameter: true, ParamName: "taskId", Position: 6},
				{Value: "comments", IsParameter: false, Position: 7},
			},
		},
	}

	detector := NewRelationshipDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segments := detector.extractPathSegments(tt.path)

			if len(segments) != len(tt.expected) {
				t.Errorf("Expected %d segments, got %d", len(tt.expected), len(segments))
				return
			}

			for i, expected := range tt.expected {
				actual := segments[i]
				if actual.Value != expected.Value ||
					actual.IsParameter != expected.IsParameter ||
					actual.ParamName != expected.ParamName ||
					actual.Position != expected.Position {
					t.Errorf("Segment %d: expected %+v, got %+v", i, expected, actual)
				}
			}
		})
	}
}

func TestForeignKeyDetection(t *testing.T) {
	tests := []struct {
		name           string
		paramName      string
		expectedMatch  bool
		expectedTarget string
	}{
		{
			name:           "user ID with camelCase",
			paramName:      "userId",
			expectedMatch:  true,
			expectedTarget: "user",
		},
		{
			name:           "project ID with snake_case",
			paramName:      "project_id",
			expectedMatch:  true,
			expectedTarget: "project",
		},
		{
			name:           "simple parameter",
			paramName:      "name",
			expectedMatch:  false,
			expectedTarget: "",
		},
		{
			name:           "ID without prefix",
			paramName:      "id",
			expectedMatch:  false,
			expectedTarget: "",
		},
		{
			name:           "complex foreign key",
			paramName:      "organizationId",
			expectedMatch:  true,
			expectedTarget: "organization",
		},
	}

	detector := NewRelationshipDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := detector.foreignKeyRegex.FindStringSubmatch(tt.paramName)
			hasMatch := len(matches) > 1

			if hasMatch != tt.expectedMatch {
				t.Errorf("Expected match: %v, got: %v for param %s", tt.expectedMatch, hasMatch, tt.paramName)
				return
			}

			if hasMatch && tt.expectedMatch {
				target := matches[1]
				if target != tt.expectedTarget {
					t.Errorf("Expected target: %s, got: %s", tt.expectedTarget, target)
				}
			}
		})
	}
}

func TestSchemaNameConversion(t *testing.T) {
	tests := []struct {
		name         string
		schemaName   string
		expectedName string
	}{
		{
			name:         "simple schema",
			schemaName:   "User",
			expectedName: "user",
		},
		{
			name:         "schema with Response suffix",
			schemaName:   "UserResponse",
			expectedName: "user",
		},
		{
			name:         "schema with Request suffix",
			schemaName:   "CreateProjectRequest",
			expectedName: "createproject",
		},
		{
			name:         "schema with Model suffix",
			schemaName:   "TaskModel",
			expectedName: "task",
		},
		{
			name:         "schema with Schema suffix",
			schemaName:   "OrderSchema",
			expectedName: "order",
		},
		{
			name:         "schema with DTO suffix",
			schemaName:   "ProductDTO",
			expectedName: "product",
		},
	}

	detector := NewRelationshipDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.schemaToResourceName(tt.schemaName)
			if result != tt.expectedName {
				t.Errorf("Expected: %s, got: %s", tt.expectedName, result)
			}
		})
	}
}
