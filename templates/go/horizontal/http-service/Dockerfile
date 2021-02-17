# BUILD STAGE
FROM golang:1.16-alpine as builder
RUN apk add --no-cache git
WORKDIR /repo
COPY . .
ARG ldflags
RUN go build -ldflags "$ldflags"

# FINAL STAGE
FROM alpine:3.13
EXPOSE 4000
RUN apk add --no-cache curl ca-certificates
HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://localhost:4000/health || exit 1
COPY --from=builder /repo/http-service /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/http-service
USER nobody
ENTRYPOINT [ "http-service" ]
