package engine

// var regex_number = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+$`)
// var regex_insert = regexp.MustCompile(`{{.*?}}`)
// var regex_attributes = regexp.MustCompile(`([-\w:])+=({{(.*?)}}|\"(.*?)\")`)

const EXT = ".hsml"
const ROOT_DIR = "./sites"
const SHARED_DIR = "./shared"
const LAYOUT_FILE = "/+layout" + EXT
const PAGE_FILE = "/+page" + EXT

var DOMAINS = [25]string{"studio", "com", "org", "net", "gov", "edu", "mil", "int", "ca", "co", "uk", "de", "jp", "fr", "au", "us", "ch", "it", "nl", "se", "no", "es", "mil", "io"}
var SELF_CLOSING_TAGS = []string{"img", "br", "hr", "input", "link", "meta", "area", "base", "col", "command", "embed", "keygen", "param", "source", "track", "wbr"}

type TemplatePageCtx struct {
	Title   string
	Head    string
	Meta    string
	Styles  string
	Scripts string
	Lang    string
	Layout  string
}

type TemplateSiteCtx struct {
	Name   string
	Domain string
}

type TemplateCtx struct {
	Page       TemplatePageCtx
	Data       map[string]any
	Components map[string]string
	Site       TemplateSiteCtx
}

func NewCtx(input TemplateCtx) *TemplateCtx {
	var result = &TemplateCtx{
		Data: input.Data,
		Page: input.Page,
	}

	if result.Page.Lang == "" {
		result.Page.Lang = "eng"
	}
	if result.Page.Layout == "" {
		result.Page.Layout = "./shared/layouts/_default" + LAYOUT_FILE
	}

	return result
}

// Render represents the rendering context and functions.
type Render struct {
	AttrFuncMap      map[string]AttributeFunc
	OperationFuncMap map[string]OperationFunc
	Ctx              *TemplateCtx
	GetComponent     func(name string) (rawBytes []byte, err error)
}

// OperationFunc defines the function signature for operations.
type OperationFunc func(r *Render, values ...string) string

// AttributeFunc defines the function signature for attribute functions.
type AttributeFunc func(name, value string) (changed bool, newAttr string, errs []error)

// TreeNode represents a node in the HTML tree.
type TreeNode struct {
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Content    string            `json:"content,omitempty"`
	Children   []*TreeNode       `json:"children,omitempty"`
}

// InitEngine initializes the rendering engine with a given context.
func InitEngine(ctx *TemplateCtx, attributeFuncMap map[string]AttributeFunc, operationFuncMap map[string]OperationFunc) Render {
	return Render{
		AttrFuncMap:      attributeFuncMap,
		OperationFuncMap: operationFuncMap,
		Ctx:              ctx,
	}
}
