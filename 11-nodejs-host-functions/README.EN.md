# Write a host function with the Extism Host SDK

In this article, we will:

- Modify the host application developed with Node.js in [the previous article](https://k33g.hashnode.dev/writing-wasm-microservices-with-nodejs-and-extism) to add a **host function** developed by us.
- Modify the Wasm plugin developed in Go in [the first article of the series](https://k33g.hashnode.dev/extism-webassembly-plugins) to use this **host function**.

## A little reminder about host functions

**Excerpt from ["Extism, WebAssembly Plugins & Host functions"](https://k33g.hashnode.dev/extism-webassembly-plugins-host-functions)** :

*It is possible for the host application to provide the guest (the Wasm module) with extra powers. We call this "host functions". It is a function developed "in the source code of the host". The host exposes (export) it to the Wasm module which will be able to execute it. For example, you can develop a host function to display messages and thus allow the Wasm module to display messages in a terminal during its execution...*

*... Extism's Plugin Development Kit (PDK) provides some ready-to-use host functions, especially for logging, HTTP requests or reading a configuration in memory.*

But with Extism's Host SDK, you can develop your own host functions. This can be useful for example for database access, interaction with MQTT or Nats brokers...

In this article, we will keep it simple and develop a host function that allows to retrieve messages stored in the memory of the host application from a key. We will use a [JavaScript Map](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map) for this.

Let's start by modifying our Node.js application.

## Development of the host functions

Modify the `server.js` file as follows:


```javascript
import Fastify from 'fastify'
import process from "node:process"

// 1️⃣
import { Context, HostFunction, ValType } from '@extism/extism'
import { readFileSync } from 'fs'

// 2️⃣
let memoryMap = new Map()

memoryMap.set("hello", "👋 Hello World 🌍")
memoryMap.set("message", "I 💜 Extism 😍")

// 3️⃣ Host function (callable by the WASM plugin)
function memoryGet(plugin, inputs, outputs, userData) { 

  // 4️⃣ Read the value of inputs from the memory
  let memKey = plugin.memory(inputs[0].v.i64)
  // memKey is a buffer, 
  // use toString() to get the string value
  
  // 5️⃣ This is the return value
  const returnValue = memoryMap.get(memKey.toString())
  
  // 6️⃣ Allocate memory
  let offs = plugin.memoryAlloc(Buffer.byteLength(returnValue))
  // 7️⃣ Copy the value into memory
  plugin.memory(offs).write(returnValue)
  
  // 8️⃣ return the position and the length for the wasm plugin
  outputs[0].v.i64 = offs 
}

// 9️⃣ Host functions list
let hostFunctions = [
  new HostFunction(
    "hostMemoryGet",
    [ValType.I64],
    [ValType.I64],
    memoryGet,
    "",
  )
]

// location of the new plugin
let wasmFile = "../12-simple-go-mem-plugin/simple.wasm"
let functionName = "say_hello"
let httpPort = 7070

let wasm = readFileSync(wasmFile)

const fastify = Fastify({
  logger: true
})

const opts = {}

// Create the WASM plugin
let ctx = new Context()

// 1️⃣0️⃣
let plugin = ctx.plugin(wasm, true, hostFunctions)

// Create and start the HTTP server
const start = async () => {

  fastify.post('/', opts, async (request, reply) => {

    // Call the WASM function, 
    // the request body is the argument of the function
    let buf = await plugin.call(functionName, request.body); 
    let result = buf.toString()

    return result
  })

  try {
    await fastify.listen({ port: httpPort, host: '0.0.0.0'})
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
start().then(r => console.log("😄 started"))
```

- 1: import `HostFunction` (which allows the host to define functions callable by the Wasm plugin) and `ValType` (an enumeration of the possible types used by the host function).
- 2: create and populate a JavaScript `Map`
- 3: define the host function `memoryGet`
- 4: when the host function is called by the Wasm plugin, the parameter passing is done using the shared memory between the plugin and the host. `plugin.memory(inputs[0].v.i64)` is used to fetch this information from the shared memory. `memKey` is a buffer that contains the key to find a value in the JavaScript `Map` (and we use `memKey.toString()` to transform the buffer into a string).
- 5: we get the value associated with the key.
- 6: we allocate memory to copy the value associated with the key. `offs` corresponds to the position and length of the value in memory (it is thanks to the bit-shifting method that we can "fit 2 values into one").
- 7: we copy the value `returnValue` into this memory at the indicated location `offs`.
- 8: we copy into the return variable `outputs` (passed to the function by reference) the value of `offs` which will allow the Wasm plugin to read in memory the result of the function.
- 9: we define an array of host functions. In our case we create only one, where `"hostMemoryGet"` will be the alias of the function "seen" by the Wasm plugin, `[ValType.I64]` represents the type of the input parameter and the type of the output parameter (we remember that Wasm functions only accept numbers - and in our case these numbers contain the positions and sizes of values in shared memory) and finally `memoryGet` which is the definition of our host function.
- 10: When instantiating the Wasm plugin, we pass as an argument the array of host functions.

Before we can run our HTTP server again, we will have to modify our Wasm plugin.

## Modify the Wasm plugin

```golang
package main

import (
	"strings"
	"github.com/extism/go-pdk"
)


//export hostMemoryGet // 1️⃣
func hostMemoryGet(x uint64) uint64

//export say_hello
func say_hello() int32 {

	// read function argument from the memory
	// this is the name passed to the function
	input := pdk.Input()

	// Call the host function
	// 2️⃣
	key1 := pdk.AllocateString("hello")
	// 3️⃣
	offs1 := hostMemoryGet(key1.Offset())

  // 4️⃣
	mem1 := pdk.FindMemory(offs1)
	/*
		mem1 is a struct instance
		type Memory struct {
			offset uint64
			length uint64
		}
	*/

	// 5️⃣
	buffMem1 := make([]byte, mem1.Length())
	mem1.Load(buffMem1)

	// 6️⃣ get the second message
	key2 := pdk.AllocateString("message")
	offs2 := hostMemoryGet(key2.Offset())
	mem2 := pdk.FindMemory(offs2)
	buffMem2 := make([]byte, mem2.Length())
	mem2.Load(buffMem2)

  // 7️⃣
	data := []string{
		"👋 Hello " + string(input),
		"key: hello, value: " + string(buffMem1),
		"key: message, value: " + string(buffMem2),
	}

	// Allocate space into the memory
	mem := pdk.AllocateString(strings.Join(data, "\n"))
	// copy output to host memory
	pdk.OutputMemory(mem)

	return 0
}

func main() {}
```

- 1: the function `hostMemoryGet` must be exported to be usable.
- 2: we want to call the host function to get the value corresponding to the key `hello`, so for that we have to copy this key into memory.
- 3: we call the host function `hostMemoryGet` (`key1.Offset()` represents the position and length in memory of the key `key1` into only one value).
- 4: `pdk.FindMemory(offs1)` allows to retrieve a structure `mem1` containing the position and length.
- 5: we can now create a buffer `buffMem1` with the size of the value to retrieve and load it with the content of the memory location (`mem1`). Then we just have to read the string with `string(buffMem1)`.
- 6: we repeat for reading the second key.
- 7: we build a slice of strings that we will then transform into a single string to return it to the host function.

> If you want to learn more about the shared memory between the host and the Wasm plugin, you can read this blog post: [WASI communication between Node.js and Wasm modules with the Wasm Buffer Memory](https://k33g.hashnode.dev/wasi-communication-between-nodejs-and-wasm-modules-with-the-wasm-buffer-memory)

### Compile the new plugin

To compile the program, use TinyGo and the command below, which will produce a file `simple.wasm` :

```bash
tinygo build -scheduler=none --no-debug \
  -o simple.wasm \
  -target wasi main.go
```

It's time to test our modifications.

## Start the server and call the MicroService

To start the server, simply use this command:

```bash
node server.js
```

Then, to call the MicroService, use this simple `curl` command :

```bash
curl -X POST http://localhost:7070 \
-H 'Content-Type: text/plain; charset=utf-8' \
-d 'Jane Doe'
```

And you will get the messages from each of the keys of the JavaScript `Map` :

```bash
👋 Hello Jane Doe
key: hello, value: 👋 Hello World 🌍
key: message, value: I 💜 Extism 😍
```

Remember that when the Wasm plugin calls the host function, it is not it that executes the processing, but rather the host application. In the case of Node.js, this will eventually slow down the execution of the plugin, because Node.js is generally slower than compiled Go. Nevertheless, the potential of host functions is very interesting.

😥 This article was a bit more complicated than the previous ones, but this concept of host functions is essential. These last two articles also show you how you can evolve your Node.js applications with other languages. Feel free to contact me for more explanations. My next article will also explain how to make host functions, but this time in Go.

