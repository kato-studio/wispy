package engine

import (
	"regexp"
)

var regex_number = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+$`)
var regex_insert = regexp.MustCompile(`{{.*?}}`)
var regex_attributes = regexp.MustCompile(`([-\w:])+=({{(.*?)}}|\"(.*?)\")`)

const EXT = ".hsml"
const ROOT_DIR = "./sites"
const LAYOUT_FILE = "/+layout" + EXT
const PAGE_FILE = "/+page" + EXT

var DOMAINS = [25]string{"studio", "com", "org", "net", "gov", "edu", "mil", "int", "ca", "co", "uk", "de", "jp", "fr", "au", "us", "ch", "it", "nl", "se", "no", "es", "mil", "io"}

type TemplatePageCtx struct {
	Title   string
	Head    string
	Meta    string
	Styles  string
	Scripts string
	Lang    string
	Layout  string
}
type TemplateCtx struct {
	Page       TemplatePageCtx
	Data       map[string]any
	Components map[string]string
	Site       map[string]string
}

func NewCtx(input TemplateCtx) *TemplateCtx {
	var result = &TemplateCtx{
		Data:       input.Data,
		Page:       input.Page,
		Components: make(map[string]string),
	}

	if result.Page.Lang == "" {
		result.Page.Lang = "eng"
	}
	if result.Page.Layout == "" {
		result.Page.Layout = "./shared/layouts/_default" + LAYOUT_FILE
	}

	return result
}
