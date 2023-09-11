function say_hello() {

	// read function argument from the memory
	let input = Host.inputString()

	let output = "👋 (From JS) Hello " + input

	// copy output to host memory
	Host.outputString(output)

	return 0
}

module.exports = {say_hello}
