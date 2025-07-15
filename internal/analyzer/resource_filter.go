package analyzer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/orchard9/api-godoc/pkg/models"
)

// ResourceFilter defines filtering criteria for resources
type ResourceFilter struct {
	Include []string // Explicit list of resources to include
	Exclude []string // List of resources to exclude
	Pattern string   // Regex pattern to match resource names
}

// ResourceFilterer filters resources based on criteria
type ResourceFilterer interface {
	FilterResources(resources []models.Resource, filter *ResourceFilter) []models.Resource
}

type resourceFilterer struct{}

// NewResourceFilterer creates a new resource filterer
func NewResourceFilterer() ResourceFilterer {
	return &resourceFilterer{}
}

// FilterResources filters resources based on the provided criteria
func (rf *resourceFilterer) FilterResources(resources []models.Resource, filter *ResourceFilter) []models.Resource {
	// No filter means include all
	if filter == nil || (len(filter.Include) == 0 && len(filter.Exclude) == 0 && filter.Pattern == "") {
		return resources
	}

	// Build include map for O(1) lookups (case insensitive)
	includeMap := make(map[string]bool)
	for _, name := range filter.Include {
		includeMap[strings.ToLower(name)] = true
	}

	// Build exclude map (case insensitive)
	excludeMap := make(map[string]bool)
	for _, name := range filter.Exclude {
		excludeMap[strings.ToLower(name)] = true
	}

	// Compile pattern if provided
	var patternRegex *regexp.Regexp
	if filter.Pattern != "" {
		patternRegex, _ = regexp.Compile(filter.Pattern)
	}

	// Filter resources
	var filtered []models.Resource
	for _, resource := range resources {
		resourceNameLower := strings.ToLower(resource.Name)

		// Check if explicitly included
		if len(includeMap) > 0 {
			if includeMap[resourceNameLower] {
				filtered = append(filtered, resource)
				continue
			}
		}

		// Check if matches pattern
		if patternRegex != nil && patternRegex.MatchString(resource.Name) {
			// Check if not excluded
			if !excludeMap[resourceNameLower] {
				filtered = append(filtered, resource)
			}
			continue
		}

		// If no include list or pattern, include unless excluded
		if len(includeMap) == 0 && patternRegex == nil {
			if !excludeMap[resourceNameLower] {
				filtered = append(filtered, resource)
			}
		}
	}

	return filtered
}

// Validate checks if the filter configuration is valid
func (f *ResourceFilter) Validate() error {
	// Validate regex pattern if provided
	if f.Pattern != "" {
		_, err := regexp.Compile(f.Pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	}
	return nil
}

// String returns a human-readable description of the filter
func (f *ResourceFilter) String() string {
	var parts []string

	if len(f.Include) > 0 {
		parts = append(parts, fmt.Sprintf("include: %s", strings.Join(f.Include, ", ")))
	}

	if len(f.Exclude) > 0 {
		parts = append(parts, fmt.Sprintf("exclude: %s", strings.Join(f.Exclude, ", ")))
	}

	if f.Pattern != "" {
		parts = append(parts, fmt.Sprintf("pattern: %s", f.Pattern))
	}

	if len(parts) == 0 {
		return "no filter"
	}

	return strings.Join(parts, "; ")
}
