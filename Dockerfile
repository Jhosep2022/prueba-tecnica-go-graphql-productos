FROM golang:1.26.5-alpine3.24 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/product-api \
    ./cmd

FROM alpine:3.24

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/product-api /app/product-api

ENV PORT=8080

USER app

EXPOSE 8080

ENTRYPOINT ["/app/product-api"]
