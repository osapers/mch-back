FROM golang:1.16.4-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add -U --no-cache ca-certificates make
WORKDIR /build
COPY . .
RUN make build

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/back /
EXPOSE 8080
ENTRYPOINT ["/back"]