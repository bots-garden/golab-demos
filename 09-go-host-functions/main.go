package main

import (
	"fmt"
	"unsafe"
	"github.com/extism/extism"
)

/*
#include <extism.h>
EXTISM_GO_FUNCTION(memory_get);
*/
import "C"

var memoryMap = map[string]string{
	"hello": "👋 Hello World 🌍",
	"message": "I 💜 Extism 😍",
}

//export memory_get
func memory_get(plugin unsafe.Pointer, inputs *C.ExtismVal, nInputs C.ExtismSize, outputs *C.ExtismVal, nOutputs C.ExtismSize, userData uintptr) {

	inputSlice := unsafe.Slice(inputs, nInputs)
	outputSlice := unsafe.Slice(outputs, nOutputs)

	// Get memory pointed to by first element of input slice
	currentPlugin := extism.GetCurrentPlugin(plugin)
	keyStr := currentPlugin.InputString(unsafe.Pointer(&inputSlice[0]))

	returnValue := memoryMap[keyStr]

	currentPlugin.ReturnString(unsafe.Pointer(&outputSlice[0]), returnValue)

	//outputSlice[0] = inputSlice[0]

}

func main() {

	// Function is used to define host functions

	hostFunctions := []extism.Function{
		extism.NewFunction(
			"hostMemoryGet",
			[]extism.ValType{extism.I64},
			[]extism.ValType{extism.I64},
			C.memory_get,
			"",
		),
	}

	ctx := extism.NewContext()

	defer ctx.Free() // this will free the context and all associated plugins

	path := "../12-simple-go-mem-plugin/simple.wasm"

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		}}

	plugin, err := ctx.PluginFromManifest(
		manifest,
		hostFunctions,
		true,
	)
	/*
	plugin, err := ctx.PluginFromManifest(
		manifest,
		[]extism.Function{},
		true,
	)
	*/

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
