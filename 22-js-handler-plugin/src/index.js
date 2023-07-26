import { setHandler } from "./core/receiver"

function handle() {
	
	setHandler(param => {
		let output = "param: " + param
		let err = null

		return [output, err]
	})
}

module.exports = {handle}
