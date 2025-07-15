// Package analyzer provides resource analysis functionality for OpenAPI specs
package analyzer

import (
	"fmt"
	"time"

	"github.com/orchard9/api-godoc/internal/parser"
	"github.com/orchard9/api-godoc/pkg/models"
)

// Analyzer defines the interface for resource analysis
type Analyzer interface {
	// Analyze processes an OpenAPI spec and extracts resources
	Analyze(spec *parser.OpenAPISpec) (*models.APIAnalysis, error)
}

// New creates a new analyzer instance
func New() Analyzer {
	return &analyzer{
		resourceAnalyzer:     NewResourceAnalyzer(),
		relationshipDetector: NewRelationshipDetector(),
	}
}

type analyzer struct {
	resourceAnalyzer     *ResourceAnalyzer
	relationshipDetector *RelationshipDetector
}

func (a *analyzer) Analyze(spec *parser.OpenAPISpec) (*models.APIAnalysis, error) {
	if spec == nil {
		return nil, fmt.Errorf("spec cannot be nil")
	}

	// Extract resources from OpenAPI paths
	resources := a.resourceAnalyzer.ExtractResources(spec)

	// Detect relationships between resources
	a.relationshipDetector.DetectRelationships(resources, spec)

	// Calculate summary statistics
	summary := a.calculateSummary(resources, spec)

	// Determine spec type
	specType := "OpenAPI 3.x"
	if spec.OpenAPI == "" || spec.OpenAPI[:1] == "2" {
		specType = "Swagger 2.0 (converted)"
	}

	// Build base URL from servers
	baseURL := ""
	if len(spec.Servers) > 0 {
		baseURL = spec.Servers[0].URL
	}

	// Create analysis result
	analysis := &models.APIAnalysis{
		Title:         spec.Info.Title,
		Version:       spec.Info.Version,
		Description:   spec.Info.Description,
		BaseURL:       baseURL,
		Resources:     resources,
		Summary:       summary,
		GeneratedAt:   time.Now(),
		SpecType:      specType,
		OriginalPaths: len(spec.Paths),
	}

	return analysis, nil
}

// calculateSummary generates summary statistics for the API analysis
func (a *analyzer) calculateSummary(resources []models.Resource, spec *parser.OpenAPISpec) models.AnalysisStat {
	totalOperations := 0
	resourcesWithOps := 0

	for _, resource := range resources {
		totalOperations += len(resource.Operations)
		if len(resource.Operations) > 0 {
			resourcesWithOps++
		}
	}

	// Calculate resource coverage (percentage of paths that map to resources)
	totalPaths := len(spec.Paths)
	coverage := 0
	if totalPaths > 0 {
		coverage = (resourcesWithOps * 100) / totalPaths
		if coverage > 100 {
			coverage = 100 // Cap at 100% since multiple operations can map to same resource
		}
	}

	return models.AnalysisStat{
		TotalResources:   len(resources),
		TotalOperations:  totalOperations,
		TotalEndpoints:   totalPaths,
		ResourceCoverage: coverage,
	}
}
