## dockzen-agent

dockzen-agent is designed to manage for Tizen with Docker and provide services to a client.

### Developer Quick-Start

To build the daemons, the following build system dependencies are required:

* go 1.7.5 or above

#### go 1.7.5

```
$ wget https://storage.googleapis.com/golang/go1.7.5.linux-amd64.tar.gz

```
If you extract the file and see the 'go' folder.<br>
Copy 'go' folder into '/usr/local/go'<br>
Set up the GOROOT, GOPATH, PATH<br>

```
$ export PATH=$PATH:/usr/local/go/bin/
$ export GOPATH=$(go env GOPATH)
$ export PATH=$PATH:$(go env GOPATH)/bin
```
#### Server URL

Put your server url in ./data/server_url.json of agent and build.<br>
The url will be created and refered in HostOS.<br>

```
/etc/dockzen/container/agent/config/server_url.json
```

### Device UUID

You can set Device unique ID which will be created and refered in HostOS.<br>

```
/etc/dockzen/container/agent/config/device_uuid.json
```

#### build

build for arm<br>
If it is done, Cotainer would be created<br>
```
$ ./build.sh arm
```
or

build for amd64<br>
If it is done, Cotainer would be created<br>
```
$ ./build.sh
```

#### build clean

```
$ make clean
```

#### Proxy
If you use proxy environment, you shoud add up the information in Dockerfile<br>

```
FROM scratch
ADD agent /
ADD ./data/server_url.json /data/
ENV http_proxy=http://10.112.1.184:8080
CMD ["/agent"]
```

If you run 'build.sh' or 'build.sh arm' it would be created container with the proxy environment.




