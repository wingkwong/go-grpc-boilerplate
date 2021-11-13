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

## MySQL

```
brew install mysql
```

```
mysql.server start
```

### Optional

Install ``vscode-proto3`` in Visual Studio Code

## Project Layout

Refer to [golang-standards/project-layout](https://github.com/golang-standards/project-layout) 

## Run protobuf Compiler 

```bash
protoc --proto_path=api/proto/v1 --proto_path=third_party --go_out=plugins=grpc:pkg/api/v1 foo-service.proto
protoc --proto_path=api/proto/v1 --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api/v1 foo-service.proto
protoc --proto_path=api/proto/v1 --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/v1 foo-service.proto
```

## Run gRPC Server with HTTP/REST gateway

```
cd cmd/server
go build .
./server -grpc-port=9090 -http-port=8080 -db-host=<HOST>:3306 -db-user=<DB_USER> -db-password=<DB_PASSWORD> -db-schema=<DB_SCHEMA> -log-level=-1
```

### Logging Level

- -1 : DebugLevel logs are typically voluminous, and are usually disabled in production.
- 0: InfoLevel is the default logging priority.
- 1: WarnLevel logs are more important than Info, but don't need individual human review.
- 2: ErrorLevel logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs.
- 3: DPanicLevel logs are particularly important errors. In development the logger panics after writing the message.
- 4: PanicLevel logs a message, then panics.
- 5: FatalLevel logs a message, then calls os.Exit(1).

## Run gRPC Client

```
cd cmd/client-grpc
go build .
./client-grpc -server=localhost:9090
```

## Run REST Client

```
cd cmd/client-rest
go build .
./client-rest -server=http://localhost:8080
```

## Run test

```
cd pkg/service/v1
go test
```