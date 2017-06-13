package services

import (
	"bytes"
	"agent/types/dockzenl"
	"lib"
	"encoding/json"
	"errors"
	"log"
	"net"
)

const (
	DockerLauncherSocket string = "/var/run/dockzen_launcher.sock"
)

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

func DZA_Mon_GetContainersInfo() {//([]byte, error) {
	log.Printf("GetContainersInfo")

	//apitest.GetContainerListsInfo_Test()
	//var info lib.ContainerListsInfo

	//lib.GetContainerListsInfo()
	lib.GetContainerListsInfo()
	/*
		stub := getDockerLauncherInfo_Stub()
		var send_stub []byte

		send_stub, _ = json.Marshal(stub)
		log.Printf(string(send_stub))

		return send_stub, nil
	*/
	//client, err := net.Dial("unix", DockerLauncherSocket)
	//if err != nil {
	//	log.Fatal("Dial error", err)
	//	return nil, err
	//}

	//defer client.Close()

	// Send Command to dockerl
	//err = writeData(client, "GetContainersInfo", nil)
	//if err != nil {
	//	return nil, err
	//}

	// Receive Command from dockerl
	//data := make([]byte, 0)
	//data, err = readData(client)
	//if err != nil {
	//	return nil, err
	//}

	//log.Printf("end\n")
	//return data, nil
}
