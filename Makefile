install-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    export PATH="$PATH:$(go env GOPATH)/bin"

run-2G-docker:
	docker run --memory 2g -v ./:/mnt -it golang:bullseye

run-http:
	@echo "Running HTTP server benchmark on port 8080. Running GET requests first and then POST requests."
	@go run http/main.go

run-grpc:
	@echo "Running gRPC server benchmark on port 8081"
	@go run grpc/main.go
