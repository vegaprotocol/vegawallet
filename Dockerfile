FROM golang:1.17-alpine AS builder
RUN apk add --no-cache git
ENV GOPROXY=direct GOSUMDB=off
WORKDIR /go/src/project
ADD . .
RUN env CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o vegawallet .

FROM alpine:3.14
ENTRYPOINT ["vegawallet"]
RUN apk add --no-cache bash ca-certificates jq
COPY --from=builder /go/src/project/vegawallet /usr/local/bin/
