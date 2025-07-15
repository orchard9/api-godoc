package test

import (
	"path/filepath"
	"testing"

	"github.com/orchard9/pg-goapi/internal/analyzer"
	"github.com/orchard9/pg-goapi/internal/parser"
	"github.com/orchard9/pg-goapi/internal/reporter"
	"github.com/orchard9/pg-goapi/pkg/models"
)

// BenchmarkStripeAPI benchmarks processing of the Stripe API
func BenchmarkStripeAPI(b *testing.B) {
	fixturePath := filepath.Join("fixtures", "stripe-openapi3.json")
	
	// Initialize components once
	p := parser.New()
	resourceAnalyzer := analyzer.NewResourceAnalyzer()
	relationshipDetector := analyzer.NewRelationshipDetector()
	patternDetector := analyzer.NewPatternDetector()
	rep := reporter.New()

	// Parse once for benchmarking processing
	spec, err := p.ParseFile(fixturePath)
	if err != nil {
		b.Fatalf("Failed to parse Stripe API: %v", err)
	}

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Extract resources
		resources := resourceAnalyzer.ExtractResources(spec)
		
		// Detect relationships
		relationshipDetector.DetectRelationships(resources, spec)
		
		// Detect patterns
		patterns := patternDetector.DetectPatterns(spec)
		
		// Create analysis
		analysis := &models.APIAnalysis{
			Title:     spec.Info.Title,
			Resources: resources,
			Patterns:  patterns,
		}
		
		// Generate markdown output
		_, err := rep.Generate(analysis, "markdown")
		if err != nil {
			b.Fatalf("Failed to generate output: %v", err)
		}
	}
}

// BenchmarkParsingOnly benchmarks just the parsing step
func BenchmarkParsingOnly(b *testing.B) {
	fixturePath := filepath.Join("fixtures", "stripe-openapi3.json")
	p := parser.New()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := p.ParseFile(fixturePath)
		if err != nil {
			b.Fatalf("Failed to parse: %v", err)
		}
	}
}

// BenchmarkResourceExtraction benchmarks resource extraction
func BenchmarkResourceExtraction(b *testing.B) {
	fixturePath := filepath.Join("fixtures", "stripe-openapi3.json")
	
	p := parser.New()
	spec, err := p.ParseFile(fixturePath)
	if err != nil {
		b.Fatalf("Failed to parse: %v", err)
	}
	
	resourceAnalyzer := analyzer.NewResourceAnalyzer()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		resourceAnalyzer.ExtractResources(spec)
	}
}

// BenchmarkPatternDetection benchmarks pattern detection
func BenchmarkPatternDetection(b *testing.B) {
	fixturePath := filepath.Join("fixtures", "stripe-openapi3.json")
	
	p := parser.New()
	spec, err := p.ParseFile(fixturePath)
	if err != nil {
		b.Fatalf("Failed to parse: %v", err)
	}
	
	detector := analyzer.NewPatternDetector()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		detector.DetectPatterns(spec)
	}
}

// BenchmarkMarkdownGeneration benchmarks markdown output generation
func BenchmarkMarkdownGeneration(b *testing.B) {
	fixturePath := filepath.Join("fixtures", "stripe-minimal.json")
	
	// Set up analysis
	p := parser.New()
	spec, _ := p.ParseFile(fixturePath)
	resourceAnalyzer := analyzer.NewResourceAnalyzer()
	resources := resourceAnalyzer.ExtractResources(spec)
	
	analysis := &models.APIAnalysis{
		Title:     "Benchmark Test",
		Resources: resources,
	}
	
	rep := reporter.New()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		rep.Generate(analysis, "markdown")
	}
}

// BenchmarkResourceFiltering benchmarks resource filtering
func BenchmarkResourceFiltering(b *testing.B) {
	fixturePath := filepath.Join("fixtures", "stripe-openapi3.json")
	
	p := parser.New()
	spec, err := p.ParseFile(fixturePath)
	if err != nil {
		b.Fatalf("Failed to parse: %v", err)
	}
	
	resourceAnalyzer := analyzer.NewResourceAnalyzer()
	resources := resourceAnalyzer.ExtractResources(spec)
	
	filter := &analyzer.ResourceFilter{
		Pattern: "^customer",
	}
	filterer := analyzer.NewResourceFilterer()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		filterer.FilterResources(resources, filter)
	}
}