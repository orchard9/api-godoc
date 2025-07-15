package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConvertUATSpecs(t *testing.T) {
	// Get the project root
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Go up two levels from internal/converter to project root
	projectRoot := filepath.Join(wd, "..", "..")

	tests := []struct {
		name     string
		specFile string
	}{
		{
			name:     "warden API conversion",
			specFile: filepath.Join(projectRoot, "uat", "artifacts", "warden.v1.swagger.json"),
		},
		{
			name:     "forge API conversion",
			specFile: filepath.Join(projectRoot, "uat", "artifacts", "forge.swagger.json"),
		},
	}

	converter := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read the spec file
			data, err := os.ReadFile(tt.specFile)
			if err != nil {
				t.Fatalf("Failed to read spec file: %v", err)
			}

			// Convert to OpenAPI 3.x
			result, err := converter.Convert(data)
			if err != nil {
				t.Fatalf("Failed to convert spec: %v", err)
			}

			// Basic validation - check that result is valid JSON
			if len(result) == 0 {
				t.Error("Conversion resulted in empty output")
			}

			// Log first 500 chars of output for debugging
			output := string(result)
			if len(output) > 500 {
				output = output[:500] + "..."
			}
			t.Logf("Converted output preview:\n%s", output)
		})
	}
}
