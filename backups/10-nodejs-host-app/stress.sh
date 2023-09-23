#!/bin/bash
# -n 10000 -c 1000
hey -n 300 -c 100 -m POST \
-H "Content-Type: text/plain" \
-d "John Doe" \
"http://localhost:7070"
