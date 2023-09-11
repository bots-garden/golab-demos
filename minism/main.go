package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/extism/extism"
	"github.com/tetratelabs/wazero"
)

func getLevel(logLevel string) extism.LogLevel {
	level := extism.Off
	switch logLevel {
	case "error":
		level = extism.Error
	case "warn":
		level = extism.Warn
	case "info":
		level = extism.Info
	case "debug":
		level = extism.Debug
	case "trace":
		level = extism.Trace
	}
	return level
}

func execute(wasmFilePath string, wasmFunctionName string, input string, logLevel string, allowHosts string, config string, wasi bool) {

	hosts := strings.Split(strings.ReplaceAll(allowHosts," ",""), ",")

	ctx := context.Background()

	level := getLevel(logLevel)

	extismConfig := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   wasi,
		LogLevel:     &level,
	}

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
		AllowedHosts: hosts,
	}

	pluginInst, err := extism.NewPlugin(ctx, manifest, extismConfig, nil)

	if err != nil {
		panic(err)
	}

	_, res, err := pluginInst.Call(
		wasmFunctionName,
		[]byte(input),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println(string(res))
		os.Exit(0)
	}

}

func parseCommand(command string, args []string) error {
	//fmt.Println("Command:", command)
	//fmt.Println("Args:", args)
	switch command {
	case "start", "listen":
		fmt.Println("start")

		return nil

	case "call":

		wasmFilePath := flag.Args()[1]     // path of the wasm file
		wasmFunctionName := flag.Args()[2] // function name

		flagSet := flag.NewFlagSet("call", flag.ExitOnError)

		input := flagSet.String("input", "", "input")
		logLevel := flagSet.String("log-level", "", "Log level")
		allowHosts := flagSet.String("allow-hosts", "*", "")
		config := flagSet.String("config", "", "")
		wasi := flagSet.Bool("wasi", true, "")

		flagSet.Parse(args[2:])

		execute(wasmFilePath, wasmFunctionName, *input, *logLevel, *allowHosts, *config, *wasi)
		return nil

	case "version":
		fmt.Println("000")
		//os.Exit(0)
		return nil

	default:
		return fmt.Errorf("🔴 invalid command")
	}
}

func main() {

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("🔴 invalid command")
		os.Exit(0)
	}

	command := flag.Args()[0]

	errCmd := parseCommand(command, flag.Args()[1:])
	if errCmd != nil {
		fmt.Println(errCmd)
		os.Exit(1)
	}

}
