echo "****************************"
export GOARCH=amd64
export CC="gcc"

mv libdockzen.a libdockzen.a_
gcc -c ../../test/dockzen_test.c
ar cr libdockzen.a dockzen_test.o
rm -rf dockzen_test.o

make
rm -rf libdockze.a
mv libdockzen.a_ libdockzen.a
