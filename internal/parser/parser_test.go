package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "valid OpenAPI 3.0 JSON",
			content: `{
				"openapi": "3.0.3",
				"info": {
					"title": "Test API",
					"version": "1.0.0"
				},
				"paths": {}
			}`,
			wantErr: false,
		},
		{
			name: "valid OpenAPI 3.0 YAML",
			content: `openapi: 3.0.3
info:
  title: Test API
  version: 1.0.0
paths: {}`,
			wantErr: false,
		},
		{
			name: "valid Swagger 2.0 JSON (should convert)",
			content: `{
				"swagger": "2.0",
				"info": {
					"title": "Test API",
					"version": "1.0.0"
				},
				"paths": {}
			}`,
			wantErr: false,
		},
		{
			name: "valid Swagger 2.0 YAML (should convert)",
			content: `swagger: "2.0"
info:
  title: Test API
  version: 1.0.0
paths: {}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			content: `{invalid json`,
			wantErr: true,
		},
		{
			name: "invalid YAML",
			content: `invalid:
  - yaml
  content
    bad indentation`,
			wantErr: true,
		},
		{
			name: "unsupported version",
			content: `{
				"swagger": "1.0",
				"info": {"title": "Old API", "version": "1.0.0"},
				"paths": {}
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpfile, err := os.CreateTemp("", "test-*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Test parsing
			p := New()
			spec, err := p.ParseFile(tmpfile.Name())

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && spec != nil {
				// Verify it's OpenAPI 3.x format
				if spec.OpenAPI == "" {
					t.Error("Expected OpenAPI version to be set")
				}
				if !strings.HasPrefix(spec.OpenAPI, "3.") {
					t.Errorf("Expected OpenAPI 3.x, got %s", spec.OpenAPI)
				}
			}
		})
	}
}

func TestParseURL(t *testing.T) {
	// Note: Actual HTTP testing would require a mock server
	t.Run("URL parsing placeholder", func(t *testing.T) {
		p := New()
		// Test that method exists
		_, err := p.ParseURL("https://example.com/openapi.json")
		if err == nil {
			t.Skip("URL parsing not yet implemented")
		}
	})
}

func TestParseStdin(t *testing.T) {
	t.Run("stdin parsing", func(t *testing.T) {
		// Create a pipe to simulate stdin
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}

		// Write test data
		testData := `{
			"openapi": "3.0.3",
			"info": {"title": "Test", "version": "1.0.0"},
			"paths": {}
		}`

		go func() {
			defer w.Close()
			_, _ = w.Write([]byte(testData))
		}()

		// Replace stdin temporarily
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = r

		p := New()
		spec, err := p.ParseStdin()
		if err != nil {
			t.Fatalf("ParseStdin() error = %v", err)
		}

		if spec.Info.Title != "Test" {
			t.Errorf("Expected title 'Test', got %s", spec.Info.Title)
		}
	})
}

func TestParseRealWorldSpecs(t *testing.T) {
	// Get project root
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	tests := []struct {
		name     string
		specFile string
	}{
		{
			name:     "warden Swagger 2.0",
			specFile: filepath.Join(projectRoot, "uat", "artifacts", "warden.v1.swagger.json"),
		},
		{
			name:     "forge Swagger 2.0",
			specFile: filepath.Join(projectRoot, "uat", "artifacts", "forge.swagger.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			spec, err := p.ParseFile(tt.specFile)
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.specFile, err)
			}

			// Verify conversion happened
			if spec.OpenAPI == "" {
				t.Error("Expected OpenAPI version after conversion")
			}
			if spec.Info.Title == "" {
				t.Error("Expected title to be preserved")
			}

			t.Logf("Successfully parsed %s: %s v%s", tt.name, spec.Info.Title, spec.Info.Version)
		})
	}
}
