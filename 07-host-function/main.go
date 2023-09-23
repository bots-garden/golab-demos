package main

import (
	"context"
	"fmt"
	"os"

	"github.com/extism/extism"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func main() {
	args := os.Args[1:]
	wasmFilePath := args[0]
	functionName := args[1]
	input := args[2]

	ctx := context.Background()

	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
	}

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{Path: wasmFilePath},
	}}

	robot_message := extism.NewHostFunctionWithStack(
		"hostRobotMessage",
		"env",
		func(ctx context.Context, plugin *extism.CurrentPlugin, stack []uint64) {

			offset := stack[0]
			bufferInput, err := plugin.ReadBytes(offset)

			if err != nil {
				fmt.Println(err.Error())
				panic(err)
			}

			message := string(bufferInput)
			fmt.Println("🤖:>", message)

			stack[0] = 0
		},
		[]api.ValueType{api.ValueTypeI64},
		api.ValueTypeI64,
	)


	hostFunctions := []extism.HostFunction{
		robot_message,
	}

	wasmPlugin, err := extism.NewPlugin(ctx, manifest, config, hostFunctions)

	if err != nil {
		panic(err)
	}

	_, result, err := wasmPlugin.Call(
		functionName,
		[]byte(input),
	)
	fmt.Println(string(result))

}
