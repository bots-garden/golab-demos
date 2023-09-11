package main

import (
	"github.com/extism/go-pdk"
)

//export hostGetString
func hostGetString() uint64

//export say_hello
func say_hello() {

	// read function argument from the memory
	// this is the name passed to the function
	//input := pdk.Input()

	// call the host function
	offset := hostGetString()
	// read the value into the memory
	// offset is the position and the length of the result (2 values into only one value)
	// get the length and the position of the result in memory
	memory := pdk.FindMemory(offset)
	/*
		mem1 is a struct instance
		type Memory struct {
			offset uint64
			length uint64
		}
	*/

	// create a buffer from the mem1
	// fill the buffer with mem1
	memoryBuffer := make([]byte, memory.Length())
	memory.Load(memoryBuffer) // the buffer contains "I 💜 Extism"

	// Allocate space into the memory
	mem := pdk.AllocateString("👋 The message is: " + string(memoryBuffer))
	// copy output to host memory
	pdk.OutputMemory(mem)

}

func main() {}
