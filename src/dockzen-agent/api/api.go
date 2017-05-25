package api

import (
	"dockzen-agent/api/types"
	"dockzen-agent/types/dockzenl"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/schorsch/go-callbacks"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RequestLauncher struct {
	WS   *websocket.Conn
	Data string
}

type ResponseLauncher struct {
	WS   *websocket.Conn
	Data string
}

var dockzen_agent_origin = "http://localhost:8080/wsAgent"
var dockzen_agent_url = "ws://localhost:8080/wsAgent"

var dockzen_agent_notify_origin = "http://localhost:8082/wsNotify"
var dockzen_agent_notify_url = "ws://localhost:8082/wsNotify"

const defaultTimeout = 30 * time.Second

var doneServer chan bool
var chSignal chan os.Signal
var sendToLauncherCh chan RequestLauncher
var responseFromLauncher chan ResponseLauncher

var notifyCh chan string
var createWsAgentConn bool = false
var createWsNotifyConn bool = false
var wsAgent *websocket.Conn
var wsNotify *websocket.Conn

var cbs callbacks.Callbacks

func checkConnection(name string) bool {
	switch name {
	case "wsAgent":
		return createWsAgentConn
	case "wsNotify":
		return createWsNotifyConn
	}

	return false
}

func createConnection(name string) (*websocket.Conn, error) {
	var ws *websocket.Conn = nil
	var err error

	switch name {
	case "wsAgent":
		if checkConnection("wsAgent") == false {
			ws, err = websocket.Dial(dockzen_agent_url, "", dockzen_agent_origin)
			if err != nil {
				log.Fatal(err)
			}
			wsAgent = ws
			createWsAgentConn = true
		} else {
			ws = wsAgent
			fmt.Println("Already Connected")
		}
	case "wsNotify":
		if checkConnection("wsNotify") == false {
			ws, err = websocket.Dial(dockzen_agent_notify_url, "", dockzen_agent_notify_origin)
			if err != nil {
				log.Fatal(err)
			}
			wsNotify = ws
			createWsNotifyConn = true
		} else {
			ws = wsNotify
			fmt.Println("Already Connected")
		}
	default:
		fmt.Println("Unknown Connection")
		err = errors.New("Unknown Connection")
	}

	return ws, err
}

func triggerCallback(data string) {

	cbs.CallbacksCall("stateChanged", data)
}

func startAPIServer(webSocketAgent, webSocketNotify *websocket.Conn) {
	fmt.Println("StartAPIServer")

	chSignal = make(chan os.Signal, 1)
	//sendToLauncherCh = make(chan string)
	sendToLauncherCh = make(chan RequestLauncher)
	responseFromLauncher = make(chan ResponseLauncher)

	notifyCh = make(chan string)

	signal.Notify(chSignal, syscall.SIGUSR1)
	go receiver(webSocketAgent, webSocketNotify)
	go sender(sendToLauncherCh, notifyCh)

	defer webSocketAgent.Close()
	defer webSocketNotify.Close()

	for {
		time.Sleep(time.Second)
	}
}

func NewAgentHndl() (*types.DockzenHndl, error) {

	agentWS, err1 := createConnection("wsAgent")
	if err1 != nil {
		fmt.Println(err1)
		return nil, err1
	}

	notifyWS, err2 := createConnection("wsNotify")
	if err2 != nil {
		fmt.Println(err2)
		return nil, err2
	}

	go startAPIServer(agentWS, notifyWS)

	fmt.Println("return NewAgnetHndl")

	return &types.DockzenHndl{agentWS}, nil
}

func RegisterStateChangedCB(theFunc func()) error {

	var err error

	if checkConnection("wsNotify") == true {
		cb := callbacks.Callback{Name: "StateChanaged", Method: theFunc}
		cbs = append(cbs, cb)
	} else {
		err = errors.New("Unable to register callback")
	}

	return err
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

func wsReceive(ws *websocket.Conn, chan_msg chan string) (err error) {

	var read_buf string

	for {
		err = websocket.Message.Receive(ws, &read_buf)

		if err != nil {
			log.Printf("wsReceive2 wsReceive : %s", err)
			syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
			break
		}

		chan_msg <- read_buf
	}

	return err
}

func GetContainersInfo(hndl *types.DockzenHndl, data string) (types.ContainerLists, error) {

	var send types.ContainerLists
	var contents []byte
	var err error

	sendCmd := types.Cmd{
		Cmd: data,
	}

	send_json, _ := json.Marshal(sendCmd)
	//send to launcher
	var sendData RequestLauncher
	sendData = RequestLauncher{
		WS:   hndl.WS,
		Data: string(send_json),
	}

	sendToLauncherCh <- sendData

	// Need to receive return from dockzen-agent
	//contents =
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

	var containerValue types.ContainerInfo

	for i := 0; i < numOfList; i++ {
		if lists.Containers[i].ContainerName != "" && lists.Containers[i].ImageName != "" {
			containerValue = types.ContainerInfo{
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

func receiver(webSocketAgent, webSocketNotify *websocket.Conn) {
	fmt.Println("Api receiver Start")
	go func() {
		<-chSignal
		return
	}()

	recvApiResponseCh := make(chan string)
	recvNotifyCh := make(chan string)

	go wsReceive(webSocketAgent, recvApiResponseCh)
	go wsReceive(webSocketNotify, recvNotifyCh)

	//defer webSocketAgent.Close()
	//defer webSocketNotify.Close()
	for {
		select {
		case msg1 := <-recvApiResponseCh:
			fmt.Println("recvApiResponseCh :" + msg1)
			rcv := types.Cmd{}
			_ = rcv
		case msg2 := <-recvNotifyCh:
			fmt.Println("recvNotifyCh :" + msg2)
			triggerCallback(msg2)
		}
	}

	fmt.Println("end api receiver")
}

func sender(sendToLaunchCh chan RequestLauncher, callBackCh chan string) {
	for {
		select {
		case msg1 := <-sendToLaunchCh:
			fmt.Println("api Sender : send to Launcher")
			n, err := msg1.WS.Write([]byte(msg1.Data))
			if err != nil {
				fmt.Printf("error : %s\n", err)
			}
			fmt.Printf("Sender[%d]: %s\n", n, msg1.Data)
		case msg2 := <-callBackCh:
			fmt.Println("api callback : callback")
			cbs.CallbacksCall("stateChanged", msg2)
			fmt.Printf("Callback[%s]\n", msg2)
		}
	}
}
