package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kato-studio/wispy/engine/ctx"
	"github.com/kato-studio/wispy/engine/template"
	"github.com/labstack/echo/v4"
)

var Log = ctx.Log
var Wispy = ctx.Wispy

const (
	colorReset = "\033[0m"
	// colorBlue    = "\033[34m"
	// colorGreen   = "\033[32m"
	// colorYellow  = "\033[33m"
	// colorMagenta = "\033[35m"
	// colorCyan    = "\033[36m"
	colorRed  = "\033[31m"
	colorGrey = "\033[90m"
)

/*
=================================================================
Core External Functions
=================================================================
*/

// RenderRoute renders a page route for a given domain and page name.
// It looks up the page in the site's route map, executes the page template,
// and then wraps it in a layout if specified.
// The data parameter can include additional dynamic values and is augmented with
// the render context under the key "_ctx".
// The route key is assumed to be in the form "domain/pageName" (e.g. "example.com/about").
func RenderRoute(site *ctx.SiteStructure, requestPath string, data map[string]any, c echo.Context) (output string, err error) {
	// Construct the route key. If route is empty, key becomes "domain/".
	routeKey := site.Domain + requestPath
	fmt.Println("Looking for \"" + requestPath + "\" as \"./sites/" + routeKey + "\"")
	route, exists := site.Routes[routeKey]
	if !exists {
		return "", fmt.Errorf("route %s not found", routeKey)
	}

	// Create the render context and inject it into the data.
	if data == nil {
		data = make(map[string]any)
	}

	// Optionally, inject additional values such as the page title.
	data["title"] = route.Title

	// Create a new template engine.
	engine := template.NewTemplateEngine()

	// Set up the rendering context using NewRenderCtx (which initializes Internal automatically).
	ctx := template.NewRenderCtx(engine, map[string]any{
		"title":       "Welcome to abc.test!!",
		"showContent": "true", // any non-empty string except "false" is truthy
		"content":     "   This is some sample content.   ",
		"items":       []string{"kwei", "apple", "banana"},
		"isTrue":      "true",
		"condition":   "10 > 5",
	}, &site.Partials)

	var sb strings.Builder
	renderErrors := template.Render(ctx, &sb, route.Template)
	for ei, err := range renderErrors {
		if ei == 0 {
			fmt.Println(colorGrey + "-------------------" + colorReset)
		}
		fmt.Println(colorGrey+"["+colorRed+"Error"+colorGrey+"] "+colorReset, err)
		if ei == len(renderErrors)-1 {
			fmt.Println(colorGrey + "-------------------")
		}
	}

	return sb.String(), err
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

// BuildSiteMap builds the host-to-site mapping by reading directories from the sites folder.
func BuildSiteMap() {
	buildStart := time.Now()
	// Read the sites directory.
	entries, err := os.ReadDir(Wispy.SITE_DIR)
	if err != nil {
		panic(fmt.Sprintf("Failed to read sites directory: %v", err))
	}

	// Process each site (directory)
	for _, entry := range entries {
		if entry.IsDir() {
			domain := entry.Name()
			siteFolderPath := filepath.Join(Wispy.SITE_DIR, domain)
			configFilePath := filepath.Join(siteFolderPath, Wispy.SITE_CONFIG_NAME)

			// Read and decode the site config.
			configBytes, err := os.ReadFile(configFilePath)
			if err != nil {
				fmt.Println(err)
				Log.Error("Could not find config for ", domain, " at ", configFilePath, ": ", err)
				continue
			}

			siteStructure := ctx.NewSiteStructure(domain)
			if _, err := toml.Decode(string(configBytes), &siteStructure); err != nil {
				fmt.Println(err)
				Log.Error("Failed to load config for ", domain, " at ", configFilePath, ": ", err)
			}

			// Build pages, layouts, and partials paths.
			pagesPath := filepath.Join(siteFolderPath, "pages")
			layoutsPath := filepath.Join(siteFolderPath, "layouts")
			partialsPath := filepath.Join(siteFolderPath, "partials")

			// Handle Pages: walk through the pages directory.
			filepath.Walk(pagesPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					fmt.Println(err)
					Log.Error("Error accessing path ", path, ": ", err)
					return err
				}
				// Only process files with the configured extension.
				if !info.IsDir() && filepath.Ext(path) == Wispy.FILE_EXT {
					// Check if file name (without extension) matches the page file name.
					baseName := strings.TrimSuffix(filepath.Base(path), Wispy.FILE_EXT)
					if baseName == Wispy.PAGE_FILE_NAME {
						// Determine the page name as the relative directory from the pages folder.
						relDir, err := filepath.Rel(pagesPath, filepath.Dir(path))
						if err != nil {
							fmt.Println(err)
							Log.Error("Error computing relative path for ", path, ": ", err)
							return err
						}
						pageName := relDir
						if pageName == "." {
							pageName = ""
						}
						templateData, err := os.ReadFile(path)
						if err != nil {
							fmt.Println(err)
							Log.Error("Failed to read page template at ", path, ": ", err)
							return err
						}
						// Use a key combining the domain and the pageName.
						routeKey := domain + "/" + pageName
						fmt.Println("Saving " + routeKey)
						siteStructure.Routes[routeKey] = ctx.PageRoutes{
							Name:     pageName,
							Title:    domain,
							Layout:   "",
							Path:     path,
							Template: string(templateData),
							MetaTags: ctx.MetaTags{
								Title:         domain + " title",
								Description:   "Page description here",
								OgTitle:       domain + " title",
								OgDescription: "Page description here",
								OgType:        "text",
								OgUrl:         domain,
							},
						}
					}
				}
				return nil
			})

			// Handle Partials: walk through the partials directory.
			filepath.WalkDir(partialsPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					Log.Error("Error accessing component path ", path, ": ", err)
					return err
				}
				if !d.IsDir() && filepath.Ext(path) == Wispy.FILE_EXT {
					templateData, err := os.ReadFile(path)
					if err != nil {
						Log.Error("Failed to read component file at ", path, ": ", err)
						return err
					}
					componentName := strings.TrimSuffix(filepath.Base(path), Wispy.FILE_EXT)
					siteStructure.Partials[componentName] = string(templateData)
				}
				return nil
			})

			// Handle Layouts: walk through the layouts directory.
			filepath.Walk(layoutsPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					Log.Error("Error accessing layout path ", path, ": ", err)
					return err
				}
				if !info.IsDir() && filepath.Ext(path) == Wispy.FILE_EXT {
					templateData, err := os.ReadFile(path)
					if err != nil {
						Log.Error("Failed to read layout file at ", path, ": ", err)
						return err
					}
					layoutName := strings.TrimSuffix(filepath.Base(path), Wispy.FILE_EXT)
					siteStructure.Layouts[layoutName] = string(templateData)
				}
				return nil
			})

			ctx.SiteMap[domain] = siteStructure
		}
	}

	fmt.Println("SiteMap Build Time: ", time.Since(buildStart))
	// Log the list of sites for confirmation.
	var domains []string
	for domain := range ctx.SiteMap {
		domains = append(domains, domain)
	}
	fmt.Println("Sites: ", domains)
}
