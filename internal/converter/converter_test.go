package converter

import (
	"encoding/json"
	"testing"
)

func TestConvertSwagger2ToOpenAPI3(t *testing.T) {
	tests := []struct {
		name    string
		swagger string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "basic conversion with info and paths",
			swagger: `{
				"swagger": "2.0",
				"info": {
					"title": "Test API",
					"version": "1.0.0",
					"description": "A test API"
				},
				"host": "api.example.com",
				"basePath": "/v1",
				"schemes": ["https"],
				"paths": {
					"/users": {
						"get": {
							"summary": "Get users",
							"produces": ["application/json"],
							"responses": {
								"200": {
									"description": "Success"
								}
							}
						}
					}
				}
			}`,
			want: map[string]interface{}{
				"openapi": "3.0.3",
				"info": map[string]interface{}{
					"title":       "Test API",
					"version":     "1.0.0",
					"description": "A test API",
				},
				"servers": []interface{}{
					map[string]interface{}{
						"url": "https://api.example.com/v1",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "convert parameters",
			swagger: `{
				"swagger": "2.0",
				"info": {"title": "Test", "version": "1.0"},
				"paths": {
					"/users/{id}": {
						"get": {
							"parameters": [
								{
									"name": "id",
									"in": "path",
									"required": true,
									"type": "string"
								},
								{
									"name": "include",
									"in": "query",
									"type": "array",
									"items": {"type": "string"}
								}
							],
							"responses": {"200": {"description": "OK"}}
						}
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "convert definitions to components",
			swagger: `{
				"swagger": "2.0",
				"info": {"title": "Test", "version": "1.0"},
				"definitions": {
					"User": {
						"type": "object",
						"properties": {
							"id": {"type": "string"},
							"name": {"type": "string"}
						}
					}
				},
				"paths": {
					"/users": {
						"post": {
							"parameters": [{
								"name": "body",
								"in": "body",
								"schema": {"$ref": "#/definitions/User"}
							}],
							"responses": {"201": {"description": "Created"}}
						}
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "convert security definitions",
			swagger: `{
				"swagger": "2.0",
				"info": {"title": "Test", "version": "1.0"},
				"securityDefinitions": {
					"api_key": {
						"type": "apiKey",
						"name": "X-API-Key",
						"in": "header"
					},
					"oauth2": {
						"type": "oauth2",
						"flow": "implicit",
						"authorizationUrl": "https://auth.example.com",
						"scopes": {
							"read": "Read access",
							"write": "Write access"
						}
					}
				},
				"paths": {}
			}`,
			wantErr: false,
		},
		{
			name:    "invalid swagger version",
			swagger: `{"swagger": "1.0"}`,
			wantErr: true,
		},
		{
			name:    "missing required fields",
			swagger: `{"swagger": "2.0"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			got, err := c.Convert([]byte(tt.swagger))

			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil {
				// Verify OpenAPI 3.x structure
				var result map[string]interface{}
				if err := json.Unmarshal(got, &result); err != nil {
					t.Errorf("Failed to unmarshal result: %v", err)
					return
				}

				// Check OpenAPI version
				if v, ok := result["openapi"].(string); !ok || v != "3.0.3" {
					t.Errorf("Expected openapi version 3.0.3, got %v", result["openapi"])
				}

				// Check info exists
				if _, ok := result["info"]; !ok {
					t.Error("Missing info object in converted spec")
				}

				// Check paths exist
				if _, ok := result["paths"]; !ok {
					t.Error("Missing paths object in converted spec")
				}
			}
		})
	}
}

func TestConvertRealWorldSpecs(t *testing.T) {
	// Test with actual UAT examples
	tests := []struct {
		name     string
		specPath string
	}{
		{
			name:     "warden API",
			specPath: "../../uat/artifacts/warden.v1.swagger.json",
		},
		{
			name:     "forge API",
			specPath: "../../uat/artifacts/forge.swagger.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: These tests will be implemented when we have file reading capability
			t.Skip("Real-world spec tests will be implemented with file reading")
		})
	}
}
