FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /gateway-msv

COPY . .

RUN go mod tidy

RUN go build -o ./gateway-msv.exe ./cmd/gateway-msv/main.go

FROM alpine:latest AS runner
WORKDIR /
COPY --from=builder /gateway-msv/gateway-msv.exe /gateway-msv.exe
CMD ["/gateway-msv.exe"]
