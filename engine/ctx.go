package engine

import "github.com/labstack/echo/v4"

// WispyConfig holds configuration options for the engine.
type WispyConfig struct {
	SITE_DIR           string
	PAGE_FILE_NAME     string
	FILE_EXT           string
	SHARED_COMP_PREFIX string
	SHARED_DIR         string
	SITE_CONFIG_NAME   string
}

// EngineCtx is the engine context which holds site mappings and configuration.

// TODO: move to internal
var Wispy = &WispyConfig{
	SITE_DIR:           "sites",
	PAGE_FILE_NAME:     "page",
	FILE_EXT:           ".html",
	SHARED_COMP_PREFIX: "wispy",
	// PUBLIC_DIR:         "public",
	SHARED_DIR:       "shared",
	SITE_CONFIG_NAME: "config.toml",
}
var Echo = echo.New()
var Log = Echo.Logger
var SiteMap = map[string]SiteStructure{}

// ------

type RenderCtx struct {
	Data map[string]interface{}
	Site SiteStructure
	Lang string
}

func NewRenderCtx(data map[string]interface{}, site SiteStructure) RenderCtx {
	return RenderCtx{
		Data: data,
		Site: site,
		Lang: "en",
	}
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
