package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Prashansa-K/serviceCatalog/config"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// Rate limit
	RPS            = 5
	BURST_REQUESTS = 10

	// Auth
	BEARER_SCHEME = "Bearer"
	AUTH_HEADER   = "header:Authorization"
)

func registerTrailingSlashRemover(app *echo.Echo) {
	// Removing trailing slash if any to ensure that api responds even if a / is present in the end
	app.Pre(middleware.RemoveTrailingSlash())
}

func registerRateLimit(app *echo.Echo) {
	config := middleware.RateLimiterConfig{
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      RPS,
				Burst:     BURST_REQUESTS,
				ExpiresIn: 1 * time.Minute,
			},
		),

		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, "rate limit exceeded")
		},
	}

	app.Use(middleware.RateLimiterWithConfig(config))
}

func registerKeyBasedAuth(app *echo.Echo) {
	app.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  AUTH_HEADER,
		AuthScheme: BEARER_SCHEME,
		// require Authorization: Bearer header to be set
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("API_AUTH_KEY"), nil
		},

		ErrorHandler: func(err error, context echo.Context) error {
			return context.JSON(http.StatusUnauthorized, "invalid key")
		},
	}))
}

func registerLogger(app *echo.Echo) {
	logFile, err := os.OpenFile(".log/log_file", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("error opening log file: %v", err)
		logFile = os.Stdout
	}

	loggerConfig := middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		Output: logFile,
	}
	app.Use(middleware.LoggerWithConfig(loggerConfig))
}

func registerMetricsServer(app *echo.Echo) {
	app.Use(echoprometheus.NewMiddleware("serviceCatalog"))

	go func() {
		// running a separate metrics server
		metricServer := echo.New()

		metricsServerConfig := config.GetMetricsServerConfig()

		// adding a separate route to serve gathered metrics
		metricServer.GET("/metrics", echoprometheus.NewHandler())
		if err := metricServer.Start(metricsServerConfig.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
}

func registerJaegarTracing(app *echo.Echo) {
	jaegertracing.New(app, nil)
}
