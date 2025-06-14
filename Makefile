VERSION = $(shell godzil show-version)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/Songmu/ghg.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

.PHONY: deps
deps:
	go get ${u}
	go mod tidy

.PHONY: devel-deps
devel-deps:
	go install github.com/tcnksm/ghr@latest
	go install github.com/Songmu/godzil/cmd/godzil@latest

.PHONY: test
test:
	go test

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghg

.PHONY: install
install:
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/ghg

CREDITS: go.sum devel-deps
	godzil credits -w

DIST_DIR = dist
.PHONY: crossbuild
crossbuild: devel-deps
	rm -rf $(DIST_DIR)
	godzil crossbuild -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
      -os=linux,darwin,windows,freebsd -d=$(DIST_DIR) ./cmd/*
	cd $(DIST_DIR) && shasum -a 256 $$(find * -type f -maxdepth 0) > SHA256SUMS

.PHONY: upload
upload:
	ghr v$(VERSION) $(DIST_DIR)
