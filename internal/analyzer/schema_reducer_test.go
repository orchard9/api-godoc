package analyzer

import (
	"strings"
	"testing"

	"github.com/orchard9/api-godoc/internal/parser"
)

func TestSchemaReduction(t *testing.T) {
	tests := []struct {
		name           string
		schema         *parser.Schema
		level          string
		wantFields     []string // Fields that should be present
		dontWantFields []string // Fields that should be filtered out
	}{
		{
			name: "essential level - keep only required and key fields",
			schema: &parser.Schema{
				Type:     "object",
				Required: []string{"id", "name"},
				Properties: map[string]parser.Schema{
					"id":          {Type: "string", Description: "Unique identifier"},
					"name":        {Type: "string", Description: "Display name"},
					"email":       {Type: "string", Format: "email"},
					"created_at":  {Type: "string", Format: "date-time"},
					"updated_at":  {Type: "string", Format: "date-time"},
					"deleted_at":  {Type: "string", Format: "date-time"},
					"metadata":    {Type: "object"},
					"internal_id": {Type: "string"},
				},
			},
			level:          "essential",
			wantFields:     []string{"id", "name", "email"}, // required + identifying fields
			dontWantFields: []string{"created_at", "updated_at", "deleted_at", "metadata", "internal_id"},
		},
		{
			name: "standard level - keep important business fields",
			schema: &parser.Schema{
				Type:     "object",
				Required: []string{"id", "title"},
				Properties: map[string]parser.Schema{
					"id":          {Type: "string"},
					"title":       {Type: "string"},
					"description": {Type: "string"},
					"status":      {Type: "string", Enum: []interface{}{"draft", "published"}},
					"price":       {Type: "number"},
					"created_at":  {Type: "string", Format: "date-time"},
					"updated_at":  {Type: "string", Format: "date-time"},
					"_links":      {Type: "object"},
					"_embedded":   {Type: "object"},
				},
			},
			level:          "standard",
			wantFields:     []string{"id", "title", "description", "status", "price"},
			dontWantFields: []string{"created_at", "updated_at", "_links", "_embedded"},
		},
		{
			name: "full level - keep all fields",
			schema: &parser.Schema{
				Type: "object",
				Properties: map[string]parser.Schema{
					"id":         {Type: "string"},
					"name":       {Type: "string"},
					"created_at": {Type: "string", Format: "date-time"},
					"metadata":   {Type: "object"},
				},
			},
			level:      "full",
			wantFields: []string{"id", "name", "created_at", "metadata"},
		},
		{
			name: "nested object reduction",
			schema: &parser.Schema{
				Type: "object",
				Properties: map[string]parser.Schema{
					"id": {Type: "string"},
					"user": {
						Type:     "object",
						Required: []string{"id", "name"},
						Properties: map[string]parser.Schema{
							"id":         {Type: "string"},
							"name":       {Type: "string"},
							"email":      {Type: "string"},
							"created_at": {Type: "string", Format: "date-time"},
							"last_login": {Type: "string", Format: "date-time"},
						},
					},
					"created_at": {Type: "string", Format: "date-time"},
				},
			},
			level:          "essential",
			wantFields:     []string{"id", "user", "user.id", "user.name", "user.email"},
			dontWantFields: []string{"created_at", "user.created_at", "user.last_login"},
		},
	}

	reducer := NewSchemaReducer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reduced := reducer.ReduceSchema(tt.schema, tt.level)

			// Check wanted fields are present
			for _, field := range tt.wantFields {
				if !hasField(reduced, field) {
					t.Errorf("Expected field %s not found in reduced schema", field)
				}
			}

			// Check unwanted fields are filtered
			for _, field := range tt.dontWantFields {
				if hasField(reduced, field) {
					t.Errorf("Unwanted field %s found in reduced schema", field)
				}
			}
		})
	}
}

func TestFieldTypeDetection(t *testing.T) {
	tests := []struct {
		name         string
		fieldName    string
		schema       parser.Schema
		expectFilter bool // Should be filtered in essential mode
	}{
		{
			name:         "timestamp field",
			fieldName:    "created_at",
			schema:       parser.Schema{Type: "string", Format: "date-time"},
			expectFilter: true,
		},
		{
			name:         "id field",
			fieldName:    "user_id",
			schema:       parser.Schema{Type: "string"},
			expectFilter: false,
		},
		{
			name:         "email field",
			fieldName:    "email",
			schema:       parser.Schema{Type: "string", Format: "email"},
			expectFilter: false,
		},
		{
			name:         "internal field",
			fieldName:    "_internal_state",
			schema:       parser.Schema{Type: "string"},
			expectFilter: true,
		},
		{
			name:         "metadata field",
			fieldName:    "metadata",
			schema:       parser.Schema{Type: "object"},
			expectFilter: true,
		},
		{
			name:         "business field",
			fieldName:    "product_name",
			schema:       parser.Schema{Type: "string"},
			expectFilter: false,
		},
	}

	reducer := &schemaReducer{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldFilter := reducer.shouldFilterField(tt.fieldName, &tt.schema, false, "essential")
			if shouldFilter != tt.expectFilter {
				t.Errorf("shouldFilterField(%s) = %v, want %v", tt.fieldName, shouldFilter, tt.expectFilter)
			}
		})
	}
}

func TestFieldConversion(t *testing.T) {
	schema := &parser.Schema{
		Type:     "object",
		Required: []string{"id", "name"},
		Properties: map[string]parser.Schema{
			"id": {
				Type:        "string",
				Description: "Unique identifier",
			},
			"name": {
				Type:      "string",
				MinLength: intPtr(1),
				MaxLength: intPtr(100),
			},
			"age": {
				Type:    "integer",
				Minimum: float64Ptr(0),
				Maximum: float64Ptr(150),
			},
			"tags": {
				Type: "array",
				Items: &parser.Schema{
					Type: "string",
				},
			},
		},
	}

	reducer := NewSchemaReducer()
	fields := reducer.SchemaToFields(schema, "standard")

	// Check field count
	if len(fields) != 4 {
		t.Errorf("Expected 4 fields, got %d", len(fields))
	}

	// Check specific fields
	for _, field := range fields {
		switch field.Name {
		case "id":
			if !field.Required {
				t.Error("id field should be required")
			}
			if field.Type.Type != "string" {
				t.Errorf("id field type = %s, want string", field.Type.Type)
			}
		case "name":
			if !field.Required {
				t.Error("name field should be required")
			}
		case "age":
			if field.Type.Type != "integer" {
				t.Errorf("age field type = %s, want integer", field.Type.Type)
			}
		case "tags":
			if field.Type.Type != "array" {
				t.Errorf("tags field type = %s, want array", field.Type.Type)
			}
			if field.Type.Items == nil || field.Type.Items.Type != "string" {
				t.Error("tags field should have string items")
			}
		}
	}
}

// Helper function to check if a field exists in schema
func hasField(schema *parser.Schema, fieldPath string) bool {
	if schema == nil || schema.Properties == nil {
		return false
	}

	// Handle nested fields
	parts := strings.Split(fieldPath, ".")

	current := schema
	for i, part := range parts {
		if prop, ok := current.Properties[part]; ok {
			if i == len(parts)-1 {
				return true
			}
			current = &prop
		} else {
			return false
		}
	}

	return false
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
