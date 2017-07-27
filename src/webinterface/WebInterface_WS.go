package webinterface

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	dockzen_h "include"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	set "services"
	"strings"
	"syscall"
	"time"
)

var wss_prefix = "ws://"

// Command structure contains command information.
type Command struct {
	Cmd string `json:"cmd"`
}

// SendChannel structure contains send channel information.
type SendChannel struct {
	containers chan ws_ContainerList_info
	updateinfo chan ws_ContainerUpdateReturn
}

// ReceiveChannel structure contains receive channel information.
type ReceiveChannel struct {
	containers chan bool
	updateinfo chan dockzen_h.ContainerUpdateInfo
}

// Containers_Channel structure contains channel information.
type Containers_Channel struct {
	receive chan bool
	send    chan ws_ContainerList_info
}

// Containers_Channel structure contains channel information for update command.
type Update_Channel struct {
	receive chan dockzen_h.ContainerUpdateInfo
	send    chan ws_ContainerUpdateReturn
}

var chSignal chan os.Signal
var done chan bool
var messagesCh chan string
var g_ws *websocket.Conn

// WI_init start ws_mainloop function.
func WI_init() {

	log.Printf("[%s] Web connection start !!!\n", __FILE__)

	for {

		ws_mainLoop()

		time.Sleep(time.Second)
	}

}

// Static ws_Server_Connect connect web server.
// Server_url param is address of web server.
// websocket.conn param is uniq id of web socket.
// This function return result of function.
func ws_Server_Connect(server_url string) (ws *websocket.Conn, err error) {

	var wss_server_url = wss_prefix + server_url
	ws, err = wsProxyDial(wss_server_url, "tcp", wss_server_url)

	if err != nil {
		log.Printf("[%s] wsProxyDial : ", __FILE__, err)
		//syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	name := getUniqueID()

	err = wsReqeustConnection(ws, name)
	if err != nil {
		log.Printf("[%s] ws_Server_Connect error = ", err)
		return ws, err
	}

	return ws, nil

}

// Static ws_MessageLoop handles incomming message from the web server.
// Param consists of message channel and receive channel.
func ws_MessageLoop(receive_channel ReceiveChannel) {

	for {
		// global channel
		msg := <-messagesCh
		log.Printf("[%s] MESSAGE !!! ", __FILE__)
		rcv := Command{}
		json.Unmarshal([]byte(msg), &rcv)
		log.Printf(rcv.Cmd)

		switch rcv.Cmd {
		case "connected":
			log.Printf("[%s] connected succefully~~", __FILE__)
		case "GetContainersInfo":
			receive_channel.containers <- true
		case "UpdateImage":
			log.Printf("[%s] command <UpdateImage>", __FILE__)
			update_msg, r := parseUpdateParam(msg)
			if r == nil {
				receive_channel.updateinfo <- update_msg
			} else {
				log.Printf("[%s] UpdateImage message null !!!")
			}

		default:
			log.Printf("[%s] add command of {%s}", __FILE__, rcv.Cmd)
			break;
		}

	}
}

// Static ws_mainLoop is main loop.
func ws_mainLoop() (err error) {

	go func() {
		<-chSignal
		done <- true
		return
	}()

	var server_url = set.GetServerURL()
	if server_url == "" {
		log.Printf("[%s] Server URL Error !! ", __FILE__)
		return
	}
	// global web socket to user in different go routine.
	g_ws, err = ws_Server_Connect(server_url)
	if err == nil {
		log.Printf("[%s] ws_Server_Connect done successfully~~~ ", __FILE__)
		messagesCh = make(chan string)
		go wsReceive(server_url, g_ws)

		var send_channel SendChannel
		send_channel.containers = make(chan ws_ContainerList_info, 5)
		send_channel.updateinfo = make(chan ws_ContainerUpdateReturn, 5)

		ws_SendMsg(send_channel)

		var container_ch Containers_Channel
		container_ch.receive = make(chan bool)
		container_ch.send = send_channel.containers

		var update_ch Update_Channel
		update_ch.receive = make(chan dockzen_h.ContainerUpdateInfo)
		update_ch.send = send_channel.updateinfo

		for i := 0; i < 3; i++ {
			go ws_GetContainerLists(container_ch)
			go ws_UpdateImage(update_ch)
		}

		var receive_channel ReceiveChannel
		receive_channel.containers = container_ch.receive
		receive_channel.updateinfo = update_ch.receive
		//go ws_UpdateImage(update_msg, send_channel.updateinfo)

		defer g_ws.Close()
		ws_MessageLoop(receive_channel)
	}

	log.Printf("[%s] return ", __FILE__)
	return err
}

// Static ws_SendMsg sends message to web server.
// Ws param is uniq id of web socket.
// send_channel param is send channel.
func ws_SendMsg(send_channel SendChannel) {
	for {
		select {
		case send_msg := <-send_channel.containers:
			log.Printf("[%s] containers sendMessage= ", __FILE__, send_msg)
			websocket.JSON.Send(g_ws, send_msg)
		case send_msg := <-send_channel.updateinfo:
			log.Printf("[%s] update sendMessage=", __FILE__, send_msg)
		}
	}
}

// Static wsReceive receives message from web server.
// Param consists of web server url, uniq id of web socket.
func wsReceive(server_url string, ws *websocket.Conn) (err error) {

	var read_buf string
/*
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		for {
			log.Printf("[%s] panic recovery !!!", __FILE__)
			g_ws, err = ws_Server_Connect(server_url)
			if err != nil {
				log.Printf("[%s] wsProxyDial : %s ", __FILE__, err)
				time.Sleep(time.Second)
				continue
			}
			go wsReceive(server_url, g_ws)
			break
		}
	}()
*/
	for {
		err = websocket.Message.Receive(ws, &read_buf)
		if err != nil {
			log.Printf("[%s] wsReceive : %s", __FILE__, err)
			messagesCh <- "err"
			break
		}
		log.Printf("[%s] received: %s", __FILE__, read_buf)
		messagesCh <- read_buf
	}

	return err
}

// Static wsReqeustConnection send device information to web server.
// Param consists of unique id of web socket and device id.
func wsReqeustConnection(ws *websocket.Conn, name string) (err error) {
	send := ConnectReq{}
	send.Cmd = "request"
	send.Name = name

	websocket.JSON.Send(ws, send)

	return nil
}

// Static wsProxyDial opens a new client connection to a websocket.
// Param consists of server url, protocol and original server url.
// This function returns unique id of web socket and result of function.
func wsProxyDial(url_, protocol, origin string) (ws *websocket.Conn, err error) {

	log.Printf("[%s] http_proxy {%s}\n", __FILE__, os.Getenv("http_proxy"))

	// comment out in case of testing without proxy(internal server)
	if strings.Contains(url_, "10.113.") {
		return websocket.Dial(url_, protocol, origin)
	}

	if os.Getenv("http_proxy") == "" {
		return websocket.Dial(url_, protocol, origin)
	}

	purl, err := url.Parse(os.Getenv("http_proxy"))
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
		log.Printf("[%s] NewConfig : ", __FILE__, err)
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
		//syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewClient")
	log.Printf("====================================")

	ret_ws, err := websocket.NewClient(config, client);
	if err != nil {
		log.Printf("[%s] NewClient ERR : ", __FILE__, err)
	}
	return ret_ws, err
}

// Static wsHttpConnect connect to web server.
// Param consists of host name, web server address.
// This function returns ReadWriteCloser and result of function.
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
		//syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}
	log.Printf("[%s] proxy_http_conn<resp> =", __FILE__, (resp))

	rwc, _ := proxy_http_conn.Hijack()

	return rwc, nil

}
