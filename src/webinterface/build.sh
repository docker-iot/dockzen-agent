echo "****************************"
export GOARCH=amd64

mv ../lib/libdockzen.a ../lib/libdockzen.a_
gcc -c -o ../lib/dockzen_test.o ../../test/dockzen_test.c
ar cr ../lib/libdockzen.a ../lib/dockzen_test.o
rm -rf ../lib/dockzen_test.o

make

rm -rf ../lib/libdockzen.a
mv ../lib/libdockzen.a_ ../lib/libdockzen.a
