# WASM Microservices with Extism and Fiber

Today, I will quickly show you how to serve Extism plugins (so Webassembly plugins) with the excellent framework [Fiber](https://docs.gofiber.io/). **Fiber** is a web framework for making HTTP servers with a similar spirit to Node.js frameworks like [Express](https://expressjs.com/) (which I have used many times in the past) or [Fastify](https://fastify.dev/).

This article will be "slightly" longer than the previous ones, because I also want to talk to you about my mistakes during my learning with Wasi.

## Prerequisites

- At best: have read all the blog posts in this series ["Discovery of Extism (The Universal Plug-in System)"](https://k33g.hashnode.dev/series/extism-discovery)
- At a minimum:
  - [Extism & WebAssembly Plugins](https://k33g.hashnode.dev/extism-webassembly-plugins)
  - [Run Extism WebAssembly plugins from a Go application](https://k33g.hashnode.dev/run-extism-webassembly-plugins-from-a-go-application)
  - [Create a Webassembly plugin with Extism and Rust](https://k33g.hashnode.dev/create-a-webassembly-plugin-with-extism-and-rust)

## Creating an HTTP server as a host application

Start by creating a `go.mod` file with the command `go mod init first-http-server`, then a `main.go` file with the following content:

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/extism/extism"
	"github.com/gofiber/fiber/v2"
)

func main() {
    // Parameters of the program 0️⃣
	wasmFilePath := os.Args[1:][0]
	wasmFunctionName := os.Args[1:][1]
    httpPort := os.Args[1:][2]

	ctx := extism.NewContext()

	defer ctx.Free()

    // Define the path to the wasm file 1️⃣
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
	}

    // Load the wasm plugin 2️⃣
	plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
	if err != nil {
		panic(err)
	}

    // Create an instance of Fiber application 3️⃣
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

    // Create a route "/" and a handler to call the wasm function 4️⃣
	app.Post("/", func(c *fiber.Ctx) error {

		params := c.Body()

        // Call the wasm function 5️⃣
        // with a string parameter
		out, err := plugin.Call(wasmFunctionName, params)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		} else {
            // Send the HTTP response to the client 6️⃣
			c.Status(http.StatusOK)
			return c.SendString(string(out))
		}

	})

    // Start the HTTP server 7️⃣
	fmt.Println("🌍 http server is listening on:", httpPort)
	app.Listen(":" + httpPort)
}
```

If you have read the previous articles, part of the code is already familiar to you.

- 0: use the program parameters to pass it the following information: the path of the wasm plugin, the name of the function to call and the HTTP port.
- 1: define a manifest with properties including the path to access the Wasm file.
- 2: load the Wasm plugin.
- 3: create a Fiber application.
- 4: create a route "/" that will be triggered by an HTTP request of type `POST`.
- 5: call the function of the plugin.
- 6: return the result (the HTTP response).
- 7: start the server.

### Start the server and serve the WASM plugin

We will use the wasm plugin developed in Rust from our previous article [Create a Webassembly plugin with Extism and Rust](https://k33g.hashnode.dev/create-a-webassembly-plugin-with-extism-and-rust).

Start the application as follows:

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go \
path_to_the_plugin/hello.wasm \
hello \
8080
```

You should get this:

```bash
🌍 http server is listening on: 8080
```

And now, make an HTTP request:

```bash
curl -X POST \
http://localhost:8080 \
-H 'content-type: text/plain; charset=utf-8' \
-d '😄 Bob Morane'
echo ""
```

And, you should get this:

```bash
{"message":"🦀 Hello 😄 Bob Morane"}
```

### Stressing the application ... It's a disaster!

I always check the behavior of my web services by "stressing" them with the utility [Hey](https://github.com/rakyll/hey) which is extremely easy to use (especially with CI jobs for example, to check the performance before and after changes).

So I'm going to stress my service for the first time with the following command:

```bash
hey -n 300 -c 1 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

So I'm going to make 300 HTTP requests to my service with 1 single connection.

And I'm going to get a report like this (it's just an excerpt):

```bash
Summary:
  Total:        0.0973 secs
  Slowest:      0.0125 secs
  Fastest:      0.0001 secs
  Average:      0.0003 secs
  Requests/sec: 3082.8745
  
  Total data:   9900 bytes
  Size/request: 33 bytes

Status code distribution:
  [200] 300 responses
```
> I work with a Mac M1 Max


Now, let's check the behavior of the service with multiple connections at the same time:

```bash
hey -n 300 -c 100 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

So I'm going to make 300 HTTP requests to my service with 100 simultaneous connections.

And this time my HTTP server will crash! And in the load test report you will see that most of the requests are in error (it's just an excerpt):

```bash
Status code distribution:
  [200] 3 responses

Error distribution:
  [3]   Post "http://localhost:8080": EOF
  [196] Post "http://localhost:8080": dial tcp 127.0.0.1:8080: connect: connection refused
  [1]   Post "http://localhost:8080": read tcp 127.0.0.1:38552->127.0.0.1:8080: read: connection reset by peer
  [1]   Post "http://localhost:8080": read tcp 127.0.0.1:38568->127.0.0.1:8080: read: connection reset by peer
```

### But what happened?

In the paragraph ["WASI" of the first article of the series](https://k33g.hashnode.dev/extism-webassembly-plugins#heading-wasi), I explained that the way to exchange values other than numbers between the host application and the wasm plugin "guest" is to use the shared webassembly memory.

> I encourage you to read this excellent article on the subject [A practical guide to WebAssembly memory](https://radu-matei.com/blog/practical-guide-to-wasm-memory/) by [Radu Matei](https://twitter.com/matei_radu) (CTO at [FermyonTech](https://twitter.com/fermyontech)).


> You can also read this one, written by yours truly : [WASI, Communication between Node.js and WASM modules with the WASM buffer memory](https://k33g.hashnode.dev/wasi-communication-between-nodejs-and-wasm-modules-with-the-wasm-buffer-memory)

But let's get back to our problem. In fact it is very simple : there were 100 connections trying simultaneously to access this shared memory and so there was a "collision", because this memory is meant to be shared between the host application and only one "guest" at a time.

We therefore need to solve this problem to make our application really usable.

## Creating a second HTTP server, the "naive" solution

My first approach was to move the loading of the plugin from the manifest and its instantiation inside the HTTP handler to guarantee that for a given request there will be only one access to the shared memory:

```golang
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/extism/extism"
	"github.com/gofiber/fiber/v2"
)

func main() {
	wasmFilePath := os.Args[1:][0]
	wasmFunctionName := os.Args[1:][1]
	httpPort := os.Args[1:][2]

	ctx := extism.NewContext()

	defer ctx.Free() // this will free the context and all associated plugins

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
	}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/", func(c *fiber.Ctx) error {

		params := c.Body()

        // Load the wasm plugin 1️⃣
		plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		}

        // Call the wasm function 2️⃣
        // with a string parameter
		out, err := plugin.Call(wasmFunctionName, params)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		} else {
			c.Status(http.StatusOK)
			return c.SendString(string(out))
		}

	})

	fmt.Println("🌍 http server is listening on:", httpPort)
	app.Listen(":" + httpPort)
}
```

- 1: load the Wasm plugin.
- 2: call the function of the plugin.

So I launched my new HTTP server:

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go \
path_to_the_plugin/hello.wasm \
hello \
8080
```

And I did some load tests:

```bash
hey -n 300 -c 100 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

And I got this report:

```bash
Summary:
  Total:        7.6182 secs
  Slowest:      4.6650 secs
  Fastest:      0.0857 secs
  Average:      2.0480 secs
  Requests/sec: 39.3794

Status code distribution:
  [200] 300 responses
```

So, it's great, everything works! 🎉 But, however, the number of requests per second seems really small. Less than 40 requests per second, compared to the 3000 requests per second of the first test, it's ridiculous 😞. But, at least my application works.

But never hesitate to ask for help (that's why Open Source is a fabulous model).

## Creating a third HTTP server, the "smart" solution

I was still annoyed by the poor performance of my MicroService. I had made a similar application with Node.js (remember the following article: [Writing Wasm MicroServices with Node.js and Extism](https://k33g.hashnode.dev/writing-wasm-microservices-with-nodejs-and-extism#heading-developing-the-application)) and the load tests gave me 1800 requests per second.

And with the Node.js version of the application, the wasm plugin was instantiated [only once](https://k33g.hashnode.dev/writing-wasm-microservices-with-nodejs-and-extism#heading-developing-the-application) and I had no memory collision problem 🤔.

This should have put me on the track, because indeed, Node.js applications use a "Single Threaded Event Loop Model" unlike Fiber which uses a "Multi-Threaded Request-Response" architecture to handle concurrent accesses. So that's why my Node.js application doesn't "crash".

It was [Steve Manuel](https://twitter.com/nilslice) (CEO of [Dylibso](https://twitter.com/dylibso), but also the creator of Extism) who gave me the solution when I explained my problem to him during a discussion on Discord:

***"So if you want thread-saftey in Go HTTP handlers re-using plugins, you need to protect them with a mutex"***

In fact, yes, it was so obvious (and also an opportunity to start studying what a mutex was).

So I followed Steve's advice, and I modified my code as follows:

```golang
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/extism/extism"
	"github.com/gofiber/fiber/v2"
)

// Store all your plugins in a normal Go hash map, 
// protected by a Mutex 1️⃣
var m sync.Mutex
var plugins = make(map[string]extism.Plugin)

// Store the plugin 2️⃣
func StorePlugin(plugin extism.Plugin) {
	plugins["code"] = plugin
}

// Retrieve the plugin 3️⃣
func GetPlugin() (extism.Plugin, error) {
	if plugin, ok := plugins["code"]; ok {
		return plugin, nil
	} else {
		return extism.Plugin{}, errors.New("🔴 no plugin")
	}
}

func main() {
	wasmFilePath := os.Args[1:][0]
	wasmFunctionName := os.Args[1:][1]
	httpPort := os.Args[1:][2]

	ctx := extism.NewContext()

	defer ctx.Free()

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
	}

    // Create an instance of the plugin 4️⃣
	plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
	if err != nil {
		log.Println("🔴 !!! Error when loading the plugin", err)
		os.Exit(1)
	}
    // Save the plugin in the map 5️⃣
	StorePlugin(plugin)
	
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/", func(c *fiber.Ctx) error {

		params := c.Body()

        // Lock the mutex 6️⃣
		m.Lock()
		defer m.Unlock()
		
        // Get the plugin 7️⃣
		plugin, err := GetPlugin()

		if err != nil {
			log.Println("🔴 !!! Error when getting the plugin", err)
			c.Status(http.StatusInternalServerError)
			return c.SendString(err.Error())
		}
		
		out, err := plugin.Call(wasmFunctionName, params)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		} else {
			c.Status(http.StatusOK)
			return c.SendString(string(out))
		}

	})

	fmt.Println("🌍 http server is listening on:", httpPort)
	app.Listen(":" + httpPort)
}
```

- 1: create a map protected by a mutex. This map will be used to "protect" the wasm plugin.
- 2: create a function to save the plugin in the map.
- 3: create a function to retrieve the plugin from the map.
- 4: create an instance of the wasm plugin.
- 5: save this instance in the map.
- 6: lock the mutex and use `defer` to unlock it at the end of execution.
- 7: get the plugin from the protected map.

Once this modification was done, I launched my HTTP server again:

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go \
path_to_the_plugin/hello.wasm \
hello \
8080
```

And I ran some load tests again:

```bash
hey -n 300 -c 100 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

And I got this:

```bash
Summary:
  Total:        0.0365 secs
  Slowest:      0.0280 secs
  Fastest:      0.0001 secs
  Average:      0.0092 secs
  Requests/sec: 8207.9604
  
  Total data:   9900 bytes
  Size/request: 33 bytes

Status code distribution:
  [200] 300 responses
```

The new version of the HTTP server could handle up to more than 8000 requests per second! 🚀

Not bad, right? That's all for today. A huge thank you to Steve Manuel for his help. I learned a lot, because I dared to ask for help. So, when you struggle with something and you can't find a solution, don't hesitate to ask around.

See you soon for the next article. 👋
