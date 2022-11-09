NAME := uistrategy
OWNER := dnitsch
GIT_TAG := "0.1.3"
VERSION := "v$(GIT_TAG)"
REVISION := $(shell git rev-parse --short HEAD)

LDFLAGS := -ldflags="-s -w -X \"github.com/$(OWNER)/$(NAME)/cmd/uistrategy.Version=$(VERSION)\" -X \"github.com/$(OWNER)/$(NAME)/cmd/uistrategy.Revision=$(REVISION)\" -extldflags -static"

install:
	go mod tidy
	go mod vendor
	
install_ci:
	go mod vendor

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/*
	mkdir -p dist

bingen:
	for os in darwin linux windows; do \
		GOOS=$$os CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/uiseeder-$$os ./cmd; \
	done

build: clean install bingen

build_ci: clean install_ci bingen

tag: 
	git tag "v$(GIT_TAG)"
	git push origin "v$(GIT_TAG)"

release:
	OWNER=$(OWNER) NAME=$(NAME) PAT=$(PAT) VERSION=$(VERSION) . hack/release.sh 

# build tag release
btr: build tag release
	echo "ran build tag release"
