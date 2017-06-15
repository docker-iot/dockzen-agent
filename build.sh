echo "****************************"

if [ "$1" = "arm" ]; then
        echo "Target Binary arch is ARM"
        export CGO_ENABLED=1
        export GOARCH=arm GOARM=7
        export CC="arm-linux-gnueabi-gcc"
        arm-linux-gnueabi-gcc -c src/lib/*.c
else
        echo "Target Binary arch is amd64"
        export GOARCH=amd64
        export CC="gcc"
        gcc -c src/lib/*.c
fi

ar cr src/lib/libdockzen.a *.o

rm -rf *.o

echo make clean
make clean

echo make build
make build
