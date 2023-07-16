# First HTTP Server: be careful to the shared memory

```bash
curl -v -X POST \
http://localhost:8080 \
-H 'content-type: text/plain; charset=utf-8' \
-d '😄 Bob Morane'
```

## Load testing

```bash
hey -n 300 -c 100 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

👋 Explain why it fails
