FROM golang:1.13-alpine as builder

RUN apk add --no-cache --update git make
RUN mkdir /build
WORKDIR /build
RUN git clone https://github.com/sebidude/yaml-renderer.git
WORKDIR /build/yaml-renderer
RUN make test build-linux

FROM scratch
COPY --from=builder /build/yaml-renderer/build/linux/yaml-renderer /usr/bin/yaml-renderer
ENTRYPOINT ["/usr/bin/yaml-renderer"]
