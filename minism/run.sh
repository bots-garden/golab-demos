#!/bin/bash

go run main.go \
call ../01-simple-go-plugin/simple.wasm \
say_hello \
--input "Bob Morane" \
--log-level info \
--allow-hosts *,*.google.com,yo.com \

