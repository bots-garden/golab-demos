package main

import (
	"github.com/extism/go-pdk"
	"github.com/valyala/fastjson"
)

//export say_hello
func say_hello() int32 {

	// read function argument from the memory
	input := pdk.Input()

	// https://jsonplaceholder.typicode.com/todos/3
	url := string(input)

	// use request (host function)
	req := pdk.NewHTTPRequest("GET", url)
	res := req.Send()

	parser := fastjson.Parser{}
	jsonValue, _ := parser.Parse(string(res.Body()))
	title := string(jsonValue.GetStringBytes("title"))

	output := "🌍 url: " + string(input) + " 📝 title: " + title

	mem := pdk.AllocateString(output)
	// copy output to host memory
	pdk.OutputMemory(mem)

	return 0
}

func main() {}
