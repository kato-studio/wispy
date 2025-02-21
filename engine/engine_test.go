package engine_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/kato-studio/wispy/atomicstyle"
	"github.com/kato-studio/wispy/engine"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *echo.Echo {
	e := echo.New()
	engine.Wispy.SITE_DIR = "../example/sites"
	engine.BuildSiteMap()

	e.GET("*", func(c echo.Context) error {
		domain := c.Request().Host
		requestPath := c.Request().URL.Path
		filename := filepath.Base(requestPath)

		fmt.Println("Handling request for domain:", domain, "Path:", requestPath)

		if _, exists := engine.ESSENTIAL_SERVE[filename]; exists {
			fmt.Println("Serving essential file:", filename)
			return c.File(engine.Wispy.SITE_DIR + "/" + domain + "/public/essential" + filename)
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

		results := bytes.NewBuffer([]byte{})
		reader := bytes.NewReader([]byte(page))

		trie := atomicstyle.BuildFullTrie()
		classes := atomicstyle.ExtractClasses(reader)
		css := atomicstyle.GenerateCSS(classes, trie)
		results.WriteString("<style>" + css + "</style>")
		results.Write(results.Bytes())

		fmt.Println("Response generated successfully:", time.Since(startTime))
		return c.HTMLBlob(http.StatusOK, results.Bytes())
	})

	return e
}

func TestServer(t *testing.T) {
	e := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "abc.test"
	res := httptest.NewRecorder()
	t.Log("Executing test request for abc.test")
	e.ServeHTTP(res, req)

	t.Log("=====OUTPUT=====")
	t.Logf(res.Body.String())
	t.Log("================")
	t.Logf("\n\n")
	t.Log("")

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Contains(t, res.Body.String(), "Welcome to abc.test")
	t.Log("Test completed successfully")
}
