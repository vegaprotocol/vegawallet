ifeq ($(CI),)
	# Not in CI
	VERSION := dev-$(USER)
	VERSION_HASH := $(shell git rev-parse HEAD | cut -b1-8)
else
	# In CI
	ifneq ($(RELEASE_VERSION),)
		VERSION := $(RELEASE_VERSION)
	else
		# No tag, so make one
		VERSION := $(shell git describe --tags 2>/dev/null)
	endif
	VERSION_HASH := $(shell echo "$(GITHUB_SHA)" | cut -b1-8)
endif

install:
	go install -v -ldflags "-X code.vegaprotocol.io/vegawallet/version.Version=${VERSION} -X code.vegaprotocol.io/vegawallet/version.VersionHash=${VERSION_HASH}"

release-windows-amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o build/vegawallet.exe -ldflags "-X code.vegaprotocol.io/vegawallet/version.Version=${VERSION} -X code.vegaprotocol.io/vegawallet/version.VersionHash=${VERSION_HASH}"
	cd build && 7z a -tzip vegawallet-windows-amd64.zip vegawallet.exe

release-windows-arm64:
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -o build/vegawallet.exe -ldflags "-X code.vegaprotocol.io/vegawallet/version.Version=${VERSION} -X code.vegaprotocol.io/vegawallet/version.VersionHash=${VERSION_HASH}"
	cd build && 7z a -tzip vegawallet-windows-arm64.zip vegawallet.exe

release-macos-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/vegawallet -ldflags "-X code.vegaprotocol.io/vegawallet/version.Version=${VERSION} -X code.vegaprotocol.io/vegawallet/version.VersionHash=${VERSION_HASH}"
	cd build && zip vegawallet-darwin-amd64.zip vegawallet

release-macos-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o build/vegawallet -ldflags "-X code.vegaprotocol.io/vegawallet/version.Version=${VERSION} -X code.vegaprotocol.io/vegawallet/version.VersionHash=${VERSION_HASH}"
	cd build && zip vegawallet-darwin-arm64.zip vegawallet

release-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/vegawallet -ldflags "-X code.vegaprotocol.io/vegawallet/version.Version=${VERSION} -X code.vegaprotocol.io/vegawallet/version.VersionHash=${VERSION_HASH}"
	cd build && zip vegawallet-linux-amd64.zip vegawallet

release-linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o build/vegawallet -ldflags "-X code.vegaprotocol.io/vegawallet/version.Version=${VERSION} -X code.vegaprotocol.io/vegawallet/version.VersionHash=${VERSION_HASH}"
	cd build && zip vegawallet-linux-arm64.zip vegawallet
