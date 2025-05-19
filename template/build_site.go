package template

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/kato-studio/wispy/wispy_common/structure"
)

// NewSiteStructure creates a new SiteStructure with initialized maps.
func NewSiteStructure(domain string) structure.SiteStructure {
	return structure.SiteStructure{
		Domain:   domain,
		Routes:   make(map[string]structure.PageRoutes),
		Layouts:  make(map[string]string),
		Partials: make(map[string]string),
	}
}

// BuildSiteMap builds the host-to-site mapping by reading directories from the sites folder.
func BuildSiteMap(engine *structure.TemplateEngine) {
	buildStart := time.Now()
	// Read the sites directory.
	entries, err := os.ReadDir(engine.SITES_DIR)
	if err != nil {
		panic(fmt.Sprintf("Failed to read sites directory: %v", err))
	}

	// Process each site (directory)
	for _, entry := range entries {
		if entry.IsDir() {
			domain := entry.Name()
			siteFolderPath := filepath.Join(engine.SITES_DIR, domain)
			configFilePath := filepath.Join(siteFolderPath, engine.SITE_CONFIG_NAME)

			// Read and decode the site config.
			configBytes, err := os.ReadFile(configFilePath)
			if err != nil {
				fmt.Println(err)
				slog.Error("Could not find config for ", domain, " at ", configFilePath, ": ", err)
				continue
			}

			siteStructure := NewSiteStructure(domain)
			if _, err := toml.Decode(string(configBytes), &siteStructure); err != nil {
				fmt.Println(err)
				slog.Error("Failed to load config for ", domain, " at ", configFilePath, ": ", err)
			}

			// Build pages, layouts, and partials paths.
			pagesPath := filepath.Join(siteFolderPath, "pages")
			layoutsPath := filepath.Join(siteFolderPath, "layouts")
			partialsPath := filepath.Join(siteFolderPath, "partials")

			// Handle Pages: walk through the pages directory.
			filepath.Walk(pagesPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					fmt.Println(err)
					slog.Error("Error accessing path ", path, ": ", err)
					return err
				}
				// Only process files with the configured extension.
				if !info.IsDir() && filepath.Ext(path) == engine.FILE_EXT {
					// Check if file name (without extension) matches the page file name.
					baseName := strings.TrimSuffix(filepath.Base(path), engine.FILE_EXT)
					if baseName == engine.PAGE_FILE_NAME {
						// Determine the page name as the relative directory from the pages folder.
						relDir, err := filepath.Rel(pagesPath, filepath.Dir(path))
						if err != nil {
							fmt.Println(err)
							slog.Error("Error computing relative path for ", path, ": ", err)
							return err
						}
						pageName := relDir
						if pageName == "." {
							pageName = ""
						}
						// Use a key combining the domain and the pageName.
						routeKey := domain + "/" + pageName
						siteStructure.Routes[routeKey] = structure.PageRoutes{
							Name:   pageName,
							Title:  domain,
							Layout: "",
							Path:   path,
							// Template: string(templateData),
							MetaTags: structure.MetaTags{
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
					slog.Error("Error accessing component path ", path, ": ", err)
					return err
				}
				if !d.IsDir() && filepath.Ext(path) == structure.Wispy.FILE_EXT {
					// templateData, err := os.ReadFile(path)
					// if err != nil {
					// 	slog.Error("Failed to read component file at ", path, ": ", err)
					// 	return err
					// }
					componentName := strings.TrimSuffix(filepath.Base(path), structure.Wispy.FILE_EXT)
					siteStructure.Partials[componentName] = path //string(templateData)
				}
				return nil
			})

			// Handle Layouts: walk through the layouts directory.
			filepath.Walk(layoutsPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					slog.Error("Error accessing layout path ", path, ": ", err)
					return err
				}
				if !info.IsDir() && filepath.Ext(path) == structure.Wispy.FILE_EXT {
					templateData, err := os.ReadFile(path)
					if err != nil {
						slog.Error("Failed to read layout file at ", path, ": ", err)
						return err
					}
					layoutName := strings.TrimSuffix(filepath.Base(path), structure.Wispy.FILE_EXT)
					siteStructure.Layouts[layoutName] = string(templateData)
				}
				return nil
			})

			engine.SiteMap[domain] = siteStructure
		}
	}
	//
	fmt.Println("SiteMap Build Time: ", time.Since(buildStart))
	// Log the list of sites for confirmation.
	var domains []string
	for domain := range engine.SiteMap {
		domains = append(domains, domain)
	}
	fmt.Println("Sites: ", domains)
}
