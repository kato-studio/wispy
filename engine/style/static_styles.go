package style

var staticStyles = map[string]string{
	//////////////////////////
	// GRID
	//////////////////////////
	"absolute":     "position: absolute;",
	"relative":     "position: relative;",
	"fixed":        "position: fixed;",
	"sticky":       "position: sticky;",
	"hidden":       "display: none;",
	"block":        "display: block;",
	"flex":         "display: flex;",
	"grid":         "display: grid;",
	"inline-block": "display: inline-block;",
	"inline":       "display: inline;",
	//
	"container":         "max-width: 100%; padding: 0 1rem;",
	"flex-row":          "flex-direction: row;",
	"flex-col":          "flex-direction: column;",
	"flex-wrap":         "flex-wrap: wrap;",
	"flex-wrap-reverse": "flex-wrap: wrap-reverse;",
	"flex-no-wrap":      "flex-wrap: nowrap;",
	"items-start":       "align-items: flex-start;",
	"items-end":         "align-items: flex-end;",
	"items-center":      "align-items: center;",
	"justify-start":     "justify-content: flex-start;",
	"justify-end":       "justify-content: flex-end;",
	"justify-center":    "justify-content: center;",
	"justify-between":   "justify-content: space-between;",
	//////////////////////////
	// TEXT
	//////////////////////////
	"text-left":     "text-align: left;",
	"text-right":    "text-align: right;",
	"text-center":   "text-align: center;",
	"text-justify":  "text-align: justify;",
	"uppercase":     "text-transform: uppercase;",
	"lowercase":     "text-transform: lowercase;",
	"capitalize":    "text-transform: capitalize;",
	"italic":        "font-style: italic;",
	"underline":     "text-decoration: underline;",
	"line-through":  "text-decoration: line-through;",
	"no-decoration": "text-decoration: none;",
	//////////////////////////
	// WIDTH / HEIGHT
	//////////////////////////
	"h-full":   "height: 100%;",
	"h-screen": "height: 100vh;",
	"h-svh":    "height: 100svh;",
	"h-lvh":    "height: 100lvh;",
	"h-dvh":    "height: 100dvh;",
	"h-min":    "height: min-content;",
	"h-max":    "height: max-content;",
	"h-fit":    "height: fit-content;",
	"w-full":   "width: 100%;",
	"w-screen": "width: 100vw;",
	"w-svw":    "width: 100svw;",
	"w-lvw":    "width: 100lvw;",
	"w-dvw":    "width: 100dvw;",
	"w-min":    "width: min-content;",
	"w-max":    "width: max-content;",
	"w-fit":    "width: fit-content;",
	//
	"min-h-full":   "min-height: 100%;",
	"min-h-screen": "min-height: 100vh;",
	"min-h-svh":    "min-height: 100svh;",
	"min-h-lvh":    "min-height: 100lvh;",
	"min-h-dvh":    "min-height: 100dvh;",
	"min-h-min":    "min-height: min-content;",
	"min-h-max":    "min-height: max-content;",
	"min-h-fit":    "min-height: fit-content;",
	"min-w-full":   "min-width: 100%;",
	"min-w-screen": "min-width: 100vw;",
	"min-w-svw":    "min-width: 100svw;",
	"min-w-lvw":    "min-width: 100lvw;",
	"min-w-dvw":    "min-width: 100dvw;",
	"min-w-min":    "min-width: min-content;",
	"min-w-max":    "min-width: max-content;",
	"min-w-fit":    "min-width: fit-content;",
	//
	"max-h-full":   "max-height: 100%;",
	"max-h-screen": "max-height: 100vh;",
	"max-h-svh":    "max-height: 100svh;",
	"max-h-lvh":    "max-height: 100lvh;",
	"max-h-dvh":    "max-height: 100dvh;",
	"max-h-min":    "max-height: min-content;",
	"max-h-max":    "max-height: max-content;",
	"max-h-fit":    "max-height: fit-content;",
	"max-w-full":   "max-width: 100%;",
	"max-w-screen": "max-width: 100vw;",
	"max-w-svw":    "max-width: 100svw;",
	"max-w-lvw":    "max-width: 100lvw;",
	"max-w-dvw":    "max-width: 100dvw;",
	"max-w-min":    "max-width: min-content;",
	"max-w-max":    "max-width: max-content;",
	"max-w-fit":    "max-width: fit-content;",
	//
	"size":        "width: 100%; height: 100%;",
	"size-screen": "width: 100vw; height: 100vh;",
	"size-svw":    "width: 100svw; height: 100svh;",
	"size-lvw":    "width: 100lvw; height: 100lvh;",
	"size-dvw":    "width: 100dvw; height: 100dvh;",

	//////////////////////////
	// OTHER
	//////////////////////////
	"border-solid":  "border-style: solid;",
	"border-dashed": "border-style: dashed;",
	"border-dotted": "border-style: dotted;",
	"border-double": "border-style: double;",
	"border-none":   "border-style: none;",

	//////////////////////////
	// ACCESSIBILITY
	//////////////////////////
	"sr-only": "position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0, 0, 0, 0); white-space: nowrap; border-width: 0;",
	"not-sr":  "position: static; width: auto; height: auto; padding: 0; margin: 0; overflow: visible; clip: auto; white-space: normal; border-width: 0;",
}

// TODO🔰 Add base styles / style reset
