echo "****************************"
export GOARCH=amd64
export CC="gcc"
export CGO_ENABLED=1
export CGO_LDFLAGS="-L${PWD}/install/amd64/lib"
make
