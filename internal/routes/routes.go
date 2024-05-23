// internal/routes/routes.go
package routes

import (
	"net/http"

	api "github.com/Prashansa-K/serviceCatalog/internal/api/v1"

	"github.com/labstack/echo/v4"
)

const (
	PONG_RESPONSE = "pong"
)

func RegisterRoutes(app *echo.Echo) {
	// Registering middlewares
	registerJaegarTracing(app)
	registerLogger(app)
	registerMetricsServer(app)
	registerTrailingSlashRemover(app)
	registerRateLimit(app)
	registerKeyBasedAuth(app)

	// Health check route
	app.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, PONG_RESPONSE)
	})

	appV1 := app.Group("/v1")

	appV1.GET("/services", api.GetServices)

	appV1.GET("/services/:serviceName", api.GetService)

	appV1.POST("/service", api.CreateService)

	appV1.POST("/service/version", api.CreateVersion)

	appV1.PATCH("/service", api.UpdateService)

	appV1.DELETE("/service/:serviceName", api.DeleteService)

	appV1.DELETE("/service/:serviceName/version/:versionName", api.DeleteVersion)
}
