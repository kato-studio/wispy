// here is my golang definitions file that is referenced to generate a trie of class names and their css properties, let's make small changes to use css variables that exist in the theme where it makes sense to do so ```
package atomicstyle

type ClassRecipe struct {
	Attribute string
}

var Recipes = map[string]ClassRecipe{
	// Width and Height
	"w-":     {Attribute: "width: calc(var(--spacing) * %s);"},
	"h-":     {Attribute: "height: calc(var(--spacing) * %s);"},
	"max-w-": {Attribute: "max-width: calc(var(--spacing) * %s);"},
	"max-h-": {Attribute: "max-height: calc(var(--spacing) * %s);"},
	"min-w-": {Attribute: "min-width: calc(var(--spacing) * %s);"},
	"min-h-": {Attribute: "min-height: calc(var(--spacing) * %s);"},

	// Padding
	"p-":  {Attribute: "padding: calc(var(--spacing) * %s);;"},
	"px-": {Attribute: "padding-left: calc(var(--spacing) * %s);; padding-right: calc(var(--spacing) * %s);;"},
	"py-": {Attribute: "padding-top: calc(var(--spacing) * %s);; padding-bottom: calc(var(--spacing) * %s);;"},
	"pt-": {Attribute: "padding-top: calc(var(--spacing) * %s);;"},
	"pb-": {Attribute: "padding-bottom: calc(var(--spacing) * %s);;"},
	"pl-": {Attribute: "padding-left: calc(var(--spacing) * %s);;"},
	"pr-": {Attribute: "padding-right: calc(var(--spacing) * %s);;"},

	// Margin (with negative value support)
	"m-":  {Attribute: "margin: calc(var(--spacing) * %s);;"},
	"mx-": {Attribute: "margin-inline: calc(var(--spacing) * %s);; margin-right: calc(var(--spacing) * %s);;"},
	"my-": {Attribute: "margin-top: calc(var(--spacing) * %s);; margin-bottom: calc(var(--spacing) * %s);;"},
	"mt-": {Attribute: "margin-top: calc(var(--spacing) * %s);;"},
	"mb-": {Attribute: "margin-bottom: calc(var(--spacing) * %s);;"},
	"ml-": {Attribute: "margin-left: calc(var(--spacing) * %s);;"},
	"mr-": {Attribute: "margin-right: calc(var(--spacing) * %s);;"},
	//
	"-m-":  {Attribute: "margin: calc(var(--spacing) * %s * -1);;"},
	"-mt-": {Attribute: "margin-top: calc(var(--spacing) * %s * -1);;"},
	"-mb-": {Attribute: "margin-bottom: calc(var(--spacing) * %s * -1);;"},
	"-ml-": {Attribute: "margin-left: calc(var(--spacing) * %s * -1);;"},
	"-mr-": {Attribute: "margin-right: calc(var(--spacing) * %s * -1);;"},

	// Grid
	"grid-":  {Attribute: "display: grid;"},
	"cols-":  {Attribute: "grid-template-columns: repeat(%s, minmax(0, 1fr));"},
	"rows-":  {Attribute: "grid-template-rows: repeat(%s, minmax(0, 1fr));"},
	"gap-":   {Attribute: "gap: calc(var(--spacing) * %s);"},
	"gap-x-": {Attribute: "column-gap: calc(var(--spacing) * %s);"},
	"gap-y-": {Attribute: "row-gap: calc(var(--spacing) * %s);"},

	// Absolute
	"top-":     {Attribute: "top: calc(var(--spacing) * %s);"},
	"bottom-":  {Attribute: "bottom: calc(var(--spacing) * %s);"},
	"left-":    {Attribute: "left: calc(var(--spacing) * %s);"},
	"right-":   {Attribute: "right: calc(var(--spacing) * %s);"},
	"-top-":    {Attribute: "top: calc(var(--spacing) * %s * -1);"},
	"-bottom-": {Attribute: "bottom: calc(var(--spacing) * %s * -1);"},
	"-left-":   {Attribute: "left: calc(var(--spacing) * %s * -1);"},
	"-right-":  {Attribute: "right: calc(var(--spacing) * %s * -1);"},
}
