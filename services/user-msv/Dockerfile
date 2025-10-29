FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /build

COPY . .

WORKDIR /build/services/user-msv

RUN go mod tidy

RUN go build -o /user-msv.exe ./cmd/user-msv/main.go

FROM alpine:latest AS runner
WORKDIR /
COPY --from=builder /user-msv.exe /user-msv.exe
CMD ["/user-msv.exe"]
