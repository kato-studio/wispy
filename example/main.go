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
	var app = echo.New()

	// Static Middleware
	app.Static("/public", "./public")

	// ### Middleware:
	// Logging and Security
	app.Use(middleware.Logger())
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
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

	var eng = engine.StartEngine(engine.WispyConfig{}, app.Logger)

	eng.BuildSiteMap()

	// Routes
	app.GET("*", func(c echo.Context) error {
		host := c.Request().Host
		requestPath := c.Request().URL.Path
		filename := filepath.Base(requestPath)

		// Handle essential site files served from "/"
		if _, exists := engine.ESSENTIAL_SERVE[filename]; exists {
			return c.File("sites/" + host + "/public/essential" + filename)
		}

		// Testing
		startTime := time.Now()
		data := map[string]interface{}{
			"title": "Welcome to " + c.Request().Host,
		}
		results := bytes.NewBuffer([]byte{})
		err := eng.RenderRoute(results, host, data, c)
		if err != nil {
			return err
		}

		fmt.Println("----------")
		fmt.Println("[] Render Time", time.Since(startTime))

		classes := atomicstyle.RegexExtractClasses(results.String())
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

	// Start server
	app.Logger.Fatal(app.Start(port))
}
