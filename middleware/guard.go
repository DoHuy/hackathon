package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BodySizeLimit(maxSizeMB int64) echo.MiddlewareFunc {
	limitBytes := maxSizeMB * 1024 * 1024
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().ContentLength > limitBytes {
				return c.String(http.StatusRequestEntityTooLarge, "Request too large")
			}
			return next(c)
		}
	}
}
