APPLICATION ?= $(shell basename $(CURDIR))
BINARIES ?= $(shell find $(PWD)/cmd/ -maxdepth 1 -mindepth 1 -type d)
BUILD_DIR ?= $(CURDIR)/bin

VERSION_VAR := main.version
VERSION := $(shell git describe --long --tags --always --abbrev=8 --dirty)
GOBUILD_VERSION_ARGS := "-X $(VERSION_VAR)=$(VERSION)"

.PHONY: all
all: clean $(BUILD_DIR)

$(BUILD_DIR):
	@for BIN in ${BINARIES}; do \
	cd $${BIN}; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags $(GOBUILD_VERSION_ARGS) -o $(BUILD_DIR)/`basename $${BIN}` .; \
	done

.PHONY: clean
clean:
	-rm -rf $(BUILD_DIR)

.PHONY: install
install:
	go install -ldflags $(GOBUILD_VERSION_ARGS) ./...

.PHONY: lint
lint:
	gometalinter --config gometalinter.cfg ./...

.PHONY: test
test:
	go test ./... -covermode=atomic -v -race || exit 1

# Docker
DOCKER_REGISTRY ?= hub.docker.com
DOCKER_REGISTRY_REPO ?= s-kostyaev
DOCKER_TAG ?= 0.0.1
DOCKER_IMAGE := $(DOCKER_REGISTRY_REPO)/$(APPLICATION):$(DOCKER_TAG)

.PHONY: docker
docker: docker-build docker-push docker-clean

.PHONY: docker-build
docker-build:
	docker build --rm --tag "$(DOCKER_IMAGE)" .

.PHONY: docker-clean
docker-clean:
	docker rmi -f "$$(docker images -q $(DOCKER_IMAGE))"

.PHONY: docker-push
docker-push:
	docker tag "$(DOCKER_IMAGE)" "$(DOCKER_REGISTRY)/$(DOCKER_IMAGE)"
	docker push "$(DOCKER_REGISTRY)/$(DOCKER_IMAGE)"

# Misc
TOOLS := \
	github.com/alecthomas/gometalinter

.PHONY: install-tools
install-tools:
	$(foreach pkg,$(TOOLS),go get -u $(pkg);)
	gometalinter --install
