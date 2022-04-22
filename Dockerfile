FROM alpine:latest AS builder
RUN apk update && apk add --no-cache ca-certificates tzdata dumb-init && update-ca-certificates

FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /usr/bin/dumb-init /usr/bin/dumb-init
ADD ./waba-coreapp-mock /opt/waba-coreapp-mock
EXPOSE 3002
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/opt/waba-coreapp-mock", "--addr=0.0.0.0:3002", "-v"]
