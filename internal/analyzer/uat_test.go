package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/orchard9/pg-goapi/internal/parser"
)

func TestUATAnalysis(t *testing.T) {
	// Get project root
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	tests := []struct {
		name         string
		specFile     string
		wantErr      bool
		minResources int
	}{
		{
			name:         "warden API analysis",
			specFile:     filepath.Join(projectRoot, "uat", "artifacts", "warden.v1.swagger.json"),
			wantErr:      false,
			minResources: 1,
		},
		{
			name:         "forge API analysis",
			specFile:     filepath.Join(projectRoot, "uat", "artifacts", "forge.swagger.json"),
			wantErr:      false,
			minResources: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse spec
			p := parser.New()
			spec, err := p.ParseFile(tt.specFile)
			if err != nil {
				t.Fatalf("ParseFile() error = %v", err)
			}

			// Analyze spec
			a := New()
			result, err := a.Analyze(spec)

			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("Expected non-nil result")
					return
				}

				if len(result.Resources) < tt.minResources {
					t.Errorf("Expected at least %d resources, got %d", tt.minResources, len(result.Resources))
				}

				// Verify API info is preserved
				if result.Title == "" {
					t.Error("Expected API title to be preserved")
				}

				// Log analysis results
				t.Logf("Analysis of %s:", tt.name)
				t.Logf("  Title: %s", result.Title)
				t.Logf("  Version: %s", result.Version)
				t.Logf("  Resources: %d", len(result.Resources))

				for i, resource := range result.Resources {
					t.Logf("  Resource %d: %s (%d operations)", i+1, resource.Name, len(resource.Operations))
					for j, op := range resource.Operations {
						if j < 3 { // Show first 3 operations
							t.Logf("    %s %s: %s", op.Method, op.Path, op.Summary)
						}
					}
					if len(resource.Operations) > 3 {
						t.Logf("    ... and %d more operations", len(resource.Operations)-3)
					}
				}
			}
		})
	}
}
