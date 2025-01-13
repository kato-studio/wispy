package engine

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/kato-studio/wispy/internal"
)

// Html processes the raw HTML bytes and applies templating logic
func (r *Render) Render(rawBytes []byte, data map[string]any) ([]byte, []error) {
	fmt.Println("Html render called with data", data)
	var nodes []*TreeNode
	var errs []error

	// Step 1: Build the node tree from raw HTML
	imports := make(map[string][]byte)
	nodes, errs = r.BuildNodeTree(rawBytes, imports)
	if len(errs) > 0 {
		return nil, errs
	}

	// Step 2: Process each node in the tree, replacing template variables
	var output bytes.Buffer
	for _, node := range nodes {
		errs := r.processNode(&output, node, data, nil)
		if len(errs) > 0 {
			return nil, errs
		}
	}

	// Step 3: Return the processed output
	return output.Bytes(), nil
}

// processNode processes each node, handling its content, attributes, and children
func (r *Render) processNode(output *bytes.Buffer, node *TreeNode, data map[string]any, slots map[string][]*TreeNode) []error {
	var errs []error

	switch node.Type {
	case "component":
		// Handle components
		rawBytes, err := r.GetComponent(node.Name)
		if err == nil {
			imports := make(map[string][]byte)
			componentTree, errs := r.BuildNodeTree(rawBytes, imports)
			if len(errs) > 0 {
				return errs
			}
			if len(componentTree) > 0 {
				bytes, errs := r.renderComponent(componentTree, node, data)
				//= wondering if this should exit or if engine should always resolve and pass any errors/warnings along?
				if len(errs) > 0 {
					return errs
				}
				output.Write(bytes)
			} else {
				output.WriteString("[Could not build node tree for \"" + node.Name + "\"]")
			}
		} else {
			output.WriteString("[Could not load \"" + node.Name + "\"]")
		}
		return nil
	case "operation":
		// Handle operations (including slots)
		if node.Name == "x:slot" {
			slotName := "default"
			if name, exists := node.Attributes["name"]; exists {
				slotName = name
			}

			passedChildren, slotExists := slots[slotName]
			if slotExists {
				// Slot has been passed process child elements
				for _, child := range passedChildren {
					r.processNode(output, child, data, nil)
				}
			} else {
				// Handle fallback content
				for _, child := range node.Children {
					r.processNode(output, child, data, nil)
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
			nodeErrs := r.processNode(output, child, data, slots)
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
func (r *Render) renderComponent(componentTree []*TreeNode, rawNode *TreeNode, data map[string]any) ([]byte, []error) {
	var output bytes.Buffer
	var errs []error
	var children = rawNode.Children

	// 1. Create merged data context with component attributes
	mergedData := make(map[string]any, len(data)+len(rawNode.Attributes))
	for k, v := range data {
		mergedData[k] = v
	}
	for k, v := range rawNode.Attributes {
		mergedData[k] = v
	}

	// 2. Build quick slot lookup
	slots := make(map[string][]*TreeNode, 2) // Pre-allocate for common case (default + named slot)
	for _, child := range children {
		name := child.Attributes["slot"]
		if name == "" {
			name = "default"
		}
		slots[name] = append(slots[name], child)
	}

	for _, compNode := range componentTree {
		r.processNode(&output, compNode, mergedData, slots)
	}

	return output.Bytes(), errs
}

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
func (r *Render) BuildNodeTree(rawBytes []byte, imports map[string][]byte) ([]*TreeNode, []error) {
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
						attributes := parseAttributes(tagContent[splitIndex:])
						tagType := detectTagType(tagName)

						// Handle imports
						if tagName == "x:imports" {
							for newPath, attr := range attributes {
								componentPath := attributes[attr]
								if _, hasPrefix := strings.CutPrefix(componentPath, "@"); hasPrefix {
									componentPath = SHARED_DIR + "/" + newPath + EXT
									fmt.Println("### " + componentPath)
								} else if newPath, hasPrefix := strings.CutPrefix(componentPath, "~"); hasPrefix {
									componentPath = ROOT_DIR + "/" + r.Ctx.Site.Name + "/" + newPath + EXT
									fmt.Println("### " + componentPath)
								}
								//
								rawBytes, err := os.ReadFile(componentPath)
								internal.IfErrPush(&errs, err)
								if err != nil {
									imports[attr] = rawBytes
								} else {
									imports[attr] = []byte{}
								}
							}
						}

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
