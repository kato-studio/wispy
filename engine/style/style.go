package style

import (
	"fmt"
	"kato-studio/go-wispy/utils"
	"regexp"
	"strings"
)

// use regex
func extractClasses(htmlContent string) utils.UniqueSet {
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

// Define supported style resolutions for responsive classes
type Styles struct {
	Colors utils.UniqueSet
	Base   []string
	Sm     []string
	Md     []string
	Lg     []string
	Xl     []string
	_2xl   []string
	_3xl   []string
}

func EscapeClassName(raw_class_name string) string {
	var escaped_class_name = strings.Replace(raw_class_name, ":", "\\:", 1)
	escaped_class_name = strings.ReplaceAll(escaped_class_name, ".", "\\.")
	escaped_class_name = strings.ReplaceAll(escaped_class_name, "]", "\\]")
	escaped_class_name = strings.ReplaceAll(escaped_class_name, "[", "\\[")
	//
	return escaped_class_name
}

func GenCssClass(class_name, value, format string, category StyleCategory) string {
	escaped_class_name := EscapeClassName(class_name)
	// Assumption is that this is a static class
	if category.Attr == "" && category.Format == "" {
		return fmt.Sprintf(".%s { %s }", escaped_class_name, value)
	}
	var attr = category.Attr
	if format != "" {
		value = fmt.Sprintf(format, value)
	}

	// handle color classes
	if strings.HasPrefix(class_name, "bg") || strings.HasPrefix(class_name, "text") {
		return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: var(--%s);", attr, value))
	}
	if category.Options != nil {
		return fmt.Sprintf(".%s { %s }", escaped_class_name, fmt.Sprintf("%s: %s;", attr, category.Options[value]))
	}
	return ""
}

func ResolveClass(raw_class_name, media_size string) string {
	// Handle classes with declarative values
	var value = "??null!!"

	// Check if the class is a static style class
	static_class := staticStyles[raw_class_name]
	if static_class != "" {
		return GenCssClass(raw_class_name, static_class, "", StyleCategory{})
	}

	if strings.Contains(raw_class_name, "-") && len(raw_class_name) > 1 {
		var split_index = strings.LastIndex(raw_class_name, "-")
		value = raw_class_name[split_index+1:]
		var class_name = raw_class_name[:split_index]
		var category, category_exists = styleCategories[class_name]

		if category_exists {
			if category.Exclude.Contains(class_name) {
				return ""
			}
			return GenCssClass(media_size+raw_class_name, value, "", category)
		} else {
			fmt.Println("[404]: " + raw_class_name)
		}
	}
	//
	return ""
}

// TODO🔰 Rename this function
// TODO🔰 Add Hover, Focus, Active, Disabled, etc. support
// TODO🔰 Add color support
// TODO🔰 Add opacity support to colors
func ResolveMediaClasses(classes utils.UniqueSet) Styles {
	var output = Styles{
		Colors: utils.UniqueSet{},
		Base:   []string{},
		Sm:     []string{},
		Md:     []string{},
		Lg:     []string{},
		Xl:     []string{},
		_2xl:   []string{},
		_3xl:   []string{},
	}

	for raw_class := range classes {
		// middleware function to validate result before appending
		ShouldAppend := func(dest *[]string, src string) {
			if src != "" {
				*dest = append(*dest, src)
			}
		}

		// Handle responsive prefixes
		split_index := strings.IndexRune(raw_class, ':')
		if split_index != -1 {
			class := raw_class[split_index+1:]
			media_size := raw_class[:split_index+1]
			fmt.Println("Media Size: ", media_size)
			if media_size == "sm:" {
				ShouldAppend(&output.Sm, ResolveClass(class, media_size))
				continue
			} else if media_size == "md:" {
				ShouldAppend(&output.Md, ResolveClass(class, media_size))
				continue
			} else if media_size == "lg:" {
				ShouldAppend(&output.Lg, ResolveClass(class, media_size))
				continue
			} else if media_size == "xl:" {
				ShouldAppend(&output.Xl, ResolveClass(class, media_size))
				continue
			} else if media_size == "2xl:" {
				ShouldAppend(&output._2xl, ResolveClass(class, media_size))
				continue
			} else if media_size == "3xl:" {
				ShouldAppend(&output._3xl, ResolveClass(class, media_size))
				continue
			}
		}

		ShouldAppend(&output.Base, ResolveClass(raw_class, ""))
	}

	return output
}

// Example usage function
func DoThing() Styles {
	htmlContent := `
		<section>
			<!-- Container -->
			<div class="mx-auto w-full max-w-7xl px-5 py-16 md:px-10 md:py-20">
				<!-- Title -->
				<p class="text-center text-sm font-bold uppercase">3 easy steps</p>
				<h2 class="text-center text-3xl font-bold md:text-5xl"> How it works </h2>
				<p class="mx-auto mb-8 mt-4 max-w-lg text-center text-sm color-gray-500 sm:text-base md:mb-12 lg:mb-16"> Lorem ipsum dolor sit amet consectetur adipiscing elit ut aliquam,purus sit amet luctus magna fringilla urna </p>
				<!-- Content -->
				<div class="grid gap-5 sm:grid-cols-2 md:grid-cols-3 lg:gap-6">
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

	classes := extractClasses(htmlContent)
	css := ResolveMediaClasses(classes)
	return css
}
