FROM golang:alpine as builder

WORKDIR /app

COPY . .
RUN go mod download && \
    go build -o app ./web/.

FROM alpine:latest as runner

WORKDIR /opt/main/

RUN apk add --no-cache tzdata

COPY --from=builder /app/app .
COPY --from=builder /app/web/templates ./templates

ENV PORT=8081

ENTRYPOINT [ "/bin/sh","-c","/opt/main/app"]
