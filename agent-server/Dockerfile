FROM golang:1.18.3 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN GOOS=linux CGO_ENABLED=0 go build -o agent-server agent-server/server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /etc/go-sro
COPY --from=builder /app/agent-server/server agent-server

EXPOSE 15882
CMD ["./agent-server"]