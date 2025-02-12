package atomicstyle

import (
	"fmt"
	"regexp"
	"strings"
	//

	styleInternal "github.com/kato-studio/wispy/atomicstyle/internal"
	dt "github.com/kato-studio/wispy/utils/datatypes"
)

// use regex
func RegexExtractClasses(htmlContent string) *dt.OrderedMap[string, struct{}] {
	classRegex := regexp.MustCompile(`class="([^"]+)"`)
	matches := classRegex.FindAllStringSubmatch(htmlContent, -1)
	classes := dt.NewOrderedMap[string, struct{}]()
	for _, match := range matches {
		for _, class_name := range strings.Split(match[1], " ") {
			classes.Set(class_name, struct{}{})
		}
	}
	return classes
}

func EscapeClassName(raw_class_name, state_string string) string {
	var escaped_class_name = raw_class_name
	var character_list = []string{":", ".", "]", "[", "/"}
	for _, character := range character_list {
		escaped_class_name = strings.ReplaceAll(escaped_class_name, character, "\\"+character)
	}
	return escaped_class_name + state_string
}

var CLASS_STATES = map[string]string{
	"hover":    "hover",
	"focus":    "focus",
	"active":   "active",
	"visited":  "visited",
	"disabled": "disabled",
	"first":    "first-child",
	"last":     "last-child",
	"odd":      "nth-child(odd)",
	"even":     "nth-child(even)",
	"after":    ":after",
	"before":   ":before",
}

func ResolveClass(raw_class_name, media_size string, Ctx StyleCTX) (value string, value_type string) {
	// Handle classes with declarative values
	var working_class_name = raw_class_name
	//
	// HANDLE STATE PREFIXES
	// Notes:
	// - this is not tested for all combinations
	// - and currently is only meant to handle single state prefix
	var state_string = ""
	if strings.Contains(raw_class_name, ":") {
		prefixes := strings.Split(raw_class_name, ":")
		var p_len = len(prefixes)
		if p_len > 1 {
			// Handle state prefixes
			for i, prefix := range prefixes {
				if i == p_len-1 {
					// final element is the class name without state prefix(s)
					working_class_name = prefix
					break
				}
				found_state := CLASS_STATES[prefix]
				if found_state != "" {
					state_string += ":" + found_state
				}
				// DEBUG
				if i == p_len-1 && found_state == "" {
					fmt.Println("[404] could not handle state: " + prefix)
				}
			}

		}
	}

	// Escape class name
	var escaped_class_name = EscapeClassName(media_size+raw_class_name, state_string)

	// STATIC CLASS
	var static_class = Ctx.StaticStyles[working_class_name]
	if static_class != "" {
		return fmt.Sprintf(".%s { %s }", escaped_class_name, static_class), "STATIC"
	}
	//
	// HANDLE DYNAMIC CLASSES
	if strings.Contains(working_class_name, "-") && len(working_class_name) > 1 {
		var IS_NEGATIVE = false
		if working_class_name[0] == '-' {
			working_class_name = working_class_name[1:]
			IS_NEGATIVE = true
		}
		var first_split = strings.IndexRune(working_class_name, '-')
		var second_split = strings.IndexRune(working_class_name[first_split+1:], '-')

		// base variables
		var color_shade = "500"
		var class_name = working_class_name[:first_split]
		var last_value = ""
		var mid_value = ""
		//
		if second_split != -1 {
			second_raw := working_class_name[first_split+1:]
			mid_value = second_raw[:second_split]
			last_value = second_raw[second_split+1:]
		} else {
			mid_value = ""
			last_value = working_class_name[first_split+1:]
		}
		//
		var category, category_exists = styleInternal.StyleCategories[class_name]
		if !category_exists {
			// Attempt to find category including the mid_value
			category, category_exists = styleInternal.StyleCategories[class_name+"-"+mid_value]
		}

		//
		// IF "CATEGORY" (predefined class options exists)
		if category_exists {
			// ---------------------
			// HANDLE COLORS
			if category.IsColor {
				//
				// Extract Opacity if opacity exists it is the last part of the class name
				// example: bg-red/50, text-red-500/50
				split_class := strings.Split(last_value, "/")
				// Check to enure color_name has been set improperly
				var color_name string
				if len(mid_value) != 0 {
					color_name = mid_value
				} else {
					color_name = last_value
				}
				opacity := ""
				if len(split_class) > 1 {
					opacity = split_class[1]
					color_shade = split_class[0]
				}

				// Find Color Shade if it exists in the color map which is passed in the context
				var found_color string
				var shade_name = last_value
				var ok = false
				found_color, ok = Ctx.Colors[color_name][last_value]
				if !ok {
					found_color = Ctx.Colors[color_name][color_shade]
					shade_name = color_shade
				}

				if found_color != "" {
					// Opacity
					if opacity != "" {
						hex_opacity_value := styleInternal.HEX_OPACITY[opacity]
						if hex_opacity_value != "" {
							color_value := found_color + hex_opacity_value
							variable_name := "--color-" + color_name + "-" + shade_name + "_" + opacity
							Ctx.AppendCssVariable(variable_name, color_value)
							return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: var(%s);", category.ColorAttr, variable_name)), "COLOR"
						} else {
							// Debug
							fmt.Println("Opacity not found: " + opacity)
						}
					} else {
						variable_name := "--color-" + color_name + "-" + shade_name
						Ctx.AppendCssVariable(variable_name, found_color)
						return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: var(%s);", category.ColorAttr, variable_name)), "COLOR"
					}
				}
			}

			// ---------------------
			// ARTIFICIAL CLASSES
			// example: .w-[100px] = width: 100px;
			if last_value[0] == '[' {
				// extract value
				artificial_value := last_value[1 : len(last_value)-1]
				artificial_value = strings.ReplaceAll(artificial_value, "_", " ")
				return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: %s;", category.Attr, artificial_value)), "ARTIFICIAL"
			}

			// DYNAMIC FLOW
			if category.Options != nil && category.Options[last_value] != "" {
				last_value = category.Options[last_value]
			}
			if last_value == "" && category.Format == "" {
				return "", ""
			}
			if IS_NEGATIVE {
				last_value = "-" + last_value
			}
			if category.Format != "" {
				return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf(category.Format, category.Attr, last_value)), "DYNAMIC"
			}
			return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: %s;", category.Attr, last_value)), "DYNAMIC"

		} else {
			// Debug
			// fmt.Println("[404]: " + class_name + " (" + raw_class_name + ")")
		}
	}
	//
	return "", ""
}

func WispyStyleGenerate(classes *dt.OrderedMap[string, struct{}], static_styles map[string]string, colors map[string]map[string]string) Styles {
	var output = Styles{
		CssVariables: map[string]string{},
		Static:       dt.NewOrderedMap[string, struct{}](),
		Base:         dt.NewOrderedMap[string, struct{}](),
		Sm:           dt.NewOrderedMap[string, struct{}](),
		Md:           dt.NewOrderedMap[string, struct{}](),
		Lg:           dt.NewOrderedMap[string, struct{}](),
		Xl:           dt.NewOrderedMap[string, struct{}](),
		_2xl:         dt.NewOrderedMap[string, struct{}](),
		_3xl:         dt.NewOrderedMap[string, struct{}](),
	}
	//
	for _, raw_class := range classes.Keys() {
		//
		AppendCssVariable := func(variable, value string) {
			if variable != "" {
				output.CssVariables[variable] = value
			}
		}
		// middleware function to validate result before appending
		ShouldAppend := func(dest *dt.OrderedMap[string, struct{}], value, value_type string) {
			if value == "" {
				return
			}
			if value_type == "STATIC" {
				output.Static.Set(value, struct{}{})
				return
			}
			dest.Set(value, struct{}{})
		}
		//
		CTX := StyleCTX{
			Colors:            colors,
			StaticStyles:      static_styles,
			AppendCssVariable: AppendCssVariable,
		}
		// Handle responsive prefixes
		split_index := strings.IndexRune(raw_class, ':')
		if split_index != -1 {
			class := raw_class[split_index+1:]
			media_size := raw_class[:split_index+1]
			clean_size := raw_class[:split_index]
			if clean_size == "sm" {
				value, value_type := ResolveClass(class, media_size, CTX)
				ShouldAppend(output.Sm, value, value_type)
				continue
			} else if clean_size == "md" {
				value, value_type := ResolveClass(class, media_size, CTX)
				ShouldAppend(output.Md, value, value_type)
				continue
			} else if clean_size == "lg" {
				value, value_type := ResolveClass(class, media_size, CTX)
				ShouldAppend(output.Lg, value, value_type)
				continue
			} else if clean_size == "xl" {
				value, value_type := ResolveClass(class, media_size, CTX)
				ShouldAppend(output.Xl, value, value_type)
				continue
			} else if clean_size == "2xl" {
				value, value_type := ResolveClass(class, media_size, CTX)
				ShouldAppend(output._2xl, value, value_type)
				continue
			} else if clean_size == "3xl" {
				value, value_type := ResolveClass(class, media_size, CTX)
				ShouldAppend(output._3xl, value, value_type)
				continue
			}
		}
		//
		value, value_type := ResolveClass(raw_class, "", CTX)
		ShouldAppend(output.Base, value, value_type)
	}
	//
	return output
}

func WispyStyleCompile(input Styles) string {
	// TODO: don't hardcode media queries if not needed
	return fmt.Sprintf(`
:root {%[1]s} %[2]s  %[3]s
@media (min-width: 640px) {%[4]s}
@media (min-width: 768px) {%[5]s}
@media (min-width: 1024px) {%[6]s}
@media (min-width: 1280px) {%[7]s}`,
		MapToCssVariables(input.CssVariables),
		// We need to ensure static classes are always at the top
		strings.Join(input.Static.Keys(), "\n"),
		strings.Join(input.Base.Keys(), "\n"),
		strings.Join(input.Sm.Keys(), "\n"),
		strings.Join(input.Md.Keys(), "\n"),
		strings.Join(input.Lg.Keys(), "\n"),
		strings.Join(input.Xl.Keys(), "\n"),
		strings.Join(input._2xl.Keys(), "\n"),
		strings.Join(input._3xl.Keys(), "\n"),
	)
}

func MapToCssVariables(input map[string]string) string {
	var output = ""
	for variable, value := range input {
		output += fmt.Sprintf("%s: %s;\n", variable, value)
	}
	return output
}
