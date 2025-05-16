package template

import (
	"log/slog"
	"os"

	"github.com/kato-studio/wispy/template/core"
	"github.com/kato-studio/wispy/template/filters"
	"github.com/kato-studio/wispy/template/structure"
	"github.com/kato-studio/wispy/template/tags"
)

var SiteMap = map[string]*structure.SiteStructure{}
var Logger *slog.JSONHandler

/*
=================================================================
Core External Functions
=================================================================
*/
func init() {
	// -------------
	// Setup Logger
	// -------------
	// logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// defer logFile.Close()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

// EngineCtx is the engine context which holds site mappings and configuration.
var Wispy = &structure.WispyConfig{
	SITE_DIR:         "./sites",
	PAGE_FILE_NAME:   "page",
	FILE_EXT:         ".hstm",
	SITE_CONFIG_NAME: "config.toml",
}

var DefaultTemplateFilters = []structure.TemplateFilter{
	filters.UpcaseFilter,
	filters.DowncaseFilter,
	filters.CapitalizeFilter,
	filters.StripFilter,
	filters.TruncateFilter,
	filters.SliceFilter,
}

var DefaultTemplateTags = []structure.TemplateTag{
	tags.IfTag,
	tags.EachTag,
	tags.CommentTag,
	tags.DefineTag,
	tags.BlockTag,
}

var DefaultEngineTags = []structure.TemplateTag{
	tags.IfTag,
	tags.EachTag,
	tags.PartialTag,
	tags.CommentTag,
	tags.DefineTag,
	tags.BlockTag,
	tags.ExtendsTag,
	tags.LayoutTag,
	tags.PassedTag,
	//
	tags.HeadTag,
	tags.FooterAssetsTag,
	tags.TitleTag,
	tags.MetaTag,
	tags.ImportJSTag,
	tags.ImportCSSTag,
	tags.CSSTag,
	tags.JSTag,
	tags.ImportTag,
}

func StartDefaultEngine() *structure.TemplateEngine {
	var engine = structure.TemplateEngine{}
	return engine.Init(DefaultEngineTags, DefaultTemplateFilters)
}

var Render = core.Render
