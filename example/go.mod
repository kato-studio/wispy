module wispy-example

go 1.23.4

// https://thewebivore.com/using-replace-in-go-mod-to-point-to-your-local-module/
replace github.com/kato-studio/wispy => /Users/theo/Desktop/kato/wispy

require (
	github.com/kato-studio/wispy v0.0.4 // indirect
	github.com/yuin/goldmark v1.7.8 // indirect
)
