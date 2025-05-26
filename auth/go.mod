module github.com/kato-studio/wispy/auth

go 1.23.7

replace github.com/kato-studio/wispy/wispy_common => ../wispy_common

replace github.com/kato-studio/wispy/template => ../template

require (
	github.com/google/uuid v1.6.0
	github.com/kato-studio/wispy/wispy_common v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.38.0
)

require (
	github.com/kato-studio/wispy/template v0.0.0-00010101000000-000000000000
	golang.org/x/oauth2 v0.30.0
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/segmentio/ksuid v1.0.4 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)
