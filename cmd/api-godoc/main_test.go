package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	// Note: flag parsing tests are complex due to global state
	// This is a simplified test for the core logic
	t.Run("parseCommaSeparated", func(t *testing.T) {
		tests := []struct {
			input string
			want  []string
		}{
			{"users", []string{"users"}},
			{"users,orders", []string{"users", "orders"}},
			{"users, orders", []string{"users", "orders"}},
			{"", []string{}},
		}

		for _, tt := range tests {
			got := parseCommaSeparated(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("parseCommaSeparated(%q) returned %d items, want %d", tt.input, len(got), len(tt.want))
			}
		}
	})
}

func TestBuildResourceFilter(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   int // expected number of includes
	}{
		{
			name: "Single include",
			config: Config{
				Include: "users",
			},
			want: 1,
		},
		{
			name: "Multiple includes",
			config: Config{
				Include: "users, orders, products",
			},
			want: 3,
		},
		{
			name: "Include and exclude",
			config: Config{
				Include: "users",
				Exclude: "admin",
			},
			want: 1,
		},
		{
			name: "Pattern filter",
			config: Config{
				ResourceFilter: "^user.*",
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := buildResourceFilter(tt.config)
			if len(filter.Include) != tt.want {
				t.Errorf("Include count = %v, want %v", len(filter.Include), tt.want)
			}
			if tt.config.ResourceFilter != "" && filter.Pattern != tt.config.ResourceFilter {
				t.Errorf("Pattern = %v, want %v", filter.Pattern, tt.config.ResourceFilter)
			}
		})
	}
}

func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "Single value",
			input: "users",
			want:  []string{"users"},
		},
		{
			name:  "Multiple values",
			input: "users,orders,products",
			want:  []string{"users", "orders", "products"},
		},
		{
			name:  "With spaces",
			input: "users, orders , products",
			want:  []string{"users", "orders", "products"},
		},
		{
			name:  "Empty string",
			input: "",
			want:  []string{},
		},
		{
			name:  "Only commas",
			input: ",,",
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommaSeparated(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("parseCommaSeparated() returned %d items, want %d", len(got), len(tt.want))
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("parseCommaSeparated()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestShowVersion(t *testing.T) {
	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showVersion()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// Check version output contains expected strings
	if !strings.Contains(output, "api-godoc version") {
		t.Error("Version output missing 'api-godoc version'")
	}
	if !strings.Contains(output, "Build time:") {
		t.Error("Version output missing 'Build time:'")
	}
	if !strings.Contains(output, "Go version:") {
		t.Error("Version output missing 'Go version:'")
	}
}

func TestProcessAPIErrors(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "Missing input file",
			config: Config{
				InputSpec: filepath.Join(tmpDir, "nonexistent.json"),
				Format:    "markdown",
			},
			wantErr: true,
		},
		{
			name: "Invalid URL",
			config: Config{
				InputSpec: "http://[invalid-url",
				Format:    "markdown",
			},
			wantErr: true,
		},
		{
			name: "Valid spec with markdown format",
			config: Config{
				InputSpec: createTestSpec(t, tmpDir),
				Format:    "markdown",
			},
			wantErr: false,
		},
		{
			name: "Invalid regex filter",
			config: Config{
				InputSpec:      createTestSpec(t, tmpDir),
				Format:         "markdown",
				ResourceFilter: "[invalid-regex",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processAPI(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("processAPI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTitleCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "Hello"},
		{"WORLD", "World"},
		{"", ""},
		{"a", "A"},
		{"ABC", "Abc"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := titleCase(tt.input); got != tt.want {
				t.Errorf("titleCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Helper function to create a test OpenAPI spec
func createTestSpec(t *testing.T, dir string) string {
	spec := `{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {
			"/test": {
				"get": {"responses": {"200": {"description": "OK"}}}
			}
		}
	}`
	path := filepath.Join(dir, "test-spec.json")
	if err := os.WriteFile(path, []byte(spec), 0644); err != nil {
		t.Fatalf("Failed to create test spec: %v", err)
	}
	return path
}
