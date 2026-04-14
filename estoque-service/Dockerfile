# ── Stage 1: build ───────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV GONOSUMDB=*
ENV GOFLAGS=-mod=mod

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CORREÇÃO: apontar para cmd/api (não cmd/server)
RUN go build -o /app/server ./cmd/api

# ── Stage 2: runtime mínimo ───────────────────────────────────────────────────
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8081
CMD ["./server"]