// here is my golang definitions file that is referenced to generate a trie of class names and their css properties, let's make small changes to use css variables that exist in the theme where it makes sense to do so ```
package atomicstyle

import (
	"fmt"
)

// --- Data Definitions ---
var (
	// Todo: review color usage
	colorNames = []string{
		// Tailwind default colors compat
		"slate", "gray", "zinc", "neutral", "stone",
		"red", "orange", "amber", "yellow", "lime",
		"green", "emerald", "teal", "cyan", "sky",
		"blue", "indigo", "violet", "purple", "fuchsia",
		"pink", "rose",
		// brand/project colors
		"primary", "accent",
	}
	opacityValues = []string{
		"", "0", "5", "10", "20", "25", "30", "40", "50", "60", "70", "75", "80", "90", "95", "100",
	}
	shades = []string{
		"50", "100", "200", "300", "400", "500", "600", "700", "800", "900", "950",
	}
	percentages = map[string]string{
		"1/2": "50%",
		"2/5": "40%",
		"3/5": "60%",
		"1/3": "33.3333%",
		"2/3": "66.6667%",
		"1/4": "25%",
		"3/4": "75%",
		"1/5": "20%",
		"4/5": "80%",
		"1/6": "16.6667%",
		"5/6": "83.3333%",
	}
	borderWidths = map[string]string{
		"px": "1",
		"1":  "1",
		"2":  "2",
		"3":  "3",
		"4":  "4",
		"5":  "5",
		"8":  "8",
	}
	borderStyles = []string{
		"solid", "dashed", "dotted", "double", "none",
	}
	transformations = map[string]string{
		"uppercase":   "uppercase",
		"lowercase":   "lowercase",
		"capitalize":  "capitalize",
		"normal-case": "none",
	}
	decorations = map[string]string{
		"underline":    "underline",
		"line-through": "line-through",
		"no-underline": "none",
	}
	sizingValues = map[string]string{
		"px":   "1px",
		"auto": "auto",
		"full": "100%",
		"w":    "100vw",
		"h":    "100vh",
		"min":  "min-content",
		"max":  "max-content",
		"fit":  "fit-content",
	}
	containerBreakpoints = []string{
		"3xs", "2xs", "xs", "sm", "md", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl", "7xl",
	}
	textSizes = []string{
		"xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl", "7xl", "8xl", "9xl",
	}
	roundedSizes = []string{
		"none", "sm", "md", "lg", "xl", "2xl", "3xl", "full",
	}
	displayValues = map[string]string{
		"block":        "block",
		"inline":       "inline",
		"inline-block": "inline-block",
		"flex":         "flex",
		"inline-flex":  "inline-flex",
		"grid":         "grid",
		"inline-grid":  "inline-grid",
		"hidden":       "none",
	}
	clearValues = []string{
		"left", "right", "both", "none",
	}
	objectFitValues = map[string]string{
		"contain":    "contain;",
		"cover":      "cover;",
		"fill":       "fill;",
		"none":       "none;",
		"scale-down": "scale-down;",
	}
	translateValues = []string{"0", "1", "2", "3", "4", "5", "6", "8", "10", "12", "16", "20", "24", "32", "40", "48", "56", "64", "px", "full"}
	boxShadowSizes  = []string{"2xs", "xs", "sm", "md", "lg", "xl", "2xl", "inner", "none"}
	mixBlendModes   = []string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"}
	objectPosValues = map[string]string{
		"bottom":       "bottom",
		"center":       "center",
		"left":         "left",
		"left-bottom":  "left bottom",
		"left-top":     "left top",
		"right":        "right",
		"right-bottom": "right bottom",
		"right-top":    "right top",
		"top":          "top",
	}
	overflowValues = map[string]string{
		"auto":    "auto",
		"hidden":  "hidden",
		"visible": "visible",
		"scroll":  "scroll",
		"clip":    "clip",
	}
	blurValues = []string{"none", "sm", "md", "lg", "xl", "2xl", "3xl"}
	flexDirs   = map[string]string{
		"row":         "row",
		"row-reverse": "row-reverse",
		"col":         "column",
		"col-reverse": "column-reverse",
	}
	flexWraps = map[string]string{
		"wrap":         "wrap",
		"wrap-reverse": "wrap-reverse",
		"nowrap":       "nowrap;",
	}
	alignValues = map[string]string{
		"start":    "flex-start",
		"end":      "flex-end",
		"center":   "center",
		"baseline": "baseline",
		"stretch":  "stretch",
		"between":  "space-between",
	}
	flexShort = map[string]string{
		"1":       "1 1 0%;",
		"auto":    "1 1 auto;",
		"initial": "0 1 auto;",
		"none":    "none;",
	}
	fontWeights = map[string]string{
		"thin":       "100",
		"extralight": "200",
		"light":      "300",
		"normal":     "400",
		"medium":     "500",
		"semibold":   "600",
		"bold":       "700",
		"extrabold":  "800",
		"black":      "900",
	}
	// Letter Spacing (Tracking)
	tracking = []string{
		"tighter", "tight", "normal", "wide", "wider", "widest",
	}
	// Line height (Leading)
	lineHeights = []string{
		"tight", "snug", "normal", "relaxed", "loose", "loose",
	}
	positions   = []string{"static", "fixed", "absolute", "relative", "sticky"}
	insetValues = []string{"0", "auto"}
)

// BuildExtendedTrie builds a trie preloaded with all of our utility CSS classes.
func BuildFullTrie() *Trie {
	trie := NewTrie()
	manualInserts(trie)
	addLayout(trie)
	addFlexGrid(trie)
	addSpacing(trie)
	addSizing(trie)
	addTypography(trie)
	addBackgrounds(trie)
	addBorders(trie)
	addRingUtils(trie)
	addDivideUtils(trie)
	addEffects(trie)
	addFilters(trie)
	addTables(trie)
	addTransitions(trie)
	addTransforms(trie)
	addInteractivity(trie)
	addSVG(trie)
	addAccessibility(trie)
	addAspectRatio(trie)
	addScrollSnap(trie)
	addGridAutoFlow(trie)
	addGridUtilities(trie)
	addPlaceholderStyling(trie)
	addAdvancedUtilities(trie)
	addContainers(trie)
	addOutlineUtils(trie)

	return trie
}

// mostly made for cases where creating logic for a small group of classes is not worth it
func manualInserts(trie *Trie) {

}

// --- Layout Utilities ---
func addLayout(trie *Trie) {
	// Container and box-sizing
	trie.Insert("box-border", "box-sizing: border-box;")
	trie.Insert("box-content", "box-sizing: content-box;")
	// Add general overflow utilities
	for name, value := range overflowValues {
		trie.Insert("overflow-"+name, "overflow: "+value+";")
	}
	// Add overflow utilities for X and Y axes
	for name, value := range overflowValues {
		trie.Insert("overflow-x-"+name, "overflow-x: "+value+";")
		trie.Insert("overflow-y-"+name, "overflow-y: "+value+";")
	}
	//
	for k, rule := range displayValues {
		trie.Insert(k, "display: "+rule+";")
	}
	// Float and clear
	for _, dir := range []string{"right", "left", "none"} {
		trie.Insert("float-"+dir, "float: "+dir+";")
	}
	for _, v := range clearValues {
		trie.Insert("clear-"+v, "clear: "+v+";")
	}
	for k, rule := range objectFitValues {
		trie.Insert("object-"+k, "object-fit: "+rule+";")
	}
	for k, rule := range objectPosValues {
		trie.Insert("object-"+k, "object-position: "+rule+";")
	}
	// Positioning utilities
	for _, pos := range positions {
		trie.Insert(pos, "position: "+pos+";")
	}
	// Inset, top/right/bottom/left
	for _, val := range insetValues {
		trie.Insert("inset-"+val, fmt.Sprintf("top: %s; right: %s; bottom: %s; left: %s;", val, val, val, val))
		trie.Insert("inset-x-"+val, fmt.Sprintf("left: %s; right: %s;", val, val))
		trie.Insert("inset-y-"+val, fmt.Sprintf("top: %s; bottom: %s;", val, val))
	}
	// Visibility and z-index
	trie.Insert("visible", "visibility: visible;")
	trie.Insert("invisible", "visibility: hidden;")
	for _, z := range []string{"0", "10", "20", "30", "40", "50", "auto"} {
		trie.Insert("z-"+z, "z-index: "+z+";")
	}
}
func addContainers(trie *Trie) {
	// TODO:
	// Dynamic containers child-width utilities with CSS variables for flexibility
}

// --- Flexbox and Grid Utilities ---
func addFlexGrid(trie *Trie) {
	// Flex direction and wrap
	for k, rule := range flexDirs {
		trie.Insert("flex-"+k, "flex-direction: "+rule)
		trie.Insert(k, "flex-direction: "+rule)
	}
	for k, rule := range flexWraps {
		trie.Insert("flex-"+k, "flex-wrap: "+rule)
		trie.Insert(k, "flex-wrap: "+rule)
	}
	// Alignment and justify utilities
	trie.Insert("justify-left", "justify-content: left;")
	trie.Insert("justify-right", "justify-content: right;")
	for k, v := range alignValues {
		trie.Insert("items-"+k, "align-items: "+v+";")
		trie.Insert("content-"+k, "align-content: "+v+";")
		trie.Insert("self-"+k, "align-self: "+v+";")
		trie.Insert("justify-"+k, "justify-content: "+v+";")
	}
	for k, rule := range flexShort {
		trie.Insert("flex-"+k, "flex: "+rule)
	}
	trie.Insert("flex-grow", "flex-grow: 1;")
	trie.Insert("grow", "flex-grow: 1;")
	trie.Insert("grow-0", "flex-grow: 0;")
	trie.Insert("shrink", "flex-shrink: 1;")
	trie.Insert("shrink-0", "flex-shrink: 0;")
	// Order utilities
	trie.Insert("order-first", "order: -9999;")
	trie.Insert("order-last", "order: 9999;")
	trie.Insert("order-none", "order: 0;")
	for i := 1; i <= 12; i++ {
		order := itoa(i)
		trie.Insert("order-"+order, "order: "+order+";")
	}
	// Grid columns/rows and gap
	for i := 1; i <= 12; i++ {
		colClass := "grid-cols-" + itoa(i)
		trie.Insert(colClass, fmt.Sprintf("grid-template-columns: repeat(%d, minmax(0, 1fr));", i))
		trie.Insert("col-start-"+itoa(i), "grid-column-start: "+itoa(i)+";")
		trie.Insert("col-end-"+itoa(i), "grid-column-end: "+itoa(i)+";")
	}
	for i := 1; i <= 6; i++ {
		rowClass := "grid-rows-" + itoa(i)
		trie.Insert(rowClass, fmt.Sprintf("grid-template-rows: repeat(%d, minmax(0, 1fr));", i))
		trie.Insert("row-start-"+itoa(i), "grid-row-start: "+itoa(i)+";")
		trie.Insert("row-end-"+itoa(i), "grid-row-end: "+itoa(i)+";")
	}

	// Justify-items/self and place-* utilities (sample)
	trie.Insert("justify-items-center", "justify-items: center;")
	trie.Insert("place-content-center", "place-content: center;")
	trie.Insert("place-items-center", "place-items: center;")
	trie.Insert("place-self-center", "place-self: center;")
}

func addGridUtilities(trie *Trie) {

	//
	for i := 0; i <= 12; i++ {
		trie.Insert(fmt.Sprintf("span-%d", i), fmt.Sprintf("grid-column: span %d", i))
	}
	// Special values
	trie.Insert("row-span-full", "grid-row: 1 / -1;")
	trie.Insert("row-auto", "grid-row: auto;")
	trie.Insert("row-start-auto", "grid-row-start: auto;")
	trie.Insert("row-end-auto", "grid-row-end: auto;")
}

// --- Grid Auto Flow & Auto Columns/Rows ---
func addGridAutoFlow(trie *Trie) {
	trie.Insert("grid-flow-row", "grid-auto-flow: row;")
	trie.Insert("grid-flow-col", "grid-auto-flow: column;")
	trie.Insert("grid-flow-row-dense", "grid-auto-flow: row dense;")
	trie.Insert("grid-flow-col-dense", "grid-auto-flow: column dense;")
	trie.Insert("auto-cols-auto", "grid-auto-columns: auto;")
	trie.Insert("auto-cols-min", "grid-auto-columns: min-content;")
	trie.Insert("auto-cols-max", "grid-auto-columns: max-content;")
	trie.Insert("auto-cols-fr", "grid-auto-columns: minmax(0, 1fr);")
	trie.Insert("auto-rows-auto", "grid-auto-rows: auto;")
	trie.Insert("auto-rows-min", "grid-auto-rows: min-content;")
	trie.Insert("auto-rows-max", "grid-auto-rows: max-content;")
	trie.Insert("auto-rows-fr", "grid-auto-rows: minmax(0, 1fr);")
}

func addSpacing(trie *Trie) {
	// Auto margins
	trie.Insert("mx-auto", "margin-inline: auto;")
	// Support `space-x-reverse` and `space-y-reverse` (for RTL handling)
	trie.Insert("space-x-reverse", "direction: rtl;")
	trie.Insert("space-y-reverse", "direction: rtl;")
}

// --- Sizing Utilities ---
func addSizing(trie *Trie) {
	trie.Insert("w-screen", "width: 100vw;")
	trie.Insert("min-w-screen", "min-width: 100vw;")
	trie.Insert("max-w-screen", "max-width: 100vw;")
	trie.Insert("h-screen", "height: 100vh;")
	trie.Insert("min-h-screen", "min-height: 100vh;")
	trie.Insert("max-h-screen", "max-height: 100vh;")

	// Sizing
	for k, val := range sizingValues {
		trie.Insert("w-"+k, "width: "+val+";")
		trie.Insert("h-"+k, "height: "+val+";")
		trie.Insert("size-"+k, "height: "+val+"; "+"width: "+val+";")
		//
		trie.Insert("min-w-"+k, "min-width: "+val+";")
		trie.Insert("max-w-"+k, "max-width: "+val+";")
		trie.Insert("min-h-"+k, "min-height: "+val+";")
		trie.Insert("max-h-"+k, "max-height: "+val+";")
	}
	// Container Breakpoints
	for _, val := range containerBreakpoints {
		// May not be needed commenting for now
		// trie.Insert("w-"+val, "width: var(--breakpoint-"+val+");")
		// trie.Insert("h-"+val, "height: var(--breakpoint-"+val+");")
		// trie.Insert("size-"+val, "height: var--breakpoint-"+val+"); "+"width: var(--breakpoint-"+val+");")
		//
		trie.Insert("min-w-screen-"+val, "min-width: var(--container-"+val+");")
		trie.Insert("max-w-screen-"+val, "max-width: var(--container-"+val+");")
	}
}

// --- Typography Utilities ---
func addTypography(trie *Trie) {
	// Font families
	for _, f := range []string{"sans", "serif", "mono"} {
		trie.Insert("font-"+f, "font-family: "+f+";")
	}
	for _, val := range textSizes {
		trie.Insert("text-"+val, "font-size: var(--text-"+val+"); line-height: var(--leading, var(--text-"+val+"--line-height));")
	}
	// Text colors
	// for _, opacity := range opacityValues {
	trie.Insert("text-white", "color: var(--color-white);")
	trie.Insert("text-black", "color: var(--color-black);")
	for _, color := range colorNames {
		class := "text-" + color
		trie.Insert(class, "color: var(--color-"+color+"-500);")
		for _, shade := range shades {
			class := class + "-" + shade
			trie.Insert(class, "color: "+toColorVar(color+"-"+shade, "")+";")
		}
	}
	// }
	for k, w := range fontWeights {
		trie.Insert("font-"+k, "font-weight: "+w+";")
	}
	// Tracking (Letter spacing)
	for _, val := range tracking {
		trie.Insert("tracking-"+val, "letter-spacing: var(--tracking-"+val+");")
	}
	for _, val := range lineHeights {
		trie.Insert("leading-"+val, "line-height: var(--leading-"+val+");")
	}
	// Text align
	for _, a := range []string{"left", "center", "right", "justify"} {
		trie.Insert("text-"+a, "text-align: "+a+";")
	}
	for _, v := range borderWidths {
		trie.Insert("underline-"+v, "text-decoration-thickness: "+v+"px;")
		trie.Insert("underline-offset-"+v, "text-underline-offset: "+v+"px;")
	}
	// Text transformations
	for key, value := range transformations {
		trie.Insert(key, "text-transform: "+value+";")
	}
	// Text decorations
	for key, value := range decorations {
		trie.Insert(key, "text-decoration: "+value+";")
	}
	// Truncation & Overflow
	trie.Insert("truncate", "overflow: hidden; text-overflow: ellipsis; white-space: nowrap;")
	trie.Insert("overflow-ellipsis", "text-overflow: ellipsis;")
	trie.Insert("overflow-clip", "text-overflow: clip;")
}

// --- Background Utilities ---
func addBackgrounds(trie *Trie) {
	trie.Insert("bg-transparent", "background-color: transparent;")

	// Background attachment
	for _, att := range []string{"fixed", "local", "scroll"} {
		trie.Insert("bg-"+att, "background-attachment: "+att+";")
	}
	// Background colors (using colors and shades)
	// for _, opacity := range opacityValues {
	trie.Insert("bg-white", "background-color: var(--color-white);")
	trie.Insert("bg-black", "background-color: var(--color-black);")
	for _, color := range colorNames {
		class := "bg-" + color
		trie.Insert(class, "background-color: "+toColorVar(color, "")+";")
		for _, shade := range shades {
			class := class + "-" + shade
			trie.Insert(class, "background-color: "+toColorVar(color+"-"+shade, "")+";")
		}
	}
}

// --- Ring Utils ---
func addRingUtils(trie *Trie) {
	// Ring widths (0, 1, 2, 4, 8, etc.)
	ringWidths := []string{"0", "1", "2", "4", "8"}
	for _, w := range ringWidths {
		trie.Insert("ring-"+w, "box-shadow: 0 0 0 "+w+"px var(--ring-color, rgba(59, 130, 246, 0.5));")
	}

	// Ring colors
	trie.Insert("ring-white", "--ring-color: var(--color-white);")
	trie.Insert("ring-black", "--ring-color: var(--color-black);")
	for _, color := range colorNames {
		trie.Insert("ring-"+color, "--ring-color: var(--color-"+color+"-500);")
		for _, shade := range shades {
			trie.Insert("ring-"+color+"-"+shade, "--ring-color: var(--color-"+color+"-"+shade+");")
		}
	}

	// Ring offset (spacing between the ring and the element)
	ringOffsets := []string{"0", "1", "2", "4", "8"}
	for _, offset := range ringOffsets {
		trie.Insert("ring-offset-"+offset, "--ring-offset-width: "+offset+"px;")
	}
}

// --- Border Ring Utils ---
func addDivideUtils(trie *Trie) {
	// Divide widths (0, 2, 4, 8, etc.)
	for _, w := range borderWidths {
		trie.Insert("divide-x-"+w, "border-right-width: "+w+"px; border-left-width: "+w+"px;")
		trie.Insert("divide-y-"+w, "border-top-width: "+w+"px; border-bottom-width: "+w+"px;")
	}
	trie.Insert("divide-white", "border-color: var(--color-white);")
	trie.Insert("divide-black", "border-color: var(--color-black);")
	for _, color := range colorNames {
		class := "divide-" + color
		trie.Insert(class, "border-color: "+toColorVar(color+"-500", "")+";")
		for _, shade := range shades {
			class := class + "-" + shade
			trie.Insert(class, "border-color: "+toColorVar(color+"-"+shade, "")+";")
		}
	}
	// Divide styles
	divideStyles := []string{"solid", "dashed", "dotted", "double", "none"}
	for _, style := range divideStyles {
		trie.Insert("divide-"+style, "border-style: "+style+";")
	}
}

// --- Outline Utils ---
func addOutlineUtils(trie *Trie) {
	// Outline widths (0, 2, 4, 8, etc.)
	for _, w := range borderWidths {
		trie.Insert("outline-"+w, "outline-width: "+w+"px;")
	}
	// Outline colors
	trie.Insert("outline-white", "outline-color: var(--color-white);")
	trie.Insert("outline-black", "outline-color: var(--color-black);")
	for _, color := range colorNames {
		class := "outline-" + color
		trie.Insert(class, "outline-color: "+toColorVar(color+"-500", "")+";")
		for _, shade := range shades {
			class := "outline-" + color + "-" + shade
			trie.Insert(class, "outline-color: "+toColorVar(color+"-"+shade, "")+";")
		}
	}
	// Outline styles
	outlineStyles := []string{"solid", "dashed", "dotted", "double", "none"}
	for _, style := range outlineStyles {
		trie.Insert("outline-"+style, "outline-style: "+style+";")
	}
}

// --- Border Utilities ---
func addBorders(trie *Trie) {
	// Border Width
	trie.Insert("border-px", "border-width: 1px;")
	trie.Insert("border", "border-width: 1px;")
	for _, w := range borderWidths {
		trie.Insert("border-"+w, "border-width: "+w+"px;")
	}
	// Border-colors
	trie.Insert("border-white", "border-color: var(--color-white);")
	trie.Insert("border-black", "border-color: var(--color-black);")
	for _, color := range colorNames {
		class := "border-" + color
		trie.Insert(class, "border-color: "+toColorVar(color+"-500", "")+";")
		for _, shade := range shades {
			trie.Insert("border-"+color+"-"+shade, "border-color: "+toColorVar(color+"-"+shade, "")+";")
		}
	}
	// Base rounded corners (e.g., rounded-sm, rounded-lg)
	trie.Insert("rounded", "border-radius: var(--radius-md);")
	for _, val := range roundedSizes {
		trie.Insert("rounded-"+val, "border-radius: var(--radius-"+val+");")
	}
	for _, style := range borderStyles {
		trie.Insert("border-"+style, "border-style: "+style+";")
	}
	// Generate rounded classes with logical properties
	for _, val := range roundedSizes {
		trie.Insert("rounded-t-"+val, "border-top-left-radius: var(--radius-"+val+");"+"border-top-right-radius: var(--radius-"+val+");")
		trie.Insert("rounded-l-"+val, "border-top-left-radius: var(--radius-"+val+");"+"border-bottom-left-radius: var(--radius-"+val+");")
		trie.Insert("rounded-r-"+val, "border-top-right-radius: var(--radius-"+val+");"+"border-bottom-right-radius: var(--radius-"+val+");")
		trie.Insert("rounded-b-"+val, "border-bottom-left-radius: var(--radius-"+val+");"+"border-bottom-right-radius: var(--radius-"+val+");")
		//
		trie.Insert("rounded-tl-"+val, "border-top-left-radius: var(--radius-"+val+");")
		trie.Insert("rounded-tr-"+val, "border-top-right-radius: var(--radius-"+val+");")
		trie.Insert("rounded-bl-"+val, "border-bottom-left-radius: var(--radius-"+val+");")
		trie.Insert("rounded-br-"+val, "border-bottom-right-radius: var(--radius-"+val+");")
	}
}

// --- Effects Utilities ---
func addEffects(trie *Trie) {
	// Shadows: sizes and with color
	for _, s := range boxShadowSizes {
		trie.Insert("shadow-"+s, "box-shadow: var(--shadow-"+s+");")
	}
	for _, color := range colorNames {
		trie.Insert("shadow-"+color, "box-shadow: var(--shadow-color, var(--color-"+color+"-500));")
	}
	// Opacity
	for _, op := range opacityValues {
		if op == "" {
			return
		}
		trie.Insert("opacity-"+op, "opacity: "+op+"%;")
	}
	// Mix-blend and background-blend modes
	for _, mode := range mixBlendModes {
		trie.Insert("mix-blend-"+mode, "mix-blend-mode: "+mode+";")
		trie.Insert("background-blend-"+mode, "background-blend-mode: "+mode+";")
	}
}

// --- Filters Utilities ---
func addFilters(trie *Trie) {
	trie.Insert("filter", "filter: blur(0);")
	trie.Insert("filter-none", "filter: none;")
	for _, b := range blurValues {
		trie.Insert("blur-"+b, "filter: blur(var(--blur-"+b+"));")
	}
	for _, b := range []string{"0", "50", "75", "90", "95", "100", "105", "110", "125", "150", "200"} {
		trie.Insert("brightness-"+b, "filter: brightness("+b+"%);")
	}
	for _, c := range []string{"0", "50", "75", "100", "125", "150", "200"} {
		trie.Insert("contrast-"+c, "filter: contrast("+c+"%);")
	}
	for _, g := range []string{"0", "25", "50", "75", "100"} {
		trie.Insert("grayscale-"+g, "filter: grayscale("+g+"%);")
	}
	for _, h := range []string{"0", "15", "30", "60", "90", "180"} {
		trie.Insert("hue-rotate-"+h, "filter: hue-rotate("+h+"deg);")
	}
	for _, inv := range []string{"0", "50", "100"} {
		trie.Insert("invert-"+inv, "filter: invert("+inv+"%);")
	}
	for _, s := range []string{"0", "50", "75", "95", "100", "150", "200"} {
		trie.Insert("saturate-"+s, "filter: saturate("+s+"%);")
	}
	for _, s := range []string{"0", "100"} {
		trie.Insert("sepia-"+s, "filter: sepia("+s+"%);")
	}
	for _, ds := range []string{"sm", "md", "lg", "xl", "2xl", "none"} {
		trie.Insert("drop-shadow-"+ds, "filter: drop-shadow("+ds+");")
	}
}

func addTables(trie *Trie) {
	// Table layout
	trie.Insert("table-auto", "table-layout: auto;")
	trie.Insert("table-fixed", "table-layout: fixed;")
	// Border collapse
	trie.Insert("border-collapse", "border-collapse: collapse;")
	trie.Insert("border-separate", "border-collapse: separate;")
	// Table spacing (border spacing)
	spacingValues := []string{"0", "1", "2", "4", "8"}
	for _, val := range spacingValues {
		trie.Insert("border-spacing-"+val, "border-spacing: "+val+"px;")
		trie.Insert("border-spacing-x-"+val, "border-spacing: "+val+"px 0;")
		trie.Insert("border-spacing-y-"+val, "border-spacing: 0 "+val+"px;")
	}
	// Table alignment
	trie.Insert("table-caption-top", "caption-side: top;")
	trie.Insert("table-caption-bottom", "caption-side: bottom;")
	// Table row and cell alignment
	trie.Insert("table-align-top", "vertical-align: top;")
	trie.Insert("table-align-middle", "vertical-align: middle;")
	trie.Insert("table-align-bottom", "vertical-align: bottom;")
	// Table borders
	for _, width := range borderWidths {
		trie.Insert("table-border-"+width, "border-width: "+width+"px;")
	}
	for _, color := range colorNames {
		trie.Insert("table-border-"+color, "border-color: var(--color-"+color+"-500);")
		for _, shade := range shades {
			trie.Insert("table-border-"+color+"-"+shade, "border-color: var(--color-"+color+"-"+shade+");")
		}
	}
	// Table row and cell padding
	paddingValues := []string{"0", "1", "2", "3", "4", "5", "8", "12", "16"}
	for _, val := range paddingValues {
		trie.Insert("table-padding-"+val, "padding: "+val+"px;")
		trie.Insert("table-padding-x-"+val, "padding-left: "+val+"px; padding-right: "+val+"px;")
		trie.Insert("table-padding-y-"+val, "padding-top: "+val+"px; padding-bottom: "+val+"px;")
	}
	// Table row and cell background colors
	for _, color := range colorNames {
		trie.Insert("table-bg-"+color, "background-color: var(--color-"+color+"-500);")
		for _, shade := range shades {
			trie.Insert("table-bg-"+color+"-"+shade, "background-color: var(--color-"+color+"-"+shade+");")
		}
	}
	// Table row and cell text colors
	for _, color := range colorNames {
		trie.Insert("table-text-"+color, "color: var(--color-"+color+"-500);")
		for _, shade := range shades {
			trie.Insert("table-text-"+color+"-"+shade, "color: var(--color-"+color+"-"+shade+");")
		}
	}
}

// --- Transitions and Animations Utilities ---
func addTransitions(trie *Trie) {
	trie.Insert("transition", "transition-property: all;")
	for _, v := range []string{"none", "all", "colors", "opacity", "shadow", "transform"} {
		trie.Insert("transition-"+v, "transition-property: "+v+";")
	}
	for _, d := range []string{"75", "100", "150", "200", "300", "500", "700", "1000"} {
		trie.Insert("duration-"+d, "transition-duration: "+d+"ms;")
		trie.Insert("delay-"+d, "transition-delay: "+d+"ms;")
	}
	for _, e := range []string{"linear", "in", "out", "in-out"} {
		trie.Insert("ease-"+e, "transition-timing-function: "+e+";")
	}
	for _, a := range []string{"none", "spin", "ping", "pulse", "bounce"} {
		trie.Insert("animate-"+a, "animation: "+a+";")
	}
}

// --- Transforms Utilities ---
func addTransforms(trie *Trie) {
	trie.Insert("transform", "transform: translate(0,0);")
	trie.Insert("transform-none", "transform: none;")
	for _, s := range []string{"0", "50", "75", "90", "95", "100", "105", "110", "125", "150"} {
		trie.Insert("scale-"+s, "transform: scale("+s+");")
	}
	for _, axis := range []string{"x", "y"} {
		for _, s := range []string{"0", "50", "75", "90", "95", "100", "105", "110", "125", "150"} {
			trie.Insert("scale-"+axis+"-"+s, "transform: scale"+axis+"("+s+");")
		}
	}
	for _, r := range []string{"0", "45", "90", "180"} {
		trie.Insert("rotate-"+r, "transform: rotate("+r+"deg);")
	}
	for _, axis := range []string{"x", "y"} {
		for _, tVal := range translateValues {
			trie.Insert("translate-"+axis+"-"+tVal, "transform: translate"+axis+"("+tVal+");")
		}
	}
	for _, axis := range []string{"x", "y"} {
		for _, s := range []string{"0", "1", "2", "3", "6", "12"} {
			trie.Insert("skew-"+axis+"-"+s, "transform: skew"+axis+"("+s+"deg);")
		}
	}
	for _, o := range []string{"center", "top", "top-right", "right", "bottom-right", "bottom", "bottom-left", "left", "top-left"} {
		trie.Insert("origin-"+o, "transform-origin: "+o+";")
	}
}

// --- Interactivity Utilities ---
func addInteractivity(trie *Trie) {
	for _, color := range colorNames {
		trie.Insert("accent-"+color, "accent-color: "+color+";")
		trie.Insert("caret-"+color, "caret-color: "+color+";")
	}
	trie.Insert("appearance-none", "appearance: none;")
	cursors := []string{"auto", "default", "pointer", "wait", "text", "move", "help", "not-allowed", "none", "context-menu", "progress", "cell", "crosshair", "vertical-text", "alias", "copy", "no-drop", "grab", "grabbing", "all-scroll", "col-resize", "row-resize", "n-resize", "e-resize", "s-resize", "w-resize", "ne-resize", "nw-resize", "se-resize", "sw-resize", "ew-resize", "ns-resize", "nesw-resize", "nwse-resize", "zoom-in", "zoom-out"}
	for _, cur := range cursors {
		trie.Insert("cursor-"+cur, "cursor: "+cur+";")
	}
	for _, pe := range []string{"none", "auto"} {
		trie.Insert("pointer-events-"+pe, "pointer-events: "+pe+";")
	}
	for _, r := range []string{"none", "x", "y", "both"} {
		trie.Insert("resize-"+r, "resize: "+r+";")
	}
	for _, s := range []string{"auto", "smooth"} {
		trie.Insert("scroll-"+s, "scroll-behavior: "+s+";")
	}
	for _, t := range []string{"auto", "none", "pinch-zoom", "manipulation"} {
		trie.Insert("touch-"+t, "touch-action: "+t+";")
	}
	for _, s := range []string{"none", "text", "all", "auto"} {
		trie.Insert("select-"+s, "user-select: "+s+";")
	}
	for _, wc := range []string{"auto", "scroll", "contents", "transform"} {
		trie.Insert("will-change-"+wc, "will-change: "+wc+";")
	}
}

// --- SVG Utilities ---
func addSVG(trie *Trie) {
	trie.Insert("fill-none", "fill: none;")
	trie.Insert("fill-current", "fill: currentColor;")
	trie.Insert("stroke-none", "stroke: none;")
	trie.Insert("stroke-current", "stroke: currentColor;")
	// for _, opacity := range opacityValues {
	for _, color := range colorNames {
		trie.Insert("fill-"+color, "fill: "+toColorVar(color, "")+";")
		trie.Insert("stroke-"+color, "stroke: "+toColorVar(color, "")+";")
		for _, shade := range shades {
			color := color + "-" + shade
			trie.Insert("stroke-"+color, "stroke: "+toColorVar(color, "")+";")
			trie.Insert("fill-"+color, "stroke: "+toColorVar(color, "")+";")
		}
	}
	// }
	for _, s := range []string{"0", "1", "2", "3"} {
		trie.Insert("stroke-"+s, "stroke-width: "+s+"px;")
	}
}

// --- Accessibility Utilities ---
func addAccessibility(trie *Trie) {
	trie.Insert("sr-only", "position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0;")
	trie.Insert("not-sr-only", "position: static; width: auto; height: auto; padding: 0; margin: 0; overflow: visible; clip: auto; white-space: normal;")
}

// Aspect Ratio Utilities
func addAspectRatio(trie *Trie) {
	trie.Insert("aspect-auto", "aspect-ratio: auto;")
	trie.Insert("aspect-square", "aspect-ratio: 1 / 1;")
	trie.Insert("aspect-video", "aspect-ratio: 16 / 9;")
}

// Scroll Snap Utilities
func addScrollSnap(trie *Trie) {
	trie.Insert("snap-start", "scroll-snap-align: start;")
	trie.Insert("snap-center", "scroll-snap-align: center;")
	trie.Insert("snap-end", "scroll-snap-align: end;")
	trie.Insert("snap-always", "scroll-snap-stop: always;")
	trie.Insert("snap-none", "scroll-snap-type: none;")
}

// Placeholder Styling
func addPlaceholderStyling(trie *Trie) {
	// TODO:
}

// Advanced Utilities (forms, typography/prose, line-clamp, advanced animations)
func addAdvancedUtilities(trie *Trie) {
	// TODO: Forms plugin (sample inputs)

	// TODO: Typography - prose classes

	// TODO: Line-clamp utilities

	// TODO: Advanced animation (sample pulse animation)
	// trie.Insert("animate-pulse", "animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;")
}

// Helper: itoa converts an integer to a string.
func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

// Helper: valVar helper for converting to a value to a css var(...)
func toColorVar(nameStr, opacityVar string) string {
	if opacityVar != "" {
		return "color-mix(in oklab, var(--color-" + nameStr + ") " + opacityVar + "%, transparent)"
	}
	return "var(--color-" + nameStr + ")"
}
