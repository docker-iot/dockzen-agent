package services

import (
	"bytes"
	"agent/types/dockzenl"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const defaultTimeout = 30 * time.Second

type CSAClient struct {
	Path       string
	HTTPClient *http.Client
}

func newHTTPClient(path string, timeout time.Duration) *http.Client {
	httpTransport := &http.Transport{}

	socketPath := path
	unixDial := func(proto, addr string) (net.Conn, error) {
		return net.DialTimeout("unix", socketPath, timeout)
	}
	httpTransport.Dial = unixDial

	return &http.Client{Transport: httpTransport}
}

func NewCSAClient() (*CSAClient, error) {

	httpClient := newHTTPClient(ContainerServiceSocket, time.Duration(defaultTimeout))
	return &CSAClient{ContainerServiceSocket, httpClient}, nil
}

func (client *CSAClient) doRequest(method string, path string, body string) ([]byte, error) {
	log.Printf("doRequest Method[%s] path[%s]", method, path)

	var resp *http.Response
	var err error

	switch method {
	case "GET":
		resp, err = client.HTTPClient.Get("http://unix" + path)
	case "POST":
		reqBody := bytes.NewBufferString(body)
		log.Printf("reqBody : [%s]\n", reqBody)
		resp, err = client.HTTPClient.Post("http://unix"+path, "text/plain", reqBody)
	default:
		return nil, errors.New("Invaild Method")
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		if method == "GET" {
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Printf("error =%s\n", err)
			}
			return contents, err
		}

		return nil, err

	} else {
		log.Printf("Error  : [%d]\n", resp.StatusCode)
	}

	return nil, errors.New(string(resp.StatusCode))
}

func GetHardwareAddress() (string, error) {

	currentNetworkHardwareName := "eth0"
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
	}

	name := netInterface.Name
	macAddress := netInterface.HardwareAddr

	log.Printf("Hardware name : %s\n", string(name))

	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		log.Printf("No able to parse MAC address : %s\n", err)
		os.Exit(-1)
	}

	log.Printf("Physical hardware address : %s \n", hwAddr.String())

	return hwAddr.String(), nil
}

func (client *CSAClient) GetContainersInfo() (ContainerLists, error) {

	var send ContainerLists

	contents, err := client.doRequest("GET", "/v1/get/GetContainersInfo", "")

	if err != nil {
		log.Printf("error [%s]", err)
		return send, err
	}

	lists := dockzenl.GetContainersInfoReturn{}

	json.Unmarshal([]byte(contents), &lists)
	var numOfList int = len(lists.Containers)
	log.Printf("numOfList[%d]\n", numOfList)

	send.Cmd = "GetContainersInfo"

	macaddress, err := GetHardwareAddress()

	send.DeviceID = macaddress
	log.Printf("send.DeviceID[%s]\n", send.DeviceID)
	var containerValue ContainerInfo

	for i := 0; i < numOfList; i++ {
		if lists.Containers[i].ContainerName != "" && lists.Containers[i].ImageName != "" {
			containerValue = ContainerInfo{
				ContainerName:   lists.Containers[i].ContainerName,
				ImageName:       lists.Containers[i].ImageName,
				ContainerStatus: lists.Containers[i].ContainerStatus,
			}
			send.ContainerCount++
		}

		send.Container = append(send.Container, containerValue)
		log.Printf("[%d]-[%s]", i, send.Container)
	}
	log.Printf("Container Count [%d]\n", send.ContainerCount)
	log.Printf("[%s]\n", send)

	return send, nil
}

func (client *CSAClient) UpdateImage(data UpdateImageParams) (UpdateImageReturn, error) {
	var send UpdateImageReturn

	send_str, _ := json.Marshal(data)
	fmt.Println(string(send_str))

	_, err := client.doRequest("POST", "/v1/post/UpdateImage", string(send_str))

	if err != nil {
		log.Printf("error [%s]", err)
		return send, err
	}

	send.Cmd = "UpdateImage"

	macaddress, err := GetHardwareAddress()

	log.Printf("macaddress[%s]\n", macaddress)
	send.DeviceID = macaddress

	return send, nil
}
