package engine

import "github.com/labstack/echo/v4"

type RenderCtx struct {
	Engine *EngineCtx
	Data   map[string]interface{}
	Site   SiteStructure
	Lang   string
}

func NewRenderCtx(e *EngineCtx, data map[string]interface{}, site SiteStructure) RenderCtx {
	return RenderCtx{
		Engine: e,
		Data:   data,
		Site:   site,
		Lang:   "en",
	}
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
	SiteMap map[string]SiteStructure // List of domains/sites from ./[SITE_DIR] for validation & routing
	Config  WispyConfig
	Log     echo.Logger // accessible list of configured settings
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
