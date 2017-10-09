# =============================
FROM golang:1.9-alpine as build
# =============================
# system deps
RUN echo "@edgecommunity http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk --no-cache add git upx@edgecommunity

# go deps
RUN go get github.com/LK4D4/vndr
ENV PKG github.com/appcelerator/amp/cluster/plugin/aws
ENV DIR /go/src/${PKG}
ENV TARGET /tmp/aws
COPY . ${DIR}
WORKDIR ${DIR}
RUN vndr

ARG LDFLAGS="-s"
# build
RUN go build -ldflags "$LDFLAGS" -o ${TARGET} cmd/main.go

# shrink the binary
RUN upx ${TARGET}

# =============================
FROM alpine
# =============================
RUN apk add --no-cache ca-certificates
COPY --from=build /tmp/aws /usr/local/bin/
ENTRYPOINT [ "aws" ]

