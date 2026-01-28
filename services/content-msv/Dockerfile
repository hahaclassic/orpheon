FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /content-msv

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./content-msv.exe ./cmd/content-msv/main.go

FROM alpine:latest AS runner
WORKDIR /
COPY --from=builder /content-msv/content-msv.exe /content-msv.exe
CMD ["/content-msv.exe"]
