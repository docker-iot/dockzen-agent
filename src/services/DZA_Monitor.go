package services

import (
	"bytes"
	"agent/types/dockzenl"
	dockzen_api "lib"
	"encoding/json"
	"errors"
	"log"
	"net"
)

const (
	DockerLauncherSocket string = "/var/run/dockzen_launcher.sock"
)
func DZA_Mon_GetContainersInfo() (dockzen_api.ContainerLists, error) {
	log.Printf("GetContainersInfo")

	var info dockzen_api.ContainerLists
	var err error
	var send_str []byte

	info, err = dockzen_api.GetContainerListsInfo()

	//send_str, err = json.Marshal(info)

	log.Printf(string(send_str))

	return info, err
}

func readData(client net.Conn) ([]byte, error) {

	data := make([]byte, 0)

	for {
		dataBuf := make([]byte, 1024)
		nr, err := client.Read(dataBuf)
		if err != nil {
			break
		}

		log.Printf("nr size [%d]", nr)
		if nr == 0 {
			break
		}

		dataBuf = dataBuf[:nr]
		data = append(data, dataBuf...)
	}

	log.Printf("receive data[%s]\n", string(data))
	//delete null character
	withoutNull := bytes.Trim(data, "\x00")

	rcv := dockzenl.Cmd{}
	err := json.Unmarshal([]byte(withoutNull), &rcv)
	log.Printf("rcv.Cmd = %s", rcv.Cmd)

	if rcv.Cmd == "GetContainersInfo" {
		log.Printf("Success\n")
		return withoutNull, nil
	} else if rcv.Cmd == "UpdateImage" {
		log.Printf("Success\n")
		return withoutNull, nil
	} else {
		log.Printf("error commnad[%s]\n", err)
	}

	return nil, errors.New("Error Cmd from Dockerl")
}

func writeData(client net.Conn, cmd string, m map[string]string) error {
	var send_str []byte
	var err error

	if cmd == "GetContainersInfo" {
		send := dockzenl.Cmd{}
		send.Cmd = "GetContainersInfo"
		send_str, err = json.Marshal(send)
	} else if cmd == "UpdateImage" {

		send := dockzenl.UpdateImageParameters{}
		send.Cmd = "UpdateImage"

		send.Param = dockzenl.UpdateParam{
			ContainerName: m["ContainerName"],
			ImageName:     m["ImageName"],
		}
		send_str, err = json.Marshal(send)

	} else {
		return errors.New("Invalid Cmd")
	}

	log.Printf(string(send_str))
	length := len(send_str)

	message := make([]byte, 0, length)
	message = append(message, send_str...)

	_, err = client.Write([]byte(message))
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}

	log.Printf("sent: %s\n", message)
	err = client.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Printf("error: %v\n", err)
		return err

	}

	return nil
}
