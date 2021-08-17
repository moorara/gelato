# BUILD STAGE
FROM golang:1.17-alpine as builder
RUN apk add --no-cache git
WORKDIR /repo
COPY . .
ARG ldflags
RUN go build -ldflags "$ldflags"

# FINAL STAGE
FROM alpine:3.14
EXPOSE 4000 5000
RUN apk add --no-cache curl ca-certificates
HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://localhost:4000/health || exit 1
COPY --from=builder /repo/grpc-service /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/grpc-service
USER nobody
ENTRYPOINT [ "grpc-service" ]
