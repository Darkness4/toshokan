module github.com/Darkness4/toshokan/scan/plugins/anchira

go 1.23

require (
	github.com/Darkness4/toshokan v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.33.0
	github.com/shamaton/msgpack/v2 v2.2.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.27.0 // indirect
)

replace github.com/Darkness4/toshokan => ../../..
