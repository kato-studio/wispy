package middleware

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"

	"github.com/kato-studio/wispy/engine"
)

// Middleware to restrict access to essential favicon files
func FaviconMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the file name from the request path
		requestPath := c.Request().URL.Path
		filename := filepath.Base(requestPath)

		// Check if the filename exists in the map
		if _, exists := engine.ESSENTIAL_SERVE[filename]; exists {
			// Proceed to serve the file
			return c.File("")
		}

		// If the file is not in the list, return 404
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "file not found",
		})
	}
}

func FileExistsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the file name from the request path
		requestPath := c.Request().URL.Path
		filename := filepath.Base(requestPath)

		// Check if the filename exists in the map
		if filename != "." {
			// Proceed to serve the file
			return next(c)
		}

		// If the file is not in the list, return 404
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "file not found",
		})
	}
}
