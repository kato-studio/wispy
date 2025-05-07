package handlers

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/kato-studio/wispy/engine"
)

// Define some example colors for logging
const (
	colorCyan  = "\033[36m"
	colorGrey  = "\033[90m"
	colorReset = "\033[0m"
)

func SitePublicFolderHnadler(w http.ResponseWriter, r *http.Request) {
	// Handle essential site files served from "/"
	domain := r.Host
	requestPath := r.URL.Path
	filename := filepath.Base(requestPath)
	if _, exists := engine.ESSENTIAL_SERVE[filename]; exists {
		targetFile := filepath.Join(engine.Wispy.SITE_DIR, domain, "public/essential", filename)
		// Serve the essential file
		http.ServeFile(w, r, targetFile)
		return
	}

	// Serve public assets
	// targetFile := filepath.Join(engine.Wispy.SITE_DIR, domain, requestPath)
	// http.StripPrefix("/public/", http.FileServer(http.Dir("static")))
	targetFile := filepath.Join(engine.Wispy.SITE_DIR, domain, requestPath)
	fmt.Println(targetFile)
	http.ServeFile(w, r, targetFile)
	// if _, err := os.Stat(targetFile); err != nil {
	// 	http.Error(w, "File not found", http.StatusNotFound)
	// 	return
	// }
	// // http.ServeFile(w, r, targetFile)
	// rawBytes, err := os.ReadFile(targetFile)
	// if err != nil {
	// 	http.Error(w, "File cloud not be read", http.StatusNotFound)
	// 	return
	// }
	// w.Write(rawBytes)
}

func SiteRouteHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	domain := r.Host
	//! ------------
	//! loacltesting addtions should never in prod!
	//! ------------
	//
	domain_cookie := r.CookiesNamed("local_dev_domain")
	if len(domain_cookie) == 1 {
		domain = domain_cookie[0].Value
		r.Host = domain
		fmt.Println("settng domain = ", domain)
	}
	requestPath := r.URL.Path
	//!---------------------------

	// Look up the site structure for the domain
	site, exists := engine.SiteMap[domain]
	if !exists {
		http.Error(w, fmt.Sprintf("domain %s not found", domain), http.StatusNotFound)
		return
	}

	// -----------
	data := map[string]any{}
	// -----------

	// Render the route
	page, err := engine.RenderRoute(site, requestPath, data, w, r)
	if err != nil {
		slog.Error("Rendering Route using \"RenderRoute()\""+err.Error(), nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Measure rendering and styling times
	renderTime := time.Now()
	var results bytes.Buffer
	results.WriteString(page)

	// Log performance metrics
	colorize := func(dur time.Duration) string {
		return fmt.Sprintf("%s%v%s", colorCyan, dur, colorGrey)
	}

	fmt.Printf("%s[Render: %s | Total: %s]%s\n",
		colorGrey,
		colorize(renderTime.Sub(startTime)),
		colorize(time.Since(startTime)),
		colorReset,
	)

	// Write the final HTML to the response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(results.Bytes())
}
