NAME := uistrategy
OWNER := dnitsch
GIT_TAG := "0.3.0"
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
		GOOS=$$os CGO_ENABLED=0 go build -mod=readonly -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/uiseeder-$$os ./cmd; \
	done

build: clean install bingen

build_ci: clean install_ci bingen

tag: 
	git tag -a $(VERSION) -m "ci tag release uistrategy" $(REVISION)
	git push origin $(VERSION)

release:
	OWNER=$(OWNER) NAME=$(NAME) PAT=$(PAT) VERSION=$(VERSION) . hack/release.sh 

# build tag release
btr: build tag release
	echo "ran build tag release"

# TEST
test: test_prereq
	go test `go list ./... | grep -v */generated/` -v -mod=readonly -race -coverprofile=.coverage/out | go-junit-report > .coverage/report-junit.xml && \
	gocov convert .coverage/out | gocov-xml > .coverage/report-cobertura.xml

test_ci:
	go test ./... -mod=readonly

test_prereq: 
	mkdir -p .coverage
	go install github.com/jstemmer/go-junit-report/v2@latest && \
	go install github.com/axw/gocov/gocov@latest && \
	go install github.com/AlekSi/gocov-xml@latest

show_coverage: test
	go tool cover -html=.coverage/out