#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run main.go ../09-rust-plugin/target/wasm32-wasi/release/simple_rust_plugin.wasm \
hello "Bob Morane"

echo ""
