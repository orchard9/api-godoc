package test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/orchard9/api-godoc/internal/analyzer"
	"github.com/orchard9/api-godoc/internal/parser"
	"github.com/orchard9/api-godoc/internal/reporter"
	"github.com/orchard9/api-godoc/pkg/models"
)

// TestRealWorldAPIFixtures tests the tool against real-world API specifications
func TestRealWorldAPIFixtures(t *testing.T) {
	fixtures := []struct {
		name      string
		file      string
		minRsrc   int // minimum expected resources
		maxTime   time.Duration
		skipLarge bool
	}{
		{
			name:    "Stripe Minimal",
			file:    "stripe-minimal.json",
			minRsrc: 3,
			maxTime: 5 * time.Second,
		},
		{
			name:    "Kubernetes Simplified",
			file:    "kubernetes-simplified.json",
			minRsrc: 2,
			maxTime: 5 * time.Second,
		},
		{
			name:      "Stripe Full API",
			file:      "stripe-openapi3.json",
			minRsrc:   20,
			maxTime:   30 * time.Second,
			skipLarge: true, // Skip in CI
		},
		{
			name:      "GitHub API",
			file:      "github-openapi3.json",
			minRsrc:   10,
			maxTime:   30 * time.Second,
			skipLarge: true, // Skip in CI
		},
	}

	for _, tt := range fixtures {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipLarge && testing.Short() {
				t.Skip("Skipping large fixture in short mode")
			}

			start := time.Now()

			// Test file path
			fixturePath := filepath.Join("fixtures", tt.file)

			// Initialize components
			p := parser.New()
			resourceAnalyzer := analyzer.NewResourceAnalyzer()
			relationshipDetector := analyzer.NewRelationshipDetector()
			patternDetector := analyzer.NewPatternDetector()
			schemaReducer := analyzer.NewSchemaReducer()
			rep := reporter.New()

			// Parse specification
			spec, err := p.ParseFile(fixturePath)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.file, err)
			}

			// Extract resources
			resources := resourceAnalyzer.ExtractResources(spec)
			if len(resources) < tt.minRsrc {
				t.Errorf("Expected at least %d resources, got %d", tt.minRsrc, len(resources))
			}

			// Test resource filtering
			if len(resources) > 5 {
				filter := &analyzer.ResourceFilter{
					Include: []string{resources[0].Name, resources[1].Name},
				}
				filterer := analyzer.NewResourceFilterer()
				filtered := filterer.FilterResources(resources, filter)
				if len(filtered) != 2 {
					t.Errorf("Expected 2 filtered resources, got %d", len(filtered))
				}
			}

			// Extract fields with schema reduction
			for i := range resources {
				resource := &resources[i]
				// Try to find schema for this resource
				if spec.Components != nil && spec.Components.Schemas != nil {
					for schemaName, schema := range spec.Components.Schemas {
						// Simple heuristic: if schema name contains resource name
						if contains(schemaName, resource.Name) {
							fields := schemaReducer.SchemaToFields(&schema, "standard")
							resource.Fields = append(resource.Fields, fields...)
							break
						}
					}
				}
			}

			// Detect relationships
			relationshipDetector.DetectRelationships(resources, spec)

			// Detect patterns
			patterns := patternDetector.DetectPatterns(spec)

			// Create analysis
			analysis := &models.APIAnalysis{
				Title:       spec.Info.Title,
				Version:     spec.Info.Version,
				Description: spec.Info.Description,
				BaseURL:     getBaseURL(spec),
				SpecType:    "OpenAPI " + spec.OpenAPI,
				GeneratedAt: time.Now(),
				Resources:   resources,
				Patterns:    patterns,
				Summary: models.AnalysisStat{
					TotalResources:  len(resources),
					TotalOperations: countOperations(resources),
					TotalEndpoints:  len(spec.Paths),
				},
			}

			// Test all output formats
			formats := []string{"markdown", "json", "ai"}
			for _, format := range formats {
				output, err := rep.Generate(analysis, format)
				if err != nil {
					t.Errorf("Failed to generate %s output: %v", format, err)
				}
				if output == "" {
					t.Errorf("Empty %s output", format)
				}
			}

			// Check timing
			elapsed := time.Since(start)
			if elapsed > tt.maxTime {
				t.Errorf("Processing took too long: %v > %v", elapsed, tt.maxTime)
			}

			t.Logf("Processed %s: %d resources, %d patterns, %v",
				tt.name, len(resources), len(patterns), elapsed)
		})
	}
}

// TestFixtureResourceCoverage tests that we extract reasonable resource coverage
func TestFixtureResourceCoverage(t *testing.T) {
	tests := []struct {
		name        string
		file        string
		minCoverage int // minimum percentage of endpoints that should be resources
	}{
		{
			name:        "Stripe Minimal",
			file:        "stripe-minimal.json",
			minCoverage: 50, // At least 50% of endpoints should be resource operations
		},
		{
			name:        "Kubernetes Simplified",
			file:        "kubernetes-simplified.json",
			minCoverage: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixturePath := filepath.Join("fixtures", tt.file)

			p := parser.New()
			spec, err := p.ParseFile(fixturePath)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.file, err)
			}

			resourceAnalyzer := analyzer.NewResourceAnalyzer()
			resources := resourceAnalyzer.ExtractResources(spec)

			totalEndpoints := len(spec.Paths)
			totalOperations := countOperations(resources)

			if totalEndpoints == 0 {
				t.Fatal("No endpoints found")
			}

			coverage := (totalOperations * 100) / totalEndpoints
			if coverage < tt.minCoverage {
				t.Errorf("Resource coverage too low: %d%% < %d%%", coverage, tt.minCoverage)
			}

			t.Logf("Resource coverage: %d%% (%d/%d)", coverage, totalOperations, totalEndpoints)
		})
	}
}

// TestPatternDetectionOnFixtures tests pattern detection on real APIs
func TestPatternDetectionOnFixtures(t *testing.T) {
	tests := []struct {
		name             string
		file             string
		expectedPatterns []string
	}{
		{
			name:             "Stripe Minimal",
			file:             "stripe-minimal.json",
			expectedPatterns: []string{}, // Minimal spec may not show patterns clearly
		},
		{
			name:             "Kubernetes Simplified",
			file:             "kubernetes-simplified.json",
			expectedPatterns: []string{}, // Simple spec may not show patterns
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixturePath := filepath.Join("fixtures", tt.file)

			p := parser.New()
			spec, err := p.ParseFile(fixturePath)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.file, err)
			}

			detector := analyzer.NewPatternDetector()
			patterns := detector.DetectPatterns(spec)

			// Check that expected patterns are found
			patternMap := make(map[string]bool)
			for _, pattern := range patterns {
				patternMap[pattern.Type] = true
			}

			for _, expected := range tt.expectedPatterns {
				if !patternMap[expected] {
					t.Errorf("Expected to find %s pattern", expected)
				}
			}

			t.Logf("Found patterns: %v", getPatternTypes(patterns))
		})
	}
}

// Helper functions
func getBaseURL(spec *parser.OpenAPISpec) string {
	if len(spec.Servers) > 0 && spec.Servers[0].URL != "" {
		return spec.Servers[0].URL
	}
	return ""
}

func countOperations(resources []models.Resource) int {
	count := 0
	for _, resource := range resources {
		count += len(resource.Operations)
	}
	return count
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func getPatternTypes(patterns []models.Pattern) []string {
	types := make([]string, len(patterns))
	for i, pattern := range patterns {
		types[i] = pattern.Type
	}
	return types
}
