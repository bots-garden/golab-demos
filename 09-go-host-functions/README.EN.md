# Writing Host Functions in Go with Extism

In this article:

- Following the same principle as the previous article: [Write a host function with the Extism Host SDK](https://k33g.hashnode.dev/write-a-host-function-with-the-extism-host-sdk), we will modify the host application developed in **Go** from the article [Run Extism WebAssembly plugins from a Go application](https://k33g.hashnode.dev/run-extism-webassembly-plugins-from-a-go-application) to add a **host function** developed by us.
- We will use the exact same plugin as the one modified in the previous article in the section [Modify the Wasm plugin](https://k33g.hashnode.dev/write-a-host-function-with-the-extism-host-sdk#heading-modify-the-wasm-plugin) to be able to call the host function.

## Prerequisites

You will need:

- Go (v1.20) and TinyGo (v0.28.1) to compile the plugins
- Extism 0.4.0: [Install Extism](https://extism.org/docs/install)
- And at least have read the previous article: [Write a host function with the Extism Host SDK](https://k33g.hashnode.dev/write-a-host-function-with-the-extism-host-sdk) (but probably also all the articles in the series).

## Modifying the host application written in Go

The objective is the same as for the previous article: to develop a host function that allows retrieving messages stored in the host application's memory based on a key. For this, we will use a [Go Map](https://go.dev/blog/maps). And this function will be used (called) by the Wasm plugin.

**Important**: To implement host functions, Extism's Go Host SDK uses the "Golang CGO" package (which allows invoking C code from Go and vice versa).
> See the documentation: [go-host-sdk/#host-functions](https://extism.org/docs/integrate-into-your-codebase/go-host-sdk/#host-functions)

Here is the modified code of the application:


```golang
package main

import (
	"fmt"
	"unsafe"
	"github.com/extism/extism"
)

// 1️⃣
/* 
#include <extism.h>
EXTISM_GO_FUNCTION(memory_get);
*/
import "C" // 2️⃣

// 3️⃣ define a map with some records
var memoryMap = map[string]string{
	"hello": "👋 Hello World 🌍",
	"message": "I 💜 Extism 😍",
}

// 4️⃣ host function definition (callable by the Wasm plugin)
//export memory_get
func memory_get(plugin unsafe.Pointer, inputs *C.ExtismVal, nInputs C.ExtismSize, outputs *C.ExtismVal, nOutputs C.ExtismSize, userData uintptr) {

    // input parameters
	inputSlice := unsafe.Slice(inputs, nInputs)
    // output value
	outputSlice := unsafe.Slice(outputs, nOutputs)

	currentPlugin := extism.GetCurrentPlugin(plugin)

    // 5️⃣ Read the value of inputs from the memory
	keyStr := currentPlugin.InputString(unsafe.Pointer(&inputSlice[0]))

    // 6️⃣ get the associated string value
	returnValue := memoryMap[keyStr]

    // 7️⃣ copy the return value to the memory
	currentPlugin.ReturnString(unsafe.Pointer(&outputSlice[0]), returnValue)

}

func main() {

	// Function is used to define host functions
    // 8️⃣ define a slice of host functions
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

    // 9️⃣ use the updated plugin
	path := "../12-simple-go-mem-plugin/simple.wasm"

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		}}

    // 1️⃣0️⃣ 
	plugin, err := ctx.PluginFromManifest(
		manifest,
		hostFunctions,
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
```

- 1: The first step is to declare the use of `EXTISM_GO_FUNCTION` with the name of the function that will be used.
- 2: Do not forget to import the package `"C"`.
- 3: Create a `map` with some elements. This `map` will be used by the host function.
- 4: Definition of the host function `memory_get`. Do not forget to export the function with `//export memory_get` (a host function will always have the same signature).
- 5: When the host function is called by the Wasm plugin, passing of parameters is done using the memory shared between the plugin and the host. `currentPlugin.InputString(unsafe.Pointer(&inputSlice[0]))` is used to fetch this information from shared memory. `keyStr` is a string that contains the key to retrieve a value from the `map`.
- 6: Fetch the value associated with the key in the `map`.
- 7: Copy the obtained value into memory to allow the Wasm plugin to read it.
- 8: Define an array of host functions. In our case, we create only one, where `"hostMemoryGet"` will be the alias of the function "seen" by the Wasm plugin, `[]extism.ValType{extism.I64}` represents the input parameter type and the return parameter type (remember that Wasm functions only accept numbers - and in our case, these numbers contain the positions and sizes of values in shared memory), and finally, `C.memory_get` which is the definition of our host function.
- 9: Use the modified Wasm plugin.
- 10: Create an instance of the Wasm plugin by passing the array of host functions as a parameter.

**Reminder**: The code of the modified Wasm plugin (written in Go) is here: [Plugin Wasm Go](https://k33g.hashnode.dev/write-a-host-function-with-the-extism-host-sdk#heading-modify-the-wasm-plugin)

## Running the Application

To test your new host application, simply run the following command:

```bash
go run main.go 
```

And you will get the following output, with the messages from each of the keys in the Go `map`:

```bash
🙂 👋 Hello 👋 Hello from the Go Host app 🤗
key: hello, value: 👋 Hello World 🌍
key: message, value: I 💜 Extism 😍
```

🎉 There you have it, we have written a host function in Go, usable with the same Wasm plugin (without modifying it). So this plugin can call a host function written in JavaScript, Go, Rust, ... as long as the application using it has implemented this host function with the same signature.

