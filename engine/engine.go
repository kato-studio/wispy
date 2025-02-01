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

type EngineCtx struct {
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
	Pages         map[string]string
	Layouts       map[string]string
	Components    map[string]string
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
	Name     string
	Url      string
	Title    string
	Layout   string
	Path     string
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

func (e *EngineCtx) NewSiteStructure(domain string) SiteStructure {
	return SiteStructure{
		Domain:        domain,
		Routes:        map[string]PageRoutes{},
		Pages:         make(map[string]string, 6),
		Layouts:       make(map[string]string, 2),
		Components:    make(map[string]string, 6),
		ContentRoutes: map[string]SiteContent{},
	}
}

func StartEngine(config WispyConfig, logger echo.Logger) EngineCtx {
	return EngineCtx{
		Templates: make(map[string]*template.Template, 20),
		SiteMap:   make(map[string]SiteStructure, 5),
		Log:       logger,
		Config: WispyConfig{
			SITE_DIR:           ".\\sites",
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
func (e *EngineCtx) SetupWispyCache() {
	cacheDir := ".wispy"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("Failed to create .wispy directory: %v", err))
		}
	}
}

// Dynamically build the host-to-site mapping based on `./sites` directory
func (e *EngineCtx) BuildSiteMap() {
	buildStart := time.Now()
	entries, err := os.ReadDir(e.Config.SITE_DIR)
	if err != nil {
		panic(fmt.Sprintf("Failed to read sites directory: %v", err))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			var domain = entry.Name()
			siteFolderPath := e.Config.SITE_DIR + "\\" + domain
			configFilePath := siteFolderPath + "\\" + e.Config.SITE_CONFIG_NAME
			file, err := os.ReadFile(configFilePath)
			if err != nil {
				e.Log.Error("Could not find config for ", domain, " at: (", configFilePath, ")")
				e.Log.Error(err)
			}

			// Create new site object
			var siteStructure SiteStructure = e.NewSiteStructure(domain)
			_, err = toml.Decode(string(file), &siteStructure)
			if err != nil {
				e.Log.Error("Failed to load config for ", domain, " at: (", configFilePath, ")")
				e.Log.Error(err)
			}

			pagesPath := siteFolderPath + "\\pages"
			layoutsPath := siteFolderPath + "\\layouts"
			componentsPath := siteFolderPath + "\\components"

			// Handles Pages
			filepath.Walk(pagesPath, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
					return err
				}
				//
				path = filepath.FromSlash(path)
				pagePath, isPageFile := strings.CutSuffix(path, e.Config.PAGE_FILE_NAME+e.Config.FILE_EXT)
				if isPageFile {
					pageName := strings.TrimPrefix(pagePath, "sites\\"+domain+"\\pages")
					pageName = strings.TrimSuffix(pageName, "\\")
					templateData, err := os.ReadFile(path)
					if err != nil {
						e.Log.Error("Failed to read templateData file at:", path)
					}
					siteStructure.Routes[domain+pageName] = PageRoutes{
						Name:     pageName,
						Title:    domain,
						Layout:   "sites\\abc.test\\layouts\\default.html",
						Path:     path,
						Template: string(templateData),
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

			// Handle Components
			filepath.WalkDir(componentsPath, func(path string, dr fs.DirEntry, err error) error {
				if err != nil {
					e.Log.Error("Error parsing comps %s: %v", path, err)
					return err
				}
				if filepath.Ext(path) == e.Config.FILE_EXT {
					templateData, err := os.ReadFile(path)
					if err != nil {
						e.Log.Error(err)
						return err
					}
					componentName := strings.TrimSuffix(filepath.Base(path), e.Config.FILE_EXT)
					siteStructure.Components[componentName] = string(templateData)
				}
				return nil
			})

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

					siteStructure.Layouts[layoutName] = string(templateData)
				}
				return nil
			})

			e.SiteMap[domain] = siteStructure
		}
	}

	fmt.Println("Site Map Build Time: ", time.Since(buildStart))
}
