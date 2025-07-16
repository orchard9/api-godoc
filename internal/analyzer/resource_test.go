package analyzer

import (
	"strings"
	"testing"

	"github.com/orchard9/api-godoc/internal/parser"
)

func TestAnalyzePathPatterns(t *testing.T) {
	ra := NewResourceAnalyzer()

	// Test paths similar to warden API
	paths := map[string]parser.PathItem{
		"/v1/auth/api-keys":      {},
		"/v1/auth/api-keys/{id}": {},
		"/v1/auth/login":         {},
		"/v1/auth/logout":        {},
		"/v1/auth/refresh":       {},
		"/v1/auth/register":      {},
	}

	patterns := ra.analyzePathPatterns(paths)

	// api-keys should have variables (true resource)
	if apiKeysPattern := patterns["api-keys"]; apiKeysPattern == nil {
		t.Error("Expected api-keys pattern to exist")
	} else if !apiKeysPattern.HasVariables {
		t.Error("Expected api-keys to have variables (true resource)")
	}

	// auth should have variables (because it contains api-keys with variables)
	if authPattern := patterns["auth"]; authPattern == nil {
		t.Error("Expected auth pattern to exist")
	} else if !authPattern.HasVariables {
		t.Error("Expected auth to have variables (contains api-keys)")
	}

	// login should not have variables (action)
	if loginPattern := patterns["login"]; loginPattern == nil {
		t.Error("Expected login pattern to exist")
	} else if loginPattern.HasVariables {
		t.Error("Expected login to not have variables (action)")
	}
}

func TestExtractResourcesForPath(t *testing.T) {
	ra := NewResourceAnalyzer()

	// Mock patterns based on warden API and deeper structures
	patterns := map[string]*PathPattern{
		"auth": {
			HasVariables: true,
			Paths:        []string{"/v1/auth/api-keys", "/v1/auth/api-keys/{id}", "/v1/auth/login"},
		},
		"api-keys": {
			HasVariables: true,
			Paths:        []string{"/v1/auth/api-keys", "/v1/auth/api-keys/{id}"},
		},
		"login": {
			HasVariables: false,
			Paths:        []string{"/v1/auth/login"},
		},
		"logout": {
			HasVariables: false,
			Paths:        []string{"/v1/auth/logout"},
		},
		"users": {
			HasVariables: true,
			Paths:        []string{"/users", "/users/{id}"},
		},
		"posts": {
			HasVariables: true,
			Paths:        []string{"/users/{id}/posts", "/users/{id}/posts/{post_id}"},
		},
		"organizations": {
			HasVariables: true,
			Paths:        []string{"/organizations", "/organizations/{org_id}"},
		},
		"projects": {
			HasVariables: true,
			Paths:        []string{"/organizations/{org_id}/projects", "/organizations/{org_id}/projects/{project_id}"},
		},
		"tasks": {
			HasVariables: true,
			Paths:        []string{"/organizations/{org_id}/projects/{project_id}/tasks", "/organizations/{org_id}/projects/{project_id}/tasks/{task_id}"},
		},
		"comments": {
			HasVariables: true,
			Paths:        []string{"/organizations/{org_id}/projects/{project_id}/tasks/{task_id}/comments", "/organizations/{org_id}/projects/{project_id}/tasks/{task_id}/comments/{comment_id}"},
		},
		"archive": {
			HasVariables: false,
			Paths:        []string{"/organizations/{org_id}/projects/{project_id}/tasks/{task_id}/archive"},
		},
		"notifications": {
			HasVariables: false,
			Paths:        []string{"/users/{id}/notifications/mark-read"},
		},
		"mark-read": {
			HasVariables: false,
			Paths:        []string{"/users/{id}/notifications/mark-read"},
		},
	}

	tests := []struct {
		path     string
		expected []string
	}{
		// Original 2-level cases
		{"/v1/auth/api-keys", []string{"api-keys"}},                 // True resource with variables
		{"/v1/auth/api-keys/{id}", []string{"api-keys"}},            // True resource with variables
		{"/v1/auth/login", []string{"auth"}},                        // Action, group under parent
		{"/v1/auth/logout", []string{"auth"}},                       // Action, group under parent
		{"/v1/auth/refresh", []string{"auth"}},                      // Action, group under parent
		{"/v1/auth/register", []string{"auth"}},                     // Action, group under parent
		{"/users/{id}/posts", []string{"users", "posts"}},           // Nested resource
		{"/users/{id}/posts/{post_id}", []string{"users", "posts"}}, // Nested resource with ID

		// Deep nested resource cases
		{"/organizations/{org_id}/projects", []string{"organizations", "projects"}},                                                                         // 2-level nested
		{"/organizations/{org_id}/projects/{project_id}", []string{"organizations", "projects"}},                                                            // 2-level nested with ID
		{"/organizations/{org_id}/projects/{project_id}/tasks", []string{"organizations", "projects", "tasks"}},                                             // 3-level nested
		{"/organizations/{org_id}/projects/{project_id}/tasks/{task_id}", []string{"organizations", "projects", "tasks"}},                                   // 3-level nested with ID
		{"/organizations/{org_id}/projects/{project_id}/tasks/{task_id}/comments", []string{"organizations", "projects", "tasks", "comments"}},              // 4-level nested
		{"/organizations/{org_id}/projects/{project_id}/tasks/{task_id}/comments/{comment_id}", []string{"organizations", "projects", "tasks", "comments"}}, // 4-level nested with ID

		// Action cases at different depths
		{"/organizations/{org_id}/projects/{project_id}/tasks/{task_id}/archive", []string{"organizations", "projects", "tasks"}}, // Action at 4th level, return full chain
		{"/users/{id}/notifications/mark-read", []string{"users", "notifications"}},                                               // Action at 3rd level, return chain

		// Single resource cases
		{"/users", []string{"users"}},      // Single resource
		{"/users/{id}", []string{"users"}}, // Single resource with ID
	}

	for _, test := range tests {
		result := ra.extractResourcesForPath(test.path, patterns)
		if !equalSlices(result, test.expected) {
			t.Errorf("For path %s, expected %v but got %v", test.path, test.expected, result)
		}
	}
}

func TestExtractResourceNames(t *testing.T) {
	ra := NewResourceAnalyzer()

	tests := []struct {
		path     string
		expected []string
	}{
		{"/v1/auth/api-keys", []string{"auth", "api-keys"}},
		{"/v1/auth/api-keys/{id}", []string{"auth", "api-keys"}},
		{"/v1/auth/login", []string{"auth", "login"}},
		{"/api/v1/tasks/{id}:move", []string{"tasks", "move"}},
		{"/v1/projects/{id}/audits", []string{"projects", "audits"}},
		{"/organizations/{org_id}/projects/{project_id}/tasks/{task_id}/comments", []string{"organizations", "projects", "tasks", "comments"}},
		{"/api/v1/organizations/{org_id}/projects/{project_id}/tasks/{task_id}/comments/{comment_id}/reactions", []string{"organizations", "projects", "tasks", "comments", "reactions"}},
	}

	for _, test := range tests {
		result := ra.extractResourceNames(test.path)
		if !equalSlices(result, test.expected) {
			t.Errorf("For path %s, expected %v but got %v", test.path, test.expected, result)
		}
	}
}

func TestHasVariablesBetweenResources(t *testing.T) {
	ra := NewResourceAnalyzer()

	tests := []struct {
		path              string
		resourcePositions []int
		expected          bool
		description       string
	}{
		{
			path:              "/v1/auth/api-keys",
			resourcePositions: []int{1, 2}, // auth at pos 1, api-keys at pos 2
			expected:          false,
			description:       "No variables between auth and api-keys",
		},
		{
			path:              "/users/{id}/posts",
			resourcePositions: []int{0, 2}, // users at pos 0, posts at pos 2
			expected:          true,
			description:       "Variable {id} between users and posts",
		},
		{
			path:              "/organizations/{org_id}/projects/{project_id}/tasks",
			resourcePositions: []int{0, 2, 4}, // organizations, projects, tasks
			expected:          true,
			description:       "Variables between nested resources",
		},
		{
			path:              "/single/resource",
			resourcePositions: []int{0, 1},
			expected:          false,
			description:       "No variables between adjacent resources",
		},
	}

	for _, test := range tests {
		pathSegments := strings.Split(strings.Trim(test.path, "/"), "/")
		result := ra.hasVariablesBetweenResources(pathSegments, test.resourcePositions)
		if result != test.expected {
			t.Errorf("For path %s (%s), expected %v but got %v", test.path, test.description, test.expected, result)
		}
	}
}

func TestFindResourcePositions(t *testing.T) {
	ra := NewResourceAnalyzer()

	tests := []struct {
		path     string
		expected []int
	}{
		{"/v1/auth/api-keys", []int{1, 2}},                                      // auth at 1, api-keys at 2
		{"/users/{id}/posts/{post_id}", []int{0, 2}},                            // users at 0, posts at 2
		{"/organizations/{org_id}/projects/{project_id}/tasks", []int{0, 2, 4}}, // organizations, projects, tasks
		{"/api/v1/single", []int{2}},                                            // single at 2 (skip api, v1)
	}

	for _, test := range tests {
		pathSegments := strings.Split(strings.Trim(test.path, "/"), "/")
		result := ra.findResourcePositions(pathSegments)
		if !equalIntSlices(result, test.expected) {
			t.Errorf("For path %s, expected %v but got %v", test.path, test.expected, result)
		}
	}
}

// Helper function to compare slices
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Helper function to compare integer slices
func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
