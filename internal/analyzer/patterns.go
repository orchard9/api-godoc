package analyzer

import (
	"fmt"
	"strings"

	"github.com/orchard9/pg-goapi/internal/parser"
	"github.com/orchard9/pg-goapi/pkg/models"
)

// PatternDetector analyzes OpenAPI specs to detect common API patterns
type PatternDetector interface {
	DetectPatterns(spec *parser.OpenAPISpec) []models.Pattern
}

// patternDetector implements the PatternDetector interface
type patternDetector struct{}

// NewPatternDetector creates a new pattern detector
func NewPatternDetector() PatternDetector {
	return &patternDetector{}
}

// DetectPatterns analyzes the OpenAPI spec to find common patterns
func (pd *patternDetector) DetectPatterns(spec *parser.OpenAPISpec) []models.Pattern {
	var patterns []models.Pattern

	// Detect pagination
	if pattern := pd.detectPagination(spec); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// Detect filtering
	if pattern := pd.detectFiltering(spec); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// Detect sorting
	if pattern := pd.detectSorting(spec); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// Detect versioning
	if pattern := pd.detectVersioning(spec); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// Detect batch operations
	if pattern := pd.detectBatchOperations(spec); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// Detect search
	if pattern := pd.detectSearch(spec); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	// Detect authentication
	if pattern := pd.detectAuthentication(spec); pattern != nil {
		patterns = append(patterns, *pattern)
	}

	return patterns
}

// detectPagination looks for pagination patterns in the API
func (pd *patternDetector) detectPagination(spec *parser.OpenAPISpec) *models.Pattern {
	var examples []string
	paginationTypes := make(map[string]bool)

	for path, pathItem := range spec.Paths {
		ops := pd.getOperations(pathItem)
		for _, op := range ops {
			if op == nil {
				continue
			}

			if pd.hasPaginationParams(op.Parameters) {
				examples = append(examples, path)
				// Detect pagination type
				for _, param := range op.Parameters {
					switch param.Name {
					case "page", "pageNumber":
						paginationTypes["page-based"] = true
					case "offset":
						paginationTypes["offset-based"] = true
					case "cursor", "next_token", "continuation_token":
						paginationTypes["cursor-based"] = true
					}
				}
			}
		}
	}

	if len(examples) == 0 {
		return nil
	}

	var typeList []string
	for t := range paginationTypes {
		typeList = append(typeList, t)
	}

	description := fmt.Sprintf("API uses pagination for list endpoints. Types detected: %s",
		strings.Join(typeList, ", "))

	return &models.Pattern{
		Type:        "pagination",
		Description: description,
		Examples:    examples,
		Confidence:  pd.calculateConfidence(len(examples)),
		Impact:      "Clients should implement pagination handling for list operations",
	}
}

// detectFiltering looks for filtering patterns
func (pd *patternDetector) detectFiltering(spec *parser.OpenAPISpec) *models.Pattern {
	var examples []string
	filterParams := make(map[string]int)

	for path, pathItem := range spec.Paths {
		ops := pd.getOperations(pathItem)
		for _, op := range ops {
			if op == nil || op.Parameters == nil {
				continue
			}

			hasFilter := false
			for _, param := range op.Parameters {
				if param.In == "query" && pd.isFilterParam(param.Name) {
					hasFilter = true
					filterParams[param.Name]++
				}
			}

			if hasFilter {
				examples = append(examples, path)
			}
		}
	}

	if len(examples) == 0 {
		return nil
	}

	// Find most common filter parameters
	var commonFilters []string
	for param, count := range filterParams {
		if count >= 2 {
			commonFilters = append(commonFilters, param)
		}
	}

	description := "API supports filtering on list endpoints"
	if len(commonFilters) > 0 {
		description += fmt.Sprintf(". Common filters: %s", strings.Join(commonFilters, ", "))
	}

	return &models.Pattern{
		Type:        "filtering",
		Description: description,
		Examples:    examples,
		Confidence:  pd.calculateConfidence(len(examples)),
		Impact:      "Clients can filter results using query parameters",
	}
}

// detectSorting looks for sorting patterns
func (pd *patternDetector) detectSorting(spec *parser.OpenAPISpec) *models.Pattern {
	var examples []string
	sortParams := make(map[string]bool)

	for path, pathItem := range spec.Paths {
		ops := pd.getOperations(pathItem)
		for _, op := range ops {
			if op == nil {
				continue
			}

			for _, param := range op.Parameters {
				if param.In == "query" && pd.isSortParam(param.Name) {
					examples = append(examples, path)
					sortParams[param.Name] = true
					break
				}
			}
		}
	}

	if len(examples) == 0 {
		return nil
	}

	var paramList []string
	for p := range sortParams {
		paramList = append(paramList, p)
	}

	description := fmt.Sprintf("API supports sorting on list endpoints using parameters: %s",
		strings.Join(paramList, ", "))

	return &models.Pattern{
		Type:        "sorting",
		Description: description,
		Examples:    examples,
		Confidence:  pd.calculateConfidence(len(examples)),
		Impact:      "Clients can sort results using query parameters",
	}
}

// detectVersioning looks for API versioning patterns
func (pd *patternDetector) detectVersioning(spec *parser.OpenAPISpec) *models.Pattern {
	versionedPaths := make(map[string][]string)

	for path := range spec.Paths {
		if version := pd.extractVersion(path); version != "" {
			versionedPaths[version] = append(versionedPaths[version], path)
		}
	}

	if len(versionedPaths) == 0 {
		return nil
	}

	var versions []string
	var examples []string
	for v, paths := range versionedPaths {
		versions = append(versions, v)
		if len(examples) < 3 {
			examples = append(examples, paths[0])
		}
	}

	description := fmt.Sprintf("API uses URL path versioning. Versions found: %s",
		strings.Join(versions, ", "))

	return &models.Pattern{
		Type:        "versioning",
		Description: description,
		Examples:    examples,
		Confidence:  "high", // URL versioning is explicit
		Impact:      "Clients should be aware of API version compatibility",
	}
}

// detectBatchOperations looks for batch operation patterns
func (pd *patternDetector) detectBatchOperations(spec *parser.OpenAPISpec) *models.Pattern {
	var examples []string

	for path, pathItem := range spec.Paths {
		// Check for common batch operation patterns
		if strings.Contains(path, ":batch") ||
			strings.Contains(path, "/batch") ||
			strings.Contains(path, "bulk") {
			examples = append(examples, path)
			continue
		}

		// Check operation summaries for batch indicators
		ops := pd.getOperations(pathItem)
		for _, op := range ops {
			if op != nil && op.Summary != "" {
				summary := strings.ToLower(op.Summary)
				if strings.Contains(summary, "batch") ||
					strings.Contains(summary, "bulk") ||
					strings.Contains(summary, "multiple") {
					examples = append(examples, path)
					break
				}
			}
		}
	}

	if len(examples) == 0 {
		return nil
	}

	return &models.Pattern{
		Type:        "batch_operations",
		Description: "API supports batch operations for bulk create/update/delete",
		Examples:    examples,
		Confidence:  pd.calculateConfidence(len(examples)),
		Impact:      "Clients can perform bulk operations for better performance",
	}
}

// detectSearch looks for search functionality
func (pd *patternDetector) detectSearch(spec *parser.OpenAPISpec) *models.Pattern {
	var examples []string
	searchParams := make(map[string]bool)

	for path, pathItem := range spec.Paths {
		// Check if path contains search
		if strings.Contains(path, "/search") {
			examples = append(examples, path)
			continue
		}

		// Check for search parameters
		ops := pd.getOperations(pathItem)
		for _, op := range ops {
			if op == nil {
				continue
			}

			for _, param := range op.Parameters {
				if param.In == "query" && pd.isSearchParam(param.Name) {
					examples = append(examples, path)
					searchParams[param.Name] = true
					break
				}
			}
		}
	}

	if len(examples) == 0 {
		return nil
	}

	var paramList []string
	for p := range searchParams {
		paramList = append(paramList, p)
	}

	description := "API provides search functionality"
	if len(paramList) > 0 {
		description += fmt.Sprintf(" using parameters: %s", strings.Join(paramList, ", "))
	}

	return &models.Pattern{
		Type:        "search",
		Description: description,
		Examples:    examples,
		Confidence:  pd.calculateConfidence(len(examples)),
		Impact:      "Clients can perform full-text or field-based searches",
	}
}

// detectAuthentication looks for authentication patterns
func (pd *patternDetector) detectAuthentication(spec *parser.OpenAPISpec) *models.Pattern {
	if spec.Components == nil || len(spec.Components.SecuritySchemes) == 0 {
		return nil
	}

	var authTypes []string
	for name, scheme := range spec.Components.SecuritySchemes {
		switch scheme.Type {
		case "http":
			if scheme.Scheme == "bearer" {
				authTypes = append(authTypes, fmt.Sprintf("Bearer token (%s)", name))
			} else {
				authTypes = append(authTypes, fmt.Sprintf("HTTP %s", scheme.Scheme))
			}
		case "apiKey":
			authTypes = append(authTypes, fmt.Sprintf("API Key in %s", scheme.In))
		case "oauth2":
			authTypes = append(authTypes, "OAuth 2.0")
		case "openIdConnect":
			authTypes = append(authTypes, "OpenID Connect")
		}
	}

	if len(authTypes) == 0 {
		return nil
	}

	description := fmt.Sprintf("API uses authentication: %s", strings.Join(authTypes, ", "))

	// Find secured endpoints
	var examples []string
	for path := range spec.Paths {
		if len(examples) < 3 {
			examples = append(examples, path)
		}
	}

	return &models.Pattern{
		Type:        "authentication",
		Description: description,
		Examples:    examples,
		Confidence:  "high", // Security schemes are explicit
		Impact:      "Clients must implement proper authentication handling",
	}
}

// Helper methods

// hasPaginationParams checks if parameters include pagination
func (pd *patternDetector) hasPaginationParams(params []parser.Parameter) bool {
	paginationParams := map[string]bool{
		"page": true, "limit": true, "offset": true, "size": true,
		"pageSize": true, "pageNumber": true, "per_page": true,
		"cursor": true, "next_token": true, "continuation_token": true,
	}

	foundParams := 0
	for _, param := range params {
		if param.In == "query" && paginationParams[param.Name] {
			foundParams++
		}
	}

	// Need at least 2 pagination-related params (e.g., page+limit)
	return foundParams >= 2 || (foundParams == 1 && pd.hasLimitParam(params))
}

// hasLimitParam checks for limit/size parameters
func (pd *patternDetector) hasLimitParam(params []parser.Parameter) bool {
	for _, param := range params {
		if param.In == "query" && (param.Name == "limit" || param.Name == "size" || param.Name == "pageSize") {
			return true
		}
	}
	return false
}

// isFilterParam checks if a parameter name suggests filtering
func (pd *patternDetector) isFilterParam(name string) bool {
	// Common filter parameter patterns
	filterPatterns := []string{
		"status", "type", "category", "tag", "role", "state",
		"created", "updated", "modified", "since", "before", "after",
		"min", "max", "from", "to", "between",
	}

	nameLower := strings.ToLower(name)
	for _, pattern := range filterPatterns {
		if strings.Contains(nameLower, pattern) {
			return true
		}
	}

	// Check for field names that end with common filter suffixes
	return strings.HasSuffix(nameLower, "_id") ||
		strings.HasSuffix(nameLower, "_at") ||
		strings.HasSuffix(nameLower, "_date")
}

// isSortParam checks if a parameter name suggests sorting
func (pd *patternDetector) isSortParam(name string) bool {
	sortParams := map[string]bool{
		"sort": true, "order": true, "orderBy": true, "sortBy": true,
		"sortOrder": true, "sortDirection": true, "orderDirection": true,
	}
	return sortParams[name]
}

// isSearchParam checks if a parameter name suggests search
func (pd *patternDetector) isSearchParam(name string) bool {
	searchParams := map[string]bool{
		"q": true, "query": true, "search": true, "keyword": true,
		"term": true, "text": true, "find": true,
	}
	return searchParams[name]
}

// extractVersion extracts version from path
func (pd *patternDetector) extractVersion(path string) string {
	// Look for version patterns like /v1/, /v2/, /api/v3/
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "v") && len(part) > 1 {
			// Check if rest is numeric
			rest := part[1:]
			if pd.isNumeric(rest) {
				return part
			}
		}
	}
	return ""
}

// isNumeric checks if string contains only digits
func (pd *patternDetector) isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// getOperations extracts all operations from a path item
func (pd *patternDetector) getOperations(pathItem parser.PathItem) []*parser.Operation {
	return []*parser.Operation{
		pathItem.Get, pathItem.Post, pathItem.Put,
		pathItem.Delete, pathItem.Patch, pathItem.Head,
		pathItem.Options, pathItem.Trace,
	}
}

// calculateConfidence determines confidence level based on examples
func (pd *patternDetector) calculateConfidence(exampleCount int) string {
	if exampleCount >= 3 {
		return "high"
	} else if exampleCount >= 2 {
		return "medium"
	}
	return "low"
}
