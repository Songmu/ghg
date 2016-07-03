test: deps
	go test

deps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get github.com/mattn/goveralls

lint: deps
	go vet
	golint -set_exit_status

cover: deps
	goveralls

release:
	_tools/releng

.PHONY: test deps lint cover release
