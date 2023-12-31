package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/extism/extism"
	"github.com/tetratelabs/wazero"
)

func getConfigFromJsonString(config string) map[string]string {
	var manifestConfig map[string]string
	err := json.Unmarshal([]byte(config), &manifestConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return manifestConfig
}

func main() {

	ctx := context.Background()

	args := os.Args[1:]
	wasmFilePath := args[0]
	functionName := args[1]
	input := args[2]
	manifestConfig := args[3]

	levelInfo := extism.Info

	pluginConfig := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
		LogLevel:     &levelInfo,
	}

	pluginManifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{Path: wasmFilePath},
		},
		AllowedHosts: []string{"*"}, // enable HTTP
		Config:       getConfigFromJsonString(manifestConfig),
	}

	wasmPlugin, err := extism.NewPlugin(ctx, pluginManifest, pluginConfig, nil)

	if err != nil {
		panic(err)
	}

	_, result, err := wasmPlugin.Call(
		functionName,
		[]byte(input),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(string(result))
		os.Exit(0)
	}

}
