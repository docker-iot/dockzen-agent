package services

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
)

func parseUpdateImageParam(request *http.Request) (ImageName, ContainerName string, err error) {

	var body UpdateImageParams

	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&body)

	log.Printf("body.ImageName = %s\n", body.ImageName)
	log.Printf("body.ContainerName = %s\n", body.ContainerName)

	ImageName = body.ImageName
	ContainerName = body.ContainerName

	return ImageName, ContainerName, err
}

func DZA_Update_Do(request *http.Request) ([]byte, error) {
	log.Printf("UpdateImageRequest")
	/*
		stub := updateImage_Stub()
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
	ImageName, ContainerName, err := parseUpdateImageParam(request)
	if err != nil {
		return nil, errors.New("Invalid Parameter")
	}
	log.Printf("ImageName : %s\n", ImageName)
	log.Printf("ContainerName : %s\n", ContainerName)
	m := make(map[string]string)
	m["ImageName"] = ImageName
	m["ContainerName"] = ContainerName

	err = writeData(client, "UpdateImage", m)
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
