package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenAPIValidation(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "valid OpenAPI 3.0 with proper validation",
			content: `{
				"openapi": "3.0.3",
				"info": {
					"title": "Test API",
					"version": "1.0.0"
				},
				"paths": {
					"/users": {
						"get": {
							"responses": {
								"200": {
									"description": "Success"
								}
							}
						}
					}
				}
			}`,
			wantErr: false,
		},
		{
			name: "OpenAPI with missing response description (handled gracefully)",
			content: `{
				"openapi": "3.0.3",
				"info": {
					"title": "Test API",
					"version": "1.0.0"
				},
				"paths": {
					"/users": {
						"get": {
							"responses": {
								"200": {}
							}
						}
					}
				}
			}`,
			wantErr: false, // go-openapi is more lenient
		},
		{
			name: "OpenAPI with missing paths (handled gracefully)",
			content: `{
				"openapi": "3.0.3",
				"info": {
					"title": "Test API",
					"version": "1.0.0"
				}
			}`,
			wantErr: false, // go-openapi is more lenient
		},
		{
			name: "valid OpenAPI with components",
			content: `{
				"openapi": "3.0.3",
				"info": {
					"title": "Test API",
					"version": "1.0.0"
				},
				"paths": {
					"/users": {
						"get": {
							"responses": {
								"200": {
									"description": "Success",
									"content": {
										"application/json": {
											"schema": {
												"$ref": "#/components/schemas/User"
											}
										}
									}
								}
							}
						}
					}
				},
				"components": {
					"schemas": {
						"User": {
							"type": "object",
							"properties": {
								"id": {
									"type": "string"
								},
								"name": {
									"type": "string"
								}
							}
						}
					}
				}
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpfile, err := os.CreateTemp("", "test-*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(tt.content)); err != nil {
				t.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			// Test enhanced parsing
			p := New()
			spec, err := p.ParseFile(tmpfile.Name())

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && spec != nil {
				// Verify parsing worked
				if spec.Info.Title != "Test API" {
					t.Errorf("Expected title 'Test API', got %s", spec.Info.Title)
				}

				// Verify paths were parsed (only for specs that actually have paths)
				if tt.name == "valid OpenAPI 3.0 with proper validation" || tt.name == "valid OpenAPI with components" || tt.name == "OpenAPI with missing response description (handled gracefully)" {
					if len(spec.Paths) == 0 {
						t.Error("Expected paths to be parsed")
					}
				}
			}
		})
	}
}

func TestOpenAPIWithLibraryValidation(t *testing.T) {
	// Get project root
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(wd, "..", "..")

	tests := []struct {
		name     string
		specFile string
		wantErr  bool
	}{
		{
			name:     "warden API with library validation",
			specFile: filepath.Join(projectRoot, "uat", "artifacts", "warden.v1.swagger.json"),
			wantErr:  false,
		},
		{
			name:     "forge API with library validation",
			specFile: filepath.Join(projectRoot, "uat", "artifacts", "forge.swagger.json"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			spec, err := p.ParseFile(tt.specFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && spec != nil {
				// Verify enhanced validation worked
				if spec.Info.Title == "" {
					t.Error("Expected title to be preserved after library validation")
				}

				t.Logf("Successfully validated %s with library: %s", tt.name, spec.Info.Title)
			}
		})
	}
}
