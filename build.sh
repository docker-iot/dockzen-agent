#! /bin/sh
set -e

# NOTE
# build script for cli interface with cross-compile environment
#		for arm without proxy 		: ./build.sh arm 
#		for arm with proxy		: ./build.sh arm proxy
#		for amd64 without proxy 	: ./build.sh amd64
#		for amd64 with proxy		: ./build.sh amd64 proxy
#                
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

export BINARY_NAME="agent"	#fixed~~~~~

export CONTAINER_VERSION=":v1.1"

if [ $# -lt 1]; then
	echo "Usage: ./build.sh [arm|amd64]"
	exit
elif [ "$1" = "arm" ]; then
        echo "****************************"
		echo "Target Binary arch is ARM"
		echo "****************************"
        export CGO_ENABLED=1
        export GOARCH=arm GOARM=7
        export CC="arm-linux-gnueabi-gcc"
        export CONTAINER_NAME="dockzen-agent-arm"
        export CGO_LDFLAGS="-L${PWD}/src/lib/install/arm/lib"
elif [ "$1" = "amd64" ]; then
		echo "****************************"
		echo "Target Binary arch is amd64"
		echo "****************************"
        export GOARCH=amd64
        export CC="gcc"
        export CONTAINER_NAME="dockzen-agent"
        export CGO_LDFLAGS="-L${PWD}/src/lib/install/amd64/lib"
else
	exit
fi


echo make clean
make clean

echo make build
make build

echo "\033[1;96mbuilt successfully... \033[0m"
file ${BINARY_NAME}

if [ "$2" = "proxy" ]; then
	echo "******dockerfile_proxy******"
	cp ${PWD}/image/Dockerfile_proxy Dockerfile
else
	echo "******dockefile*************"
	pwd
	cp ${PWD}/image/Dockerfile Dockerfile	
fi

docker build -t ${CONTAINER_NAME}${CONTAINER_VERSION} .
echo "\033[1;96mcreated docker image successfully... \033[0m"
rm -f Dockerfile
docker images ${CONTAINER_NAME}${CONTAINER_VERSION}

echo "run as \"docker run ${CONTAINER_NAME}${CONTAINER_VERSION}\""
