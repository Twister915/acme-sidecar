FROM alpine:edge
RUN apk --no-cache add ca-certificates
RUN mkdir /app/
WORKDIR /app/
COPY acme-sidecar /app/
ENTRYPOINT /app/acme-sidecar