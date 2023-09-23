#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
minism call hello.wasm hello \
--input "Bob Morane"

