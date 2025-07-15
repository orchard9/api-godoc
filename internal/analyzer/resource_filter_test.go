package analyzer

import (
	"testing"

	"github.com/orchard9/pg-goapi/pkg/models"
)

func TestResourceFilter(t *testing.T) {
	// Sample resources for testing
	resources := []models.Resource{
		{Name: "User", Operations: []models.Operation{{Method: "GET", Path: "/users"}}},
		{Name: "UserProfile", Operations: []models.Operation{{Method: "GET", Path: "/users/{id}/profile"}}},
		{Name: "Order", Operations: []models.Operation{{Method: "GET", Path: "/orders"}}},
		{Name: "OrderItem", Operations: []models.Operation{{Method: "GET", Path: "/orders/{id}/items"}}},
		{Name: "Product", Operations: []models.Operation{{Method: "GET", Path: "/products"}}},
		{Name: "Category", Operations: []models.Operation{{Method: "GET", Path: "/categories"}}},
		{Name: "Payment", Operations: []models.Operation{{Method: "POST", Path: "/payments"}}},
		{Name: "PaymentMethod", Operations: []models.Operation{{Method: "GET", Path: "/payment-methods"}}},
	}

	tests := []struct {
		name      string
		resources []models.Resource
		filter    ResourceFilter
		want      []string // expected resource names
	}{
		{
			name:      "Include specific resources",
			resources: resources,
			filter: ResourceFilter{
				Include: []string{"User", "Order"},
			},
			want: []string{"User", "Order"},
		},
		{
			name:      "Exclude specific resources",
			resources: resources,
			filter: ResourceFilter{
				Exclude: []string{"Payment", "PaymentMethod"},
			},
			want: []string{"User", "UserProfile", "Order", "OrderItem", "Product", "Category"},
		},
		{
			name:      "Pattern matching with regex",
			resources: resources,
			filter: ResourceFilter{
				Pattern: "^User.*",
			},
			want: []string{"User", "UserProfile"},
		},
		{
			name:      "Pattern matching Order resources",
			resources: resources,
			filter: ResourceFilter{
				Pattern: "Order",
			},
			want: []string{"Order", "OrderItem"},
		},
		{
			name:      "Include overrides exclude",
			resources: resources,
			filter: ResourceFilter{
				Include: []string{"Payment"},
				Exclude: []string{"Payment", "User"},
			},
			want: []string{"Payment"},
		},
		{
			name:      "No filter returns all",
			resources: resources,
			filter:    ResourceFilter{},
			want:      []string{"User", "UserProfile", "Order", "OrderItem", "Product", "Category", "Payment", "PaymentMethod"},
		},
		{
			name:      "Case insensitive matching",
			resources: resources,
			filter: ResourceFilter{
				Include: []string{"user", "PRODUCT"},
			},
			want: []string{"User", "Product"},
		},
		{
			name:      "Pattern with include",
			resources: resources,
			filter: ResourceFilter{
				Pattern: "^User",
				Include: []string{"Product"},
			},
			want: []string{"User", "UserProfile", "Product"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filterer := NewResourceFilterer()
			filtered := filterer.FilterResources(tt.resources, &tt.filter)

			// Convert to names for easier comparison
			got := make([]string, len(filtered))
			for i, r := range filtered {
				got[i] = r.Name
			}

			// Check lengths
			if len(got) != len(tt.want) {
				t.Errorf("FilterResources() returned %d resources, want %d", len(got), len(tt.want))
				t.Errorf("Got: %v, Want: %v", got, tt.want)
				return
			}

			// Check each resource is present
			wantMap := make(map[string]bool)
			for _, name := range tt.want {
				wantMap[name] = true
			}

			for _, name := range got {
				if !wantMap[name] {
					t.Errorf("FilterResources() returned unexpected resource: %s", name)
				}
				delete(wantMap, name)
			}

			// Check all expected resources were found
			for name := range wantMap {
				t.Errorf("FilterResources() missing expected resource: %s", name)
			}
		})
	}
}

func TestResourceFilterWithOperations(t *testing.T) {
	// Test that filtered resources retain their operations
	resources := []models.Resource{
		{
			Name: "User",
			Operations: []models.Operation{
				{Method: "GET", Path: "/users"},
				{Method: "POST", Path: "/users"},
				{Method: "GET", Path: "/users/{id}"},
			},
		},
		{
			Name: "Order",
			Operations: []models.Operation{
				{Method: "GET", Path: "/orders"},
				{Method: "POST", Path: "/orders"},
			},
		},
	}

	filter := ResourceFilter{
		Include: []string{"User"},
	}

	filterer := NewResourceFilterer()
	filtered := filterer.FilterResources(resources, &filter)

	if len(filtered) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(filtered))
	}

	if filtered[0].Name != "User" {
		t.Errorf("Expected User resource, got %s", filtered[0].Name)
	}

	if len(filtered[0].Operations) != 3 {
		t.Errorf("Expected 3 operations, got %d", len(filtered[0].Operations))
	}
}

func TestResourceFilterValidation(t *testing.T) {
	tests := []struct {
		name    string
		filter  ResourceFilter
		wantErr bool
	}{
		{
			name:    "Valid include filter",
			filter:  ResourceFilter{Include: []string{"User"}},
			wantErr: false,
		},
		{
			name:    "Valid pattern filter",
			filter:  ResourceFilter{Pattern: "^User.*"},
			wantErr: false,
		},
		{
			name:    "Invalid regex pattern",
			filter:  ResourceFilter{Pattern: "["},
			wantErr: true,
		},
		{
			name:    "Empty filter is valid",
			filter:  ResourceFilter{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}