.PHONY: all
all: binary

ORG=karlkfi
REPO=$(shell git rev-parse --show-toplevel | xargs basename)
REPO_PATH=github.com/$(ORG)/$(REPO)

.PHONY: clean
clean:
	rm -rf _output

.PHONY: godep
godep:
	go get github.com/tools/godep

.PHONY: binary
binary:
	godep go build -v -o _output/bin/oinker $(REPO_PATH)

.PHONY: binary-alpine
binary-alpine:
	docker run --rm -v "$(CURDIR):/go/src/$(REPO_PATH)" -w /go/src/$(REPO_PATH) $(ORG)/$(REPO)-build:latest make godep binary

.PHONY: image
image: binary-alpine
	docker build --no-cache -t $(ORG)/$(REPO):latest .

.PHONY: build-image
build-image:
	cd build && docker build -t $(ORG)/$(REPO)-build:latest .
