package style

import (
	"fmt"
	"kato-studio/go-wispy/utils"
	"regexp"
	"strings"
)

// use regex
func ExtractClasses(htmlContent string) utils.UniqueSet {
	classRegex := regexp.MustCompile(`class="([^"]+)"`)
	matches := classRegex.FindAllStringSubmatch(htmlContent, -1)
	classes := utils.NewUniqueSet()
	for _, match := range matches {
		for _, class_name := range strings.Split(match[1], " ") {
			classes.Add(class_name)
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

func GenCssClass(raw_class_name, state_string, value string, category StyleCategory) string {
	escaped_class_name := EscapeClassName(raw_class_name, state_string)
	// Assumption is that this is a static class
	if category.Attr == "" && category.Format == "" {
		return fmt.Sprintf(".%s { %s }", escaped_class_name, value)
	}
	if category.Format != "" {
		if category.Options != nil {
			return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf(category.Format, category.Attr, category.Options[value]))
		}
		return fmt.Sprintf(".%s { %s }", escaped_class_name, value)
	}
	if category.Options != nil {
		return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: %s;", category.Attr, category.Options[value]))
	}
	return ""
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

func ResolveClass(raw_class_name, media_size string, Ctx StyleCTX) string {
	// Handle classes with declarative values
	var raw_value = "NULLISH"
	var working_class_name = raw_class_name
	//
	// HANDLE STATE PREFIXES
	// notes:
	// - this is not tested for all combinations
	// - and currently is only meant to handle single state prefix
	var state_string = ""
	if strings.Contains(raw_class_name, ":") {
		prefixes := strings.Split(raw_class_name, ":")
		var p_len = len(prefixes)
		if p_len > 1 {
			// Handle state prefixes
			fmt.Println("-------------")
			for i, prefix := range prefixes {
				if i == p_len-1 {
					// final element is the class name without state prefix(s)
					working_class_name = prefix
					break
				}
				fmt.Println("State Prefix: " + prefix)
				found_state := CLASS_STATES[prefix]
				if found_state != "" {
					state_string += ":" + found_state
				}
			}

			fmt.Println("State Class: " + raw_class_name)
			fmt.Println("State->: " + state_string)
			fmt.Println("-------------")
		}
	}

	// STATIC CLASS
	static_class := Ctx.StaticStyles[working_class_name]
	if static_class != "" {
		return GenCssClass(raw_class_name, state_string, static_class, StyleCategory{})
	}
	//
	if strings.Contains(working_class_name, "-") && len(working_class_name) > 1 {
		var split_index = strings.LastIndex(working_class_name, "-")
		raw_value = working_class_name[split_index+1:]
		var class_name = working_class_name[:split_index]
		var category, category_exists = styleCategories[class_name]
		//
		if category_exists {
			if category.Exclude.Contains(class_name) {
				return ""
			}
			//
			// HANDLE COLORS
			if category.IsColor {
				//
				escaped_class_name := EscapeClassName(media_size+raw_class_name, state_string)
				// Opacity
				split_class := strings.Split(raw_value, "/")
				value := split_class[0]
				opacity := ""
				if len(split_class) > 1 {
					opacity = split_class[1]
				}
				// ------------------
				// Default Color
				var sub_value = ""
				var current_sub_value = "500"
				sub_index := strings.LastIndex(value, "-")
				// if there is no sub value, then use 500 as default
				// example: .color-primary = .color-primary-500
				if sub_index != -1 {
					sub_value = value[sub_index+1:]
					current_sub_value = sub_value
				}
				//
				found_color := Ctx.Colors[value][current_sub_value]
				if found_color != "" {
					// Opacity
					if opacity != "" {
						hex_opacity_value := HEX_OPACITY[opacity]
						if hex_opacity_value != "" {
							found_color = found_color + hex_opacity_value
							variable_name := "--color-" + value + sub_value + "_" + opacity
							Ctx.AppendCssVariable(variable_name, found_color)
							return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: var(%s);", category.Attr, variable_name))
						} else {
							fmt.Println("Opacity not found: " + opacity)
						}
					} else {
						variable_name := "--color-" + value + sub_value
						Ctx.AppendCssVariable(variable_name, found_color)
						return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: var(%s);", category.Attr, variable_name))
					}
				}
			}
			// DYNAMIC CLASS
			return GenCssClass(media_size+raw_class_name, state_string, raw_value, category)
		} else {
			fmt.Println("[404]: " + raw_class_name)
		}
	}
	//
	return ""
}

// TODO🔰 Add Hover, Focus, Active, Disabled, etc. support
func WispyStyleGenerate(classes utils.UniqueSet, static_styles map[string]string, colors map[string]map[string]string) Styles {
	var output = Styles{
		CssVariables: map[string]string{},
		Base:         []string{},
		Sm:           []string{},
		Md:           []string{},
		Lg:           []string{},
		Xl:           []string{},
		_2xl:         []string{},
		_3xl:         []string{},
	}
	//
	for raw_class := range classes {
		// middleware function to validate result before appending
		ShouldAppend := func(dest *[]string, src string) {
			if src != "" {
				*dest = append(*dest, src)
			}
		}
		//
		AppendCssVariable := func(variable, value string) {
			if variable != "" {
				output.CssVariables[variable] = value
			}
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
			if media_size == "sm:" {
				ShouldAppend(&output.Sm, ResolveClass(class, media_size, CTX))
				continue
			} else if media_size == "md:" {
				ShouldAppend(&output.Md, ResolveClass(class, media_size, CTX))
				continue
			} else if media_size == "lg:" {
				ShouldAppend(&output.Lg, ResolveClass(class, media_size, CTX))
				continue
			} else if media_size == "xl:" {
				ShouldAppend(&output.Xl, ResolveClass(class, media_size, CTX))
				continue
			} else if media_size == "2xl:" {
				ShouldAppend(&output._2xl, ResolveClass(class, media_size, CTX))
				continue
			} else if media_size == "3xl:" {
				ShouldAppend(&output._3xl, ResolveClass(class, media_size, CTX))
				continue
			}
		}
		//
		ShouldAppend(&output.Base, ResolveClass(raw_class, "", CTX))
	}
	//
	return output
}

func WispyStyleCompile(input Styles) string {
	return fmt.Sprintf(`
		:root {
			%[1]s
		}
		%[2]s

		@media (min-width: 640px) {%[3]s}

		@media (min-width: 768px) {%[4]s}

		@media (min-width: 1024px) {%[5]s}

		@media (min-width: 1280px) {%[6]s}
	`,
		MapToCssVariables(input.CssVariables),
		strings.Join(input.Base, "\n"),
		strings.Join(input.Sm, "\n"),
		strings.Join(input.Md, "\n"),
		strings.Join(input.Lg, "\n"),
		strings.Join(input.Xl, "\n"),
	)
}

func MapToCssVariables(input map[string]string) string {
	var output = ""
	for variable, value := range input {
		output += fmt.Sprintf("%s: %s;\n", variable, value)
	}
	return output
}

// Example usage function
func DoThing() Styles {
	htmlContent := `
		<section>
			<!-- Container -->
			<div class="mx-auto w-full max-w-7xl px-5 py-16 md:px-10 md:py-20 color-red/50">
				<!-- Title -->
				<p class="text-center text-sm font-bold uppercase">3 easy steps</p>
				<h2 class="text-center text-3xl font-bold md:text-5xl hover:color-red"> How it works </h2>
				<p class="mx-auto mb-8 mt-4 max-w-lg text-center text-sm color-gray-500 sm:text-base md:mb-12 lg:mb-16"> Lorem ipsum dolor sit amet consectetur adipiscing elit ut aliquam,purus sit amet luctus magna fringilla urna </p>
				<!-- Content -->
				<div class="grid gap-5 sm:cols-2 md:cols-3 lg:gap-6">
					<!-- Item -->
					<div class="grid gap-4 rounded-md border border-solid border-gray-300 p-8 md:p-10">
						<div class="flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
							<p class="text-sm font-bold sm:text-xl">1</p>
						</div>
						<p class="text-xl font-semibold">Find Component</p>
						<p class="text-sm color-gray-500"> Lorem ipsum dolor sit amet consectetur adipiscing elit ut aliquam, purus sit. </p>
					</div>
					<!-- Item -->
					<div class="grid gap-4 rounded-md border border-solid border-gray-300 p-8 md:p-10">
						<div class="flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
							<p class="text-sm font-bold sm:text-xl">2</p>
						</div>
						<p class="text-xl font-semibold">Copy and Paste</p>
						<p class="text-sm color-gray-500"> Lorem ipsum dolor sit amet consectetur adipiscing elit ut aliquam, purus sit. </p>
					</div>
					<!-- Item -->
					<div class="grid gap-4 rounded-md border border-solid border-gray-300 p-8 md:p-10">
						<div class="flex h-12 w-12 items-center justify-center rounded-full bg-gray-100">
							<p class="text-sm font-bold sm:text-xl">3</p>
						</div>
						<p class="text-xl font-semibold">Done</p>
						<p class="text-sm color-gray-500"> Lorem ipsum dolor sit amet consectetur adipiscing elit ut aliquam, purus sit. </p>
					</div>
				</div>
			</div>
		</section>
		<div class="container mx-10">
			<div class="bg-blue-500 color-white p-4 ml-3 mr-6">Hello World</div>
			<div class="flex">
				<div class="bg-red-500 color-white size-6 p-4">Goodbye World</div>
				<div class="bg-red-500 color-white p-4">Goodbye World</div>
				<div class="bg-red-500 color-white size-6 p-4">Goodbye World</div>
			</div>
			<div class="blasd--ah_b__oop bsdasdlah blhasasd:123">Hello World</div>
		</div>
	`
	//
	classes := ExtractClasses(htmlContent)
	css := WispyStyleGenerate(classes, WispyStaticStyles, WispyColors)
	return css
}
