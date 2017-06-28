echo "****************************"
export GOARCH=amd64
export CC="gcc"

gcc -c dockzen_test.c
ar cr libdockzen.a *.o
rm -rf *.o

make
