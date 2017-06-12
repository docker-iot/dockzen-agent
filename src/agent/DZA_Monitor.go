package main

import (
	"log"
	"net"
)

func DZA_Mon_GetContainersInfo() ([]byte, error) {
	log.Printf("GetContainersInfo")
	/*
		stub := getDockerLauncherInfo_Stub()
		var send_stub []byte

		send_stub, _ = json.Marshal(stub)
		log.Printf(string(send_stub))

		return send_stub, nil
	*/
	client, err := net.Dial("unix", DockerLauncherSocket)
	if err != nil {
		log.Fatal("Dial error", err)
		return nil, err
	}

	defer client.Close()

	// Send Command to dockerl
	err = writeData(client, "GetContainersInfo", nil)
	if err != nil {
		return nil, err
	}

	// Receive Command from dockerl
	data := make([]byte, 0)
	data, err = readData(client)
	if err != nil {
		return nil, err
	}

	log.Printf("end\n")
	return data, nil
}
