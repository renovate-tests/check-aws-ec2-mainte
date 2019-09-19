RELEASE_URL := https://github.com/ntrv/check-aws-ec2-mainte/releases/download/$(shell git describe --tags | cut -d- -f1)/check-aws-ec2-mainte_darwin_amd64.zip
RELEASE_FILE := $(notdir $(RELEASE_URL))

.PHONY: default
default: build


.PHONY: bootstrap
bootstrap:
	go get -v github.com/laher/goxc

.PHONY: build
build:
	goxc -build-ldflags="-w -s -X github.com/ntrv/check-aws-ec2-mainte/lib.version=$(shell git describe --tags)"

.PHONY: test
test:
	go test -v -race -cover -covermode=atomic ./... -coverprofile=cover.profile

.PHONY: checksum
checksum:
	curl $(RELEASE_URL) -o $(TMPDIR)$(RELEASE_FILE)
	sha256sum $(TMPDIR)$(RELEASE_FILE)
