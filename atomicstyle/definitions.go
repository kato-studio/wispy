package atomicstyle

import "fmt"

// --- Data Definitions ---

var (
	// Sample color names (expand as needed)
	colorNames = []string{"red", "blue", "green", "black", "white", "gray", "indigo", "purple", "pink"}

	// Common numeric scales
	opacityValues = []string{"0", "5", "10", "20", "25", "30", "40", "50", "60", "70", "75", "80", "90", "95", "100"}
	borderWidths  = []string{"0", "2", "4", "8"}
	roundedSizes  = []string{"none", "sm", "md", "lg", "xl", "2xl", "3xl", "full"}
	borderStyles  = []string{"solid", "dashed", "dotted", "double", "none"}
	// For divide, we reuse borderWidths for the numeric ones.
	ringWidths  = []string{"0", "1", "2", "4", "8"}
	ringOffsets = []string{"0", "1", "2", "4", "8"}

	// Spacing scale (from our previous definition)
	spacingScale = map[string]string{
		"0":  "0",
		"1":  "0.25rem",
		"2":  "0.5rem",
		"3":  "0.75rem",
		"4":  "1rem",
		"5":  "1.25rem",
		"6":  "1.5rem",
		"8":  "2rem",
		"10": "2.5rem",
		"12": "3rem",
		"16": "4rem",
		"20": "5rem",
		"24": "6rem",
		"32": "8rem",
		"40": "10rem",
		"48": "12rem",
		"56": "14rem",
		"64": "16rem",
		"px": "1px",
	}
)

// BuildExtendedTrie builds a trie preloaded with all of our utility CSS classes.
func BuildExtendedTrie() *Trie {
	trie := NewTrie()

	addLayout(trie)
	addFlexGrid(trie)
	addSpacing(trie)
	addSizing(trie)
	addTypography(trie)
	addBackgrounds(trie)

	addBorders(trie)
	addEffects(trie)
	addFilters(trie)
	addTables(trie)
	addTransitions(trie)
	addTransforms(trie)
	addInteractivity(trie)
	addSVG(trie)
	addAccessibility(trie)

	return trie
}

// --- Layout Utilities ---
func addLayout(trie *Trie) {
	// Container and box-sizing
	trie.Insert("container", "width: 100%; margin-left: auto; margin-right: auto; padding: 1rem;")
	trie.Insert("box-border", "box-sizing: border-box;")
	trie.Insert("box-content", "box-sizing: content-box;")

	// Display utilities (block, inline, inline-block, flex, inline-flex, grid, inline-grid, hidden)
	displayValues := map[string]string{
		"block":        "display: block;",
		"inline":       "display: inline;",
		"inline-block": "display: inline-block;",
		"flex":         "display: flex;",
		"inline-flex":  "display: inline-flex;",
		"grid":         "display: grid;",
		"inline-grid":  "display: inline-grid;",
		"hidden":       "display: none;",
	}
	for k, rule := range displayValues {
		trie.Insert("display-"+k, rule)
	}

	// Float and clear
	for _, dir := range []string{"right", "left", "none"} {
		trie.Insert("float-"+dir, "float: "+dir+";")
	}
	clearValues := []string{"left", "right", "both", "none"}
	for _, v := range clearValues {
		trie.Insert("clear-"+v, "clear: "+v+";")
	}

	// Object fit and position
	objectFitValues := map[string]string{
		"contain":    "object-fit: contain;",
		"cover":      "object-fit: cover;",
		"fill":       "object-fit: fill;",
		"none":       "object-fit: none;",
		"scale-down": "object-fit: scale-down;",
	}
	for k, rule := range objectFitValues {
		trie.Insert("object-"+k, rule)
	}
	objectPosValues := map[string]string{
		"bottom":       "object-position: bottom;",
		"center":       "object-position: center;",
		"left":         "object-position: left;",
		"left-bottom":  "object-position: left bottom;",
		"left-top":     "object-position: left top;",
		"right":        "object-position: right;",
		"right-bottom": "object-position: right bottom;",
		"right-top":    "object-position: right top;",
		"top":          "object-position: top;",
	}
	for k, rule := range objectPosValues {
		trie.Insert("object-"+k, rule)
	}

	// Positioning utilities
	positions := []string{"static", "fixed", "absolute", "relative", "sticky"}
	for _, pos := range positions {
		trie.Insert(pos, "position: "+pos+";")
	}

	// Inset, top/right/bottom/left
	insetValues := []string{"0", "auto"}
	for _, val := range insetValues {
		trie.Insert("inset-"+val, fmt.Sprintf("top: %s; right: %s; bottom: %s; left: %s;", val, val, val, val))
		trie.Insert("inset-x-"+val, fmt.Sprintf("left: %s; right: %s;", val, val))
		trie.Insert("inset-y-"+val, fmt.Sprintf("top: %s; bottom: %s;", val, val))
		trie.Insert("top-"+val, "top: "+val+";")
		trie.Insert("right-"+val, "right: "+val+";")
		trie.Insert("bottom-"+val, "bottom: "+val+";")
		trie.Insert("left-"+val, "left: "+val+";")
	}

	// Visibility and z-index
	trie.Insert("visible", "visibility: visible;")
	trie.Insert("invisible", "visibility: hidden;")
	for _, z := range []string{"0", "10", "20", "30", "40", "50", "auto"} {
		trie.Insert("z-"+z, "z-index: "+z+";")
	}
}

// --- Flexbox and Grid Utilities ---
func addFlexGrid(trie *Trie) {
	// Flex direction and wrap
	flexDirs := map[string]string{
		"row":         "flex-direction: row;",
		"row-reverse": "flex-direction: row-reverse;",
		"col":         "flex-direction: column;",
		"col-reverse": "flex-direction: column-reverse;",
	}
	for k, rule := range flexDirs {
		trie.Insert("flex-"+k, rule)
	}
	flexWraps := map[string]string{
		"wrap":         "flex-wrap: wrap;",
		"wrap-reverse": "flex-wrap: wrap-reverse;",
		"nowrap":       "flex-wrap: nowrap;",
	}
	for k, rule := range flexWraps {
		trie.Insert("flex-"+k, rule)
	}

	// Alignment and justify utilities
	alignValues := map[string]string{
		"start":    "flex-start",
		"end":      "flex-end",
		"center":   "center",
		"baseline": "baseline",
		"stretch":  "stretch",
	}
	for k, v := range alignValues {
		trie.Insert("items-"+k, "align-items: "+v+";")
		trie.Insert("content-"+k, "align-content: "+v+";")
		trie.Insert("self-"+k, "align-self: "+v+";")
		trie.Insert("justify-"+k, "justify-content: "+v+";")
	}

	// Flex shorthand and grow/shrink
	flexShort := map[string]string{
		"1":       "flex: 1 1 0%;",
		"auto":    "flex: 1 1 auto;",
		"initial": "flex: 0 1 auto;",
		"none":    "flex: none;",
	}
	for k, rule := range flexShort {
		trie.Insert("flex-"+k, rule)
	}
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
	axes := []string{"x", "y"}
	gapValues := []string{"0", "1", "2", "3", "4", "5", "6", "8", "10", "12", "16", "20", "24", "32", "40", "48", "56", "64", "px"}
	for _, axis := range axes {
		for _, val := range gapValues {
			trie.Insert("gap-"+axis+"-"+val, "gap-"+axis+": "+val+";")
		}
	}
	// Justify-items/self and place-* utilities (sample)
	trie.Insert("justify-items-center", "justify-items: center;")
	trie.Insert("place-content-center", "place-content: center;")
	trie.Insert("place-items-center", "place-items: center;")
	trie.Insert("place-self-center", "place-self: center;")
}

// --- Spacing Utilities ---
func addSpacing(trie *Trie) {
	for k, val := range spacingScale {
		trie.Insert("p-"+k, "padding: "+val+";")
		trie.Insert("px-"+k, "padding-left: "+val+"; padding-right: "+val+";")
		trie.Insert("py-"+k, "padding-top: "+val+"; padding-bottom: "+val+";")
		trie.Insert("m-"+k, "margin: "+val+";")
		trie.Insert("mx-"+k, "margin-left: "+val+"; margin-right: "+val+";")
		trie.Insert("my-"+k, "margin-top: "+val+"; margin-bottom: "+val+";")
	}
	// For space between siblings (note: actual implementation would need complex selectors)
	trie.Insert("space-x-4", ">* + * { margin-left: 1rem; }")
	trie.Insert("space-y-4", ">* + * { margin-top: 1rem; }")
}

// --- Sizing Utilities ---
func addSizing(trie *Trie) {
	wValues := []string{"0", "1", "2", "3", "4", "5", "6", "8", "10", "12", "16", "20", "24", "32", "40", "48", "56", "64", "auto", "px", "full", "screen", "min", "max", "fit"}
	for _, v := range wValues {
		trie.Insert("w-"+v, "width: "+v+";")
		trie.Insert("h-"+v, "height: "+v+";")
	}
	// Similarly for min-w, max-w, min-h, max-h (simplified)
	for _, v := range []string{"0", "full", "min", "max", "fit"} {
		trie.Insert("min-w-"+v, "min-width: "+v+";")
		trie.Insert("max-w-"+v, "max-width: "+v+";")
		trie.Insert("min-h-"+v, "min-height: "+v+";")
		trie.Insert("max-h-"+v, "max-height: "+v+";")
	}
}

// --- Typography Utilities ---
func addTypography(trie *Trie) {
	// Font families
	for _, f := range []string{"sans", "serif", "mono"} {
		trie.Insert("font-"+f, "font-family: "+f+";")
	}
	// Text sizes
	for _, s := range []string{"xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl", "7xl", "8xl", "9xl"} {
		trie.Insert("text-"+s, "font-size: "+s+";")
	}
	// Font weights
	for _, w := range []string{"thin", "extralight", "light", "normal", "medium", "semibold", "bold", "extrabold", "black"} {
		trie.Insert("font-"+w, "font-weight: "+w+";")
	}
	// Tracking and leading (sample)
	trie.Insert("tracking-tight", "letter-spacing: -0.05em;")
	trie.Insert("leading-snug", "line-height: 1.375;")
	// Text align
	for _, a := range []string{"left", "center", "right", "justify"} {
		trie.Insert("text-"+a, "text-align: "+a+";")
	}
	// Text color (using our colors and a sample scale)
	shades := []string{"50", "100", "200", "300", "400", "500", "600", "700", "800", "900"}
	for _, color := range colorNames {
		for _, shade := range shades {
			class := "text-" + color + "-" + shade
			rule := "color: " + color + shade + ";" // placeholder value
			trie.Insert(class, rule)
		}
	}
	// Underline, uppercase, etc.
	trie.Insert("underline", "text-decoration: underline;")
	trie.Insert("line-through", "text-decoration: line-through;")
	trie.Insert("no-underline", "text-decoration: none;")
	trie.Insert("uppercase", "text-transform: uppercase;")
	trie.Insert("lowercase", "text-transform: lowercase;")
	trie.Insert("capitalize", "text-transform: capitalize;")
	trie.Insert("normal-case", "text-transform: none;")
	trie.Insert("truncate", "overflow: hidden; text-overflow: ellipsis; white-space: nowrap;")
	trie.Insert("overflow-ellipsis", "text-overflow: ellipsis;")
	trie.Insert("overflow-clip", "text-overflow: clip;")
}

// --- Background Utilities ---
func addBackgrounds(trie *Trie) {
	// Background attachment
	for _, att := range []string{"fixed", "local", "scroll"} {
		trie.Insert("bg-"+att, "background-attachment: "+att+";")
	}
	// Background colors (using colors and shades)
	shades := []string{"50", "100", "200", "300", "400", "500", "600", "700", "800", "900"}
	for _, color := range colorNames {
		for _, shade := range shades {
			class := "bg-" + color + "-" + shade
			rule := "background-color: " + color + shade + ";" // placeholder value
			trie.Insert(class, rule)
		}
	}
	// Background opacity
	for _, op := range opacityValues {
		trie.Insert("bg-opacity-"+op, "background-opacity: "+op+"%;")
	}
	// Gradient direction
	for _, d := range []string{"t", "tr", "r", "br", "b", "bl", "l", "tl"} {
		trie.Insert("bg-gradient-to-"+d, "background-image: linear-gradient(to "+d+", var(--tw-gradient-stops));")
	}
	// Gradient color stops
	for _, color := range colorNames {
		for _, shade := range shades {
			trie.Insert("from-"+color+"-"+shade, "/* from-color: "+color+shade+" */")
			trie.Insert("via-"+color+"-"+shade, "/* via-color: "+color+shade+" */")
			trie.Insert("to-"+color+"-"+shade, "/* to-color: "+color+shade+" */")
		}
	}
}

// --- Border Utilities ---
func addBorders(trie *Trie) {
	// Border widths
	for _, w := range borderWidths {
		trie.Insert("border-"+w, "border-width: "+w+"px;")
	}
	// Border colors (using sample colors)
	for _, color := range colorNames {
		trie.Insert("border-"+color, "border-color: "+color+";")
		trie.Insert("divide-"+color, "border-color: "+color+";")
		trie.Insert("ring-"+color, "ring-color: "+color+";")
	}
	// Border opacity, divide opacity, and ring opacity
	for _, op := range opacityValues {
		trie.Insert("border-opacity-"+op, "border-opacity: "+op+"%;")
		trie.Insert("divide-opacity-"+op, "divide-opacity: "+op+"%;")
		trie.Insert("ring-opacity-"+op, "ring-opacity: "+op+"%;")
	}
	// Rounded corners
	for _, r := range roundedSizes {
		trie.Insert("rounded-"+r, "border-radius: "+r+";")
	}
	// Rounded by side (t, r, b, l, tl, tr, br, bl)
	for _, side := range []string{"t", "r", "b", "l", "tl", "tr", "br", "bl"} {
		for _, r := range roundedSizes {
			trie.Insert("rounded-"+side+"-"+r, "/* border-radius on "+side+": "+r+" */")
		}
	}
	// Border styles
	for _, s := range borderStyles {
		trie.Insert("border-"+s, "border-style: "+s+";")
	}
	// Divide utilities: numeric and reverse
	for _, axis := range []string{"x", "y"} {
		for _, w := range borderWidths {
			trie.Insert("divide-"+axis+"-"+w, "divide-"+axis+": "+w+"px;")
		}
		trie.Insert("divide-"+axis+"-reverse", "divide-"+axis+"-reverse: 1;")
	}
	// Ring widths and offsets
	for _, rw := range ringWidths {
		trie.Insert("ring-"+rw, "ring-width: "+rw+"px;")
	}
	for _, ro := range ringOffsets {
		trie.Insert("ring-offset-"+ro, "ring-offset-width: "+ro+"px;")
		// Also support ring-offset color using our colors:
		for _, color := range colorNames {
			trie.Insert("ring-offset-"+color, "ring-offset-color: "+color+";")
		}
	}
}

// --- Effects Utilities ---
func addEffects(trie *Trie) {
	// Shadows: sizes and with color
	for _, s := range []string{"sm", "md", "lg", "xl", "2xl", "inner", "none"} {
		trie.Insert("shadow-"+s, "box-shadow: "+s+";")
	}
	for _, color := range colorNames {
		trie.Insert("shadow-"+color, "box-shadow: "+color+";")
	}
	// Opacity
	for _, op := range opacityValues {
		trie.Insert("opacity-"+op, "opacity: "+op+"%;")
	}
	// Mix-blend and background-blend modes
	mixBlendModes := []string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"}
	for _, mode := range mixBlendModes {
		trie.Insert("mix-blend-"+mode, "mix-blend-mode: "+mode+";")
		trie.Insert("background-blend-"+mode, "background-blend-mode: "+mode+";")
	}
}

// --- Filters Utilities ---
func addFilters(trie *Trie) {
	trie.Insert("filter", "filter: blur(0);")
	trie.Insert("filter-none", "filter: none;")
	for _, b := range []string{"none", "sm", "md", "lg", "xl", "2xl", "3xl"} {
		trie.Insert("blur-"+b, "filter: blur("+b+");")
	}
	for _, b := range []string{"0", "50", "75", "90", "95", "100", "105", "110", "125", "150", "200"} {
		trie.Insert("brightness-"+b, "filter: brightness("+b+"%);")
	}
	for _, c := range []string{"0", "50", "75", "100", "125", "150", "200"} {
		trie.Insert("contrast-"+c, "filter: contrast("+c+"%);")
	}
	for _, g := range []string{"0", "100"} {
		trie.Insert("grayscale-"+g, "filter: grayscale("+g+"%);")
	}
	for _, h := range []string{"0", "15", "30", "60", "90", "180"} {
		trie.Insert("hue-rotate-"+h, "filter: hue-rotate("+h+"deg);")
	}
	for _, inv := range []string{"0", "100"} {
		trie.Insert("invert-"+inv, "filter: invert("+inv+"%);")
	}
	for _, s := range []string{"0", "50", "100", "150", "200"} {
		trie.Insert("saturate-"+s, "filter: saturate("+s+"%);")
	}
	for _, s := range []string{"0", "100"} {
		trie.Insert("sepia-"+s, "filter: sepia("+s+"%);")
	}
	for _, ds := range []string{"sm", "md", "lg", "xl", "2xl", "none"} {
		trie.Insert("drop-shadow-"+ds, "filter: drop-shadow("+ds+");")
	}
}

// --- Tables Utilities ---
func addTables(trie *Trie) {
	trie.Insert("border-collapse", "border-collapse: collapse;")
	trie.Insert("border-separate", "border-collapse: separate;")
	for _, t := range []string{"auto", "fixed"} {
		trie.Insert("table-"+t, "table-layout: "+t+";")
	}
	for _, pos := range []string{"top", "bottom"} {
		trie.Insert("caption-"+pos, "caption-side: "+pos+";")
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
	translateValues := []string{"0", "1", "2", "3", "4", "5", "6", "8", "10", "12", "16", "20", "24", "32", "40", "48", "56", "64", "px", "full", "1/2", "1/3", "2/3", "1/4", "2/4", "3/4"}
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
	for _, color := range colorNames {
		trie.Insert("fill-"+color, "fill: "+color+";")
		trie.Insert("stroke-"+color, "stroke: "+color+";")
	}
	trie.Insert("fill-none", "fill: none;")
	trie.Insert("fill-current", "fill: currentColor;")
	trie.Insert("stroke-none", "stroke: none;")
	trie.Insert("stroke-current", "stroke: currentColor;")
	for _, s := range []string{"0", "1", "2"} {
		trie.Insert("stroke-"+s, "stroke-width: "+s+"px;")
	}
}

// --- Accessibility Utilities ---
func addAccessibility(trie *Trie) {
	trie.Insert("sr-only", "position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border: 0;")
	trie.Insert("not-sr-only", "position: static; width: auto; height: auto; padding: 0; margin: 0; overflow: visible; clip: auto; white-space: normal;")
}

// Helper: itoa converts an integer to a string.
func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
