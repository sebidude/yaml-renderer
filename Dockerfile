FROM alpine as builder

RUN apk add shadow ca-certificates
RUN useradd -u 10001 yaml-renderer

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY build/linux/yaml-renderer /yaml-renderer

USER yaml-renderer
ENTRYPOINT ["/yaml-renderer"]
