ARG GOLANG_VERSION=1.17
FROM golang:${GOLANG_VERSION} AS builder

LABEL VERSION="0.1.1"
LABEL org.opencontainers.image.description "Simple CD tool for docker swarm"
RUN apt-get -qq update && apt-get -yqq install upx

ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux

WORKDIR /src

COPY . .
RUN go build \
  -a \
  -trimpath \
  -ldflags "-s -w -extldflags '-static'" \
  -tags 'osusergo netgo static_build' \
  -o /bin/swarm_deploy \
  ./main.go


RUN upx -q -9 /bin/swarm_deploy

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/swarm_deploy /bin/swarm_deploy

ENTRYPOINT ["/bin/swarm_deploy"]