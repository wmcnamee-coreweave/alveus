SHELL := bash

.PHONY: setup
setup:
	mise install

.PHONY: test
test: build
	go run github.com/onsi/ginkgo/v2/ginkgo run ./...

.PHONY: build
build: gen
	go build -o alveus ./cmd/alveus/main.go

.PHONY: gen
gen:
	go generate ./...
	go run ./cmd/gen-examples/main.go

.PHONY: git-push-tag
git-push-tag:
ifeq ($(TAG_TYPE),)
TAG_TYPE := unknown
endif
ifeq ($(TAG_SUFFIX),)
TAG_SUFFIX := unknown
endif
git-push-tag:
ifeq ($(CI),true)
	git config user.name "$${GITHUB_ACTOR}"
	git config user.email "$${GITHUB_ACTOR}@users.noreply.github.com"
endif
	# A previous release was created using a lightweight tag
	# git describe by default includes only annotated tags
	# git describe --tags includes lightweight tags as well
	DESCRIBE=$$(git tag -l --sort=-v:refname | grep -v nightly | head -n 1) && \
	MAJOR_VERSION=$$(echo "$${DESCRIBE}" | awk '{split($$0,a,"."); print a[1]}') && \
	MINOR_VERSION=$$(echo "$${DESCRIBE}" | awk '{split($$0,a,"."); print a[2]}') && \
	MINOR_VERSION="$$((${MINOR_VERSION:-0} + 1))" && \
	TAG="$${MAJOR_VERSION:-0}.$${MINOR_VERSION}.0-$(TAG_TYPE).$(TAG_SUFFIX)" && \
	echo "tag $${TAG}"
	git tag -a $TAG -m "$TAG: $(TAG_TYPE) build" && \
	git push origin $TAG

.PHONY: git-push-tag-pr
git-push-tag-pr: override TAG_TYPE = pr
git-push-tag-pr: override TAG_SUFFIX = $(shell git rev-parse --short HEAD)
git-push-tag-pr: git-push-tag

.PHONY: git-push-tag-nightly
git-push-tag-nightly: override TAG_TYPE = nightly
git-push-tag-nightly: override TAG_SUFFIX = $(shell date +'%Y%m%d')
git-push-tag-nightly: git-push-tag

.PHONY: release
release:
	mkdir -p $(CURDIR)/bin
	rm -rf $(CURDIR)/bin/docker
	ln -s $(shell which podman) $(CURDIR)/bin/docker && \
	export PATH="$(CURDIR)/bin:$(PATH)" && \
	goreleaser release --snapshot --clean

.PHONY: example
example: build
	cd examples && ../alveus generate \
		-s example-service.yaml \
		-r github.com/ghostsquad/fake
