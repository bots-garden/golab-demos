# WebAssembly Plugin in JavaScript with Extism

In the last few days we have used the Go PDK (Plugin Development Kit) to develop WebAssembly applications and run them with the Extism CLI.

As I said in the first article, there are PDKs for several languages, including [JavaScript](https://github.com/extism/js-pdk) (and I'm a big fan of JavaScript). But how is this possible? Indeed, we can't compile JavaScript to native Wasm code.

In fact, this PDK uses (among other things) the [QuickJS](https://bellard.org/quickjs/) project to run JavaScript code in a Wasm program.

Of course, it won't run as fast as a Wasm program compiled with TinyGo or Rust, but it allows you to run JavaScript functions in a completely "sandboxed" environment. Like Shopify does (this PDK is a fork of the [Javy](https://github.com/bytecodealliance/javy) project initiated by [Shopify](https://www.shopify.com/).

> Shopify developed the Javy project to bring JavaScript support to Shopify Functions. Shopify Functions allow developers to create custom extensions and features for the specific needs of merchants in JavaScript, which is a popular and familiar language for many web developers.

So, imagine that you decide to create a FaaS platform oriented JavaScript and that you want to give your users the possibility to create and publish their own functions, going through a similar mechanism will have at least two advantages:

- Encourage adoption (JavaScritp is well known)
- Guarantee the integrity of your platform (functions are executed in a sandboxed environment)

> Some reading: [Bringing Javascript to WebAssembly for Shopify Functions](https://shopify.engineering/javascript-in-webassembly-for-shopify-functions)

But let's get back to our sheep and let me explain how to create an Extism plugin in JavaScript.

## Prerequisites

You will need Extism 0.4.0 and Extims-js PDK to build Wasm modules with JavaScript:

  - [Install Extism](https://extism.org/docs/install)
  - [Install Extism-js PDK](https://extism.org/docs/write-a-plugin/js-pdk#how-to-install-and-use-the-extism-js-pdk)

> **Note**, if you encounter a problem installing the PDK, you can do it manually this way (modify according to your environment):
> ```bash
> export TAG="v0.5.0"
> export ARCH="aarch64"
> export  OS="linux"
> curl -L -O "https://github.com/extism/js-pdk/releases/download/$TAG/extism-js-$ARCH-$OS-$TAG.gz"
> gunzip extism-js*.gz
> sudo mv extism-js-* /usr/local/bin/extism-js
> chmod +x /usr/local/bin/extism-js
> ```

## The simplest Extism plugin

In a directory, create a file `index.js` with the following content:

```javascript
function say_hello() {

	// read function argument from the memory
	let input = Host.inputString()

	let output = "param: " + input

	console.log("👋 Hey, I'm a JS function into a wasm module 💜")

	// copy output to host memory
	Host.outputString(output)

	return 0
}

module.exports = {say_hello}
```

You can see that the code is very simple and has a completely similar logic to what we have seen in the previous articles.


### Compile the Wasm plugin

To compile the program, use the command below, which will produce a file `hello-js.wasm`:

```bash
extism-js index.js -o hello-js.wasm
```

And now we are going to run our Wasm plugin as we did in the previous examples.

### Run the `say_hello` function of the Wasm plugin

For this we will use the Extism CLI.

To run the `say_hello` function with the parameter `"😀 Hello World 🌍! (from JavaScript)"`, use the following command:

```bash
extism call ./hello-js.wasm \
  say_hello --input "😀 Hello World 🌍! (from JavaScript)" \
  --wasi \
  --log-level info
```

You will get:

```bash
extism_runtime::pdk INFO 2023-07-18T07:08:34.347325607+02:00 - 👋 Hey, I'm a JS function into a wasm module 💜
param: 😀 Hello World 🌍! (from JavaScript)
```

If you try to run the plugin without specifying the log level, like this:

```bash
extism call ./hello-js.wasm \
  say_hello --input "😀 Hello World 🌍! (from JavaScript)" \
  --wasi
```

you will only get this:

```bash
param: 😀 Hello World 🌍! (from JavaScript)
```

> In the case of the JavaScript PDK `console.log()` is a kind of alias for calling the host log function of the Extism CLI.


That's all for today. I'll let you familiarize yourself with the development of Extism plugins. If all goes well, tomorrow we'll tackle writing a host application in Go.