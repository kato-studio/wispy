package engine

import "github.com/tidwall/gjson"

const LAYOUT_FILE = "+layout.hstm"
const PAGE_FILE = "+page.hstm"

type RenderCTX struct {
	Json      gjson.Result
	Snippet   map[string]string
	Variables map[string]string
}
