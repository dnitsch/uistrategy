NAME := uistrategy
OWNER := dnitsch
VERSION := "v0.0.0"
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
	    [ $$os = "windows" ] && EXT=".exe"; \
		GOOS=$$os CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/uiseeder-$$os ./cmd; \
	done

build: clean install bingen

build_ci: clean install_ci bingen
