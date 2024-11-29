package style

import (
	"kato-studio/go-wispy/utils"
	"maps"
)

type StyleCategory struct {
	Attr            string
	Directions      map[string]string
	Options         map[string]string
	Exclude         utils.UniqueSet
	Format          string
	IsColor         bool
	PrefixHasStatic bool
}

// All hex value from 100% to 0% alpha in 10% increments
var HEX_OPACITY = map[string]string{
	"100": "FF", "90": "E6", "80": "CC", "70": "B3", "60": "99", "50": "80", "40": "66", "30": "4D", "20": "33", "10": "1A", "0": "00",
}
var MARGIN_PADDING = map[string]string{
	"px": "1px", "0.5": "0.125rem", "1": "0.25rem", "1.5": "0.375rem", "2": "0.5rem", "2.5": "0.625rem", "3": "0.75rem", "4": "1rem", "5": "1.25rem", "6": "1.5rem", "7": "1.75rem", "8": "2rem", "9": "2.25rem", "10": "2.5rem", "11": "2.75rem", "12": "3rem", "14": "3.5rem", "16": "4rem", "20": "5rem", "24": "6rem", "28": "7rem", "32": "8rem", "36": "9rem", "40": "10rem", "44": "11rem", "48": "12rem", "52": "13rem", "56": "14rem", "60": "15rem", "64": "16rem", "72": "18rem", "80": "20rem", "96": "24rem", "100": "26rem",
}
var PERCENTAGE = map[string]string{
	"auto": "auto", "full": "100%", "1/2": "50%", "1/3": "33.333333%", "2/3": "66.666667%", "1/4": "25%", "2/4": "50%", "3/4": "75%", "1/5": "20%", "2/5": "40%", "3/5": "60%", "4/5": "80%", "1/6": "16.666667%", "2/6": "33.333333%", "3/6": "50%", "4/6": "66.666667%", "5/6": "83.333333%", "1/12": "8.333333%", "2/12": "16.666667%", "3/12": "25%", "4/12": "33.333333%", "5/12": "41.666667%", "6/12": "50%", "7/12": "58.333333%", "8/12": "66.666667%", "9/12": "75%", "10/12": "83.333333%", "11/12": "91.666667%",
}
var TWELVE = map[string]string{
	"1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7", "8": "8", "9": "9", "10": "10", "11": "11", "12": "12",
}
var GAP = map[string]string{
	"0": "0", "1": "0.25rem", "2": "0.5rem", "3": "0.75rem", "4": "1rem", "5": "1.25rem", "6": "1.5rem", "7": "1.75rem", "8": "2rem", "9": "2.25rem", "10": "2.5rem", "12": "3rem", "14": "3.5rem", "16": "4rem", "20": "5rem", "24": "6rem", "28": "7rem", "32": "8rem", "36": "9rem", "40": "10rem", "44": "11rem", "48": "12rem", "52": "13rem", "56": "14rem", "60": "15rem", "64": "16rem", "72": "18rem", "80": "20rem", "96": "24rem",
}
var INCREMENT = map[string]string{
	"1": "1px", "2": "2px", "3": "3px", "4": "4px", "5": "5px", "6": "6px", "8": "8px", "10": "10px", "12": "12px", "16": "16px",
}
var SIZES = map[string]string{
	"sm": "24rem", "md": "28rem", "lg": "32rem", "xl": "36rem", "2xl": "42rem", "3xl": "48rem", "4xl": "56rem", "5xl": "64rem", "6xl": "72rem", "7xl": "80rem", "8xl": "90rem", "9xl": "100rem",
}
var SCREENS = map[string]string{
	"sm": "640px", "md": "768px", "lg": "1024px", "xl": "1280px", "2xl": "1536px",
}

var MARGIN_PADDING_PERCENTAGE = maps.Clone(MARGIN_PADDING)
var MARGIN_PADDING_PERCENTAGE_SIZES = maps.Clone(MARGIN_PADDING_PERCENTAGE)
var MARGIN_PADDING_PERCENTAGE_WIDTHS_SCREENS = maps.Clone(MARGIN_PADDING_PERCENTAGE_SIZES)

func init() {
	maps.Insert(MARGIN_PADDING_PERCENTAGE, maps.All(PERCENTAGE))
	maps.Insert(MARGIN_PADDING_PERCENTAGE_SIZES, maps.All(SIZES))
	maps.Insert(MARGIN_PADDING_PERCENTAGE_WIDTHS_SCREENS, maps.All(SCREENS))
}

var styleCategories = map[string]StyleCategory{
	"p": {
		Attr:    "padding",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"pt": {
		Attr:    "padding-top",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"pb": {
		Attr:    "padding-bottom",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"pl": {
		Attr:    "padding-left",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"pr": {
		Attr:    "padding-right",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"px": {
		Attr:    "padding-left",
		Options: MARGIN_PADDING_PERCENTAGE,
		Format:  "padding-right: %[2]s;",
	},
	"py": {
		Attr:    "padding-top",
		Options: MARGIN_PADDING_PERCENTAGE,
		Format:  "padding-bottom: %[2]s;",
	},
	"m": {
		Attr:    "margin",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"mt": {
		Attr:    "margin-top",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"mb": {
		Attr:    "margin-bottom",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"ml": {
		Attr:    "margin-left",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"mr": {
		Attr:    "margin-right",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"mx": {
		Attr:    "",
		Options: MARGIN_PADDING_PERCENTAGE,
		Format:  "margin-left: %[2]s; margin-right: %[2]s;",
	},
	"my": {
		Attr:    "",
		Options: MARGIN_PADDING_PERCENTAGE,
		Format:  "margin-top: %[2]s; margin-bottom: %[2]s;",
	},
	"top": {
		Attr:    "top",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"bottom": {
		Attr:    "bottom",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"left": {
		Attr:    "left",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"right": {
		Attr:    "right",
		Options: MARGIN_PADDING_PERCENTAGE,
	},
	"bg": {
		IsColor: true,
		Attr:    "background",
	},
	"color": {
		IsColor: true,
		Attr:    "color",
	},
	"text": {
		Attr:    "font-size",
		Exclude: utils.NewUniqueSet("left", "right", "center", "justify"),
		Options: map[string]string{
			"xs": "0.75rem", "sm": "0.875rem", "base": "1rem", "lg": "1.125rem", "xl": "1.25rem", "2xl": "1.5rem", "3xl": "1.875rem", "4xl": "2.25rem", "5xl": "3rem", "6xl": "4rem",
		},
	},
	"font": {
		Attr: "font-weight",
		Options: map[string]string{
			"thin": "100", "slime": "200", "light": "300", "normal": "400", "medium": "500", "semibold": "600", "bold": "700", "black": "900",
		},
	},
	"rounded": {
		Attr: "border-radius",
		Options: map[string]string{
			"none": "0", "sm": "0.125rem", "md": "0.375rem", "lg": "0.5rem", "full": "9999px",
		},
	},
	"cols": {
		Attr:    "grid-template-columns",
		Options: TWELVE,
		Format:  "%[1]s: repeat(%[2]s, 1fr);",
	},
	"border": {
		IsColor: true,
		Attr:    "border",
		Options: INCREMENT,
	},
	"border-t": {
		Attr:    "border-top",
		Options: INCREMENT,
	},
	"border-b": {
		Attr:    "border-bottom",
		Options: INCREMENT,
	},
	"border-l": {
		Attr:    "border-left",
		Options: INCREMENT,
	},
	"border-r": {
		Attr:    "border-right",
		Options: INCREMENT,
	},
	"border-color": {
		IsColor: true,
		Attr:    "border-color",
	},
	"rows": {
		Attr:    "grid-template-rows",
		Options: TWELVE,
		Format:  "%[1]s: repeat(%[2]s, 1fr);",
	},
	"gap": {
		Attr:    "gap",
		Options: GAP,
	},
	"row-gap": {
		Attr:    "row-gap",
		Options: GAP,
	},
	"col-gap": {
		Attr:    "column-gap",
		Options: GAP,
	},
	"w": {
		Attr:            "width",
		Options:         MARGIN_PADDING_PERCENTAGE_WIDTHS_SCREENS,
		PrefixHasStatic: true,
	},
	"h": {
		Attr:            "height",
		Options:         MARGIN_PADDING_PERCENTAGE_WIDTHS_SCREENS,
		PrefixHasStatic: true,
	},
	"min-w": {
		Attr:            "min-width",
		Options:         MARGIN_PADDING_PERCENTAGE_WIDTHS_SCREENS,
		PrefixHasStatic: true,
	},
	"min-h": {
		Attr:            "min-height",
		Options:         MARGIN_PADDING_PERCENTAGE_WIDTHS_SCREENS,
		PrefixHasStatic: true,
	},
	"max-w": {
		Attr:            "max-width",
		Options:         MARGIN_PADDING_PERCENTAGE_WIDTHS_SCREENS,
		PrefixHasStatic: true,
	},
	"max-h": {
		Attr:            "max-height",
		Options:         MARGIN_PADDING_PERCENTAGE_WIDTHS_SCREENS,
		PrefixHasStatic: true,
	},
	"size": {
		Attr:    "",
		Format:  "width: %[2]s; height: %[2]s;",
		Options: MARGIN_PADDING_PERCENTAGE_SIZES,
	},
}