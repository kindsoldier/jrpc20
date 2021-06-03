
## JSON RPC 2.0 implementation

### Function call
```
$ curl -H 'Content-Type: application/json' -X POST -d '{"method": "add", "params": {"first": 2, "second": 2} }' http://localhost:8090/fc
{
    "jsonrpc":"2.0",
    "result":4
}
```
### Function calls with wrong request schema
```
# curl -H 'Content-Type: application/json' -X POST -d '{"method": "add" }' http://localhost:8090/fc
{
    "jsonrpc":"2.0",
    "error":{
        "code":-32700,
        "message":"validation error: (root):params is required"
    }
}
```

```
$ curl -H 'Content-Type: application/json' -X POST -d '{"method": "add", "params": {"first": 2 } }' http://localhost:8090/fc
{
    "jsonrpc":"2.0",
    "error":{
        "code":-32700,
        "message":"validation error: (root).params:second is required"
    }
}
```
### Benachmarks over HTTP

```
$ ab -T 'application/json' -p add-body -n 10000 -c10 http://localhost:8090/fc
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>

Server Hostname:        localhost
Server Port:            8090

Document Path:          /fc
Document Length:        41 bytes

Concurrency Level:      10
Time taken for tests:   0.765 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1490000 bytes
Total body sent:        2010000
HTML transferred:       410000 bytes
Requests per second:    13079.91 [#/sec] (mean)
Time per request:       0.765 [ms] (mean)
Time per request:       0.076 [ms] (mean, across all concurrent requests)
Transfer rate:          1903.23 [Kbytes/sec] received
                        2567.44 kb/s sent
                        4470.67 kb/s total
```

```
$ ab -T 'application/json' -p add-body -n 10000 -c10 http://localhost:8090/fc
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>

Server Hostname:        localhost
Server Port:            8090

Document Path:          /fc
Document Length:        41 bytes

Concurrency Level:      10
Time taken for tests:   1.170 seconds
Complete requests:      10000
Failed requests:        0
Total transferred:      1490000 bytes
Total body sent:        2010000
HTML transferred:       410000 bytes
Requests per second:    8549.56 [#/sec] (mean)
Time per request:       1.170 [ms] (mean)
Time per request:       0.117 [ms] (mean, across all concurrent requests)
Transfer rate:          1244.03 [Kbytes/sec] received
                        1678.18 kb/s sent
                        2922.21 kb/s total
```

### Benchmark over web socket

```
$ (for f in `seq 1 1000`; do echo '{"method":"add","params":{"first":2,"second":2}}';done) | (time websocat ws://localhost:8090/ws) | tail -5

{"jsonrpc":"2.0","result":4}
{"jsonrpc":"2.0","result":4}
{"jsonrpc":"2.0","result":4}
{"jsonrpc":"2.0","result":4}
{"jsonrpc":"2.0","result":4}

real	0m0.069s
user	0m0.017s
sys	0m0.012s

```
