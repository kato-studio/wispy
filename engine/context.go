package engine

import "github.com/tidwall/gjson"

type RenderCTX struct {
	Json      gjson.Result
	Snippet   map[string]string
	Variables map[string]string
}
