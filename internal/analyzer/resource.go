package analyzer

import (
	"regexp"
	"strings"

	"github.com/orchard9/api-godoc/internal/parser"
	"github.com/orchard9/api-godoc/pkg/models"
)

// ResourceAnalyzer handles resource extraction from OpenAPI paths
type ResourceAnalyzer struct {
	pathVariableRegex *regexp.Regexp
}

// NewResourceAnalyzer creates a new resource analyzer
func NewResourceAnalyzer() *ResourceAnalyzer {
	return &ResourceAnalyzer{
		pathVariableRegex: regexp.MustCompile(`\{[^}]+\}`),
	}
}

// ExtractResources analyzes OpenAPI paths to identify business resources
func (ra *ResourceAnalyzer) ExtractResources(spec *parser.OpenAPISpec) []models.Resource {
	resourceMap := make(map[string]*models.Resource)

	// Process each path in the spec
	for path, pathItem := range spec.Paths {
		resourceNames := ra.extractResourceNames(path)

		// Add operations to resources
		ra.addOperations(resourceMap, resourceNames, path, pathItem)
	}

	// Convert map to slice
	resources := make([]models.Resource, 0, len(resourceMap))
	for _, resource := range resourceMap {
		resources = append(resources, *resource)
	}

	return resources
}

// extractResourceNames identifies resource names from a path
func (ra *ResourceAnalyzer) extractResourceNames(path string) []string {
	// Remove path variables (e.g., {id}, {user_id})
	cleanPath := ra.pathVariableRegex.ReplaceAllString(path, "")

	// Split by slashes and filter out empty strings
	segments := strings.Split(cleanPath, "/")
	var resourceNames []string

	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment != "" && ra.isResourceName(segment) {
			resourceNames = append(resourceNames, segment)
		}
	}

	return resourceNames
}

// isResourceName determines if a path segment represents a resource
func (ra *ResourceAnalyzer) isResourceName(segment string) bool {
	// Skip common non-resource segments
	skipSegments := map[string]bool{
		"api":     true,
		"v1":      true,
		"v2":      true,
		"v3":      true,
		"version": true,
		"health":  true,
		"status":  true,
		"ping":    true,
	}

	return !skipSegments[segment]
}

// addOperations adds operations from a PathItem to the appropriate resources
func (ra *ResourceAnalyzer) addOperations(resourceMap map[string]*models.Resource, resourceNames []string, path string, pathItem parser.PathItem) {
	// For each resource name in the path
	for _, resourceName := range resourceNames {
		// Get or create resource
		if resourceMap[resourceName] == nil {
			resourceMap[resourceName] = &models.Resource{
				Name:        resourceName,
				Description: ra.generateResourceDescription(resourceName),
				Operations:  []models.Operation{},
			}
		}

		// Add all operations from the path item
		if pathItem.Get != nil {
			resourceMap[resourceName].Operations = append(resourceMap[resourceName].Operations, ra.createOperation("GET", path, pathItem.Get))
		}
		if pathItem.Post != nil {
			resourceMap[resourceName].Operations = append(resourceMap[resourceName].Operations, ra.createOperation("POST", path, pathItem.Post))
		}
		if pathItem.Put != nil {
			resourceMap[resourceName].Operations = append(resourceMap[resourceName].Operations, ra.createOperation("PUT", path, pathItem.Put))
		}
		if pathItem.Delete != nil {
			resourceMap[resourceName].Operations = append(resourceMap[resourceName].Operations, ra.createOperation("DELETE", path, pathItem.Delete))
		}
		if pathItem.Patch != nil {
			resourceMap[resourceName].Operations = append(resourceMap[resourceName].Operations, ra.createOperation("PATCH", path, pathItem.Patch))
		}
		if pathItem.Head != nil {
			resourceMap[resourceName].Operations = append(resourceMap[resourceName].Operations, ra.createOperation("HEAD", path, pathItem.Head))
		}
		if pathItem.Options != nil {
			resourceMap[resourceName].Operations = append(resourceMap[resourceName].Operations, ra.createOperation("OPTIONS", path, pathItem.Options))
		}
	}
}

// createOperation converts a parser.Operation to models.Operation
func (ra *ResourceAnalyzer) createOperation(method, path string, op *parser.Operation) models.Operation {
	return models.Operation{
		Method:      method,
		Path:        path,
		Summary:     op.Summary,
		Description: op.Description,
		OperationID: op.OperationID,
	}
}

// generateResourceDescription creates a basic description for a resource
func (ra *ResourceAnalyzer) generateResourceDescription(resourceName string) string {
	// Simple capitalization and description generation
	if len(resourceName) == 0 {
		return "Resource operations"
	}

	capitalized := strings.ToUpper(resourceName[:1]) + resourceName[1:]
	return capitalized + " resource operations"
}
