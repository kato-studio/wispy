package template

import (
	"fmt"
	"strings"
)

// RenderCommand simply returns a formatted string showing its arguments.
func RenderCommand(ctx *RenderCtx, args ...any) string {
	return fmt.Sprintf("[render: %v]", args)
}

// ComponentCommand returns a block command function that renders a component.
func ComponentCommand(engine *TemplateEngine) func(ctx *RenderCtx, args ...any) string {
	return func(ctx *RenderCtx, args ...any) string {
		containerClass := ""
		if len(args) > 1 {
			if class, ok := args[0].(string); ok {
				containerClass = strings.Trim(class, "\"")
			}
		}
		body, _ := args[len(args)-1].([]Node)
		ctx.Internal["inComponent"] = true
		defer delete(ctx.Internal, "inComponent")
		renderedBody := engine.compileNodes(body)(ctx)
		if containerClass != "" {
			return fmt.Sprintf("<div class=\"%s\">%s</div>", containerClass, renderedBody)
		}
		return renderedBody
	}
}

// IfCommand returns a block command function that renders its body if the condition is true.
func IfCommand(engine *TemplateEngine) func(ctx *RenderCtx, args ...any) string {
	return func(ctx *RenderCtx, args ...any) string {
		ctx.Internal["inIf"] = true
		defer delete(ctx.Internal, "inIf")
		if len(args) < 2 {
			return ""
		}
		condExpr, ok := args[0].(string)
		if !ok {
			return ""
		}
		condTokens := tokenize(condExpr)
		condVal := fmt.Sprintf("%v", engine.evaluateExpression(condTokens, ctx))
		if condVal != "" && condVal != "false" {
			if body, ok := args[len(args)-1].([]Node); ok {
				return engine.compileNodes(body)(ctx)
			}
		}
		return ""
	}
}

// ForCommand returns a block command function that iterates over a collection.
func ForCommand(engine *TemplateEngine) func(ctx *RenderCtx, args ...any) string {
	return func(ctx *RenderCtx, args ...any) string {
		ctx.Internal["inFor"] = true
		defer delete(ctx.Internal, "inFor")
		if len(args) < 3 {
			return ""
		}
		loopVar, ok := args[0].(string)
		if !ok {
			return ""
		}
		inStr, ok := args[1].(string)
		if !ok || inStr != "in" {
			return ""
		}
		// Join all tokens from index 2 onward for the collection expression.
		collExpr := ""
		for i := 2; i < len(args); i++ {
			if s, ok := args[i].(string); ok {
				collExpr += s + " "
			}
		}
		collExpr = strings.TrimSpace(collExpr)
		collExpr = strings.TrimPrefix(collExpr, "(")
		collExpr = strings.TrimSuffix(collExpr, ")")
		collTokens := tokenize(collExpr)
		collVal := engine.evaluateExpression(collTokens, ctx)
		var items []string
		switch v := collVal.(type) {
		case []string:
			items = v
		case string:
			items = strings.Split(v, ",")
			for i, item := range items {
				items[i] = strings.TrimSpace(item)
			}
		default:
			items = []string{fmt.Sprintf("%v", v)}
		}
		var sb strings.Builder
		oldVal, hadOld := ctx.Internal[loopVar]
		for _, item := range items {
			ctx.Internal[loopVar] = item
			if body, ok := args[len(args)-1].([]Node); ok {
				sb.WriteString(engine.compileNodes(body)(ctx))
			}
		}
		if hadOld {
			ctx.Internal[loopVar] = oldVal
		} else {
			delete(ctx.Internal, loopVar)
		}
		return sb.String()
	}
}
