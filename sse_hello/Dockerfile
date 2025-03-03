# Build stage
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

WORKDIR /build

COPY . .

ARG TARGETPLATFORM
ARG BUILDPLATFORM

RUN echo "Building on $BUILDPLATFORM, targeting $TARGETPLATFORM"

# Set CGO_ENABLED=0 for static linking and GOOS/GOARCH based on target platform
RUN case "$TARGETPLATFORM" in \
    "linux/amd64") GOARCH=amd64 ;; \
    "linux/arm64") GOARCH=arm64 ;; \
    *) GOARCH=arm64 ;; \
    esac && \
    CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o app main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary and static files from builder
COPY --from=builder /build/app .
COPY --from=builder /build/public ./public

# 기본 prefix 설정
ENV URL_PATH_PREFIX=/app/sse-hello

EXPOSE 8080

CMD ["./app"] 