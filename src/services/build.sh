echo "****************************"
export GOARCH=amd64
export CGO_ENABLED=1
export CGO_LDFLAGS="-L${PWD}/../lib/install/amd64/lib"
make
