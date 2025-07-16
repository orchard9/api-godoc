package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type TestCase struct {
	Name        string
	SpecFile    string
	Args        []string
	Expected    ExpectedResult
	Description string
}

type ExpectedResult struct {
	// Content checks
	ContainsText    []string
	NotContainsText []string
	ContainsRegex   []string

	// Structural checks
	MinResources      int
	MaxResources      int
	SpecificResources []string
	MinOperations     int

	// Relationship checks
	HasRelationships  bool
	RelationshipPairs [][]string // [["projects", "tasks"], ["tasks", "dependencies"]]

	// Pattern checks
	HasPatterns      bool
	SpecificPatterns []string // ["versioning", "pagination"]

	// Diagram checks
	HasMermaidDiagram bool
	DiagramNodes      []string

	// Format specific
	OutputFormat  string
	JSONStructure map[string]string // path -> expected type, e.g. "resources" -> "array"

	// Exit code
	ExitCode int
}

type TestResult struct {
	TestCase TestCase
	Passed   bool
	Output   string
	Errors   []string
	Duration time.Duration
}

var testCases = []TestCase{
	{
		Name:        "Warden API - Complete Markdown Analysis",
		SpecFile:    "artifacts/warden.v1.swagger.json",
		Description: "Verify complete markdown output with all sections",
		Expected: ExpectedResult{
			ContainsText: []string{
				"# warden/v1/common.proto",
				"## Overview",
				"## API Statistics",
				"- **Total Resources**:",
				"- **Total Operations**:",
				"## Resources",
				"### Api-keys",
				"| Method | Path | Summary |",
				"| GET | `/v1/auth/api-keys` |",
				"| POST | `/v1/auth/api-keys` | API key management |",
				"| DELETE | `/v1/auth/api-keys/{id}` |",
				"## Detected Patterns",
				"### Versioning",
				"**Confidence**: high",
			},
			SpecificResources: []string{"Api-keys", "Auth", "Login", "Logout", "Refresh", "Register"},
			MinResources:      6,
			MinOperations:     8,
			HasPatterns:       true,
			SpecificPatterns:  []string{"versioning"},
			OutputFormat:      "markdown",
			ExitCode:          0,
		},
	},
	{
		Name:        "Forge API - Resource Relationships",
		SpecFile:    "artifacts/forge.swagger.json",
		Description: "Test relationship detection and Mermaid diagram generation",
		Expected: ExpectedResult{
			ContainsText: []string{
				"## Resource Relationships",
				"### Relationship Diagram",
				"```mermaid",
				"graph TD",
			},
			HasRelationships: true,
			RelationshipPairs: [][]string{
				{"projects", "understanding"},
				{"projects", "audits"},
				{"tasks", "dependencies"},
			},
			HasMermaidDiagram: true,
			DiagramNodes:      []string{"projects", "tasks", "dependencies", "understanding"},
			MinResources:      50,
			MinOperations:     100,
			ExitCode:          0,
		},
	},
	{
		Name:        "Forge API - JSON Structure Validation",
		SpecFile:    "artifacts/forge.swagger.json",
		Args:        []string{"-f", "json"},
		Description: "Validate JSON output structure and content",
		Expected: ExpectedResult{
			OutputFormat: "json",
			JSONStructure: map[string]string{
				"title":     "string",
				"version":   "string",
				"resources": "array",
				"patterns":  "array",
			},
			MinResources:  50,
			MinOperations: 100,
			ExitCode:      0,
		},
	},
	{
		Name:        "AI Format - Conciseness Test",
		SpecFile:    "artifacts/forge.swagger.json",
		Args:        []string{"-f", "ai"},
		Description: "Verify AI format is concise and well-structured",
		Expected: ExpectedResult{
			ContainsText: []string{
				"API: Forge Development System API",
				"Stats: 64 resources, 156 operations",
				"RESOURCES:",
				"tasks (",
				"projects (",
			},
			ContainsRegex: []string{
				`tasks \(\d+ ops\)`,
				`projects.*has_many`,
			},
			OutputFormat: "ai",
			ExitCode:     0,
		},
	},
	{
		Name:        "Resource Filtering - Complex Include",
		SpecFile:    "artifacts/forge.swagger.json",
		Args:        []string{"--include", "tasks"},
		Description: "Test filtering includes only specified resources",
		Expected: ExpectedResult{
			SpecificResources: []string{"Tasks"},
			NotContainsText: []string{
				"### Agents",
				"### Docker",
				"### Git",
				"### Releases",
			},
			MinResources: 1,
			MaxResources: 10, // Include relationship sections and patterns
			ExitCode:     0,
		},
	},
	{
		Name:        "Pattern Detection - Forge API",
		SpecFile:    "artifacts/forge.swagger.json",
		Description: "Verify pattern detection in complex API",
		Expected: ExpectedResult{
			HasPatterns:      true,
			SpecificPatterns: []string{"versioning"},
			ContainsText: []string{
				"### Versioning",
				"API uses URL path versioning",
			},
			ExitCode: 0,
		},
	},
	{
		Name:        "Schema Reduction Levels",
		SpecFile:    "artifacts/forge.swagger.json",
		Args:        []string{"-s", "essential"},
		Description: "Test schema reduction flag is accepted",
		Expected: ExpectedResult{
			ContainsText: []string{
				"## Resources",
			},
			ExitCode: 0,
		},
	},
	{
		Name:        "Operation Details - Markdown",
		SpecFile:    "artifacts/forge.swagger.json",
		Args:        []string{"--include", "tasks"},
		Description: "Verify detailed operation information",
		Expected: ExpectedResult{
			ContainsText: []string{
				"### Tasks",
				"| Method | Path | Summary |",
				"| GET | `/api/v1/tasks` | List tasks |",
				"| POST | `/api/v1/tasks` | CreateTask creates a fully specified task",
				"| GET | `/api/v1/tasks/{taskId}` | Get task by ID |",
				"| PUT | `/api/v1/tasks/{task.id}` | UpdateTask updates a task's content",
				"| DELETE | `/api/v1/tasks/{taskId}` | DeleteTask removes a task",
				"#### GET /api/v1/tasks",
				"#### GET /api/v1/tasks/{taskId}",
			},
			ExitCode: 0,
		},
	},
	{
		Name:        "Regex Pattern Filter",
		SpecFile:    "artifacts/forge.swagger.json",
		Args:        []string{"--filter", "^tasks?$"},
		Description: "Test regex filtering for resources",
		Expected: ExpectedResult{
			SpecificResources: []string{"Task", "Tasks"},
			NotContainsText: []string{
				"### TaskBranches", // Should not match
				"### Projects",     // Should not match
			},
			MinResources: 2,
			MaxResources: 10, // Account for relationship sections and patterns
			ExitCode:     0,
		},
	},
	{
		Name:        "Statistics Accuracy",
		SpecFile:    "artifacts/warden.v1.swagger.json",
		Args:        []string{"-f", "json"},
		Description: "Verify statistics are calculated correctly",
		Expected: ExpectedResult{
			OutputFormat: "json",
			ContainsRegex: []string{
				`"totalResources":\s*6`,
				`"totalOperations":\s*16`,
				`"totalEndpoints":\s*6`,
			},
			ExitCode: 0,
		},
	},
	{
		Name:        "Relationship Strength Detection",
		SpecFile:    "artifacts/forge.swagger.json",
		Args:        []string{"-f", "json"},
		Description: "Verify relationship strength indicators",
		Expected: ExpectedResult{
			OutputFormat: "json",
			ContainsRegex: []string{
				`"strength":\s*"strong"`,
				`"type":\s*"has_many"`,
				`"type":\s*"belongs_to"`,
			},
			ExitCode: 0,
		},
	},
	{
		Name:        "Error Handling - Invalid Regex",
		SpecFile:    "artifacts/forge.swagger.json",
		Args:        []string{"--filter", "[invalid(regex"},
		Description: "Test error handling for invalid regex",
		Expected: ExpectedResult{
			ExitCode:     1,
			ContainsText: []string{"error", "regex"},
		},
	},
}

func main() {
	fmt.Println("=== API GoDoc UAT Runner ===")
	fmt.Printf("Running %d test cases...\n\n", len(testCases))

	// Build the binary first
	if err := buildBinary(); err != nil {
		log.Fatalf("Failed to build binary: %v", err)
	}

	results := runTests()

	// Print summary
	printSummary(results)

	// Exit with failure if any tests failed
	for _, result := range results {
		if !result.Passed {
			os.Exit(1)
		}
	}
}

func buildBinary() error {
	fmt.Println("Building api-godoc binary...")
	cmd := exec.Command("go", "build", "-o", "../build/api-godoc", "../cmd/api-godoc")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %v\nOutput: %s", err, output)
	}
	fmt.Println("Build successful!")
	return nil
}

func runTests() []TestResult {
	var results []TestResult

	for i, tc := range testCases {
		fmt.Printf("[%d/%d] Running: %s\n", i+1, len(testCases), tc.Name)
		result := runTest(tc)
		results = append(results, result)

		if result.Passed {
			fmt.Printf("✓ PASSED (%.2fs)\n", result.Duration.Seconds())
		} else {
			fmt.Printf("✗ FAILED (%.2fs)\n", result.Duration.Seconds())
			for _, err := range result.Errors {
				fmt.Printf("  - %s\n", err)
			}
		}
		fmt.Println()
	}

	return results
}

func runTest(tc TestCase) TestResult {
	start := time.Now()
	result := TestResult{
		TestCase: tc,
		Passed:   true,
		Errors:   []string{},
	}

	// Build command
	args := tc.Args
	if tc.SpecFile != "" {
		args = append(args, tc.SpecFile)
	}

	cmd := exec.Command("../build/api-godoc", args...)
	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(start)

	// Check exit code
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	if exitCode != tc.Expected.ExitCode {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("expected exit code %d, got %d", tc.Expected.ExitCode, exitCode))
		if tc.Expected.ExitCode != 0 {
			// For error cases, we might still want to check error messages
			result.Output = string(output)
		} else {
			return result
		}
	}

	// For successful runs, perform detailed checks
	if tc.Expected.ExitCode == 0 {
		validateOutput(&result, tc.Expected)
	}

	return result
}

func validateOutput(result *TestResult, expected ExpectedResult) {
	output := result.Output

	// Check contains text
	for _, text := range expected.ContainsText {
		if !strings.Contains(output, text) {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("missing expected text: %q", text))
		}
	}

	// Check not contains text
	for _, text := range expected.NotContainsText {
		if strings.Contains(output, text) {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("contains unexpected text: %q", text))
		}
	}

	// Check regex patterns
	for _, pattern := range expected.ContainsRegex {
		re, err := regexp.Compile(pattern)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("invalid test regex: %s", pattern))
			continue
		}
		if !re.MatchString(output) {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("output doesn't match regex: %s", pattern))
		}
	}

	// JSON specific validation
	if expected.OutputFormat == "json" {
		validateJSON(result, expected)
	}

	// Markdown specific validation
	if expected.OutputFormat == "markdown" || expected.OutputFormat == "" {
		validateMarkdown(result, expected)
	}

	// Check specific resources
	if len(expected.SpecificResources) > 0 {
		for _, resource := range expected.SpecificResources {
			// Look for resource in different formats
			patterns := []string{
				fmt.Sprintf("### %s", titleCase(resource)),
				fmt.Sprintf(`"%s"`, resource),
				fmt.Sprintf("### %s\n", titleCase(resource)),
				fmt.Sprintf("### %s", resource), // exact match
			}
			found := false
			for _, pattern := range patterns {
				if strings.Contains(output, pattern) {
					found = true
					break
				}
			}
			if !found {
				result.Passed = false
				result.Errors = append(result.Errors, fmt.Sprintf("missing expected resource: %s", resource))
			}
		}
	}

	// Check patterns
	if expected.HasPatterns && len(expected.SpecificPatterns) > 0 {
		for _, pattern := range expected.SpecificPatterns {
			if !strings.Contains(strings.ToLower(output), strings.ToLower(pattern)) {
				result.Passed = false
				result.Errors = append(result.Errors, fmt.Sprintf("missing expected pattern: %s", pattern))
			}
		}
	}

	// Check Mermaid diagram
	if expected.HasMermaidDiagram {
		if !strings.Contains(output, "```mermaid") {
			result.Passed = false
			result.Errors = append(result.Errors, "missing Mermaid diagram")
		}

		// Check diagram nodes
		for _, node := range expected.DiagramNodes {
			if !strings.Contains(output, node) {
				result.Passed = false
				result.Errors = append(result.Errors, fmt.Sprintf("Mermaid diagram missing node: %s", node))
			}
		}
	}

	// Check relationship pairs
	if len(expected.RelationshipPairs) > 0 {
		for _, pair := range expected.RelationshipPairs {
			// Look for relationship in various formats
			patterns := []string{
				fmt.Sprintf("%s.*→.*%s", pair[0], pair[1]),
				fmt.Sprintf("%s.*has_many.*%s", pair[0], pair[1]),
				fmt.Sprintf("%s.*belongs_to.*%s", pair[1], pair[0]),
			}
			found := false
			for _, pattern := range patterns {
				re := regexp.MustCompile(pattern)
				if re.MatchString(output) {
					found = true
					break
				}
			}
			if !found {
				result.Passed = false
				result.Errors = append(result.Errors, fmt.Sprintf("missing relationship: %s -> %s", pair[0], pair[1]))
			}
		}
	}
}

func validateJSON(result *TestResult, expected ExpectedResult) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(result.Output), &data); err != nil {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("invalid JSON: %v", err))
		return
	}

	// Check structure
	for path, expectedType := range expected.JSONStructure {
		value, exists := data[path]
		if !exists {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("JSON missing field: %s", path))
			continue
		}

		actualType := getJSONType(value)
		if actualType != expectedType {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("JSON field %s: expected %s, got %s", path, expectedType, actualType))
		}
	}

	// Check resources
	if resources, ok := data["resources"].([]interface{}); ok {
		if expected.MinResources > 0 && len(resources) < expected.MinResources {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("expected at least %d resources, got %d", expected.MinResources, len(resources)))
		}
		if expected.MaxResources > 0 && len(resources) > expected.MaxResources {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("expected at most %d resources, got %d", expected.MaxResources, len(resources)))
		}

		// Count total operations
		totalOps := 0
		for _, r := range resources {
			if resource, ok := r.(map[string]interface{}); ok {
				if ops, ok := resource["operations"].([]interface{}); ok {
					totalOps += len(ops)
				}
			}
		}
		if expected.MinOperations > 0 && totalOps < expected.MinOperations {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("expected at least %d operations, got %d", expected.MinOperations, totalOps))
		}
	}

	// Check statistics
	if stats, ok := data["statistics"].(map[string]interface{}); ok {
		if totalRes, ok := stats["totalResources"].(float64); ok && expected.MinResources > 0 {
			if int(totalRes) < expected.MinResources {
				result.Passed = false
				result.Errors = append(result.Errors, fmt.Sprintf("statistics shows %d resources, expected at least %d", int(totalRes), expected.MinResources))
			}
		}
	}
}

func validateMarkdown(result *TestResult, expected ExpectedResult) {
	output := result.Output

	// Count resources by looking for ### headings
	resourceMatches := regexp.MustCompile(`(?m)^### [A-Z]`).FindAllString(output, -1)
	resourceCount := len(resourceMatches)

	if expected.MinResources > 0 && resourceCount < expected.MinResources {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("found %d resources, expected at least %d", resourceCount, expected.MinResources))
	}

	if expected.MaxResources > 0 && resourceCount > expected.MaxResources {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("found %d resources, expected at most %d", resourceCount, expected.MaxResources))
	}

	// Count operations in tables
	operationMatches := regexp.MustCompile(`\| (GET|POST|PUT|DELETE|PATCH) \|`).FindAllString(output, -1)
	operationCount := len(operationMatches)

	if expected.MinOperations > 0 && operationCount < expected.MinOperations {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("found %d operations, expected at least %d", operationCount, expected.MinOperations))
	}
}

func getJSONType(v interface{}) string {
	switch v.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	case nil:
		return "null"
	default:
		return "unknown"
	}
}

func printSummary(results []TestResult) {
	fmt.Println("\n=== Test Summary ===")

	passed := 0
	failed := 0
	totalDuration := time.Duration(0)
	var failedTests []TestResult

	for _, result := range results {
		if result.Passed {
			passed++
		} else {
			failed++
			failedTests = append(failedTests, result)
		}
		totalDuration += result.Duration
	}

	fmt.Printf("\nTotal: %d tests\n", len(results))
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Duration: %.2fs\n", totalDuration.Seconds())

	if failed > 0 {
		fmt.Println("\nFailed tests:")
		for _, result := range failedTests {
			fmt.Printf("\n❌ %s\n", result.TestCase.Name)
			for _, err := range result.Errors {
				fmt.Printf("   - %s\n", err)
			}
			if len(result.Output) > 0 && len(result.Output) < 500 {
				fmt.Printf("   Output preview:\n%s\n", indent(result.Output, "     "))
			}
		}
	}

	if passed == len(results) {
		fmt.Println("\n✅ All tests passed!")
	} else {
		fmt.Println("\n❌ Some tests failed")
	}
}

func indent(text, prefix string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if len(line) > 100 {
			lines[i] = prefix + line[:97] + "..."
		} else {
			lines[i] = prefix + line
		}
	}
	if len(lines) > 10 {
		lines = lines[:10]
		lines = append(lines, prefix+"...")
	}
	return strings.Join(lines, "\n")
}

// titleCase converts a string to title case, handling hyphens properly
func titleCase(s string) string {
	if s == "" {
		return s
	}

	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
		}
	}
	return strings.Join(parts, "-")
}
