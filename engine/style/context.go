package style

import "kato-studio/go-wispy/utils"

type StyleCTX struct {
	Colors            map[string]map[string]string
	StaticStyles      map[string]string
	Styles            Styles
	AppendCssVariable func(color_name, color string)
	StyleCategories   map[string]map[string]string
}

func NewStyleCTX(colors map[string]map[string]string) StyleCTX {
	if colors == nil {
		return StyleCTX{
			Colors: map[string]map[string]string{},
			Styles: Styles{},
		}
	}
	return StyleCTX{
		Colors: colors,
		Styles: Styles{},
	}
}

type Styles struct {
	CssVariables map[string]string
	Static       utils.UniqueSet
	Base         utils.UniqueSet
	Sm           utils.UniqueSet
	Md           utils.UniqueSet
	Lg           utils.UniqueSet
	Xl           utils.UniqueSet
	_2xl         utils.UniqueSet
	_3xl         utils.UniqueSet
}

// type ColorsStruct struct {
// 	Primary   map[string]string
// 	Secondary map[string]string
// 	Accent    map[string]string
// 	Neutral   map[string]string
// 	Red       map[string]string
// 	Blue      map[string]string
// 	Green     map[string]string
// 	Yellow    map[string]string
// 	Indigo    map[string]string
// 	Purple    map[string]string
// 	Pink      map[string]string
// 	Grey      map[string]string
// }

var WispyColors = map[string]map[string]string{
	"white": {
		"500": "#FFFFFF",
	},
	"black": {
		"500": "#121212",
	},
	"primary": {
		"50":  "#FEFBFB",
		"100": "#FBE8E4",
		"200": "#F4BEB3",
		"300": "#EE9886",
		"400": "#E8715A",
		"500": "#E1472A",
		"600": "#C5381C",
		"700": "#9C2C16",
		"800": "#742110",
		"900": "#4C160B",
		"950": "#360F08",
	},
	"secondary": {
		"50":  "#FEF6EC",
		"100": "#FEEED8",
		"200": "#FDDEAF",
		"300": "#FDD18C",
		"400": "#FDC562",
		"500": "#FFBB38",
		"600": "#F79E02",
		"700": "#BE7704",
		"800": "#814F03",
		"900": "#4A2C03",
		"950": "#2C1A02",
	},
	"accent": {
		"50":  "#F6FEF1",
		"100": "#EBFCDE",
		"200": "#DAF9C2",
		"300": "#C6F7A1",
		"400": "#B5F485",
		"500": "#A2F164",
		"600": "#7DEC28",
		"700": "#5BBB11",
		"800": "#3D7E0B",
		"900": "#1D3D06",
		"950": "#102103",
	},
	"neutral": {
		"50":  "#FBFAF9",
		"100": "#FBFAF9",
		"200": "#F7F4F2",
		"300": "#F6F1EF",
		"400": "#F4EEEB",
		"500": "#F2E9E4",
		"600": "#E1D0C6",
		"700": "#D1BAAD",
		"800": "#BEA293",
		"900": "#A88A7A",
		"950": "#997E70",
	},
	"red": {
		"50":  "#FFE8E5",
		"100": "#FFD4D1",
		"200": "#FFADA8",
		"300": "#FF8985",
		"400": "#FF5F5C",
		"500": "#FF3334",
		"600": "#FF0905",
		"700": "#D60700",
		"800": "#AD0900",
		"900": "#800800",
		"950": "#660900",
	},
	"blue": {
		"50":  "#E7E6FF",
		"100": "#CFCDFE",
		"200": "#A8A5FD",
		"300": "#7C78FD",
		"400": "#514BFC",
		"500": "#2721FB",
		"600": "#0C04E6",
		"700": "#0903AF",
		"800": "#06027D",
		"900": "#04014B",
		"950": "#030132",
	},
	"green": {
		"50":  "#C2FAC8",
		"100": "#AFF8B7",
		"200": "#93F59D",
		"300": "#73F27F",
		"400": "#56F066",
		"500": "#35ED47",
		"600": "#13D727",
		"700": "#0E9F1D",
		"800": "#0A6C13",
		"900": "#053309",
		"950": "#031C05",
	},
	"yellow": {
		"50":  "#FDFBC8",
		"100": "#FDF9B5",
		"200": "#FCF688",
		"300": "#FAF360",
		"400": "#F9EF34",
		"500": "#F8EC0C",
		"600": "#ded109",
		"800": "#a16207",
		"900": "#854d0e",
		"950": "#713f12",
	},
	"indigo": {
		"50":  "#F1F2FE",
		"100": "#E7E9FE",
		"200": "#CFD3FC",
		"300": "#B2B9FB",
		"400": "#9AA3F9",
		"500": "#818CF8",
		"600": "#4354F5",
		"700": "#0D23E8",
		"800": "#0A1AAE",
		"900": "#06116F",
		"950": "#040C4D",
	},
	"purple": {
		"50":  "#DCCAED",
		"100": "#D1BBE7",
		"200": "#B999DB",
		"300": "#A37BD1",
		"400": "#8C5CC6",
		"500": "#7440B7",
		"600": "#66379F",
		"700": "#562E84",
		"800": "#46256A",
		"900": "#361C4F",
		"950": "#2C1640",
	},
	"pink": {
		"50":  "#fdf2f8",
		"100": "#fce7f3",
		"200": "#fbcfe8",
		"300": "#f9a8d4",
		"400": "#f472b6",
		"500": "#ec4899",
		"600": "#db2777",
		"700": "#be185d",
		"800": "#9d174d",
		"900": "#831843",
	},
	"gray": {
		"50":  "#FCFCFC",
		"100": "#F0F0F0",
		"200": "#D4D4D4",
		"300": "#BABABA",
		"400": "#A1A1A1",
		"500": "#868686",
		"600": "#707070",
		"700": "#595959",
		"800": "#424242",
		"900": "#2B2B2B",
		"950": "#1F1F1F",
	},
}
