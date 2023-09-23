package main

import (
	"context"
	"fmt"

	"github.com/extism/extism"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

var memoryMap = map[string]string{
	"hello":   "👋 Hello World 🌍",
	"message": "I 💜 Extism 😍",
}

func main() {

	ctx := context.Background() // new

	path := "../12-simple-go-mem-plugin/simple.wasm"

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
	

	memory_get := extism.NewHostFunctionWithStack(
		"hostMemoryGet",
		"env",
		func(ctx context.Context, plugin *extism.CurrentPlugin , stack []uint64) {

			// Read the value put in memory by the wasm module
			offset := stack[0]
			bufferInput, err := plugin.ReadBytes(offset)

			if err != nil {
				fmt.Println("🥵", err.Error())
				panic(err)
			}

			// Decode the value to string
			keyStr := string(bufferInput)
			fmt.Println("🟢 keyStr:", keyStr)

			// Read the value from the map
			returnValue := memoryMap[keyStr]

			plugin.Free(offset)

			// write the string result into memory
			offset, err = plugin.WriteBytes([]byte(returnValue))
			if err != nil {
				fmt.Println("😡", err.Error())
				panic(err)
			}

			stack[0] = offset
		},
		[]api.ValueType{api.ValueTypeI64},
		api.ValueTypeI64,
	)

	hostFunctions := []extism.HostFunction{
		memory_get,
	}

	pluginInst, err := extism.NewPlugin(ctx, manifest, config, hostFunctions)

	if err != nil {
		panic(err)
	}

	_, res, err := pluginInst.Call(
		"say_hello",
		[]byte("👋 Hello from the Go Host app 🤗"),
	)

	if err != nil {
		fmt.Println("😡", err)
		//os.Exit(1)
	} else {
		//fmt.Println("🙂", res)
		fmt.Println("🙂", string(res))
	}
}
