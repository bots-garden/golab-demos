#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
go run main.go ../08-go-plugin/simple.wasm \
say_hello "Bob Morane"

echo ""
