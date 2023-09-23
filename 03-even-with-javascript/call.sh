#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
minism call ./hello-js.wasm say_hello \
  --input "😀 Hello World 🌍! (from JavaScript)"
  