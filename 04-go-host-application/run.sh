#!/bin/bash

go run main.go ../01-simple-go-plugin/simple.wasm say_hello "Bob Morane"

go run main.go ../03-even-with-javascript/hello-js.wasm say_hello "Bob Morane"

go run main.go ../05-hello-rust-plugin/target/wasm32-wasi/release/hello_rust_plugin.wasm hello "Bob Morane"


go run main.go ../02-ready-to-use-host-functions/host-functions.wasm say_hello https://jsonplaceholder.typicode.com/todos/1
go run main.go ../02-ready-to-use-host-functions/host-functions.wasm say_hello https://jsonplaceholder.typicode.com/todos/3

