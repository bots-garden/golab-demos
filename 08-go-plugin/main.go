package main

import (
	"github.com/extism/go-pdk"
)


//export hostRobotMessage
func hostRobotMessage(offset uint64) uint64

func RobotMessage(message string) {
	messageMemory := pdk.AllocateString(message)
	hostRobotMessage(messageMemory.Offset())
}

//export say_hello
func say_hello() {
	input := pdk.Input()
	
	RobotMessage("hello " + string(input))

}


func main() {}
