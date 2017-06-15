package webinterface

import (
	dockzen_h "include"
	"fmt"
	"log"
	"services"
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
)

var wss_server_url = "ws://10.113.62.204:4000"
var wss_server_origin = "ws://10.113.62.204:4000"

//var wss_server_url = "ws://13.124.64.10:4000"
//var wss_server_origin = "ws://13.124.64.10:4000"

type Command struct {
	Cmd string `json:"cmd"`
}

var chSignal chan os.Signal
var done chan bool

func WI_init(){

	log.Println("Web connection start !!!\n")

	for {

		go ws_mainLoop()

		<-done
		time.Sleep(time.Second)
	}

}

func wsGetContainerLists(ws *websocket.Conn) (err error) {
	var containersInfo dockzen_h.Containers_info
	var ret = services.DZA_Mon_GetContainersInfo(&containersInfo)

	if ret != 0 {
		fmt.Println("GetContainersInfo error = ", ret)
	} else {
		var send_info ContainerList_info
		send_info.Cmd = "GetContainersInfo"
		send_info.ContainerCount = int(containersInfo.Count)
		send_info.DeviceID, err = services.GetHardwareAddress()

		fmt.Println("DevicedID = ", send_info.DeviceID)

		for i := 0; i < send_info.ContainerCount; i++ {
			send_info.Container[i].ID = containersInfo.Containerinfo[i].ID
			send_info.Container[i].Name = containersInfo.Containerinfo[i].Name
			send_info.Container[i].ImageName = containersInfo.Containerinfo[i].ImageName
			send_info.Container[i].Status = containersInfo.Containerinfo[i].Status
		}
		fmt.Println("ContainerInfo -> ", send_info)
		websocket.JSON.Send(ws, send_info)
	}

	return nil
}

func wsUpdateImage(ws *websocket.Conn, data dockzen_h.ContainerUpdateInfo) (err error) {

	var ret = services.DZA_Update_Do(data)

	fmt.Println("ret = ", ret )
//	if err1 != nil {
//		log.Printf("error = %s", err1)
//		return err1
//	} else {
//		log.Printf("send = %s", send)
//		websocket.JSON.Send(ws, send)
//	}

	return nil
}

func parseUpdateParam(msg string) dockzen_h.ContainerUpdateInfo {
	send := dockzen_h.ContainerUpdateInfo{}
	json.Unmarshal([]byte(msg), &send)
	fmt.Println("parsing ContainerName: " + send.ContainerName)
	fmt.Println("parsing ImageName: " + send.ImageName)

	return send
}

func ws_mainLoop() (err error) {

	go func() {
		<-chSignal
		done <- true
		return
	}()

	ws, err := wsProxyDial(wss_server_url, "tcp", wss_server_origin)

	if err != nil {
		log.Println("wsProxyDial : ", err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return err
	}

	defer ws.Close()

	/* connect test2 : message driven
	 */
	messages := make(chan string)
	go wsReceive(ws, messages)

	name, _ := services.GetHardwareAddress()

	err = wsReqeustConnection(ws, name)

	for {
		msg := <-messages

		rcv := Command{}
		json.Unmarshal([]byte(msg), &rcv)
		fmt.Println(rcv.Cmd)

		switch rcv.Cmd {
		case "connected":
			log.Printf("connected succefully~~")
		case "GetContainersInfo":
			wsGetContainerLists(ws)
		case "UpdateImage":
			log.Printf("command <UpdateImage>")
			wsUpdateImage(ws, parseUpdateParam(msg))
		default:
			log.Printf("add command of {%s}", rcv.Cmd)
		}

	}
}

func wsReceive(ws *websocket.Conn, chan_msg chan string) (err error) {

	var read_buf string

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		for {
			log.Printf("panic recovery !!!")
			ws, err = wsProxyDial(wss_server_url, "tcp", wss_server_origin)

			if err != nil {
				log.Printf("wsProxyDial : %s ", err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}()

	for {
		err = websocket.Message.Receive(ws, &read_buf)
		if err != nil {
			log.Printf("wsReceive : %s", err)
			syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
			break
		}
		log.Printf("received: %s", read_buf)
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

	log.Printf("http_proxy {%s}\n", os.Getenv("HTTP_PROXY"))

	// comment out in case of testing without proxy
	if strings.Contains(url_, "10.113.") {
		return websocket.Dial(url_, protocol, origin)
	}

	if os.Getenv("HTTP_PROXY") == "" {
		return websocket.Dial(url_, protocol, origin)
	}

	purl, err := url.Parse(os.Getenv("HTTP_PROXY"))
	if err != nil {
		log.Println("Parse : ", err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewConfig")
	log.Printf("====================================")
	config, err := websocket.NewConfig(url_, origin)
	if err != nil {
		log.Println("NewConfig : ", err)
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
		log.Println("HttpConnect : ", err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewClient")
	log.Printf("====================================")
	return websocket.NewClient(config, client)
}

func wsHttpConnect(proxy, url_ string) (io.ReadWriteCloser, error) {
	log.Println("proxy =", proxy)
	proxy_tcp_conn, err := net.Dial("tcp", proxy)
	if err != nil {
		return nil, err
	}
	log.Println("proxy_tcp_conn =", proxy_tcp_conn)
	log.Println("url_ =", url_)

	turl, err := url.Parse(url_)
	if err != nil {
		log.Println("Parse : ", err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Println("proxy turl.Host =", string(turl.Host))

	req := http.Request{
		Method: "CONNECT",
		URL:    &url.URL{},
		Host:   turl.Host,
	}

	proxy_http_conn := httputil.NewProxyClientConn(proxy_tcp_conn, nil)
	//cc := http.NewClientConn(proxy_tcp_conn, nil)

	log.Println("proxy_http_conn =", proxy_http_conn)

	resp, err := proxy_http_conn.Do(&req)
	if err != nil && err != httputil.ErrPersistEOF {
		log.Println("ErrPersistEOF : ", err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}
	log.Println("proxy_http_conn<resp> =", (resp))

	rwc, _ := proxy_http_conn.Hijack()

	return rwc, nil

}
