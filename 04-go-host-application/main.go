package main

import (
	"context"
	"fmt"
	"os"

	"github.com/extism/extism"
	"github.com/tetratelabs/wazero"
)

func main() {
	//argsWithProg := os.Args
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
			extism.WasmFile{
				Path: wasmFilePath},
		},
		AllowedHosts:  []string{"*"}, 
		/*
			HTTP calls are disallowed by default. 
			If you want to enable HTTP you need 
			to specify the hosts that the plug-in is allowed 
			to communicate with. 
		*/
	}

	pluginInst, err := extism.NewPlugin(ctx, manifest, config, nil)

	if err != nil {
		panic(err)
	}

	_, res, err := pluginInst.Call(
		functionName,
		[]byte(input),
	)

	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	} else {
		fmt.Println(string(res))
	}

}
