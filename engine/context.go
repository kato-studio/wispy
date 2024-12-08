package engine

import (
	"regexp"

	"github.com/tidwall/gjson"
)

var regex_number = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+$`)
var regex_insert = regexp.MustCompile(`{{.*?}}`)
var regex_attributes = regexp.MustCompile(`([-\w:])+=({{(.*?)}}|\"(.*?)\")`)

const EXT = ".htsm"
const ROOT_DIR = "./sites"
const LAYOUT_FILE = "/+layout" + EXT
const PAGE_FILE = "/+page" + EXT
const DOCUMENT_FILE = "/_document" + EXT

var DOMAINS = [25]string{"studio", "com", "org", "net", "gov", "edu", "mil", "int", "ca", "co", "uk", "de", "jp", "fr", "au", "us", "ch", "it", "nl", "se", "no", "es", "mil", "io"}

type RenderCTX struct {
	Json       gjson.Result
	Components map[string]string
	Head       map[string]string
	Variables  map[string]string
}
