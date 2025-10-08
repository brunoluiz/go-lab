# syntax=docker/dockerfile:1

FROM scratch
ARG GOBIN_PATH
LABEL maintainer="Bruno Luiz da Silva <github@brunoluiz.net>"
LABEL org.opencontainers.image.source="https://github.com/brunoluiz/go-lab"

COPY ${GOBIN_PATH} /app
USER nobody
ENTRYPOINT ["/app"]
