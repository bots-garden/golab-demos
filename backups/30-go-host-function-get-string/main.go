package main

import (
	"context"
	"fmt"

	"github.com/extism/extism"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)


var extismMessage = "I 💜 Extism"

func main() {

	ctx := context.Background() // new

	path := "../40-simple-go-plugin/simple.wasm"

	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
	}

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		},
		AllowedHosts:  []string{"*"}, 
	}
	

	get_string := extism.NewHostFunctionWithStack(
		"hostGetString",
		"env",
		func(ctx context.Context, plugin *extism.CurrentPlugin , stack []uint64) {

			// Read the value from the map
			returnValue := extismMessage

			// Write the string result into memory
			offset, err := plugin.WriteBytes([]byte(returnValue))
			if err != nil {
				fmt.Println("😡", err.Error())
				panic(err)
			}
			// The only way to share data is via writing memory 
			// and sharing offsets (containing pos and length of the data)
			// Return the offset of the string (position and size)
			stack[0] = offset
		},
		nil, // no parameters (before []api.ValueType{api.ValueTypeI64})
		api.ValueTypeI64, // offset
	)

	hostFunctions := []extism.HostFunction{
		get_string,
	}

	pluginInst, err := extism.NewPlugin(ctx, manifest, config, hostFunctions)

	if err != nil {
		panic(err)
	}

	_, res, err := pluginInst.Call(
		"say_hello",
		nil, // []byte("👋 Hello from the Go Host app 🤗")
	)

	if err != nil {
		fmt.Println("😡", err)
	} else {
		fmt.Println(string(res))
	}
}
