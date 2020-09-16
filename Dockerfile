FROM golang:1.15-alpine as builder

WORKDIR /src/app
COPY . .

RUN apk add --no-cache \
        git \
#        ca-certificates \
        upx

#RUN go get -u github.com/gobuffalo/packr/v2/packr2 \
#    && go generate .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
        go build -ldflags="-w -s" -mod vendor -o /app .

RUN upx -q /app && \
    upx -t /app

# ---

FROM scratch

WORKDIR /

#RUN adduser -S -D -H -h /srv appuser
#USER appuser

#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app /app

EXPOSE 8080

CMD ["/app"]
