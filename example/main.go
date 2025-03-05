package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/kato-studio/wispy/atomicstyle"
	"github.com/kato-studio/wispy/engine"
	"github.com/kato-studio/wispy/engine/ctx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	colorCyan  = "\033[36m"
	colorGrey  = "\033[90m"
	colorReset = "\033[0m"
)

func main() {
	const port = ":80"

	// Static Middleware
	ctx.Echo.Static("/public", "./public")
	// ### Middleware:
	// Logging and Security
	ctx.Echo.Use(middleware.Logger())
	ctx.Echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:             middleware.DefaultSkipper,
		StackSize:           4 << 10, // 4 KB
		DisableStackAll:     false,
		DisablePrintStack:   false,
		LogLevel:            0,
		LogErrorFunc:        nil,
		DisableErrorHandler: false,
	}))

	// Compression
	// app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	// 	Level: 5,
	// }))

	// Security
	// app.Use(middleware.Secure())
	// app.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(12))))
	// - Upload
	// app.Use(middleware.BodyLimit("2M"))
	// Routes

	trie := atomicstyle.BuildFullTrie()

	// brackets syntax
	ctx.Echo.GET("*", func(c echo.Context) error {
		domain := c.Request().Host
		requestPath := c.Request().URL.Path
		filename := filepath.Base(requestPath)

		// Handle essential site files served from "/"
		if _, exists := ctx.ESSENTIAL_SERVE[filename]; exists {
			return c.File("sites/" + domain + "/public/essential" + filename)
		}

		// Testing
		startTime := time.Now()
		data := map[string]any{
			"domain": c.Request().Host,
			"title":  "Welcome to " + c.Request().Host,
			"text":   "go to" + c.Request().Host,
			"link":   "/" + c.Request().Host,
			"name":   "john doe",
			"x":      42,
			"items":  []string{"apple", "banana", "cherry"},
			"bar":    "baz",
		}

		// Look up the site structure for the domain.
		site, exists := ctx.SiteMap[domain]
		if !exists {
			return fmt.Errorf("domain %s not found", domain)
		}

		page, err := engine.RenderRoute(&site, requestPath, data, c)
		if err != nil {
			engine.Log.Error(err)
			return err
		}

		// -----------
		renderTime := time.Now()

		var results = bytes.Buffer{}
		results.WriteString(page)

		styleTime := time.Now()
		reader := bytes.NewReader([]byte(page))

		styleExtTime := time.Now()
		classes := atomicstyle.ExtractClasses(reader)

		styleGenTime := time.Now()
		css := atomicstyle.GenerateCSS(classes, trie)
		// -----------

		colorize := func(dur time.Duration) string {
			return fmt.Sprintf("%s%v%s", colorCyan, dur, colorGrey)
		}

		fmt.Printf("%s[Render: %s | Style: %s | Extract: %s | CSS: %s | Total: %s]%s\n",
			colorGrey,
			colorize(renderTime.Sub(startTime)),
			colorize(styleTime.Sub(startTime)),
			colorize(styleExtTime.Sub(styleTime)),
			colorize(styleGenTime.Sub(styleExtTime)),
			colorize(time.Since(startTime)),
			colorReset)

		results.WriteString("<style>" + css + "</style>")
		//
		// return c.Render(http.StatusOK, templateName, data)
		return c.HTMLBlob(http.StatusOK, results.Bytes())
	})

	// Internal System setup
	// - Initial Build and Start
	// --- Build
	engine.BuildSiteMap()
	// --- Start
	engine.Log.Fatal(ctx.Echo.Start(port))
}
