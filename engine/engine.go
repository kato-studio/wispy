package engine

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v4"

	"github.com/kato-studio/wispy/engine/templateFuncs"
)

// CONSTS
// Define the essential favicon filenames
var ESSENTIAL_SERVE = map[string]struct{}{
	"about.txt":                  {},
	"android-chrome-192x192.png": {},
	"android-chrome-512x512.png": {},
	"apple-touch-icon.png":       {},
	"favicon-16x16.png":          {},
	"favicon-32x32.png":          {},
	"favicon.ico":                {},
	"site.webmanifest":           {},
}

// TemplateRenderer supports multi-site template rendering
type WispyConfig struct {
	SITE_DIR           string
	SITE_CONFIG_NAME   string
	PAGE_FILE_NAME     string
	FILE_EXT           string
	SHARED_COMP_PREFIX string
	// PUBLIC_DIR         string
	SHARED_DIR string
}

type Engine struct {
	// Templates map[string]map[string]*template.Template // Templates per site
	Templates map[string]*template.Template // Templates per site
	SiteMap   map[string]SiteStructure      // List of domains/sites from ./[SITE_DIR] for validation & routing
	Config    WispyConfig
	Log       echo.Logger // accessible list of configured settings
}

type SiteStructure struct {
	// import from config
	Name  string                       `toml:"name"`
	Theme map[string]map[string]string `toml:"theme"`
	//
	Domain        string // set based on directory name
	Pages         map[string]*template.Template
	Layouts       map[string]*template.Template
	Components    *template.Template
	Routes        map[string]PageRoutes
	ContentRoutes map[string]SiteContent
}
type MetaTags struct {
	Title         string
	Description   string
	OgTitle       string
	OgDescription string
	OgType        string
	OgUrl         string
	OtherTags     map[string]string
}
type PageRoutes struct {
	Url      string
	Title    string
	Layout   string
	Template string
	MetaTags
}
type ContentChange struct {
	Author  string
	Date    string
	Changes map[string]string
}
type SiteContent struct {
	Name        string
	Title       string
	Description string
	Slug        string
	Category    string
	Tags        string
	Author      string
	LastUpdate  string
	Changes     map[string]ContentChange
}

func NewSiteStructure(domain string) SiteStructure {
	return SiteStructure{
		Domain:        domain,
		Routes:        map[string]PageRoutes{},
		Pages:         make(map[string]*template.Template, 6),
		Layouts:       make(map[string]*template.Template),
		Components:    template.New(domain + "-components").Funcs(templateFuncs.GetDefaults()),
		ContentRoutes: map[string]SiteContent{},
	}
}

func StartEngine(config WispyConfig, logger echo.Logger) Engine {
	return Engine{
		Templates: make(map[string]*template.Template, 20),
		SiteMap:   make(map[string]SiteStructure, 5),
		Log:       logger,
		Config: WispyConfig{
			SITE_DIR:           "./sites",
			PAGE_FILE_NAME:     "page",
			FILE_EXT:           ".html",
			SHARED_COMP_PREFIX: "wispy",
			// PUBLIC_DIR:         "public",
			SHARED_DIR:       "shared",
			SITE_CONFIG_NAME: "config.toml",
		},
	}
}

// ---====----
// Setup / Initialization functions
// ---====----

// Initialize and configure the .wispy cache directory
func (e *Engine) SetupWispyCache() {
	cacheDir := ".wispy"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Failed to create .wispy directory: %v", err))
		}
	}
}

// Dynamically build the host-to-site mapping based on `./sites` directory
func (e *Engine) BuildSiteMap() {
	buildStart := time.Now()
	entries, err := os.ReadDir(e.Config.SITE_DIR)
	if err != nil {
		panic(fmt.Sprintf("Failed to read sites directory: %v", err))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			var domain = entry.Name()
			siteFolderPath := e.Config.SITE_DIR + "/" + domain
			configFilePath := siteFolderPath + "/" + e.Config.SITE_CONFIG_NAME
			file, err := os.ReadFile(configFilePath)
			if err != nil {
				e.Log.Error("Could not find config for ", domain, " at: (", configFilePath, ")")
				e.Log.Error(err)
			}

			// Create new site object
			var siteStructure SiteStructure = NewSiteStructure(domain)
			_, err = toml.Decode(string(file), &siteStructure)
			if err != nil {
				e.Log.Error("Failed to load config for ", domain, " at: (", configFilePath, ")")
				e.Log.Error(err)
			}

			pagesPath := siteFolderPath + "/pages"
			layoutsPath := siteFolderPath + "/layouts"
			componentsPath := siteFolderPath + "/components"

			// Handles Pages
			filepath.Walk(pagesPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
					return err
				}
				//
				pagePath, isPageFile := strings.CutSuffix(path, e.Config.PAGE_FILE_NAME+e.Config.FILE_EXT)
				if isPageFile {
					templateName := strings.Replace(pagePath, "sites/"+domain+"/pages", domain, 1)
					templateName = strings.TrimSuffix(templateName, "/")
					templateData, err := os.ReadFile(path)
					if err != nil {
						e.Log.Error("Failed to read templateData file at:", path)
					}

					newTemplate, err := template.New(templateName).Parse(string(templateData))
					if err != nil {
						e.Log.Error("Failed to create template from file at:", path, err)
					}
					siteStructure.Pages[templateName] = newTemplate
					siteStructure.Routes[templateName] = PageRoutes{
						Url:      templateName,
						Title:    domain,
						Layout:   "sites/abc.test/layouts/default.html",
						Template: path,
						MetaTags: MetaTags{
							Title:         "domain title",
							Description:   "page description here boop",
							OgTitle:       "domain title",
							OgDescription: "page description here boop",
							OgType:        "text",
							OgUrl:         domain,
						},
					}
				}
				return nil
			})

			fmt.Println("---")
			// Handle Components
			filepath.WalkDir(componentsPath, func(path string, dr fs.DirEntry, err error) error {
				fmt.Println("Walking... ", path)
				if err != nil {
					e.Log.Error("Error parsing comps %s: %v", path, err)
					return err
				}
				fmt.Println(filepath.Ext(path), filepath.Ext(path) == e.Config.FILE_EXT)
				if filepath.Ext(path) == e.Config.FILE_EXT {
					_, err := siteStructure.Components.ParseFiles(path)
					if err != nil {
						e.Log.Error("Error parsing comps %s: %v", path, err)
					}
				}
				return nil
			})
			fmt.Println("Components")
			fmt.Println(siteStructure.Components.DefinedTemplates())
			fmt.Println("---")

			// Handle Layouts
			filepath.Walk(layoutsPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
					return err
				}

				thisLayoutPath, isLayoutFile := strings.CutSuffix(path, e.Config.FILE_EXT)
				layoutName := filepath.Base(thisLayoutPath)
				if isLayoutFile {
					templateData, err := os.ReadFile(path)
					if err != nil {
						e.Log.Error("Failed to read templateData file at:", path)
					}

					newTemplate, err := template.New(layoutName).Parse(string(templateData))
					if err != nil {
						e.Log.Error("Failed to create template from file at:", path, err)
					}
					siteStructure.Layouts[layoutName] = newTemplate
				}
				return nil
			})

			e.SiteMap[domain] = siteStructure
		}
	}

	fmt.Println("Site Map Build Time: ", time.Since(buildStart))
}
