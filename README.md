# gRPC Project
A sample project to understand gRPC (HTTP/2) using Go language

## Install Proto Compiler (protoc) on MacOS
```shell
brew install protobuf-c
```

## Setup gPRC for GoLang
https://grpc.io/docs/languages/go/quickstart/

## Commands
### Compile proto files
```shell
protoc --go-grpc_out=. --go_out=. *.proto
```
