package main

import (
	"fmt"

	"github.com/extism/extism"
)

func main() {

	ctx := extism.NewContext()

	defer ctx.Free() // this will free the context and all associated plugins

	//path := "../03-even-with-javascript/hello-js.wasm"
	path := "../01-simple-go-plugin/simple.wasm"

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		}}

	plugin, err := ctx.PluginFromManifest(
		manifest,
		[]extism.Function{},
		true,
	)

	if err != nil {
		panic(err)
	}

	res, err := plugin.Call(
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
