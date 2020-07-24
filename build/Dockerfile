FROM golang:1.14.6-alpine AS builder

ARG PROJECT_NAME="registry-auth-server"
ARG VERSION="dev"
ARG COMMIT="none"
ENV CGO_ENABLED=0
ENV GOPROXY=https://proxy.golang.org

COPY . /src/
WORKDIR /src
RUN go mod vendor
RUN DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') go build -ldflags="-s -w -X github.com/angelbarrera92/registry-auth-server/cmd.version=${VERSION} -X github.com/angelbarrera92/registry-auth-server/cmd.commit=${COMMIT} -X github.com/angelbarrera92/registry-auth-server/cmd.date=${DATE}" -mod vendor -o /registry-auth-server

FROM alpine:3.12

LABEL org.opencontainers.image.name=${PROJECT_NAME}
LABEL org.opencontainers.image.revision=${COMMIT}
LABEL org.opencontainers.image.version=${VERSION}

COPY --from=builder /registry-auth-server /

RUN mkdir /etc/registry-auth-server && \
    chown 1001:1001 /etc/registry-auth-server

USER 1001
ENTRYPOINT ["/registry-auth-server"]
