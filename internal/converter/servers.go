package converter

import "fmt"

// convertServers creates OpenAPI 3.x servers from Swagger 2.0 host, basePath, and schemes
func (c *converter) convertServers(swagger map[string]interface{}) []map[string]interface{} {
	var servers []map[string]interface{}

	host, hasHost := swagger["host"].(string)
	basePath, hasBasePath := swagger["basePath"].(string)
	schemes, hasSchemes := swagger["schemes"].([]interface{})

	if !hasHost && !hasBasePath {
		// No server information available
		return servers
	}

	// Default values
	if !hasHost {
		host = "localhost"
	}
	if !hasBasePath {
		basePath = ""
	}
	if !hasSchemes || len(schemes) == 0 {
		schemes = []interface{}{"https"}
	}

	// Create a server entry for each scheme
	for _, scheme := range schemes {
		if schemeStr, ok := scheme.(string); ok {
			server := map[string]interface{}{
				"url": fmt.Sprintf("%s://%s%s", schemeStr, host, basePath),
			}

			// Add description if available
			if desc, ok := swagger["info"].(map[string]interface{})["x-server-description"].(string); ok {
				server["description"] = desc
			}

			servers = append(servers, server)
		}
	}

	return servers
}
