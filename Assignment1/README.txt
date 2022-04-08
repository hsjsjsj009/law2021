# Prerequisites
1. Golang
2. Install dependencies with command : go mod download

# How To Run
## Using docker-compose
Command : docker-compose -f docker-compose.yml up -d
URLs :
Rest -> http://localhost:8080
Web -> http://localhost:8081
Web with reverse proxy -> http://localhost:8082
## Using go command
### Rest Server
Command : cd rest && go run main.go
URL -> http://localhost:8080 (You can change the port by giving environment variable PORT)
### Web Server
Command : cd web && REST_SERVER=<rest-server-host>:<rest-server-port> go run main.go
example : cd web && REST_SERVER=localhost:8080 go run main.go
URL -> http://localhost:8081 (You can change the port by giving environment variable PORT)
