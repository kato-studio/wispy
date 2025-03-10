package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

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
	// Static Middleware
	ctx.Echo.Static("/static", "./static")
	// ### Middleware:
	// Logging and Security
	ctx.Echo.Use(middleware.Logger())
	ctx.Echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:             middleware.DefaultSkipper,
		StackSize:           4 << 10, // 4 KB
		DisableStackAll:     false,
		DisablePrintStack:   true,
		LogLevel:            0,
		LogErrorFunc:        nil,
		DisableErrorHandler: false,
	}))
	// ### SSL:
	// Staging / Testing Config
	ctx.Echo.AutoTLSManager.Client = &acme.Client{
		DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
	ctx.Echo.AutoTLSManager.HostPolicy = autocert.HostWhitelist("www.kato.studio", "kato.studio", "www.odt.agency", "odt.agency")
	// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
	ctx.Echo.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
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

	ctx.Echo.GET("/-run-site-build", func(c echo.Context) error {
		engine.BuildSiteMap()
		return c.HTML(200, `<body>site rebuilt!</body>`)
	})

	// brackets syntax
	ctx.Echo.GET("*", func(c echo.Context) error {
		//!TEMP CODE!!
		engine.BuildSiteMap()
		//!----

		domain := c.Request().Host
		requestPath := c.Request().URL.Path
		filename := filepath.Base(requestPath)

		// Handle essential site files served from "/"
		fmt.Println("requestPath", requestPath)
		fmt.Println("== hasPrefix ", strings.HasPrefix(requestPath, "/public/"))
		if _, exists := ctx.ESSENTIAL_SERVE[filename]; exists {
			return c.File("sites/" + domain + "/public/essential" + filename)
		} else if strings.HasPrefix(requestPath, "/public/") {
			fmt.Println("return file?")
			return c.File(filepath.Join("sites", domain, requestPath))
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
	engine.Log.Fatal(ctx.Echo.StartAutoTLS(":443"))
}
