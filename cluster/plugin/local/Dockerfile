# =============================
FROM golang:1.9-alpine as build
# =============================
# system deps
RUN apk --no-cache add git binutils

ENV PKG github.com/appcelerator/amp/cluster/plugin/local
ENV DIR /go/src/${PKG}
ENV TARGET /tmp/amp-local
COPY . ${DIR}
WORKDIR ${DIR}

ARG LDFLAGS="-s"
# build
RUN go build -ldflags "${LDFLAGS}" -o ${TARGET} cmd/main.go

# shrink the binary
RUN strip ${TARGET}

# =============================
FROM alpine
# =============================
RUN apk add --no-cache ca-certificates
COPY --from=build /tmp/amp-local /usr/local/bin/
ENTRYPOINT [ "amp-local" ]
