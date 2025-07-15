package reporter

import (
	"strings"
	"testing"
	"time"

	"github.com/orchard9/pg-goapi/pkg/models"
)

func TestMarkdownGeneration(t *testing.T) {
	// Create test analysis data
	analysis := &models.APIAnalysis{
		Title:       "Test API",
		Version:     "1.0.0",
		Description: "A test API for validation",
		BaseURL:     "https://api.example.com",
		SpecType:    "OpenAPI 3.x",
		GeneratedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Summary: models.AnalysisStat{
			TotalResources:   2,
			TotalOperations:  4,
			TotalEndpoints:   4,
			ResourceCoverage: 100,
		},
		Resources: []models.Resource{
			{
				Name:        "users",
				Description: "User management operations",
				Operations: []models.Operation{
					{
						Method:      "GET",
						Path:        "/users",
						Summary:     "List all users",
						Description: "Retrieve a list of all users in the system",
						Parameters: []models.Parameter{
							{
								Name:        "limit",
								In:          "query",
								Type:        "integer",
								Required:    false,
								Description: "Maximum number of users to return",
							},
						},
						Responses: []models.Response{
							{
								StatusCode:  "200",
								Description: "Successful response",
								ContentType: "application/json",
							},
						},
					},
					{
						Method:  "POST",
						Path:    "/users",
						Summary: "Create a new user",
						RequestBody: &models.RequestBody{
							Description: "User data",
							Required:    true,
							ContentType: "application/json",
						},
					},
				},
				Relationships: []models.Relationship{
					{
						Resource:    "posts",
						Type:        "has_many",
						Via:         "path hierarchy",
						Description: "users contains multiple posts resources",
						Strength:    "strong",
					},
				},
			},
			{
				Name:        "posts",
				Description: "Blog post operations",
				Operations: []models.Operation{
					{
						Method:  "GET",
						Path:    "/users/{userId}/posts",
						Summary: "Get user posts",
					},
					{
						Method:  "POST",
						Path:    "/users/{userId}/posts",
						Summary: "Create user post",
					},
				},
				Relationships: []models.Relationship{
					{
						Resource:    "users",
						Type:        "belongs_to",
						Via:         "path hierarchy",
						Description: "posts belongs to a users resource",
						Strength:    "strong",
					},
				},
			},
		},
		Patterns: []models.Pattern{
			{
				Type:        "pagination",
				Description: "API uses limit parameter for pagination",
				Confidence:  "high",
				Impact:      "Affects list operations",
				Examples:    []string{"/users?limit=10"},
			},
		},
	}

	reporter := New()
	result, err := reporter.Generate(analysis, "markdown")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify key sections are present
	expectedSections := []string{
		"# Test API",
		"## Overview",
		"## API Statistics",
		"## Resources",
		"### Users",
		"### Posts",
		"## Resource Relationships",
		"### Users Relationships",
		"### Posts Relationships",
		"## Detected Patterns",
		"### Pagination",
	}

	for _, section := range expectedSections {
		if !strings.Contains(result, section) {
			t.Errorf("Expected section '%s' not found in markdown output", section)
		}
	}

	// Verify specific content
	if !strings.Contains(result, "Total Resources**: 2") {
		t.Error("Expected resource count not found")
	}

	if !strings.Contains(result, "GET | `/users` | List all users") {
		t.Error("Expected operation table entry not found")
	}

	if !strings.Contains(result, "**has_many** posts (strong strength") {
		t.Error("Expected relationship not found")
	}

	if !strings.Contains(result, "**Confidence**: high") {
		t.Error("Expected pattern confidence not found")
	}
}

func TestJSONGeneration(t *testing.T) {
	analysis := &models.APIAnalysis{
		Title:   "Test API",
		Version: "1.0.0",
		Summary: models.AnalysisStat{
			TotalResources: 1,
		},
		Resources: []models.Resource{
			{Name: "test", Operations: []models.Operation{}},
		},
	}

	reporter := New()
	result, err := reporter.Generate(analysis, "json")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify it's valid JSON with expected content
	if !strings.Contains(result, `"title": "Test API"`) {
		t.Error("Expected JSON title not found")
	}

	if !strings.Contains(result, `"version": "1.0.0"`) {
		t.Error("Expected JSON version not found")
	}

	if !strings.Contains(result, `"totalResources": 1`) {
		t.Error("Expected JSON resource count not found")
	}
}

func TestAIOptimizedGeneration(t *testing.T) {
	analysis := &models.APIAnalysis{
		Title:    "Test API",
		Version:  "1.0.0",
		SpecType: "OpenAPI 3.x",
		Summary: models.AnalysisStat{
			TotalResources:  2,
			TotalOperations: 3,
			TotalEndpoints:  3,
		},
		Resources: []models.Resource{
			{
				Name: "users",
				Operations: []models.Operation{
					{Method: "GET", Path: "/users", Summary: "List users"},
					{Method: "POST", Path: "/users", Summary: "Create user"},
				},
				Relationships: []models.Relationship{
					{Resource: "posts", Type: "has_many"},
				},
			},
			{
				Name: "posts",
				Operations: []models.Operation{
					{Method: "GET", Path: "/posts", Summary: "List posts"},
				},
			},
		},
	}

	reporter := New()
	result, err := reporter.Generate(analysis, "ai")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify condensed format
	expectedContent := []string{
		"API: Test API v1.0.0 (OpenAPI 3.x)",
		"Stats: 2 resources, 3 operations, 3 endpoints",
		"RESOURCES:",
		"- users (2 ops) -> has_many:posts",
		"- posts (1 ops)",
		"KEY OPERATIONS:",
		"- GET /users",
		"- POST /users",
		"- GET /posts",
	}

	for _, content := range expectedContent {
		if !strings.Contains(result, content) {
			t.Errorf("Expected AI content '%s' not found in output", content)
		}
	}
}

func TestOperationSorting(t *testing.T) {
	operations := []models.Operation{
		{Method: "DELETE", Path: "/users/1"},
		{Method: "GET", Path: "/users"},
		{Method: "POST", Path: "/users"},
		{Method: "PUT", Path: "/users/1"},
		{Method: "PATCH", Path: "/users/1"},
	}

	resource := models.Resource{
		Name:       "users",
		Operations: operations,
	}

	analysis := &models.APIAnalysis{
		Title:     "Test API",
		Version:   "1.0.0",
		Resources: []models.Resource{resource},
	}

	reporter := New()
	result, err := reporter.Generate(analysis, "markdown")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Find the operations table
	lines := strings.Split(result, "\n")
	var tableStart int
	for i, line := range lines {
		if strings.Contains(line, "| Method | Path | Summary |") {
			tableStart = i + 2 // Skip the header separator
			break
		}
	}

	if tableStart == 0 {
		t.Fatal("Operations table not found")
	}

	// Verify operations are sorted: GET, POST, PUT, PATCH, DELETE
	expectedOrder := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for i, expectedMethod := range expectedOrder {
		if tableStart+i >= len(lines) {
			t.Fatal("Not enough table rows")
		}

		line := lines[tableStart+i]
		if !strings.HasPrefix(line, "| "+expectedMethod+" ") {
			t.Errorf("Expected method %s at position %d, got line: %s", expectedMethod, i, line)
		}
	}
}

func TestMethodOrder(t *testing.T) {
	reporter := &reporter{}

	testCases := []struct {
		method string
		order  int
	}{
		{"GET", 0},
		{"POST", 1},
		{"PUT", 2},
		{"PATCH", 3},
		{"DELETE", 4},
		{"HEAD", 5},
		{"UNKNOWN", 5},
	}

	for _, tc := range testCases {
		result := reporter.methodOrder(tc.method)
		if result != tc.order {
			t.Errorf("Expected order %d for method %s, got %d", tc.order, tc.method, result)
		}
	}
}

func TestKeyOperationDetection(t *testing.T) {
	reporter := &reporter{}

	testCases := []struct {
		operation models.Operation
		isKey     bool
		name      string
	}{
		{
			operation: models.Operation{Method: "GET", Path: "/users"},
			isKey:     true,
			name:      "GET operation",
		},
		{
			operation: models.Operation{Method: "POST", Path: "/users"},
			isKey:     true,
			name:      "POST operation",
		},
		{
			operation: models.Operation{Method: "HEAD", Path: "/health"},
			isKey:     false,
			name:      "HEAD operation",
		},
		{
			operation: models.Operation{Method: "OPTIONS", Path: "/api", Summary: "Very detailed summary"},
			isKey:     true,
			name:      "Operation with meaningful summary",
		},
		{
			operation: models.Operation{Method: "TRACE", Path: "/debug", Summary: "Short"},
			isKey:     false,
			name:      "Operation with short summary",
		},
	}

	for _, tc := range testCases {
		result := reporter.isKeyOperation(tc.operation)
		if result != tc.isKey {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.isKey, result)
		}
	}
}

func TestEmptyAnalysis(t *testing.T) {
	analysis := &models.APIAnalysis{
		Title:     "Empty API",
		Version:   "1.0.0",
		Resources: []models.Resource{},
	}

	reporter := New()
	result, err := reporter.Generate(analysis, "markdown")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Should still have basic structure
	if !strings.Contains(result, "# Empty API") {
		t.Error("Expected title not found")
	}

	if !strings.Contains(result, "## Resources") {
		t.Error("Expected resources section not found")
	}

	// Should not have relationships section for empty analysis
	if strings.Contains(result, "## Resource Relationships") {
		t.Error("Should not have relationships section for empty analysis")
	}
}
