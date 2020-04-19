install:
	go install ./...

proto:
	protoc --go_out=paths=source_relative,plugins=grpc:. ./proto/*.proto
