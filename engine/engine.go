package engine

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/kato-studio/wispy/internal"
)

// Html processes the raw HTML bytes and applies templating logic
// ctxId is used as a unique identifier
func (r *Render) Render(rawBytes []byte, data map[string]any, attrs map[string]string, children []*TreeNode, ctxId string) ([]byte, []error) {
	// fmt.Println("Html render; data:", data)
	var nodes []*TreeNode
	var errs []error
	var namedImportPath = make(map[string]string)
	var output bytes.Buffer

	// Step 1: Build the node tree from raw HTML
	nodes, errs = r.BuildNodeTree(rawBytes, ctxId, namedImportPath)

	// Component case!
	fmt.Println("<> ", ctxId)
	if len(ctxId) > len(ROOT_DIR) && ctxId[:len(ROOT_DIR)] == ROOT_DIR ||
		len(ctxId) > len(SHARED_DIR) && ctxId[:len(SHARED_DIR)] == SHARED_DIR {
		fmt.Println("ctxId", ctxId[:2])
		fmt.Println(ctxId)
		// 2. Create merged data context with component attributes
		mergedData := maps.Clone(data)
		for k, v := range attrs {
			mergedData[k] = v
		}

		// 3. Build quick slot lookup
		slots := make(map[string][]*TreeNode, 2) // Pre-allocate for common case (default + named slot)
		for _, directChild := range nodes {
			name := directChild.Attributes["slot"]
			if name == "" {
				name = "default"
			}
			slots[name] = append(slots[name], directChild)
		}

		for _, compNode := range children {
			r.processNode(&output, compNode, mergedData, slots, namedImportPath)
		}
		return output.Bytes(), errs
	}

	// Non-component
	for _, child := range nodes {
		r.processNode(&output, child, data, nil, namedImportPath)
	}
	return output.Bytes(), errs
}

// processNode processes each node, handling its content, attributes, and children
func (r *Render) processNode(output *bytes.Buffer, node *TreeNode, data map[string]any, slots map[string][]*TreeNode, namedImportPath map[string]string) []error {
	var errs []error

	switch node.Type {
	case "component":
		// Handle components
		compPath, compExists := namedImportPath[node.Name]
		rawBytes, impExists := r.Imports[compPath]
		if compExists && impExists {
			//
			bytes, compErrs := r.Render(rawBytes, data, node.Attributes, node.Children, compPath)
			internal.IfErrPush(&errs, compErrs...)
			output.Write(bytes)
		} else {
			// TODO only write string in "dev mode" ("dev mode" not yet implemented)
			output.WriteString("[Could not load \"" + node.Name + "\"]")
		}
		return nil
	case "operation":
		// Handle Imports (x:imports)
		if node.Name == "x:imports" {
			for name, rawImportPath := range node.Attributes {
				//
				var importPath = rawImportPath
				fmt.Println("r.Imports for ::", importPath)
				if _, exists := r.Imports[importPath]; exists {
					fmt.Printf("[%s] Already in instance cache \"r.Imports\"", name)
				} else {
					var hasSharedPrefix bool
					importPath, hasSharedPrefix = strings.CutPrefix(importPath, "@")
					if hasSharedPrefix {
						importPath = SHARED_DIR + importPath + EXT
					}
					var hasSitePrefix bool
					importPath, hasSitePrefix = strings.CutPrefix(importPath, "~")
					if hasSitePrefix {
						importPath = ROOT_DIR + "/" + r.Ctx.Site.Name + importPath + EXT
					}
					//
					rawBytes, err := os.ReadFile(importPath)
					internal.IfErrPush(&errs, err)
					if err == nil && (hasSharedPrefix || hasSitePrefix) {
						namedImportPath[name] = importPath
						r.Imports[importPath] = rawBytes
					}
				}
			}
		} else
		// Handle Slots (x:slots)
		if node.Name == "x:slot" {
			slotName := "default"
			if name, exists := node.Attributes["name"]; exists {
				slotName = name
			}

			passedChildren, slotExists := slots[slotName]
			if slotExists {
				// Slot has been passed process child elements
				for _, child := range passedChildren {
					childErrors := r.processNode(output, child, data, nil, namedImportPath)
					internal.IfErrPush(&errs, childErrors...)
				}
			} else {
				// Handle fallback content
				for _, child := range node.Children {
					childErrors := r.processNode(output, child, data, nil, namedImportPath)
					internal.IfErrPush(&errs, childErrors...)
				}
			}

			return nil
		}
		// Handle other operations...
		return nil
	case "element":
		// Process regular html elements
		isSelfClosing := slices.Contains(SELF_CLOSING_TAGS, node.Name)

		output.WriteString("<" + node.Name)
		for attr, value := range node.Attributes {
			processedAttrValue, err := r.processTemplateVariables(value, data)
			internal.IfErrPush(&errs, err)
			if attrFunc, exists := r.AttrFuncMap[attr]; exists {
				changed, newAttr, attrErrs := attrFunc(attr, value)
				internal.IfErrPush(&errs, attrErrs...)
				if changed {
					output.WriteString(fmt.Sprintf(` %s`, newAttr))
				} else {
					output.WriteString(fmt.Sprintf(` %s="%s"`, attr, processedAttrValue))
				}
			} else {
				output.WriteString(fmt.Sprintf(` %s="%s"`, attr, processedAttrValue))
			}
			output.WriteString(fmt.Sprintf(` %s="%s"`, attr, processedAttrValue))
		}
		output.WriteString(">")

		for _, child := range node.Children {
			nodeErrs := r.processNode(output, child, data, slots, namedImportPath)
			internal.IfErrPush(&errs, nodeErrs...)
		}

		if !isSelfClosing {
			output.WriteString(fmt.Sprintf("</%s>", node.Name))
		}
	case "text":
		processedContent, err := r.processTemplateVariables(node.Content, data)
		internal.IfErrPush(&errs, err)
		output.WriteString(processedContent)
	}

	return errs
}

// renderComponent handles component rendering with slot management
// func (r *Render) renderComponent(componentTree []*TreeNode, rawNode *TreeNode, data map[string]any) ([]byte, []error) {

// }

// processTemplateVariables replaces template placeholders (e.g., {{.variable}}) with actual values from r.ctx.Data
func (r *Render) processTemplateVariables(input string, data map[string]any) (string, error) {
	result := input
	for {
		start := strings.Index(result, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "}}")
		if end == -1 {
			return result[:start], errors.New("mismatched template placeholder syntax")
		}
		end = start + end // Adjust end to be relative to full string

		placeholder := result[start+2 : end]
		if placeholder[0] == '.' {
			placeholder = placeholder[1:]
			if val, found := data[placeholder]; found {
				result = result[:start] + fmt.Sprintf("%v", val) + result[end+2:]
			} else {
				result = result[:start] + result[end+2:]
			}
		} else {
			split := smartSplit([]byte(placeholder))
			if len(split) > 1 {
				result = result[:start] + callTemplateFunction(split) + result[end+2:]
			} else {
				result = result[:start] + result[end+2:]
			}
		}
	}

	return result, nil
}

// Dummy template function for now (can be replaced later with actual logic)
func callTemplateFunction(args ...interface{}) string {
	// This is a placeholder function; replace it with actual logic as needed
	fmt.Println("args passed to dummyFunction", args)
	return "[[Dummy Op]]"
}

// Builds a tree structure.
// ctxId is used to prevent the same JS/CSS from being hoisted multiple times.
// for components this is the components path (other sources TBD)
func (r *Render) BuildNodeTree(rawBytes []byte, ctxId string, namedImportPath map[string]string) ([]*TreeNode, []error) {
	var root TreeNode
	var stack []*TreeNode
	stack = append(stack, &root) // Start with a dummy root node
	var currentContent []byte
	var validContent = false
	var inTag = false
	var i = 0
	var errs []error
	var lineNumber = 1

	reset := func() {
		inTag = false
		validContent = false
		currentContent = nil
	}

	for ; i < len(rawBytes); i++ {

		b := rawBytes[i]

		// skip line break
		if b == '\n' || b == '\r' || (isWhitespace(b) && isWhitespace(rawBytes[i+1])) {
			if b == '\n' || b == '\r' {
				// "temp fix" for line number not sure why the count is off by 1 from vscode file
				// lineNumber is set 1 by default, so here we increment it to 2 if we're past the first line
				if lineNumber == 1 {
					lineNumber++
				}
				lineNumber++
			}
			continue
		}
		//
		if !isWhitespace(b) {
			validContent = true
		}

		// Core Logic
		// --------
		if b == '<' {
			if len(currentContent) > 0 && validContent {
				// A new tag/node has been found create a text node for any previously cached bytes
				textNode := &TreeNode{
					Type:    "text",
					Content: string(currentContent),
					// Content: string(bytes.TrimSuffix(currentContent, []byte{' '})),
				}
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, textNode)
			}
			reset()
			inTag = true

			// Handle closing comments
			// we need to handle this first so that if we're in a comment
			// we can continue without call other unnecessary checks
			if bytes.HasPrefix(rawBytes[i:], []byte("<!--")) {
				closingTag := []byte("-->")
				closingIndex := bytes.Index(rawBytes[i:], closingTag)
				if closingIndex == -1 {
					fmt.Println("[Error] could not find comment closing tag attempting to resume parsing (unsafe comment parsing not implemented)")
					fmt.Println("[Warn] skipping all content!")
					i = len(rawBytes)
				} else {
					// add comment node to tree
					commentNode := &TreeNode{
						Type: "comment",
						// "+4" used to exclude opening "<!--" bytes from comment contents
						Content: string(rawBytes[i+4 : i+closingIndex]),
					}
					parent := stack[len(stack)-1]
					parent.Children = append(parent.Children, commentNode)

					// "+3" ensures we do not have hyphens left over from comment written as content
					i += (closingIndex + len(closingTag))
					reset()
				}
			}

			// Handling style tags
			var stylePrefix = []byte("<style")
			if bytes.HasPrefix(rawBytes[i:], stylePrefix) {
				// Find the closing ">" of the <style> tag.
				tagEndIndex := bytes.Index(rawBytes[i:], []byte(">"))
				if tagEndIndex == -1 {
					fmt.Print("[Error]", "Malformed HTML; no closing for <style> tag.")
					break
				} else {
					tagEndIndex += i
				}
				// move index forward
				i += len(stylePrefix)
				// remove closing bracket
				styleTag := rawBytes[i:tagEndIndex]
				attributes := parseAttributes(styleTag)
				scopeNamespace, isScoped := attributes["scoped"]
				if scopeNamespace == "" && isScoped {
					scopeNamespace = strings.ReplaceAll(ctxId, "/", "-")
				}
				//
				closingTag := "</style>"
				closingIndex := bytes.Index(rawBytes[tagEndIndex:], []byte(closingTag))
				if closingIndex == -1 {
					fmt.Print("[Error]", "Malformed HTML; no closing </style>.")
					break
				} else {
					closingIndex += tagEndIndex
				}

				// Handle non-scoped style tag
				if !isScoped {
					commentNode := &TreeNode{
						Name: "style",
						Type: "element",
						// "+4" used to exclude opening "<!--" bytes from comment contents
						Content: string(rawBytes[i+4 : i+closingIndex]),
					}
					parent := stack[len(stack)-1]
					parent.Children = append(parent.Children, commentNode)

					i = closingIndex + len(closingTag)
					reset()
					continue
				}

				// Handle hoisting scoped styles
				_, exists := r.Css[ctxId]
				if !exists {
					// extract all styles within <style> tag
					styles := rawBytes[tagEndIndex+1 : closingIndex]
					// Scope the styles inside the <style> block.
					var selector = scopeNamespace
					scopedStyles := scopeStyles(selector, string(styles))
					r.Css[ctxId] = []byte(scopedStyles)
					// Append the processed styles and update the offset.
				}

				i = closingIndex + len(closingTag)

				reset()
			}

			continue
		}

		if b == '>' {
			if inTag {
				tagContent := currentContent
				reset()
				inTag = false

				if len(tagContent) > 0 {
					if tagContent[0] == '/' {
						// Closing tag
						// Handle closing tag does not have start tag
						if len(stack[:len(stack)-1]) == 0 {
							fmt.Printf("[Error] broken end tag line:[%d] <%s> \n", lineNumber, string(tagContent))
							continue
						} else {
							// Happy path
							stack = stack[:len(stack)-1]
						}
					} else {
						// Opening or Self-closing tag
						isSelfClosing := tagContent[len(tagContent)-1] == '/'
						tagContent = bytes.TrimSuffix([]byte(tagContent), []byte("/"))
						splitIndex := bytes.IndexByte(tagContent, ' ')
						tagName := ""
						if splitIndex == -1 {
							splitIndex = 0
							tagName = string(tagContent)
						} else {
							tagName = string(tagContent[:splitIndex])
						}

						// Node creation logic
						attributes := parseAttributes(tagContent[splitIndex:])
						tagType := detectTagType(tagName)

						newNode := &TreeNode{
							Name:       tagName,
							Type:       tagType,
							Attributes: attributes,
						}

						if len(stack) != 0 {
							parent := stack[len(stack)-1]
							parent.Children = append(parent.Children, newNode)
						} else {
							fmt.Println("[Error] Could not resolve closing tag ", string(tagContent))
						}
						// Only push non-self-closing tags onto the stack
						if !isSelfClosing && !slices.Contains(SELF_CLOSING_TAGS, tagName) {
							stack = append(stack, newNode)
						}
					}
				}
			}
			continue
		}
		// Don't write empty nodes
		if b != ' ' || (b == ' ' && len(currentContent) > 1) {
			currentContent = append(currentContent, b)
		}
	}

	return root.Children, errs
}
