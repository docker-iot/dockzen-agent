echo "****************************"
export GOARCH=amd64
export CC="gcc"

gcc -c ../lib/dockzen_test.c
ar cr ../lib/libdockzen.a ../lib/*.o
rm -rf *.o

make
