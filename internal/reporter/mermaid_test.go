package reporter

import (
	"strings"
	"testing"

	"github.com/orchard9/api-godoc/pkg/models"
)

func TestGenerateMermaidDiagram(t *testing.T) {
	tests := []struct {
		name      string
		resources []models.Resource
		want      []string // Check for presence of these strings
	}{
		{
			name: "simple two resource relationship",
			resources: []models.Resource{
				{
					Name: "User",
					Relationships: []models.Relationship{
						{
							Type:        "has_many",
							Resource:    "Post",
							Description: "User has many posts",
							Strength:    "strong",
						},
					},
				},
				{
					Name: "Post",
					Relationships: []models.Relationship{
						{
							Type:        "belongs_to",
							Resource:    "User",
							Description: "Post belongs to user",
							Strength:    "strong",
						},
					},
				},
			},
			want: []string{
				"graph TD",
				"User",
				"Post",
				"User -->|has_many| Post",
				"Post -->|belongs_to| User",
			},
		},
		{
			name: "complex relationships with multiple resources",
			resources: []models.Resource{
				{
					Name: "Organization",
					Relationships: []models.Relationship{
						{
							Type:     "has_many",
							Resource: "Team",
							Strength: "strong",
						},
						{
							Type:     "has_many",
							Resource: "Project",
							Strength: "medium",
						},
					},
				},
				{
					Name: "Team",
					Relationships: []models.Relationship{
						{
							Type:     "belongs_to",
							Resource: "Organization",
							Strength: "strong",
						},
						{
							Type:     "has_many",
							Resource: "User",
							Strength: "strong",
						},
					},
				},
				{
					Name: "User",
					Relationships: []models.Relationship{
						{
							Type:     "belongs_to",
							Resource: "Team",
							Strength: "strong",
						},
						{
							Type:     "references",
							Resource: "Project",
							Strength: "weak",
						},
					},
				},
				{
					Name: "Project",
					Relationships: []models.Relationship{
						{
							Type:     "belongs_to",
							Resource: "Organization",
							Strength: "medium",
						},
						{
							Type:     "referenced_by",
							Resource: "User",
							Strength: "weak",
						},
					},
				},
			},
			want: []string{
				"graph TD",
				"Organization",
				"Team",
				"User",
				"Project",
				"Organization -->|has_many| Team",
				"Organization -.->|has_many| Project",
				"Team -->|belongs_to| Organization",
				"Team -->|has_many| User",
				"User -->|belongs_to| Team",
				"User -.->|references| Project",
				"Project -.->|belongs_to| Organization",
				"Project -.->|referenced_by| User",
			},
		},
		{
			name:      "no relationships",
			resources: []models.Resource{{Name: "User"}, {Name: "Post"}},
			want: []string{
				"graph TD",
				"User",
				"Post",
			},
		},
	}

	r := &reporter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.generateMermaidDiagram(tt.resources)
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("generateMermaidDiagram() missing expected content:\nwant: %s\ngot:\n%s", want, got)
				}
			}
		})
	}
}

func TestMermaidRelationshipArrow(t *testing.T) {
	tests := []struct {
		strength string
		want     string
	}{
		{"strong", "-->"},
		{"medium", "-.->"},
		{"weak", "-.->"},
		{"", "-->"},
	}

	r := &reporter{}
	for _, tt := range tests {
		t.Run(tt.strength, func(t *testing.T) {
			got := r.getMermaidArrow(tt.strength)
			if got != tt.want {
				t.Errorf("getMermaidArrow(%s) = %s, want %s", tt.strength, got, tt.want)
			}
		})
	}
}
