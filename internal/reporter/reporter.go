// Package reporter provides output formatting functionality
package reporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/orchard9/pg-goapi/pkg/models"
)

// Reporter defines the interface for generating output
type Reporter interface {
	// Generate creates formatted output from analysis results
	Generate(analysis *models.APIAnalysis, format string) (string, error)
}

// New creates a new reporter instance
func New() Reporter {
	return &reporter{}
}

type reporter struct{}

func (r *reporter) Generate(analysis *models.APIAnalysis, format string) (string, error) {
	switch format {
	case "json":
		return r.generateJSON(analysis)
	case "ai":
		return r.generateAIOptimized(analysis)
	default: // markdown
		return r.generateMarkdown(analysis)
	}
}

// generateMarkdown creates a comprehensive markdown report
func (r *reporter) generateMarkdown(analysis *models.APIAnalysis) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# %s\n\n", analysis.Title))

	if analysis.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", analysis.Description))
	}

	// Overview section
	sb.WriteString("## Overview\n\n")
	sb.WriteString(fmt.Sprintf("- **API Version**: %s\n", analysis.Version))
	sb.WriteString(fmt.Sprintf("- **Specification Type**: %s\n", analysis.SpecType))
	if analysis.BaseURL != "" {
		sb.WriteString(fmt.Sprintf("- **Base URL**: %s\n", analysis.BaseURL))
	}
	sb.WriteString(fmt.Sprintf("- **Generated**: %s\n\n", analysis.GeneratedAt.Format("2006-01-02 15:04:05")))

	// Statistics
	sb.WriteString("## API Statistics\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Resources**: %d\n", analysis.Summary.TotalResources))
	sb.WriteString(fmt.Sprintf("- **Total Operations**: %d\n", analysis.Summary.TotalOperations))
	sb.WriteString(fmt.Sprintf("- **Total Endpoints**: %d\n", analysis.Summary.TotalEndpoints))
	sb.WriteString(fmt.Sprintf("- **Resource Coverage**: %d%%\n\n", analysis.Summary.ResourceCoverage))

	// Resources section
	sb.WriteString("## Resources\n\n")
	sb.WriteString("This section groups API endpoints by business resources for better understanding.\n\n")

	// Sort resources by name for consistent output
	sortedResources := make([]models.Resource, len(analysis.Resources))
	copy(sortedResources, analysis.Resources)
	sort.Slice(sortedResources, func(i, j int) bool {
		return sortedResources[i].Name < sortedResources[j].Name
	})

	for _, resource := range sortedResources {
		r.writeResourceSection(&sb, resource)
	}

	// Relationships section
	if r.hasRelationships(analysis.Resources) {
		sb.WriteString("## Resource Relationships\n\n")
		sb.WriteString("This section shows how resources relate to each other.\n\n")

		// Add Mermaid diagram
		sb.WriteString("### Relationship Diagram\n\n")
		sb.WriteString("```mermaid\n")
		sb.WriteString(r.generateMermaidDiagram(analysis.Resources))
		sb.WriteString("```\n\n")

		sb.WriteString("### Relationship Details\n\n")
		r.writeRelationshipsSection(&sb, analysis.Resources)
	}

	// Patterns section
	if len(analysis.Patterns) > 0 {
		sb.WriteString("## Detected Patterns\n\n")
		for _, pattern := range analysis.Patterns {
			sb.WriteString(fmt.Sprintf("### %s\n\n", titleCase(pattern.Type)))
			sb.WriteString(fmt.Sprintf("**Confidence**: %s  \n", pattern.Confidence))
			sb.WriteString(fmt.Sprintf("**Impact**: %s\n\n", pattern.Impact))
			sb.WriteString(fmt.Sprintf("%s\n\n", pattern.Description))

			if len(pattern.Examples) > 0 {
				sb.WriteString("**Examples**:\n")
				for _, example := range pattern.Examples {
					sb.WriteString(fmt.Sprintf("- %s\n", example))
				}
				sb.WriteString("\n")
			}
		}
	}

	return sb.String(), nil
}

// writeResourceSection writes a resource and its operations
func (r *reporter) writeResourceSection(sb *strings.Builder, resource models.Resource) {
	sb.WriteString(fmt.Sprintf("### %s\n\n", titleCase(resource.Name)))

	if resource.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", resource.Description))
	}

	if resource.Category != "" {
		sb.WriteString(fmt.Sprintf("**Category**: %s  \n", resource.Category))
	}

	if resource.IsCollection {
		sb.WriteString("**Type**: Collection Resource  \n")
	}

	sb.WriteString(fmt.Sprintf("**Operations**: %d\n\n", len(resource.Operations)))

	// Sort operations by method and path
	sortedOps := make([]models.Operation, len(resource.Operations))
	copy(sortedOps, resource.Operations)
	sort.Slice(sortedOps, func(i, j int) bool {
		if sortedOps[i].Method != sortedOps[j].Method {
			return r.methodOrder(sortedOps[i].Method) < r.methodOrder(sortedOps[j].Method)
		}
		return sortedOps[i].Path < sortedOps[j].Path
	})

	// Write operations table
	sb.WriteString("| Method | Path | Summary |\n")
	sb.WriteString("|--------|------|----------|\n")

	for _, op := range sortedOps {
		summary := op.Summary
		if summary == "" {
			summary = op.Description
		}
		if len(summary) > 80 {
			summary = summary[:77] + "..."
		}

		// Escape pipe characters in summary
		summary = strings.ReplaceAll(summary, "|", "\\|")

		sb.WriteString(fmt.Sprintf("| %s | `%s` | %s |\n", op.Method, op.Path, summary))
	}

	sb.WriteString("\n")

	// Write detailed operation information
	for _, op := range sortedOps {
		r.writeOperationDetails(sb, op)
	}
}

// writeOperationDetails writes detailed operation information
func (r *reporter) writeOperationDetails(sb *strings.Builder, op models.Operation) {
	if op.Description == "" && len(op.Parameters) == 0 && len(op.Responses) == 0 {
		return // Skip if no additional details
	}

	sb.WriteString(fmt.Sprintf("#### %s %s\n\n", op.Method, op.Path))

	if op.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", op.Description))
	}

	if len(op.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("**Tags**: %s\n\n", strings.Join(op.Tags, ", ")))
	}

	if op.Deprecated {
		sb.WriteString("⚠️ **Deprecated**\n\n")
	}

	// Parameters
	if len(op.Parameters) > 0 {
		sb.WriteString("**Parameters**:\n\n")
		sb.WriteString("| Name | In | Type | Required | Description |\n")
		sb.WriteString("|------|----|----- |----------|-------------|\n")

		for _, param := range op.Parameters {
			required := "No"
			if param.Required {
				required = "Yes"
			}

			desc := strings.ReplaceAll(param.Description, "|", "\\|")
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}

			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				param.Name, param.In, param.Type, required, desc))
		}
		sb.WriteString("\n")
	}

	// Request Body
	if op.RequestBody != nil {
		sb.WriteString("**Request Body**:\n\n")
		if op.RequestBody.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", op.RequestBody.Description))
		}
		sb.WriteString(fmt.Sprintf("- **Content Type**: %s\n", op.RequestBody.ContentType))
		if op.RequestBody.Required {
			sb.WriteString("- **Required**: Yes\n")
		}
		sb.WriteString("\n")
	}

	// Responses
	if len(op.Responses) > 0 {
		sb.WriteString("**Responses**:\n\n")
		sb.WriteString("| Status | Description | Content Type |\n")
		sb.WriteString("|--------|-------------|---------------|\n")

		for _, resp := range op.Responses {
			desc := strings.ReplaceAll(resp.Description, "|", "\\|")
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}

			contentType := resp.ContentType
			if contentType == "" {
				contentType = "-"
			}

			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n",
				resp.StatusCode, desc, contentType))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
}

// writeRelationshipsSection writes resource relationships
func (r *reporter) writeRelationshipsSection(sb *strings.Builder, resources []models.Resource) {
	for _, resource := range resources {
		if len(resource.Relationships) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("### %s Relationships\n\n", titleCase(resource.Name)))

		for _, rel := range resource.Relationships {
			sb.WriteString(fmt.Sprintf("- **%s** %s (%s strength", rel.Type, rel.Resource, rel.Strength))
			if rel.Via != "" {
				sb.WriteString(fmt.Sprintf(" via `%s`", rel.Via))
			}
			sb.WriteString(")\n")

			if rel.Description != "" {
				sb.WriteString(fmt.Sprintf("  - %s\n", rel.Description))
			}
		}
		sb.WriteString("\n")
	}
}

// hasRelationships checks if any resources have relationships
func (r *reporter) hasRelationships(resources []models.Resource) bool {
	for _, resource := range resources {
		if len(resource.Relationships) > 0 {
			return true
		}
	}
	return false
}

// generateMermaidDiagram creates a Mermaid diagram for resource relationships
func (r *reporter) generateMermaidDiagram(resources []models.Resource) string {
	var sb strings.Builder
	sb.WriteString("graph TD\n")

	// First, list all resources as nodes
	for _, resource := range resources {
		sb.WriteString(fmt.Sprintf("    %s\n", resource.Name))
	}

	// Then, add relationships as edges
	for _, resource := range resources {
		for _, rel := range resource.Relationships {
			arrow := r.getMermaidArrow(rel.Strength)
			sb.WriteString(fmt.Sprintf("    %s %s|%s| %s\n",
				resource.Name, arrow, rel.Type, rel.Resource))
		}
	}

	return sb.String()
}

// getMermaidArrow returns the appropriate arrow style based on relationship strength
func (r *reporter) getMermaidArrow(strength string) string {
	switch strength {
	case "strong":
		return "-->"
	case "medium", "weak":
		return "-.->"
	default:
		return "-->"
	}
}

// methodOrder returns sort order for HTTP methods
func (r *reporter) methodOrder(method string) int {
	switch method {
	case "GET":
		return 0
	case "POST":
		return 1
	case "PUT":
		return 2
	case "PATCH":
		return 3
	case "DELETE":
		return 4
	default:
		return 5
	}
}

// generateJSON creates JSON output
func (r *reporter) generateJSON(analysis *models.APIAnalysis) (string, error) {
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

// generateAIOptimized creates AI-optimized condensed output
func (r *reporter) generateAIOptimized(analysis *models.APIAnalysis) (string, error) {
	var sb strings.Builder

	// Condensed header
	sb.WriteString(fmt.Sprintf("API: %s v%s (%s)\n", analysis.Title, analysis.Version, analysis.SpecType))
	sb.WriteString(fmt.Sprintf("Stats: %d resources, %d operations, %d endpoints\n\n",
		analysis.Summary.TotalResources, analysis.Summary.TotalOperations, analysis.Summary.TotalEndpoints))

	// Condensed resources
	sb.WriteString("RESOURCES:\n")
	for _, resource := range analysis.Resources {
		sb.WriteString(fmt.Sprintf("- %s (%d ops)", resource.Name, len(resource.Operations)))

		if len(resource.Relationships) > 0 {
			var relTypes []string
			for _, rel := range resource.Relationships {
				relTypes = append(relTypes, fmt.Sprintf("%s:%s", rel.Type, rel.Resource))
			}
			sb.WriteString(fmt.Sprintf(" -> %s", strings.Join(relTypes, ", ")))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\nKEY OPERATIONS:\n")
	for _, resource := range analysis.Resources {
		for _, op := range resource.Operations {
			if r.isKeyOperation(op) {
				sb.WriteString(fmt.Sprintf("- %s %s", op.Method, op.Path))
				if op.Summary != "" {
					summary := op.Summary
					if len(summary) > 60 {
						summary = summary[:57] + "..."
					}
					sb.WriteString(fmt.Sprintf(" (%s)", summary))
				}
				sb.WriteString("\n")
			}
		}
	}

	return sb.String(), nil
}

// isKeyOperation determines if an operation is important for AI summary
func (r *reporter) isKeyOperation(op models.Operation) bool {
	// Include CRUD operations and operations with clear business logic
	if op.Method == "GET" || op.Method == "POST" || op.Method == "PUT" || op.Method == "DELETE" {
		return true
	}

	// Include operations with meaningful summaries
	if op.Summary != "" && len(op.Summary) > 10 {
		return true
	}

	return false
}

// titleCase capitalizes the first letter of a string
func titleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
