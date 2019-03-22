.PHONY: default
default: build

.PHONY: build
build:
	goxc -build-ldflags="-w -s -X github.com/ntrv/check-aws-ec2-mainte/lib.version=$(git describe --tags)"

.PHONY: test
test:
	go test -v -race -cover -covermode=atomic ./lib -coverprofile=cover.profile
