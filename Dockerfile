FROM golang:1.18-alpine as builder

WORKDIR /src/app
COPY . .

RUN apk add --no-cache \
        git \
        upx

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
        go build -ldflags="-w -s" -mod vendor -o /app ./cmd/app/...

RUN upx -q /app && \
    upx -t /app

FROM alpine:3.13

WORKDIR /

RUN apk add --no-cache bash curl ca-certificates

COPY --from=builder /app /app

CMD ["/app"]
