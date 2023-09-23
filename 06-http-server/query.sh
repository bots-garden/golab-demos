#!/bin/bash
clear
bat $0 --line-range 5:
echo ""
curl -X POST \
http://localhost:8080 \
-H 'content-type: text/plain; charset=utf-8' \
-d '😄 Bob Morane'

