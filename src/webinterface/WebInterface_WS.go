package webinterface

import (
	"log"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"encoding/json"
	"syscall"
	"time"
	set "services"
	dockzen_h "include"
)

var wss_prefix = "ws://"

type Command struct {
	Cmd string `json:"cmd"`
}

type SendChannel struct{
	containers chan ws_ContainerList_info
	updateinfo chan ws_ContainerUpdateReturn
}

type ReceiveChannel struct{
	containers chan bool
	updateinfo chan dockzen_h.ContainerUpdateInfo
}

type Containers_Channel struct{
	receive chan bool
	send chan ws_ContainerList_info
}

type Update_Channel struct{
	receive chan dockzen_h.ContainerUpdateInfo
	send chan ws_ContainerUpdateReturn
}

var chSignal chan os.Signal
var done chan bool

func WI_init(){

	log.Printf("[%s] Web connection start !!!\n", __FILE__)

	for {

		go ws_mainLoop()

		<-done
		time.Sleep(time.Second)
	}

}

func WS_Server_Connect(server_url string) (ws *websocket.Conn, err error) {

	var wss_server_url = wss_prefix + server_url
	ws, err = wsProxyDial(wss_server_url, "tcp", wss_server_url)

	if err != nil {
		log.Printf("[%s] wsProxyDial : ",__FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	name, _ := GetHardwareAddress()

	err = wsReqeustConnection(ws, name)
	if err != nil {
		log.Printf("[%s] WS_Server_Connect error = ", err)
		return ws, err
	}

	return ws, nil

}

func WS_MessageLoop(messages chan string, receive_channel ReceiveChannel){

	for {
		msg := <-messages
		log.Printf("[%s] MESSAGE !!! ", __FILE__)
		rcv := Command{}
		json.Unmarshal([]byte(msg), &rcv)
		log.Printf(rcv.Cmd)

		switch rcv.Cmd {
		case "connected":
			log.Printf("[%s] connected succefully~~", __FILE__)
		case "GetContainersInfo":
				receive_channel.containers <-true
		case "UpdateImage":
			log.Printf("[%s] command <UpdateImage>", __FILE__)
			update_msg, r := ParseUpdateParam(msg)
			if r == nil {
					receive_channel.updateinfo <- update_msg
			} else {
				log.Printf("[%s] UpdateImage message null !!!")
			}

		default:
			log.Printf("[%s] add command of {%s}", __FILE__, rcv.Cmd)
		}

	}
}

func ws_mainLoop() (err error) {

	go func() {
		<-chSignal
		done <- true
		return
	}()

	var server_url = set.GetServerURL("")
	if server_url == "" {
		log.Printf("[%s] Server URL Error !! ", __FILE__)
		return
	}

	ws, err := WS_Server_Connect(server_url)

	messages := make(chan string)
	go wsReceive(server_url, ws, messages)

	var send_channel SendChannel
	send_channel.containers = make(chan ws_ContainerList_info)
	send_channel.updateinfo = make(chan ws_ContainerUpdateReturn)

	go WS_SendMsg(ws, send_channel)

	var container_ch Containers_Channel
	container_ch.receive = make(chan bool)
	container_ch.send = send_channel.containers

	var update_ch Update_Channel
	update_ch.receive = make(chan dockzen_h.ContainerUpdateInfo)
	update_ch.send = send_channel.updateinfo

	for i:= 0; i<3;i++{
		go WS_GetContainerLists(container_ch)
		go WS_UpdateImage(update_ch)
	}

	var receive_channel ReceiveChannel
	receive_channel.containers = container_ch.receive
	receive_channel.updateinfo = update_ch.receive
	//go WS_UpdateImage(update_msg, send_channel.updateinfo)

	defer ws.Close()
	WS_MessageLoop(messages, receive_channel)

	return nil
}

func WS_SendMsg(ws *websocket.Conn, send_channel SendChannel){
	for{
		select{
		case send_msg:= <-send_channel.containers:
			log.Printf("[%s] containers sendMessage= ", __FILE__, send_msg)
			websocket.JSON.Send(ws, send_msg)
		case send_msg:= <-send_channel.updateinfo:
			log.Printf("[%s] update sendMessage=", __FILE__, send_msg)
		}
	}
}

func wsReceive(server_url string, ws *websocket.Conn, chan_msg chan string) (err error) {

	var read_buf string

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		for {
			log.Printf("[%s] panic recovery !!!", __FILE__)
			ws, err = wsProxyDial(server_url, "tcp", server_url)

			if err != nil {
				log.Printf("[%s] wsProxyDial : %s ", __FILE__, err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}()

	for {
		err = websocket.Message.Receive(ws, &read_buf)
		if err != nil {
			log.Printf("[%s] wsReceive : %s", __FILE__, err)
			syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
			break
		}
		log.Printf("[%s] received: %s", __FILE__, read_buf)
		chan_msg <- read_buf
	}

	return err
}

func wsReqeustConnection(ws *websocket.Conn, name string) (err error) {
	send := ConnectReq{}
	send.Cmd = "request"
	send.Name = name

	websocket.JSON.Send(ws, send)

	return nil
}


func wsProxyDial(url_, protocol, origin string) (ws *websocket.Conn, err error) {

	log.Printf("[%s] http_proxy {%s}\n", __FILE__, os.Getenv("HTTP_PROXY"))

	// comment out in case of testing without proxy
	if strings.Contains(url_, "10.113.") {
		return websocket.Dial(url_, protocol, origin)
	}

	if os.Getenv("HTTP_PROXY") == "" {
		return websocket.Dial(url_, protocol, origin)
	}

	purl, err := url.Parse(os.Getenv("HTTP_PROXY"))
	if err != nil {
		log.Printf("[%s] Parse : ", __FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewConfig")
	log.Printf("====================================")
	config, err := websocket.NewConfig(url_, origin)
	if err != nil {
		log.Printf("[%s] NewConfig : ", __FILE__,  err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	if protocol != "" {
		config.Protocol = []string{protocol}
	}

	log.Printf("====================================")
	log.Printf("    HttpConnect")
	log.Printf("====================================")
	client, err := wsHttpConnect(purl.Host, url_)
	if err != nil {
		log.Printf("[%s] HttpConnect : ", __FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewClient")
	log.Printf("====================================")
	return websocket.NewClient(config, client)
}

func wsHttpConnect(proxy, url_ string) (io.ReadWriteCloser, error) {
	log.Printf("[%s] proxy =", __FILE__, proxy)
	proxy_tcp_conn, err := net.Dial("tcp", proxy)
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] proxy_tcp_conn =", __FILE__, proxy_tcp_conn)
	log.Printf("[%s] url_ =", __FILE__, url_)

	turl, err := url.Parse(url_)
	if err != nil {
		log.Printf("[%s] Parse : ", __FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("[%s] proxy turl.Host =", __FILE__, string(turl.Host))

	req := http.Request{
		Method: "CONNECT",
		URL:    &url.URL{},
		Host:   turl.Host,
	}

	proxy_http_conn := httputil.NewProxyClientConn(proxy_tcp_conn, nil)
	//cc := http.NewClientConn(proxy_tcp_conn, nil)

	log.Printf("[%s] proxy_http_conn =", __FILE__, proxy_http_conn)

	resp, err := proxy_http_conn.Do(&req)
	if err != nil && err != httputil.ErrPersistEOF {
		log.Printf("[%s] ErrPersistEOF : ", __FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}
	log.Printf("[%s] proxy_http_conn<resp> =", __FILE__, (resp))

	rwc, _ := proxy_http_conn.Hijack()

	return rwc, nil

}
