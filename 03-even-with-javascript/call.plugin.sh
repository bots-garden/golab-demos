#!/bin/bash
extism call ./hello-js.wasm \
  say_hello --input "😀 Hello GoLab! (from JavaScript)" \
  --wasi \
  --log-level info

echo ""