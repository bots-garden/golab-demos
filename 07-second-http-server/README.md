# First HTTP Server

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

```bash
Summary:
  Total:	7.8349 secs
  Slowest:	5.4285 secs
  Fastest:	0.1125 secs
  Average:	2.2750 secs
  Requests/sec:	38.2901
  
  Total data:	9900 bytes
  Size/request:	33 bytes

Response time histogram:
  0.113 [1]	|
  0.644 [22]	|■■■■■■■■■
  1.176 [77]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  1.707 [0]	|
  2.239 [1]	|
  2.771 [99]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  3.302 [49]	|■■■■■■■■■■■■■■■■■■■■
  3.834 [14]	|■■■■■■
  4.365 [20]	|■■■■■■■■
  4.897 [10]	|■■■■
  5.429 [7]	|■■■


Latency distribution:
  10% in 0.6721 secs
  25% in 0.8899 secs
  50% in 2.4412 secs
  75% in 2.8735 secs
  90% in 4.0453 secs
  95% in 4.3824 secs
  99% in 5.3329 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0007 secs, 0.1125 secs, 5.4285 secs
  DNS-lookup:	0.0003 secs, 0.0000 secs, 0.0028 secs
  req write:	0.0002 secs, 0.0001 secs, 0.0030 secs
  resp wait:	2.2740 secs, 0.1122 secs, 5.4264 secs
  resp read:	0.0001 secs, 0.0000 secs, 0.0010 secs

Status code distribution:
  [200]	300 responses
```