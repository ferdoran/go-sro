FROM golang:1.18.3 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN GOOS=linux CGO_ENABLED=0 go build -o gateway-server gateway-server/server.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /etc/go-sro/
COPY --from=builder /app/gateway-server/server gateway-server

EXPOSE 15779
CMD ["./gateway-server"]