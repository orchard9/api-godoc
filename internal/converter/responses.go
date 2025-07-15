package converter

// convertResponses converts Swagger 2.0 responses to OpenAPI 3.x format
func (c *converter) convertResponses(responses map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})

	for status, response := range responses {
		if respObj, ok := response.(map[string]interface{}); ok {
			converted[status] = c.convertResponse(respObj)
		}
	}

	return converted
}

// convertResponse converts a single response
func (c *converter) convertResponse(response map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})

	// Copy description (required in OpenAPI 3.x)
	if desc, ok := response["description"]; ok {
		converted["description"] = desc
	} else {
		converted["description"] = "Response"
	}

	// Convert headers
	if headers, ok := response["headers"].(map[string]interface{}); ok {
		converted["headers"] = c.convertHeaders(headers)
	}

	// Convert schema to content
	if schema, ok := response["schema"]; ok {
		content := make(map[string]interface{})
		mediaType := map[string]interface{}{
			"schema": c.convertSchemaRef(schema),
		}

		// Use default content type
		content["application/json"] = mediaType
		converted["content"] = content
	}

	// Convert examples if present
	if examples, ok := response["examples"]; ok {
		if content, ok := converted["content"].(map[string]interface{}); ok {
			for contentType, mediaType := range content {
				if mt, ok := mediaType.(map[string]interface{}); ok {
					if exampleForType, ok := examples.(map[string]interface{})[contentType]; ok {
						mt["example"] = exampleForType
					}
				}
			}
		}
	}

	return converted
}

// convertHeaders converts response headers
func (c *converter) convertHeaders(headers map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})

	for name, header := range headers {
		if h, ok := header.(map[string]interface{}); ok {
			convertedHeader := make(map[string]interface{})

			// Copy description
			if desc, ok := h["description"]; ok {
				convertedHeader["description"] = desc
			}

			// Convert type to schema
			schema := make(map[string]interface{})
			if t, ok := h["type"]; ok {
				schema["type"] = t
			}
			if f, ok := h["format"]; ok {
				schema["format"] = f
			}

			convertedHeader["schema"] = schema
			converted[name] = convertedHeader
		}
	}

	return converted
}
