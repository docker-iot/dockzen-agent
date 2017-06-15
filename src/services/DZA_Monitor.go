package services

import (
	"bytes"
	"agent/types/dockzenl"
	dockzen_api "lib"
	dockzen_h "include"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	DockerLauncherSocket string = "/var/run/dockzen_launcher.sock"
)

func DZA_Mon_GetContainersInfo(containersInfo *dockzen_h.Containers_info) int {

	fmt.Println("DZA_Mon_GetContainersInfo call !! ")

	var ret = dockzen_api.GetContainerListsInfo(containersInfo)

	return ret
}

func GetHardwareAddress() (string, error) {

	currentNetworkHardwareName := "eth0"
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
	}

	name := netInterface.Name
	macAddress := netInterface.HardwareAddr

	fmt.Println("Hardware name : ", string(name))

	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		fmt.Println("No able to parse MAC address : ", err)
		os.Exit(-1)
	}

	fmt.Println("Physical hardware address : ", hwAddr.String())

	return hwAddr.String(), nil
}

func readData(client net.Conn) ([]byte, error) {

	data := make([]byte, 0)

	for {
		dataBuf := make([]byte, 1024)
		nr, err := client.Read(dataBuf)
		if err != nil {
			break
		}

		fmt.Println("nr size ", nr)
		if nr == 0 {
			break
		}

		dataBuf = dataBuf[:nr]
		data = append(data, dataBuf...)
	}

	fmt.Println("receive data : ", string(data))
	//delete null character
	withoutNull := bytes.Trim(data, "\x00")

	rcv := dockzenl.Cmd{}
	err := json.Unmarshal([]byte(withoutNull), &rcv)
	fmt.Println("rcv.Cmd = ", rcv.Cmd)

	if rcv.Cmd == "GetContainersInfo" {
		fmt.Println("Success")
		return withoutNull, nil
	} else if rcv.Cmd == "UpdateImage" {
		fmt.Println("Success")
		return withoutNull, nil
	} else {
		fmt.Println("error commnad = ", err)
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

	//log.Printf(string(send_str))
	length := len(send_str)

	message := make([]byte, 0, length)
	message = append(message, send_str...)

	_, err = client.Write([]byte(message))
	if err != nil {
		//log.Printf("error: %v\n", err)
		return err
	}

	//log.Printf("sent: %s\n", message)
	err = client.(*net.UnixConn).CloseWrite()
	if err != nil {
		//log.Printf("error: %v\n", err)
		return err

	}

	return nil
}
