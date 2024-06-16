# Build stage
FROM golang:1.22.4-bookworm AS builder

WORKDIR /build/

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o rest-server github.com/colmmurphy91/muzz/cmd/server

# Final stage
FROM debian:12.5-slim

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates && \
      rm -rf /var/lib/apt/lists/*

WORKDIR /api/
ENV PATH=/api/bin/:$PATH

COPY --from=builder /build/rest-server ./bin/rest-server
COPY --from=builder /build/env.example .


EXPOSE 8080

CMD ["./bin/rest-server", "-env", "/api/env.example"]