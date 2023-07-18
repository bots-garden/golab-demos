#!/bin/bash
extism call ./hello-js.wasm \
  say_hello --input "😀 Hello World 🌍! (from JavaScript)" \
  --wasi \
  --log-level info

echo ""