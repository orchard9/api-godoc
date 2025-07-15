package analyzer

import (
	"regexp"
	"strings"

	"github.com/orchard9/pg-goapi/internal/parser"
	"github.com/orchard9/pg-goapi/pkg/models"
)

// RelationshipDetector analyzes API specifications to identify relationships between resources
type RelationshipDetector struct {
	pathParamRegex  *regexp.Regexp
	idPatternRegex  *regexp.Regexp
	foreignKeyRegex *regexp.Regexp
}

// NewRelationshipDetector creates a new relationship detector
func NewRelationshipDetector() *RelationshipDetector {
	return &RelationshipDetector{
		pathParamRegex:  regexp.MustCompile(`\{([^}]+)\}`),
		idPatternRegex:  regexp.MustCompile(`^(.+)Id$|^(.+)_id$`),
		foreignKeyRegex: regexp.MustCompile(`^(.+)(Id|_id)$`),
	}
}

// DetectRelationships analyzes resources to identify relationships between them
func (rd *RelationshipDetector) DetectRelationships(resources []models.Resource, spec *parser.OpenAPISpec) {
	resourceMap := make(map[string]*models.Resource)
	for i := range resources {
		resourceMap[resources[i].Name] = &resources[i]
	}

	// Analyze path-based relationships
	rd.detectPathRelationships(resourceMap, spec)

	// Analyze parameter-based relationships
	rd.detectParameterRelationships(resourceMap)

	// Analyze schema-based relationships
	rd.detectSchemaRelationships(resourceMap, spec)
}

// detectPathRelationships identifies relationships from nested paths
func (rd *RelationshipDetector) detectPathRelationships(resourceMap map[string]*models.Resource, spec *parser.OpenAPISpec) {
	for path := range spec.Paths {
		segments := rd.extractPathSegments(path)
		rd.analyzePathHierarchy(segments, resourceMap)
	}
}

// extractPathSegments extracts meaningful segments from a path
func (rd *RelationshipDetector) extractPathSegments(path string) []PathSegment {
	parts := strings.Split(path, "/")
	var segments []PathSegment

	for i, part := range parts {
		if part == "" {
			continue
		}

		if rd.pathParamRegex.MatchString(part) {
			// This is a path parameter
			paramName := rd.pathParamRegex.ReplaceAllString(part, "$1")
			segments = append(segments, PathSegment{
				Value:       part,
				IsParameter: true,
				ParamName:   paramName,
				Position:    i,
			})
		} else {
			// This is a resource name or path segment
			segments = append(segments, PathSegment{
				Value:       part,
				IsParameter: false,
				Position:    i,
			})
		}
	}

	return segments
}

// PathSegment represents a segment of an API path
type PathSegment struct {
	Value       string
	IsParameter bool
	ParamName   string
	Position    int
}

// analyzePathHierarchy detects parent-child relationships from path structure
func (rd *RelationshipDetector) analyzePathHierarchy(segments []PathSegment, resourceMap map[string]*models.Resource) {
	for i := 0; i < len(segments)-1; i++ {
		current := segments[i]
		next := segments[i+1]

		// Skip if current segment is a parameter
		if current.IsParameter {
			continue
		}

		// Look for pattern: /resource/{id}/childResource
		if i+2 < len(segments) && next.IsParameter {
			childSegment := segments[i+2]
			if !childSegment.IsParameter {
				rd.addRelationship(resourceMap, current.Value, childSegment.Value, "has_many", "path hierarchy", "strong")
				rd.addRelationship(resourceMap, childSegment.Value, current.Value, "belongs_to", "path hierarchy", "strong")
			}
		}
	}
}

// detectParameterRelationships identifies relationships from path and query parameters
func (rd *RelationshipDetector) detectParameterRelationships(resourceMap map[string]*models.Resource) {
	for _, resource := range resourceMap {
		for _, operation := range resource.Operations {
			rd.analyzeOperationParameters(resource, operation, resourceMap)
		}
	}
}

// analyzeOperationParameters examines parameters for foreign key relationships
func (rd *RelationshipDetector) analyzeOperationParameters(resource *models.Resource, operation models.Operation, resourceMap map[string]*models.Resource) {
	// Extract path parameters
	pathParams := rd.pathParamRegex.FindAllStringSubmatch(operation.Path, -1)
	for _, match := range pathParams {
		if len(match) > 1 {
			paramName := match[1]
			rd.analyzeParameterForRelationship(resource.Name, paramName, resourceMap)
		}
	}
}

// analyzeParameterForRelationship checks if a parameter indicates a relationship
func (rd *RelationshipDetector) analyzeParameterForRelationship(resourceName, paramName string, resourceMap map[string]*models.Resource) {
	// Check for foreign key patterns like userId, user_id, projectId, etc.
	matches := rd.foreignKeyRegex.FindStringSubmatch(paramName)
	if len(matches) > 1 {
		// Extract the referenced resource name
		referencedResource := matches[1]
		referencedResource = strings.ToLower(referencedResource)

		// Check if this resource exists
		if targetResource, exists := resourceMap[referencedResource]; exists {
			rd.addRelationship(resourceMap, targetResource.Name, resourceName, "has_many", paramName, "medium")
			rd.addRelationship(resourceMap, resourceName, targetResource.Name, "references", paramName, "medium")
		}
	}
}

// detectSchemaRelationships identifies relationships from request/response schemas
func (rd *RelationshipDetector) detectSchemaRelationships(resourceMap map[string]*models.Resource, spec *parser.OpenAPISpec) {
	// Analyze components/schemas for cross-references
	if spec.Components != nil && spec.Components.Schemas != nil {
		for schemaName, schema := range spec.Components.Schemas {
			rd.analyzeSchemaReferences(schemaName, schema, resourceMap, spec)
		}
	}
}

// analyzeSchemaReferences examines schema properties for references to other resources
func (rd *RelationshipDetector) analyzeSchemaReferences(schemaName string, schema parser.Schema, resourceMap map[string]*models.Resource, spec *parser.OpenAPISpec) {
	if schema.Properties == nil {
		return
	}

	resourceName := rd.schemaToResourceName(schemaName)
	sourceResource := resourceMap[resourceName]
	if sourceResource == nil {
		return
	}

	for propertyName, property := range schema.Properties {
		// Check for $ref to other schemas
		if property.Ref != "" {
			referencedSchema := rd.extractSchemaNameFromRef(property.Ref)
			referencedResource := rd.schemaToResourceName(referencedSchema)

			if targetResource := resourceMap[referencedResource]; targetResource != nil {
				rd.addRelationship(resourceMap, sourceResource.Name, targetResource.Name, "references", propertyName, "strong")
			}
		}

		// Check for array references
		if property.Type == "array" && property.Items != nil && property.Items.Ref != "" {
			referencedSchema := rd.extractSchemaNameFromRef(property.Items.Ref)
			referencedResource := rd.schemaToResourceName(referencedSchema)

			if targetResource := resourceMap[referencedResource]; targetResource != nil {
				rd.addRelationship(resourceMap, sourceResource.Name, targetResource.Name, "has_many", propertyName, "strong")
				rd.addRelationship(resourceMap, targetResource.Name, sourceResource.Name, "belongs_to", propertyName, "strong")
			}
		}

		// Check for foreign key naming patterns in properties
		rd.analyzePropertyForForeignKey(sourceResource.Name, propertyName, resourceMap)
	}
}

// analyzePropertyForForeignKey checks property names for foreign key patterns
func (rd *RelationshipDetector) analyzePropertyForForeignKey(resourceName, propertyName string, resourceMap map[string]*models.Resource) {
	matches := rd.foreignKeyRegex.FindStringSubmatch(propertyName)
	if len(matches) > 1 {
		referencedResource := matches[1]
		referencedResource = strings.ToLower(referencedResource)

		if targetResource := resourceMap[referencedResource]; targetResource != nil {
			rd.addRelationship(resourceMap, resourceName, targetResource.Name, "references", propertyName, "medium")
			rd.addRelationship(resourceMap, targetResource.Name, resourceName, "referenced_by", propertyName, "weak")
		}
	}
}

// schemaToResourceName converts a schema name to a likely resource name
func (rd *RelationshipDetector) schemaToResourceName(schemaName string) string {
	// Remove common suffixes
	name := strings.ToLower(schemaName)
	suffixes := []string{"response", "request", "schema", "model", "dto"}

	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) {
			name = strings.TrimSuffix(name, suffix)
			break
		}
	}

	return name
}

// extractSchemaNameFromRef extracts schema name from a $ref path
func (rd *RelationshipDetector) extractSchemaNameFromRef(ref string) string {
	// Handle #/components/schemas/SchemaName format
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}

// addRelationship safely adds a relationship to a resource
func (rd *RelationshipDetector) addRelationship(resourceMap map[string]*models.Resource, fromResource, toResource, relType, via, strength string) {
	source := resourceMap[fromResource]
	if source == nil {
		return
	}

	// Check if relationship already exists
	for _, existing := range source.Relationships {
		if existing.Resource == toResource && existing.Type == relType && existing.Via == via {
			return // Relationship already exists
		}
	}

	// Add new relationship
	relationship := models.Relationship{
		Resource:    toResource,
		Type:        relType,
		Via:         via,
		Description: rd.generateRelationshipDescription(fromResource, toResource, relType),
		Strength:    strength,
	}

	source.Relationships = append(source.Relationships, relationship)
}

// generateRelationshipDescription creates a human-readable description
func (rd *RelationshipDetector) generateRelationshipDescription(from, to, relType string) string {
	switch relType {
	case "has_many":
		return from + " contains multiple " + to + " resources"
	case "belongs_to":
		return from + " belongs to a " + to + " resource"
	case "references":
		return from + " references a " + to + " resource"
	case "referenced_by":
		return from + " is referenced by " + to + " resources"
	default:
		return from + " has a " + relType + " relationship with " + to
	}
}
