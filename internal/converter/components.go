package converter

// convertComponents converts Swagger 2.0 definitions and other reusable components
func (c *converter) convertComponents(swagger map[string]interface{}) map[string]interface{} {
	components := make(map[string]interface{})

	// Convert definitions to schemas
	if definitions, ok := swagger["definitions"].(map[string]interface{}); ok {
		schemas := make(map[string]interface{})
		for name, definition := range definitions {
			schemas[name] = c.convertSchema(definition)
		}
		components["schemas"] = schemas
	}

	// Convert parameters
	if parameters, ok := swagger["parameters"].(map[string]interface{}); ok {
		convertedParams := make(map[string]interface{})
		for name, param := range parameters {
			if p, ok := param.(map[string]interface{}); ok {
				if p["in"] == "body" {
					// Body parameters become request bodies in OpenAPI 3.x
					// Store them separately for now
					continue
				}
				convertedParams[name] = c.convertParameter(p)
			}
		}
		if len(convertedParams) > 0 {
			components["parameters"] = convertedParams
		}
	}

	// Convert responses
	if responses, ok := swagger["responses"].(map[string]interface{}); ok {
		convertedResponses := make(map[string]interface{})
		for name, response := range responses {
			if r, ok := response.(map[string]interface{}); ok {
				convertedResponses[name] = c.convertResponse(r)
			}
		}
		if len(convertedResponses) > 0 {
			components["responses"] = convertedResponses
		}
	}

	// Convert securityDefinitions to securitySchemes
	if secDefs, ok := swagger["securityDefinitions"].(map[string]interface{}); ok {
		securitySchemes := make(map[string]interface{})
		for name, secDef := range secDefs {
			if sd, ok := secDef.(map[string]interface{}); ok {
				securitySchemes[name] = c.convertSecurityScheme(sd)
			}
		}
		components["securitySchemes"] = securitySchemes
	}

	return components
}

// convertSchema converts a schema definition
func (c *converter) convertSchema(schema interface{}) interface{} {
	switch s := schema.(type) {
	case map[string]interface{}:
		converted := make(map[string]interface{})

		// Handle $ref
		if ref, ok := s["$ref"].(string); ok {
			converted["$ref"] = c.convertRef(ref)
			return converted
		}

		// Copy schema properties
		schemaProps := []string{
			"type", "format", "title", "description", "default",
			"multipleOf", "maximum", "exclusiveMaximum", "minimum",
			"exclusiveMinimum", "maxLength", "minLength", "pattern",
			"maxItems", "minItems", "uniqueItems", "maxProperties",
			"minProperties", "required", "enum", "nullable",
			"discriminator", "readOnly", "writeOnly", "xml",
			"externalDocs", "example", "deprecated",
		}

		for _, prop := range schemaProps {
			if value, ok := s[prop]; ok {
				converted[prop] = value
			}
		}

		// Convert nested schemas
		if props, ok := s["properties"].(map[string]interface{}); ok {
			convertedProps := make(map[string]interface{})
			for name, prop := range props {
				convertedProps[name] = c.convertSchema(prop)
			}
			converted["properties"] = convertedProps
		}

		// Convert items
		if items, ok := s["items"]; ok {
			converted["items"] = c.convertSchema(items)
		}

		// Convert additionalProperties
		if addProps, ok := s["additionalProperties"]; ok {
			converted["additionalProperties"] = c.convertSchema(addProps)
		}

		// Convert allOf, anyOf, oneOf
		for _, key := range []string{"allOf", "anyOf", "oneOf"} {
			if schemas, ok := s[key].([]interface{}); ok {
				var convertedSchemas []interface{}
				for _, schema := range schemas {
					convertedSchemas = append(convertedSchemas, c.convertSchema(schema))
				}
				converted[key] = convertedSchemas
			}
		}

		return converted
	default:
		return s
	}
}

// convertSchemaRef handles schema references
func (c *converter) convertSchemaRef(schema interface{}) interface{} {
	if s, ok := schema.(map[string]interface{}); ok {
		if ref, ok := s["$ref"].(string); ok {
			return map[string]interface{}{
				"$ref": c.convertRef(ref),
			}
		}
	}
	return c.convertSchema(schema)
}

// convertRef converts Swagger 2.0 references to OpenAPI 3.x format
func (c *converter) convertRef(ref string) string {
	// Convert #/definitions/ to #/components/schemas/
	if len(ref) > 14 && ref[:14] == "#/definitions/" {
		return "#/components/schemas/" + ref[14:]
	}

	// Convert #/parameters/ to #/components/parameters/
	if len(ref) > 13 && ref[:13] == "#/parameters/" {
		return "#/components/parameters/" + ref[13:]
	}

	// Convert #/responses/ to #/components/responses/
	if len(ref) > 12 && ref[:12] == "#/responses/" {
		return "#/components/responses/" + ref[12:]
	}

	return ref
}

// convertSecurityScheme converts security definitions
func (c *converter) convertSecurityScheme(secDef map[string]interface{}) map[string]interface{} {
	converted := make(map[string]interface{})

	secType, _ := secDef["type"].(string)

	switch secType {
	case "basic":
		converted["type"] = "http"
		converted["scheme"] = "basic"

	case "apiKey":
		converted["type"] = "apiKey"
		if name, ok := secDef["name"]; ok {
			converted["name"] = name
		}
		if in, ok := secDef["in"]; ok {
			converted["in"] = in
		}

	case "oauth2":
		converted["type"] = "oauth2"
		flows := make(map[string]interface{})

		flow, _ := secDef["flow"].(string)
		switch flow {
		case "implicit":
			implicit := make(map[string]interface{})
			if authUrl, ok := secDef["authorizationUrl"]; ok {
				implicit["authorizationUrl"] = authUrl
			}
			if scopes, ok := secDef["scopes"]; ok {
				implicit["scopes"] = scopes
			}
			flows["implicit"] = implicit

		case "password":
			password := make(map[string]interface{})
			if tokenUrl, ok := secDef["tokenUrl"]; ok {
				password["tokenUrl"] = tokenUrl
			}
			if scopes, ok := secDef["scopes"]; ok {
				password["scopes"] = scopes
			}
			flows["password"] = password

		case "application":
			clientCreds := make(map[string]interface{})
			if tokenUrl, ok := secDef["tokenUrl"]; ok {
				clientCreds["tokenUrl"] = tokenUrl
			}
			if scopes, ok := secDef["scopes"]; ok {
				clientCreds["scopes"] = scopes
			}
			flows["clientCredentials"] = clientCreds

		case "accessCode":
			authCode := make(map[string]interface{})
			if authUrl, ok := secDef["authorizationUrl"]; ok {
				authCode["authorizationUrl"] = authUrl
			}
			if tokenUrl, ok := secDef["tokenUrl"]; ok {
				authCode["tokenUrl"] = tokenUrl
			}
			if scopes, ok := secDef["scopes"]; ok {
				authCode["scopes"] = scopes
			}
			flows["authorizationCode"] = authCode
		}

		converted["flows"] = flows
	}

	// Copy description
	if desc, ok := secDef["description"]; ok {
		converted["description"] = desc
	}

	return converted
}
