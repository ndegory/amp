# =============================
FROM golang:1.9-alpine as build
# =============================
RUN echo "@edgecommunity http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk update
RUN apk add build-base
RUN apk add ca-certificates
RUN apk add git
RUN apk add upx@edgecommunity

# go
# ensure any binaries are placed in /tmp/bin
ENV TMPBIN /tmp/bin
ENV TARGET ${TMPBIN}/ampctl
ENV PKG github.com/appcelerator/amp/cluster/agent
ENV PKGCMD ${PKG}/cmd
ENV PKGROOT /go/src/${PKG}
ENV PKGCMD ${PKG}/cmd
ARG LDFLAGS="-s -w"

#RUN go get github.com/LK4D4/vndr

COPY . ${PKGROOT}
WORKDIR ${PKGROOT}

# get vendor packages
# TODO: uncomment after development (convenient to just existing vendor for now)
#RUN vndr

# test
RUN go test -v -timeout 30m ${PKG}/admin

# build
RUN echo "Building ampagent with ${LDFLAGS}" && \
    CGO_ENABLED=0 go build -a -ldflags "${LDFLAGS}" -o ${TARGET} ${PKGCMD}

# shrink the binary
RUN upx ${TARGET}

# =============================
FROM scratch
# =============================
ENV SWARM_SOCKET "/var/run/docker/swarm/control.sock"
ENV PKG github.com/appcelerator/amp/cluster/agent
ENV PKGROOT /go/src/${PKG}

COPY --from=build /etc/ssl/certs /etc/ssl/certs
COPY --from=build /tmp/bin/ /bin
COPY --from=build ${PKGROOT}/defaults /defaults
COPY ./stacks /stacks

ENTRYPOINT [ "/bin/ampctl" ]
