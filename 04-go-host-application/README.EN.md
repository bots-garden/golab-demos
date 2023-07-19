# Run Extism WebAssembly plugins from a Go application

For a few days, we have seen that it is possible to develop WebAssembly plugins with the Extism Plugin Development Kit and run them with the Extism CLI. Today, it's time to move up a level: we're going to create an application in Go that can load these plugins and run them as the CLI does.

To do this we will use the **Host SDK** of Extism for the Go language. As a reminder, Extism provides Host SDKs for many languages (https://extism.org/docs/category/integrate-into-your-codebase).

As a reminder, a host application is an application that thanks to a Wasm runtime SDK, is capable of running WebAssembly programs. The **Host SDKs** of Extism are "overlays" on the Wasm runtime SDK to make your life easier (avoid complicated plumbing).

Currently, Extism uses the **[WasmTime](https://wasmtime.dev/)** runtime.

> If I refer to this [issue (WASI threads support)](https://github.com/extism/extism/issues/357), it is not impossible that the support of other Wasm runtimes will be taken into account, and in particular [Wazero](https://wazero.io/).

But enough talk, let's get down to business.

## Prerequisites

You will need

- Go (v1.20)
- Extism 0.4.0: [Install Extism](https://extism.org/docs/install)

## Creating the application

Start by creating a `go.mod` file with the command `go mod init go-host-application`, then a `main.go` file with the following content:

```golang
package main

import (
	"fmt"

	"github.com/extism/extism"
)

func main() {

	ctx := extism.NewContext()

    // This will free the context and all associated plugins
	defer ctx.Free() 

    // Path to the wasm file 0️⃣
	path := "../03-even-with-javascript/hello-js.wasm"
    
    // Define the path to the wasm file 1️⃣
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		}}

    // Load the wasm plugin 2️⃣
	plugin, err := ctx.PluginFromManifest(
        manifest, 
        []extism.Function{}, 
        true,
    )

	if err != nil {
		panic(err)
	}

    // Call the `say_hello` function 3️⃣
    // with a string parameter
	res, err := plugin.Call(
		"say_hello",
		[]byte("👋 Hello from the Go Host app 🤗"),
	)

	if err != nil {
		fmt.Println("😡", err)
	} else {
        // Display the return value 4️⃣
		fmt.Println("🙂", string(res))
	}

}
```

You see, the code is very very simple:

- 0: let's use the JavaScript Wasm plugin that we developed in the previous article.
- 1: define a manifest with properties including the path to access the Wasm file.
- 2: load the Wasm plugin.
- 3: call the `say_hello` function of the plugin.
- 4: display the result (the type of `res` is `[]byte`).

### Run the program

Use simply this command:

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go
```
> You need to set the linker lookup path env var explicitly.


And you will get this:

```bash
🙂 param: 👋 Hello from the Go Host app 🤗
```

You can of course do the test with the first plugin developed with TinyGo. Change the value of the variable `	path := "../01-simple-go-plugin/simple.wasm"` and run again:

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go
```

And you should get this:

```bash
🙂 👋 Hello 👋 Hello from the Go Host app 🤗
```

🎉 you see, it's easy to create Go applications that can run Wasm plugins written in different languages.

If I can keep up the pace, tomorrow I'll explain how to do the same thing but this time with Node.js.

