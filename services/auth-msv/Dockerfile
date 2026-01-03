FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /auth-msv

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./auth-msv.exe ./cmd/auth-msv/main.go

FROM alpine:latest AS runner
WORKDIR /
COPY --from=builder /auth-msv/auth-msv.exe /auth-msv.exe
CMD ["/auth-msv.exe"]
