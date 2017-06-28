echo "****************************"
export GOARCH=amd64
export CC="gcc"

mv libdockzen.a libdockzen.a_
gcc -c dockzen_test.c
ar cr libdockzen.a *.o
rm -rf *.o

make
rm -rf libdockze.a
mv libdockzen.a_ libdockzen.a
