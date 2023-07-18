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
