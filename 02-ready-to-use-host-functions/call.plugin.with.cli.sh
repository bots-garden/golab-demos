#!/bin/bash
extism call ./host-functions.wasm \
  say_hello --input "😀 Hello World 🌍! (from TinyGo)" \
  --wasi \
  --log-level info \
  --allow-host '*' \
  --config route=https://jsonplaceholder.typicode.com/todos/1 

echo ""
