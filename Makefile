
export GOPATH := $(shell pwd):$(shell pwd)/vendor

build:
	go build -o ${BINARY_NAME} -a -v -ldflags '-extldflags "--static"' main

clean:
	rm -rf ${BINARY_NAME}
