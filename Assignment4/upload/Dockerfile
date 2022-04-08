FROM golang:alpine as builder

WORKDIR /app

COPY . .
RUN go mod download && \
    go build -o app ./upload/.

FROM alpine:latest as runner

WORKDIR /opt/main/

COPY --from=builder /app/app .
COPY --from=builder /app/html html

ENTRYPOINT [ "/bin/sh","-c","/opt/main/app"]
