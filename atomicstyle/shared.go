package atomicstyle

import (
	dt "github.com/kato-studio/wispy/utils/datatypes"
)

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
	Static       *dt.OrderedMap[string, struct{}]
	Base         *dt.OrderedMap[string, struct{}]
	Sm           *dt.OrderedMap[string, struct{}]
	Md           *dt.OrderedMap[string, struct{}]
	Lg           *dt.OrderedMap[string, struct{}]
	Xl           *dt.OrderedMap[string, struct{}]
	_2xl         *dt.OrderedMap[string, struct{}]
	_3xl         *dt.OrderedMap[string, struct{}]
}

var WispyColors = map[string]map[string]string{
	"white": {
		"500": "rgba(255, 255, 255, var(--opacity))", // #FFFFFF
	},
	"black": {
		"500": "rgba(18, 18, 18, var(--opacity))", // #121212
	},
	"primary": {
		"50":  "rgba(254, 251, 251, var(--opacity))", // #FEFBFB
		"100": "rgba(251, 232, 228, var(--opacity))", // #FBE8E4
		"200": "rgba(244, 190, 179, var(--opacity))", // #F4BEB3
		"300": "rgba(238, 152, 134, var(--opacity))", // #EE9886
		"400": "rgba(232, 113, 90, var(--opacity))",  // #E8715A
		"500": "rgba(225, 71, 42, var(--opacity))",   // #E1472A
		"600": "rgba(197, 56, 28, var(--opacity))",   // #C5381C
		"700": "rgba(156, 44, 22, var(--opacity))",   // #9C2C16
		"800": "rgba(116, 33, 16, var(--opacity))",   // #742110
		"900": "rgba(76, 22, 11, var(--opacity))",    // #4C160B
		"950": "rgba(54, 15, 8, var(--opacity))",     // #360F08
	},
	"secondary": {
		"50":  "rgba(254, 246, 236, var(--opacity))", // #FEF6EC
		"100": "rgba(254, 238, 216, var(--opacity))", // #FEEED8
		"200": "rgba(253, 222, 175, var(--opacity))", // #FDDEAF
		"300": "rgba(253, 209, 140, var(--opacity))", // #FDD18C
		"400": "rgba(253, 197, 98, var(--opacity))",  // #FDC562
		"500": "rgba(255, 187, 56, var(--opacity))",  // #FFBB38
		"600": "rgba(247, 158, 2, var(--opacity))",   // #F79E02
		"700": "rgba(190, 119, 4, var(--opacity))",   // #BE7704
		"800": "rgba(129, 79, 3, var(--opacity))",    // #814F03
		"900": "rgba(74, 44, 3, var(--opacity))",     // #4A2C03
		"950": "rgba(44, 26, 2, var(--opacity))",     // #2C1A02
	},
	"accent": {
		"50":  "rgba(246, 254, 241, var(--opacity))", // #F6FEF1
		"100": "rgba(235, 252, 222, var(--opacity))", // #EBFCDE
		"200": "rgba(218, 249, 194, var(--opacity))", // #DAF9C2
		"300": "rgba(198, 247, 161, var(--opacity))", // #C6F7A1
		"400": "rgba(181, 244, 133, var(--opacity))", // #B5F485
		"500": "rgba(162, 241, 100, var(--opacity))", // #A2F164
		"600": "rgba(125, 236, 40, var(--opacity))",  // #7DEC28
		"700": "rgba(91, 187, 17, var(--opacity))",   // #5BBB11
		"800": "rgba(61, 126, 11, var(--opacity))",   // #3D7E0B
		"900": "rgba(29, 61, 6, var(--opacity))",     // #1D3D06
		"950": "rgba(16, 33, 3, var(--opacity))",     // #102103
	},
	"neutral": {
		"50":  "rgba(251, 250, 249, var(--opacity))", // #FBFAF9
		"100": "rgba(251, 250, 249, var(--opacity))", // #FBFAF9
		"200": "rgba(247, 244, 242, var(--opacity))", // #F7F4F2
		"300": "rgba(246, 241, 239, var(--opacity))", // #F6F1EF
		"400": "rgba(244, 238, 235, var(--opacity))", // #F4EEEB
		"500": "rgba(242, 233, 228, var(--opacity))", // #F2E9E4
		"600": "rgba(225, 208, 198, var(--opacity))", // #E1D0C6
		"700": "rgba(209, 186, 173, var(--opacity))", // #D1BAAD
		"800": "rgba(190, 162, 147, var(--opacity))", // #BEA293
		"900": "rgba(168, 138, 122, var(--opacity))", // #A88A7A
		"950": "rgba(153, 126, 112, var(--opacity))", // #997E70
	},
	"red": {
		"50":  "rgba(255, 232, 229, var(--opacity))", // #FFE8E5
		"100": "rgba(255, 212, 209, var(--opacity))", // #FFD4D1
		"200": "rgba(255, 173, 168, var(--opacity))", // #FFADA8
		"300": "rgba(255, 137, 133, var(--opacity))", // #FF8985
		"400": "rgba(255, 95, 92, var(--opacity))",   // #FF5F5C
		"500": "rgba(255, 51, 52, var(--opacity))",   // #FF3334
		"600": "rgba(255, 9, 5, var(--opacity))",     // #FF0905
		"700": "rgba(214, 7, 0, var(--opacity))",     // #D60700
		"800": "rgba(173, 9, 0, var(--opacity))",     // #AD0900
		"900": "rgba(128, 8, 0, var(--opacity))",     // #800800
		"950": "rgba(102, 9, 0, var(--opacity))",     // #660900
	},
	"blue": {
		"50":  "rgba(231, 230, 255, var(--opacity))", // #E7E6FF
		"100": "rgba(207, 205, 254, var(--opacity))", // #CFCDFE
		"200": "rgba(168, 165, 253, var(--opacity))", // #A8A5FD
		"300": "rgba(124, 120, 253, var(--opacity))", // #7C78FD
		"400": "rgba(81, 75, 252, var(--opacity))",   // #514BFC
		"500": "rgba(39, 33, 251, var(--opacity))",   // #2721FB
		"600": "rgba(12, 4, 230, var(--opacity))",    // #0C04E6
		"700": "rgba(9, 3, 175, var(--opacity))",     // #0903AF
		"800": "rgba(6, 2, 125, var(--opacity))",     // #06027D
		"900": "rgba(4, 1, 75, var(--opacity))",      // #04014B
		"950": "rgba(3, 1, 50, var(--opacity))",      // #030132
	},
	"green": {
		"50":  "rgba(194, 250, 200, var(--opacity))", // #C2FAC8
		"100": "rgba(175, 248, 183, var(--opacity))", // #AFF8B7
		"200": "rgba(147, 245, 157, var(--opacity))", // #93F59D
		"300": "rgba(115, 242, 127, var(--opacity))", // #73F27F
		"400": "rgba(86, 240, 102, var(--opacity))",  // #56F066
		"500": "rgba(53, 237, 71, var(--opacity))",   // #35ED47
		"600": "rgba(19, 215, 39, var(--opacity))",   // #13D727
		"700": "rgba(14, 159, 29, var(--opacity))",   // #0E9F1D
		"800": "rgba(10, 108, 19, var(--opacity))",   // #0A6C13
		"900": "rgba(5, 51, 9, var(--opacity))",      // #053309
		"950": "rgba(3, 28, 5, var(--opacity))",      // #031C05
	},
	"yellow": {
		"50":  "rgba(253, 251, 200, var(--opacity))", // #FDFBC8
		"100": "rgba(253, 249, 181, var(--opacity))", // #FDF9B5
		"200": "rgba(252, 246, 136, var(--opacity))", // #FCF688
		"300": "rgba(250, 243, 96, var(--opacity))",  // #FAF360
		"400": "rgba(249, 239, 52, var(--opacity))",  // #F9EF34
		"500": "rgba(248, 236, 12, var(--opacity))",  // #F8EC0C
		"600": "rgba(222, 209, 9, var(--opacity))",   // #ded109
		"800": "rgba(161, 98, 7, var(--opacity))",    // #a16207
		"900": "rgba(133, 77, 14, var(--opacity))",   // #854d0e
		"950": "rgba(113, 63, 18, var(--opacity))",   // #713f12
	},
	"indigo": {
		"50":  "rgba(241, 242, 254, var(--opacity))", // #F1F2FE
		"100": "rgba(231, 233, 254, var(--opacity))", // #E7E9FE
		"200": "rgba(207, 211, 252, var(--opacity))", // #CFD3FC
		"300": "rgba(178, 185, 251, var(--opacity))", // #B2B9FB
		"400": "rgba(154, 163, 249, var(--opacity))", // #9AA3F9
		"500": "rgba(129, 140, 248, var(--opacity))", // #818CF8
		"600": "rgba(67, 84, 245, var(--opacity))",   // #4354F5
		"700": "rgba(13, 35, 232, var(--opacity))",   // #0D23E8
		"800": "rgba(10, 26, 174, var(--opacity))",   // #0A1AAE
		"900": "rgba(6, 17, 111, var(--opacity))",    // #06116F
		"950": "rgba(4, 12, 77, var(--opacity))",     // #040C4D
	},
	"purple": {
		"50":  "rgba(220, 202, 237, var(--opacity))", // #DCCAED
		"100": "rgba(209, 187, 231, var(--opacity))", // #D1BBE7
		"200": "rgba(185, 153, 219, var(--opacity))", // #B999DB
		"300": "rgba(163, 123, 209, var(--opacity))", // #A37BD1
		"400": "rgba(140, 92, 198, var(--opacity))",  // #8C5CC6
		"500": "rgba(116, 64, 183, var(--opacity))",  // #7440B7
		"600": "rgba(102, 55, 159, var(--opacity))",  // #66379F
		"700": "rgba(86, 46, 132, var(--opacity))",   // #562E84
		"800": "rgba(70, 37, 106, var(--opacity))",   // #46256A
		"900": "rgba(54, 28, 79, var(--opacity))",    // #361C4F
		"950": "rgba(44, 22, 64, var(--opacity))",    // #2C1640
	},
	"pink": {
		"50":  "rgba(253, 242, 248, var(--opacity))", // #fdf2f8
		"100": "rgba(252, 231, 243, var(--opacity))", // #fce7f3
		"200": "rgba(251, 207, 232, var(--opacity))", // #fbcfe8
		"300": "rgba(249, 168, 212, var(--opacity))", // #f9a8d4
		"400": "rgba(244, 114, 182, var(--opacity))", // #f472b6
		"500": "rgba(236, 72, 153, var(--opacity))",  // #ec4899
		"600": "rgba(219, 39, 119, var(--opacity))",  // #db2777
		"700": "rgba(190, 24, 93, var(--opacity))",   // #be185d
		"800": "rgba(157, 23, 77, var(--opacity))",   // #9d174d
		"900": "rgba(131, 24, 67, var(--opacity))",   // #831843
	},
	"gray": {
		"50":  "rgba(252, 252, 252, var(--opacity))", // #FCFCFC
		"100": "rgba(240, 240, 240, var(--opacity))", // #F0F0F0
		"200": "rgba(212, 212, 212, var(--opacity))", // #D4D4D4
		"300": "rgba(186, 186, 186, var(--opacity))", // #BABABA
		"400": "rgba(161, 161, 161, var(--opacity))", // #A1A1A1
		"500": "rgba(134, 134, 134, var(--opacity))", // #868686
		"600": "rgba(112, 112, 112, var(--opacity))", // #707070
		"700": "rgba(89, 89, 89, var(--opacity))",    // #595959
		"800": "rgba(66, 66, 66, var(--opacity))",    // #424242
		"900": "rgba(43, 43, 43, var(--opacity))",    // #2B2B2B
		"950": "rgba(31, 31, 31, var(--opacity))",    // #1F1F1F
	},
}
