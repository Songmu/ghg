VERSION = $(shell gobump show -r)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/Songmu/ghg.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

GO ?= GO111MODULE=on go

deps:
	env GO111MODULE=on go mod download

devel-deps: deps
	$(GO) get ${u} \
	  golang.org/x/lint/golint             \
	  github.com/mattn/goveralls           \
	  github.com/motemen/gobump/cmd/gobump \
	  github.com/Songmu/goxz/cmd/goxz      \
	  github.com/Songmu/ghch/cmd/ghch      \
	  github.com/tcnksm/ghr

test: deps
	$(GO) test

lint: devel-deps
	$(GO) vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build: deps
	$(GO) build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghg

crossbuild: devel-deps
	goxz -pv=v$(shell gobump show -r) -build-ldflags=$(BUILD_LDFLAGS) \
	  -os=linux,darwin,windows,freebsd -arch=amd64 -d=./dist/v$(shell gobump show -r) \
	  ./cmd/ghg

bump: devel-deps
	_tools/releng

upload:
	ghr v$(VERSION) dist/v$(VERSION)

release: bump crossbuild upload

.PHONY: deps devel-deps test lint cover build crossbuild build upload release
