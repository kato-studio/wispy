// template.go
package template

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

// ANSI color codes for log messages.
const (
	colorReset   = "\033[0m"
	colorBlue    = "\033[34m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorRed     = "\033[31m"
)

// TemplateFunc defines a command's behavior.
type TemplateFunc struct {
	IsBlock bool
	Render  func(ctx *RenderCtx, args ...any) string
}

// TemplateEngine holds the function and filter maps along with an optional custom get handler.
type TemplateEngine struct {
	FuncMap   map[string]TemplateFunc
	FilterMap map[string]func(interface{}, []string) interface{}
	// GetFunc returns the value for a variable; the default checks both Data and Internal.
	GetFunc func(ctx *RenderCtx, key string) interface{}
}

// NewTemplateEngine creates a new TemplateEngine instance with default functions and filters.
func NewTemplateEngine() *TemplateEngine {
	eng := &TemplateEngine{
		FuncMap:   make(map[string]TemplateFunc),
		FilterMap: make(map[string]func(interface{}, []string) interface{}),
	}
	eng.GetFunc = func(ctx *RenderCtx, key string) interface{} {
		if val, ok := ctx.Data[key]; ok {
			return val
		}
		if val, ok := ctx.Internal[key]; ok {
			return val
		}
		return ""
	}

	// Register filters from filters.go.
	eng.FilterMap["upcase"] = UpcaseFilter
	eng.FilterMap["downcase"] = DowncaseFilter
	eng.FilterMap["capitalize"] = CapitalizeFilter
	eng.FilterMap["strip"] = StripFilter
	eng.FilterMap["truncate"] = TruncateFilter
	eng.FilterMap["slice"] = SliceFilter

	// Register command functions from operations.go.
	eng.FuncMap["render"] = TemplateFunc{
		IsBlock: false,
		Render:  RenderCommand,
	}
	eng.FuncMap["component"] = TemplateFunc{
		IsBlock: true,
		Render:  ComponentCommand(eng),
	}
	eng.FuncMap["if"] = TemplateFunc{
		IsBlock: true,
		Render:  IfCommand(eng),
	}
	eng.FuncMap["for"] = TemplateFunc{
		IsBlock: true,
		Render:  ForCommand(eng),
	}

	return eng
}

// RenderCtx carries rendering data and internal state.
type RenderCtx struct {
	Data     map[string]interface{}
	Internal map[string]interface{}
}

// NewRenderCtx creates a new RenderCtx with initialized internal state.
func NewRenderCtx(data map[string]interface{}) *RenderCtx {
	return &RenderCtx{
		Data:     data,
		Internal: make(map[string]interface{}),
	}
}

// nextTag returns the index and type ("var" or "code") of the next tag.
func nextTag(input string, pos int) (int, string) {
	idxVar := strings.Index(input[pos:], "{{")
	idxCode := strings.Index(input[pos:], "{%")
	if idxVar != -1 {
		idxVar += pos
	}
	if idxCode != -1 {
		idxCode += pos
	}
	if idxVar == -1 && idxCode == -1 {
		return -1, ""
	} else if idxVar == -1 {
		return idxCode, "code"
	} else if idxCode == -1 {
		return idxVar, "var"
	} else if idxVar < idxCode {
		return idxVar, "var"
	} else {
		return idxCode, "code"
	}
}

// evaluateExpression propagates an initial value (from GetFunc) through a chain of filters.
func (engine *TemplateEngine) evaluateExpression(tokens []string, ctx *RenderCtx) interface{} {
	if len(tokens) == 0 {
		return ""
	}
	varName := strings.TrimPrefix(tokens[0], ".")
	var value interface{}
	if engine.GetFunc != nil {
		value = engine.GetFunc(ctx, varName)
	} else {
		if v, ok := ctx.Data[varName]; ok {
			value = v
		}
	}
	i := 1
	for i < len(tokens) {
		if tokens[i] == "|" && i+1 < len(tokens) {
			filterToken := tokens[i+1]
			filterName, filterArgs := parseFilter(filterToken)
			if filterFunc, ok := engine.FilterMap[filterName]; ok {
				value = filterFunc(value, filterArgs)
			}
			i += 2
		} else {
			i++
		}
	}
	return value
}

// Node represents either literal text or a command.
type Node struct {
	Text    string
	Command *Command
}

// Command represents a parsed command from a tag.
type Command struct {
	Name string
	Args []string
	Body []Node
}

// parseTemplate scans the input string and returns a slice of Nodes.
// It uses "{{" ... "}}" for variables and "{%" ... "%}" for code.
func parseTemplate(input string, pos int, expectedEnd string, engine *TemplateEngine) ([]Node, int, error) {
	var nodes []Node
	for pos < len(input) {
		idx, tagType := nextTag(input, pos)
		if idx == -1 {
			nodes = append(nodes, Node{Text: input[pos:]})
			pos = len(input)
			break
		}
		if idx > pos {
			nodes = append(nodes, Node{Text: input[pos:idx]})
			pos = idx
		}
		var closeDelim string
		if tagType == "var" {
			closeDelim = "}}"
		} else {
			closeDelim = "%}"
		}
		endTagPos := strings.Index(input[pos:], closeDelim)
		if endTagPos == -1 {
			return nil, pos, fmt.Errorf("unterminated tag at pos %d", pos)
		}
		// Always skip 2 characters for the opening delimiter.
		tagContent := strings.TrimSpace(input[pos+2 : pos+endTagPos])
		pos += endTagPos + len(closeDelim)

		// If this is a code tag and an end tag is indicated, handle accordingly.
		if tagType == "code" && strings.HasPrefix(tagContent, "/") {
			tagContent = strings.TrimSpace(tagContent[1:])
			if expectedEnd != "" {
				if tagContent != expectedEnd {
					return nil, pos, fmt.Errorf("unexpected end tag '{%%/%s%%}', expected '{%%/%s%%}'", tagContent, expectedEnd)
				}
				return nodes, pos, nil
			}
			return nil, pos, fmt.Errorf("unexpected end tag '{%%/%s%%}' at pos %d", tagContent, pos)
		}

		tokens := tokenize(tagContent)
		if len(tokens) == 0 {
			return nil, pos, fmt.Errorf("empty tag at pos %d", pos)
		}
		cmdName := tokens[0]
		args := tokens[1:]
		cmd := Command{
			Name: cmdName,
			Args: args,
		}
		if tagType == "code" {
			if tf, ok := engine.FuncMap[cmdName]; ok && tf.IsBlock {
				body, newPos, err := parseTemplate(input, pos, cmdName, engine)
				if err != nil {
					return nil, pos, err
				}
				cmd.Body = body
				pos = newPos
			}
		}
		nodes = append(nodes, Node{Command: &cmd})
	}
	if expectedEnd != "" {
		return nil, pos, fmt.Errorf("expected end tag '{%%/%s%%}' not found", expectedEnd)
	}
	return nodes, pos, nil
}

// tokenize splits a string into tokens by whitespace while respecting quoted substrings.
func tokenize(s string) []string {
	var tokens []string
	var token strings.Builder
	inQuote := false
	var quoteChar rune
	for _, r := range s {
		if inQuote {
			token.WriteRune(r)
			if r == quoteChar {
				inQuote = false
				tokens = append(tokens, token.String())
				token.Reset()
			}
		} else {
			if unicode.IsSpace(r) {
				if token.Len() > 0 {
					tokens = append(tokens, token.String())
					token.Reset()
				}
			} else if r == '"' || r == '\'' {
				inQuote = true
				quoteChar = r
				token.WriteRune(r)
			} else {
				token.WriteRune(r)
			}
		}
	}
	if token.Len() > 0 {
		tokens = append(tokens, token.String())
	}
	return tokens
}

// parseFilter splits a filter token into its name and arguments.
func parseFilter(token string) (string, []string) {
	parts := strings.SplitN(token, ":", 2)
	name := parts[0]
	var args []string
	if len(parts) == 2 {
		args = strings.Split(parts[1], ",")
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)
		}
	}
	return name, args
}

// compileCommand compiles a Command node into a render function.
func (engine *TemplateEngine) compileCommand(cmd *Command) func(ctx *RenderCtx) string {
	// For variable tags, treat as inline expressions.
	if strings.HasPrefix(cmd.Name, ".") {
		tokens := append([]string{cmd.Name}, cmd.Args...)
		return func(ctx *RenderCtx) string {
			return fmt.Sprintf("%v", engine.evaluateExpression(tokens, ctx))
		}
	}
	convertArgs := func(args []string) []any {
		r := make([]any, len(args))
		for i, a := range args {
			r[i] = a
		}
		return r
	}
	if tf, ok := engine.FuncMap[cmd.Name]; ok {
		if tf.IsBlock {
			return func(ctx *RenderCtx) string {
				args := convertArgs(cmd.Args)
				args = append(args, cmd.Body)
				return tf.Render(ctx, args...)
			}
		}
		return func(ctx *RenderCtx) string {
			return tf.Render(ctx, convertArgs(cmd.Args)...)
		}
	}
	return func(ctx *RenderCtx) string {
		return fmt.Sprintf("{%% %s %%}", strings.Join(cmd.Args, " "))
	}
}

// compileNodes compiles a slice of Nodes into a render function.
func (engine *TemplateEngine) compileNodes(nodes []Node) func(ctx *RenderCtx) string {
	compiledFuncs := make([]func(ctx *RenderCtx) string, len(nodes))
	for i, node := range nodes {
		if node.Command != nil {
			compiledFuncs[i] = engine.compileCommand(node.Command)
		} else {
			text := node.Text
			compiledFuncs[i] = func(ctx *RenderCtx) string { return text }
		}
	}
	return func(ctx *RenderCtx) string {
		var sb strings.Builder
		for _, f := range compiledFuncs {
			sb.WriteString(f(ctx))
		}
		return sb.String()
	}
}

// RenderTemplateFromString renders a template from a string.
func RenderString(input string, engine *TemplateEngine, ctx *RenderCtx) (string, error) {
	nodes, pos, err := parseTemplate(input, 0, "", engine)
	if err != nil {
		return "", err
	}
	if pos != len(input) {
		return "", fmt.Errorf("not all input was consumed (stopped at pos %d)", pos)
	}
	renderFunc := engine.compileNodes(nodes)
	return renderFunc(ctx), nil
}

// RenderTemplateFromFile reads a template file, compiles it using the engine, and returns the rendered output.
func RenderFile(filePath string, engine *TemplateEngine, ctx *RenderCtx, print bool) (string, error) {
	start := time.Now()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	log.Printf("%sRead file in %v%s", colorCyan, time.Since(start), colorReset)
	input := string(data)
	parseStart := time.Now()
	nodes, pos, err := parseTemplate(input, 0, "", engine)
	if print {
		printNodes(nodes, 2)
	}
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}
	if pos != len(input) {
		log.Printf("%sWarning: not all input was consumed, stopped at pos %d%s", colorYellow, pos, colorReset)
	}
	log.Printf("%sParsed template in %v%s", colorGreen, time.Since(parseStart), colorReset)
	compileStart := time.Now()
	renderFunc := engine.compileNodes(nodes)
	log.Printf("%sCompiled template in %v%s", colorYellow, time.Since(compileStart), colorReset)
	renderStart := time.Now()
	result := renderFunc(ctx)
	log.Printf("%sRendered template in %v%s", colorMagenta, time.Since(renderStart), colorReset)
	total := time.Since(start)
	log.Printf("%sTotal execution time: %v%s", colorRed, total, colorReset)
	return result, nil
}

// printNodes is a helper to pretty-print the parsed AST.
func printNodes(nodes []Node, indent int) {
	indentStr := strings.Repeat("  ", indent)
	for _, n := range nodes {
		if n.Command != nil {
			fmt.Printf("%sCommand: %s\n", indentStr, n.Command.Name)
			if len(n.Command.Args) > 0 {
				fmt.Printf("%s  Args: %v\n", indentStr, n.Command.Args)
			}
			if len(n.Command.Body) > 0 {
				fmt.Printf("%s  Body:\n", indentStr)
				printNodes(n.Command.Body, indent+2)
			}
		} else if len(n.Text) > 0 {
			trimmed := strings.TrimSpace(n.Text)
			if trimmed != "" {
				fmt.Printf("%sText: %q\n", indentStr, trimmed)
			}
		}
	}
}
