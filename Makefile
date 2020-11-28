BINARY_NAME=tfcloud

.PHONY: install test lint crossbuild clean build deps

deps:
	go mod download
	go get -u golang.org/x/lint/golint

install: deps
	go install

lint: deps
	gofmt -s -l .
	golint ./...
	go vet ./...

test: lint deps
	go test -v ./...

build: deps
	go build -o $(GOPATH)/bin/$(BINARY_NAME)

crossbuild: test
	gox -os="linux darwin windows" -arch="386 amd64" -output "bin/$(BINARY_NAME)_{{.OS}}_{{.Arch}}/{{.Dir}}"

clean:
	go mod tidy
	go clean
	rm -f $(BINARY_NAME)
	rm -f bin/
