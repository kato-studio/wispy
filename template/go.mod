module github.com/kato-studio/wispy/template

go 1.23.7

replace github.com/kato-studio/wispy/wispy_common => ../wispy_common

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/kato-studio/wispy/wispy_common v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/crypto v0.38.0
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)
