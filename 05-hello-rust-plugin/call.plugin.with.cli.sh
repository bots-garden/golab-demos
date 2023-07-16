#!/bin/bash
extism call ./target/wasm32-wasi/release/hello_rust_plugin.wasm \
  hello --input "Lisa" \
  --wasi \

echo ""

extism call ./target/wasm32-wasi/release/hello_rust_plugin.wasm \
  hello --input "Bob"  \
  --wasi \

echo ""
