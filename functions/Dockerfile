# Stage 1: Build the Go binaries
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build each function's binary
# CGO_ENABLED=0 is important for static binaries, useful for scratch/alpine images
# -a -installsuffix cgo reduces image size
# -o specifies the output file name
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/a ./cmd/a
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/b ./cmd/b
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/c ./cmd/c
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/w2 ./cmd/w2
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/rw ./cmd/rw
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/d ./cmd/d

# Stage 2: Create a minimal image for each function
# Using alpine for a slightly larger but more debuggable image than scratch

# --- Function A Image ---
FROM alpine:latest AS function-a
EXPOSE 8080
WORKDIR /app

ENV W2_URL="http://function-w2.openfaas-fn.svc.cluster.local:8080"
ENV RW_URL="http://function-rw.openfaas-fn.svc.cluster.local:8080"

COPY --from=builder /app/bin/a .
ENTRYPOINT ["./a"]

# --- Function B Image ---
FROM alpine:latest AS function-b
EXPOSE 8080
WORKDIR /app

ENV W2_URL="http://function-w2.openfaas-fn.svc.cluster.local:8080"
ENV RW_URL="http://function-rw.openfaas-fn.svc.cluster.local:8080"

COPY --from=builder /app/bin/b .
ENTRYPOINT ["./b"]

# --- Function C Image ---
FROM alpine:latest AS function-c
EXPOSE 8080
WORKDIR /app

ENV D_URL="http://function-d.openfaas-fn.svc.cluster.local:8080"
ENV W2_URL="http://function-w2.openfaas-fn.svc.cluster.local:8080"

COPY --from=builder /app/bin/c .
ENTRYPOINT ["./c"]

# --- Function W2 Image ---
FROM alpine:latest AS function-w2
EXPOSE 8080
WORKDIR /app
COPY --from=builder /app/bin/w2 .
ENTRYPOINT ["./w2"]

# --- Function RW Image ---
FROM alpine:latest AS function-rw
EXPOSE 8080
WORKDIR /app
COPY --from=builder /app/bin/rw .
ENTRYPOINT ["./rw"]

# --- Function D Image ---
FROM alpine:latest AS function-d
EXPOSE 8080
WORKDIR /app

ENV RW_URL="http://function-rw.openfaas-fn.svc.cluster.local:8080"

COPY --from=builder /app/bin/d .
ENTRYPOINT ["./d"]