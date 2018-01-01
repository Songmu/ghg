CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/Songmu/ghg.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

deps:
	go get ${u} github.com/golang/dep/cmd/dep
	dep ensure

devel-deps: deps
	go get ${u} github.com/golang/lint/golint
	go get ${u} github.com/mattn/goveralls
	go get ${u} github.com/motemen/gobump/cmd/gobump
	go get ${u} github.com/Songmu/goxz/cmd/goxz
	go get ${u} github.com/Songmu/ghch/cmd/ghch
	go get ${u} github.com/tcnksm/ghr

test: deps
	go test

lint: devel-deps
	go vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghg

crossbuild: devel-deps
	goxz -pv=v$(shell gobump show -r) -build-ldflags=$(BUILD_LDFLAGS) \
	  -os=linux,darwin,windows,freebsd -arch=amd64 -d=./dist/v$(shell gobump show -r) \
	  ./cmd/ghg

release: devel-deps
	_tools/releng
	_tools/upload_artifacts

.PHONY: deps devel-deps test lint cover build crossbuild release
