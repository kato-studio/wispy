package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/kato-studio/wispy/atomicstyle"
	"github.com/kato-studio/wispy/engine"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	const port = ":80"

	// Static Middleware
	engine.Echo.Static("/public", "./public")

	// ### Middleware:
	// Logging and Security
	engine.Echo.Use(middleware.Logger())
	engine.Echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
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

	engine.Echo.GET("*", func(c echo.Context) error {
		domain := c.Request().Host
		requestPath := c.Request().URL.Path
		filename := filepath.Base(requestPath)

		// Handle essential site files served from "/"
		if _, exists := engine.ESSENTIAL_SERVE[filename]; exists {
			return c.File("sites/" + domain + "/public/essential" + filename)
		}

		// Testing
		startTime := time.Now()
		data := map[string]interface{}{
			"title": "Welcome to " + c.Request().Host,
		}

		// Look up the site structure for the domain.
		site, exists := engine.SiteMap[domain]
		if !exists {
			return fmt.Errorf("domain %s not found", domain)
		}

		page, err := site.RenderRoute(requestPath, data, c)
		if err != nil {
			engine.Log.Error(err)
			return err
		}

		var results = bytes.Buffer{}
		results.WriteString(page)

		fmt.Println("----------")
		fmt.Println("[] Render Time", time.Since(startTime))

		classes := atomicstyle.RegexExtractClasses(page)
		compiledCss := atomicstyle.WispyStyleGenerate(classes, atomicstyle.WispyStaticStyles, atomicstyle.WispyColors)
		styleTime := time.Now()
		results.WriteString("<style>")
		results.WriteString(atomicstyle.WispyStyleCompile(compiledCss))
		results.WriteString("</style>")
		fmt.Println("[] Style Time", time.Since(styleTime))
		fmt.Println("[] Total Time", time.Since(startTime))
		fmt.Println("----------")

		// return c.Render(http.StatusOK, templateName, data)
		return c.HTMLBlob(http.StatusOK, results.Bytes())
	})

	// Internal System setup
	engine.BuildSiteMap()
	fmt.Println("-----")
	fmt.Print(len(engine.SiteMap))

	// Start server
	engine.Log.Fatal(engine.Echo.Start(port))

}
