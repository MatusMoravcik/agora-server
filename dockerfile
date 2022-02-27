# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS builder
RUN mkdir /build
ADD go.mod go.sum main.go /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/agora-server /app/
COPY views/ /app/views
WORKDIR /app
CMD ["./agora-server"]