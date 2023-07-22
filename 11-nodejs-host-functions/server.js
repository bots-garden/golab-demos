import Fastify from 'fastify'
import process from "node:process"

import { Context, HostFunction, ValType } from '@extism/extism'
import { readFileSync } from 'fs'

let memoryMap = new Map()

memoryMap.set("hello", "👋 Hello World 🌍")
memoryMap.set("message", "I 💜 Extism 😍")


// Host function (callable by the WASM plugin)
function memoryGet(plugin, inputs, outputs, userData) {

  // Read the value of inputs from the memory
  let memKey = plugin.memory(inputs[0].v.i64)
  // memKey is a buffer, use toString() to get the string value
  
  // This is the return value
  const returnValue = memoryMap.get(memKey.toString())
  
  // Allocate memory
  // Copy the value into memory
  let offs = plugin.memoryAlloc(Buffer.byteLength(returnValue))
  plugin.memory(offs).write(returnValue)
  
  //console.log("👋", offs)
  outputs[0].v.i64 = offs 
}

// Host functions list
let hostFunctions = [
  new HostFunction(
    "hostMemoryGet",
    [ValType.I64],
    [ValType.I64],
    memoryGet,
    "",
  )
]

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
//let plugin = ctx.plugin(wasm, true, [])
let plugin = ctx.plugin(wasm, true, hostFunctions)
// Thanks to the event queue we can instantiate the plugin only one time
// https://www.geeksforgeeks.org/how-to-handle-concurrency-in-node-js/

// Create and start the HTTP server
const start = async () => {

  fastify.post('/', opts, async (request, reply) => {

    // Call the WASM function, the request body is the argument of the function
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

