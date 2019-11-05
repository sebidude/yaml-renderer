FROM scratch

COPY build/linux/yaml-renderer /usr/bin/yaml-renderer
ENTRYPOINT ["/usr/bin/yaml-renderer"]
