# Creator
Dipta Laksmana Baswara<br>
1806235832

# How To Run
## Prerequisites
- Prepare rabbitmq server in local with port 5672 for amqp and port 15674 for STOMP Web Socket
## Server 1
- Run server1/main.go
```bash
go run server1/main.go
```
## Server 2
- Run server2/main.go
```bash
go run server2/main.go
```

# How To Compile The .proto File
```shell
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    downloader/downloader.proto
```