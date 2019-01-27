VERSION = $(shell gobump show -r)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/Songmu/ghg.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

deps:
	env GO111MODULE=on go mod download

devel-deps: deps
	env GO111MODULE=on go get ${u} golang.org/x/lint/golint
	env GO111MODULE=on go get ${u} github.com/mattn/goveralls
	env GO111MODULE=on go get ${u} github.com/motemen/gobump/cmd/gobump
	env GO111MODULE=on go get ${u} github.com/Songmu/goxz/cmd/goxz
	env GO111MODULE=on go get ${u} github.com/Songmu/ghch/cmd/ghch
	env GO111MODULE=on go get ${u} github.com/tcnksm/ghr

test: deps
	env GO111MODULE=on go test

lint: devel-deps
	env GO111MODULE=on go vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build: deps
	env GO111MODULE=on go build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghg

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
