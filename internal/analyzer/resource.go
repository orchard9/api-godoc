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

	// First pass: analyze all paths to identify true resources vs actions
	resourcePatterns := ra.analyzePathPatterns(spec.Paths)

	// Process each path in the spec
	for path, pathItem := range spec.Paths {
		resourceNames := ra.extractResourcesForPath(path, resourcePatterns)

		// Add operations to the appropriate resources
		for _, resourceName := range resourceNames {
			ra.addOperationsToResource(resourceMap, resourceName, path, pathItem)
		}
	}

	// Convert map to slice
	resources := make([]models.Resource, 0, len(resourceMap))
	for _, resource := range resourceMap {
		resources = append(resources, *resource)
	}

	return resources
}

// PathPattern represents analysis of path patterns to identify true resources
type PathPattern struct {
	HasVariables bool     // Whether this resource has paths with variables
	Paths        []string // All paths for this resource
}

// analyzePathPatterns analyzes all paths to identify true resources vs actions
func (ra *ResourceAnalyzer) analyzePathPatterns(paths map[string]parser.PathItem) map[string]*PathPattern {
	patterns := make(map[string]*PathPattern)

	for path := range paths {
		// Extract all potential resource names from the path
		segments := ra.extractResourceNames(path)

		// Check if this path has variables
		hasVariables := ra.pathVariableRegex.MatchString(path)

		// Update patterns for each resource segment
		for _, segment := range segments {
			if patterns[segment] == nil {
				patterns[segment] = &PathPattern{
					HasVariables: false,
					Paths:        []string{},
				}
			}

			patterns[segment].Paths = append(patterns[segment].Paths, path)
			if hasVariables {
				patterns[segment].HasVariables = true
			}
		}
	}

	return patterns
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
			// Handle action syntax like "tasks/{id}:move" -> extract "move" from ":move"
			segment = strings.TrimPrefix(segment, ":")
			resourceNames = append(resourceNames, segment)
		}
	}

	return resourceNames
}

// extractResourcesForPath determines which resources a path should belong to
func (ra *ResourceAnalyzer) extractResourcesForPath(path string, patterns map[string]*PathPattern) []string {
	segments := ra.extractResourceNames(path)

	// If no segments found, return empty
	if len(segments) == 0 {
		return []string{}
	}

	// Strategy for arbitrary depth:
	// 1. Identify if this is a nested resource path (has variables between resource segments)
	// 2. For actions without variables, group under the closest parent resource
	// 3. For true resources with variables, create all resources in the chain
	// 4. For non-nested paths, prioritize the most specific resource

	pathSegments := strings.Split(strings.Trim(path, "/"), "/")
	resourcePositions := ra.findResourcePositions(pathSegments)

	// If we only have one resource segment, it's straightforward
	if len(resourcePositions) == 1 {
		return []string{segments[0]}
	}

	// Check if this is a nested resource path by looking for variables between resource segments
	hasNestedStructure := ra.hasVariablesBetweenResources(pathSegments, resourcePositions)

	if hasNestedStructure {
		// Check if the last segment is an action (has no variables)
		if len(segments) > 0 {
			lastSegment := segments[len(segments)-1]
			if pattern, exists := patterns[lastSegment]; exists && !pattern.HasVariables {
				// This is an action on a nested resource, return all segments except the action
				return segments[:len(segments)-1]
			}
		}

		// This is a true nested resource path - create all resources in the chain
		return segments
	}

	// This is a non-nested path - find the most appropriate resource
	return ra.selectPrimaryResource(segments, patterns, path)
}

// findResourcePositions identifies the positions of resource segments in the path
func (ra *ResourceAnalyzer) findResourcePositions(pathSegments []string) []int {
	var resourcePositions []int
	for i, pathSegment := range pathSegments {
		cleanSegment := ra.pathVariableRegex.ReplaceAllString(pathSegment, "")
		cleanSegment = strings.TrimPrefix(cleanSegment, ":")
		cleanSegment = strings.TrimSpace(cleanSegment)

		if cleanSegment != "" && ra.isResourceName(cleanSegment) {
			resourcePositions = append(resourcePositions, i)
		}
	}
	return resourcePositions
}

// hasVariablesBetweenResources checks if there are variables between any consecutive resource segments
func (ra *ResourceAnalyzer) hasVariablesBetweenResources(pathSegments []string, resourcePositions []int) bool {
	if len(resourcePositions) < 2 {
		return false
	}

	// Check each consecutive pair of resource positions
	for i := 0; i < len(resourcePositions)-1; i++ {
		pos1, pos2 := resourcePositions[i], resourcePositions[i+1]

		// Check if there are variables between these resource positions
		for j := pos1 + 1; j < pos2; j++ {
			if ra.pathVariableRegex.MatchString(pathSegments[j]) {
				return true
			}
		}
	}
	return false
}

// selectPrimaryResource chooses the most appropriate resource for non-nested paths
func (ra *ResourceAnalyzer) selectPrimaryResource(segments []string, patterns map[string]*PathPattern, path string) []string {
	// Check if the rightmost segment is an action (has no variables)
	if len(segments) > 0 {
		lastSegment := segments[len(segments)-1]
		if pattern, exists := patterns[lastSegment]; exists && !pattern.HasVariables {
			// This is an action, return its parent resource
			if len(segments) > 1 {
				return []string{segments[len(segments)-2]}
			}
		}
	}

	// For paths with variables or true resources, find the rightmost segment that has variables
	for i := len(segments) - 1; i >= 0; i-- {
		segment := segments[i]
		if pattern, exists := patterns[segment]; exists && pattern.HasVariables {
			return []string{segment}
		}
	}

	// Default: return the most specific (rightmost) segment
	return []string{segments[len(segments)-1]}
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

// addOperationsToResource adds operations from a PathItem to a single resource
func (ra *ResourceAnalyzer) addOperationsToResource(resourceMap map[string]*models.Resource, resourceName string, path string, pathItem parser.PathItem) {
	// Skip empty resource names
	if resourceName == "" {
		return
	}

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
