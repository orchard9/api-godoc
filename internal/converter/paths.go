package converter

// convertPaths converts Swagger 2.0 paths to OpenAPI 3.x format
func (c *converter) convertPaths(paths map[string]interface{}) map[string]interface{} {
	convertedPaths := make(map[string]interface{})

	for path, pathItem := range paths {
		if pathObj, ok := pathItem.(map[string]interface{}); ok {
			convertedPath := make(map[string]interface{})

			// Convert each operation (get, post, put, delete, etc.)
			for method, operation := range pathObj {
				if opObj, ok := operation.(map[string]interface{}); ok {
					convertedPath[method] = c.convertOperation(opObj)
				}
			}

			convertedPaths[path] = convertedPath
		}
	}

	return convertedPaths
}

// convertOperation converts a single operation from Swagger 2.0 to OpenAPI 3.x
func (c *converter) convertOperation(operation map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})

	// Copy simple fields
	simpleFields := []string{"summary", "description", "operationId", "tags", "deprecated", "security"}
	for _, field := range simpleFields {
		if value, ok := operation[field]; ok {
			converted[field] = value
		}
	}

	// Convert parameters
	if params, ok := operation["parameters"].([]interface{}); ok {
		convertedParams, requestBody := c.convertParameters(params)
		if len(convertedParams) > 0 {
			converted["parameters"] = convertedParams
		}
		if requestBody != nil {
			converted["requestBody"] = requestBody
		}
	}

	// Convert responses
	if responses, ok := operation["responses"].(map[string]interface{}); ok {
		converted["responses"] = c.convertResponses(responses)
	}

	// Convert produces to content types in responses
	if produces, ok := operation["produces"].([]interface{}); ok && len(produces) > 0 {
		// Store produces for use in response conversion
		converted["x-produces"] = produces
	}

	// Convert consumes to content types in requestBody
	if consumes, ok := operation["consumes"].([]interface{}); ok && len(consumes) > 0 {
		// Store consumes for use in requestBody conversion
		if rb, ok := converted["requestBody"].(map[string]interface{}); ok {
			c.addContentTypes(rb, consumes)
		}
	}

	return converted
}

// convertParameters converts Swagger 2.0 parameters to OpenAPI 3.x format
func (c *converter) convertParameters(params []interface{}) ([]interface{}, map[string]interface{}) {
	var parameters []interface{}
	var requestBody map[string]interface{}

	for _, param := range params {
		if p, ok := param.(map[string]interface{}); ok {
			if p["in"] == "body" {
				// Convert body parameter to requestBody
				requestBody = c.convertBodyParameter(p)
			} else {
				// Convert other parameters
				parameters = append(parameters, c.convertParameter(p))
			}
		}
	}

	return parameters, requestBody
}

// convertParameter converts a non-body parameter
func (c *converter) convertParameter(param map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})

	// Copy basic fields
	fields := []string{"name", "in", "description", "required", "deprecated", "allowEmptyValue"}
	for _, field := range fields {
		if value, ok := param[field]; ok {
			converted[field] = value
		}
	}

	// Convert schema
	schema := make(map[string]interface{})

	// Type and format
	if t, ok := param["type"]; ok {
		schema["type"] = t
		if f, ok := param["format"]; ok {
			schema["format"] = f
		}
	}

	// Array items
	if items, ok := param["items"]; ok {
		schema["items"] = items
	}

	// Other schema properties
	schemaProps := []string{"minimum", "maximum", "pattern", "enum", "default"}
	for _, prop := range schemaProps {
		if value, ok := param[prop]; ok {
			schema[prop] = value
		}
	}

	converted["schema"] = schema

	return converted
}

// convertBodyParameter converts a body parameter to requestBody
func (c *converter) convertBodyParameter(param map[string]interface{}) map[string]interface{} {
	requestBody := make(map[string]interface{})

	if desc, ok := param["description"]; ok {
		requestBody["description"] = desc
	}

	if req, ok := param["required"]; ok {
		requestBody["required"] = req
	}

	// Create content with application/json by default
	content := make(map[string]interface{})
	mediaType := make(map[string]interface{})

	if schema, ok := param["schema"]; ok {
		mediaType["schema"] = c.convertSchemaRef(schema)
	}

	content["application/json"] = mediaType
	requestBody["content"] = content

	return requestBody
}

// addContentTypes adds content types to a requestBody
func (c *converter) addContentTypes(requestBody map[string]interface{}, consumes []interface{}) {
	if content, ok := requestBody["content"].(map[string]interface{}); ok {
		// Get the schema from existing content
		var schema interface{}
		for _, mediaType := range content {
			if mt, ok := mediaType.(map[string]interface{}); ok {
				if s, ok := mt["schema"]; ok {
					schema = s
					break
				}
			}
		}

		// Clear and recreate content with all media types
		requestBody["content"] = make(map[string]interface{})
		content = requestBody["content"].(map[string]interface{})

		for _, contentType := range consumes {
			if ct, ok := contentType.(string); ok {
				content[ct] = map[string]interface{}{
					"schema": schema,
				}
			}
		}
	}
}
