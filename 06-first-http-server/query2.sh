#!/bin/bash
curl -X POST \
http://localhost:8080 \
-H 'content-type: text/plain; charset=utf-8' \
-d 'https://jsonplaceholder.typicode.com/todos/3'
echo ""

curl -X POST \
http://localhost:8080 \
-H 'content-type: text/plain; charset=utf-8' \
-d 'https://jsonplaceholder.typicode.com/todos/2'
echo ""
