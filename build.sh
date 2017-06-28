#! /bin/sh
set -e

# NOTE
# build script for cli interface with cross-compile environment
#		for arm 	: ./build.sh arm
#		for amd64 	: ./build.sh
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

export BINARY_NAME="agent"	#fixed~~~~~

export CONTAINER_NAME="agent"
export CONTAINER_VERSION=":v1.0"

if [ "$1" = "arm" ]; then
        echo "****************************"
		echo "Target Binary arch is ARM"
		echo "****************************"
        export CGO_ENABLED=1
        export GOARCH=arm GOARM=7
        export CC="arm-linux-gnueabi-gcc"

        export CGO_LDFLAGS="-L${PWD}/src/lib/install/arm/lib"
else
		echo "****************************"
		echo "Target Binary arch is amd64"
		echo "****************************"
        export GOARCH=amd64
        export CC="gcc"

        export CGO_LDFLAGS="-L${PWD}/src/lib/install/amd64/lib"
fi

echo make clean
make clean

echo make build
make build

echo "\033[1;96mbuilt successfully... \033[0m"
file ${BINARY_NAME}

docker build -t ${CONTAINER_NAME}${CONTAINER_VERSION} .
echo "\033[1;96mcreated docker image successfully... \033[0m"
docker images ${CONTAINER_NAME}${CONTAINER_VERSION}

echo "run as \"docker run ${CONTAINER_NAME}${CONTAINER_VERSION}\""
