ifeq ($(CI),)
	# Not in CI
	VERSION := dev-$(USER)
	VERSION_HASH := $(shell git rev-parse HEAD | cut -b1-8)
else
	# In CI
	ifneq ($(GITHUB_REF),)
		VERSION := $(GITHUB_REF)
	else
		# No tag, so make one
		VERSION := $(shell git describe --tags 2>/dev/null)
	endif
	VERSION_HASH := $(shell echo "$(GITHUB_SHA)" | cut -b1-8)
endif

install:
	go install -v -ldflags "-X main.Version=${VERSION} -X main.VersionHash=${VERSION_HASH}" ./...

proto:
	protoc --go_out=paths=source_relative,plugins=grpc:. ./proto/*.proto

release:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/vegawallet-linux-amd64 -ldflags "-X main.Version=${VERSION} -X main.VersionHash=${VERSION_HASH}" ./cmd/vegawallet
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o build/vegawallet-linux-386 -ldflags "-X main.Version=${VERSION} -X main.VersionHash=${VERSION_HASH}" ./cmd/vegawallet
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/vegawallet-windows-amd64 -ldflags "-X main.Version=${VERSION} -X main.VersionHash=${VERSION_HASH}" ./cmd/vegawallet
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o build/vegawallet-windows-386 -ldflags "-X main.Version=${VERSION} -X main.VersionHash=${VERSION_HASH}" ./cmd/vegawallet
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/vegawallet-darwin-amd64 -ldflags "-X main.Version=${VERSION} -X main.VersionHash=${VERSION_HASH}" ./cmd/vegawallet
	GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -o build/vegawallet-darwin-386 -ldflags "-X main.Version=${VERSION} -X main.VersionHash=${VERSION_HASH}" ./cmd/vegawallet
