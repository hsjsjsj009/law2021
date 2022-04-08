FROM golang:alpine as builder

WORKDIR /app

COPY . .
RUN go mod download && \
    go build -o app ./server/.

FROM alpine:latest as runner

WORKDIR /opt/main/

COPY --from=builder /app/app .

ENV DB_HOST=postgres
ENV DB_PORT=5432

ENTRYPOINT [ "/bin/sh","-c","/opt/main/app"]
