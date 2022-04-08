FROM golang:alpine as builder

WORKDIR /app

COPY . .
RUN go mod download && \
    go build -o app ./download/.

FROM alpine:latest as runner

WORKDIR /opt/main/

COPY --from=builder /app/app .

ENTRYPOINT [ "/bin/sh","-c","/opt/main/app"]
