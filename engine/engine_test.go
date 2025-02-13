package engine_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/kato-studio/wispy/engine"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *echo.Echo {
	e := echo.New()
	engine.BuildSiteMap()

	e.GET("*", func(c echo.Context) error {
		domain := c.Request().Host
		requestPath := c.Request().URL.Path
		filename := filepath.Base(requestPath)

		fmt.Println("Handling request for domain:", domain, "Path:", requestPath)

		if _, exists := engine.ESSENTIAL_SERVE[filename]; exists {
			fmt.Println("Serving essential file:", filename)
			return c.File("sites/" + domain + "/public/essential" + filename)
		}

		startTime := time.Now()
		data := map[string]interface{}{
			"title": "Welcome to " + c.Request().Host,
		}

		fmt.Println("Checking site map for domain:", domain)
		site, exists := engine.SiteMap[domain]
		if !exists {
			fmt.Println("Domain not found in site map")
			return c.String(http.StatusNotFound, "Domain not found")
		}

		page, err := site.RenderRoute(requestPath, data, c)
		if err != nil {
			fmt.Println("Error rendering page:", err)
			return c.String(http.StatusInternalServerError, "Render error")
		}

		var results = bytes.Buffer{}
		results.WriteString(page)

		// classes := atomicstyle.RegexExtractClasses(page)
		// compiledCss := atomicstyle.WispyStyleGenerate(classes, atomicstyle.WispyStaticStyles, atomicstyle.WispyColors)
		// results.WriteString("<style>" + atomicstyle.WispyStyleCompile(compiledCss) + "</style>")

		fmt.Println("Response generated successfully:", time.Since(startTime))
		return c.HTMLBlob(http.StatusOK, results.Bytes())
	})

	return e
}

func TestServer(t *testing.T) {
	e := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/test-route", nil)
	req.Host = "abc.test"
	res := httptest.NewRecorder()
	fmt.Println("Executing test request for abc.test")
	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Welcome to abc.test")
	fmt.Println("Test completed successfully")
}
