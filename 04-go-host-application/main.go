package main

import (
	"context"
	"fmt"

	"github.com/extism/extism"
	"github.com/tetratelabs/wazero"
)

func main() {

	//ctx := extism.NewContext()
	ctx := context.Background()

	//defer ctx.Free() // this will free the context and all associated plugins
	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
	}

	//path := "../03-even-with-javascript/hello-js.wasm"
	path := "../01-simple-go-plugin/simple.wasm"

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		}}
	
	/*
	plugin, err := ctx.PluginFromManifest(
		manifest,
		[]extism.Function{},
		true,
	)
	*/

	pluginInst, err := extism.NewPlugin(ctx, manifest, config, nil) // new


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
