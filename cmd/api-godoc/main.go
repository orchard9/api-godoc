package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/orchard9/api-godoc/internal/analyzer"
	"github.com/orchard9/api-godoc/internal/parser"
	"github.com/orchard9/api-godoc/internal/reporter"
	"github.com/orchard9/api-godoc/pkg/models"
)

// Build-time variables (set by Makefile)
var (
	version   = "dev"
	buildTime = "unknown"
	buildHash = "unknown"
)

// Config holds the CLI configuration
type Config struct {
	InputSpec      string
	OutputFile     string
	Format         string
	SchemaLevel    string
	Include        string
	Exclude        string
	ResourceFilter string
	Verbose        bool
	ShowVersion    bool
	ShowHelp       bool
}

func main() {
	config := parseFlags()

	if config.ShowVersion {
		showVersion()
		return
	}

	if config.ShowHelp {
		showHelp()
		return
	}

	if config.InputSpec == "" {
		fmt.Fprintf(os.Stderr, "Error: OpenAPI specification file or URL is required\n\n")
		showHelp()
		os.Exit(1)
	}

	// Set up logging
	if config.Verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(os.Stderr)
	}

	// Process the API specification
	if err := processAPI(config); err != nil {
		log.Fatalf("Error processing API: %v", err)
	}
}

// parseFlags parses command-line flags and returns configuration
func parseFlags() Config {
	var config Config

	// Define flags
	flag.StringVar(&config.OutputFile, "output", "", "Output file (default: stdout)")
	flag.StringVar(&config.OutputFile, "o", "", "Output file (default: stdout)")
	flag.StringVar(&config.Format, "format", "markdown", "Output format: markdown, json, ai")
	flag.StringVar(&config.Format, "f", "markdown", "Output format: markdown, json, ai")
	flag.StringVar(&config.SchemaLevel, "schema", "standard", "Schema detail level: essential, standard, full")
	flag.StringVar(&config.SchemaLevel, "s", "standard", "Schema detail level: essential, standard, full")
	flag.StringVar(&config.Include, "include", "", "Comma-separated list of resources to include")
	flag.StringVar(&config.Include, "i", "", "Comma-separated list of resources to include")
	flag.StringVar(&config.Exclude, "exclude", "", "Comma-separated list of resources to exclude")
	flag.StringVar(&config.Exclude, "e", "", "Comma-separated list of resources to exclude")
	flag.StringVar(&config.ResourceFilter, "filter", "", "Regex pattern to filter resources")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose logging")
	flag.BoolVar(&config.Verbose, "v", false, "Enable verbose logging")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version information")
	flag.BoolVar(&config.ShowHelp, "help", false, "Show help message")
	flag.BoolVar(&config.ShowHelp, "h", false, "Show help message")

	// Custom usage function
	flag.Usage = showHelp

	// Parse flags
	flag.Parse()

	// Get positional argument (input spec)
	if flag.NArg() > 0 {
		config.InputSpec = flag.Arg(0)
	}

	return config
}

// processAPI processes the OpenAPI specification and generates output
func processAPI(config Config) error {
	if config.Verbose {
		log.Printf("Processing OpenAPI specification: %s", config.InputSpec)
		log.Printf("Output format: %s", config.Format)
		if config.OutputFile != "" {
			log.Printf("Output file: %s", config.OutputFile)
		}
	}

	// Initialize components
	p := parser.New()
	resourceAnalyzer := analyzer.NewResourceAnalyzer()
	relationshipDetector := analyzer.NewRelationshipDetector()
	patternDetector := analyzer.NewPatternDetector()
	schemaReducer := analyzer.NewSchemaReducer()
	rep := reporter.New()

	// Parse the OpenAPI specification
	var spec *parser.OpenAPISpec
	var err error

	// Check if input is a URL or file path
	if strings.HasPrefix(config.InputSpec, "http://") || strings.HasPrefix(config.InputSpec, "https://") {
		spec, err = p.ParseURL(config.InputSpec)
	} else {
		spec, err = p.ParseFile(config.InputSpec)
	}

	if err != nil {
		return fmt.Errorf("failed to parse specification: %w", err)
	}

	// Extract resources
	if config.Verbose {
		log.Println("Extracting resources from specification")
	}
	resources := resourceAnalyzer.ExtractResources(spec)

	// Apply resource filtering
	if config.Include != "" || config.Exclude != "" || config.ResourceFilter != "" {
		if config.Verbose {
			log.Println("Applying resource filters")
		}
		filter := buildResourceFilter(config)
		if err := filter.Validate(); err != nil {
			return fmt.Errorf("invalid resource filter: %w", err)
		}
		filterer := analyzer.NewResourceFilterer()
		resources = filterer.FilterResources(resources, filter)
		if config.Verbose {
			log.Printf("Filtered to %d resources", len(resources))
		}
	}

	// Extract fields for resources
	if config.Verbose {
		log.Println("Extracting resource fields with schema reduction")
	}
	extractResourceFields(resources, spec, schemaReducer, config.SchemaLevel)

	// Detect relationships
	if config.Verbose {
		log.Println("Detecting resource relationships")
	}
	relationshipDetector.DetectRelationships(resources, spec)

	// Detect patterns
	if config.Verbose {
		log.Println("Detecting API patterns")
	}
	patterns := patternDetector.DetectPatterns(spec)

	// Create analysis result
	analysis := &models.APIAnalysis{
		Title:       getAPITitle(spec),
		Version:     getAPIVersion(spec),
		Description: getAPIDescription(spec),
		BaseURL:     getBaseURL(spec),
		SpecType:    getSpecType(spec),
		GeneratedAt: time.Now(),
		Resources:   resources,
		Summary:     calculateSummary(resources, spec),
		Patterns:    patterns,
	}

	// Generate output
	if config.Verbose {
		log.Printf("Generating %s output", config.Format)
	}
	output, err := rep.Generate(analysis, config.Format)
	if err != nil {
		return fmt.Errorf("failed to generate output: %w", err)
	}

	// Write output
	if config.OutputFile != "" {
		if config.Verbose {
			log.Printf("Writing output to file: %s", config.OutputFile)
		}
		err = os.WriteFile(config.OutputFile, []byte(output), 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Documentation generated: %s\n", config.OutputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

// Helper functions for extracting API metadata
func getAPITitle(spec *parser.OpenAPISpec) string {
	if spec.Info.Title != "" {
		return spec.Info.Title
	}
	return "API Documentation"
}

func getAPIVersion(spec *parser.OpenAPISpec) string {
	if spec.Info.Version != "" {
		return spec.Info.Version
	}
	return "1.0.0"
}

func getAPIDescription(spec *parser.OpenAPISpec) string {
	return spec.Info.Description
}

func getBaseURL(spec *parser.OpenAPISpec) string {
	if len(spec.Servers) > 0 && spec.Servers[0].URL != "" {
		return spec.Servers[0].URL
	}
	return ""
}

func getSpecType(spec *parser.OpenAPISpec) string {
	if spec.OpenAPI != "" {
		return "OpenAPI " + spec.OpenAPI
	}
	return "Unknown"
}

func calculateSummary(resources []models.Resource, spec *parser.OpenAPISpec) models.AnalysisStat {
	totalOperations := 0
	totalEndpoints := len(spec.Paths)

	for _, resource := range resources {
		totalOperations += len(resource.Operations)
	}

	// Calculate resource coverage (what percentage of endpoints are resource operations)
	resourceCoverage := 0
	if totalEndpoints > 0 {
		resourceCoverage = (totalOperations * 100) / totalEndpoints
	}

	return models.AnalysisStat{
		TotalResources:   len(resources),
		TotalOperations:  totalOperations,
		TotalEndpoints:   totalEndpoints,
		ResourceCoverage: resourceCoverage,
	}
}

func showVersion() {
	fmt.Printf("api-godoc version %s\n", version)
	fmt.Printf("Build time: %s\n", buildTime)
	fmt.Printf("Build hash: %s\n", buildHash)
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

// titleCase converts a string to title case (replacement for deprecated strings.Title)
func titleCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

// extractResourceFields extracts fields from request/response schemas for resources
func extractResourceFields(resources []models.Resource, spec *parser.OpenAPISpec, reducer analyzer.SchemaReducer, level string) {
	// For each resource, try to find associated schemas
	for i := range resources {
		resource := &resources[i]

		// Check component schemas that match resource name
		if spec.Components != nil && spec.Components.Schemas != nil {
			// Try common schema naming patterns
			schemaNames := []string{
				resource.Name,
				titleCase(strings.ToLower(resource.Name)),
				resource.Name + "Response",
				resource.Name + "Request",
				resource.Name + "Model",
				resource.Name + "DTO",
			}

			for _, name := range schemaNames {
				if schema, ok := spec.Components.Schemas[name]; ok {
					fields := reducer.SchemaToFields(&schema, level)
					resource.Fields = mergeFields(resource.Fields, fields)
				}
			}
		}
	}
}

// mergeFields merges new fields into existing fields, avoiding duplicates
func mergeFields(existing, new []models.Field) []models.Field {
	fieldMap := make(map[string]models.Field)

	// Add existing fields
	for _, field := range existing {
		fieldMap[field.Name] = field
	}

	// Add new fields
	for _, field := range new {
		if _, exists := fieldMap[field.Name]; !exists {
			fieldMap[field.Name] = field
		}
	}

	// Convert back to slice
	result := make([]models.Field, 0, len(fieldMap))
	for _, field := range fieldMap {
		result = append(result, field)
	}

	return result
}

// buildResourceFilter creates a ResourceFilter from config
func buildResourceFilter(config Config) *analyzer.ResourceFilter {
	filter := &analyzer.ResourceFilter{
		Pattern: config.ResourceFilter,
	}

	// Parse include list
	if config.Include != "" {
		filter.Include = parseCommaSeparated(config.Include)
	}

	// Parse exclude list
	if config.Exclude != "" {
		filter.Exclude = parseCommaSeparated(config.Exclude)
	}

	return filter
}

// parseCommaSeparated splits a comma-separated string and trims spaces
func parseCommaSeparated(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func showHelp() {
	fmt.Println("API GoDoc - OpenAPI Documentation Generator")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  api-godoc [options] <openapi-spec>")
	fmt.Println("")
	fmt.Println("ARGUMENTS:")
	fmt.Println("  <openapi-spec>    OpenAPI specification file (JSON/YAML) or URL")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("  -o, --output <file>    Output file (default: api-docs.md)")
	fmt.Println("  -f, --format <format>  Output format: markdown, json, ai (default: markdown)")
	fmt.Println("  -s, --schema <level>   Schema detail level: essential, standard, full (default: standard)")
	fmt.Println("  -i, --include <list>   Comma-separated list of resources to include")
	fmt.Println("  -e, --exclude <list>   Comma-separated list of resources to exclude")
	fmt.Println("      --filter <regex>   Regex pattern to filter resources")
	fmt.Println("  -v, --verbose          Enable verbose logging")
	fmt.Println("      --version          Show version information")
	fmt.Println("  -h, --help             Show this help message")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  api-godoc api-spec.json")
	fmt.Println("  api-godoc -f json -o analysis.json api-spec.json")
	fmt.Println("  api-godoc https://api.example.com/openapi.json")
	fmt.Println("")
	fmt.Println("For more information, visit: https://github.com/orchard9/api-godoc")
}
