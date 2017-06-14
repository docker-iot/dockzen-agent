export GOPATH := $(shell pwd)

build:
	go build -a -ldflags '-extldflags "--static"' agent

clean:
	rm -rf agent
