#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run main.go \
../04-hello-rust-plugin/target/wasm32-wasi/release/hello_rust_plugin.wasm \
hello \
8080
