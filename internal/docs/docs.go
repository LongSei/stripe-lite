package docs

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/thep1ckaxe91/stripe-lite/internal/api/v1"
	"gopkg.in/yaml.v3"
)

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Stripe-lite API Gateway - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5/favicon-32x32.png" sizes="32x32" />
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
    .topbar { display: none !important; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" charset="UTF-8"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`

// RegisterRoutes mounts the embedded Swagger UI and OpenAPI raw spec endpoints.
func RegisterRoutes(r *gin.Engine) error {
	swagger, err := v1.GetSwagger()
	if err != nil {
		return fmt.Errorf("failed to retrieve swagger spec: %w", err)
	}

	specJSON, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal swagger to json: %w", err)
	}

	rawJSON, err := v1.RawSpecData()
	if err != nil {
		return fmt.Errorf("failed to retrieve raw spec: %w", err)
	}

	var rawObj any
	if err := json.Unmarshal(rawJSON, &rawObj); err != nil {
		return fmt.Errorf("failed to unmarshal spec json: %w", err)
	}

	rawYAML, err := jsonToYAML(rawObj)
	if err != nil {
		return fmt.Errorf("failed to marshal spec yaml: %w", err)
	}

	// Interactive Documentation UIs
	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs")
	})

	// Machine-readable OpenAPI spec endpoints
	r.GET("/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", specJSON)
	})
	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/yaml; charset=utf-8", rawYAML)
	})

	return nil
}

func jsonToYAML(obj any) ([]byte, error) {
	return yaml.Marshal(obj)
}
