#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
minism call host-functions.wasm say_hello \
--input "😀 Hello World 🌍! (from TinyGo)" \
--log-level info \
--allow-hosts '["*"]' \
--config '{"route":"https://jsonplaceholder.typicode.com/todos/1"}'

