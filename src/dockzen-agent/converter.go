package main

import (
	"bytes"
	"dockzen-agent/types/dockzenl"
	"encoding/json"
	"fmt"
)

var recvDone chan bool

func getContainersInfo_Stub() dockzenl.GetContainersInfoReturn {
	send := dockzenl.GetContainersInfoReturn{
		Containers: []dockzenl.Container{
			{
				ContainerName:   "aaaa",
				ImageName:       "tizen1",
				ContainerStatus: "created",
			},
			{
				ContainerName:   "bbbb",
				ImageName:       "tizen2",
				ContainerStatus: "exited",
			},
		},
	}

	return send
}

func updateImage_Stub() dockzenl.Cmd {
	send := dockzenl.Cmd{
		Cmd: "UpdateImage",
	}

	return send
}

func updateImage(ch chan string, data []byte) ([]byte, error) {
	fmt.Println("UpdateImage")
	stub := updateImage_Stub()
	var send_stub []byte

	send_stub, _ = json.Marshal(stub)
	fmt.Println(string(send_stub))

	return send_stub, nil
}

func parseUpdateImageParam(data []byte) (ImageName, ContainerName string, err error) {

	/*decoder := json.NewDecoder(request.Body)
	decoder.Decode(&body)

	fmt.Println("body.ImageName = " + body.ImageName)
	lfmt.Println("body.ContainerName = " + body.ContainerName)

	ImageName = body.ImageName
	ContainerName = body.ContainerName
	*/

	return ImageName, ContainerName, err
}

func getContainersInfo(ch chan string) ([]byte, error) {
	fmt.Println("GetContainersInfo")

	stub := getContainersInfo_Stub()
	var send_stub []byte

	send_stub, _ = json.Marshal(stub)
	fmt.Println(string(send_stub))

	return send_stub, nil

	send := dockzenl.Cmd{
		Cmd: "GetContainersInfo",
	}

	send_json, _ := json.Marshal(send)

	ch <- string(send_json)
	fmt.Println("Send getContainersInfo command to launcher")

	//ret := <-recvDone
	//fmt.Println(ret)

	data := make([]byte, 0)

	for {
		dataBuf := make([]byte, 1024)
		nr, err := dockzenLauncherClient.Read(dataBuf)
		if err != nil {
			break
		}

		fmt.Printf("nr size [%d]\n", nr)
		if nr == 0 {
			break
		}

		dataBuf = dataBuf[:nr]
		data = append(data, dataBuf...)
	}

	fmt.Printf("receive data[%s]\n", string(data))
	withoutNull := bytes.Trim(data, "\x00")

	fmt.Println("end")
	return withoutNull, nil
}
