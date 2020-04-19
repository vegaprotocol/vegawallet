ifeq ($(CI),)
	# Not in CI
	VERSION := dev-$(USER)
	VERSION_HASH := $(shell git rev-parse HEAD | cut -b1-8)
else
	# In CI
	ifneq ($(DRONE),)
		# In Drone: https://docker-runner.docs.drone.io/configuration/environment/variables/

		ifneq ($(DRONE_TAG),)
			VERSION := $(DRONE_TAG)
		else
			# No tag, so make one
			VERSION := $(shell git describe --tags 2>/dev/null)
		endif
		VERSION_HASH := $(shell echo "$(CI_COMMIT_SHA)" | cut -b1-8)

	else
		# In an unknown CI
		VERSION := unknown-CI
		VERSION_HASH := unknown-CI
	endif
endif

install:
	go install -v -ldflags "-X main.Version=${VERSION} -X main.VersionHash=${VERSION_HASH}" ./...

proto:
	protoc --go_out=paths=source_relative,plugins=grpc:. ./proto/*.proto
