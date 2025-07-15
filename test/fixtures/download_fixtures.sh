#!/bin/bash
# Download real-world OpenAPI specifications for testing

set -e

echo "Downloading real-world API fixtures..."

# Stripe API (OpenAPI 3.0)
echo "Downloading Stripe API spec..."
curl -L -o stripe-openapi3.json \
  "https://raw.githubusercontent.com/stripe/openapi/master/openapi/spec3.json" || \
  echo "Failed to download Stripe API spec"

# GitHub API (OpenAPI 3.0) 
echo "Downloading GitHub API spec..."
curl -L -o github-openapi3.json \
  "https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json" || \
  echo "Failed to download GitHub API spec"

# Kubernetes API (OpenAPI 3.0)
echo "Downloading Kubernetes API spec..."
# Note: K8s spec is very large (>20MB), so we'll create a simplified version
cat > kubernetes-simplified.json << 'EOF'
{
  "openapi": "3.0.0",
  "info": {
    "title": "Kubernetes API (Simplified)",
    "version": "v1.28.0",
    "description": "Simplified Kubernetes API for testing"
  },
  "servers": [
    {"url": "https://kubernetes.default.svc"}
  ],
  "paths": {
    "/api/v1/namespaces": {
      "get": {
        "summary": "List namespaces",
        "operationId": "listNamespaces",
        "responses": {
          "200": {"description": "OK"}
        }
      },
      "post": {
        "summary": "Create namespace",
        "operationId": "createNamespace",
        "responses": {
          "201": {"description": "Created"}
        }
      }
    },
    "/api/v1/namespaces/{namespace}/pods": {
      "get": {
        "summary": "List pods",
        "operationId": "listPods",
        "parameters": [
          {"name": "namespace", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {
          "200": {"description": "OK"}
        }
      }
    },
    "/api/v1/namespaces/{namespace}/services": {
      "get": {
        "summary": "List services",
        "operationId": "listServices",
        "parameters": [
          {"name": "namespace", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {
          "200": {"description": "OK"}
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Namespace": {
        "type": "object",
        "required": ["metadata"],
        "properties": {
          "metadata": {"$ref": "#/components/schemas/ObjectMeta"},
          "spec": {"type": "object"},
          "status": {"type": "object"}
        }
      },
      "Pod": {
        "type": "object",
        "required": ["metadata", "spec"],
        "properties": {
          "metadata": {"$ref": "#/components/schemas/ObjectMeta"},
          "spec": {"$ref": "#/components/schemas/PodSpec"},
          "status": {"type": "object"}
        }
      },
      "ObjectMeta": {
        "type": "object",
        "properties": {
          "name": {"type": "string"},
          "namespace": {"type": "string"},
          "labels": {"type": "object"},
          "annotations": {"type": "object"}
        }
      },
      "PodSpec": {
        "type": "object",
        "required": ["containers"],
        "properties": {
          "containers": {
            "type": "array",
            "items": {"$ref": "#/components/schemas/Container"}
          }
        }
      },
      "Container": {
        "type": "object",
        "required": ["name", "image"],
        "properties": {
          "name": {"type": "string"},
          "image": {"type": "string"},
          "ports": {
            "type": "array",
            "items": {"type": "object"}
          }
        }
      }
    }
  }
}
EOF

echo "Creating smaller test fixtures..."

# Create a minimal Stripe-like API for faster tests
cat > stripe-minimal.json << 'EOF'
{
  "openapi": "3.0.0",
  "info": {
    "title": "Stripe API (Minimal)",
    "version": "2024.01.01",
    "description": "Minimal Stripe-like API for testing"
  },
  "servers": [
    {"url": "https://api.stripe.com/v1"}
  ],
  "paths": {
    "/customers": {
      "get": {
        "summary": "List customers",
        "parameters": [
          {"name": "limit", "in": "query", "schema": {"type": "integer"}},
          {"name": "starting_after", "in": "query", "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "OK"}}
      },
      "post": {
        "summary": "Create customer",
        "responses": {"200": {"description": "OK"}}
      }
    },
    "/customers/{id}": {
      "get": {
        "summary": "Retrieve customer",
        "parameters": [
          {"name": "id", "in": "path", "required": true, "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "OK"}}
      }
    },
    "/charges": {
      "post": {
        "summary": "Create charge",
        "requestBody": {
          "content": {
            "application/x-www-form-urlencoded": {
              "schema": {"$ref": "#/components/schemas/ChargeRequest"}
            }
          }
        },
        "responses": {"200": {"description": "OK"}}
      }
    },
    "/subscriptions": {
      "get": {
        "summary": "List subscriptions",
        "parameters": [
          {"name": "customer", "in": "query", "schema": {"type": "string"}},
          {"name": "status", "in": "query", "schema": {"type": "string"}}
        ],
        "responses": {"200": {"description": "OK"}}
      }
    }
  },
  "components": {
    "schemas": {
      "Customer": {
        "type": "object",
        "properties": {
          "id": {"type": "string"},
          "email": {"type": "string", "format": "email"},
          "name": {"type": "string"},
          "created": {"type": "integer"}
        }
      },
      "ChargeRequest": {
        "type": "object",
        "required": ["amount", "currency"],
        "properties": {
          "amount": {"type": "integer"},
          "currency": {"type": "string"},
          "customer": {"type": "string"},
          "description": {"type": "string"}
        }
      }
    }
  }
}
EOF

echo "Done! Created fixtures:"
ls -la *.json