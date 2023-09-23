import Fastify from 'fastify'
import process from "node:process"

import { Context } from '@extism/extism'
import { readFileSync } from 'fs'

let wasmFile = "../01-simple-go-plugin/simple.wasm"
let functionName = "say_hello"
let httpPort = 7070

let wasm = readFileSync(wasmFile)

const fastify = Fastify({
  logger: true
})

const opts = {}


// Create the WASM plugin
let ctx = new Context()
let plugin = ctx.plugin(wasm, true, [])
// Thanks to the event queue we can instantiate the plugin only one time
// https://www.geeksforgeeks.org/how-to-handle-concurrency-in-node-js/

// Create and start the HTTP server
const start = async () => {

  fastify.post('/', opts, async (request, reply) => {

    // Call the WASM function, the request body is the argulent of the function
    let buf = await plugin.call(functionName, request.body); 
    let result = buf.toString()

    //return JSON.parse(result)
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

