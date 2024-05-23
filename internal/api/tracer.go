package api

import (
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
)

func TraceServiceFunc(c echo.Context, fn interface{}, params ...interface{}) interface{} {
	return jaegertracing.TraceFunction(c, fn, params...)
}
