module github.com/kato-studio/wispy/auth

go 1.23.7

replace github.com/kato-studio/wispy/utilities => ../utilities

require (
	github.com/go-oauth2/oauth2/v4 v4.5.3
	github.com/google/uuid v1.6.0
	github.com/kato-studio/wispy/utilities v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.28
	golang.org/x/crypto v0.33.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/tidwall/btree v0.0.0-20191029221954-400434d76274 // indirect
	github.com/tidwall/buntdb v1.1.2 // indirect
	github.com/tidwall/gjson v1.12.1 // indirect
	github.com/tidwall/grect v0.0.0-20161006141115-ba9a043346eb // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/rtree v0.0.0-20180113144539-6cd427091e0e // indirect
	github.com/tidwall/tinyqueue v0.0.0-20180302190814-1e39f5511563 // indirect
	golang.org/x/net v0.33.0 // indirect
)
