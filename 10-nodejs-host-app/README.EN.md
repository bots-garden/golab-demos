# Writing Wasm MicroServices with Node.js and Extism

> How to write a host application with Node.js

With the help of Extism, writing a host application (i.e. an application capable of running WebAssembly plugins) is rather easy. We have seen in a [previous article](https://k33g.hashnode.dev/run-extism-webassembly-plugins-from-a-go-application) how to do it in Go. Today we will do it with Node.js. You will see that it is very simple, but this example will allow us to go further and discover how to write host functions.

**This application is an HTTP server that will serve the WebAssembly plugin as a MicroService**. And we will use [the WebAssembly plugin developed with TinyGo](https://k33g.hashnode.dev/extism-webassembly-plugins) that we did in a previous article.

## Prerequisites

You will need

- Go (v1.20) and TinyGo (v0.28.1) to compile the plugins
- Extism 0.4.0: [Install Extism](https://extism.org/docs/install)
- Node.js (v19.9.0) (this is the version I use)

## Creating the application

### Installing dependencies

In a directory, create a file `package.json` with the following content:

```json
{
  "dependencies": {
    "@extism/extism": "^0.4.0",
    "fastify": "^4.20.0"
  },
  "type": "module"
}
```

> **[Fastify](https://fastify.dev/)** is a Node.js project that allows you to develop web application servers (server-side), like [Express.js](https://expressjs.com/). But you can use whatever you want.

Then, type the command below to install the necessary dependencies:

```bash
npm install
```

### Developing the application

Create a file `server.js` with the following content:

```javascript
import Fastify from 'fastify'
import process from "node:process"

import { Context } from '@extism/extism'
import { readFileSync } from 'fs'

let wasmFile = "../01-simple-go-plugin/simple.wasm"
let functionName = "say_hello"
let httpPort = 7070

let wasm = readFileSync(wasmFile) // 1️⃣

const fastify = Fastify({
  logger: true
})

const opts = {}


// 2️⃣
let ctx = new Context()
let plugin = ctx.plugin(wasm, true, [])

// Create and start the HTTP server
const start = async () => {

  fastify.post('/', opts, async (request, reply) => { // 3️⃣

    // 4️⃣
    let buf = await plugin.call(functionName, request.body); 
    let result = buf.toString()

    return result
  })

  try { // 5️⃣
    await fastify.listen({ port: httpPort, host: '0.0.0.0'})
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
start().then(r => console.log("😄 started"))
```

- 1: load the WebAssembly plugin file.
- 2: create an Extism context and use it to initialize the Wasm plugin.
- 3: define a route (endpoint) of the HTTP server. The code will be executed on each HTTP POST call of `http://localhost:7070`.
- 4: call the function of the module with as parameters the name of the function (`say_hello`) and the data posted by the HTTP request, and return the result.
- 5: start the HTTP server.

### Start the HTTP server

Simply use this command:

```bash
node server.js
```

### Call the MicroService

To call the MicroService, use this simple `curl` command:

```bash
curl -X POST http://localhost:7070 \
-H 'Content-Type: text/plain; charset=utf-8' \
-d 'Jane Doe'
```

And you will get:

```bash
👋 Hello Jane Doe
```

😍 You see that with the plugin system proposed by Extism, it becomes very easy to write polyglot MicroServices and to offer them using Node.js. You are not far from having the basics to write a FaaS (but that will be another story, probably later 😉). I'll let you experiment with what we've seen today.

The next article to come will reuse this example, and I will remind the concept of **host function** and how to write them to bring additional features to WebAssembly plugins.
