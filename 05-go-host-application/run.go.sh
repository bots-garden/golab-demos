#!/bin/bash
clear
bat $0 --line-range 10:
echo ""
#go run main.go ../01-simple-go-plugin/simple.wasm \
#say_hello "Bob Morane" \
#'{"firstName":"Jane","lastName":"Doe"}'

#echo ""
# args: wasm_file function_name config
./hostapp ../02-ready-to-use-host-functions/host-functions.wasm \
say_hello Bob \
'{"route": "https://jsonplaceholder.typicode.com/todos/3"}'

echo ""
