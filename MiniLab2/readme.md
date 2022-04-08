# How To Run
## Prerequisites
- Prepare rabbitmq server in local with port 5672 for amqp and port 15674 for STOMP Web Socket
## Server 1
- Go to the server1 directory and run main.go
```bash
cd server1
go run main.go
```
## Server 2
- Go to the server2 directory and run main.go
```bash
cd server1
go run main.go
```

# How To Compile The .proto File
```shell
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    compression/compression.proto
```