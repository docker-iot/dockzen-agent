export GOPATH := $(shell pwd)/../../:$(shell pwd)/../../vendor

build:
	go test WS_Services_test.go WS_Services.go device_info.go WebInterface_WS.go ./services ./lib ./include
