FROM golang:alpine as builder

WORKDIR /app

COPY . .
RUN go mod download && \
    go build -o app ./rest/.

FROM alpine:latest as runner

WORKDIR /opt/main/

RUN apk add --no-cache tzdata

COPY --from=builder /app/app .

ENV PORT=8080

ENTRYPOINT [ "/bin/sh","-c","/opt/main/app"]
