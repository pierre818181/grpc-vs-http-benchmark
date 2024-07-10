# Steps to benchmark:

```bash
    make install-proto
```

OR:

```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    export PATH="$PATH:$(go env GOPATH)/bin"
```

# Benchmark HTTP requests:
```bash
    make run-http
```

# Benchmark GRPC requests:
```bash
    make run-grpc
```

GIST: At high throughput, GRPC is quite fast even in a local server. Last time, I benchmarked with a ngrok https server in between, GRPC was >10 times faster even for a GET request. So most likely if we move from the current setup to GRPC, we will cut network latency 10x.# grpc-vs-http-benchmark