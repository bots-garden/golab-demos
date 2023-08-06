package main

import (
	"github.com/extism/go-pdk"
)
//export hostPrintMessage
func hostPrintMessage(offset uint64) uint64

func printMessage(message string) {
	messageMemory := pdk.AllocateString(message)
	hostPrintMessage(messageMemory.Offset())
}
//export hostDisplayMessage
func hostDisplayMessage(offset uint64) uint64

func displayMessage(message string) {
	messageMemory := pdk.AllocateString(message)
	hostDisplayMessage(messageMemory.Offset())
}

//export say_hello
func say_hello() int32 {
	input := pdk.Input()
	output := "👋 Hello " + string(input)

	printMessage("from say_hello")

	mem := pdk.AllocateString(output)
	pdk.OutputMemory(mem)
	return 0

}

//export say_hey
func say_hey() int32 {
	input := pdk.Input()
	output := "🫱 Hey " + string(input)

	displayMessage("from say_hey")

	mem := pdk.AllocateString(output)
	pdk.OutputMemory(mem)
	return 0

}

func main() {}
