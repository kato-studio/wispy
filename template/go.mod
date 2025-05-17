module github.com/kato-studio/wispy/template

go 1.23.7

replace github.com/kato-studio/wispy/utilities => ../utilities

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/kato-studio/wispy/utilities v0.0.0-00010101000000-000000000000
	github.com/tursodatabase/go-libsql v0.0.0-20250416102726-983f7e9acb0e
)

require (
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/libsql/sqlite-antlr4-parser v0.0.0-20240327125255-dbf53b6cbf06 // indirect
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc // indirect
)
