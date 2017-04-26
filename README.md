## dockzen-agent

dockzen-agent is designed to manage for Tizen with Docker and provide services to a client.

### Developer Quick-Start

To build the daemons, the following build system dependencies are required:

* go 1.7.5 or above
* gb tool due to library dependecy

#### go 1.7.5

```
$ wget https://storage.googleapis.com/golang/go1.7.5.linux-amd64.tar.gz

```
If you extract the file and see the 'go' folder.
Copy 'go' folder into '/usr/local/go'
Set up the GOROOT, GOPATH, PATH

```
$ export PATH=$PATH:/usr/local/go/bin/
$ export GOPATH=$(go env GOPATH)
$ export PATH=$PATH:$(go env GOPATH)/bin
```

#### gb

```
$ go get github.com/constabulary/gb/...
```

If you install gb correctly, you can check whether it is installed or not using below command

```
$ gb info
```

#### build

build for arm
```
$ ./build.sh arm
```
or

build for amd64
```
$ ./build.sh
```

#### build clean

```
$ make clean
```
/bin/ folder is created and **two binaries** you can see in the folder.
**dockzen-agentconn** is client which can receive a command from web server
**dockzen-agent** is main daemon to check request form dockzen-agentconn and request it to dockzen-launcher

