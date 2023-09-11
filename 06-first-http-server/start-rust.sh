#!/bin/bash
#LD_LIBRARY_PATH=/usr/local/lib \
go run main.go \
../05-hello-rust-plugin/target/wasm32-wasi/release/hello_rust_plugin.wasm \
hello \
8080