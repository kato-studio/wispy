package template

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

// Define some example colors for logging
const (
	colorCyan  = "\033[36m"
	colorGrey  = "\033[90m"
	colorReset = "\033[0m"
)

func SitePublicFolderHandler(engine *structure.TemplateEngine, w http.ResponseWriter, r *http.Request) {
	// Handle essential site files served from "/"
	domain := r.Host
	requestPath := r.URL.Path
	filename := filepath.Base(requestPath)
	if _, exists := core.ESSENTIAL_SERVE[filename]; exists {
		targetFile := filepath.Join(engine.SITES_DIR, domain, "public/essential", filename)
		// Serve the essential file
		http.ServeFile(w, r, targetFile)
		return
	}

	// Serve public assets
	targetFile := filepath.Join(engine.SITES_DIR, domain, requestPath)
	fmt.Println(targetFile)
	http.ServeFile(w, r, targetFile)

}
func SiteAuthRouteHandler(engine *structure.TemplateEngine, w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	domain := r.Host

	// Look up the site structure for the domain
	site, exists := engine.SiteMap[domain]
	if !exists {
		http.Error(w, fmt.Sprintf("domain %s not found", domain), http.StatusNotFound)
		return
	}

	scopedDirectory := filepath.Join(engine.SITES_DIR, site.Domain)
	// Handle public content
	path := filepath.Clean(r.URL.Path)
	// if file extension check if there is a valid file in public directory to serve
	if filepath.Ext(path) != "" {
		// Serve public content if available
		root := os.DirFS(filepath.Join(scopedDirectory, "public"))
		f, err := root.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				// File doesn't exist, continue
			}
			fmt.Println("error opening public file:", path)
		} else {
			f.Close()
			stat, err := f.Stat()
			if err != nil && !stat.IsDir() {
				// File exists and is not a directory - serve it
				http.ServeFile(w, r, path)
				return
			}
		}
	}
	//
	data := map[string]any{}
	ctx := engine.InitCtx(scopedDirectory, &site, data)

	//
	page, err := RenderRoute(engine, ctx, r.URL.Path, data, w, r)
	if err != nil {
		slog.Error("Rendering Route using \"RenderRoute()\"" + err.Error())
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
