echo "****************************"

if [ "$1" = "arm" ]; then
        echo "Target Binary arch is ARM"
        export CGO_ENABLED=1
        export GOARCH=arm GOARM=7
        export CC="arm-linux-gnueabi-gcc"
        arm-linux-gnueabi-gcc -c -o src/lib/dockzen_api.o src/lib/dockzen_api.c
else
        echo "Target Binary arch is amd64"
        export GOARCH=amd64
        export CC="gcc"
        gcc -c -o src/lib/dockzen_api.o src/lib/dockzen_api.c
fi

ar cr src/lib/libdockzen.a src/lib/dockzen_api.o

rm -rf src/lib/dockzen_api.o

echo make clean
make clean

echo make build
make build
