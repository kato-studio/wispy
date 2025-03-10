package template

import (
	"fmt"
	"strings"
)

type TemplateEngine struct {
	// starting deliminator - default "{%"
	DelimStart string
	// starting deliminator  - default "%}"
	DelimEnd string
	// used to trim the start and end delim as well as leading and trailing whitespace
	CutSet string // generated by NewTemplateEngine() using default deliminator's
	// map to check template tags against when rendering
	TagMap map[string]TemplateTag
	// map to check template filters against when rendering
	FilterMap map[string]EngineFilter
}

// Base function to create TemplateEngine instance used to to control base template settings
// RenderCtx references the Template engine to for engine settings like registered tags & filters
func NewTemplateEngine() *TemplateEngine {
	eng := &TemplateEngine{
		DelimStart: "{%",
		DelimEnd:   "%}",
		TagMap:     make(map[string]TemplateTag),
		FilterMap:  make(map[string]EngineFilter),
	}
	eng.GenCutSet(" \n\t") // generate initial CutSet whitespace

	eng.RegisterTags(DefaultTemplateTags)
	eng.RegisterFilters(DefaultTemplateFilters)

	return eng
}

// Resolves variables from the RenderCtx's Data map.
func (eng *TemplateEngine) GetFunc(ctx *RenderCtx, key string) any {
	if val, ok := ctx.Data[key]; ok {
		return val
	}
	return ""
}

// used to generate CutSet string for trimming tags when parsing
func (eng *TemplateEngine) GenCutSet(initial_cutset string) {
	var str = ConcatStrings(initial_cutset, eng.DelimStart, eng.DelimEnd)
	ia := len(str) / 2                  // initial allocation
	seen := make(map[rune]struct{}, ia) //fast lookup for weather a run has been seen if i'm not mistaken "struct{}" has a smaller memory footprint then "bool"
	cutset := make([]rune, ia)
	for _, r := range str {
		if _, ok := seen[r]; !ok {
			cutset = append(cutset, r)
			seen[r] = struct{}{}
		}
	}
	eng.CutSet = string(cutset)
}

// ----
// Registration functions
// ----
// Register Tag used in template {% <TAG> ... %} (I.E. {% if .user.loggedIn %} ... {% endif %}
func (eng *TemplateEngine) RegisterSingleTag(newTag TemplateTag) {
	_, exists := eng.TagMap[newTag.Name]
	if exists {
		fmt.Println("[Warning] \"" + newTag.Name + "\" a Tag with the same name had already been registered and was overridden")
	}
	eng.TagMap[newTag.Name] = newTag
}

// Register slice/array of Tags used in template {% <TAG> ... %} (I.E. {% if .user.loggedIn %} ... {% endif %}
func (eng *TemplateEngine) RegisterTags(tags []TemplateTag) {
	for _, t := range tags {
		eng.RegisterSingleTag(t)
	}
}

// Register Filter used in template {% .<VAR> | <FILTER> %} (I.E. {% .data.title | uppercase %} )
func (eng *TemplateEngine) RegisterSingleFilter(newFilter EngineFilter) {
	_, exists := eng.FilterMap[newFilter.Name]
	if exists {
		fmt.Println("[Warning] \"" + newFilter.Name + "\" a Filter with the same name had already been registered and was overridden")
	}
	eng.FilterMap[newFilter.Name] = newFilter
}

// Register slice/array of Filters used in template {% .<VAR> | <FILTER> %} (I.E. {% .data.title | uppercase %} )
func (eng *TemplateEngine) RegisterFilters(filters []EngineFilter) {
	for _, f := range filters {
		eng.RegisterSingleFilter(f)
	}
}

// RenderCtx represents the rendering context.
type RenderCtx struct {
	Engine        *TemplateEngine    // Reference to the TemplateEngine.
	Props         map[string]any     // Props passed to the component.
	Blocks        map[string]string  // Slots for block content.
	Data          map[string]any     // Data available during rendering.
	Partials      *map[string]string // Component like templates
	InternalFlags map[string]any     // For Tag handlers set flags for reference later in render operation
}

// TODO add methods for changing ctx instead of "allowing" directly setting variables
// NewRenderCtx creates a new RenderCtx with initialized internal state.
func NewRenderCtx(engine *TemplateEngine, data map[string]any, partials *map[string]string) *RenderCtx {
	return &RenderCtx{
		Engine:        engine,
		Data:          data,
		Blocks:        make(map[string]string),
		Props:         make(map[string]any),
		Partials:      partials,
		// 
		InternalFlags: make(map[string]any, 5),
	}
}

// related micro-utils
func CleanTemplateTag(ctx *RenderCtx, s string) string {
	return strings.Trim(s, ctx.Engine.CutSet)
}
