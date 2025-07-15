package analyzer

import (
	"strings"

	"github.com/orchard9/api-godoc/internal/parser"
	"github.com/orchard9/api-godoc/pkg/models"
)

// SchemaReducer reduces schema complexity by filtering fields
type SchemaReducer interface {
	ReduceSchema(schema *parser.Schema, level string) *parser.Schema
	SchemaToFields(schema *parser.Schema, level string) []models.Field
}

// schemaReducer implements the SchemaReducer interface
type schemaReducer struct{}

// NewSchemaReducer creates a new schema reducer
func NewSchemaReducer() SchemaReducer {
	return &schemaReducer{}
}

// ReduceSchema filters schema fields based on reduction level
func (sr *schemaReducer) ReduceSchema(schema *parser.Schema, level string) *parser.Schema {
	if schema == nil || level == "full" {
		return schema
	}

	// Create a deep copy to avoid modifying the original
	reduced := sr.copySchema(schema)

	// Apply reduction
	sr.reduceSchemaRecursive(reduced, level)

	return reduced
}

// SchemaToFields converts a schema to a list of fields with reduction
func (sr *schemaReducer) SchemaToFields(schema *parser.Schema, level string) []models.Field {
	if schema == nil {
		return nil
	}

	var fields []models.Field

	// First reduce the schema if needed
	reducedSchema := sr.ReduceSchema(schema, level)

	// Then convert to fields
	sr.extractFields(reducedSchema, "", &fields)

	return fields
}

// reduceSchemaRecursive applies reduction recursively
func (sr *schemaReducer) reduceSchemaRecursive(schema *parser.Schema, level string) {
	if schema == nil || schema.Properties == nil {
		return
	}

	// Filter properties
	filteredProps := make(map[string]parser.Schema)
	for name, prop := range schema.Properties {
		isRequired := sr.isRequired(name, schema.Required)
		if !sr.shouldFilterField(name, &prop, isRequired, level) {
			propCopy := prop
			// Recursively reduce nested objects
			if prop.Type == "object" {
				sr.reduceSchemaRecursive(&propCopy, level)
			}
			filteredProps[name] = propCopy
		}
	}

	schema.Properties = filteredProps
}

// shouldFilterField determines if a field should be filtered out
func (sr *schemaReducer) shouldFilterField(name string, schema *parser.Schema, isRequired bool, level string) bool {
	// Always keep required fields
	if isRequired {
		return false
	}

	switch level {
	case "essential":
		// Keep only essential fields
		return !sr.isEssentialField(name, schema)
	case "standard":
		// Filter out technical/internal fields
		return sr.isTechnicalField(name, schema)
	default:
		return false
	}
}

// isEssentialField checks if a field is essential
func (sr *schemaReducer) isEssentialField(name string, schema *parser.Schema) bool {
	nameLower := strings.ToLower(name)

	// Always keep ID fields (but not internal ones)
	if !strings.HasPrefix(nameLower, "internal") &&
		(strings.Contains(nameLower, "id") || nameLower == "uuid") {
		return true
	}

	// Keep name/title fields
	if nameLower == "name" || nameLower == "title" || nameLower == "label" {
		return true
	}

	// Keep email fields
	if nameLower == "email" || schema.Format == "email" {
		return true
	}

	// Keep type/status/state fields
	if nameLower == "type" || nameLower == "status" || nameLower == "state" {
		return true
	}

	// Keep key business fields
	if strings.HasSuffix(nameLower, "_name") || strings.HasSuffix(nameLower, "_type") {
		return true
	}

	// Keep objects that have essential nested fields
	if schema.Type == "object" && schema.Properties != nil {
		for childName, childSchema := range schema.Properties {
			if sr.isRequired(childName, schema.Required) || sr.isEssentialField(childName, &childSchema) {
				return true
			}
		}
	}

	return false
}

// isTechnicalField checks if a field is technical/internal
func (sr *schemaReducer) isTechnicalField(name string, schema *parser.Schema) bool {
	nameLower := strings.ToLower(name)

	// Filter timestamp fields
	if sr.isTimestampField(name, schema) {
		return true
	}

	// Filter internal fields (starting with _)
	if strings.HasPrefix(name, "_") {
		return true
	}

	// Filter metadata fields
	if nameLower == "metadata" || nameLower == "meta" {
		return true
	}

	// Filter links/embedded fields (HAL/HATEOAS)
	if nameLower == "links" || nameLower == "embedded" ||
		name == "_links" || name == "_embedded" {
		return true
	}

	// Filter version/revision fields
	if nameLower == "version" || nameLower == "revision" ||
		nameLower == "etag" || nameLower == "__v" {
		return true
	}

	return false
}

// isTimestampField checks if a field is a timestamp
func (sr *schemaReducer) isTimestampField(name string, schema *parser.Schema) bool {
	nameLower := strings.ToLower(name)

	// Check format
	if schema.Format == "date-time" || schema.Format == "date" {
		return true
	}

	// Check common timestamp names
	timestampSuffixes := []string{"_at", "_on", "_date", "_time"}
	for _, suffix := range timestampSuffixes {
		if strings.HasSuffix(nameLower, suffix) {
			return true
		}
	}

	// Check common timestamp names
	timestampNames := []string{"created", "updated", "modified", "deleted", "timestamp"}
	for _, tsName := range timestampNames {
		if strings.Contains(nameLower, tsName) {
			return true
		}
	}

	return false
}

// extractFields recursively extracts fields from schema
func (sr *schemaReducer) extractFields(schema *parser.Schema, prefix string, fields *[]models.Field) {
	if schema == nil {
		return
	}

	// Handle object properties
	if schema.Type == "object" && schema.Properties != nil {
		for name, prop := range schema.Properties {
			fieldName := name
			if prefix != "" {
				fieldName = prefix + "." + name
			}

			field := models.Field{
				Name:        fieldName,
				Type:        sr.buildFieldType(&prop),
				Required:    sr.isRequired(name, schema.Required),
				Description: prop.Description,
			}

			*fields = append(*fields, field)

			// Recursively extract nested object fields
			if prop.Type == "object" {
				sr.extractFields(&prop, fieldName, fields)
			}
		}
	}
}

// buildFieldType builds a FieldType from schema
func (sr *schemaReducer) buildFieldType(schema *parser.Schema) models.FieldType {
	fieldType := models.FieldType{
		Type:   schema.Type,
		Format: schema.Format,
	}

	// Handle array types
	if schema.Type == "array" && schema.Items != nil {
		itemType := sr.buildFieldType(schema.Items)
		fieldType.Items = &itemType
	}

	// Handle enum
	if len(schema.Enum) > 0 {
		for _, e := range schema.Enum {
			if str, ok := e.(string); ok {
				fieldType.Enum = append(fieldType.Enum, str)
			}
		}
	}

	// Handle $ref
	if schema.Ref != "" {
		fieldType.Reference = schema.Ref
		// Extract type name from ref
		parts := strings.Split(schema.Ref, "/")
		if len(parts) > 0 {
			fieldType.Type = parts[len(parts)-1]
		}
	}

	// Default to object if no type
	if fieldType.Type == "" {
		fieldType.Type = "object"
	}

	return fieldType
}

// isRequired checks if a field name is in the required list
func (sr *schemaReducer) isRequired(fieldName string, required []string) bool {
	for _, req := range required {
		if req == fieldName {
			return true
		}
	}
	return false
}

// copySchema creates a deep copy of a schema
func (sr *schemaReducer) copySchema(schema *parser.Schema) *parser.Schema {
	if schema == nil {
		return nil
	}

	copied := *schema

	// Deep copy properties
	if schema.Properties != nil {
		copied.Properties = make(map[string]parser.Schema)
		for k, v := range schema.Properties {
			copied.Properties[k] = *sr.copySchema(&v)
		}
	}

	// Deep copy required array
	if schema.Required != nil {
		copied.Required = make([]string, len(schema.Required))
		copy(copied.Required, schema.Required)
	}

	// Deep copy items
	if schema.Items != nil {
		copied.Items = sr.copySchema(schema.Items)
	}

	return &copied
}
