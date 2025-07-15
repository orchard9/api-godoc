package internal

import (
	"testing"

	"github.com/orchard9/api-godoc/internal/analyzer"
	"github.com/orchard9/api-godoc/internal/reporter"
	"github.com/orchard9/api-godoc/pkg/models"
)

// TestAnalyzerBasics tests basic analyzer functionality
func TestAnalyzerBasics(t *testing.T) {
	t.Run("ResourceFilter validation", func(t *testing.T) {
		filter := &analyzer.ResourceFilter{
			Pattern: "[invalid",
		}
		err := filter.Validate()
		if err == nil {
			t.Error("Expected error for invalid regex pattern")
		}
	})

	t.Run("ResourceFilter with valid pattern", func(t *testing.T) {
		filter := &analyzer.ResourceFilter{
			Pattern: "^user.*",
			Include: []string{"admin"},
			Exclude: []string{"test"},
		}
		err := filter.Validate()
		if err != nil {
			t.Errorf("Unexpected error for valid filter: %v", err)
		}
	})
}

// TestReporterBasics tests basic reporter functionality
func TestReporterBasics(t *testing.T) {
	rep := reporter.New()

	t.Run("Generate markdown", func(t *testing.T) {
		analysis := &models.APIAnalysis{
			Title:       "Test API",
			Version:     "1.0.0",
			Description: "Test description",
			Resources: []models.Resource{
				{
					Name: "users",
					Operations: []models.Operation{
						{Method: "GET", Path: "/users"},
					},
				},
			},
		}
		output, err := rep.Generate(analysis, "markdown")
		if err != nil {
			t.Errorf("Failed to generate markdown: %v", err)
		}
		if output == "" {
			t.Error("Expected non-empty markdown output")
		}
	})

	t.Run("Generate JSON", func(t *testing.T) {
		analysis := &models.APIAnalysis{
			Title: "Test API",
		}
		output, err := rep.Generate(analysis, "json")
		if err != nil {
			t.Errorf("Failed to generate JSON: %v", err)
		}
		if output == "" {
			t.Error("Expected non-empty JSON output")
		}
	})

	t.Run("Generate AI optimized", func(t *testing.T) {
		analysis := &models.APIAnalysis{
			Title: "Test API",
			Resources: []models.Resource{
				{Name: "users"},
			},
		}
		output, err := rep.Generate(analysis, "ai")
		if err != nil {
			t.Errorf("Failed to generate AI output: %v", err)
		}
		if output == "" {
			t.Error("Expected non-empty AI output")
		}
	})
}
