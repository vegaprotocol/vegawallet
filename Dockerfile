FROM golang:1.16.2 AS builder
ENV GOPROXY=direct GOSUMDB=off
WORKDIR /go/src/project
ADD *.go go.* ./
ADD cmd cmd
ADD crypto crypto
ADD fsutil fsutil
ADD logger logger
ADD service service
ADD version version
ADD wallet wallet
RUN env CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o go-wallet .

# # #

FROM ubuntu:20.04
ENTRYPOINT ["/go-wallet"]
RUN \
	apt update && \
	DEBIAN_FRONTEND=noninteractive apt install -y ca-certificates && \
	rm -rf /var/lib/apt/lists
COPY --from=builder /go/src/project/go-wallet /
