package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kato-studio/wispy/engine/template"
	"github.com/labstack/echo/v4"
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
func (site *SiteStructure) RenderRoute(requestPath string, data map[string]interface{}, c echo.Context) (output string, err error) {
	// Construct the route key. If route is empty, key becomes "domain/".
	routeKey := site.Domain + requestPath
	fmt.Println("Looking for...", requestPath)
	fmt.Println("Looking for as ...", routeKey)
	route, exists := site.Routes[routeKey]
	if !exists {
		return "", fmt.Errorf("route %s not found", routeKey)
	}

	// Create the render context and inject it into the data.
	if data == nil {
		data = make(map[string]interface{})
	}

	// Optionally, inject additional values such as the page title.
	data["title"] = route.Title

	// Create a new template engine.
	engine := template.NewTemplateEngine()

	// Set up the rendering context using NewRenderCtx (which initializes Internal automatically).
	ctx := template.NewRenderCtx(map[string]interface{}{
		"Domain":      "example.com",
		"title":       "Welcome to abc.test!!",
		"showContent": "true", // any non-empty string except "false" is truthy
		"content":     "   This is some sample content.   ",
		"items":       "apple, banana, cherry",
		"condition":   "true",
	})

	// Render the sample template from file with timing logs.
	fmt.Println("\nRendered Template:")
	output, err = template.RenderString(route.Template, engine, ctx)

	return output, err
}

// NewSiteStructure creates a new SiteStructure with initialized maps.
func NewSiteStructure(domain string) SiteStructure {
	return SiteStructure{
		Domain:     domain,
		Routes:     make(map[string]PageRoutes),
		Layouts:    make(map[string]string),
		Components: make(map[string]string),
	}
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

			siteStructure := NewSiteStructure(domain)
			if _, err := toml.Decode(string(configBytes), &siteStructure); err != nil {
				fmt.Println(err)
				Log.Error("Failed to load config for ", domain, " at ", configFilePath, ": ", err)
			}

			// Build pages, layouts, and components paths.
			pagesPath := filepath.Join(siteFolderPath, "pages")
			layoutsPath := filepath.Join(siteFolderPath, "layouts")
			componentsPath := filepath.Join(siteFolderPath, "components")

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
						// For now, default the layout to default.html in layouts.
						defaultLayoutPath := filepath.Join(siteFolderPath, "layouts", "default"+Wispy.FILE_EXT)
						siteStructure.Routes[routeKey] = PageRoutes{
							Name:     pageName,
							Title:    domain,
							Layout:   defaultLayoutPath,
							Path:     path,
							Template: string(templateData),
							MetaTags: MetaTags{
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

			// Handle Components: walk through the components directory.
			filepath.WalkDir(componentsPath, func(path string, d fs.DirEntry, err error) error {
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
					siteStructure.Components[componentName] = string(templateData)
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

			SiteMap[domain] = siteStructure
		}
	}

	fmt.Println("SiteMap Build Time: ", time.Since(buildStart))
	// Log the list of sites for confirmation.
	var domains []string
	for domain := range SiteMap {
		domains = append(domains, domain)
	}
	fmt.Println("Sites: ", domains)
}
