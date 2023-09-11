package main

import (
	"context"
	"fmt"

	"github.com/extism/extism"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func main() {

	ctx := context.Background() // new

	path := "../100-go-plugin/simple.wasm"

	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
	}

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		}}

	print_message := extism.NewHostFunctionWithStack(
		"hostPrintMessage",
		"env",
		func(ctx context.Context, plugin *extism.CurrentPlugin, stack []uint64) {

			offset := stack[0]
			bufferInput, err := plugin.ReadBytes(offset)

			if err != nil {
				fmt.Println("🥵", err.Error())
				panic(err)
			}

			message := string(bufferInput)
			fmt.Println("🟢:", message)

			stack[0] = 0
		},
		[]api.ValueType{api.ValueTypeI64},
		api.ValueTypeI64,
	)

	display_message := extism.NewHostFunctionWithStack(
		"hostDisplayMessage",
		"env",
		func(ctx context.Context, plugin *extism.CurrentPlugin, stack []uint64) {

			offset := stack[0]
			bufferInput, err := plugin.ReadBytes(offset)

			if err != nil {
				fmt.Println("🥵", err.Error())
				panic(err)
			}

			message := string(bufferInput)
			fmt.Println("🟣:", message)

			stack[0] = 0
		},
		[]api.ValueType{api.ValueTypeI64},
		api.ValueTypeI64,
	)

	hostFunctions := []extism.HostFunction{
		print_message,
		display_message,

	}
	//fmt.Println("📦", hostFunctions)

	//pluginInst, err := extism.NewPlugin(ctx, manifest, config, hostFunctions)
	pluginInst, err := extism.NewPlugin(ctx, manifest, config, hostFunctions)

	if err != nil {
		panic(err)
	}

	_, res, err := pluginInst.Call(
		"say_hello",
		[]byte("John Doe"),
	)
	fmt.Println("🙂", string(res))

	/*
		if err != nil {
			fmt.Println("😡", err)
			//os.Exit(1)
		} else {
			//fmt.Println("🙂", res)
			fmt.Println("🙂", string(res))
		}
	*/

	_, res, err = pluginInst.Call(
		"say_hey",
		[]byte("Jane Doe"),
	)
	fmt.Println("🙂", string(res))

	/*
		if err != nil {
			fmt.Println("😡", err)
			//os.Exit(1)
		} else {
			//fmt.Println("🙂", res)
			fmt.Println("🙂", string(res))
		}
	*/

}
