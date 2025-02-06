package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"strings"
)

/*
=================================================================
Core Template Engine Structure and Initialization
=================================================================
*/

// TemplateEngine handles template parsing and execution with registered functions
type TemplateEngine struct {
	funcs map[string]interface{} // Map of registered template functions

}

// NewTemplateEngine creates a new template engine with core functions
// Example:
//
//	ctx := &RenderCtx{...}
//	engine := NewTemplateEngine(ctx)
func NewTemplateEngine() *TemplateEngine {
	engine := &TemplateEngine{
		funcs: make(map[string]interface{}),
	}
	// Register core functions.
	engine.RegisterFunc("component", engine.componentFunc)
	engine.RegisterFunc("include", engine.includeFunc)
	engine.RegisterFunc("props", Props)
	engine.RegisterFunc("slot", engine.slotFunc)
	engine.RegisterFunc("safeHTML", engine.safeHTMLFunc)
	// Register the minimal toJSON for test purposes.
	engine.RegisterFunc("toJSON", toJSONFunc)
	return engine
}

// RegisterFunc adds a custom function to the template engine
// Example:
//
//	engine.RegisterFunc("upper", func(s string) string { return strings.ToUpper(s) })
func (e *TemplateEngine) RegisterFunc(name string, fn interface{}) {
	e.funcs[name] = fn
}

/*
=================================================================
Main Execution Flow
=================================================================
*/

// Execute processes the template with given data
// Example:
//
//	result, err := engine.Execute("<h1>{{.Title}}</h1>", map[string]interface{}{"Title": "Hello"})
func (e *TemplateEngine) Execute(template string, data interface{}) (string, error) {
	var result bytes.Buffer
	lines := strings.Split(template, "\n")

	for i := 0; i < len(lines); i++ {
		processed, newIndex, err := e.processLine(lines, i, data)
		if err != nil {
			return "", err
		}
		result.WriteString(processed)
		i = newIndex
	}

	return result.String(), nil
}

/*
=================================================================
Template Processing Pipeline
=================================================================
*/

// processLine handles individual template lines and blocks
func (e *TemplateEngine) processLine(lines []string, index int, data interface{}) (string, int, error) {
	line := lines[index]

	// Handle control structures
	if strings.Contains(line, "{{if") {
		return e.processIfBlock(lines, index, data)
	}
	if strings.Contains(line, "{{include") {
		return e.processIncludeBlock(lines, index, data)
	}

	// Process inline tags and variables
	processed, err := e.processInlineTags(line, data)
	return processed + "\n", index, err
}

// processInlineTags handles inline template tags and variables
func (e *TemplateEngine) processInlineTags(line string, data interface{}) (string, error) {
	var result strings.Builder
	remaining := line

	for {
		start := strings.Index(remaining, "{{")
		if start == -1 {
			result.WriteString(remaining)
			break
		}

		end := strings.Index(remaining, "}}")
		if end == -1 {
			result.WriteString(remaining)
			break
		}

		// Process content before tag
		result.WriteString(remaining[:start])

		// Extract and process tag content
		tagContent := strings.TrimSpace(remaining[start+2 : end])
		processed, err := e.processTag(tagContent, data)
		if err != nil {
			return "", err
		}
		result.WriteString(processed)

		// Move to next tag
		remaining = remaining[end+2:]
	}

	return result.String(), nil
}

/*
=================================================================
Control Structure Handlers
=================================================================
*/

// processIfBlock handles conditional blocks
func (e *TemplateEngine) processIfBlock(lines []string, index int, data interface{}) (string, int, error) {
	var result bytes.Buffer
	startLine := lines[index]

	// Parse condition expression
	conditionExpr := strings.TrimSpace(strings.TrimPrefix(strings.Split(startLine, "{{if ")[1], "}}"))
	condition, err := e.evaluateCondition(conditionExpr, data)
	if err != nil {
		return "", index, err
	}

	// Process block content based on condition
	var content bytes.Buffer
	depth := 1
	currentIndex := index + 1

	for ; currentIndex < len(lines); currentIndex++ {
		line := lines[currentIndex]
		if strings.Contains(line, "{{if") {
			depth++
		}
		if strings.Contains(line, "{{end}}") {
			depth--
			if depth == 0 {
				break
			}
		}
		if condition {
			content.WriteString(line + "\n")
		}
	}

	// Execute conditional content
	if condition {
		processed, err := e.Execute(content.String(), data)
		if err != nil {
			return "", currentIndex, err
		}
		result.WriteString(processed)
	}

	return result.String(), currentIndex, nil
}

// processIncludeBlock handles component inclusion with slot content
func (e *TemplateEngine) processIncludeBlock(lines []string, index int, data interface{}) (string, int, error) {
	line := lines[index]
	var result bytes.Buffer

	// Parse include arguments. Expect function-call style: e.g., include "Header" (props("title" "Demo Page"))
	argsStr := strings.TrimSpace(strings.TrimSuffix(strings.Split(line, "{{include")[1], "}}"))
	// Here we use our parseArgs which supports whitespace-separated tokens if no commas are found.
	args, err := parseArgs(argsStr)
	if err != nil {
		return "", index, err
	}

	// Validate component name
	if len(args) < 1 {
		return "", index, fmt.Errorf("include requires component name")
	}
	componentName := strings.Trim(args[0], "\"")

	// Process props data if provided
	var propsData interface{}
	if len(args) > 1 {
		propsData, err = e.evalExpression(strings.Join(args[1:], " "), data)
		if err != nil {
			return "", index, err
		}
	}

	// Capture slot content
	var slotContent bytes.Buffer
	currentIndex := index + 1
	for ; currentIndex < len(lines); currentIndex++ {
		if strings.Contains(lines[currentIndex], "{{end}}") {
			break
		}
		slotContent.WriteString(lines[currentIndex] + "\n")
	}

	// Prepare combined data for included component
	newData := make(map[string]interface{})
	if dataMap, ok := data.(map[string]interface{}); ok {
		maps.Copy(newData, dataMap)
	}
	if propsMap, ok := propsData.(map[string]interface{}); ok {
		maps.Copy(newData, propsMap)
	}
	newData["slot"] = slotContent.String()

	// Render included component
	included, err := e.includeFunc(componentName, newData)
	if err != nil {
		return "", currentIndex, err
	}
	result.WriteString(included)

	return result.String(), currentIndex, nil
}

/*
=================================================================
Core Template Functions
=================================================================
*/

// componentFunc handles component rendering
// Example: {{component "Button" (props("variant" "primary"))}}
func (e *TemplateEngine) componentFunc(name string, data interface{}) (string, error) {
	return e.handleComponent(name, data, "component")
}

// includeFunc handles template inclusion
// Example: {{include "Header" (props("title" .PageTitle))}}
func (e *TemplateEngine) includeFunc(name string, data interface{}) (string, error) {
	return e.handleComponent(name, data, "include")
}

// slotFunc renders slot content
// Example: <div>{{slot}}</div>
func (e *TemplateEngine) slotFunc(args ...interface{}) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("slot requires data argument")
	}

	data, ok := args[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid slot data")
	}

	if content, exists := data["slot"]; exists {
		return fmt.Sprint(content), nil
	}
	return "", nil
}

// safeHTMLFunc marks HTML as safe for rendering
// Example: {{safeHTML("<strong>Bold</strong>")}}
func (e *TemplateEngine) safeHTMLFunc(html string) (string, error) {
	return html, nil
}

// Props creates a property map from key-value pairs
// Example: {{props("id" "main" "class" "container")}}
func Props(args ...interface{}) (map[string]interface{}, error) {
	if len(args)%2 != 0 {
		return nil, fmt.Errorf("props requires even number of arguments")
	}

	props := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			return nil, fmt.Errorf("prop key at position %d must be string", i)
		}
		props[key] = args[i+1]
	}
	return props, nil
}

func toJSONFunc(arg interface{}) (interface{}, error) {
	b, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

/*
=================================================================
Helper Functions and Utilities
=================================================================
*/

// handleComponent is the unified component handler
func (e *TemplateEngine) handleComponent(name string, data interface{}, funcType string) (string, error) {
	ctx, ok := data.(map[string]interface{})["_ctx"].(*RenderCtx)
	if !ok {
		return "", fmt.Errorf("missing render context")
	}

	// Get component content
	var component string
	var exists bool
	switch funcType {
	case "component", "include":
		component, exists = ctx.Site.Components[name]
	default:
		return "", fmt.Errorf("invalid component type: %s", funcType)
	}

	if !exists {
		return "", fmt.Errorf("%s %s not found", funcType, name)
	}

	// Create new engine instance for component
	engine := NewTemplateEngine()
	return engine.Execute(component, data)
}

// evaluateCondition handles boolean logic in templates
func (e *TemplateEngine) evaluateCondition(expr string, data interface{}) (bool, error) {
	if strings.Contains(expr, "eq(") {
		return e.handleEquality(expr, data)
	}

	value, err := e.evalExpression(expr, data)
	if err != nil {
		return false, err
	}

	if boolVal, ok := value.(bool); ok {
		return boolVal, nil
	}

	return false, fmt.Errorf("non-boolean condition: %v", value)
}

// handleEquality handles equality comparisons
func (e *TemplateEngine) handleEquality(expr string, data interface{}) (bool, error) {
	argsStr := expr[strings.Index(expr, "(")+1 : strings.LastIndex(expr, ")")]
	args, err := parseArgs(argsStr)
	if err != nil {
		return false, err
	}

	if len(args) != 2 {
		return false, fmt.Errorf("eq requires 2 arguments")
	}

	left, err := e.evalExpression(args[0], data)
	if err != nil {
		return false, err
	}

	right, err := e.evalExpression(args[1], data)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(left, right), nil
}

// evalExpression evaluates template expressions.
// If the expression contains parentheses, it is treated as a function call.
func (e *TemplateEngine) evalExpression(expr string, data interface{}) (interface{}, error) {
	expr = strings.TrimSpace(expr)

	if strings.Contains(expr, "(") {
		return e.evalFunctionCall(expr, data)
	}

	if strings.HasPrefix(expr, ".") {
		return e.evalDataAccess(expr, data)
	}

	return e.evalLiteral(expr)
}

// evalPipeline evaluates an expression with a pipeline operator.
// Each function in the pipeline must use the parentheses syntax.
func (e *TemplateEngine) evalPipeline(expr string, data interface{}) (string, error) {
	parts := strings.Split(expr, "|")
	// Evaluate the left-most expression.
	current, err := e.evalExpression(strings.TrimSpace(parts[0]), data)
	if err != nil {
		return "", err
	}
	// Process each subsequent function.
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		var fnName string
		var args []string
		if strings.Contains(part, "(") && strings.Contains(part, ")") {
			fnName = part[:strings.Index(part, "(")]
			argsStr := part[strings.Index(part, "(")+1 : strings.LastIndex(part, ")")]
			args, err = parseArgs(argsStr)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("function call %s must include parentheses", part)
		}
		fn, exists := e.funcs[fnName]
		if !exists {
			return "", fmt.Errorf("function %s not found", fnName)
		}
		// Build argument list: first argument is the current value.
		evaluatedArgs := []interface{}{current}
		for _, arg := range args {
			evalArg, err := e.evalExpression(arg, data)
			if err != nil {
				return "", err
			}
			evaluatedArgs = append(evaluatedArgs, evalArg)
		}
		fnValue := reflect.ValueOf(fn)
		in := make([]reflect.Value, len(evaluatedArgs))
		for i, a := range evaluatedArgs {
			in[i] = reflect.ValueOf(a)
		}
		results := fnValue.Call(in)
		if len(results) != 2 {
			return "", fmt.Errorf("function %s returned %d values, expected 2", fnName, len(results))
		}
		if errVal, ok := results[1].Interface().(error); ok && errVal != nil {
			return "", errVal
		}
		current = results[0].Interface()
	}
	return fmt.Sprint(current), nil
}

// evalDataAccess handles dot-notation data access.
func (e *TemplateEngine) evalDataAccess(expr string, data interface{}) (interface{}, error) {
	parts := strings.Split(expr[1:], ".")
	value := reflect.ValueOf(data)

	for _, part := range parts {
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		if value.Kind() != reflect.Struct {
			return nil, fmt.Errorf("invalid data access: %s", expr)
		}
		value = value.FieldByName(part)
		if !value.IsValid() {
			return nil, fmt.Errorf("field %s not found", part)
		}
	}
	return value.Interface(), nil
}

// evalLiteral handles literals and simple values.
func (e *TemplateEngine) evalLiteral(expr string) (interface{}, error) {
	if strings.HasPrefix(expr, `"`) && strings.HasSuffix(expr, `"`) {
		return strings.Trim(expr, `"`), nil
	}
	if val, err := strconv.Atoi(expr); err == nil {
		return val, nil
	}
	if val, err := strconv.ParseBool(expr); err == nil {
		return val, nil
	}
	return nil, fmt.Errorf("unhandled expression: %s", expr)
}

// parseArgs handles argument parsing for function calls.
// If no commas are present, arguments are split on whitespace.
func parseArgs(input string) ([]string, error) {
	if strings.Contains(input, ",") {
		var args []string
		var current bytes.Buffer
		inQuotes := false
		escape := false

		for _, r := range input {
			switch {
			case escape:
				current.WriteRune(r)
				escape = false
			case r == '\\':
				escape = true
			case r == '"':
				inQuotes = !inQuotes
				current.WriteRune(r)
			case r == ',' && !inQuotes:
				if current.Len() > 0 {
					args = append(args, strings.TrimSpace(current.String()))
					current.Reset()
				}
			default:
				current.WriteRune(r)
			}
		}
		if current.Len() > 0 {
			args = append(args, strings.TrimSpace(current.String()))
		}
		return args, nil
	}
	// Otherwise, split on whitespace while preserving quoted strings.
	var args []string
	var current bytes.Buffer
	inQuotes := false
	escape := false

	for _, r := range input {
		switch {
		case escape:
			current.WriteRune(r)
			escape = false
		case r == '\\':
			escape = true
		case r == '"':
			inQuotes = !inQuotes
			current.WriteRune(r)
		case !inQuotes && (r == ' ' || r == '\t' || r == '\n'):
			if current.Len() > 0 {
				args = append(args, strings.TrimSpace(current.String()))
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, strings.TrimSpace(current.String()))
	}
	return args, nil
}

/*
=================================================================
Missing Functions Implementation
=================================================================
*/

// processTag handles individual template tags
func (e *TemplateEngine) processTag(content string, data interface{}) (string, error) {
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return "", nil
	}

	switch parts[0] {
	case "if":
		return "", fmt.Errorf("if blocks should be handled in processIfBlock")
	case "include":
		return "", fmt.Errorf("include blocks should be handled in processIncludeBlock")
	case "slot":
		return e.handleSlot(data)
	default:
		return e.handleFunctionOrVariable(content, data)
	}
}

// evalFunctionCall evaluates function calls in templates.
// The expression must include parentheses.
func (e *TemplateEngine) evalFunctionCall(expr string, data interface{}) (interface{}, error) {
	fnName := expr[:strings.Index(expr, "(")]
	argsStr := expr[strings.Index(expr, "(")+1 : strings.LastIndex(expr, ")")]
	args, err := parseArgs(argsStr)
	if err != nil {
		return nil, err
	}
	evaluatedArgs := make([]interface{}, len(args))
	for i, arg := range args {
		evaluatedArg, err := e.evalExpression(arg, data)
		if err != nil {
			return nil, err
		}
		evaluatedArgs[i] = evaluatedArg
	}
	fn, exists := e.funcs[fnName]
	if !exists {
		return nil, fmt.Errorf("function %s not found", fnName)
	}
	fnValue := reflect.ValueOf(fn)
	in := make([]reflect.Value, len(evaluatedArgs))
	for i, arg := range evaluatedArgs {
		in[i] = reflect.ValueOf(arg)
	}
	results := fnValue.Call(in)
	if len(results) != 2 {
		return nil, fmt.Errorf("function %s returned %d values, expected 2", fnName, len(results))
	}
	if errVal, ok := results[1].Interface().(error); ok && errVal != nil {
		return nil, errVal
	}
	return results[0].Interface(), nil
}

// handleFunctionOrVariable processes function calls (which must include parentheses)
// or variable access/literal values.
func (e *TemplateEngine) handleFunctionOrVariable(expr string, data interface{}) (string, error) {
	// If the expression contains parentheses, treat it as a function call.
	if strings.Contains(expr, "(") {
		result, err := e.handleFunctionCall(expr, data)
		if err != nil {
			return "", err
		}
		return fmt.Sprint(result), nil
	}

	// Otherwise, treat as variable access.
	if strings.HasPrefix(expr, ".") {
		value, err := e.evalExpression(expr[1:], data)
		if err != nil {
			return "", err
		}
		return fmt.Sprint(value), nil
	}

	// Try literal values.
	if val, err := strconv.ParseBool(expr); err == nil {
		return strconv.FormatBool(val), nil
	}
	if _, err := strconv.Atoi(expr); err == nil {
		return expr, nil
	}
	if strings.HasPrefix(expr, `"`) && strings.HasSuffix(expr, `"`) {
		return strings.Trim(expr, `"`), nil
	}

	return "", fmt.Errorf("unhandled expression: %s", expr)
}

// handleSlot processes slot content.
func (e *TemplateEngine) handleSlot(data interface{}) (string, error) {
	if dataMap, ok := data.(map[string]interface{}); ok {
		if slot, exists := dataMap["slot"]; exists {
			return fmt.Sprint(slot), nil
		}
	}
	return "", nil
}

// handleFunctionCall processes function execution.
func (e *TemplateEngine) handleFunctionCall(expr string, data interface{}) (string, error) {
	result, err := e.evalFunctionCall(expr, data)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(result), nil
}

/*
=================================================================
Example Usage
=================================================================
*/

/*
func main() {
	// Initialize engine context
	ctx := &RenderCtx{
		Site: SiteStructure{
			Components: map[string]string{
				"Button": `<button class="{{.class}}">{{.text}}</button>`,
			},
		},
	}

	// Create template engine
	engine := NewTemplateEngine(ctx)

	// Register custom function
	engine.RegisterFunc("upper", func(args []interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("upper expects 1 argument")
		}
		return strings.ToUpper(fmt.Sprint(args[0])), nil
	})

	// Define template with multiple features
	template := `
	<div>
		{{include "Header" (props("title" "Demo Page"))}}
		<h1>{{upper(.title)}}</h1>
		{{if eq(.showButton "true")}}
			{{component "Button" (props("class" "primary" "text" "Click Me"))}}
		{{end}}
	</div>`

	// Execute template with data
	data := map[string]interface{}{
		"title":      "welcome",
		"showButton": "true",
		"_ctx":       ctx,
	}

	result, err := engine.Execute(template, data)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
*/
