module github.com/kato-studio/wispy/auth

go 1.23.7

replace github.com/kato-studio/wispy/wispy_common => ../wispy_common

replace github.com/kato-studio/wispy/template => ../template

require (
	github.com/google/uuid v1.6.0
	github.com/kato-studio/wispy/template v0.0.0-00010101000000-000000000000
	github.com/kato-studio/wispy/wispy_common v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.33.0
)

require golang.org/x/oauth2 v0.30.0
