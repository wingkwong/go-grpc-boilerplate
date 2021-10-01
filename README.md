# go-grpc-boilerplate
A gRPC microservice boilerplate written in Go

## Installation

Install [protoc](http://google.github.io/proto-lens/installing-protoc.html)

```bash
PROTOC_ZIP=protoc-3.14.0-osx-x86_64.zip
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/$PROTOC_ZIP
sudo unzip -o $PROTOC_ZIP -d /usr/local bin/protoc
sudo unzip -o $PROTOC_ZIP -d /usr/local 'include/*'
rm -f $PROTOC_ZIP
```

### Optional

Install ``vscode-proto3`` in Visual Studio Code

## Project Layout

Refer to [golang-standards/project-layout](https://github.com/golang-standards/project-layout) 

## Run protobuf Compiler 

```bash
protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/api/v1 foo-service.proto
```