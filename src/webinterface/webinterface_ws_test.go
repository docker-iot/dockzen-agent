package webinterface

type ConnectedResp struct {
	Cmd       string `json:"cmd"`
	Token     string `json:"token"`
	Clinetnum int    `json:"clientnum"`
}

func wsTest1(ws *websocket.Conn) (err error) {
	name, _ := os.Hostname()
	err = wsReqeustConnection(ws, name)

	// receive connection token
	Token, err := wsReceiveConnection(ws)
	log.Printf("recv.Token = '%s'", Token)

	return err
}

func wsReceiveConnection(ws *websocket.Conn) (Token string, err error) {
	recv := ConnectedResp{}

	err = websocket.JSON.Receive(ws, &recv)
	if err != nil {
		log.Println("wsReceiveConnection : ", err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return "", err
	}

	return recv.Token, err
}

func json_marshal() {
	// convert from struct to string
	send := ConnectedResp{}
	send.Cmd = "request"
	send.Token = "1234"
	send.Clinetnum = 88

	send_str, _ := json.Marshal(send)
	fmt.Println(string(send_str))
}

func json_unmarshal() {
	// convert from string to struct
	rcv_str := `{"cmd": "connected"
			, "token": "test-token"
			, "clinetnum": 3}`
	rcv := ConnectedResp{}
	json.Unmarshal([]byte(rcv_str), &rcv)
	fmt.Println(rcv)
	fmt.Println(rcv.Cmd)
}

type DockerInfo struct {
	ID      string        `json:"Id"`
	Names   []string      `json:"Names"`
	Image   string        `json:"Image"`
	ImageID string        `json:"ImageID"`
	Command string        `json:"Command"`
	Created int           `json:"Created"`
	Ports   []interface{} `json:"Ports"`
	Labels  struct {
	} `json:"Labels"`
	State      string `json:"State"`
	Status     string `json:"Status"`
	HostConfig struct {
		NetworkMode string `json:"NetworkMode"`
	} `json:"HostConfig"`
	NetworkSettings struct {
		Networks struct {
			Bridge struct {
				IPAMConfig          interface{} `json:"IPAMConfig"`
				Links               interface{} `json:"Links"`
				Aliases             interface{} `json:"Aliases"`
				NetworkID           string      `json:"NetworkID"`
				EndpointID          string      `json:"EndpointID"`
				Gateway             string      `json:"Gateway"`
				IPAddress           string      `json:"IPAddress"`
				IPPrefixLen         int         `json:"IPPrefixLen"`
				IPv6Gateway         string      `json:"IPv6Gateway"`
				GlobalIPv6Address   string      `json:"GlobalIPv6Address"`
				GlobalIPv6PrefixLen int         `json:"GlobalIPv6PrefixLen"`
				MacAddress          string      `json:"MacAddress"`
			} `json:"bridge"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
	Mounts []struct {
		Type        string `json:"Type"`
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
		Mode        string `json:"Mode"`
		RW          bool   `json:"RW"`
		Propagation string `json:"Propagation"`
	} `json:"Mounts"`
}

func dockertest() {
	// test docker daemon response
	inputstring := `[{"Id": "8433735be769c5787965fdbd3d8bd7635d793d5fff968b0626ea04f5ee80a755"
				,"Names": "/poc1"
				,"Image": "13.124.64.10:443/minimal"
				,"ImageID":"sha256:8502bca5fca7a2a8ea6e5434a1a5462cc4cf84c116cbdacef4aab078b2571dc8"
				,"Command":"/sbin/init"
				,"Created":1491832614
				,"Ports":[]
				,"Labels":{}
				,"State":"created"
				,"Status":"Created"
				,"HostConfig":{"NetworkMode":"bridge"}
				,"NetworkSettings":{"Networks":{"bridge":{"IPAMConfig":null
														,"Links":null
														,"Aliases":null
														,"NetworkID":"7fdceaa2ab5188435e6d9f553d54cd530a1b2b8396e7ccc66e66b55faab51223"
														,"EndpointID":"b89fd1e7ab9786b9080904b455b6e14e60648749f66cfac5d02c81dd59876df8"
														,"Gateway":"172.17.0.1"
														,"IPAddress":"172.17.0.2"
														,"IPPrefixLen":16
														,"IPv6Gateway":""
														,"GlobalIPv6Address":""
														,"GlobalIPv6PrefixLen":0
														,"MacAddress":"02:42:ac:11:00:02"}}}
				,"Mounts":[{"Type":"bind"
						,"Source":"/sys/fs/cgroup"
						,"Destination":"/sys/fs/cgroup"
						,"Mode":""
						,"RW":true
						,"Propagation":""}]
				},
				{"Id": "8433735be769c5787965fdbd3d8bd7635d793d5fff968b0626ea04f5ee80a755"
				,"Names": "/poc1"
				,"Image": "13.124.64.10:443/minimal"
				,"ImageID":"sha256:8502bca5fca7a2a8ea6e5434a1a5462cc4cf84c116cbdacef4aab078b2571dc8"
				,"Command":"/sbin/init"
				,"Created":1491832614
				,"Ports":[]
				,"Labels":{}
				,"State":"created"
				,"Status":"Created"
				,"HostConfig":{"NetworkMode":"bridge"}
				,"NetworkSettings":{"Networks":{"bridge":{"IPAMConfig":null
														,"Links":null
														,"Aliases":null
														,"NetworkID":"7fdceaa2ab5188435e6d9f553d54cd530a1b2b8396e7ccc66e66b55faab51223"
														,"EndpointID":"b89fd1e7ab9786b9080904b455b6e14e60648749f66cfac5d02c81dd59876df8"
														,"Gateway":"172.17.0.1"
														,"IPAddress":"172.17.0.2"
														,"IPPrefixLen":16
														,"IPv6Gateway":""
														,"GlobalIPv6Address":""
														,"GlobalIPv6PrefixLen":0
														,"MacAddress":"02:42:ac:11:00:02"}}}
				,"Mounts":[{"Type":"bind"
						,"Source":"/sys/fs/cgroup"
						,"Destination":"/sys/fs/cgroup"
						,"Mode":""
						,"RW":true
						,"Propagation":""}]}]`

	dockerinfo := make([]DockerInfo, 0)
	json.Unmarshal([]byte(inputstring), &dockerinfo)

	fmt.Printf("[0] id =%s\n", dockerinfo[0].ID)
	fmt.Printf("[0] Image =%s\n", dockerinfo[0].Image)
	fmt.Printf("[1] id =%s\n", dockerinfo[1].ID)
	fmt.Printf("[1] Image =%s\n", dockerinfo[1].Image)
}
