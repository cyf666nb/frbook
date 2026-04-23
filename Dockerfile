FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
RUN GOTOOLCHAIN=auto GOPROXY=https://goproxy.cn,direct go mod download

COPY . .
RUN GOTOOLCHAIN=auto GOPROXY=https://goproxy.cn,direct go mod tidy

RUN GOTOOLCHAIN=auto GOPROXY=https://goproxy.cn,direct CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config config

EXPOSE 8080

CMD ["./main"]
