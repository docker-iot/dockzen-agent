
export GOPATH := $(shell pwd):$(shell pwd)/vendor

build:
	go build -a -ldflags '-extldflags "--static"' ${BINARY_NAME}

clean:
	rm -rf ${BINARY_NAME}
