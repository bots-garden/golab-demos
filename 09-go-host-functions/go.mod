module go-host-application-host-function

go 1.20

require (
	github.com/extism/extism v0.5.0
	github.com/tetratelabs/wazero v1.5.0
)

require github.com/gobwas/glob v0.2.3 // indirect

replace github.com/extism/extism => ../go-sdk
