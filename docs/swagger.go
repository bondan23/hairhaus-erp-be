package docs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupSwagger registers Swagger UI and OpenAPI spec routes.
// - GET /swagger/doc.yaml   → serves the raw OpenAPI spec
// - GET /swagger/*any       → serves Swagger UI
func SetupSwagger(r *gin.Engine) {
	// Serve the OpenAPI YAML spec file
	r.StaticFile("/swagger/doc.yaml", "./docs/swagger.yaml")

	// Serve Swagger UI using embedded HTML pointing to our YAML
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	r.GET("/swagger/index.html", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>HAIRHAUS ERP API — Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: "/swagger/doc.yaml",
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.SwaggerUIStandalonePreset
            ],
            layout: "BaseLayout"
        });
    </script>
</body>
</html>`
