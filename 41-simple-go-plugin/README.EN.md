# Extism & WebAssembly Plugins

Extism is a set of SDK projects that allow you to develop applications that run WebAssembly plugins, as well as develop WebAssembly plugins themselves.

Extism provides multiple SDKs for building host applications (those that load and execute WebAssembly plugins) in various languages such as Rust, Go, Ruby, PHP, JavaScript, Java, Erlang, Haskell, Zig, .Net, C, Swift, and OCaml.

As mentioned earlier, Extism also provides Plug-in Development Kits (PDKs) for developing WebAssembly plugins in Go, Rust, Haskell, C, Zig, AssemblyScript, and JavaScript.

## Wasi

To execute WebAssembly plugins from non-browser host applications, you need to use the Wasi specification. WebAssembly runtimes like WasmEdge, WasmTime, Wazero, Wasmer, and others implement this specification and provide SDKs for creating host applications. However, the Wasi specification is still a work in progress and comes with some limitations.

For example, functions in a Wasm program can only accept numbers as parameters and can only return a single number. This means that using strings as parameters and return values is not straightforward. It's worth noting that by manipulating the shared memory between the host and the Wasm program, it is possible to work around this limitation (and it can be a valuable learning experience).

Another example of a limitation is that a Wasm function cannot make HTTP requests or write to the console (display a result). The workaround is to create functions in the host application to perform these operations and expose them to the Wasm module for it to use. Again, this is not a trivial task.

## Fortunately, we have Extism!

Extism provides all the necessary "plumbing" to overcome the limitations of the Wasi specification, making it easy to develop, for instance, a Go application that can execute Wasm plugins developed in Go, Rust, and even JavaScript (and other languages, of course).

Extism also comes with a CLI that allows you to test your plugins. So, today, we will focus solely on the development of WebAssembly plugins.

## Prerequisites

To reproduce the examples in this article, you will need:

- Go (v1.20) & TinyGo (v0.28.1)
- Node.js (v19.9.0)
- Extism 0.4.0 & Extism-js PDK for building Wasm modules with JavaScript
  - [Install Extism](https://extism.org/docs/install)
  - [Install Extism-js PDK](https://extism.org/docs/write-a-plugin/js-pdk#how-to-install-and-use-the-extism-js-pdk)

But let's see how to create our first Wasm plugin.

## First plugin in Go

Start by creating a `go.mod` file with the command `go mod init simple-go-plugin`, and then create a `main.go` file with the following content:

```golang
package main

import (
	"github.com/extism/go-pdk"
)

//export say_hello 1️⃣
func say_hello() int32 {

	// read function argument from the memory
	input := pdk.Input() //2️⃣

	output := "👋 Hello " + string(input)

	mem := pdk.AllocateString(output) //3️⃣
	// copy output to host memory
	pdk.OutputMemory(mem) //4️⃣

	return 0
}

func main() {}
```

**Remarks**:
- 1️⃣: The `//export say_hello` annotation is necessary for the `say_hello` function to be "visible" to the host application (which will be the Extism CLI).
- 2️⃣: `pdk.Input()` allows reading the shared memory between the Wasm module and the host application to extract a buffer (`[]byte`) containing the parameter sent by the host function.
- 3️⃣: Allocate memory for the return value.
- 4️⃣: Copy the value into memory (it will be usable by the host application).

### Compile the Wasm plugin

To compile the program, use TinyGo and the following command, which will produce a `simple.wasm` file:

```bash
tinygo build -scheduler=none --no-debug \
  -o simple.wasm \
  -target wasi main.go
```

### Execute the `say_hello` function of the Wasm plugin

To do this, we will use the Extism CLI (we will see how to develop our own host application in a future article).

To execute the `say_hello` function with the string parameter `"Lisa"`, use the following command:

```bash
extism call ./simple.wasm \
  say_hello --input "Lisa" \
  --wasi
```

And you will get the output:

```bash
👋 Hello Lisa
```

That's all for today. In the next articles, we will see:
- How to use the "ready-to-use" host functions provided by Extism
- How to create a Wasm plugin with JavaScript
- How to develop a host application in Go
- And probably more 😉
