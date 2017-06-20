
export GOPATH := $(shell pwd):$(shell pwd)/vendor

build:
	go build -a -ldflags '-extldflags "--static"' agent

clean:
	rm -rf agent
