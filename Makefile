BIN = maven-runner

build: deps
	go build -o $(BIN)

lint: deps
	go vet
	golint

test: testdeps
	go test ./...

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .

install: deps
	go install

clean:
	go clean

.PHONY: build lint test deps testdeps install clean
