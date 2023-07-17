# Extism, WebAssembly Plugins & Host Functions

As I explained in the previous article, Wasm programs are "limited" by default (it's also a security principle). The WASI specification is a set of APIs that will allow the development of WebAssembly (Wasm) programs to access system resources (if permitted by the host application). Currently, the WASI specification is still being written, and only a few APIs are available. Therefore, even though it's possible to use the socket or FileSystem API, the capabilities of a Wasm program in terms of accessing system resources are limited: no terminal display, no HTTP access, and so on.

## We are saved, we have host functions

However, to make our lives easier, the host application can provide additional powers to the guest (the Wasm module). We call these "host functions." A host function is a function developed within the host application's source code. The host application exposes it (exports it) to the Wasm module, which can then execute it. For example, you can develop a host function to display messages and allow the Wasm module to display messages in a terminal during its execution.

> **Note:** You should note that when you use host functions, your Wasm module will only be executable by your host application.

Yesterday, I explained that the type system in the WASI specification for passing function parameters and return values is very limited (only numbers are supported). This implies some "acrobatics" to develop a host function.

## Extism provides ready-to-use host functions

To help you develop Wasm programs without worrying about complexity, the Extism Plugin Development Kit (PDK) provides some ready-to-use host functions, including logging, HTTP requests, and reading an in-memory configuration.

### Creating a new Wasm plugin (with the Extism PDK)

Start by creating a `go.mod` file using the command `go mod init ready-to-use-host-functions`, and then create a `main.go` file with the following content:

```golang
package main

import (
	"github.com/extism/go-pdk"
	"github.com/valyala/fastjson"
)

//export say_hello
func say_hello() int32 {

	// read function argument from memory
	input := pdk.Input()

    // 1️⃣ write information to the logs
	pdk.Log(pdk.LogInfo, "👋 hello this is wasm 💜") 

    // 2️⃣ get the value associated with the `route` key 
    // in the config object
	route, _ := pdk.GetConfig("route")
    // the value of `route` is
    // https://jsonplaceholder.typicode.com/todos/1

    // 3️⃣ write information to the logs
	pdk.Log(pdk.LogInfo, "🌍 calling "+route)

    // 4️⃣ make an HTTP request
	req := pdk.NewHTTPRequest("GET", route)
	res := req.Send()
	
    // Read the result of the request
	parser := fastjson.Parser{}
	jsonValue, _ := parser.Parse(string(res.Body()))
	title := string(jsonValue.GetStringBytes("title"))

    // Prepare the return value
	output := "param: " + string(input) + " title: " + title

	mem := pdk.AllocateString(output)
	// copy output to host memory
	pdk.OutputMemory(mem)

	return 0
}

func main() {}
```

- 1: `pdk.Log` is a host function provided by the Extism CLI, allowing sending messages to the logs.
- 2: The `pdk.GetConfig` host function allows reading values from a configuration passed in memory by the host application. In this example, we retrieve a URL to use for an HTTP request.
- 3: We use `pdk.Log` again to send the value associated with the `route` configuration key to the logs.
- 4: The `pdk.NewHTTPRequest` host function allows making HTTP requests.

When you develop your own applications with the Extism SDK, they will also provide these same host functions (it's also possible to develop custom host functions, but that's for later).

> **Note:** TinyGo has built-in JSON serialization/deserialization support, but I continue to use `fastjson` for its faster performance in my use cases.

Now, let's test our new Wasm plugin.

### Compiling the Wasm plugin

To compile the program, use TinyGo and the following command, which will produce a `host-functions.wasm` file:

```bash
tinygo build -scheduler=none --no-debug \
  -o host-functions.wasm \
  -target wasi main.go
```

### Executing the `say_hello` function of the Wasm plugin

For this, we will use the Extism CLI (we will see how to develop our own host application in a future article).

To execute the `say_hello` function with the string parameter `"😀 Hello World 🌍! (from TinyGo)"`, use the following command:

```bash
extism call ./host-functions.wasm \
  say_hello --input "😀 Hello World 🌍! (from TinyGo)" \
  --wasi \
  --log-level info \
  --allow-host '*' \
  --config route=https://jsonplaceholder.typicode.com/todos/1 
```

- To display the logs, you need to specify the log level with `--log-level info`.
- To allow the Wasm module to make an HTTP request "outside," you need to give it permission by specifying `--allow-host '*'`.
- Finally, with the `--config` flag, you can "push" configuration information to the Wasm program.

You will get the following output:

```bash
extism_runtime::pdk INFO 2023-07-17T06:45:27.063609583+02:00 - 👋 hello this is wasm 💜
extism_runtime::pdk INFO 2023-07-17T06:45:27.063691542+02:00 - 🌍 calling https://jsonplaceholder.typicode.com/todos/1
param: 😀 Hello World 🌍! (from TinyGo) title: delectus aut autem
```

🎉 And there you have it! This concludes the second Extism discovery article.
👋 See you soon for the next article on how to create a JavaScript plugin.