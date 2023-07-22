package main

import (
	"strings"
	"github.com/extism/go-pdk"
)

//export hostMemoryGet
func hostMemoryGet(x uint64) uint64

//export say_hello
func say_hello() int32 {

	// read function argument from the memory
	// this is the name passed to the function
	input := pdk.Input()

	// Call the host function
	// 1- copy the key to the shared memory
	key1 := pdk.AllocateString("hello")
	// call the host function
	// key1.Offset() is the position and the length of key1 into the memory (2 values into only one value)
	// read https://k33g.hashnode.dev/wasi-communication-between-nodejs-and-wasm-modules-with-the-wasm-buffer-memory
	offs1 := hostMemoryGet(key1.Offset())
	// read the value into the memory
	// offs1 is the position and the length of the result (2 values into only one value)
	// get the length and the position of the result in memory
	mem1 := pdk.FindMemory(offs1)
	/*
		mem1 is a struct instance
		type Memory struct {
			offset uint64
			length uint64
		}
	*/

	// create a buffer from the mem1
	// fill the buffer with mem1
	buffMem1 := make([]byte, mem1.Length())
	mem1.Load(buffMem1)

	// get the second message
	key2 := pdk.AllocateString("message")
	offs2 := hostMemoryGet(key2.Offset())
	mem2 := pdk.FindMemory(offs2)
	buffMem2 := make([]byte, mem2.Length())
	mem2.Load(buffMem2)

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
