FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o server benchmark.go

FROM alpine
WORKDIR /build
COPY --from=builder /build/benchmark /build/benchmark
ENTRYPOINT ["./benchmark"]