package template

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kato-studio/wispy/wispy_common/structure"
)

const (
	// colorReset = "\033[0m"
	// colorBlue    = "\033[34m"
	// colorGreen   = "\033[32m"
	// colorYellow  = "\033[33m"
	// colorMagenta = "\033[35m"
	// colorCyan    = "\033[36m"
	colorRed = "\033[31m"
	// colorGrey = "\033[90m"
)

// RenderRoute renders a page route for a given domain and page name.
// It looks up the page in the site's route map
// The route key is assumed to be in the form "domain/pageName" (e.g. "example.com/about").
func RenderRoute(engine *structure.TemplateEngine, ctx *structure.RenderCtx, requestPath string, data map[string]any, w http.ResponseWriter, r *http.Request) (output string, err error) {
	ctx.ResponseWriter = &w
	ctx.Request = r

	// Construct the route key. If route is empty, key becomes "domain/".
	site := ctx.Site
	routeKey := site.Domain + requestPath
	route, exists := site.Routes[routeKey]
	if !exists {
		http.Redirect(*ctx.ResponseWriter, ctx.Request, "/error?code=404&source="+requestPath+"&error=route not found", http.StatusSeeOther)
		return "", fmt.Errorf("route %s not found", routeKey)
	}

	// Create the render context and inject it into the data.
	if data == nil {
		data = make(map[string]any)
	}

	// Add easy access to URL Query Params
	// Initialize URL.Query structure if not exists
	if _, ok := data["URL"].(map[string]any); !ok {
		data["URL"] = make(map[string]any)
	}
	if _, ok := data["URL"].(map[string]any)["Query"].(map[string]any); !ok {
		data["URL"].(map[string]any)["Query"] = make(map[string]any)
	}
	queryMap := data["URL"].(map[string]any)["Query"].(map[string]any)
	// Add query parameters
	for key, values := range r.URL.Query() {
		if len(values) > 1 {
			queryMap[key] = values
		} else {
			queryMap[key] = values[0]
		}
	}

	// Read and merge JSON data if the file exist
	jsonPath := filepath.Join(filepath.Dir(route.Path), "data_en.json")
	jsonAsBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Println("data_en.json not found")
		log.Printf("failed to read JSON file: %s \n", err.Error())
	} else {
		var jsonData map[string]any
		if err := json.Unmarshal(jsonAsBytes, &jsonData); err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		// Merge JSON data with existing data
		for k, v := range jsonData {
			data[k] = v
		}
	}

	rootLayoutPath := path.Join(ctx.ScopedDirectory, "layouts", "root.hstm")
	rootLayoutAsBytes, err := os.ReadFile(rootLayoutPath)
	if err != nil {
		fmt.Println(err)
		slog.Error("Failed to read root layout at ", rootLayoutPath, ": ", err)
		return "", fmt.Errorf("Failed to read root layout at %s", rootLayoutPath)
	}

	templateAsBytes, err := os.ReadFile(route.Path)
	if err != nil {
		fmt.Println(err)
		slog.Error("Failed to read page template at ", route.Path, ": ", err)
		return "", fmt.Errorf("route %s not found", routeKey)
	}
	// Update for use in asset imports
	ctx.CurrentTemplatePath = strings.TrimSuffix(route.Path, ctx.Engine.PAGE_FILE_NAME)
	//
	var sb strings.Builder
	ctx.Data = data
	renderErrors := Render(ctx, &sb, string(templateAsBytes))

	// Render the layouts/root.hstm
	ctx.Passed = sb.String()
	var rootSb strings.Builder
	rootRenderErrs := Render(ctx, &rootSb, string(rootLayoutAsBytes))
	if len(rootRenderErrs) > 0 {
		renderErrors = append(renderErrors, rootRenderErrs...)
	}

	// TODO: better error logging using built-in go logger with better highlighting
	for ei, err := range renderErrors {
		if ei == 0 {
			fmt.Println(colorGrey + "-------------------" + colorReset)
		}
		fmt.Println(colorGrey+"["+colorRed+"Error"+colorGrey+"] "+colorReset, err)
		if ei == len(renderErrors)-1 {
			fmt.Println(colorGrey + "-------------------")
		}
	}

	return rootSb.String(), err
}

// SetupWispyCache ensures the .wispy cache directory exists.
func SetupWispyCache() {
	cacheDir := ".wispy"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Failed to create .wispy directory: %v", err))
		}
	}
}
