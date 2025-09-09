# Stage 1: Build the Go binaries
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/program ./pkg

FROM alpine:latest AS function
EXPOSE 8080
WORKDIR /app

ENV PRIME_NUMBERS_URL_PARALLEL="http://dispatcher.default.svc.cluster.local/function/openfaas-fn/prime-numbers-invoked-parallel/"
ENV PRIME_NUMBERS_URL_SEQUENTIAL="http://dispatcher.default.svc.cluster.local/function/openfaas-fn/prime-numbers-invoked-sequential/"

COPY --from=builder /app/bin/program .
ENTRYPOINT ["./program"]