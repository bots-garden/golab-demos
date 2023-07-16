#!/bin/bash
extism call ./simple.wasm \
  say_hello --input "Lisa" \
  --wasi

echo ""

extism call ./simple.wasm \
  say_hello --input "Bob" \
  --wasi

echo ""

