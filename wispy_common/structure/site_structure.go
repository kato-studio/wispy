package structure

// WispyConfig holds configuration options for the engine.
type WispyConfig struct {
	SITE_DIR         string
	PAGE_FILE_NAME   string
	FILE_EXT         string
	SITE_CONFIG_NAME string
}

// EngineCtx is the engine context which holds site mappings and configuration.
var Wispy = &WispyConfig{
	SITE_DIR:         "./sites",
	PAGE_FILE_NAME:   "page",
	FILE_EXT:         ".hstm",
	SITE_CONFIG_NAME: "config.toml",
}

// SiteStructure represents a single site, with its pages, layouts, and partials.
type SiteStructure struct {
	Domain   string
	Routes   map[string]PageRoutes
	Layouts  map[string]string
	Partials map[string]string
	// Add other fields as needed (e.g., site-specific config settings)
}

// PageRoutes holds information about a page.
type PageRoutes struct {
	Name     string
	Title    string
	Layout   string
	Path     string
	Template string
	MetaTags MetaTags
}

// MetaTags holds metadata information for a page.
type MetaTags struct {
	Title         string
	Description   string
	OgTitle       string
	OgDescription string
	OgType        string
	OgUrl         string
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

// ----------------------
//
//	CONSTS
//
// ----------------------
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
